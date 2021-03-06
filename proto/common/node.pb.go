// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/node.proto

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

type NodeStatus int32

const (
	NodeStatus_INACTIVE NodeStatus = 0
	NodeStatus_ACTIVE   NodeStatus = 1
)

var NodeStatus_name = map[int32]string{
	0: "INACTIVE",
	1: "ACTIVE",
}

var NodeStatus_value = map[string]int32{
	"INACTIVE": 0,
	"ACTIVE":   1,
}

func (x NodeStatus) String() string {
	return proto.EnumName(NodeStatus_name, int32(x))
}

func (NodeStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fe0d67cf988b3585, []int{0}
}

type Node struct {
	Endpoint             string     `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Status               NodeStatus `protobuf:"varint,2,opt,name=status,proto3,enum=common.NodeStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe0d67cf988b3585, []int{0}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *Node) GetStatus() NodeStatus {
	if m != nil {
		return m.Status
	}
	return NodeStatus_INACTIVE
}

type SelfNode struct {
	Node                 *Node    `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	AllNodes             uint32   `protobuf:"varint,2,opt,name=all_nodes,json=allNodes,proto3" json:"all_nodes,omitempty"`
	ActiveNodes          uint32   `protobuf:"varint,3,opt,name=active_nodes,json=activeNodes,proto3" json:"active_nodes,omitempty"`
	InactiveNodes        uint32   `protobuf:"varint,4,opt,name=inactive_nodes,json=inactiveNodes,proto3" json:"inactive_nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SelfNode) Reset()         { *m = SelfNode{} }
func (m *SelfNode) String() string { return proto.CompactTextString(m) }
func (*SelfNode) ProtoMessage()    {}
func (*SelfNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe0d67cf988b3585, []int{1}
}

func (m *SelfNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfNode.Unmarshal(m, b)
}
func (m *SelfNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfNode.Marshal(b, m, deterministic)
}
func (m *SelfNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfNode.Merge(m, src)
}
func (m *SelfNode) XXX_Size() int {
	return xxx_messageInfo_SelfNode.Size(m)
}
func (m *SelfNode) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfNode.DiscardUnknown(m)
}

var xxx_messageInfo_SelfNode proto.InternalMessageInfo

func (m *SelfNode) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *SelfNode) GetAllNodes() uint32 {
	if m != nil {
		return m.AllNodes
	}
	return 0
}

func (m *SelfNode) GetActiveNodes() uint32 {
	if m != nil {
		return m.ActiveNodes
	}
	return 0
}

func (m *SelfNode) GetInactiveNodes() uint32 {
	if m != nil {
		return m.InactiveNodes
	}
	return 0
}

