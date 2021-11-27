package cmd

type VoteRequest struct {
	Term         int
	CandicateId  int
	LastLogIndex int
	LastLogTerm  int
	RpcIP        string
	RpcPort      string
}

type VoteResponse struct {
	Term       int
	VoteGanted bool
}

func (m *StateMachine) StartElection() {

}

func (m *StateMachine) HandleElection(request *VoteRequest) (*VoteResponse, error) {
	return &VoteResponse{}, nil
}
