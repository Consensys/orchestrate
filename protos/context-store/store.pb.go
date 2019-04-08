// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/context-store/store.proto

package context_store

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	envelope "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StoreRequest struct {
	Envelope             *envelope.Envelope `protobuf:"bytes,1,opt,name=envelope,proto3" json:"envelope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *StoreRequest) Reset()         { *m = StoreRequest{} }
func (m *StoreRequest) String() string { return proto.CompactTextString(m) }
func (*StoreRequest) ProtoMessage()    {}
func (*StoreRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_79926f01850c7b5c, []int{0}
}

func (m *StoreRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StoreRequest.Unmarshal(m, b)
}
func (m *StoreRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StoreRequest.Marshal(b, m, deterministic)
}
func (m *StoreRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StoreRequest.Merge(m, src)
}
func (m *StoreRequest) XXX_Size() int {
	return xxx_messageInfo_StoreRequest.Size(m)
}
func (m *StoreRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StoreRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StoreRequest proto.InternalMessageInfo

func (m *StoreRequest) GetEnvelope() *envelope.Envelope {
	if m != nil {
		return m.Envelope
	}
	return nil
}

type StoreResponse struct {
	// Status of trace element
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// Status of trace element
	LastUpdated *timestamp.Timestamp `protobuf:"bytes,2,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	// Trace object
	Envelope *envelope.Envelope `protobuf:"bytes,3,opt,name=envelope,proto3" json:"envelope,omitempty"`
	// Error
	Err                  *common.Error `protobuf:"bytes,4,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *StoreResponse) Reset()         { *m = StoreResponse{} }
func (m *StoreResponse) String() string { return proto.CompactTextString(m) }
func (*StoreResponse) ProtoMessage()    {}
func (*StoreResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_79926f01850c7b5c, []int{1}
}

func (m *StoreResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StoreResponse.Unmarshal(m, b)
}
func (m *StoreResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StoreResponse.Marshal(b, m, deterministic)
}
func (m *StoreResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StoreResponse.Merge(m, src)
}
func (m *StoreResponse) XXX_Size() int {
	return xxx_messageInfo_StoreResponse.Size(m)
}
func (m *StoreResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StoreResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StoreResponse proto.InternalMessageInfo

func (m *StoreResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *StoreResponse) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

func (m *StoreResponse) GetEnvelope() *envelope.Envelope {
	if m != nil {
		return m.Envelope
	}
	return nil
}

func (m *StoreResponse) GetErr() *common.Error {
	if m != nil {
		return m.Err
	}
	return nil
}

type TxHashRequest struct {
	// Chain ID the transaction has been sent to
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// Hash of the transaction
	TxHash               string   `protobuf:"bytes,2,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxHashRequest) Reset()         { *m = TxHashRequest{} }
func (m *TxHashRequest) String() string { return proto.CompactTextString(m) }
func (*TxHashRequest) ProtoMessage()    {}
func (*TxHashRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_79926f01850c7b5c, []int{2}
}

