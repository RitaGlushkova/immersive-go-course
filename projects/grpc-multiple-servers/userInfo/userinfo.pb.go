// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: userInfo/userInfo.proto

package userinfo

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

// The request message containing the user's name and date of birth.
type UserInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Dob  string `protobuf:"bytes,2,opt,name=dob,proto3" json:"dob,omitempty"`
}

func (x *UserInfoRequest) Reset() {
	*x = UserInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userInfo_userInfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfoRequest) ProtoMessage() {}

func (x *UserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userInfo_userInfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfoRequest.ProtoReflect.Descriptor instead.
func (*UserInfoRequest) Descriptor() ([]byte, []int) {
	return file_userInfo_userInfo_proto_rawDescGZIP(), []int{0}
}

func (x *UserInfoRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserInfoRequest) GetDob() string {
	if x != nil {
		return x.Dob
	}
	return ""
}

// The response message containing user information
type UserInfoReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Dob        string `protobuf:"bytes,2,opt,name=dob,proto3" json:"dob,omitempty"`
	Email      string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	College    string `protobuf:"bytes,4,opt,name=college,proto3" json:"college,omitempty"`
	Occupation string `protobuf:"bytes,5,opt,name=occupation,proto3" json:"occupation,omitempty"`
	Age        int32  `protobuf:"varint,6,opt,name=age,proto3" json:"age,omitempty"`
	Redirect   string `protobuf:"bytes,7,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Notfound   string `protobuf:"bytes,8,opt,name=notfound,proto3" json:"notfound,omitempty"`
}

func (x *UserInfoReply) Reset() {
	*x = UserInfoReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userInfo_userInfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInfoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfoReply) ProtoMessage() {}

func (x *UserInfoReply) ProtoReflect() protoreflect.Message {
	mi := &file_userInfo_userInfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfoReply.ProtoReflect.Descriptor instead.
func (*UserInfoReply) Descriptor() ([]byte, []int) {
	return file_userInfo_userInfo_proto_rawDescGZIP(), []int{1}
}

func (x *UserInfoReply) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserInfoReply) GetDob() string {
	if x != nil {
		return x.Dob
	}
	return ""
}

func (x *UserInfoReply) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserInfoReply) GetCollege() string {
	if x != nil {
		return x.College
	}
	return ""
}

func (x *UserInfoReply) GetOccupation() string {
	if x != nil {
		return x.Occupation
	}
	return ""
}

func (x *UserInfoReply) GetAge() int32 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *UserInfoReply) GetRedirect() string {
	if x != nil {
		return x.Redirect
	}
	return ""
}

func (x *UserInfoReply) GetNotfound() string {
	if x != nil {
		return x.Notfound
	}
	return ""
}

var File_userInfo_userInfo_proto protoreflect.FileDescriptor

var file_userInfo_userInfo_proto_rawDesc = []byte{
	0x0a, 0x17, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x22, 0x37, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f,
	0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x6f, 0x62, 0x22, 0xcf, 0x01, 0x0a,
	0x0d, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x64, 0x6f, 0x62, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f,
	0x6c, 0x6c, 0x65, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x6f, 0x74, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x74, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x32, 0x50,
	0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x44, 0x0a, 0x0c, 0x53, 0x65,
	0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x42, 0x56, 0x5a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x52,
	0x69, 0x74, 0x61, 0x47, 0x6c, 0x75, 0x73, 0x68, 0x6b, 0x6f, 0x76, 0x61, 0x2f, 0x69, 0x6d, 0x6d,
	0x65, 0x72, 0x73, 0x69, 0x76, 0x65, 0x2d, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x75, 0x72, 0x73, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x6d,
	0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_userInfo_userInfo_proto_rawDescOnce sync.Once
	file_userInfo_userInfo_proto_rawDescData = file_userInfo_userInfo_proto_rawDesc
)

func file_userInfo_userInfo_proto_rawDescGZIP() []byte {
	file_userInfo_userInfo_proto_rawDescOnce.Do(func() {
		file_userInfo_userInfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_userInfo_userInfo_proto_rawDescData)
	})
	return file_userInfo_userInfo_proto_rawDescData
}

var file_userInfo_userInfo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_userInfo_userInfo_proto_goTypes = []interface{}{
	(*UserInfoRequest)(nil), // 0: protobuf.UserInfoRequest
	(*UserInfoReply)(nil),   // 1: protobuf.UserInfoReply
}
var file_userInfo_userInfo_proto_depIdxs = []int32{
	0, // 0: protobuf.UserInfo.SendUserInfo:input_type -> protobuf.UserInfoRequest
	1, // 1: protobuf.UserInfo.SendUserInfo:output_type -> protobuf.UserInfoReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_userInfo_userInfo_proto_init() }
func file_userInfo_userInfo_proto_init() {
	if File_userInfo_userInfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_userInfo_userInfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInfoRequest); i {
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
		file_userInfo_userInfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInfoReply); i {
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
			RawDescriptor: file_userInfo_userInfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_userInfo_userInfo_proto_goTypes,
		DependencyIndexes: file_userInfo_userInfo_proto_depIdxs,
		MessageInfos:      file_userInfo_userInfo_proto_msgTypes,
	}.Build()
	File_userInfo_userInfo_proto = out.File
	file_userInfo_userInfo_proto_rawDesc = nil
	file_userInfo_userInfo_proto_goTypes = nil
	file_userInfo_userInfo_proto_depIdxs = nil
}
