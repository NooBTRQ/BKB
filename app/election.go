package app

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/infrastructure"
	"fmt"
	"strconv"
	"sync"
)

type VoteRequest struct {
	Term         int64
	CandidateId  int8
	LastLogIndex int64
	LastLogTerm  int64
}

type VoteResponse struct {
	Term        int64
	VoteGranted bool
}

func (raft *StateMachine) StartElection() {

	raft.lock.Lock()
	defer raft.lock.Unlock()
	raft.BecomeCandicate()
	for !raft.sendVote() {
		raft.ResetElectionTimeout()
	}
}

func (raft *StateMachine) sendVote() bool {
	fmt.Println("选举timerout，" + strconv.FormatInt(int64(raft.ElectionTimeout), 10))
	voteCount := 1
	var lock sync.Mutex
	cfg := infrastructure.CfgInstance
	raft.wg.Add(len(cfg.ClusterMembers))
	for _, node := range cfg.ClusterMembers {
		go func(node infrastructure.ClusterMember) {
			rpcRequest := &rpcProto.VoteReq{}
			rpcRequest.CandidateId = int32(raft.CandidateId)
			rpcRequest.Term = int64(raft.CurrentTerm)
			if len(raft.Log) > 0 {
				rpcRequest.LastLogTerm = int64(raft.Log[len(raft.Log)-1].Term)
				rpcRequest.LastLogIndex = int64(len(raft.Log) - 1)
			}
			res, err := infrastructure.SendVoteRequest(rpcRequest, node.RpcIP, node.RpcPort)
			if err == nil && res.VoteGranted {
				lock.Lock()
				voteCount++
				lock.Unlock()
			}

			raft.wg.Done()
		}(node)
	}

	raft.wg.Wait()
	fmt.Println("收到票数：" + strconv.FormatInt(int64(voteCount), 10))
	if voteCount > len(cfg.ClusterMembers)/2 {

		raft.BecomeLeader()
		return true
	} else {
		return false
	}
}

func (raft *StateMachine) HandleElection(request *VoteRequest) (*VoteResponse, error) {

	raft.lock.Lock()
	defer raft.lock.Unlock()
	raft.ElectionTimer.Reset(raft.ElectionTimeout)
	res := &VoteResponse{}
	var err error
	if request.Term < raft.CurrentTerm {
		res.Term = raft.CurrentTerm
		return res, nil
	} else if request.Term > raft.CurrentTerm {
		raft.VoteFor = request.CandidateId
		raft.CurrentTerm = request.Term
		res.Term = request.Term
		res.VoteGranted = true
		raft.BecomeFollower()
		return res, nil
	}

	if (raft.VoteFor == 0 || raft.VoteFor == request.CandidateId) && isUpToDateLog(raft, request) {

		raft.VoteFor = raft.CandidateId
		res.Term = request.Term
		res.VoteGranted = true
		raft.BecomeFollower()
	}

	return res, err
}

func isUpToDateLog(raft *StateMachine, req *VoteRequest) bool {

	return len(raft.Log) == 0 || (int64(len(raft.Log))-1 <= req.LastLogIndex && raft.Log[len(raft.Log)-1].Term <= req.Term)
}
