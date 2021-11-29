package cmd

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

func (m *StateMachine) StartElection() {

	m.BecomeCandicate()

	var wg sync.WaitGroup
	var voteCount int
	var lock sync.Mutex
	cfg := infrastructure.CfgInstance
	wg.Add(len(cfg.ClusterMembers))
	//给其它节点发送voteRequest
	for _, node := range cfg.ClusterMembers {
		go func(node infrastructure.ClusterMember) {
			rpcRequest := &rpcProto.VoteReq{}
			rpcRequest.CandidateId = int64(m.CandidateId)
			rpcRequest.Term = int64(m.CurrentTerm)
			if len(m.Log) > 0 {
				rpcRequest.LastLogTerm = int64(m.Log[len(m.Log)-1].Term)
				rpcRequest.LastLogIndex = int64(m.Log[len(m.Log)-1].Index)
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

		m.BecomeLeader()
	}
}

func (m *StateMachine) HandleElection(request *VoteRequest) (*VoteResponse, error) {
	m.ElectionTimer.Reset(m.ElectionTimeout)
	res := &VoteResponse{}
	var err error
	if request.Term < m.CurrentTerm {
		res.Term = m.CurrentTerm

	} else if request.Term >= m.CurrentTerm {
		fmt.Println("已投票,candidateId:" + strconv.FormatInt(int64(request.CandidateId), 10))
		m.VoteFor = request.CandidateId
		res.Term = request.Term
		res.VoteGranted = true

		return res, nil
	} //else

	if (m.VoteFor == 0 || m.VoteFor == request.CandidateId) && (len(m.Log) == 0 || (m.Log[len(m.Log)].Term <= request.LastLogTerm && m.Log[len(m.Log)].Index <= request.LastLogIndex)) {

		m.VoteFor = m.CandidateId
		res.Term = request.Term
		res.VoteGranted = true
	}

	return res, err
}
