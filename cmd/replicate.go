package cmd

type RaftLog struct {
	Term  int
	Index int
}

type AppendRequest struct {
	Term         int
	PrevLogTerm  int
	PrevLogIndex int
	LeaderCommit int
	Entries      []RaftLog
}

type AppendResponse struct {
	Term    int
	Success bool
}

func (m *StateMachine) SendHeartBeat() {

}
