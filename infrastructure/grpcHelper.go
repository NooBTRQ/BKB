package infrastructure

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/cmd"
	"context"

	"github.com/copier"
	"google.golang.org/grpc"
)

func SendVoteRequest(req *cmd.VoteRequest) (*cmd.VoteResponse, error) {
	conn, err := grpc.Dial(req.RpcIP+":"+req.RpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rpcProto.NewElectionClient(conn)
	rpcReq := &rpcProto.VoteReq{}
	copier.Copy(rpcReq, req)
	rpcRes, err := client.RequestVote(context.Background(), rpcReq)
	res := &cmd.VoteResponse{}
	copier.Copy(res, rpcRes)
	return res, err
}

func SendLogReplicate() {

}
