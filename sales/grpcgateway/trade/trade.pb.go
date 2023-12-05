// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: trade/trade.proto

package trade

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The request message for creating a sale.
type CreateSaleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LineItems      []*LineItem `protobuf:"bytes,1,rep,name=line_items,json=lineItems,proto3" json:"line_items,omitempty"`
	DiscountAmount float32     `protobuf:"fixed32,2,opt,name=discountAmount,proto3" json:"discountAmount,omitempty"` // Flat discount amount on the total sale
}

func (x *CreateSaleRequest) Reset() {
	*x = CreateSaleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_trade_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSaleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSaleRequest) ProtoMessage() {}

func (x *CreateSaleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trade_trade_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSaleRequest.ProtoReflect.Descriptor instead.
func (*CreateSaleRequest) Descriptor() ([]byte, []int) {
	return file_trade_trade_proto_rawDescGZIP(), []int{0}
}

func (x *CreateSaleRequest) GetLineItems() []*LineItem {
	if x != nil {
		return x.LineItems
	}
	return nil
}

func (x *CreateSaleRequest) GetDiscountAmount() float32 {
	if x != nil {
		return x.DiscountAmount
	}
	return 0
}

// Represents an item in a sale.
type LineItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId string `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity  int32  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *LineItem) Reset() {
	*x = LineItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_trade_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LineItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LineItem) ProtoMessage() {}

func (x *LineItem) ProtoReflect() protoreflect.Message {
	mi := &file_trade_trade_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LineItem.ProtoReflect.Descriptor instead.
func (*LineItem) Descriptor() ([]byte, []int) {
	return file_trade_trade_proto_rawDescGZIP(), []int{1}
}

func (x *LineItem) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *LineItem) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

// The response message for a sale creation.
type CreateSaleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SaleId     string                 `protobuf:"bytes,1,opt,name=sale_id,json=saleId,proto3" json:"sale_id,omitempty"`
	LineItems  []*LineItem            `protobuf:"bytes,2,rep,name=line_items,json=lineItems,proto3" json:"line_items,omitempty"`
	TotalPrice *wrapperspb.FloatValue `protobuf:"bytes,3,opt,name=total_price,json=totalPrice,proto3" json:"total_price,omitempty"`
}

func (x *CreateSaleResponse) Reset() {
	*x = CreateSaleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_trade_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSaleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSaleResponse) ProtoMessage() {}

func (x *CreateSaleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trade_trade_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSaleResponse.ProtoReflect.Descriptor instead.
func (*CreateSaleResponse) Descriptor() ([]byte, []int) {
	return file_trade_trade_proto_rawDescGZIP(), []int{2}
}

func (x *CreateSaleResponse) GetSaleId() string {
	if x != nil {
		return x.SaleId
	}
	return ""
}

func (x *CreateSaleResponse) GetLineItems() []*LineItem {
	if x != nil {
		return x.LineItems
	}
	return nil
}

func (x *CreateSaleResponse) GetTotalPrice() *wrapperspb.FloatValue {
	if x != nil {
		return x.TotalPrice
	}
	return nil
}

var File_trade_trade_proto protoreflect.FileDescriptor

var file_trade_trade_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x72, 0x61, 0x64, 0x65, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x72, 0x61, 0x64, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6b, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x61, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a,
	0x0a, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x65, 0x2e, 0x4c, 0x69, 0x6e, 0x65, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x09, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x26, 0x0a,
	0x0e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x45, 0x0a, 0x08, 0x4c, 0x69, 0x6e, 0x65, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x9b, 0x01, 0x0a,
	0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x61, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x61, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x0a,
	0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x65, 0x2e, 0x4c, 0x69, 0x6e, 0x65, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x09, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x3c, 0x0a, 0x0b,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x32, 0x67, 0x0a, 0x0c, 0x53, 0x61,
	0x6c, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x0a, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x65, 0x12, 0x18, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x53, 0x61, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x3a, 0x01, 0x2a, 0x22, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x61,
	0x6c, 0x65, 0x73, 0x42, 0x08, 0x5a, 0x06, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trade_trade_proto_rawDescOnce sync.Once
	file_trade_trade_proto_rawDescData = file_trade_trade_proto_rawDesc
)

func file_trade_trade_proto_rawDescGZIP() []byte {
	file_trade_trade_proto_rawDescOnce.Do(func() {
		file_trade_trade_proto_rawDescData = protoimpl.X.CompressGZIP(file_trade_trade_proto_rawDescData)
	})
	return file_trade_trade_proto_rawDescData
}

var file_trade_trade_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_trade_trade_proto_goTypes = []interface{}{
	(*CreateSaleRequest)(nil),     // 0: trade.CreateSaleRequest
	(*LineItem)(nil),              // 1: trade.LineItem
	(*CreateSaleResponse)(nil),    // 2: trade.CreateSaleResponse
	(*wrapperspb.FloatValue)(nil), // 3: google.protobuf.FloatValue
}
var file_trade_trade_proto_depIdxs = []int32{
	1, // 0: trade.CreateSaleRequest.line_items:type_name -> trade.LineItem
	1, // 1: trade.CreateSaleResponse.line_items:type_name -> trade.LineItem
	3, // 2: trade.CreateSaleResponse.total_price:type_name -> google.protobuf.FloatValue
	0, // 3: trade.SalesService.CreateSale:input_type -> trade.CreateSaleRequest
	2, // 4: trade.SalesService.CreateSale:output_type -> trade.CreateSaleResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_trade_trade_proto_init() }
func file_trade_trade_proto_init() {
	if File_trade_trade_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trade_trade_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSaleRequest); i {
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
		file_trade_trade_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LineItem); i {
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
		file_trade_trade_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSaleResponse); i {
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
			RawDescriptor: file_trade_trade_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trade_trade_proto_goTypes,
		DependencyIndexes: file_trade_trade_proto_depIdxs,
		MessageInfos:      file_trade_trade_proto_msgTypes,
	}.Build()
	File_trade_trade_proto = out.File
	file_trade_trade_proto_rawDesc = nil
	file_trade_trade_proto_goTypes = nil
	file_trade_trade_proto_depIdxs = nil
}
