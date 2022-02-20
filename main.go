package main

import (
	apiServer "BlackKingBar/api"
	"BlackKingBar/infrastructure"
	"BlackKingBar/raft"
	"context"
	"fmt"
)

const (
	Quit   = "quit"
	Data   = "data"
	Log    = "log"
	Term   = "term"
	Leader = "leader"
)

// raft集群入口
func main() {

	//胶水代码
	// 1. 读取配置
	err := infrastructure.InitConfig()

	if err != nil {
		panic("start raft server err,init config err," + err.Error())
	}
	// 2. 启动服务(利用context关掉其它服务)
	ctx, cancel := context.WithCancel(context.Background())
	err = raft.InitStateMachine(ctx)
	if err != nil {
		panic("start raft server err,init stateMachine err," + err.Error())
	}
	go apiServer.StartRpc(ctx)

	go apiServer.StartHttp(ctx)

	// 启动消费协程
	go raft.R.Process(ctx)

	// 3. 处理输入
	var cmd string
	for {
		fmt.Scanln(&cmd)

		switch cmd {
		case Quit:
			cancel()
			return
		case Term:
			fmt.Println(raft.R.CurrentTerm)
		case Leader:
			fmt.Println(raft.R.LeaderId)
		case Data:
			for k, v := range raft.R.Data {
				fmt.Printf("key:%s,value:%s\n", k, v)
			}
		case Log:
			for _, log := range raft.R.Log {
				fmt.Printf("term:%d,index:%d,key:%s,op:%d,value:%s\n", log.Term, log.Index, log.Cmd.Key, log.Cmd.Operation, log.Cmd.Value)
			}

		}
	}
}
