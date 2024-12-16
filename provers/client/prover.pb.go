// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: prover/v1/prover.proto

package client

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

type ProveStateTransitionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProveStateTransitionRequest) Reset() {
	*x = ProveStateTransitionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prover_v1_prover_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProveStateTransitionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProveStateTransitionRequest) ProtoMessage() {}

func (x *ProveStateTransitionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_prover_v1_prover_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProveStateTransitionRequest.ProtoReflect.Descriptor instead.
func (*ProveStateTransitionRequest) Descriptor() ([]byte, []int) {
	return file_prover_v1_prover_proto_rawDescGZIP(), []int{0}
}

type ProveStateTransitionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proof        []byte `protobuf:"bytes,1,opt,name=proof,proto3" json:"proof,omitempty"`
	PublicValues []byte `protobuf:"bytes,2,opt,name=public_values,json=publicValues,proto3" json:"public_values,omitempty"`
}

func (x *ProveStateTransitionResponse) Reset() {
	*x = ProveStateTransitionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prover_v1_prover_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProveStateTransitionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProveStateTransitionResponse) ProtoMessage() {}

func (x *ProveStateTransitionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_prover_v1_prover_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProveStateTransitionResponse.ProtoReflect.Descriptor instead.
func (*ProveStateTransitionResponse) Descriptor() ([]byte, []int) {
	return file_prover_v1_prover_proto_rawDescGZIP(), []int{1}
}

func (x *ProveStateTransitionResponse) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *ProveStateTransitionResponse) GetPublicValues() []byte {
	if x != nil {
		return x.PublicValues
	}
	return nil
}

type KeyValuePair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KeyValuePair) Reset() {
	*x = KeyValuePair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prover_v1_prover_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValuePair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValuePair) ProtoMessage() {}

func (x *KeyValuePair) ProtoReflect() protoreflect.Message {
	mi := &file_prover_v1_prover_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValuePair.ProtoReflect.Descriptor instead.
func (*KeyValuePair) Descriptor() ([]byte, []int) {
	return file_prover_v1_prover_proto_rawDescGZIP(), []int{2}
}

func (x *KeyValuePair) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KeyValuePair) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type ProveMembershipRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height        int64           `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	KeyValuePairs []*KeyValuePair `protobuf:"bytes,2,rep,name=key_value_pairs,json=keyValuePairs,proto3" json:"key_value_pairs,omitempty"`
}

func (x *ProveMembershipRequest) Reset() {
	*x = ProveMembershipRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prover_v1_prover_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProveMembershipRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProveMembershipRequest) ProtoMessage() {}

func (x *ProveMembershipRequest) ProtoReflect() protoreflect.Message {
	mi := &file_prover_v1_prover_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProveMembershipRequest.ProtoReflect.Descriptor instead.
func (*ProveMembershipRequest) Descriptor() ([]byte, []int) {
	return file_prover_v1_prover_proto_rawDescGZIP(), []int{3}
}

func (x *ProveMembershipRequest) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *ProveMembershipRequest) GetKeyValuePairs() []*KeyValuePair {
	if x != nil {
		return x.KeyValuePairs
	}
	return nil
}

type ProveMembershipResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proof  []byte `protobuf:"bytes,1,opt,name=proof,proto3" json:"proof,omitempty"`
	Height int64  `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *ProveMembershipResponse) Reset() {
	*x = ProveMembershipResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prover_v1_prover_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProveMembershipResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProveMembershipResponse) ProtoMessage() {}

func (x *ProveMembershipResponse) ProtoReflect() protoreflect.Message {
	mi := &file_prover_v1_prover_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProveMembershipResponse.ProtoReflect.Descriptor instead.
func (*ProveMembershipResponse) Descriptor() ([]byte, []int) {
	return file_prover_v1_prover_proto_rawDescGZIP(), []int{4}
}

func (x *ProveMembershipResponse) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *ProveMembershipResponse) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

var File_prover_v1_prover_proto protoreflect.FileDescriptor

