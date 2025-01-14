// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v5.29.0--rc2
// source: trim.proto

package pb

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

type TrimRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	S             string                 `protobuf:"bytes,1,opt,name=s,proto3" json:"s,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TrimRequest) Reset() {
	*x = TrimRequest{}
	mi := &file_trim_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TrimRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrimRequest) ProtoMessage() {}

func (x *TrimRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trim_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrimRequest.ProtoReflect.Descriptor instead.
func (*TrimRequest) Descriptor() ([]byte, []int) {
	return file_trim_proto_rawDescGZIP(), []int{0}
}

func (x *TrimRequest) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

type TrimResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	S             string                 `protobuf:"bytes,1,opt,name=s,proto3" json:"s,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TrimResponse) Reset() {
	*x = TrimResponse{}
	mi := &file_trim_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TrimResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrimResponse) ProtoMessage() {}

func (x *TrimResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trim_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrimResponse.ProtoReflect.Descriptor instead.
func (*TrimResponse) Descriptor() ([]byte, []int) {
	return file_trim_proto_rawDescGZIP(), []int{1}
}

func (x *TrimResponse) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

var File_trim_proto protoreflect.FileDescriptor

var file_trim_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x72, 0x69, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62,
	0x22, 0x1b, 0x0a, 0x0b, 0x54, 0x72, 0x69, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0c, 0x0a, 0x01, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73, 0x22, 0x1c, 0x0a,
	0x0c, 0x54, 0x72, 0x69, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0c, 0x0a,
	0x01, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73, 0x32, 0x38, 0x0a, 0x04, 0x54,
	0x72, 0x69, 0x6d, 0x12, 0x30, 0x0a, 0x09, 0x54, 0x72, 0x69, 0x6d, 0x53, 0x70, 0x61, 0x63, 0x65,
	0x12, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x69, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x69, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x61, 0x64, 0x64, 0x73, 0x72, 0x76, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trim_proto_rawDescOnce sync.Once
	file_trim_proto_rawDescData = file_trim_proto_rawDesc
)

func file_trim_proto_rawDescGZIP() []byte {
	file_trim_proto_rawDescOnce.Do(func() {
		file_trim_proto_rawDescData = protoimpl.X.CompressGZIP(file_trim_proto_rawDescData)
	})
	return file_trim_proto_rawDescData
}

var file_trim_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_trim_proto_goTypes = []any{
	(*TrimRequest)(nil),  // 0: pb.TrimRequest
	(*TrimResponse)(nil), // 1: pb.TrimResponse
}
var file_trim_proto_depIdxs = []int32{
	0, // 0: pb.Trim.TrimSpace:input_type -> pb.TrimRequest
	1, // 1: pb.Trim.TrimSpace:output_type -> pb.TrimResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_trim_proto_init() }
func file_trim_proto_init() {
	if File_trim_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_trim_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trim_proto_goTypes,
		DependencyIndexes: file_trim_proto_depIdxs,
		MessageInfos:      file_trim_proto_msgTypes,
	}.Build()
	File_trim_proto = out.File
	file_trim_proto_rawDesc = nil
	file_trim_proto_goTypes = nil
	file_trim_proto_depIdxs = nil
}
