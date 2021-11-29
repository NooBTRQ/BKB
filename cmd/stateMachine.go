package cmd

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

var MachineInstance *StateMachine

func InitStateMachine() error {

	cfg := infrastructure.CfgInstance
	MachineInstance = &StateMachine{}
	MachineInstance.CandidateId = cfg.CandidateId

	//从持久化文件中读取
	MachineInstance.Log = make([]RaftLog, 0)
	MachineInstance.CurrentTerm = 0
	//
	MachineInstance.BecomeFollower()
	return nil
}

func (m *StateMachine) BecomeFollower() {
	fmt.Println(time.Now())
	fmt.Println("成为follower")
	m.State = Follower
	m.ElectionTimeout = time.Duration(rand.Intn(150)+150) * time.Millisecond
	m.ElectionTimer = time.NewTicker(m.ElectionTimeout)
	m.stopHeartBeat()
	go func() {
		for {
			<-m.ElectionTimer.C
			m.StartElection()
		}
	}()
}

func (m *StateMachine) BecomeCandicate() {
	fmt.Println(time.Now())
	fmt.Println("成为candidate")
	m.State = Candidate
	m.CurrentTerm++
	m.VoteFor = m.CandidateId
}

func (m *StateMachine) BecomeLeader() {
	fmt.Println(time.Now())
	fmt.Println("成为leader")
	m.State = Leader
	m.ElectionTimer.Stop()
	m.HeartBeatDuration = time.Duration(50) * time.Millisecond
	m.HeartBeatTimer = time.NewTicker(m.HeartBeatDuration)

	go func() {
		for {
			<-m.HeartBeatTimer.C
			m.SendHeartBeat()
		}
	}()
}

func (m *StateMachine) stopHeartBeat() {
	if m.HeartBeatTimer != nil {
		m.HeartBeatTimer.Stop()
	}
}
