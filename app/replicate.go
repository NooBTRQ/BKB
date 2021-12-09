package app

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
}

type AppendEntriesRequest struct {
	Term         int64
	PrevLogTerm  int64
	PrevLogIndex int64
	LeaderCommit int64
	Entries      []RaftLog
}

type AppendEntriesResponse struct {
	Term    int64
	Success bool
}

func (raft *StateMachine) SendHeartBeat() {
	raft.lock.Lock()
	defer raft.lock.Unlock()
	cfg := infrastructure.CfgInstance
	var wg sync.WaitGroup
	wg.Add(len(cfg.ClusterMembers))

	for _, node := range cfg.ClusterMembers {
		go func(node infrastructure.ClusterMember) {
			request := &rpcProto.AppendEntriesReq{}
			request.Term = int64(raft.CurrentTerm)
			request.LeaderId = int32(raft.CandidateId)
			request.LeaderCommit = int64(raft.CommitIndex)
			request.Entries = make([]*rpcProto.LogEntry, 0)
			if len(raft.Log) > 0 {
				request.PrevLogIndex = int64(raft.NextIndex[node.CandidateId]) - 1
				request.PrevLogTerm = int64(raft.Log[request.PrevLogIndex].Term)
				entry := &rpcProto.LogEntry{}
				entry.Index = int64(raft.NextIndex[node.CandidateId])
				entry.Term = int64(raft.Log[raft.NextIndex[node.CandidateId]].Term)
				request.Entries = append(request.Entries, entry)
			}
			res, err := infrastructure.SendLogReplicate(request, node.RpcIP, node.RpcPort)
			if err != nil {
				return
			}

			if res.Success {

				raft.NextIndex[node.CandidateId]++
				raft.MatchIndex[node.CandidateId]++
			} else {
				raft.NextIndex[node.CandidateId]--
			}
			wg.Done()
		}(node)
	}

	wg.Wait()
}

func (raft *StateMachine) HandleAppendEntries(req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	raft.lock.Lock()
	defer raft.lock.Unlock()
	fmt.Println("收到心跳，currntTerm:" + strconv.FormatInt(int64(raft.CurrentTerm), 10))

	raft.ResetElectionTimeout()
	res := &AppendEntriesResponse{}
	if raft.CurrentTerm > req.Term {
		res.Term = raft.CurrentTerm
		return res, nil
	} else if !isMatchPrevLog(req, raft) {
		res.Term = raft.CurrentTerm
		return res, nil
	}

	raft.CommitIndex = req.LeaderCommit
	for _, entry := range req.Entries {
		if int64(len(raft.Log))-1 >= entry.Index && raft.Log[entry.Index].Term != entry.Term {
			raft.Log = raft.Log[:entry.Index]
		}

		raft.Log = append(raft.Log, entry)
		if raft.CommitIndex > entry.Index {
			raft.CommitIndex = entry.Index
		}
	}

	res.Success = true
	return res, nil
}

func isMatchPrevLog(req *AppendEntriesRequest, raft *StateMachine) bool {
	if int64(len(raft.Log)) > req.PrevLogIndex && raft.Log[req.PrevLogIndex].Term == req.Term {
		return true
	}
	return false
}
