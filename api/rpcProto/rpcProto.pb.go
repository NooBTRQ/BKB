// raft 集群间rpc协议

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.16.0
// source: rpcProto.proto

package rpcProto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VoteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId  int32 `protobuf:"varint,2,opt,name=candidateId,proto3" json:"candidateId,omitempty"`
	LastLogIndex int64 `protobuf:"varint,3,opt,name=lastLogIndex,proto3" json:"lastLogIndex,omitempty"`
	LastLogTerm  int64 `protobuf:"varint,4,opt,name=lastLogTerm,proto3" json:"lastLogTerm,omitempty"`
}

func (x *VoteReq) Reset() {
	*x = VoteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpcProto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteReq) ProtoMessage() {}

func (x *VoteReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpcProto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteReq.ProtoReflect.Descriptor instead.
func (*VoteReq) Descriptor() ([]byte, []int) {
	return file_rpcProto_proto_rawDescGZIP(), []int{0}
}

func (x *VoteReq) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteReq) GetCandidateId() int32 {
	if x != nil {
		return x.CandidateId
	}
	return 0
}

func (x *VoteReq) GetLastLogIndex() int64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *VoteReq) GetLastLogTerm() int64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

type VoteRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted bool  `protobuf:"varint,2,opt,name=voteGranted,proto3" json:"voteGranted,omitempty"`
}

func (x *VoteRes) Reset() {
	*x = VoteRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpcProto_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRes) ProtoMessage() {}

func (x *VoteRes) ProtoReflect() protoreflect.Message {
	mi := &file_rpcProto_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRes.ProtoReflect.Descriptor instead.
func (*VoteRes) Descriptor() ([]byte, []int) {
	return file_rpcProto_proto_rawDescGZIP(), []int{1}
}

func (x *VoteRes) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteRes) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

type AppendEntriesReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64       `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId     int32       `protobuf:"varint,2,opt,name=leaderId,proto3" json:"leaderId,omitempty"`
	PrevLogTerm  int64       `protobuf:"varint,3,opt,name=prevLogTerm,proto3" json:"prevLogTerm,omitempty"`
	PrevLogIndex int64       `protobuf:"varint,4,opt,name=prevLogIndex,proto3" json:"prevLogIndex,omitempty"`
	LeaderCommit int64       `protobuf:"varint,5,opt,name=leaderCommit,proto3" json:"leaderCommit,omitempty"`
	Entries      []*LogEntry `protobuf:"bytes,6,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *AppendEntriesReq) Reset() {
	*x = AppendEntriesReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpcProto_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesReq) ProtoMessage() {}

func (x *AppendEntriesReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpcProto_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesReq.ProtoReflect.Descriptor instead.
func (*AppendEntriesReq) Descriptor() ([]byte, []int) {
	return file_rpcProto_proto_rawDescGZIP(), []int{2}
}

func (x *AppendEntriesReq) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesReq) GetLeaderId() int32 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *AppendEntriesReq) GetPrevLogTerm() int64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendEntriesReq) GetPrevLogIndex() int64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendEntriesReq) GetLeaderCommit() int64 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

