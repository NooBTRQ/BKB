package cmd

import (
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

var ins *StateMachine

func GetSingleton() *StateMachine {

	if ins == nil {

		ins = &StateMachine{}
		ins.CurrentTerm = 0
		ins.ElectionTimerDuration = rand.Intn(150) + 150
		ins.State = Follower
		ins.ElectionTimer = time.NewTimer(time.Millisecond * time.Duration(ins.ElectionTimerDuration))

		for _ = range ins.ElectionTimer.C {

			ins.election()
		}
	}

	return ins
}
