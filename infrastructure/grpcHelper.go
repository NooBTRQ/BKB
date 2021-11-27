package infrastructure

import (
	"BlackKingBar/api/rpcProto"
	"context"

	"google.golang.org/grpc"
)

func SendVoteRequest(req *rpcProto.VoteReq, rpcIp, rpcPort string) (*rpcProto.VoteRes, error) {
	conn, err := grpc.Dial(rpcIp+":"+rpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rpcProto.NewElectionClient(conn)
	res, err := client.RequestVote(context.Background(), req)
	return res, err
}

func SendLogReplicate() {

}
