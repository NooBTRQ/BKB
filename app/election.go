package app

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/infrastructure"
	"fmt"
	"strconv"
	"sync"
)

type VoteRequest struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type VoteResponse struct {
	Term        int
	VoteGranted bool
}

func (raft *StateMachine) StartElection() {

	raft.BecomeCandicate()

	var wg sync.WaitGroup
	var voteCount int
	var lock sync.Mutex
	cfg := infrastructure.CfgInstance
	wg.Add(len(cfg.ClusterMembers))
	//给其它节点发送voteRequest
	for _, node := range cfg.ClusterMembers {
		go func(node infrastructure.ClusterMember) {
			rpcRequest := &rpcProto.VoteReq{}
			rpcRequest.CandidateId = int64(raft.CandidateId)
			rpcRequest.Term = int64(raft.CurrentTerm)
			if len(raft.Log) > 0 {
				rpcRequest.LastLogTerm = int64(raft.Log[len(raft.Log)-1].Term)
				rpcRequest.LastLogIndex = int64(raft.Log[len(raft.Log)-1].Index)
			}
			res, err := infrastructure.SendVoteRequest(rpcRequest, node.RpcIP, node.RpcPort)
			if err == nil && res.VoteGranted {
				lock.Lock()
				voteCount++
				lock.Unlock()
			}

			wg.Done()
		}(node)
	}

	wg.Wait()
	fmt.Println("收到票数：" + strconv.FormatInt(int64(voteCount), 10))
	if voteCount > len(cfg.ClusterMembers)/2 {

		raft.BecomeLeader()
	}
}

func (raft *StateMachine) HandleElection(request *VoteRequest) (*VoteResponse, error) {
	raft.ElectionTimer.Reset(raft.ElectionTimeout)
	res := &VoteResponse{}
	var err error
	if request.Term < raft.CurrentTerm {
		res.Term = raft.CurrentTerm

	} else if request.Term >= raft.CurrentTerm {
		fmt.Println("已投票,candidateId:" + strconv.FormatInt(int64(request.CandidateId), 10))
		raft.VoteFor = request.CandidateId
		res.Term = request.Term
		res.VoteGranted = true

		return res, nil
	} //else

	if (raft.VoteFor == 0 || raft.VoteFor == request.CandidateId) && (len(raft.Log) == 0 || (raft.Log[len(raft.Log)].Term <= request.LastLogTerm && raft.Log[len(raft.Log)].Index <= request.LastLogIndex)) {

		raft.VoteFor = raft.CandidateId
		res.Term = request.Term
		res.VoteGranted = true
	}

	return res, err
}