type OrgNodes struct {
	Myself               *SelfNode `protobuf:"bytes,1,opt,name=myself,proto3" json:"myself,omitempty"`
	Nodes                []*Node   `protobuf:"bytes,2,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *OrgNodes) Reset()         { *m = OrgNodes{} }
func (m *OrgNodes) String() string { return proto.CompactTextString(m) }
func (*OrgNodes) ProtoMessage()    {}
func (*OrgNodes) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe0d67cf988b3585, []int{2}
}

func (m *OrgNodes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrgNodes.Unmarshal(m, b)
}
func (m *OrgNodes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrgNodes.Marshal(b, m, deterministic)
}
func (m *OrgNodes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrgNodes.Merge(m, src)
}
func (m *OrgNodes) XXX_Size() int {
	return xxx_messageInfo_OrgNodes.Size(m)
}
func (m *OrgNodes) XXX_DiscardUnknown() {
	xxx_messageInfo_OrgNodes.DiscardUnknown(m)
}

var xxx_messageInfo_OrgNodes proto.InternalMessageInfo

func (m *OrgNodes) GetMyself() *SelfNode {
	if m != nil {
		return m.Myself
	}
	return nil
}

func (m *OrgNodes) GetNodes() []*Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func init() {
	proto.RegisterEnum("common.NodeStatus", NodeStatus_name, NodeStatus_value)
	proto.RegisterType((*Node)(nil), "common.Node")
	proto.RegisterType((*SelfNode)(nil), "common.SelfNode")
	proto.RegisterType((*OrgNodes)(nil), "common.OrgNodes")
}

func init() {
	proto.RegisterFile("common/node.proto", fileDescriptor_fe0d67cf988b3585)
}

var fileDescriptor_fe0d67cf988b3585 = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x5f, 0x4b, 0xc3, 0x30,
	0x14, 0xc5, 0xad, 0x9b, 0xa5, 0xbb, 0xfb, 0xc3, 0xcc, 0xd3, 0xd0, 0x97, 0x5a, 0x51, 0xca, 0x1e,
	0x5a, 0x98, 0xf8, 0x01, 0x54, 0x7c, 0xd8, 0x4b, 0x85, 0x4e, 0x44, 0x44, 0x90, 0x98, 0x66, 0x6b,
	0x20, 0x4d, 0x4a, 0x93, 0x39, 0xfc, 0x1c, 0x7e, 0x61, 0x69, 0x92, 0xea, 0xc4, 0xb7, 0x7b, 0xef,
	0x39, 0xbf, 0xd3, 0x43, 0x03, 0xc7, 0x44, 0x56, 0x95, 0x14, 0xa9, 0x90, 0x05, 0x4d, 0xea, 0x46,
	0x6a, 0x89, 0x7c, 0x7b, 0x8a, 0x32, 0xe8, 0x67, 0xb2, 0xa0, 0xe8, 0x04, 0x02, 0x2a, 0x8a, 0x5a,
	0x32, 0xa1, 0x67, 0x5e, 0xe8, 0xc5, 0x83, 0xfc, 0x67, 0x47, 0x73, 0xf0, 0x95, 0xc6, 0x7a, 0xab,
	0x66, 0x87, 0xa1, 0x17, 0x4f, 0x16, 0x28, 0xb1, 0x70, 0xd2, 0x92, 0x2b, 0xa3, 0xe4, 0xce, 0x11,
	0x7d, 0x79, 0x10, 0xac, 0x28, 0x5f, 0x9b, 0xd0, 0x10, 0xfa, 0xed, 0x27, 0x4d, 0xe0, 0x70, 0x31,
	0xda, 0xc7, 0x72, 0xa3, 0xa0, 0x53, 0x18, 0x60, 0xce, 0xdf, 0xda, 0xd9, 0xa6, 0x8f, 0xf3, 0x00,
	0x73, 0xde, 0x3a, 0x14, 0x3a, 0x83, 0x11, 0x26, 0x9a, 0x7d, 0x50, 0xa7, 0xf7, 0x8c, 0x3e, 0xb4,
	0x37, 0x6b, 0xb9, 0x80, 0x09, 0x13, 0x7f, 0x4c, 0x7d, 0x63, 0x1a, 0x77, 0x57, 0x63, 0x8b, 0x9e,
	0x21, 0x78, 0x68, 0x36, 0x16, 0x89, 0xc1, 0xaf, 0x3e, 0x15, 0xe5, 0x6b, 0x57, 0x6b, 0xda, 0xd5,
	0xea, 0x6a, 0xe7, 0x4e, 0x47, 0x11, 0x1c, 0x75, 0xc5, 0x7a, 0xff, 0xfa, 0x5b, 0x69, 0x7e, 0x09,
	0xf0, 0xfb, 0x17, 0xd0, 0x08, 0x82, 0x65, 0x76, 0x73, 0xf7, 0xb8, 0x7c, 0xba, 0x9f, 0x1e, 0x20,
	0x00, 0xdf, 0xcd, 0xde, 0xed, 0x2b, 0x9c, 0x13, 0x59, 0x25, 0xe5, 0x16, 0xef, 0x28, 0x4b, 0x76,
	0x8c, 0x0a, 0xda, 0x90, 0x12, 0x33, 0x61, 0xdf, 0xc3, 0x45, 0xbf, 0x5c, 0x6f, 0x98, 0xee, 0x4c,
	0x44, 0x56, 0x69, 0x2d, 0x99, 0x52, 0x52, 0x28, 0x8a, 0x1b, 0x52, 0xa6, 0x7b, 0x58, 0x6a, 0xb0,
	0xd4, 0x62, 0xef, 0xbe, 0xd9, 0xae, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0x87, 0xf3, 0xea, 0xfe,
	0xe9, 0x01, 0x00, 0x00,
}
