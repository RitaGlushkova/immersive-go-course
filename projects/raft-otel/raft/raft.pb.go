// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: raft/raft.proto

package raft

import (
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

// The request message from a candidate
type CandidateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId  int64 `protobuf:"varint,2,opt,name=candidateId,proto3" json:"candidateId,omitempty"`
	LastLogIndex int64 `protobuf:"varint,3,opt,name=lastLogIndex,proto3" json:"lastLogIndex,omitempty"`
	LastLogTerm  int64 `protobuf:"varint,4,opt,name=lastLogTerm,proto3" json:"lastLogTerm,omitempty"`
}

func (x *CandidateRequest) Reset() {
	*x = CandidateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_raft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CandidateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CandidateRequest) ProtoMessage() {}

func (x *CandidateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_raft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CandidateRequest.ProtoReflect.Descriptor instead.
func (*CandidateRequest) Descriptor() ([]byte, []int) {
	return file_raft_raft_proto_rawDescGZIP(), []int{0}
}

func (x *CandidateRequest) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *CandidateRequest) GetCandidateId() int64 {
	if x != nil {
		return x.CandidateId
	}
	return 0
}

func (x *CandidateRequest) GetLastLogIndex() int64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *CandidateRequest) GetLastLogTerm() int64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

// Candidate received vote results
type VoteResults struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted bool  `protobuf:"varint,2,opt,name=voteGranted,proto3" json:"voteGranted,omitempty"`
}

func (x *VoteResults) Reset() {
	*x = VoteResults{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_raft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResults) ProtoMessage() {}

func (x *VoteResults) ProtoReflect() protoreflect.Message {
	mi := &file_raft_raft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResults.ProtoReflect.Descriptor instead.
func (*VoteResults) Descriptor() ([]byte, []int) {
	return file_raft_raft_proto_rawDescGZIP(), []int{1}
}

func (x *VoteResults) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteResults) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

// Request to Append Entries from leader to followers
type RequestAppend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId     int64  `protobuf:"varint,2,opt,name=leaderId,proto3" json:"leaderId,omitempty"`
	PrevLogIndex int64  `protobuf:"varint,3,opt,name=prevLogIndex,proto3" json:"prevLogIndex,omitempty"`
	Entry        *Entry `protobuf:"bytes,4,opt,name=entry,proto3" json:"entry,omitempty"`
	LeaderCommit int64  `protobuf:"varint,5,opt,name=leaderCommit,proto3" json:"leaderCommit,omitempty"`
}

func (x *RequestAppend) Reset() {
	*x = RequestAppend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_raft_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestAppend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestAppend) ProtoMessage() {}

func (x *RequestAppend) ProtoReflect() protoreflect.Message {
	mi := &file_raft_raft_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestAppend.ProtoReflect.Descriptor instead.
func (*RequestAppend) Descriptor() ([]byte, []int) {
	return file_raft_raft_proto_rawDescGZIP(), []int{2}
}

func (x *RequestAppend) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestAppend) GetLeaderId() int64 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *RequestAppend) GetPrevLogIndex() int64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *RequestAppend) GetEntry() *Entry {
	if x != nil {
		return x.Entry
	}
	return nil
}

func (x *RequestAppend) GetLeaderCommit() int64 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

// Result of trying to Append Entries from leader to followers
type ResultAppend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term     int64  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Succeeds bool   `protobuf:"varint,2,opt,name=succeeds,proto3" json:"succeeds,omitempty"`
	Entry    *Entry `protobuf:"bytes,3,opt,name=entry,proto3" json:"entry,omitempty"`
}

func (x *ResultAppend) Reset() {
	*x = ResultAppend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_raft_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResultAppend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResultAppend) ProtoMessage() {}

func (x *ResultAppend) ProtoReflect() protoreflect.Message {
	mi := &file_raft_raft_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResultAppend.ProtoReflect.Descriptor instead.
func (*ResultAppend) Descriptor() ([]byte, []int) {
	return file_raft_raft_proto_rawDescGZIP(), []int{3}
}

func (x *ResultAppend) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *ResultAppend) GetSucceeds() bool {
	if x != nil {
		return x.Succeeds
	}
	return false
}

