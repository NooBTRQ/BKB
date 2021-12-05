// raft集群间grpc通信
package apiServer

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/app"
	"BlackKingBar/infrastructure"
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/copier"
	"google.golang.org/grpc"
)

type RaftServer struct {
}

var rpcServer *grpc.Server

func StartRpc() error {

	rpcServer := grpc.NewServer()
	s := RaftServer{}
	rpcProto.RegisterElectionServer(rpcServer, &s)
	rpcProto.RegisterReplicateServer(rpcServer, &s)
	cfg := infrastructure.CfgInstance

	listener, err := net.Listen("tcp", cfg.HttpIP+":"+cfg.RpcPort)
	if err != nil {
		return err
	}
	return rpcServer.Serve(listener)
}

func StopRpc() {
	rpcServer.Stop()
}

func (s *RaftServer) RequestVote(ctx context.Context, request *rpcProto.VoteReq) (*rpcProto.VoteRes, error) {

	raft := app.Raft
	req := &app.VoteRequest{}
	copier.Copy(req, request)
	fmt.Println("收到投票请求,candidateId:" + strconv.FormatInt(int64(req.CandidateId), 10) + "termId:" + strconv.FormatInt(int64(req.Term), 10))
	res, err := raft.HandleElection(req)
	resPonse := &rpcProto.VoteRes{}
	copier.Copy(resPonse, res)

	return resPonse, err
}

func (s *RaftServer) AppendEntries(ctx context.Context, request *rpcProto.AppendEntriesReq) (*rpcProto.AppendEntriesRes, error) {
	raft := app.Raft
	req := &app.AppendEntriesRequest{}
	copier.Copy(req, request)
	res, err := raft.HandleAppendEntries(req)
	resPonse := &rpcProto.AppendEntriesRes{}
	copier.Copy(resPonse, res)

	return resPonse, err
}
