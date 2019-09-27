// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types/common/account_instance.proto

package common

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	chain "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
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

// Information about an account on a chain
type AccountInstance struct {
	// Chain ID of the chain
	Chain *chain.Chain `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	// Deployment address of the contract
	Account              *ethereum.Account `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AccountInstance) Reset()         { *m = AccountInstance{} }
func (m *AccountInstance) String() string { return proto.CompactTextString(m) }
func (*AccountInstance) ProtoMessage()    {}
func (*AccountInstance) Descriptor() ([]byte, []int) {
	return fileDescriptor_916420242451c119, []int{0}
}

func (m *AccountInstance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountInstance.Unmarshal(m, b)
}
func (m *AccountInstance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountInstance.Marshal(b, m, deterministic)
}
func (m *AccountInstance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountInstance.Merge(m, src)
}
func (m *AccountInstance) XXX_Size() int {
	return xxx_messageInfo_AccountInstance.Size(m)
}
func (m *AccountInstance) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountInstance.DiscardUnknown(m)
}

var xxx_messageInfo_AccountInstance proto.InternalMessageInfo

func (m *AccountInstance) GetChain() *chain.Chain {
	if m != nil {
		return m.Chain
	}
	return nil
}

func (m *AccountInstance) GetAccount() *ethereum.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func init() {
	proto.RegisterType((*AccountInstance)(nil), "common.AccountInstance")
}

func init() {
	proto.RegisterFile("types/common/account_instance.proto", fileDescriptor_916420242451c119)
}

var fileDescriptor_916420242451c119 = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x31, 0x4b, 0x04, 0x31,
	0x10, 0x85, 0x39, 0xc1, 0x13, 0xa2, 0x20, 0x6e, 0xe3, 0x79, 0x95, 0x9c, 0x8d, 0x20, 0x66, 0x40,
	0x7b, 0x51, 0xaf, 0xb2, 0x3d, 0x3b, 0x1b, 0x49, 0x86, 0x71, 0x2f, 0xec, 0x66, 0x66, 0x49, 0x66,
	0x8b, 0xfd, 0xf7, 0xe2, 0x26, 0x0b, 0x36, 0xaf, 0x79, 0x1f, 0xdf, 0xbc, 0x31, 0x77, 0x3a, 0x0d,
	0x94, 0x01, 0x25, 0x46, 0x61, 0x70, 0x88, 0x32, 0xb2, 0x7e, 0x07, 0xce, 0xea, 0x18, 0xc9, 0x0e,
	0x49, 0x54, 0x9a, 0x75, 0xa9, 0xb7, 0x37, 0x05, 0x26, 0x3d, 0x52, 0xa2, 0x31, 0x82, 0x77, 0xb9,
	0x22, 0xdb, 0xeb, 0xea, 0x39, 0xba, 0xc0, 0x25, 0x4b, 0xb1, 0xf3, 0xe6, 0xf2, 0xad, 0x58, 0x3f,
	0xaa, 0xb4, 0xd9, 0x99, 0xd3, 0x99, 0xd8, 0xac, 0x6e, 0x57, 0xf7, 0xe7, 0x4f, 0x17, 0xb6, 0xf0,
	0xfb, 0xbf, 0x3c, 0x94, 0xaa, 0x79, 0x30, 0x67, 0x75, 0xcc, 0xe6, 0x64, 0xa6, 0xae, 0xec, 0x72,
	0xd6, 0x56, 0xdf, 0x61, 0x21, 0xde, 0x5f, 0xbf, 0x5e, 0xda, 0xa0, 0xbd, 0xf3, 0x16, 0x25, 0xc2,
	0x5e, 0x38, 0x13, 0x7f, 0x4e, 0x19, 0xb0, 0x0f, 0xc4, 0x0a, 0x3f, 0x09, 0x50, 0x12, 0x3d, 0x66,
	0x75, 0xd8, 0xc1, 0xd0, 0xb5, 0xb6, 0x0d, 0x0a, 0xff, 0x1f, 0xf7, 0xeb, 0x79, 0xec, 0xf3, 0x6f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x83, 0x0e, 0xee, 0xa0, 0x0f, 0x01, 0x00, 0x00,
}
