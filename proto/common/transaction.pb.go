// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/transaction.proto

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

type TxType int32

const (
	TxType_COMMON_TRANSACTION TxType = 0
	TxType_VOTE_TRANSACTION   TxType = 1
)

var TxType_name = map[int32]string{
	0: "COMMON_TRANSACTION",
	1: "VOTE_TRANSACTION",
}

var TxType_value = map[string]int32{
	"COMMON_TRANSACTION": 0,
	"VOTE_TRANSACTION":   1,
}

func (x TxType) String() string {
	return proto.EnumName(TxType_name, int32(x))
}

func (TxType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{0}
}

type TxStatus int32

const (
	TxStatus_VALID                        TxStatus = 0
	TxStatus_INVALID_TX_PAYLOAD           TxStatus = 1
	TxStatus_INVALID_TX_HEADER            TxStatus = 2
	TxStatus_INVALID_TX_TYPE              TxStatus = 3
	TxStatus_UNKNOWN_TX_TYPE              TxStatus = 4
	TxStatus_DUPLICATE_TX_ID              TxStatus = 5
	TxStatus_INVALID_STATE_UPDATES_TYPE   TxStatus = 6
	TxStatus_INVALID_INSERT_STATEMENTS    TxStatus = 7
	TxStatus_INVALID_SELECT_STATEMENTS    TxStatus = 8
	TxStatus_INVALID_SQL_SYNTAX           TxStatus = 9
	TxStatus_INVALID_SQL_SCHEMA           TxStatus = 10
	TxStatus_INVALID_COMMON_TX_DATA       TxStatus = 11
	TxStatus_INVALID_VOTE_TX_DATA         TxStatus = 12
	TxStatus_INVALID_CHAIN_CONFIG         TxStatus = 13
	TxStatus_INVALID_CONTRACT_INVOCATION  TxStatus = 21
	TxStatus_INVALID_APPROVAL_SIGNATURE   TxStatus = 22
	TxStatus_APPROVALS_VALIDATION_FAILURE TxStatus = 23
	TxStatus_UNKNOWN_VOTE_HANDLER         TxStatus = 31
	TxStatus_INVALID_VOTE_PAYLOAD         TxStatus = 32
	TxStatus_INVALID_MVCC                 TxStatus = 41
	TxStatus_INVALID                      TxStatus = 255
)

var TxStatus_name = map[int32]string{
	0:   "VALID",
	1:   "INVALID_TX_PAYLOAD",
	2:   "INVALID_TX_HEADER",
	3:   "INVALID_TX_TYPE",
	4:   "UNKNOWN_TX_TYPE",
	5:   "DUPLICATE_TX_ID",
	6:   "INVALID_STATE_UPDATES_TYPE",
	7:   "INVALID_INSERT_STATEMENTS",
	8:   "INVALID_SELECT_STATEMENTS",
	9:   "INVALID_SQL_SYNTAX",
	10:  "INVALID_SQL_SCHEMA",
	11:  "INVALID_COMMON_TX_DATA",
	12:  "INVALID_VOTE_TX_DATA",
	13:  "INVALID_CHAIN_CONFIG",
	21:  "INVALID_CONTRACT_INVOCATION",
	22:  "INVALID_APPROVAL_SIGNATURE",
	23:  "APPROVALS_VALIDATION_FAILURE",
	31:  "UNKNOWN_VOTE_HANDLER",
	32:  "INVALID_VOTE_PAYLOAD",
	41:  "INVALID_MVCC",
	255: "INVALID",
}

var TxStatus_value = map[string]int32{
	"VALID":                        0,
	"INVALID_TX_PAYLOAD":           1,
	"INVALID_TX_HEADER":            2,
	"INVALID_TX_TYPE":              3,
	"UNKNOWN_TX_TYPE":              4,
	"DUPLICATE_TX_ID":              5,
	"INVALID_STATE_UPDATES_TYPE":   6,
	"INVALID_INSERT_STATEMENTS":    7,
	"INVALID_SELECT_STATEMENTS":    8,
	"INVALID_SQL_SYNTAX":           9,
	"INVALID_SQL_SCHEMA":           10,
	"INVALID_COMMON_TX_DATA":       11,
	"INVALID_VOTE_TX_DATA":         12,
	"INVALID_CHAIN_CONFIG":         13,
	"INVALID_CONTRACT_INVOCATION":  21,
	"INVALID_APPROVAL_SIGNATURE":   22,
	"APPROVALS_VALIDATION_FAILURE": 23,
	"UNKNOWN_VOTE_HANDLER":         31,
	"INVALID_VOTE_PAYLOAD":         32,
	"INVALID_MVCC":                 41,
	"INVALID":                      255,
}

