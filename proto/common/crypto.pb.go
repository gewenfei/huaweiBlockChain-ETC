// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/crypto.proto

package common

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Certificate struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	PemCert              []byte   `protobuf:"bytes,2,opt,name=pem_cert,json=pemCert,proto3" json:"pem_cert,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Certificate) Reset()         { *m = Certificate{} }
func (m *Certificate) String() string { return proto.CompactTextString(m) }
func (*Certificate) ProtoMessage()    {}
func (*Certificate) Descriptor() ([]byte, []int) {
	return fileDescriptor_0967a5b0afda3b53, []int{0}
}

func (m *Certificate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Certificate.Unmarshal(m, b)
}
func (m *Certificate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Certificate.Marshal(b, m, deterministic)
}
func (m *Certificate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Certificate.Merge(m, src)
}
func (m *Certificate) XXX_Size() int {
	return xxx_messageInfo_Certificate.Size(m)
}
func (m *Certificate) XXX_DiscardUnknown() {
	xxx_messageInfo_Certificate.DiscardUnknown(m)
}

var xxx_messageInfo_Certificate proto.InternalMessageInfo

func (m *Certificate) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Certificate) GetPemCert() []byte {
	if m != nil {
		return m.PemCert
	}
	return nil
}

func init() {
	proto.RegisterType((*Certificate)(nil), "common.Certificate")
}

func init() {
	proto.RegisterFile("common/crypto.proto", fileDescriptor_0967a5b0afda3b53)
}

var fileDescriptor_0967a5b0afda3b53 = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x24, 0xcd, 0xb1, 0xca, 0xc2, 0x40,
	0x10, 0x04, 0x60, 0xf2, 0xf3, 0x13, 0xf5, 0xb4, 0x3a, 0x9b, 0xd8, 0x05, 0xab, 0x54, 0xb9, 0x42,
	0xc4, 0xc6, 0x4a, 0xdf, 0x20, 0xa5, 0x8d, 0x9c, 0xcb, 0x6a, 0xb6, 0xd8, 0xdb, 0x63, 0xb3, 0x12,
	0x7c, 0x7b, 0x31, 0xe9, 0x66, 0x06, 0x86, 0xcf, 0x6d, 0x41, 0x98, 0x25, 0x05, 0xd0, 0x4f, 0x36,
	0x69, 0xb3, 0x8a, 0x89, 0x2f, 0xe7, 0x71, 0x7f, 0x76, 0xeb, 0x2b, 0xaa, 0xd1, 0x93, 0x20, 0x1a,
	0x7a, 0xef, 0xfe, 0x53, 0x64, 0xac, 0x8a, 0xba, 0x68, 0x56, 0xdd, 0x94, 0xfd, 0xce, 0x2d, 0x33,
	0xf2, 0x1d, 0x50, 0xad, 0xfa, 0xab, 0x8b, 0x66, 0xd3, 0x2d, 0x32, 0xf2, 0xef, 0x75, 0x39, 0xdd,
	0x8e, 0x2f, 0xb2, 0xb6, 0x7f, 0xc7, 0x11, 0xa9, 0x05, 0xe1, 0x90, 0x85, 0x86, 0x41, 0xd2, 0x80,
	0x51, 0xa1, 0x0f, 0x23, 0x61, 0x42, 0x85, 0x3e, 0x52, 0x0a, 0x93, 0x1b, 0x66, 0xf6, 0x51, 0x4e,
	0xed, 0xf0, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x81, 0x04, 0x96, 0xce, 0x9c, 0x00, 0x00, 0x00,
}