func (x *AppendEntriesReq) GetEntries() []*LogEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type AppendEntriesRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term    int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success bool  `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *AppendEntriesRes) Reset() {
	*x = AppendEntriesRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpcProto_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesRes) ProtoMessage() {}

func (x *AppendEntriesRes) ProtoReflect() protoreflect.Message {
	mi := &file_rpcProto_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesRes.ProtoReflect.Descriptor instead.
func (*AppendEntriesRes) Descriptor() ([]byte, []int) {
	return file_rpcProto_proto_rawDescGZIP(), []int{3}
}

func (x *AppendEntriesRes) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesRes) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term      int64  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Index     int64  `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Key       string `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value     string `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	Operation int64  `protobuf:"varint,5,opt,name=operation,proto3" json:"operation,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpcProto_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_rpcProto_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_rpcProto_proto_rawDescGZIP(), []int{4}
}

func (x *LogEntry) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogEntry) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *LogEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *LogEntry) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *LogEntry) GetOperation() int64 {
	if x != nil {
		return x.Operation
	}
	return 0
}

var File_rpcProto_proto protoreflect.FileDescriptor

var file_rpcProto_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x70, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x72, 0x70, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x01, 0x0a, 0x07, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x61,
	0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c,
	0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65,
	0x72, 0x6d, 0x22, 0x3f, 0x0a, 0x07, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72,
	0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e,
	0x74, 0x65, 0x64, 0x22, 0xda, 0x01, 0x0a, 0x10, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76,
	0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x70,
	0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x72,
	0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x22,
	0x0a, 0x0c, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x12, 0x2c, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x70, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6c,
	0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x22, 0x40, 0x0a, 0x10, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x22, 0x7a, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65,
	0x72, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0x41,
	0x0a, 0x08, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x0b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x70, 0x63, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e, 0x72,
	0x70, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x22,
	0x00, 0x32, 0x56, 0x0a, 0x09, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x49,
	0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12,
	0x1a, 0x2e, 0x72, 0x70, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e,
	0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x1a, 0x2e, 0x72, 0x70,
	0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x3b,
	0x72, 0x70, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpcProto_proto_rawDescOnce sync.Once
	file_rpcProto_proto_rawDescData = file_rpcProto_proto_rawDesc
)

func file_rpcProto_proto_rawDescGZIP() []byte {
	file_rpcProto_proto_rawDescOnce.Do(func() {
		file_rpcProto_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpcProto_proto_rawDescData)
	})
	return file_rpcProto_proto_rawDescData
}

var file_rpcProto_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_rpcProto_proto_goTypes = []interface{}{
	(*VoteReq)(nil),          // 0: rpcProto.VoteReq
	(*VoteRes)(nil),          // 1: rpcProto.VoteRes
	(*AppendEntriesReq)(nil), // 2: rpcProto.AppendEntriesReq
	(*AppendEntriesRes)(nil), // 3: rpcProto.AppendEntriesRes
	(*LogEntry)(nil),         // 4: rpcProto.logEntry
}
var file_rpcProto_proto_depIdxs = []int32{
	4, // 0: rpcProto.AppendEntriesReq.entries:type_name -> rpcProto.logEntry
	0, // 1: rpcProto.election.RequestVote:input_type -> rpcProto.VoteReq
	2, // 2: rpcProto.replicate.AppendEntries:input_type -> rpcProto.AppendEntriesReq
	1, // 3: rpcProto.election.RequestVote:output_type -> rpcProto.VoteRes
	3, // 4: rpcProto.replicate.AppendEntries:output_type -> rpcProto.AppendEntriesRes
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpcProto_proto_init() }
func file_rpcProto_proto_init() {
	if File_rpcProto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpcProto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpcProto_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpcProto_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpcProto_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpcProto_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpcProto_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_rpcProto_proto_goTypes,
		DependencyIndexes: file_rpcProto_proto_depIdxs,
		MessageInfos:      file_rpcProto_proto_msgTypes,
	}.Build()
	File_rpcProto_proto = out.File
	file_rpcProto_proto_rawDesc = nil
	file_rpcProto_proto_goTypes = nil
	file_rpcProto_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ElectionClient is the client API for Election service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ElectionClient interface {
	RequestVote(ctx context.Context, in *VoteReq, opts ...grpc.CallOption) (*VoteRes, error)
}

type electionClient struct {
	cc grpc.ClientConnInterface
}

func NewElectionClient(cc grpc.ClientConnInterface) ElectionClient {
	return &electionClient{cc}
}

func (c *electionClient) RequestVote(ctx context.Context, in *VoteReq, opts ...grpc.CallOption) (*VoteRes, error) {
	out := new(VoteRes)
	err := c.cc.Invoke(ctx, "/rpcProto.election/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ElectionServer is the server API for Election service.
type ElectionServer interface {
	RequestVote(context.Context, *VoteReq) (*VoteRes, error)
}

// UnimplementedElectionServer can be embedded to have forward compatible implementations.
type UnimplementedElectionServer struct {
}

func (*UnimplementedElectionServer) RequestVote(context.Context, *VoteReq) (*VoteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}

func RegisterElectionServer(s *grpc.Server, srv ElectionServer) {
	s.RegisterService(&_Election_serviceDesc, srv)
}

func _Election_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ElectionServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcProto.election/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ElectionServer).RequestVote(ctx, req.(*VoteReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Election_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpcProto.election",
	HandlerType: (*ElectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVote",
			Handler:    _Election_RequestVote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcProto.proto",
}

// ReplicateClient is the client API for Replicate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReplicateClient interface {
	AppendEntries(ctx context.Context, in *AppendEntriesReq, opts ...grpc.CallOption) (*AppendEntriesRes, error)
}

type replicateClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicateClient(cc grpc.ClientConnInterface) ReplicateClient {
	return &replicateClient{cc}
}

func (c *replicateClient) AppendEntries(ctx context.Context, in *AppendEntriesReq, opts ...grpc.CallOption) (*AppendEntriesRes, error) {
	out := new(AppendEntriesRes)
	err := c.cc.Invoke(ctx, "/rpcProto.replicate/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicateServer is the server API for Replicate service.
type ReplicateServer interface {
	AppendEntries(context.Context, *AppendEntriesReq) (*AppendEntriesRes, error)
}

// UnimplementedReplicateServer can be embedded to have forward compatible implementations.
type UnimplementedReplicateServer struct {
}

func (*UnimplementedReplicateServer) AppendEntries(context.Context, *AppendEntriesReq) (*AppendEntriesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}

func RegisterReplicateServer(s *grpc.Server, srv ReplicateServer) {
	s.RegisterService(&_Replicate_serviceDesc, srv)
}

func _Replicate_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicateServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcProto.replicate/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicateServer).AppendEntries(ctx, req.(*AppendEntriesReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Replicate_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpcProto.replicate",
	HandlerType: (*ReplicateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _Replicate_AppendEntries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcProto.proto",
}