func (x TxStatus) String() string {
	return proto.EnumName(TxStatus_name, int32(x))
}

func (TxStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{1}
}

type TxHeader struct {
	Type                 TxType   `protobuf:"varint,1,opt,name=type,proto3,enum=common.TxType" json:"type,omitempty"`
	TxId                 string   `protobuf:"bytes,2,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	ChainId              string   `protobuf:"bytes,3,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxHeader) Reset()         { *m = TxHeader{} }
func (m *TxHeader) String() string { return proto.CompactTextString(m) }
func (*TxHeader) ProtoMessage()    {}
func (*TxHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{0}
}

func (m *TxHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxHeader.Unmarshal(m, b)
}
func (m *TxHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxHeader.Marshal(b, m, deterministic)
}
func (m *TxHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxHeader.Merge(m, src)
}
func (m *TxHeader) XXX_Size() int {
	return xxx_messageInfo_TxHeader.Size(m)
}
func (m *TxHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_TxHeader.DiscardUnknown(m)
}

var xxx_messageInfo_TxHeader proto.InternalMessageInfo

func (m *TxHeader) GetType() TxType {
	if m != nil {
		return m.Type
	}
	return TxType_COMMON_TRANSACTION
}

func (m *TxHeader) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *TxHeader) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

type Transaction struct {
	Payload              []byte            `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Approvals            []*Approval       `protobuf:"bytes,2,rep,name=approvals,proto3" json:"approvals,omitempty"`
	Extensions           map[string][]byte `protobuf:"bytes,3,rep,name=extensions,proto3" json:"extensions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{1}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Transaction) GetApprovals() []*Approval {
	if m != nil {
		return m.Approvals
	}
	return nil
}

func (m *Transaction) GetExtensions() map[string][]byte {
	if m != nil {
		return m.Extensions
	}
	return nil
}

type TxPayload struct {
	Header               *TxHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Data                 []byte    `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *TxPayload) Reset()         { *m = TxPayload{} }
func (m *TxPayload) String() string { return proto.CompactTextString(m) }
func (*TxPayload) ProtoMessage()    {}
func (*TxPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{2}
}

func (m *TxPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxPayload.Unmarshal(m, b)
}
func (m *TxPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxPayload.Marshal(b, m, deterministic)
}
func (m *TxPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxPayload.Merge(m, src)
}
func (m *TxPayload) XXX_Size() int {
	return xxx_messageInfo_TxPayload.Size(m)
}
func (m *TxPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_TxPayload.DiscardUnknown(m)
}

var xxx_messageInfo_TxPayload proto.InternalMessageInfo

func (m *TxPayload) GetHeader() *TxHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *TxPayload) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// if this is a config update, the handler is "config"
// and payload is Chainconfig
type VoteTxData struct {
	Handler              string   `protobuf:"bytes,1,opt,name=handler,proto3" json:"handler,omitempty"`
	Subject              string   `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Payload              []byte   `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteTxData) Reset()         { *m = VoteTxData{} }
func (m *VoteTxData) String() string { return proto.CompactTextString(m) }
func (*VoteTxData) ProtoMessage()    {}
func (*VoteTxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{3}
}

func (m *VoteTxData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteTxData.Unmarshal(m, b)
}
func (m *VoteTxData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteTxData.Marshal(b, m, deterministic)
}
func (m *VoteTxData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteTxData.Merge(m, src)
}
func (m *VoteTxData) XXX_Size() int {
	return xxx_messageInfo_VoteTxData.Size(m)
}
func (m *VoteTxData) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteTxData.DiscardUnknown(m)
}

var xxx_messageInfo_VoteTxData proto.InternalMessageInfo

func (m *VoteTxData) GetHandler() string {
	if m != nil {
		return m.Handler
	}
	return ""
}

func (m *VoteTxData) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *VoteTxData) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type CommonTxData struct {
	ContractInvocation   []byte              `protobuf:"bytes,1,opt,name=contractInvocation,proto3" json:"contractInvocation,omitempty"`
	Response             *InvocationResponse `protobuf:"bytes,2,opt,name=response,proto3" json:"response,omitempty"`
	StateUpdates         []*StateUpdates     `protobuf:"bytes,3,rep,name=stateUpdates,proto3" json:"stateUpdates,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CommonTxData) Reset()         { *m = CommonTxData{} }
func (m *CommonTxData) String() string { return proto.CompactTextString(m) }
func (*CommonTxData) ProtoMessage()    {}
func (*CommonTxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{4}
}

func (m *CommonTxData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonTxData.Unmarshal(m, b)
}
func (m *CommonTxData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonTxData.Marshal(b, m, deterministic)
}
func (m *CommonTxData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonTxData.Merge(m, src)
}
func (m *CommonTxData) XXX_Size() int {
	return xxx_messageInfo_CommonTxData.Size(m)
}
func (m *CommonTxData) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonTxData.DiscardUnknown(m)
}

