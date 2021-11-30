package infrastructure

import (
	"BlackKingBar/api/rpcProto"
	"context"
	"time"

	"google.golang.org/grpc"
)

func SendVoteRequest(req *rpcProto.VoteReq, rpcIp, rpcPort string) (*rpcProto.VoteRes, error) {

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(150*time.Millisecond)))
	conn, err := grpc.DialContext(ctx, rpcIp+":"+rpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rpcProto.NewElectionClient(conn)
	res, err := client.RequestVote(ctx, req)
	return res, err
}

func SendLogReplicate(req *rpcProto.AppendReq, rpcIp, rpcPort string) (*rpcProto.AppendRes, error) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(50*time.Millisecond)))
	conn, err := grpc.DialContext(ctx, rpcIp+":"+rpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rpcProto.NewReplicateClient(conn)
	res, err := client.AppendEntries(ctx, req)
	return res, err
}
