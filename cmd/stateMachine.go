package cmd

// raft状态机

const (
	Follower  = 0
	Candidate = 1
	Leader    = 2
)

type StateMachine struct {
	candicateId   int32
	state         int32
	currentTerm   int
	voteFor       int32
	log           []int
	commitIndex   int
	lastApplied   int
	nextIndex     []int
	matchIndex    []int
	electionTimer []int8
}

var ins *StateMachine

func GetSingleton() *StateMachine {

	if ins == nil {

		ins = &StateMachine{}
	}

	return ins
}
