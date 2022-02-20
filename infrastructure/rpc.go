package infrastructure

import (
	"BlackKingBar/api/rpcProto"
	"context"
	"time"

	"google.golang.org/grpc"
)

func SendVoteRequest(req *rpcProto.VoteReq, rpcIp, rpcPort string) (*rpcProto.VoteRes, error) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(10*time.Millisecond)))
	defer cancel()
	conn, err := grpc.DialContext(ctx, rpcIp+":"+rpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rpcProto.NewElectionClient(conn)
	res, err := client.RequestVote(ctx, req)
	return res, err
}

func AppendEntries(req *rpcProto.AppendEntriesReq, rpcIp, rpcPort string) (*rpcProto.AppendEntriesRes, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(10*time.Millisecond)))
	defer cancel()
	conn, err := grpc.DialContext(ctx, rpcIp+":"+rpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rpcProto.NewReplicateClient(conn)
	res, err := client.AppendEntries(ctx, req)
	return res, err
}
