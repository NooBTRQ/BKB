// raft集群间grpc通信
package apiServer

import (
	"BlackKingBar/api/rpcProto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
}

func StartRpc() error {

	rpcServer := grpc.NewServer()
	s := Server{}
	rpcProto.RegisterElectionServer(rpcServer, &s)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	// 4. 运行rpcServer，传入listener
	return rpcServer.Serve(listener)
}

func (s *Server) RequestVote(ctx context.Context, request *rpcProto.VoteReq) (*rpcProto.VoteRes, error) {

	return &rpcProto.VoteRes{}, nil
}
