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

func (s *RaftServer) RequestVote(ctx context.Context, request *rpcProto.VoteReq) (*rpcProto.VoteRes, error) {

	raft := app.Raft
	req := &app.VoteRequest{}
	copier.Copy(req, request)
	fmt.Println("收到投票请求,candidateId:" + strconv.FormatInt(int64(req.CandidateId), 10))
	res, err := raft.HandleElection(req)
	resPonse := &rpcProto.VoteRes{}
	copier.Copy(resPonse, res)

	return resPonse, err
}

func (s *RaftServer) AppendEntries(ctx context.Context, request *rpcProto.AppendReq) (*rpcProto.AppendRes, error) {
	raft := app.Raft
	req := &app.AppendRequest{}
	copier.Copy(req, request)
	res, err := raft.HandleAppendEntries(req)
	resPonse := &rpcProto.AppendRes{}
	copier.Copy(resPonse, res)

	return resPonse, err
}