func (m *TxHashRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxHashRequest.Unmarshal(m, b)
}
func (m *TxHashRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxHashRequest.Marshal(b, m, deterministic)
}
func (m *TxHashRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxHashRequest.Merge(m, src)
}
func (m *TxHashRequest) XXX_Size() int {
	return xxx_messageInfo_TxHashRequest.Size(m)
}
func (m *TxHashRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TxHashRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TxHashRequest proto.InternalMessageInfo

func (m *TxHashRequest) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

func (m *TxHashRequest) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

type IDRequest struct {
	// Envelope identifier
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDRequest) Reset()         { *m = IDRequest{} }
func (m *IDRequest) String() string { return proto.CompactTextString(m) }
func (*IDRequest) ProtoMessage()    {}
func (*IDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_79926f01850c7b5c, []int{3}
}

func (m *IDRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDRequest.Unmarshal(m, b)
}
func (m *IDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDRequest.Marshal(b, m, deterministic)
}
func (m *IDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDRequest.Merge(m, src)
}
func (m *IDRequest) XXX_Size() int {
	return xxx_messageInfo_IDRequest.Size(m)
}
func (m *IDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IDRequest proto.InternalMessageInfo

func (m *IDRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SetStatusRequest struct {
	// Trace identifier
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Status
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetStatusRequest) Reset()         { *m = SetStatusRequest{} }
func (m *SetStatusRequest) String() string { return proto.CompactTextString(m) }
func (*SetStatusRequest) ProtoMessage()    {}
func (*SetStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_79926f01850c7b5c, []int{4}
}

func (m *SetStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetStatusRequest.Unmarshal(m, b)
}
func (m *SetStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetStatusRequest.Marshal(b, m, deterministic)
}
func (m *SetStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetStatusRequest.Merge(m, src)
}
func (m *SetStatusRequest) XXX_Size() int {
	return xxx_messageInfo_SetStatusRequest.Size(m)
}
func (m *SetStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetStatusRequest proto.InternalMessageInfo

func (m *SetStatusRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SetStatusRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type PendingTracesRequest struct {
	// Pending duration in nanoseconds
	Duration             int64    `protobuf:"varint,1,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PendingTracesRequest) Reset()         { *m = PendingTracesRequest{} }
func (m *PendingTracesRequest) String() string { return proto.CompactTextString(m) }
func (*PendingTracesRequest) ProtoMessage()    {}
func (*PendingTracesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_79926f01850c7b5c, []int{5}
}

func (m *PendingTracesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PendingTracesRequest.Unmarshal(m, b)
}
func (m *PendingTracesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PendingTracesRequest.Marshal(b, m, deterministic)
}
func (m *PendingTracesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingTracesRequest.Merge(m, src)
}
func (m *PendingTracesRequest) XXX_Size() int {
	return xxx_messageInfo_PendingTracesRequest.Size(m)
}
func (m *PendingTracesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingTracesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PendingTracesRequest proto.InternalMessageInfo

func (m *PendingTracesRequest) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

type PendingTracesResponse struct {
	// Pending traces
	Envelopes []*envelope.Envelope `protobuf:"bytes,1,rep,name=envelopes,proto3" json:"envelopes,omitempty"`
	// Error
	Err                  *common.Error `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *PendingTracesResponse) Reset()         { *m = PendingTracesResponse{} }
func (m *PendingTracesResponse) String() string { return proto.CompactTextString(m) }
func (*PendingTracesResponse) ProtoMessage()    {}
func (*PendingTracesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_79926f01850c7b5c, []int{6}
}

func (m *PendingTracesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PendingTracesResponse.Unmarshal(m, b)
}
func (m *PendingTracesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PendingTracesResponse.Marshal(b, m, deterministic)
}
func (m *PendingTracesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingTracesResponse.Merge(m, src)
}
func (m *PendingTracesResponse) XXX_Size() int {
	return xxx_messageInfo_PendingTracesResponse.Size(m)
}
func (m *PendingTracesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingTracesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PendingTracesResponse proto.InternalMessageInfo

func (m *PendingTracesResponse) GetEnvelopes() []*envelope.Envelope {
	if m != nil {
		return m.Envelopes
	}
	return nil
}

func (m *PendingTracesResponse) GetErr() *common.Error {
	if m != nil {
		return m.Err
	}
	return nil
}

func init() {
	proto.RegisterType((*StoreRequest)(nil), "contextstore.StoreRequest")
	proto.RegisterType((*StoreResponse)(nil), "contextstore.StoreResponse")
	proto.RegisterType((*TxHashRequest)(nil), "contextstore.TxHashRequest")
	proto.RegisterType((*IDRequest)(nil), "contextstore.IDRequest")
	proto.RegisterType((*SetStatusRequest)(nil), "contextstore.SetStatusRequest")
	proto.RegisterType((*PendingTracesRequest)(nil), "contextstore.PendingTracesRequest")
	proto.RegisterType((*PendingTracesResponse)(nil), "contextstore.PendingTracesResponse")
}

func init() { proto.RegisterFile("protos/context-store/store.proto", fileDescriptor_79926f01850c7b5c) }

var fileDescriptor_79926f01850c7b5c = []byte{
	// 522 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x41, 0x6f, 0xd3, 0x30,
	0x18, 0xdd, 0x5a, 0xe8, 0x9a, 0xaf, 0x2d, 0x42, 0x16, 0xb0, 0x2e, 0x93, 0xb6, 0xc9, 0x5c, 0xb8,
	0xcc, 0x46, 0xe5, 0x86, 0x80, 0x43, 0xbb, 0xc1, 0x8a, 0x38, 0xa0, 0xb4, 0x48, 0x88, 0x4b, 0xe5,
	0x26, 0x5e, 0x1a, 0xd6, 0xda, 0xc1, 0x76, 0x51, 0xf7, 0xcb, 0xe0, 0xe7, 0xa1, 0x38, 0x71, 0x68,
	0xaa, 0x96, 0x4a, 0x5c, 0x22, 0xdb, 0xdf, 0xfb, 0x9e, 0x9f, 0xdf, 0xfb, 0x14, 0xb8, 0x48, 0x95,
	0x34, 0x52, 0xd3, 0x50, 0x0a, 0xc3, 0x57, 0xe6, 0x52, 0x1b, 0xa9, 0x38, 0xb5, 0x5f, 0x62, 0x4b,
	0xa8, 0x5d, 0x94, 0xec, 0x99, 0x7f, 0x56, 0xe0, 0xb9, 0xf8, 0xc9, 0xe7, 0x32, 0xe5, 0xe5, 0x22,
	0x47, 0xfb, 0x27, 0x25, 0xdf, 0x62, 0x21, 0x05, 0xe5, 0x4a, 0x49, 0x55, 0x94, 0xce, 0x63, 0x29,
	0xe3, 0x39, 0xa7, 0x76, 0x37, 0x5d, 0xde, 0x52, 0x93, 0x2c, 0xb8, 0x36, 0x6c, 0x91, 0xe6, 0x00,
	0xfc, 0x0e, 0xda, 0xa3, 0xec, 0x92, 0x80, 0xff, 0x58, 0x72, 0x6d, 0x10, 0x81, 0xa6, 0x63, 0xef,
	0x1e, 0x5e, 0x1c, 0xbe, 0x68, 0xf5, 0x10, 0x29, 0xaf, 0xbb, 0x2e, 0x16, 0x41, 0x89, 0xc1, 0xbf,
	0x0e, 0xa1, 0x53, 0x10, 0xe8, 0x54, 0x0a, 0xcd, 0xd1, 0x33, 0x68, 0x68, 0xc3, 0xcc, 0x52, 0xdb,
	0x7e, 0x2f, 0x28, 0x76, 0xe8, 0x2d, 0xb4, 0xe7, 0x4c, 0x9b, 0xc9, 0x32, 0x8d, 0x98, 0xe1, 0x51,
	0xb7, 0x66, 0xd9, 0x7d, 0x92, 0x2b, 0x24, 0x4e, 0x21, 0x19, 0x3b, 0x85, 0x41, 0x2b, 0xc3, 0x7f,
	0xc9, 0xe1, 0x15, 0x61, 0xf5, 0xfd, 0xc2, 0xd0, 0x39, 0xd4, 0xb9, 0x52, 0xdd, 0x07, 0x16, 0xda,
	0x21, 0xb9, 0x37, 0xe4, 0x3a, 0xf3, 0x26, 0xc8, 0x2a, 0x78, 0x00, 0x9d, 0xf1, 0xea, 0x86, 0xe9,
	0x99, 0x7b, 0xfa, 0x09, 0x34, 0xc3, 0x19, 0x4b, 0xc4, 0x24, 0x89, 0x0a, 0xe9, 0x47, 0x76, 0x3f,
	0x8c, 0xd0, 0x31, 0x1c, 0x99, 0xd5, 0x64, 0xc6, 0xf4, 0xcc, 0xca, 0xf6, 0x82, 0x86, 0xb1, 0xad,
	0xf8, 0x14, 0xbc, 0xe1, 0x95, 0x23, 0x78, 0x04, 0xb5, 0xb2, 0xb5, 0x96, 0x44, 0xf8, 0x35, 0x3c,
	0x1e, 0x71, 0x33, 0xb2, 0xcf, 0xdf, 0x81, 0x59, 0x73, 0xab, 0xb6, 0xee, 0x16, 0xee, 0xc1, 0x93,
	0xcf, 0x5c, 0x44, 0x89, 0x88, 0xc7, 0x8a, 0x85, 0xbc, 0xec, 0xf7, 0xa1, 0x19, 0x2d, 0x15, 0x33,
	0x89, 0x14, 0x96, 0xa5, 0x1e, 0x94, 0x7b, 0xfc, 0x1d, 0x9e, 0x6e, 0xf4, 0x14, 0x91, 0xbc, 0x04,
	0xcf, 0xf9, 0x92, 0xa5, 0x52, 0xdf, 0x61, 0xde, 0x5f, 0x90, 0x73, 0xaf, 0xb6, 0xcb, 0xbd, 0xde,
	0xef, 0x3a, 0x3c, 0xb4, 0xb9, 0xa3, 0xbe, 0x5b, 0xf8, 0x64, 0x7d, 0x6a, 0xc9, 0xfa, 0x58, 0xf9,
	0xa7, 0x5b, 0x6b, 0xb9, 0x3c, 0x7c, 0x80, 0x3e, 0x42, 0xfb, 0x93, 0x64, 0x51, 0xff, 0x3e, 0x4f,
	0x04, 0x6d, 0xc0, 0x2b, 0x39, 0xed, 0xe3, 0xea, 0x43, 0x33, 0xe7, 0x1a, 0x5e, 0xa1, 0xe3, 0x2a,
	0xb4, 0x8c, 0x6a, 0x1f, 0xc7, 0x00, 0xbc, 0x0f, 0x2e, 0xb9, 0xff, 0x26, 0x79, 0x03, 0x5e, 0x19,
	0x3f, 0x3a, 0xdb, 0xc0, 0x6e, 0xcc, 0x85, 0x5f, 0xf5, 0x18, 0x1f, 0xa0, 0xaf, 0xd0, 0xca, 0x9e,
	0x51, 0x04, 0x8a, 0x70, 0xb5, 0x7f, 0xdb, 0x6c, 0xf8, 0xcf, 0xff, 0x89, 0x71, 0xba, 0xfa, 0x37,
	0xdf, 0xde, 0xc7, 0x89, 0x99, 0xb3, 0x69, 0x76, 0x25, 0x1d, 0x64, 0xa7, 0x62, 0x74, 0xaf, 0x69,
	0x38, 0x4f, 0xb8, 0x30, 0xf4, 0x56, 0xd1, 0x50, 0x2a, 0x7e, 0xa9, 0x0d, 0x0b, 0xef, 0x68, 0x7a,
	0x17, 0x93, 0x38, 0x31, 0x74, 0xdb, 0x2f, 0x6b, 0xda, 0xb0, 0xa7, 0xaf, 0xfe, 0x04, 0x00, 0x00,
	0xff, 0xff, 0x44, 0x9a, 0x1c, 0x89, 0xd1, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StoreClient is the client API for Store service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StoreClient interface {
	// Store an envelope
	Store(ctx context.Context, in *StoreRequest, opts ...grpc.CallOption) (*StoreResponse, error)
	// LoadByTxHash load an envelope by transaction hash
	LoadByTxHash(ctx context.Context, in *TxHashRequest, opts ...grpc.CallOption) (*StoreResponse, error)
	// LoadByID load an envelope by identifier
	LoadByID(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*StoreResponse, error)
	// GetStatus returns trace status
	GetStatus(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*StoreResponse, error)
	// SetStatus set a trace status
	SetStatus(ctx context.Context, in *SetStatusRequest, opts ...grpc.CallOption) (*common.Error, error)
	// LoadPending load pending traces
	LoadPending(ctx context.Context, in *PendingTracesRequest, opts ...grpc.CallOption) (*PendingTracesResponse, error)
}

type storeClient struct {
	cc *grpc.ClientConn
}

func NewStoreClient(cc *grpc.ClientConn) StoreClient {
	return &storeClient{cc}
}

func (c *storeClient) Store(ctx context.Context, in *StoreRequest, opts ...grpc.CallOption) (*StoreResponse, error) {
	out := new(StoreResponse)
	err := c.cc.Invoke(ctx, "/contextstore.Store/Store", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) LoadByTxHash(ctx context.Context, in *TxHashRequest, opts ...grpc.CallOption) (*StoreResponse, error) {
	out := new(StoreResponse)
	err := c.cc.Invoke(ctx, "/contextstore.Store/LoadByTxHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) LoadByID(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*StoreResponse, error) {
	out := new(StoreResponse)
	err := c.cc.Invoke(ctx, "/contextstore.Store/LoadByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) GetStatus(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*StoreResponse, error) {
	out := new(StoreResponse)
	err := c.cc.Invoke(ctx, "/contextstore.Store/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) SetStatus(ctx context.Context, in *SetStatusRequest, opts ...grpc.CallOption) (*common.Error, error) {
	out := new(common.Error)
	err := c.cc.Invoke(ctx, "/contextstore.Store/SetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) LoadPending(ctx context.Context, in *PendingTracesRequest, opts ...grpc.CallOption) (*PendingTracesResponse, error) {
	out := new(PendingTracesResponse)
	err := c.cc.Invoke(ctx, "/contextstore.Store/LoadPending", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StoreServer is the server API for Store service.
type StoreServer interface {
	// Store an envelope
	Store(context.Context, *StoreRequest) (*StoreResponse, error)
	// LoadByTxHash load an envelope by transaction hash
	LoadByTxHash(context.Context, *TxHashRequest) (*StoreResponse, error)
	// LoadByID load an envelope by identifier
	LoadByID(context.Context, *IDRequest) (*StoreResponse, error)
	// GetStatus returns trace status
	GetStatus(context.Context, *IDRequest) (*StoreResponse, error)
	// SetStatus set a trace status
	SetStatus(context.Context, *SetStatusRequest) (*common.Error, error)
	// LoadPending load pending traces
	LoadPending(context.Context, *PendingTracesRequest) (*PendingTracesResponse, error)
}

func RegisterStoreServer(s *grpc.Server, srv StoreServer) {
	s.RegisterService(&_Store_serviceDesc, srv)
}

func _Store_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contextstore.Store/Store",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).Store(ctx, req.(*StoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_LoadByTxHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).LoadByTxHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contextstore.Store/LoadByTxHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).LoadByTxHash(ctx, req.(*TxHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_LoadByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).LoadByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contextstore.Store/LoadByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).LoadByID(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contextstore.Store/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).GetStatus(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_SetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).SetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contextstore.Store/SetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).SetStatus(ctx, req.(*SetStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_LoadPending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PendingTracesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).LoadPending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contextstore.Store/LoadPending",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).LoadPending(ctx, req.(*PendingTracesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Store_serviceDesc = grpc.ServiceDesc{
	ServiceName: "contextstore.Store",
	HandlerType: (*StoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Store",
			Handler:    _Store_Store_Handler,
		},
		{
			MethodName: "LoadByTxHash",
			Handler:    _Store_LoadByTxHash_Handler,
		},
		{
			MethodName: "LoadByID",
			Handler:    _Store_LoadByID_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _Store_GetStatus_Handler,
		},
		{
			MethodName: "SetStatus",
			Handler:    _Store_SetStatus_Handler,
		},
		{
			MethodName: "LoadPending",
			Handler:    _Store_LoadPending_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/context-store/store.proto",
}
