package app

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/infrastructure"
	"fmt"
	"strconv"
	"sync"
)

type RaftLog struct {
	Term  int
	Index int
}

type AppendRequest struct {
	Term         int
	PrevLogTerm  int
	PrevLogIndex int
	LeaderCommit int
	Entries      []RaftLog
}

type AppendResponse struct {
	Term    int
	Success bool
}

func (raft *StateMachine) SendHeartBeat() {

	cfg := infrastructure.CfgInstance
	var wg sync.WaitGroup
	wg.Add(len(cfg.ClusterMembers))

	for _, node := range cfg.ClusterMembers {
		go func(node infrastructure.ClusterMember) {
			request := &rpcProto.AppendReq{}
			request.Term = int64(raft.CurrentTerm)
			request.LeaderId = int64(raft.CandidateId)
			request.LeaderCommit = int64(raft.CommitIndex)
			if len(raft.Log) > 0 {
				request.PrevLogTerm = int64(raft.Log[len(raft.Log)-1].Term)
				request.PrevLogIndex = int64(raft.Log[len(raft.Log)-1].Index)
			}
			infrastructure.SendLogReplicate(request, node.RpcIP, node.RpcPort)
			wg.Done()
		}(node)
	}

	wg.Wait()
}

func (raft *StateMachine) HandleAppendEntries(req *AppendRequest) (*AppendResponse, error) {

	fmt.Println("收到心跳，currntTerm:" + strconv.FormatInt(int64(raft.CurrentTerm), 10))

	res := &AppendResponse{}
	if raft.CurrentTerm > req.Term {
		res.Term = raft.CurrentTerm
		return res, nil
	} else if raft.CurrentTerm <= req.Term {

		raft.CurrentTerm = req.Term
		if raft.State == Leader || raft.State == Candidate {
			raft.BecomeFollower()
		} else {
			raft.ElectionTimer.Reset(Raft.ElectionTimeout)
		}
	}

	return res, nil
}