func (x *ResultAppend) GetEntry() *Entry {
	if x != nil {
		return x.Entry
	}
	return nil
}

// Entry from Client
type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_raft_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_raft_raft_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_raft_raft_proto_rawDescGZIP(), []int{4}
}

func (x *Entry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Entry) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_raft_raft_proto protoreflect.FileDescriptor

var file_raft_raft_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x72, 0x61, 0x66, 0x74, 0x22, 0x8e, 0x01, 0x0a, 0x10, 0x43, 0x61, 0x6e, 0x64,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d,
	0x12, 0x20, 0x0a, 0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6c, 0x61, 0x73,
	0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x43, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x76,
	0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x22, 0xaa, 0x01,
	0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x22, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x21, 0x0a, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x22, 0x61, 0x0a, 0x0c, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x73, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x73, 0x12, 0x21, 0x0a, 0x05, 0x65, 0x6e,
	0x74, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x72, 0x61, 0x66, 0x74,
	0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x2f, 0x0a,
	0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x7c,
	0x0a, 0x04, 0x52, 0x61, 0x66, 0x74, 0x12, 0x38, 0x0a, 0x0b, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x13, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x1a, 0x12, 0x2e, 0x72, 0x61, 0x66,
	0x74, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x22, 0x00,
	0x12, 0x3a, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12,
	0x16, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22, 0x00, 0x42, 0x29, 0x5a, 0x27,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x52, 0x69, 0x74, 0x61, 0x47,
	0x6c, 0x75, 0x73, 0x68, 0x6b, 0x6f, 0x76, 0x61, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2d, 0x6f, 0x74,
	0x65, 0x6c, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_raft_raft_proto_rawDescOnce sync.Once
	file_raft_raft_proto_rawDescData = file_raft_raft_proto_rawDesc
)

func file_raft_raft_proto_rawDescGZIP() []byte {
	file_raft_raft_proto_rawDescOnce.Do(func() {
		file_raft_raft_proto_rawDescData = protoimpl.X.CompressGZIP(file_raft_raft_proto_rawDescData)
	})
	return file_raft_raft_proto_rawDescData
}

var file_raft_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_raft_raft_proto_goTypes = []interface{}{
	(*CandidateRequest)(nil), // 0: raft.CandidateRequest
	(*VoteResults)(nil),      // 1: raft.VoteResults
	(*RequestAppend)(nil),    // 2: raft.RequestAppend
	(*ResultAppend)(nil),     // 3: raft.ResultAppend
	(*Entry)(nil),            // 4: raft.Entry
}
var file_raft_raft_proto_depIdxs = []int32{
	4, // 0: raft.RequestAppend.entry:type_name -> raft.Entry
	4, // 1: raft.ResultAppend.entry:type_name -> raft.Entry
	2, // 2: raft.Raft.AppendEntry:input_type -> raft.RequestAppend
	0, // 3: raft.Raft.RequestVote:input_type -> raft.CandidateRequest
	3, // 4: raft.Raft.AppendEntry:output_type -> raft.ResultAppend
	1, // 5: raft.Raft.RequestVote:output_type -> raft.VoteResults
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_raft_raft_proto_init() }
func file_raft_raft_proto_init() {
	if File_raft_raft_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_raft_raft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CandidateRequest); i {
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
		file_raft_raft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResults); i {
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
		file_raft_raft_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestAppend); i {
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
		file_raft_raft_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResultAppend); i {
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
		file_raft_raft_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
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
			RawDescriptor: file_raft_raft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_raft_raft_proto_goTypes,
		DependencyIndexes: file_raft_raft_proto_depIdxs,
		MessageInfos:      file_raft_raft_proto_msgTypes,
	}.Build()
	File_raft_raft_proto = out.File
	file_raft_raft_proto_rawDesc = nil
	file_raft_raft_proto_goTypes = nil
	file_raft_raft_proto_depIdxs = nil
}
