package raft

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/infrastructure"
	"sync"
	"sync/atomic"
)

const (
	Set    = 1
	Delete = 2
)

type RaftLog struct {
	Term  int64
	Index int64
	Cmd   Command
}

type Command struct {
	Operation int8
	Key       string
	Value     string
	Done      chan bool
}

type AppendEntriesRequest struct {
	LeaderId     int8
	Term         int64
	PrevLogTerm  int64
	PrevLogIndex int64
	LeaderCommit int64
	Entries      []RaftLog
	Done         chan *AppendEntriesResponse
}

type AppendEntriesResponse struct {
	Term    int64
	Success bool
	Err     error
}

type HeartBeat struct{}

func (raft *Raft) HandleClientOpeartion(op *Command) bool {
	cfg := infrastructure.CfgInstance
	log := &RaftLog{}
	log.Term = raft.CurrentTerm
	log.Index = raft.LastLog().Index + 1
	log.Cmd = *op

	raft.Log = append(raft.Log, *log)
	err := raft.SaveLog(log)

	if err != nil {
		return false
	}
	raft.sendLogEntries()

	replicateCount := 1

	raft.MatchIndex.Range(func(key, value interface{}) bool {
		idx := value.(int64)
		if idx >= log.Index {
			replicateCount++
		}
		return true
	})

	if replicateCount > len(cfg.ClusterMembers)/2 {
		atomic.StoreInt64(&raft.CommitIndex, log.Index)
		raft.CommitCh <- struct{}{}
		return true
	}
	return false
}

func (raft *Raft) HeartBeat() {
	//fmt.Println("发送心跳:" + strconv.FormatInt(int64(raft.CurrentTerm), 10))
	raft.sendLogEntries()
}

func (raft *Raft) sendLogEntries() {

	cfg := infrastructure.CfgInstance
	var wg sync.WaitGroup
	wg.Add(len(cfg.ClusterMembers))

	for _, node := range cfg.ClusterMembers {
		go func(node infrastructure.ClusterMember) {
			request := &rpcProto.AppendEntriesReq{}
			request.Term = int64(raft.CurrentTerm)
			request.LeaderId = int32(raft.NodeId)
			request.LeaderCommit = int64(raft.CommitIndex)
			request.Entries = make([]*rpcProto.LogEntry, 0)
			v, _ := raft.NextIndex.Load(node.NodeId)
			nextIdx := v.(int64)
			request.PrevLogIndex = nextIdx - 1
			request.PrevLogTerm = int64(raft.Log[request.PrevLogIndex].Term)
			lastIndex := nextIdx
			if raft.LastLog().Index >= nextIdx {
				sendLogs := raft.Log[nextIdx:]
				for _, log := range sendLogs {
					entry := &rpcProto.LogEntry{}
					entry.Index = log.Index
					entry.Term = log.Term
					entry.Key = log.Cmd.Key
					entry.Operation = int64(log.Cmd.Operation)
					entry.Value = log.Cmd.Value
					request.Entries = append(request.Entries, entry)
					lastIndex = log.Index
				}
			}

			res, err := infrastructure.AppendEntries(request, node.RpcIP, node.RpcPort)
			if err != nil {
				//fmt.Println("复制日志:" + err.Error())
			} else if res.Success && len(request.Entries) > 0 {
				raft.MatchIndex.Store(node.NodeId, lastIndex)
				raft.NextIndex.Store(node.NodeId, lastIndex+1)
			} else if !res.Success {
				raft.NextIndex.Store(node.NodeId, nextIdx-1)
			}
			wg.Done()
		}(node)
	}
	wg.Wait()
}

func (raft *Raft) HandleAppendEntries(req *AppendEntriesRequest) *AppendEntriesResponse {
	res := &AppendEntriesResponse{}
	if req.Term > raft.CurrentTerm {
		raft.BecomeFollowerWithLogAppend(req)
	}

	if raft.CurrentTerm > req.Term {
		res.Term = raft.CurrentTerm
		return res
	} else if !raft.isMatchPrevLog(req) {
		res.Term = raft.CurrentTerm
		return res
	}

	raft.CommitIndex = req.LeaderCommit
	raft.LeaderId = req.LeaderId
	lastEntryIdx := req.LeaderCommit
	for _, entry := range req.Entries {

		if raft.LastLog().Index >= entry.Index && raft.Log[entry.Index].Term == entry.Term {
			continue
		} else if raft.LastLog().Index >= entry.Index && raft.Log[entry.Index].Term != entry.Term {
			raft.Log = raft.Log[:entry.Index]
		}

		raft.Log = append(raft.Log, entry)
		err := raft.SaveLog(&entry)
		if err != nil {
			res.Err = err
			return res
		}

		lastEntryIdx = entry.Index
	}

	raft.CommitIndex = int64(infrastructure.Min(int(req.LeaderCommit), int(lastEntryIdx)))
	raft.CommitCh <- struct{}{}
	res.Success = true
	return res
}

func (raft *Raft) isMatchPrevLog(req *AppendEntriesRequest) bool {
	return raft.LastLog().Index >= req.PrevLogIndex && raft.Log[req.PrevLogIndex].Term == req.PrevLogTerm
}
