package app

import (
	"BlackKingBar/infrastructure"
	"fmt"
	"math/rand"
	"time"
)

// raft状态机

const (
	Follower  = 0
	Candidate = 1
	Leader    = 2
)

type StateMachine struct {
	CandidateId       int
	State             int
	CurrentTerm       int
	VoteFor           int
	CommitIndex       int
	LastApplied       int
	NextIndex         []int
	MatchIndex        []int
	Log               []RaftLog
	ElectionTimeout   time.Duration
	HeartBeatDuration time.Duration
	ElectionTimer     *time.Ticker
	HeartBeatTimer    *time.Ticker
}

var Raft *StateMachine

func InitStateMachine() error {

	cfg := infrastructure.CfgInstance
	Raft = &StateMachine{}
	Raft.CandidateId = cfg.CandidateId

	//从持久化文件中读取
	Raft.Log = make([]RaftLog, 0)
	Raft.CurrentTerm = 0
	//服务启动时状态为follower
	Raft.BecomeFollower()
	return nil
}

func (raft *StateMachine) BecomeFollower() {
	fmt.Println(time.Now())
	fmt.Println("成为follower")
	raft.State = Follower
	raft.ElectionTimeout = time.Duration(rand.Intn(15)+5) * time.Second
	raft.ElectionTimer = time.NewTicker(raft.ElectionTimeout)
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
	raft.HeartBeatDuration = time.Duration(5) * time.Second
	raft.HeartBeatTimer = time.NewTicker(raft.HeartBeatDuration)

	go func() {
		for {
			<-raft.HeartBeatTimer.C
			raft.SendHeartBeat()
		}
	}()
}

func (raft *StateMachine) stopHeartBeat() {
	if raft.HeartBeatTimer != nil {
		raft.HeartBeatTimer.Stop()
	}
}
