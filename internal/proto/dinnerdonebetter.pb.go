// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: dinnerdonebetter.proto

package proto

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_dinnerdonebetter_proto protoreflect.FileDescriptor

var file_dinnerdonebetter_proto_rawDesc = []byte{
	0x0a, 0x16, 0x64, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x64, 0x6f, 0x6e, 0x65, 0x62, 0x65, 0x74, 0x74,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x64, 0x69, 0x6e, 0x6e, 0x65, 0x72,
	0x64, 0x6f, 0x6e, 0x65, 0x62, 0x65, 0x74, 0x74, 0x65, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x69,
	0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32,
	0x6b, 0x0a, 0x10, 0x44, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x44, 0x6f, 0x6e, 0x65, 0x42, 0x65, 0x74,
	0x74, 0x65, 0x72, 0x12, 0x57, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x64, 0x69, 0x6e, 0x6e, 0x65, 0x72,
	0x64, 0x6f, 0x6e, 0x65, 0x62, 0x65, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x42, 0x34, 0x5a, 0x32,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x69, 0x6e, 0x6e, 0x65,
	0x72, 0x64, 0x6f, 0x6e, 0x65, 0x62, 0x65, 0x74, 0x74, 0x65, 0x72, 0x2f, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_dinnerdonebetter_proto_goTypes = []interface{}{
	(*empty.Empty)(nil),     // 0: google.protobuf.Empty
	(*ValidIngredient)(nil), // 1: dinnerdonebetter.ValidIngredient
}
var file_dinnerdonebetter_proto_depIdxs = []int32{
	0, // 0: dinnerdonebetter.DinnerDoneBetter.GetRandomValidIngredient:input_type -> google.protobuf.Empty
	1, // 1: dinnerdonebetter.DinnerDoneBetter.GetRandomValidIngredient:output_type -> dinnerdonebetter.ValidIngredient
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_dinnerdonebetter_proto_init() }
func file_dinnerdonebetter_proto_init() {
	if File_dinnerdonebetter_proto != nil {
		return
	}
	file_valid_ingredient_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dinnerdonebetter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dinnerdonebetter_proto_goTypes,
		DependencyIndexes: file_dinnerdonebetter_proto_depIdxs,
	}.Build()
	File_dinnerdonebetter_proto = out.File
	file_dinnerdonebetter_proto_rawDesc = nil
	file_dinnerdonebetter_proto_goTypes = nil
	file_dinnerdonebetter_proto_depIdxs = nil
}