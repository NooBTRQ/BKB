package cmd

import (
	"BlackKingBar/config"
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
	CandidateId           int
	State                 int
	CurrentTerm           int
	VoteFor               int
	Log                   []RaftLog
	CommitIndex           int
	LastApplied           int
	NextIndex             []int
	MatchIndex            []int
	ElectionTimerDuration time.Duration
	ElectionTimer         *time.Timer
}

type RaftLog struct {
	Term  int
	Index int
}

var MachineInstance *StateMachine

func InitStateMachine() error {

	cfg := config.CfgInstance
	MachineInstance = &StateMachine{}
	MachineInstance.CurrentTerm = 0
	MachineInstance.ElectionTimerDuration = time.Duration(rand.Intn(150)+150) * time.Millisecond
	MachineInstance.State = Follower
	MachineInstance.ElectionTimer = time.NewTimer(MachineInstance.ElectionTimerDuration)
	MachineInstance.CandidateId = cfg.CandidateId
	MachineInstance.Log = make([]RaftLog, 0)
	for _ = range MachineInstance.ElectionTimer.C {

		go MachineInstance.StartElection()
	}

	return nil
}
