// raft集群间grpc通信
package apiServer

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/cmd"
	"BlackKingBar/config"
	"context"
	"log"
	"net"

	"github.com/copier"
	"google.golang.org/grpc"
)

type RaftServer struct {
}

func StartRpc() error {

	rpcServer := grpc.NewServer()
	s := RaftServer{}
	rpcProto.RegisterElectionServer(rpcServer, &s)

	cfg := config.CfgInstance

	listener, err := net.Listen("tcp", cfg.HttpIP+":"+cfg.RpcPort)
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	return rpcServer.Serve(listener)
}

func (s *RaftServer) RequestVote(ctx context.Context, request *rpcProto.VoteReq) (*rpcProto.VoteRes, error) {

	machine := cmd.MachineInstance
	req := &cmd.VoteRequest{}
	copier.Copy(req, request)
	//req.Term = int(request.Term)
	//req.CandicateId = int(request.Candidate)
	//req.LastLogTerm = int(request.LastLogTerm)
	//req.LastLogIndex = int(request.LastLogIndex)

	res, err := machine.HandleElection(req)

	resPonse := &rpcProto.VoteRes{}
	copier.Copy(resPonse, res)
	//resPonse.Term = int64(res.Term)
	//resPonse.VoteGranted = res.VoteGanted

	return resPonse, err
}