var file_prover_v1_prover_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x76,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x63, 0x65, 0x6c, 0x65, 0x73, 0x74,
	0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x1d, 0x0a, 0x1b,
	0x50, 0x72, 0x6f, 0x76, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x59, 0x0a, 0x1c, 0x50,
	0x72, 0x6f, 0x76, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70,
	0x72, 0x6f, 0x6f, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f,
	0x66, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x36, 0x0a, 0x0c, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x50, 0x61, 0x69, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x7a,
	0x0a, 0x16, 0x50, 0x72, 0x6f, 0x76, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x48, 0x0a, 0x0f, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x70, 0x61,
	0x69, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x65, 0x6c, 0x65,
	0x73, 0x74, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4b,
	0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x50, 0x61, 0x69, 0x72, 0x52, 0x0d, 0x6b, 0x65, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x50, 0x61, 0x69, 0x72, 0x73, 0x22, 0x47, 0x0a, 0x17, 0x50, 0x72,
	0x6f, 0x76, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x16, 0x0a, 0x06, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x68, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x32, 0xef, 0x01, 0x0a, 0x06, 0x50, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x79,
	0x0a, 0x14, 0x50, 0x72, 0x6f, 0x76, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x2e, 0x63, 0x65, 0x6c, 0x65, 0x73, 0x74, 0x69,
	0x61, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x76,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x63, 0x65, 0x6c, 0x65, 0x73, 0x74,
	0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f,
	0x76, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6a, 0x0a, 0x0f, 0x50, 0x72, 0x6f,
	0x76, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x12, 0x2a, 0x2e, 0x63,
	0x65, 0x6c, 0x65, 0x73, 0x74, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x63, 0x65, 0x6c, 0x65, 0x73,
	0x74, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72,
	0x6f, 0x76, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x73,
	0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_prover_v1_prover_proto_rawDescOnce sync.Once
	file_prover_v1_prover_proto_rawDescData = file_prover_v1_prover_proto_rawDesc
)

func file_prover_v1_prover_proto_rawDescGZIP() []byte {
	file_prover_v1_prover_proto_rawDescOnce.Do(func() {
		file_prover_v1_prover_proto_rawDescData = protoimpl.X.CompressGZIP(file_prover_v1_prover_proto_rawDescData)
	})
	return file_prover_v1_prover_proto_rawDescData
}

var file_prover_v1_prover_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_prover_v1_prover_proto_goTypes = []any{
	(*ProveStateTransitionRequest)(nil),  // 0: celestia.prover.v1.ProveStateTransitionRequest
	(*ProveStateTransitionResponse)(nil), // 1: celestia.prover.v1.ProveStateTransitionResponse
	(*KeyValuePair)(nil),                 // 2: celestia.prover.v1.KeyValuePair
	(*ProveMembershipRequest)(nil),       // 3: celestia.prover.v1.ProveMembershipRequest
	(*ProveMembershipResponse)(nil),      // 4: celestia.prover.v1.ProveMembershipResponse
}
var file_prover_v1_prover_proto_depIdxs = []int32{
	2, // 0: celestia.prover.v1.ProveMembershipRequest.key_value_pairs:type_name -> celestia.prover.v1.KeyValuePair
	0, // 1: celestia.prover.v1.Prover.ProveStateTransition:input_type -> celestia.prover.v1.ProveStateTransitionRequest
	3, // 2: celestia.prover.v1.Prover.ProveMembership:input_type -> celestia.prover.v1.ProveMembershipRequest
	1, // 3: celestia.prover.v1.Prover.ProveStateTransition:output_type -> celestia.prover.v1.ProveStateTransitionResponse
	4, // 4: celestia.prover.v1.Prover.ProveMembership:output_type -> celestia.prover.v1.ProveMembershipResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_prover_v1_prover_proto_init() }
func file_prover_v1_prover_proto_init() {
	if File_prover_v1_prover_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_prover_v1_prover_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ProveStateTransitionRequest); i {
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
		file_prover_v1_prover_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProveStateTransitionResponse); i {
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
		file_prover_v1_prover_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*KeyValuePair); i {
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
		file_prover_v1_prover_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ProveMembershipRequest); i {
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
		file_prover_v1_prover_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ProveMembershipResponse); i {
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
			RawDescriptor: file_prover_v1_prover_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_prover_v1_prover_proto_goTypes,
		DependencyIndexes: file_prover_v1_prover_proto_depIdxs,
		MessageInfos:      file_prover_v1_prover_proto_msgTypes,
	}.Build()
	File_prover_v1_prover_proto = out.File
	file_prover_v1_prover_proto_rawDesc = nil
	file_prover_v1_prover_proto_goTypes = nil
	file_prover_v1_prover_proto_depIdxs = nil
}
