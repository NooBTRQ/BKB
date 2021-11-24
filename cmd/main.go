package cmd

import (
	apiServer "BlackKingBar/api"
)

// raft集群入口
func main() {

	//胶水代码
	// 1. 读取配置
	// 2. 启动服务(利用context关掉其它服务)
	apiServer.StartHttp()
	apiServer.StartRpc()
	// 3. 处理输入
}