var xxx_messageInfo_CommonTxData proto.InternalMessageInfo

func (m *CommonTxData) GetContractInvocation() []byte {
	if m != nil {
		return m.ContractInvocation
	}
	return nil
}

func (m *CommonTxData) GetResponse() *InvocationResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *CommonTxData) GetStateUpdates() []*StateUpdates {
	if m != nil {
		return m.StateUpdates
	}
	return nil
}

type ContractInvocation struct {
	ContractName         string   `protobuf:"bytes,1,opt,name=contract_name,json=contractName,proto3" json:"contract_name,omitempty"`
	FuncName             string   `protobuf:"bytes,2,opt,name=func_name,json=funcName,proto3" json:"func_name,omitempty"`
	Args                 [][]byte `protobuf:"bytes,3,rep,name=args,proto3" json:"args,omitempty"`
	EncryptedKey         []byte   `protobuf:"bytes,4,opt,name=encryptedKey,proto3" json:"encryptedKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContractInvocation) Reset()         { *m = ContractInvocation{} }
func (m *ContractInvocation) String() string { return proto.CompactTextString(m) }
func (*ContractInvocation) ProtoMessage()    {}
func (*ContractInvocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{5}
}

func (m *ContractInvocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContractInvocation.Unmarshal(m, b)
}
func (m *ContractInvocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContractInvocation.Marshal(b, m, deterministic)
}
func (m *ContractInvocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContractInvocation.Merge(m, src)
}
func (m *ContractInvocation) XXX_Size() int {
	return xxx_messageInfo_ContractInvocation.Size(m)
}
func (m *ContractInvocation) XXX_DiscardUnknown() {
	xxx_messageInfo_ContractInvocation.DiscardUnknown(m)
}

var xxx_messageInfo_ContractInvocation proto.InternalMessageInfo

func (m *ContractInvocation) GetContractName() string {
	if m != nil {
		return m.ContractName
	}
	return ""
}

func (m *ContractInvocation) GetFuncName() string {
	if m != nil {
		return m.FuncName
	}
	return ""
}

func (m *ContractInvocation) GetArgs() [][]byte {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *ContractInvocation) GetEncryptedKey() []byte {
	if m != nil {
		return m.EncryptedKey
	}
	return nil
}

type InvocationResponse struct {
	Status               Status   `protobuf:"varint,1,opt,name=status,proto3,enum=common.Status" json:"status,omitempty"`
	StatusInfo           string   `protobuf:"bytes,2,opt,name=status_info,json=statusInfo,proto3" json:"status_info,omitempty"`
	Payload              []byte   `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InvocationResponse) Reset()         { *m = InvocationResponse{} }
func (m *InvocationResponse) String() string { return proto.CompactTextString(m) }
func (*InvocationResponse) ProtoMessage()    {}
func (*InvocationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{6}
}

