/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

package sdk

import (
	"io"

	"git.huawei.com/poissonsearch/wienerchain/proto/common"
)

// Context is the context for smart contract stub.
type Context interface {
	GetStub(txID string) ContractStub
	AddStub(stubInterface ContractStub) error
	DelStub(stubInterface ContractStub)
	GetContractInfo() Contract
	SetContractInfo(contract Contract)
}

// SQLContext is the context for SQL smart contract stub.
type SQLContext interface {
	GetStub(txID string) SQLContractStub
	AddStub(stubInterface SQLContractStub) error
	DelStub(stubInterface SQLContractStub)
	GetContractInfo() SQLContract
	SetContractInfo(contract SQLContract)
}

// Stub is common contract stub.
type Stub interface {
	// Get the invoke function name
	FuncName() string

	// Get the invoke parameters
	Parameters() [][]byte

	// Get the tx id for this transaction
	TxID() string

	// Get the chain id for this transaction
	ChainID() string

	// Get the contract id for this transaction
	ContractName() string

	GenerateStateUpdates() ([]*common.StateUpdates, error)
	GenerateSQLStateUpdates() ([]*common.StateUpdates, error)
	GenerateExtension() (map[string][]byte, error)

	io.Closer
}

// SQLContract is the SQL smart contract interface .
type SQLContract interface {
	InitSQL(stub SQLContractStub) common.InvocationResponse
	InvokeSQL(stub SQLContractStub) common.InvocationResponse
}

// SQLContractStub is the SQL smart contract stub .
type SQLContractStub interface {
	Stub
	Query(query string, args ...interface{}) (Rows, error)
	Exec(query string, args ...interface{}) error
}

// Rows is the SQL rows interface .
type Rows interface {
	Next() (bool, error)
	Scan(dest ...interface{}) error
	Close() error
}

// Contract is the leveldb interface .
type Contract interface {
	Init(stub ContractStub) common.InvocationResponse
	Invoke(stub ContractStub) common.InvocationResponse
}

// ContractStub is the leveldb stub interface .
type ContractStub interface {
	Stub
	// Get value by key from state DB
	GetKV(key string) ([]byte, error)

	// This action will only generate read write set,
	// the key value will not be put into stateDB until the transaction is validated.
	PutKV(key string, value []byte) error

	// The value should be base type or implement the Marshal(v interface{}) ([]byte, error)
	// and Unmarshal(data []byte, v interface{}) error interface
	PutKVCommon(key string, value interface{}) error

	// This action will only generate read write set, the key value will not be del until the transaction is validated.
	DelKV(key string) error

	// Get a K、V iterator for startKey and endKey. Users do not need to perceive the internal buffering process,
	// even if the result are too large
	// [stareKey, endKey) for example: 11--13, you will get 11, 12, not including 13
	// for example: 11--11,you will get nothing
	GetIterator(startKey, endKey string) (Iterator, error)
	// Get history iterator of the key
	GetKeyHistoryIterator(key string) (HistoryIterator, error)
	// 以下两个接口，作为GenerateComKey、DivideComKey、GetComKeyIterator的替代方案
	// Save composite index for objectKey.
	// The mark is the common mark for the object index.
	// The attributes is the attributes for the object value
	// The objectKey is the key for the origin object
	// 直接在此接口内部，完成索引的生成和保存，相当于GenerateComKey和putKV的组合使用
	SaveComIndex(indexName string, attributes []string, objectKey string) error

	// Get the object value by composite index
	// We can get the index by mark and attributes, then we can get object key and value by index
	GetKVByComIndex(indexName string, attributes []string) (Iterator, error)

	// Delete the index for certain object key
	DelComIndexOneRow(indexName string, attributes []string, objectKey string) error
}

// Iterator  interface for user to get key and value one by one
type Iterator interface {
	Next() bool
	Key() string
	Value() []byte
	Close()
}

// Iterator  interface for user to get key history one by one
type HistoryIterator interface {
	Iterator
	// returns BlockNum & TxNum
	Version() (uint64, int32)
	TxID() string
	// whether the key has been deleted
	IsDeleted() bool
}

// ValueSerialization when using "PutKVCommon(key string, value interface{}) error",
// the parameter "value " should be implemented with ValueSerialization interface
type ValueSerialization interface {
	Marshal() ([]byte, error)
}
