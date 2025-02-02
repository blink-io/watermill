// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.24.4
// source: internal/api/events.proto

package cqrs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TestProtobufEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	When *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=when,proto3" json:"when,omitempty"`
}

func (x *TestProtobufEvent) Reset() {
	*x = TestProtobufEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestProtobufEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestProtobufEvent) ProtoMessage() {}

func (x *TestProtobufEvent) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestProtobufEvent.ProtoReflect.Descriptor instead.
func (*TestProtobufEvent) Descriptor() ([]byte, []int) {
	return file_internal_api_events_proto_rawDescGZIP(), []int{0}
}

func (x *TestProtobufEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TestProtobufEvent) GetWhen() *timestamppb.Timestamp {
	if x != nil {
		return x.When
	}
	return nil
}

var File_internal_api_events_proto protoreflect.FileDescriptor

var file_internal_api_events_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x53, 0x0a, 0x11,
	0x54, 0x65, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x2e, 0x0a, 0x04, 0x77, 0x68, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x77, 0x68, 0x65,
	0x6e, 0x42, 0x26, 0x5a, 0x24, 0x62, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x66, 0x75, 0x6e, 0x2f, 0x69,
	0x64, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_internal_api_events_proto_rawDescOnce sync.Once
	file_internal_api_events_proto_rawDescData = file_internal_api_events_proto_rawDesc
)

func file_internal_api_events_proto_rawDescGZIP() []byte {
	file_internal_api_events_proto_rawDescOnce.Do(func() {
		file_internal_api_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_api_events_proto_rawDescData)
	})
	return file_internal_api_events_proto_rawDescData
}

var file_internal_api_events_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_api_events_proto_goTypes = []interface{}{
	(*TestProtobufEvent)(nil),     // 0: TestProtobufEvent
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_internal_api_events_proto_depIdxs = []int32{
	1, // 0: TestProtobufEvent.when:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_api_events_proto_init() }
func file_internal_api_events_proto_init() {
	if File_internal_api_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_api_events_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestProtobufEvent); i {
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
			RawDescriptor: file_internal_api_events_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_api_events_proto_goTypes,
		DependencyIndexes: file_internal_api_events_proto_depIdxs,
		MessageInfos:      file_internal_api_events_proto_msgTypes,
	}.Build()
	File_internal_api_events_proto = out.File
	file_internal_api_events_proto_rawDesc = nil
	file_internal_api_events_proto_goTypes = nil
	file_internal_api_events_proto_depIdxs = nil
}
