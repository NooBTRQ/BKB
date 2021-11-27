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
	CandicateId           int
	State                 int
	CurrentTerm           int
	VoteFor               int
	Log                   []int
	CommitIndex           int
	LastApplied           int
	NextIndex             []int
	MatchIndex            []int
	ElectionTimerDuration int
	ElectionTimer         *time.Timer
}

var MachineInstance *StateMachine

func InitStateMachine() error {

	cfg := config.CfgInstance
	MachineInstance = &StateMachine{}
	MachineInstance.CurrentTerm = 0
	MachineInstance.ElectionTimerDuration = rand.Intn(150) + 150
	MachineInstance.State = Follower
	MachineInstance.ElectionTimer = time.NewTimer(time.Millisecond * time.Duration(MachineInstance.ElectionTimerDuration))
	MachineInstance.CandicateId = cfg.CandicateId
	for _ = range MachineInstance.ElectionTimer.C {

		go MachineInstance.StartElection()
	}

	return nil
}
