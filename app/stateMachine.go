package app

import (
	"BlackKingBar/infrastructure"
	"fmt"
	"sync"
	"time"
)

// raft状态机

const (
	Follower  = 0
	Candidate = 1
	Leader    = 2
)

type StateMachine struct {
	CandidateId       int8
	State             int8
	CurrentTerm       int64
	VoteFor           int8
	CommitIndex       int64
	LastApplied       int64
	NextIndex         map[int8]int64
	MatchIndex        map[int8]int64
	Log               []RaftLog
	ElectionTimeout   time.Duration
	HeartBeatDuration time.Duration
	ElectionTimer     *time.Ticker
	HeartBeatTimer    *time.Ticker
	sync.RWMutex
	sync.WaitGroup
}

var Raft *StateMachine

func InitStateMachine() error {
	cfg := infrastructure.CfgInstance
	Raft = &StateMachine{}
	Raft.CandidateId = cfg.CandidateId
	Raft.ElectionTimer = time.NewTicker(time.Second * 3600)
	Raft.HeartBeatTimer = time.NewTicker(time.Second * 3600)
	//从持久化存储中读取重启前数据
	Raft.Log = make([]RaftLog, 0)
	Raft.Log = append(Raft.Log, RaftLog{})
	Raft.CurrentTerm = 0

	//服务启动时状态为follower
	Raft.BecomeFollower()
	return nil
}

func (raft *StateMachine) BecomeFollower() {
	fmt.Println(time.Now())
	fmt.Println("成为follower")
	raft.State = Follower

	raft.ElectionTimeout = time.Duration(infrastructure.GetRandNumber(150, 300, int(raft.CandidateId))) * time.Millisecond
	raft.ElectionTimer.Reset(raft.ElectionTimeout)
	raft.stopHeartBeat()
	go func() {
		for {
			<-raft.ElectionTimer.C
			raft.StartElection()
		}
	}()
}

func (raft *StateMachine) BecomeCandicate() {
	fmt.Println(time.Now())
	fmt.Println("成为candidate")
	raft.State = Candidate
	raft.CurrentTerm++
	raft.VoteFor = raft.CandidateId
}

func (raft *StateMachine) BecomeLeader() {
	fmt.Println(time.Now())
	fmt.Println("成为leader")
	raft.State = Leader
	raft.ElectionTimer.Stop()
	raft.HeartBeatDuration = time.Duration(10) * time.Millisecond
	raft.HeartBeatTimer.Reset(raft.HeartBeatDuration)
	raft.MatchIndex = make(map[int8]int64)
	raft.NextIndex = make(map[int8]int64)
	cfg := infrastructure.CfgInstance

	for _, member := range cfg.ClusterMembers {
		raft.NextIndex[member.CandidateId] = int64(len(raft.Log))
	}
	go func() {
		for {
			<-raft.HeartBeatTimer.C
			raft.SendHeartBeat()
		}
	}()
}

func (raft *StateMachine) stopHeartBeat() {
	raft.HeartBeatTimer.Stop()
}

func (raft *StateMachine) ResetElectionTimeout() {
	raft.ElectionTimer.Stop()
	raft.ElectionTimeout = time.Duration(infrastructure.GetRandNumber(150, 300, int(raft.CandidateId))) * time.Millisecond
	raft.ElectionTimer.Reset(raft.ElectionTimeout)
}
