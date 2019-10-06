// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/services/envelope-store/store.proto

package envelope_store

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	chain "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/chain"
	envelope "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/envelope"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/ethereum"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Possible statuses for a transaction
type Status int32

const (
	Status_STORED  Status = 0
	Status_PENDING Status = 1
	Status_MINED   Status = 2
	Status_ERROR   Status = 3
)

var Status_name = map[int32]string{
	0: "STORED",
	1: "PENDING",
	2: "MINED",
	3: "ERROR",
}

var Status_value = map[string]int32{
	"STORED":  0,
	"PENDING": 1,
	"MINED":   2,
	"ERROR":   3,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{0}
}

type StatusInfo struct {
	// Status of Transaction
	Status Status `protobuf:"varint,1,opt,name=status,proto3,enum=envelopestore.Status" json:"status,omitempty"`
	// Date events for envelope
	StoredAt             *timestamp.Timestamp `protobuf:"bytes,2,opt,name=storedAt,proto3" json:"storedAt,omitempty"`
	SentAt               *timestamp.Timestamp `protobuf:"bytes,3,opt,name=sentAt,proto3" json:"sentAt,omitempty"`
	MinedAt              *timestamp.Timestamp `protobuf:"bytes,4,opt,name=minedAt,proto3" json:"minedAt,omitempty"`
	ErrorAt              *timestamp.Timestamp `protobuf:"bytes,5,opt,name=errorAt,proto3" json:"errorAt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *StatusInfo) Reset()         { *m = StatusInfo{} }
func (m *StatusInfo) String() string { return proto.CompactTextString(m) }
func (*StatusInfo) ProtoMessage()    {}
func (*StatusInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{0}
}

func (m *StatusInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusInfo.Unmarshal(m, b)
}
func (m *StatusInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusInfo.Marshal(b, m, deterministic)
}
func (m *StatusInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusInfo.Merge(m, src)
}
func (m *StatusInfo) XXX_Size() int {
	return xxx_messageInfo_StatusInfo.Size(m)
}
func (m *StatusInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusInfo.DiscardUnknown(m)
}

var xxx_messageInfo_StatusInfo proto.InternalMessageInfo

func (m *StatusInfo) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_STORED
}

func (m *StatusInfo) GetStoredAt() *timestamp.Timestamp {
	if m != nil {
		return m.StoredAt
	}
	return nil
}

func (m *StatusInfo) GetSentAt() *timestamp.Timestamp {
	if m != nil {
		return m.SentAt
	}
	return nil
}

func (m *StatusInfo) GetMinedAt() *timestamp.Timestamp {
	if m != nil {
		return m.MinedAt
	}
	return nil
}

func (m *StatusInfo) GetErrorAt() *timestamp.Timestamp {
	if m != nil {
		return m.ErrorAt
	}
	return nil
}

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
	return fileDescriptor_77281f742fc35d14, []int{1}
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
	// Envelope
	Envelope             *envelope.Envelope `protobuf:"bytes,1,opt,name=envelope,proto3" json:"envelope,omitempty"`
	StatusInfo           *StatusInfo        `protobuf:"bytes,2,opt,name=status_info,json=statusInfo,proto3" json:"status_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *StoreResponse) Reset()         { *m = StoreResponse{} }
func (m *StoreResponse) String() string { return proto.CompactTextString(m) }
func (*StoreResponse) ProtoMessage()    {}
func (*StoreResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{2}
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

func (m *StoreResponse) GetEnvelope() *envelope.Envelope {
	if m != nil {
		return m.Envelope
	}
	return nil
}

func (m *StoreResponse) GetStatusInfo() *StatusInfo {
	if m != nil {
		return m.StatusInfo
	}
	return nil
}

type LoadByIDRequest struct {
	// Envelope identifier
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoadByIDRequest) Reset()         { *m = LoadByIDRequest{} }
func (m *LoadByIDRequest) String() string { return proto.CompactTextString(m) }
func (*LoadByIDRequest) ProtoMessage()    {}
func (*LoadByIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{3}
}

func (m *LoadByIDRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadByIDRequest.Unmarshal(m, b)
}
func (m *LoadByIDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadByIDRequest.Marshal(b, m, deterministic)
}
func (m *LoadByIDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadByIDRequest.Merge(m, src)
}
func (m *LoadByIDRequest) XXX_Size() int {
	return xxx_messageInfo_LoadByIDRequest.Size(m)
}
func (m *LoadByIDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadByIDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoadByIDRequest proto.InternalMessageInfo

func (m *LoadByIDRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type LoadByTxHashRequest struct {
	// Chain the transaction has been sent to
	Chain *chain.Chain `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
	// Hash of the transaction
	TxHash               *ethereum.Hash `protobuf:"bytes,3,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *LoadByTxHashRequest) Reset()         { *m = LoadByTxHashRequest{} }
func (m *LoadByTxHashRequest) String() string { return proto.CompactTextString(m) }
func (*LoadByTxHashRequest) ProtoMessage()    {}
func (*LoadByTxHashRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{4}
}

func (m *LoadByTxHashRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadByTxHashRequest.Unmarshal(m, b)
}
func (m *LoadByTxHashRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadByTxHashRequest.Marshal(b, m, deterministic)
}
func (m *LoadByTxHashRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadByTxHashRequest.Merge(m, src)
}
func (m *LoadByTxHashRequest) XXX_Size() int {
	return xxx_messageInfo_LoadByTxHashRequest.Size(m)
}
func (m *LoadByTxHashRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadByTxHashRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoadByTxHashRequest proto.InternalMessageInfo

func (m *LoadByTxHashRequest) GetChain() *chain.Chain {
	if m != nil {
		return m.Chain
	}
	return nil
}

func (m *LoadByTxHashRequest) GetTxHash() *ethereum.Hash {
	if m != nil {
		return m.TxHash
	}
	return nil
}

type SetStatusRequest struct {
	// Envelope identifier
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Status
	Status               Status   `protobuf:"varint,2,opt,name=status,proto3,enum=envelopestore.Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetStatusRequest) Reset()         { *m = SetStatusRequest{} }
func (m *SetStatusRequest) String() string { return proto.CompactTextString(m) }
func (*SetStatusRequest) ProtoMessage()    {}
func (*SetStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{5}
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

func (m *SetStatusRequest) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_STORED
}

type StatusResponse struct {
	StatusInfo           *StatusInfo `protobuf:"bytes,1,opt,name=status_info,json=statusInfo,proto3" json:"status_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{6}
}

func (m *StatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResponse.Unmarshal(m, b)
}
func (m *StatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResponse.Marshal(b, m, deterministic)
}
func (m *StatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResponse.Merge(m, src)
}
func (m *StatusResponse) XXX_Size() int {
	return xxx_messageInfo_StatusResponse.Size(m)
}
func (m *StatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetStatusInfo() *StatusInfo {
	if m != nil {
		return m.StatusInfo
	}
	return nil
}

type LoadPendingRequest struct {
	// Pending duration
	Duration             *duration.Duration `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *LoadPendingRequest) Reset()         { *m = LoadPendingRequest{} }
func (m *LoadPendingRequest) String() string { return proto.CompactTextString(m) }
func (*LoadPendingRequest) ProtoMessage()    {}
func (*LoadPendingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{7}
}

func (m *LoadPendingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadPendingRequest.Unmarshal(m, b)
}
func (m *LoadPendingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadPendingRequest.Marshal(b, m, deterministic)
}
func (m *LoadPendingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadPendingRequest.Merge(m, src)
}
func (m *LoadPendingRequest) XXX_Size() int {
	return xxx_messageInfo_LoadPendingRequest.Size(m)
}
func (m *LoadPendingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadPendingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoadPendingRequest proto.InternalMessageInfo

func (m *LoadPendingRequest) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

type LoadPendingResponse struct {
	// Pending envelopes
	Responses            []*StoreResponse `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *LoadPendingResponse) Reset()         { *m = LoadPendingResponse{} }
func (m *LoadPendingResponse) String() string { return proto.CompactTextString(m) }
func (*LoadPendingResponse) ProtoMessage()    {}
func (*LoadPendingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77281f742fc35d14, []int{8}
}

func (m *LoadPendingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadPendingResponse.Unmarshal(m, b)
}
func (m *LoadPendingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadPendingResponse.Marshal(b, m, deterministic)
}
func (m *LoadPendingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadPendingResponse.Merge(m, src)
}
func (m *LoadPendingResponse) XXX_Size() int {
	return xxx_messageInfo_LoadPendingResponse.Size(m)
}
func (m *LoadPendingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadPendingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LoadPendingResponse proto.InternalMessageInfo

func (m *LoadPendingResponse) GetResponses() []*StoreResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

func init() {
	proto.RegisterEnum("envelopestore.Status", Status_name, Status_value)
	proto.RegisterType((*StatusInfo)(nil), "envelopestore.StatusInfo")
	proto.RegisterType((*StoreRequest)(nil), "envelopestore.StoreRequest")
	proto.RegisterType((*StoreResponse)(nil), "envelopestore.StoreResponse")
	proto.RegisterType((*LoadByIDRequest)(nil), "envelopestore.LoadByIDRequest")
	proto.RegisterType((*LoadByTxHashRequest)(nil), "envelopestore.LoadByTxHashRequest")
	proto.RegisterType((*SetStatusRequest)(nil), "envelopestore.SetStatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "envelopestore.StatusResponse")
	proto.RegisterType((*LoadPendingRequest)(nil), "envelopestore.LoadPendingRequest")
	proto.RegisterType((*LoadPendingResponse)(nil), "envelopestore.LoadPendingResponse")
}

func init() {
	proto.RegisterFile("pkg/services/envelope-store/store.proto", fileDescriptor_77281f742fc35d14)
}

var fileDescriptor_77281f742fc35d14 = []byte{
	// 654 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x5b, 0x6f, 0x1a, 0x3d,
	0x10, 0xfd, 0x80, 0x0f, 0x02, 0x43, 0x42, 0x91, 0xab, 0xaa, 0x84, 0x36, 0xb7, 0x7d, 0x49, 0x54,
	0x29, 0x5e, 0x89, 0xde, 0xa4, 0x3e, 0x54, 0x4a, 0x02, 0x6a, 0x50, 0x52, 0x92, 0x18, 0x9e, 0xfa,
	0x12, 0x2d, 0x60, 0xc0, 0x0a, 0xac, 0xb7, 0x6b, 0x13, 0x05, 0xf5, 0x5f, 0xb4, 0x7f, 0xb8, 0x5a,
	0x5f, 0xb6, 0xb0, 0x5d, 0x85, 0xe6, 0xc5, 0x9a, 0xf5, 0x9c, 0x39, 0x33, 0x7b, 0x7c, 0x2c, 0xc3,
	0x61, 0x70, 0x37, 0x76, 0x05, 0x0d, 0xef, 0xd9, 0x80, 0x0a, 0x97, 0xfa, 0xf7, 0x74, 0xca, 0x03,
	0x7a, 0x2c, 0x24, 0x0f, 0xa9, 0xab, 0x56, 0x1c, 0x84, 0x5c, 0x72, 0xb4, 0x65, 0x73, 0x6a, 0xb3,
	0xbe, 0x23, 0x17, 0xc1, 0x52, 0x41, 0x1c, 0x68, 0x74, 0x7d, 0xdb, 0xa4, 0xe5, 0x84, 0x86, 0x74,
	0x3e, 0x73, 0xfb, 0x9e, 0xb0, 0xa9, 0x97, 0x3a, 0x35, 0x98, 0x78, 0xcc, 0xd7, 0xab, 0x49, 0xec,
	0x8d, 0x39, 0x1f, 0x4f, 0xa9, 0xab, 0xbe, 0xfa, 0xf3, 0x91, 0x2b, 0xd9, 0x8c, 0x0a, 0xe9, 0xcd,
	0x02, 0x03, 0xd8, 0x4d, 0x02, 0x86, 0xf3, 0xd0, 0x93, 0x8c, 0x1b, 0x02, 0xe7, 0x57, 0x16, 0xa0,
	0x2b, 0x3d, 0x39, 0x17, 0x6d, 0x7f, 0xc4, 0xd1, 0x31, 0x14, 0x84, 0xfa, 0xaa, 0x65, 0xf6, 0x33,
	0x47, 0x95, 0xc6, 0x0b, 0xbc, 0xf2, 0x0b, 0x58, 0x43, 0x89, 0x01, 0xa1, 0x0f, 0x50, 0x54, 0xfb,
	0xc3, 0x13, 0x59, 0xcb, 0xee, 0x67, 0x8e, 0xca, 0x8d, 0x3a, 0xd6, 0x0d, 0xb1, 0x6d, 0x88, 0x7b,
	0x76, 0x22, 0x12, 0x63, 0x51, 0x03, 0x0a, 0x82, 0xfa, 0xf2, 0x44, 0xd6, 0x72, 0x6b, 0xab, 0x0c,
	0x12, 0xbd, 0x83, 0x8d, 0x19, 0xf3, 0x55, 0xab, 0xff, 0xd7, 0x16, 0x59, 0x68, 0x54, 0x45, 0xc3,
	0x90, 0x87, 0x27, 0xb2, 0x96, 0x5f, 0x5f, 0x65, 0xa0, 0xce, 0x67, 0xd8, 0xec, 0x46, 0xb3, 0x12,
	0xfa, 0x7d, 0x4e, 0x85, 0x44, 0x18, 0x8a, 0x56, 0x07, 0x25, 0x4c, 0xb9, 0x81, 0x62, 0x61, 0x70,
	0xcb, 0x04, 0x24, 0xc6, 0x38, 0x3f, 0x60, 0xcb, 0xd4, 0x8b, 0x80, 0xfb, 0x82, 0x3e, 0x95, 0x00,
	0x7d, 0x82, 0xb2, 0x96, 0xf8, 0x96, 0xf9, 0x23, 0x6e, 0xb4, 0xdd, 0x4e, 0x3d, 0x8c, 0xe8, 0xdc,
	0x08, 0x88, 0x38, 0x76, 0x0e, 0xe0, 0xd9, 0x25, 0xf7, 0x86, 0xa7, 0x8b, 0x76, 0xd3, 0xce, 0x5f,
	0x81, 0x2c, 0x1b, 0xaa, 0xc6, 0x25, 0x92, 0x65, 0x43, 0xa7, 0x0f, 0xcf, 0x35, 0xa4, 0xf7, 0x70,
	0xee, 0x89, 0x89, 0x85, 0x39, 0x90, 0x57, 0xe6, 0x32, 0xfd, 0x36, 0xb1, 0xb6, 0xda, 0x59, 0xb4,
	0x12, 0x9d, 0x42, 0x87, 0xb0, 0x21, 0x1f, 0x6e, 0x27, 0x9e, 0x98, 0x98, 0xb3, 0xab, 0x60, 0xeb,
	0x58, 0xac, 0xb8, 0x0a, 0x52, 0x71, 0x3a, 0x37, 0x50, 0xed, 0x52, 0x69, 0x0c, 0x93, 0x3e, 0xc7,
	0x92, 0xdd, 0xb2, 0xff, 0x60, 0x37, 0xe7, 0x12, 0x2a, 0x96, 0xcf, 0xe8, 0x9a, 0xd0, 0x29, 0xf3,
	0x14, 0x9d, 0x2e, 0x00, 0x45, 0x22, 0x5c, 0x53, 0x7f, 0xc8, 0xfc, 0xb1, 0x1d, 0xf1, 0x3d, 0x14,
	0xed, 0x15, 0x89, 0xe9, 0x92, 0x8e, 0x69, 0x1a, 0x00, 0x89, 0xa1, 0xce, 0x8d, 0x56, 0x34, 0x26,
	0x8b, 0xe7, 0x2b, 0x85, 0x26, 0x8e, 0xae, 0x54, 0xee, 0xa8, 0xdc, 0x78, 0xfd, 0xd7, 0x74, 0x4b,
	0x46, 0x21, 0x7f, 0xe0, 0x6f, 0x3e, 0x42, 0x41, 0x4f, 0x8e, 0x00, 0x0a, 0xdd, 0xde, 0x15, 0x69,
	0x35, 0xab, 0xff, 0xa1, 0x32, 0x6c, 0x5c, 0xb7, 0x3a, 0xcd, 0x76, 0xe7, 0x4b, 0x35, 0x83, 0x4a,
	0x90, 0xff, 0xda, 0xee, 0xb4, 0x9a, 0xd5, 0x6c, 0x14, 0xb6, 0x08, 0xb9, 0x22, 0xd5, 0x5c, 0xe3,
	0x67, 0x0e, 0xb6, 0xac, 0xa7, 0x14, 0x3b, 0x3a, 0x85, 0xbc, 0x0e, 0x5e, 0xa5, 0x37, 0x57, 0xbf,
	0x5e, 0x7f, 0x74, 0x32, 0x74, 0x0e, 0x45, 0x6b, 0x2b, 0xb4, 0x9b, 0x40, 0x26, 0xfc, 0xb6, 0x86,
	0xe9, 0x1a, 0x36, 0x97, 0xdd, 0x87, 0x9c, 0x54, 0xb6, 0x15, 0x6b, 0xae, 0x61, 0xbc, 0x80, 0x52,
	0xec, 0x35, 0xb4, 0x97, 0x84, 0x26, 0x5c, 0x58, 0xdf, 0x49, 0x77, 0x99, 0x25, 0xeb, 0x41, 0x79,
	0xe9, 0x28, 0xd1, 0x41, 0xca, 0x74, 0xab, 0x9e, 0xa9, 0x3b, 0x8f, 0x41, 0x34, 0xeb, 0x69, 0xe7,
	0xdb, 0xe5, 0x98, 0xc9, 0xa9, 0xd7, 0xc7, 0x03, 0x3e, 0x73, 0xcf, 0xa2, 0x3d, 0xbf, 0xbb, 0x10,
	0xee, 0x60, 0xca, 0xa8, 0x2f, 0xdd, 0x51, 0xe8, 0x0e, 0x78, 0x18, 0x3d, 0x22, 0xde, 0xe0, 0x4e,
	0x85, 0x2a, 0xc2, 0x63, 0x26, 0xdd, 0xd5, 0x67, 0x43, 0xbf, 0x33, 0xfd, 0x82, 0x72, 0xe3, 0xdb,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8a, 0xc2, 0xdf, 0x75, 0x8d, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EnvelopeStoreClient is the client API for EnvelopeStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EnvelopeStoreClient interface {
	// Store an envelope
	Store(ctx context.Context, in *StoreRequest, opts ...grpc.CallOption) (*StoreResponse, error)
	// Load envelope by identifier
	LoadByID(ctx context.Context, in *LoadByIDRequest, opts ...grpc.CallOption) (*StoreResponse, error)
	// Load Envelope by transaction hash
	LoadByTxHash(ctx context.Context, in *LoadByTxHashRequest, opts ...grpc.CallOption) (*StoreResponse, error)
	// SetStatus set an envelope status
	SetStatus(ctx context.Context, in *SetStatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	// LoadPending load envelopes of pending transactions
	LoadPending(ctx context.Context, in *LoadPendingRequest, opts ...grpc.CallOption) (*LoadPendingResponse, error)
}

type envelopeStoreClient struct {
	cc *grpc.ClientConn
}

func NewEnvelopeStoreClient(cc *grpc.ClientConn) EnvelopeStoreClient {
	return &envelopeStoreClient{cc}
}

func (c *envelopeStoreClient) Store(ctx context.Context, in *StoreRequest, opts ...grpc.CallOption) (*StoreResponse, error) {
	out := new(StoreResponse)
	err := c.cc.Invoke(ctx, "/envelopestore.EnvelopeStore/Store", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envelopeStoreClient) LoadByID(ctx context.Context, in *LoadByIDRequest, opts ...grpc.CallOption) (*StoreResponse, error) {
	out := new(StoreResponse)
	err := c.cc.Invoke(ctx, "/envelopestore.EnvelopeStore/LoadByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envelopeStoreClient) LoadByTxHash(ctx context.Context, in *LoadByTxHashRequest, opts ...grpc.CallOption) (*StoreResponse, error) {
	out := new(StoreResponse)
	err := c.cc.Invoke(ctx, "/envelopestore.EnvelopeStore/LoadByTxHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envelopeStoreClient) SetStatus(ctx context.Context, in *SetStatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/envelopestore.EnvelopeStore/SetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envelopeStoreClient) LoadPending(ctx context.Context, in *LoadPendingRequest, opts ...grpc.CallOption) (*LoadPendingResponse, error) {
	out := new(LoadPendingResponse)
	err := c.cc.Invoke(ctx, "/envelopestore.EnvelopeStore/LoadPending", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnvelopeStoreServer is the server API for EnvelopeStore service.
type EnvelopeStoreServer interface {
	// Store an envelope
	Store(context.Context, *StoreRequest) (*StoreResponse, error)
	// Load envelope by identifier
	LoadByID(context.Context, *LoadByIDRequest) (*StoreResponse, error)
	// Load Envelope by transaction hash
	LoadByTxHash(context.Context, *LoadByTxHashRequest) (*StoreResponse, error)
	// SetStatus set an envelope status
	SetStatus(context.Context, *SetStatusRequest) (*StatusResponse, error)
	// LoadPending load envelopes of pending transactions
	LoadPending(context.Context, *LoadPendingRequest) (*LoadPendingResponse, error)
}

// UnimplementedEnvelopeStoreServer can be embedded to have forward compatible implementations.
type UnimplementedEnvelopeStoreServer struct {
}

func (*UnimplementedEnvelopeStoreServer) Store(ctx context.Context, req *StoreRequest) (*StoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}
func (*UnimplementedEnvelopeStoreServer) LoadByID(ctx context.Context, req *LoadByIDRequest) (*StoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadByID not implemented")
}
func (*UnimplementedEnvelopeStoreServer) LoadByTxHash(ctx context.Context, req *LoadByTxHashRequest) (*StoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadByTxHash not implemented")
}
func (*UnimplementedEnvelopeStoreServer) SetStatus(ctx context.Context, req *SetStatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStatus not implemented")
}
func (*UnimplementedEnvelopeStoreServer) LoadPending(ctx context.Context, req *LoadPendingRequest) (*LoadPendingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadPending not implemented")
}

func RegisterEnvelopeStoreServer(s *grpc.Server, srv EnvelopeStoreServer) {
	s.RegisterService(&_EnvelopeStore_serviceDesc, srv)
}

func _EnvelopeStore_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvelopeStoreServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/envelopestore.EnvelopeStore/Store",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvelopeStoreServer).Store(ctx, req.(*StoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnvelopeStore_LoadByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvelopeStoreServer).LoadByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/envelopestore.EnvelopeStore/LoadByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvelopeStoreServer).LoadByID(ctx, req.(*LoadByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnvelopeStore_LoadByTxHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadByTxHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvelopeStoreServer).LoadByTxHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/envelopestore.EnvelopeStore/LoadByTxHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvelopeStoreServer).LoadByTxHash(ctx, req.(*LoadByTxHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnvelopeStore_SetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvelopeStoreServer).SetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/envelopestore.EnvelopeStore/SetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvelopeStoreServer).SetStatus(ctx, req.(*SetStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnvelopeStore_LoadPending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadPendingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvelopeStoreServer).LoadPending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/envelopestore.EnvelopeStore/LoadPending",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvelopeStoreServer).LoadPending(ctx, req.(*LoadPendingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EnvelopeStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "envelopestore.EnvelopeStore",
	HandlerType: (*EnvelopeStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Store",
			Handler:    _EnvelopeStore_Store_Handler,
		},
		{
			MethodName: "LoadByID",
			Handler:    _EnvelopeStore_LoadByID_Handler,
		},
		{
			MethodName: "LoadByTxHash",
			Handler:    _EnvelopeStore_LoadByTxHash_Handler,
		},
		{
			MethodName: "SetStatus",
			Handler:    _EnvelopeStore_SetStatus_Handler,
		},
		{
			MethodName: "LoadPending",
			Handler:    _EnvelopeStore_LoadPending_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/services/envelope-store/store.proto",
}
