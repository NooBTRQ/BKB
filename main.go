package main

import (
	apiServer "BlackKingBar/api"
	"BlackKingBar/cmd"
	"BlackKingBar/config"
	"fmt"
)

// raft集群入口
func main() {

	//胶水代码
	// 1. 读取配置
	err := config.InitConfig()

	if err != nil {

		panic("启动服务失败！，获取配置文件出错")
	}
	// 2. 启动服务(利用context关掉其它服务)
	err = cmd.InitStateMachine()
	if err != nil {
		panic("启动服务失败！，启动raft状态机失败")
	}

	go apiServer.StartRpc()

	go apiServer.StartHttp()
	// 3. 处理输入

	fmt.Scanln()
}
