# BlackKingBar
a distribute linearizability k-v storge.
practice raft and golang

main function:
1. election ✔
2. heartBeat ✔
3. logReplicate ✔
4. k-v storge ✔
5. member change
6. snapshot
7. log compaction

TODO:
1. 服务下线相关功能 ✔
2. 补充test(节点超时、下线等测试用例)
3. prevote
4. Batching,pipelining
5. readindex ✔
6. lead lease
