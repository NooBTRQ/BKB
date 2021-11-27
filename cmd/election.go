package cmd

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/config"
	"BlackKingBar/infrastructure"
	"sync"
)

type VoteRequest struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type VoteResponse struct {
	Term       int
	VoteGanted bool
}

func (m *StateMachine) StartElection() {
	m.State = Candidate
	m.CurrentTerm++
	m.VoteFor = m.CandidateId
	m.ElectionTimer.Stop()

	var wg sync.WaitGroup
	var voteCount int
	var lock sync.Mutex
	cfg := config.CfgInstance
	wg.Add(len(cfg.ClusterMember))
	//给其它节点发送voteRequest
	for _, node := range cfg.ClusterMember {
		go func() {
			rpcRequest := &rpcProto.VoteReq{}
			rpcRequest.CandidateId = int32(m.CandidateId)
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
		}()
	}

	wg.Wait()

	if voteCount > len(cfg.ClusterMember)/2 {

		m.State = Leader
	}
}

func (m *StateMachine) HandleElection(request *VoteRequest) (*VoteResponse, error) {

	res := &VoteResponse{}
	var err error
	if request.Term < m.CurrentTerm {
		res.Term = m.CurrentTerm

	} else if request.Term > m.CurrentTerm {
		m.State = Follower
		m.VoteFor = m.CandidateId
		res.Term = request.Term
		res.VoteGanted = true
		//持久化状态
		//重置定时器
		_ = m.ElectionTimer.Reset(m.ElectionTimerDuration)
		return res, nil
	} else if (m.VoteFor == 0 || m.VoteFor == request.CandidateId) && (len(m.Log) == 0 || (m.Log[len(m.Log)].Term <= request.LastLogTerm && m.Log[len(m.Log)].Index <= request.LastLogIndex)) {
		m.State = Follower
		m.VoteFor = m.CandidateId
		res.Term = request.Term
		res.VoteGanted = true
		//持久化状态
		//重置定时器
		_ = m.ElectionTimer.Reset(m.ElectionTimerDuration)
	}

	return res, err
}
