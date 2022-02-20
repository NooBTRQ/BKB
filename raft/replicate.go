package raft

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/infrastructure"
	"fmt"
	"strconv"
	"sync"
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
		if op.Operation == Set {
			raft.Data[op.Key] = op.Value
		} else if op.Operation == Delete {
			delete(raft.Data, op.Key)
		}
		raft.LastApplied = log.Index
		return true
	}
	return false
}

func (raft *Raft) HeartBeat() {
	fmt.Println("发送心跳:" + strconv.FormatInt(int64(raft.CurrentTerm), 10))
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
			if len(raft.Log) > 1 {

				request.PrevLogIndex = nextIdx - 1
				request.PrevLogTerm = int64(raft.Log[request.PrevLogIndex].Term)
				if raft.LastLog().Index >= nextIdx {
					sendLogs := raft.Log[nextIdx:]
					for _, log := range sendLogs {
						entry := &rpcProto.LogEntry{}
						entry.Index = nextIdx
						entry.Term = int64(raft.Log[nextIdx].Term)
						entry.Key = log.Cmd.Key
						entry.Operation = int64(log.Cmd.Operation)
						entry.Value = log.Cmd.Value
						request.Entries = append(request.Entries, entry)
					}
				}
			}

			res, err := infrastructure.AppendEntries(request, node.RpcIP, node.RpcPort)
			if err != nil {
				fmt.Printf("复制日志:" + err.Error())
			} else if res.Success {
				raft.MatchIndex.Store(node.NodeId, nextIdx)
				raft.NextIndex.Store(node.NodeId, nextIdx+1)
			} else {
				raft.NextIndex.Store(node.NodeId, nextIdx-1)
			}
			wg.Done()
		}(node)
	}
	wg.Wait()
}

func (raft *Raft) HandleAppendEntries(req *AppendEntriesRequest) *AppendEntriesResponse {
	res := &AppendEntriesResponse{}
	if req.Term > raft.CurrentTerm && raft.State != Follower {
		raft.BecomeFollower(req)
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
	for _, entry := range req.Entries {
		if raft.LastLog().Index > entry.Index && raft.Log[entry.Index].Term != entry.Term {
			raft.Log = raft.Log[:entry.Index]
		}

		raft.Log = append(raft.Log, entry)
		err := raft.SaveLog(&entry)
		if err != nil {
			res.Err = err
			return res
		}
		if raft.CommitIndex > entry.Index {
			raft.CommitIndex = entry.Index
			if raft.CommitIndex > raft.LastApplied {

				if entry.Cmd.Operation == Set {
					raft.Data[entry.Cmd.Key] = entry.Cmd.Value
				} else if entry.Cmd.Operation == Delete {
					delete(raft.Data, entry.Cmd.Key)
				}
				raft.LastApplied = raft.CommitIndex
			}
		}
	}

	res.Success = true
	return res
}

func (raft *Raft) isMatchPrevLog(req *AppendEntriesRequest) bool {
	return raft.LastLog().Index >= req.PrevLogIndex && raft.Log[req.PrevLogIndex].Term == req.PrevLogTerm
}
