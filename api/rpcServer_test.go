package apiServer

import (
	"BlackKingBar/api/rpcProto"
	"BlackKingBar/raft"
	"fmt"
	"testing"

	"github.com/copier"
)

func TestDataConvert(t *testing.T) {
	request := rpcProto.AppendEntriesReq{}
	request.Term = 1
	request.LeaderId = 2
	request.Entries = make([]*rpcProto.LogEntry, 0)
	request.Entries = append(request.Entries, &rpcProto.LogEntry{Operation: 1, Key: "1111", Value: "aaaaa"})

	req := &raft.AppendEntriesRequest{}
	copier.Copy(req, &request)
	a := req.Entries[0]
	fmt.Println(a.Cmd.Key)
}
