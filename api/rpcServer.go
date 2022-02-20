// raft集群间grpc通信
package apiServer

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/infrastructure"
	"BlackKingBar/raft"
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/copier"
	"google.golang.org/grpc"
)

type RaftServer struct {
}

func StartRpc(ctx context.Context) {

	rpcServer := grpc.NewServer()
	s := RaftServer{}
	rpcProto.RegisterElectionServer(rpcServer, &s)
	rpcProto.RegisterReplicateServer(rpcServer, &s)
	cfg := infrastructure.CfgInstance

	listener, err := net.Listen("tcp", cfg.HttpIP+":"+cfg.RpcPort)
	if err != nil {
		if err != nil {
			panic("start rpcServer err," + err.Error())
		}
	}
	go func() {
		err = rpcServer.Serve(listener)
		if err != nil {
			panic("start rpcServer err," + err.Error())
		}
	}()

	<-ctx.Done()

	rpcServer.Stop()
}

func (s *RaftServer) RequestVote(ctx context.Context, request *rpcProto.VoteReq) (*rpcProto.VoteRes, error) {

	if ctx.Err() == context.Canceled {
		return nil, ctx.Err()
	}

	rf := raft.R
	req := &raft.VoteRequest{}
	req.Done = make(chan *raft.VoteResponse)
	copier.Copy(req, request)

	if ctx.Err() == context.Canceled {
		return nil, ctx.Err()
	}

	rf.EventCh <- req
	res := <-req.Done
	resPonse := &rpcProto.VoteRes{}
	copier.Copy(resPonse, res)

	return resPonse, res.Err
}

func (s *RaftServer) AppendEntries(ctx context.Context, request *rpcProto.AppendEntriesReq) (*rpcProto.AppendEntriesRes, error) {
	if ctx.Err() == context.Canceled {
		return nil, ctx.Err()
	}
	fmt.Println("收到appendEntries" + strconv.FormatInt(int64(request.Term), 10))
	rf := raft.R
	rf.ResetElectionTicker()
	rf.ReceiveHeartBeatTime = time.Now()
	req := &raft.AppendEntriesRequest{}
	copier.Copy(req, request)
	req.Done = make(chan *raft.AppendEntriesResponse)

	if ctx.Err() == context.Canceled {
		return nil, ctx.Err()
	}
	rf.EventCh <- req
	res := <-req.Done
	resPonse := &rpcProto.AppendEntriesRes{}
	copier.Copy(resPonse, res)

	return resPonse, res.Err
}