func (m *InvocationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvocationResponse.Unmarshal(m, b)
}
func (m *InvocationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvocationResponse.Marshal(b, m, deterministic)
}
func (m *InvocationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvocationResponse.Merge(m, src)
}
func (m *InvocationResponse) XXX_Size() int {
	return xxx_messageInfo_InvocationResponse.Size(m)
}
func (m *InvocationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InvocationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InvocationResponse proto.InternalMessageInfo

func (m *InvocationResponse) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_UNKNOWN
}

func (m *InvocationResponse) GetStatusInfo() string {
	if m != nil {
		return m.StatusInfo
	}
	return ""
}

func (m *InvocationResponse) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type Approval struct {
	Cert                 []byte   `protobuf:"bytes,1,opt,name=cert,proto3" json:"cert,omitempty"`
	Sign                 []byte   `protobuf:"bytes,2,opt,name=sign,proto3" json:"sign,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Approval) Reset()         { *m = Approval{} }
func (m *Approval) String() string { return proto.CompactTextString(m) }
func (*Approval) ProtoMessage()    {}
func (*Approval) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{7}
}

func (m *Approval) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Approval.Unmarshal(m, b)
}
func (m *Approval) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Approval.Marshal(b, m, deterministic)
}
func (m *Approval) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Approval.Merge(m, src)
}
func (m *Approval) XXX_Size() int {
	return xxx_messageInfo_Approval.Size(m)
}
func (m *Approval) XXX_DiscardUnknown() {
	xxx_messageInfo_Approval.DiscardUnknown(m)
}

var xxx_messageInfo_Approval proto.InternalMessageInfo

func (m *Approval) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *Approval) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

type TxResult struct {
	TxId                 string   `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	Status               TxStatus `protobuf:"varint,2,opt,name=status,proto3,enum=common.TxStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxResult) Reset()         { *m = TxResult{} }
func (m *TxResult) String() string { return proto.CompactTextString(m) }
func (*TxResult) ProtoMessage()    {}
func (*TxResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{8}
}

func (m *TxResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxResult.Unmarshal(m, b)
}
func (m *TxResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxResult.Marshal(b, m, deterministic)
}
func (m *TxResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxResult.Merge(m, src)
}
func (m *TxResult) XXX_Size() int {
	return xxx_messageInfo_TxResult.Size(m)
}
func (m *TxResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TxResult.DiscardUnknown(m)
}

var xxx_messageInfo_TxResult proto.InternalMessageInfo

func (m *TxResult) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *TxResult) GetStatus() TxStatus {
	if m != nil {
		return m.Status
	}
	return TxStatus_VALID
}

type BlockResult struct {
	BlockNum             uint64      `protobuf:"varint,1,opt,name=blockNum,proto3" json:"blockNum,omitempty"`
	TxResults            []*TxResult `protobuf:"bytes,2,rep,name=txResults,proto3" json:"txResults,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *BlockResult) Reset()         { *m = BlockResult{} }
func (m *BlockResult) String() string { return proto.CompactTextString(m) }
func (*BlockResult) ProtoMessage()    {}
func (*BlockResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{9}
}

func (m *BlockResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockResult.Unmarshal(m, b)
}
func (m *BlockResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockResult.Marshal(b, m, deterministic)
}
func (m *BlockResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockResult.Merge(m, src)
}
func (m *BlockResult) XXX_Size() int {
	return xxx_messageInfo_BlockResult.Size(m)
}
func (m *BlockResult) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockResult.DiscardUnknown(m)
}

var xxx_messageInfo_BlockResult proto.InternalMessageInfo

func (m *BlockResult) GetBlockNum() uint64 {
	if m != nil {
		return m.BlockNum
	}
	return 0
}

func (m *BlockResult) GetTxResults() []*TxResult {
	if m != nil {
		return m.TxResults
	}
	return nil
}

