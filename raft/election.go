package raft

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/infrastructure"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

type VoteRequest struct {
	Term         int64
	CandidateId  int8
	LastLogIndex int64
	LastLogTerm  int64
	Done         chan *VoteResponse
}

type VoteResponse struct {
	Term        int64
	VoteGranted bool
	Err         error
}

type Election struct{}

func (raft *Raft) StartElection() {
	raft.BecomeCandicate()
	raft.ResetElectionTicker()
	raft.RequetVote()
	raft.SaveMachine()

}

// 发起投票，若获得半数以上选票，则变为leader
func (raft *Raft) RequetVote() {
	raft.StopElectionTicker()
	fmt.Println("发起选举" + strconv.FormatInt(int64(raft.CurrentTerm), 10))
	var voteCount int64 = 1
	cfg := infrastructure.CfgInstance
	var wg sync.WaitGroup
	wg.Add(len(cfg.ClusterMembers))
	for _, node := range cfg.ClusterMembers {
		go func(node infrastructure.ClusterMember) {
			rpcRequest := &rpcProto.VoteReq{}
			rpcRequest.CandidateId = int32(raft.NodeId)
			rpcRequest.Term = int64(raft.CurrentTerm)
			if len(raft.Log) > 0 {
				rpcRequest.LastLogTerm = raft.LastLog().Term
				rpcRequest.LastLogIndex = raft.LastLog().Index
			}
			res, err := infrastructure.SendVoteRequest(rpcRequest, node.RpcIP, node.RpcPort)
			if err == nil && res.VoteGranted && res.Term <= raft.CurrentTerm {
				atomic.AddInt64(&voteCount, 1)
			} else if err == nil {

				raft.CurrentTerm = res.Term
			} else if err != nil {
				fmt.Print(rpcRequest.Term)
				fmt.Println(": " + err.Error())
			}
			wg.Done()
		}(node)
	}
	wg.Wait()
	if atomic.LoadInt64(&voteCount) > cfg.HalfCount() {

		raft.BecomeLeader()
	} else {
		raft.StartElection()
	}
}

func (raft *Raft) HandleVoteRequest(request *VoteRequest) *VoteResponse {
	fmt.Println("处理投票" + strconv.FormatInt(int64(request.Term), 10))
	res := &VoteResponse{}
	var err error
	if request.Term < raft.CurrentTerm {
		res.Term = raft.CurrentTerm
		return res
	} else if request.Term > raft.CurrentTerm {
		raft.VoteFor = request.CandidateId
		raft.CurrentTerm = request.Term
		res.Term = request.Term
		res.VoteGranted = true
		raft.BecomeFollower(nil)
		return res
	}

	if (raft.VoteFor == 0 || raft.VoteFor == request.CandidateId) && isUpToDateLog(raft, request) {

		raft.VoteFor = request.CandidateId
		raft.CurrentTerm = request.Term
		res.Term = request.Term
		res.VoteGranted = true
		raft.BecomeFollower(nil)
	}
	err = raft.SaveMachine()
	res.Err = err
	fmt.Println(res.VoteGranted)
	return res
}

func isUpToDateLog(raft *Raft, req *VoteRequest) bool {

	return len(raft.Log) == 0 || (raft.LastLog().Index <= req.LastLogIndex && raft.LastLog().Term <= req.Term)
}