type KeyModification struct {
	Value                []byte   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	BlockNum             uint64   `protobuf:"varint,2,opt,name=blockNum,proto3" json:"blockNum,omitempty"`
	TxNum                int32    `protobuf:"varint,3,opt,name=txNum,proto3" json:"txNum,omitempty"`
	TxID                 string   `protobuf:"bytes,4,opt,name=txID,proto3" json:"txID,omitempty"`
	IsDeleted            bool     `protobuf:"varint,5,opt,name=isDeleted,proto3" json:"isDeleted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeyModification) Reset()         { *m = KeyModification{} }
func (m *KeyModification) String() string { return proto.CompactTextString(m) }
func (*KeyModification) ProtoMessage()    {}
func (*KeyModification) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6296da495f91d72, []int{10}
}

func (m *KeyModification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyModification.Unmarshal(m, b)
}
func (m *KeyModification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyModification.Marshal(b, m, deterministic)
}
func (m *KeyModification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyModification.Merge(m, src)
}
func (m *KeyModification) XXX_Size() int {
	return xxx_messageInfo_KeyModification.Size(m)
}
func (m *KeyModification) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyModification.DiscardUnknown(m)
}

var xxx_messageInfo_KeyModification proto.InternalMessageInfo

func (m *KeyModification) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *KeyModification) GetBlockNum() uint64 {
	if m != nil {
		return m.BlockNum
	}
	return 0
}

func (m *KeyModification) GetTxNum() int32 {
	if m != nil {
		return m.TxNum
	}
	return 0
}

func (m *KeyModification) GetTxID() string {
	if m != nil {
		return m.TxID
	}
	return ""
}

func (m *KeyModification) GetIsDeleted() bool {
	if m != nil {
		return m.IsDeleted
	}
	return false
}

func init() {
	proto.RegisterEnum("common.TxType", TxType_name, TxType_value)
	proto.RegisterEnum("common.TxStatus", TxStatus_name, TxStatus_value)
	proto.RegisterType((*TxHeader)(nil), "common.TxHeader")
	proto.RegisterType((*Transaction)(nil), "common.Transaction")
	proto.RegisterMapType((map[string][]byte)(nil), "common.Transaction.ExtensionsEntry")
	proto.RegisterType((*TxPayload)(nil), "common.TxPayload")
	proto.RegisterType((*VoteTxData)(nil), "common.VoteTxData")
	proto.RegisterType((*CommonTxData)(nil), "common.CommonTxData")
	proto.RegisterType((*ContractInvocation)(nil), "common.ContractInvocation")
	proto.RegisterType((*InvocationResponse)(nil), "common.InvocationResponse")
	proto.RegisterType((*Approval)(nil), "common.Approval")
	proto.RegisterType((*TxResult)(nil), "common.TxResult")
	proto.RegisterType((*BlockResult)(nil), "common.BlockResult")
	proto.RegisterType((*KeyModification)(nil), "common.KeyModification")
}

func init() {
	proto.RegisterFile("common/transaction.proto", fileDescriptor_f6296da495f91d72)
}

var fileDescriptor_f6296da495f91d72 = []byte{
	// 1037 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0xdb, 0x6e, 0xdb, 0x46,
	0x10, 0x0d, 0x75, 0xb3, 0x34, 0x52, 0x12, 0x76, 0xad, 0xb8, 0x8a, 0x92, 0xd6, 0x02, 0x0d, 0x14,
	0x6a, 0x1e, 0x64, 0xc0, 0x45, 0x83, 0xa0, 0x40, 0x1f, 0x18, 0x92, 0x89, 0x09, 0x4b, 0x94, 0xba,
	0xa4, 0x55, 0x3b, 0x28, 0x4a, 0xac, 0xc9, 0xb5, 0xcd, 0x46, 0x22, 0x05, 0x72, 0x65, 0x4b, 0xdf,
	0x50, 0xa0, 0xcf, 0xfd, 0x8c, 0xfe, 0x4f, 0x3f, 0xa6, 0xc5, 0x2e, 0x49, 0x89, 0xb2, 0x8b, 0xbe,
	0xed, 0x9c, 0x33, 0xb3, 0x7b, 0xe6, 0x0c, 0x47, 0x82, 0x8e, 0x17, 0xcd, 0xe7, 0x51, 0x78, 0xcc,
	0x62, 0x12, 0x26, 0xc4, 0x63, 0x41, 0x14, 0x0e, 0x16, 0x71, 0xc4, 0x22, 0x54, 0x4b, 0x99, 0xee,
	0x7e, 0x96, 0x31, 0xa3, 0xfe, 0x0d, 0x8d, 0x53, 0xb2, 0xdb, 0xce, 0xc0, 0x39, 0x4d, 0x12, 0x72,
	0x43, 0x53, 0x54, 0xf9, 0x15, 0xea, 0xce, 0xea, 0x94, 0x12, 0x9f, 0xc6, 0x48, 0x81, 0x0a, 0x5b,
	0x2f, 0x68, 0x47, 0xea, 0x49, 0xfd, 0x67, 0x27, 0xcf, 0x06, 0x69, 0xc1, 0xc0, 0x59, 0x39, 0xeb,
	0x05, 0xc5, 0x82, 0x43, 0xfb, 0x50, 0x65, 0x2b, 0x37, 0xf0, 0x3b, 0xa5, 0x9e, 0xd4, 0x6f, 0xe0,
	0x0a, 0x5b, 0x99, 0x3e, 0x7a, 0x09, 0x75, 0xef, 0x96, 0x04, 0x21, 0xc7, 0xcb, 0x02, 0xdf, 0x13,
	0xb1, 0xe9, 0x2b, 0x7f, 0x4b, 0xd0, 0x74, 0xb6, 0x42, 0x51, 0x07, 0xf6, 0x16, 0x64, 0x3d, 0x8b,
	0x88, 0x2f, 0x9e, 0x69, 0xe1, 0x3c, 0x44, 0x03, 0x68, 0x90, 0xc5, 0x22, 0x8e, 0xee, 0xc8, 0x2c,
	0xe9, 0x94, 0x7a, 0xe5, 0x7e, 0xf3, 0x44, 0xce, 0x25, 0xa8, 0x19, 0x81, 0xb7, 0x29, 0x48, 0x03,
	0xa0, 0x2b, 0x46, 0xc3, 0x24, 0x88, 0xc2, 0xa4, 0x53, 0x16, 0x05, 0x47, 0x1b, 0xcd, 0x05, 0x6f,
	0x8c, 0x4d, 0x96, 0x11, 0xb2, 0x78, 0x8d, 0x0b, 0x65, 0xdd, 0x1f, 0xe1, 0xf9, 0x03, 0x1a, 0xc9,
	0x50, 0xfe, 0x4c, 0xd7, 0x42, 0x5d, 0x03, 0xf3, 0x23, 0x6a, 0x43, 0xf5, 0x8e, 0xcc, 0x96, 0x54,
	0xf4, 0xdc, 0xc2, 0x69, 0xf0, 0x43, 0xe9, 0x9d, 0xa4, 0x98, 0xd0, 0x70, 0x56, 0x93, 0xac, 0x81,
	0x3e, 0xd4, 0x6e, 0x85, 0x91, 0xa2, 0xb6, 0xa0, 0x3e, 0x37, 0x18, 0x67, 0x3c, 0x42, 0x50, 0xf1,
	0x09, 0x23, 0xd9, 0x7d, 0xe2, 0xac, 0x7c, 0x02, 0x98, 0x46, 0x8c, 0x3a, 0x2b, 0x9d, 0x30, 0xc2,
	0x6d, 0xba, 0x25, 0xa1, 0x3f, 0xcb, 0x2e, 0x6b, 0xe0, 0x3c, 0xe4, 0x4c, 0xb2, 0xbc, 0xfa, 0x8d,
	0x7a, 0x2c, 0x1b, 0x41, 0x1e, 0x16, 0xad, 0x2d, 0xef, 0x58, 0xab, 0xfc, 0x25, 0x41, 0x4b, 0x13,
	0x5a, 0xb2, 0xeb, 0x07, 0x80, 0xbc, 0x28, 0x64, 0x31, 0xf1, 0x98, 0x19, 0xde, 0x45, 0x1e, 0xe1,
	0x46, 0x65, 0x03, 0xf9, 0x0f, 0x06, 0xbd, 0x85, 0x7a, 0x4c, 0x93, 0x45, 0x14, 0x26, 0xa9, 0x09,
	0xcd, 0x93, 0x6e, 0xde, 0xdc, 0x36, 0x0b, 0x67, 0x19, 0x78, 0x93, 0x8b, 0xde, 0x41, 0x2b, 0x61,
	0x84, 0xd1, 0xf3, 0x85, 0x4f, 0x18, 0xcd, 0xa7, 0xd4, 0xce, 0x6b, 0xed, 0x02, 0x87, 0x77, 0x32,
	0x95, 0x3f, 0x24, 0x40, 0xda, 0x63, 0x21, 0x47, 0xf0, 0x34, 0x97, 0xe7, 0x86, 0x64, 0x4e, 0x33,
	0x77, 0x5a, 0x39, 0x68, 0x91, 0x39, 0x45, 0xaf, 0xa0, 0x71, 0xbd, 0x0c, 0xbd, 0x34, 0x21, 0x35,
	0xa9, 0xce, 0x01, 0x41, 0x22, 0xa8, 0x90, 0xf8, 0x26, 0x95, 0xd2, 0xc2, 0xe2, 0x8c, 0x14, 0x68,
	0xd1, 0xd0, 0x8b, 0xd7, 0x0b, 0x46, 0xfd, 0x33, 0xba, 0xee, 0x54, 0x84, 0x11, 0x3b, 0x98, 0x72,
	0x0f, 0xe8, 0x71, 0xab, 0xe8, 0x1b, 0xa8, 0x71, 0xd9, 0xcb, 0xe4, 0xe1, 0xd2, 0xd8, 0x02, 0xc5,
	0x19, 0x8b, 0x0e, 0xa1, 0x99, 0x9e, 0xdc, 0x20, 0xbc, 0x8e, 0x32, 0x51, 0x90, 0x42, 0x66, 0x78,
	0x1d, 0xfd, 0xcf, 0xf0, 0x4e, 0xa0, 0x9e, 0x7f, 0xfe, 0x5c, 0xbc, 0x47, 0x63, 0x96, 0x4d, 0x4a,
	0x9c, 0x39, 0x96, 0x04, 0x37, 0x61, 0xfe, 0x31, 0xf1, 0xb3, 0x62, 0xf2, 0xad, 0xc6, 0x34, 0x59,
	0xce, 0xd8, 0x76, 0x63, 0xa5, 0xc2, 0xc6, 0xf6, 0x37, 0xba, 0x4b, 0x42, 0x77, 0xe1, 0x5b, 0xdd,
	0x55, 0xae, 0x5c, 0x42, 0xf3, 0xfd, 0x2c, 0xf2, 0x3e, 0x67, 0xb7, 0x75, 0xa1, 0x7e, 0xc5, 0x43,
	0x6b, 0x39, 0x17, 0x17, 0x56, 0xf0, 0x26, 0xe6, 0x1b, 0xcc, 0xb2, 0x57, 0x1f, 0x6d, 0x70, 0x2e,
	0x07, 0x6f, 0x53, 0x94, 0xdf, 0x25, 0x78, 0x7e, 0x46, 0xd7, 0xa3, 0xc8, 0x0f, 0xae, 0x83, 0x6c,
	0xc0, 0x9b, 0x5d, 0x93, 0x0a, 0xbb, 0xb6, 0xf3, 0x6a, 0xe9, 0xc1, 0xab, 0x6d, 0xde, 0x1f, 0x27,
	0xb8, 0x6f, 0x55, 0x9c, 0x06, 0xdc, 0x15, 0xb6, 0x32, 0x75, 0x31, 0xca, 0xb4, 0x69, 0x1d, 0xbd,
	0x86, 0x46, 0x90, 0xe8, 0x74, 0x46, 0x19, 0xf5, 0x3b, 0xd5, 0x9e, 0xd4, 0xaf, 0xe3, 0x2d, 0xf0,
	0xe6, 0x2d, 0xd4, 0xd2, 0x5f, 0x3a, 0x74, 0x00, 0x48, 0x1b, 0x8f, 0x46, 0x63, 0xcb, 0x75, 0xb0,
	0x6a, 0xd9, 0xaa, 0xe6, 0x98, 0x63, 0x4b, 0x7e, 0x82, 0xda, 0x20, 0x4f, 0xc7, 0x8e, 0xb1, 0x83,
	0x4a, 0x6f, 0xfe, 0xac, 0x70, 0xb3, 0x53, 0xd7, 0x50, 0x03, 0xaa, 0x53, 0x75, 0x68, 0xea, 0xf2,
	0x13, 0x7e, 0x8b, 0x69, 0x89, 0xc0, 0x75, 0x2e, 0xdc, 0x89, 0x7a, 0x39, 0x1c, 0xab, 0xba, 0x2c,
	0xa1, 0x17, 0xf0, 0x45, 0x01, 0x3f, 0x35, 0x54, 0xdd, 0xc0, 0x72, 0x09, 0xed, 0xc3, 0xf3, 0x02,
	0xec, 0x5c, 0x4e, 0x0c, 0xb9, 0xcc, 0xc1, 0x73, 0xeb, 0xcc, 0x1a, 0xff, 0x6c, 0x6d, 0xc0, 0x0a,
	0x07, 0xf5, 0xf3, 0xc9, 0xd0, 0xd4, 0x54, 0xae, 0xe5, 0xc2, 0x35, 0x75, 0xb9, 0x8a, 0xbe, 0x86,
	0x6e, 0x5e, 0x6e, 0x3b, 0x9c, 0x38, 0x9f, 0xe8, 0xaa, 0x63, 0xd8, 0x69, 0x51, 0x0d, 0x7d, 0x05,
	0x2f, 0x73, 0xde, 0xb4, 0x6c, 0x03, 0x3b, 0x69, 0xda, 0xc8, 0xb0, 0x1c, 0x5b, 0xde, 0x2b, 0xd2,
	0xb6, 0x31, 0x34, 0xb4, 0x1d, 0xba, 0x5e, 0xec, 0xc5, 0xfe, 0x69, 0xe8, 0xda, 0x97, 0x96, 0xa3,
	0x5e, 0xc8, 0x8d, 0x47, 0xb8, 0x76, 0x6a, 0x8c, 0x54, 0x19, 0x50, 0x17, 0x0e, 0x72, 0x3c, 0x77,
	0xf2, 0xc2, 0xd5, 0x55, 0x47, 0x95, 0x9b, 0xa8, 0x03, 0xed, 0x9c, 0x4b, 0xdd, 0xcc, 0x98, 0x56,
	0x91, 0xd1, 0x4e, 0x55, 0xd3, 0x72, 0xb5, 0xb1, 0xf5, 0xc1, 0xfc, 0x28, 0x3f, 0x45, 0x87, 0xf0,
	0x6a, 0x7b, 0x9f, 0xe5, 0x60, 0x55, 0x73, 0x5c, 0xd3, 0x9a, 0x8e, 0x35, 0x55, 0x0c, 0xe1, 0x45,
	0xb1, 0x7d, 0x75, 0x32, 0xc1, 0xe3, 0xa9, 0x3a, 0x74, 0x6d, 0xf3, 0xa3, 0xa5, 0x3a, 0xe7, 0xd8,
	0x90, 0x0f, 0x50, 0x0f, 0x5e, 0xe7, 0xb8, 0xed, 0x8a, 0x3c, 0x51, 0xe9, 0x7e, 0x50, 0xcd, 0x21,
	0xcf, 0xf8, 0x92, 0x3f, 0x9e, 0x5b, 0x2d, 0x64, 0x9d, 0xaa, 0x96, 0x3e, 0x34, 0xb0, 0x7c, 0xf8,
	0x48, 0x70, 0x3e, 0xca, 0x1e, 0x92, 0xa1, 0x95, 0x33, 0xa3, 0xa9, 0xa6, 0xc9, 0xdf, 0xa2, 0x16,
	0xec, 0x65, 0x88, 0xfc, 0x8f, 0xf4, 0xfe, 0x17, 0x38, 0xf2, 0xa2, 0xf9, 0xe0, 0x76, 0x49, 0xee,
	0x69, 0x30, 0xb8, 0x0f, 0x68, 0x48, 0x63, 0xf1, 0xc7, 0x98, 0xfe, 0xf9, 0x66, 0xbb, 0xf1, 0xe9,
	0xfb, 0x9b, 0x80, 0xe5, 0x49, 0x5e, 0x34, 0x3f, 0x5e, 0x44, 0x41, 0x92, 0xf0, 0x9f, 0x17, 0x12,
	0x7b, 0xb7, 0xc7, 0x85, 0xb2, 0x63, 0x51, 0x76, 0x9c, 0x96, 0x5d, 0xd5, 0x44, 0xf4, 0xdd, 0xbf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xb0, 0xd3, 0x84, 0x02, 0x10, 0x08, 0x00, 0x00,
}
