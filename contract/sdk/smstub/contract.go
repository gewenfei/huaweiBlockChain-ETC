/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

// Package smstub is the stub for smart contract api
package smstub

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/big"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"

	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/utils"

	"git.huawei.com/poissonsearch/wienerchain/contract/sdk"
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/logger"
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/smcontext"
	"git.huawei.com/poissonsearch/wienerchain/proto/common"
	"git.huawei.com/poissonsearch/wienerchain/proto/contract"
	"github.com/golang/protobuf/proto"
)

const (
	minUnicodeValue       = 0
	maxUnicodeValue       = utf8.MaxRune
	compositeKeyNamespace = "\x00"
)

const (
	bufferSize    = 100
	timeDelay     = 20
	queryIDMaxNum = 10000
)

// TxExecution transaction exec struct.
type TxExecution struct {
	TransactionID  string
	ContractID     string
	ChainName      string
	Invocation     *common.ContractInvocation
	RecvChan       chan []byte
	StreamSendChan chan contract.ContractMsg
}

// QueryIterator iterator for query.
type QueryIterator struct {
	queryID  string
	beginKey string
	endKey   string
	Position int

	//When over buffer size, the point of query response should be changed to a new one
	queryRes  *contract.QueryResponse
	onceLen   int
	txContext *TxExecution
}

// QueryKeyHistoryIterator iterator for key history
type QueryKeyHistoryIterator struct {
	queryID   string
	key       string
	position  int
	queryRes  *contract.HistoryForKey
	onceLen   int
	txContext *TxExecution
}

// ComKeyIterator composite key iterator, we can get the actual data, not the index.
type ComKeyIterator struct {
	iterator     sdk.Iterator
	txContext    *TxExecution
	currentKey   string
	currentValue []byte
}

// QueryRows query rows struct.
type QueryRows struct {
	queryID  string
	query    string
	args     *common.PrimitiveValues
	position int

	//When over buffer size, the point of query response should be changed to a new one
	queryRes  *contract.SqlQueryResponse
	onceLen   int
	txContext *TxExecution
}

var ctlog = logger.GetDefaultLogger()

// getObjectKV
func (c *ComKeyIterator) getObjectKV() (string, []byte) {
	compositeKey := c.iterator.Key()

	ctlog.Debugf("In getObjectKV,  get compositeKey is %v", compositeKey)

	componentIndex := 1

	components := []string{}
	for i := 1; i < len(compositeKey); i++ {
		if compositeKey[i] == minUnicodeValue {
			components = append(components, compositeKey[componentIndex:i])
			componentIndex = i + 1
		}
	}

	ctlog.Debugf("The components is %v", components)

	objectKey := components[len(components)-1]

	c.currentKey = objectKey
	value, e := c.txContext.GetKV(objectKey)
	if e != nil {
		return "", nil
	}
	c.currentValue = value
	return objectKey, value
}

// Next judge is there next element for iterator.
func (c *ComKeyIterator) Next() bool {
	c.currentKey = ""
	c.currentValue = nil
	var hasNext bool
	var key string
	var value []byte

	for {
		hasNext = c.iterator.Next()
		if !hasNext {
			return false
		}
		key, value = c.getObjectKV()
		if value == nil {
			continue
		} else {
			c.currentKey = key
			c.currentValue = value
			return hasNext
		}
	}
}

// Key get key for iterator element.
func (c *ComKeyIterator) Key() string {
	if c.currentKey != "" {
		return c.currentKey
	}
	key, _ := c.getObjectKV()
	return key
}

// Value get value for iterator element.
func (c *ComKeyIterator) Value() []byte {
	if c.currentKey != "" {
		return c.currentValue
	}

	_, value := c.getObjectKV()
	return value
}

// Close close the iterator.
func (c ComKeyIterator) Close() {
	c.iterator.Close()
}

// Next judge is there next element for iterator.
func (qi *QueryIterator) Next() bool {
	if qi.onceLen <= 0 {
		return false
	}
	if qi.Position < qi.onceLen-1 {
		qi.Position++
		return true
	}
	if qi.queryRes.IsOver {
		return false
	}

	// Query from wnode
	err := qi.queryFromWnode(bufferSize)
	if err != nil {
		ctlog.Errorf("Query from wnode failed when get next")
		return false
	}
	if qi.onceLen == 0 {
		ctlog.Errorf("There is no data for next req")
		return false
	}
	qi.Position = 0

	return true
}

// Key get key for iterator element.
func (qi *QueryIterator) Key() string {
	return qi.queryRes.KvArray[qi.Position].Key
}

// Value get value for iterator element.
func (qi *QueryIterator) Value() []byte {
	return qi.queryRes.KvArray[qi.Position].Value
}

// Close close the iterator
func (qi *QueryIterator) Close() {
	err := qi.queryFromWnode(0)
	if err != nil {
		ctlog.Errorf("Query from wnode failed when close the iterator")
	}
}
func (qi *QueryKeyHistoryIterator) queryFromWnode(bufferSize uint32) error {
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = qi.txContext.TransactionID
	contractMsg.Contractid = qi.txContext.ContractID
	contractMsg.Chainid = qi.txContext.ChainName
	contractMsg.Type = contract.ContractMsg_GET_HISTORY_FOR_KEY

	queryHistoryInfo := &contract.GetHistoryForKey{}
	queryHistoryInfo.Key = qi.key
	queryHistoryInfo.QueryId = qi.queryID
	queryHistoryInfo.BufferSize = bufferSize

	queryHistoryBytes, err := proto.Marshal(queryHistoryInfo)
	if err != nil {
		ctlog.Errorf("marshal fail!")
		return err
	}

	contractMsg.Payload = queryHistoryBytes
	payload, err := qi.txContext.CallWnode(contractMsg)
	if err != nil {
		ctlog.Errorf("call wnode failed when Get KeyHistoryIterator.")
		return err
	}

	queryRes := &contract.HistoryForKey{}
	err = proto.Unmarshal(payload, queryRes)
	if err != nil {
		ctlog.Errorf("failed to unmarshal query response.")
		return err
	}

	qi.position = 0
	qi.queryRes = queryRes
	qi.onceLen = len(queryRes.HistoryArray)

	ctlog.Debugf("the once len is %d", qi.onceLen)
	return nil
}

// Next judge is there next element for iterator.
func (qi *QueryKeyHistoryIterator) Next() bool {
	if qi.onceLen <= 0 {
		return false
	}
	if qi.position < qi.onceLen-1 {
		qi.position++
		return true
	}
	if qi.queryRes.IsOver {
		return false
	}

	// Query from wnode
	err := qi.queryFromWnode(bufferSize)
	if err != nil {
		ctlog.Errorf("query from wnode failed when get next")
		return false
	}
	if qi.onceLen == 0 {
		ctlog.Errorf("there is no data for next req")
		return false
	}
	qi.position = 0

	return true
}

// Key get history key.
func (qi *QueryKeyHistoryIterator) Key() string {
	return qi.key
}

// Value get history Value.
func (qi *QueryKeyHistoryIterator) Value() []byte {
	return qi.queryRes.HistoryArray[qi.position].Value
}

// Close get history Close.
func (qi *QueryKeyHistoryIterator) Close() {
	err := qi.queryFromWnode(0)
	if err != nil {
		ctlog.Errorf("query from wnode failed when close the iterator")
	}
}

// Version get history Version.
func (qi *QueryKeyHistoryIterator) Version() (uint64, int32) {
	return qi.queryRes.HistoryArray[qi.position].BlockNum, qi.queryRes.HistoryArray[qi.position].TxNum
}

// TxID get history TxID.
func (qi *QueryKeyHistoryIterator) TxID() string {
	return qi.queryRes.HistoryArray[qi.position].TxID
}

// IsDeleted get history IsDeleted.
func (qi *QueryKeyHistoryIterator) IsDeleted() bool {
	return qi.queryRes.HistoryArray[qi.position].IsDeleted
}

func (qi *QueryIterator) queryFromWnode(bufferSize uint32) error {
	// The default buffer size is 100
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = qi.txContext.TransactionID
	contractMsg.Contractid = qi.txContext.ContractID
	contractMsg.Chainid = qi.txContext.ChainName
	contractMsg.Type = contract.ContractMsg_GET_KV_BY_RANGE
	queryInfo := &contract.GetStateByRange{}
	queryInfo.BeginKey = qi.beginKey
	queryInfo.EndKey = qi.endKey
	queryInfo.QueryId = qi.queryID
	queryInfo.BufferSize = bufferSize
	queryBytes, err := proto.Marshal(queryInfo)
	if err != nil {
		return err
	}
	contractMsg.Payload = queryBytes
	queryResBytes, err := qi.txContext.CallWnode(contractMsg)
	if err != nil {
		return err
	}
	queryRes := &contract.QueryResponse{}
	err = proto.Unmarshal(queryResBytes, queryRes)
	if err != nil {
		return err
	}
	qi.Position = 0
	qi.queryRes = queryRes
	qi.onceLen = len(queryRes.KvArray)
	ctlog.Debugf("In queryFromWnode, the once len is %d", qi.onceLen)

	return nil
}

// GetTxID get transaction ID.
func (te *TxExecution) GetTxID() string {
	return te.TransactionID
}

// FuncName the function name for this tx.
func (te *TxExecution) FuncName() string {
	return te.Invocation.FuncName
}

// Parameters the parameter for this tx
func (te *TxExecution) Parameters() [][]byte {
	return te.Invocation.Args
}

// GetKV get the value for certain key.
func (te *TxExecution) GetKV(key string) ([]byte, error) {
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = te.TransactionID
	contractMsg.Contractid = te.ContractID
	contractMsg.Chainid = te.ChainName
	contractMsg.Type = contract.ContractMsg_GET_KV
	getstateInfo := &contract.GetState{}
	getstateInfo.Key = key
	payload, err := proto.Marshal(getstateInfo)
	if err != nil {
		ctlog.Errorf("marshal failed: %s", err)
		return nil, err
	}
	contractMsg.Payload = payload

	value, err := te.CallWnode(contractMsg)
	if err != nil {
		return nil, fmt.Errorf("get value failed")
	}
	return value, nil
}

// PutKVCommon the value should be base type or implement
func (te *TxExecution) PutKVCommon(key string, value interface{}) error {
	var valueBytes []byte

	switch v := value.(type) {
	case string:
		valueBytes = []byte(v)
	case []byte:
		valueBytes = v
	case int:
		i := v
		valueBytes = IntToBytes(i)
	case int32:
		i := v
		valueBytes = IntToBytes(int(i))
	case int64:
		i := v
		valueBytes = IntToBytes(int(i))
	case sdk.ValueSerialization:
		object := v
		valueMarshal, err := object.Marshal()
		if err != nil {
			return err
		}
		valueBytes = valueMarshal
	default:
		ctlog.Infof("Unknown type for put kv")
		return fmt.Errorf("unknown type for put kv")
	}

	return te.PutKV(key, valueBytes)
}

// IntToBytes change int to bytes.
func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	err := binary.Write(bytebuf, binary.BigEndian, data)
	if err != nil {
		return nil
	}
	return bytebuf.Bytes()
}

// PutKV put the key value to state DB.
func (te *TxExecution) PutKV(key string, value []byte) error {
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = te.TransactionID
	contractMsg.Contractid = te.ContractID
	contractMsg.Chainid = te.ChainName
	contractMsg.Type = contract.ContractMsg_PUT_KV
	putstateInfo := &contract.PutState{}
	putstateInfo.Key = key
	putstateInfo.Value = value

	payload, err := proto.Marshal(putstateInfo)
	if err != nil {
		return err
	}
	contractMsg.Payload = payload
	_, err = te.CallWnode(contractMsg)
	if err != nil {
		err1 := fmt.Errorf("putKV failed : %s", err.Error())
		return err1
	}
	return nil
}

// CallWnode call wnode to get state value.
func (te *TxExecution) CallWnode(msg *contract.ContractMsg) ([]byte, error) {
	te.StreamSendChan <- *msg
	for {
		select {
		case <-time.After(timeDelay * time.Second):
			ctlog.Infof("Time out for call wnode for tx:%s", msg.Txid)
			return nil, fmt.Errorf("time out for call wnode for tx:%s", msg.Txid)
		case info := <-te.RecvChan:

			return info, nil
		}
	}
}

// DelKV delete certain key value in state DB.
func (te *TxExecution) DelKV(key string) error {
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = te.TransactionID
	contractMsg.Contractid = te.ContractID
	contractMsg.Chainid = te.ChainName
	contractMsg.Type = contract.ContractMsg_DEL_KV
	deleteInfo := contract.DelState{}
	deleteInfo.Key = key
	payload, err := proto.Marshal(&deleteInfo)
	if err != nil {
		return err
	}
	contractMsg.Payload = payload
	if _, err := te.CallWnode(contractMsg); err != nil {
		ctlog.Errorf("Delete key failed: %s", err)
		return err
	}
	return nil
}

// TxID get transaction ID.
func (te *TxExecution) TxID() string {
	return te.TransactionID
}

// ChainID get chain ID.
func (te *TxExecution) ChainID() string {
	return te.ChainName
}

// ContractName get contract name.
func (te *TxExecution) ContractName() string {
	return te.ContractID
}

// Close sql close interface .
func (te *TxExecution) Close() error {
	panic("implement me")
}

// GenerateQueryID generate query ID
func GenerateQueryID() string {
	result, err := rand.Int(rand.Reader, big.NewInt(queryIDMaxNum))
	if err != nil {
		ctlog.Errorf("can not get a uniform random value")
		return ""
	}

	txStr := fmt.Sprintf("%d%s", time.Now().UnixNano(), result.String())

	h := sha256.New()
	_, err = h.Write([]byte(txStr))
	if err != nil {
		return txStr
	}
	hashInfo := h.Sum(nil)

	hexStr := fmt.Sprintf("%x", hashInfo)
	return hexStr
}

// GetIterator get iterator from begin key to end key.
func (te *TxExecution) GetIterator(beginKey, endKey string) (sdk.Iterator, error) {
	iterator, err := te.getIteratorWithBufferSize(beginKey, endKey, bufferSize)
	return iterator, err
}

// GetKeyHistoryIterator get history iterator from the key
func (te *TxExecution) GetKeyHistoryIterator(key string) (sdk.HistoryIterator, error) {
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = te.TransactionID
	contractMsg.Contractid = te.ContractID
	contractMsg.Chainid = te.ChainName
	contractMsg.Type = contract.ContractMsg_GET_HISTORY_FOR_KEY

	queryHistoryInfo := &contract.GetHistoryForKey{}
	queryHistoryInfo.Key = key
	queryHistoryInfo.QueryId = GenerateQueryID()
	queryHistoryInfo.BufferSize = bufferSize

	queryHistoryBytes, err := proto.Marshal(queryHistoryInfo)
	if err != nil {
		ctlog.Errorf("marshal fail!")
		return nil, err
	}

	contractMsg.Payload = queryHistoryBytes
	payload, err := te.CallWnode(contractMsg)
	if err != nil {
		ctlog.Errorf("call wnode failed when Get KeyHistoryIterator.")
		return nil, err
	}

	queryRes := &contract.HistoryForKey{}
	err = proto.Unmarshal(payload, queryRes)
	if err != nil {
		ctlog.Errorf("failed to unmarshal query response.")
		return nil, err
	}

	iterator := &QueryKeyHistoryIterator{}
	iterator.queryRes = queryRes
	iterator.key = key
	iterator.txContext = te
	iterator.onceLen = len(queryRes.HistoryArray)
	iterator.position = -1
	iterator.queryID = queryHistoryInfo.QueryId

	ctlog.Debugf("the once len is %d", iterator.onceLen)
	return iterator, nil
}

// getIteratorWithBufferSize get iterator with bufferSize
func (te *TxExecution) getIteratorWithBufferSize(begin, end string, bufferSize uint32) (sdk.Iterator, error) {
	// The default buffer size is 100
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = te.TransactionID
	contractMsg.Contractid = te.ContractID
	contractMsg.Chainid = te.ChainName
	contractMsg.Type = contract.ContractMsg_GET_KV_BY_RANGE
	QueryRangeInfo := &contract.GetStateByRange{}
	QueryRangeInfo.BeginKey = begin
	QueryRangeInfo.EndKey = end

	// Default buffer size
	QueryRangeInfo.BufferSize = bufferSize

	QueryRangeInfo.QueryId = GenerateQueryID()

	queryRangeBytes, err := proto.Marshal(QueryRangeInfo)
	if err != nil {
		return nil, err
	}
	contractMsg.Payload = queryRangeBytes

	payload, err := te.CallWnode(contractMsg)
	if err != nil {
		ctlog.Errorf("Call wnode failed when get state by range.")
		return nil, err
	}
	queryRes := &contract.QueryResponse{}
	err = proto.Unmarshal(payload, queryRes)
	if err != nil {
		ctlog.Errorf("Failed to unmarshal query response.")
		return nil, err
	}

	iterator := &QueryIterator{}
	iterator.queryID = QueryRangeInfo.QueryId
	iterator.queryRes = queryRes
	iterator.Position = -1
	iterator.txContext = te
	iterator.onceLen = len(queryRes.KvArray)
	ctlog.Debugf("The oncelen is %d", iterator.onceLen)

	return iterator, nil
}

// SaveComIndex save composite key index.
func (te *TxExecution) SaveComIndex(indexName string, attributes []string, objectKey string) error {
	index, err := GenerateIndex(indexName, attributes, objectKey)
	if err != nil {
		return err
	}
	err = te.PutKV(index, []byte("0"))
	if err != nil {
		return err
	}

	return nil
}

// DelComIndexOneRow delete ComIndexOneRow.
func (te *TxExecution) DelComIndexOneRow(indexName string, attributes []string, objectKey string) error {
	index, err := GenerateIndex(indexName, attributes, objectKey)
	if err != nil {
		ctlog.Infof("generate index failed:%s", err.Error())
		return err
	}
	err = te.DelKV(index)
	if err != nil {
		ctlog.Errorf("del key value failed:%s", err.Error())
		return err
	}

	return nil
}

// GenerateIndex generate composite key index.
func GenerateIndex(indexName string, attributes []string, objectKey string) (string, error) {
	ctlog.Debugf("The utf8.MaxRune string len is %d", len(string(utf8.MaxRune)))

	comKey := compositeKeyNamespace + indexName
	for _, v := range attributes {
		comKey = comKey + string(minUnicodeValue) + v
	}
	if objectKey != "" {
		comKey = comKey + string(minUnicodeValue) + objectKey + string(minUnicodeValue)
	}

	ctlog.Debugf("The generated com key is %s\n", comKey)

	return comKey, nil
}

// GetKVByComIndex get key value by composite key index. 看能否直接返回用户最终想要的信息的迭代，而不是复合键的迭代，方便用户使用
func (te *TxExecution) GetKVByComIndex(indexName string, attributes []string) (sdk.Iterator, error) {
	ctlog.Debugf("The utf8.MaxRune string len is %d", len(string(utf8.MaxRune)))
	comKey, err := GenerateIndex(indexName, attributes, "")
	if err != nil {
		return nil, err
	}
	beginKey := comKey
	endKey := comKey + string(utf8.MaxRune)
	iterator, e := te.GetIterator(beginKey, endKey)
	if e != nil {
		return nil, e
	}

	comKeyIte := &ComKeyIterator{}
	comKeyIte.txContext = te
	comKeyIte.iterator = iterator

	return comKeyIte, nil
}

// GenerateComKey generate composite key.
func (te *TxExecution) GenerateComKey(objectType string, attributes []string) (string, error) {
	return createCompositeKey(objectType, attributes)
}

// DivideComKey divide composite key.
func (te *TxExecution) DivideComKey(compositeKey string) (string, []string, error) {
	return splitCompositeKey(compositeKey)
}

func createCompositeKey(objectType string, attributes []string) (string, error) {
	if err := validateAttribute(objectType); err != nil {
		return "", err
	}
	ck := compositeKeyNamespace + objectType + string(minUnicodeValue)
	for _, att := range attributes {
		if err := validateAttribute(att); err != nil {
			return "", err
		}
		ck += att + string(minUnicodeValue)
	}
	return ck, nil
}

func splitCompositeKey(compositeKey string) (string, []string, error) {
	componentIndex := 1
	components := []string{}
	for i := 1; i < len(compositeKey); i++ {
		if compositeKey[i] == minUnicodeValue {
			components = append(components, compositeKey[componentIndex:i])
			componentIndex = i + 1
		}
	}
	return components[0], components[1:], nil
}

func validateAttribute(str string) error {
	if !utf8.ValidString(str) {
		return fmt.Errorf("not a valid utf8 string: [%x]", str)
	}
	for index, value := range str {
		if value == minUnicodeValue || value == maxUnicodeValue {
			return fmt.Errorf(`unicode %#U is error at position [%d]. %#U and %#U are not allowed \
                             in the input attribute of a composite key`,
				value, index, minUnicodeValue, maxUnicodeValue)
		}
	}
	return nil
}

// GetComKeyIterator get composite key iterator.
func (te *TxExecution) GetComKeyIterator(objectType string, attributes []string) (sdk.Iterator, error) {
	ctlog.Debugf("The objectType is %s, the attributes is %v", objectType, attributes)
	startKey, endKey, err := te.createRangeKeysForPartCompKey(objectType, attributes)
	if err != nil {
		return nil, err
	}
	ctlog.Debugf("In GetComKeyIterator, the startKey is %s and endKey is %s ", startKey, endKey)

	iterator, err := te.GetIterator(startKey, endKey)

	return iterator, err
}

// createRangeKeysForPartCompKey create range key for part compKey.
func (te *TxExecution) createRangeKeysForPartCompKey(objectType string, attributes []string) (string, string, error) {
	partialCompositeKey, err := te.GenerateComKey(objectType, attributes)
	if err != nil {
		return "", "", err
	}
	startKey := partialCompositeKey
	endKey := partialCompositeKey + string(utf8.MaxRune)
	return startKey, endKey, nil
}

// GenerateStateUpdates generate state updates.
func (te *TxExecution) GenerateStateUpdates() ([]*common.StateUpdates, error) {
	return nil, nil
}

// GenerateSQLStateUpdates generate SQL state updates.
func (te *TxExecution) GenerateSQLStateUpdates() ([]*common.StateUpdates, error) {
	return nil, nil
}

// GenerateExtension generate extension.
func (te *TxExecution) GenerateExtension() (map[string][]byte, error) {
	return nil, nil
}

// Start leveldb contract.
func Start(userContract sdk.Contract) {
	//var config *contract.ContractConfig
	config, err := readConfig()
	if err != nil {
		ctlog.Errorf("Start readConfig failed")
	}
	smcontext.GetSmContext().SetContractInfo(userContract)
	NewGrpcServer(config)
}

// StartSQL is start SQL contract
func StartSQL(userContract sdk.SQLContract) {
	// read config
	config, err := readConfig()
	if err != nil {
		ctlog.Errorf("StartSQL readConfig failed")
	}
	smcontext.GetSmSQLContext().SetContractInfo(userContract)
	NewGrpcServer(config)
}

// readConfig is read config
func readConfig() (*contract.ContractConfig, error) {
	// read config
	file, err := ioutil.ReadFile("/opt/contractconfig")
	if err != nil {
		ctlog.Errorf("Read config file failed,:%s", err.Error())
		return nil, err
	}
	config := &contract.ContractConfig{}
	err = proto.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	cfg := &logger.LogCfg{}
	cfg.Type = "pretty"
	cfg.Level = config.LogLevel
	err = logger.Init(cfg)
	if err != nil {
		ctlog.Errorf("Logger init failed")
		return nil, err
	}
	ctlog.Debugf("The config port is %s.", config.Port)
	return config, nil
}

// Exec is exec sql query
func (te *TxExecution) Exec(query string, args ...interface{}) error {
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = te.TransactionID
	contractMsg.Contractid = te.ContractID
	contractMsg.Chainid = te.ChainID()
	contractMsg.Type = contract.ContractMsg_SQL_EXEC
	sqlExecInfo := &contract.SqlExec{}
	sqlExecInfo.Query = query
	var err error
	sqlExecInfo.Args, err = utils.ConvertToProtoPrimitives(args...)
	if err != nil {
		ctlog.Errorf("ConvertToProtoPrimitive failed : %s\n", err)
		return err
	}
	payload, err := proto.Marshal(sqlExecInfo)
	if err != nil {
		return err
	}
	contractMsg.Payload = payload
	_, err = te.CallWnode(contractMsg)
	if err != nil {
		ctlog.Errorf("Call Wnode failed : %s\n", err)
		return err
	}
	return nil
}

// Query is sql query result .
func (te *TxExecution) Query(query string, args ...interface{}) (sdk.Rows, error) {
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = te.TransactionID
	contractMsg.Contractid = te.ContractID
	contractMsg.Chainid = te.ChainID()
	contractMsg.Type = contract.ContractMsg_SQL_QUERY
	sqlQueryInfo := &contract.SqlQuery{}
	sqlQueryInfo.Query = query
	var err error
	sqlQueryInfo.Args, err = utils.ConvertToProtoPrimitives(args...)
	if err != nil {
		ctlog.Errorf("ConvertToProtoPrimitive failed : %s\n", err)
	}
	sqlQueryInfo.BufferSize = bufferSize
	sqlQueryInfo.QId = GenerateQueryID()
	payload, err := proto.Marshal(sqlQueryInfo)
	if err != nil {
		return nil, err
	}
	contractMsg.Payload = payload
	payload, err = te.CallWnode(contractMsg)
	if err != nil {
		ctlog.Errorf("query failed : %s", err.Error())
		return nil, err
	}
	rows := &contract.SqlQueryResponse{}
	err = proto.Unmarshal(payload, rows)
	if err != nil {
		return nil, err
	}

	queryRows := &QueryRows{}
	queryRows.queryID = sqlQueryInfo.QId
	queryRows.query = query
	queryRows.args = sqlQueryInfo.Args
	queryRows.queryRes = rows
	queryRows.position = -1
	queryRows.txContext = te
	queryRows.onceLen = len(rows.RowsArray)
	ctlog.Debugf("The oncelen is %d", queryRows.onceLen)

	return queryRows, nil
}

// Next is sql query result next .
func (rows *QueryRows) Next() (bool, error) {
	var err error
	rows.position++
	if rows.position < rows.onceLen {
		return true, nil
	}
	if rows.queryRes.IsOver {
		return false, nil
	}
	// Query from wnode
	err = rows.queryRows(bufferSize)
	if err != nil {
		ctlog.Errorf("Query from wnode failed when get next")
		return false, err
	}
	if rows.onceLen == 0 {
		ctlog.Errorf("There is no data for next req")
		return false, nil
	}
	rows.position = 0
	return true, nil
}

func (rows *QueryRows) queryRows(bufferSize uint32) error {
	// The default buffer size is 100
	contractMsg := &contract.ContractMsg{}
	contractMsg.Txid = rows.txContext.TransactionID
	contractMsg.Contractid = rows.txContext.ContractID
	contractMsg.Chainid = rows.txContext.ChainID()
	contractMsg.Type = contract.ContractMsg_SQL_QUERY
	queryInfo := &contract.SqlQuery{}
	queryInfo.Query = rows.query
	queryInfo.Args = rows.args
	queryInfo.BufferSize = bufferSize
	queryInfo.QId = rows.queryID
	queryBytes, err := proto.Marshal(queryInfo)
	if err != nil {
		return err
	}
	contractMsg.Payload = queryBytes
	queryResBytes, err := rows.txContext.CallWnode(contractMsg)
	if err != nil {
		return err
	}
	queryRes := &contract.SqlQueryResponse{}
	err = proto.Unmarshal(queryResBytes, queryRes)
	if err != nil {
		return err
	}
	rows.position = 0
	rows.queryRes = queryRes
	rows.onceLen = len(queryRes.RowsArray)
	ctlog.Debugf("In queryFromWnode, the once len is %d", rows.onceLen)
	return nil
}

// Scan is query scan rows .
func (rows *QueryRows) Scan(dest ...interface{}) error {
	if len(rows.queryRes.RowsArray) == 0 || rows.position < 0 {
		return fmt.Errorf("next function is should be called")
	}
	values, err := utils.ConvertToGolangPrimitives(rows.queryRes.RowsArray[rows.position])
	if err != nil {
		ctlog.Errorf("ConvertToGolangPrimitives failed when scan: %s", err)
		return err
	}
	if len(dest) != len(values) {
		return fmt.Errorf("expected variables is %d, got %d", len(values), len(dest))
	}

	for i, d := range dest {
		err = convertQueryRows(d, values[i])
		if err != nil {
			ctlog.Errorf("convertQueryRows failed when scan: %s", err)
			return err
		}
	}
	return nil
}

// Close is sql query close .
func (rows *QueryRows) Close() error {
	err := rows.queryRows(0)
	if err != nil {
		ctlog.Errorf("Query from wnode failed when close the rows")
		return err
	}
	return nil
}

// convertQueryRows is convert query rows into interface{} .
func convertQueryRows(dest, src interface{}) error { /// nolint
	// Common cases, without reflect.
	ctlog.Debugf("convert query rows: %v, %v", dest, src)
	err, isContinue := convertTypePhase(dest, src)
	if err != nil {
		return err
	}
	if !isContinue {
		return nil
	}

	var sv reflect.Value
	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	default:
		ctlog.Debugf("dest pointer is : %v", d)
	}

	dpv := reflect.ValueOf(dest)
	dv, err := validateData(sv, dpv, src)
	if err != nil {
		return err
	}
	// The following conversions use a string value as an intermediate representation
	// to convert between various numeric types.
	ctlog.Debugf("convert query rows dv: %v, src: %v", dv.Kind(), src)
	return convertNumericType(dv, src)
}

func validateData(sv reflect.Value, dpv reflect.Value, src interface{}) (reflect.Value, error) {
	if dpv.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("dpv.kind() != reflect.ptr, dpv.kind(): %v", dpv.Kind())
	}
	if dpv.IsNil() {
		return reflect.Value{}, fmt.Errorf("dpv is nil(): %v", dpv.IsNil())
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}
		return dv, nil
	}
	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
	}
	return dv, nil
}

func convertNumericType(dv reflect.Value, src interface{}) error { /// nolint
	switch dv.Kind() {
	case reflect.Ptr:
		if src == nil {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		}
		dv.Set(reflect.New(dv.Type().Elem()))
		return convertQueryRows(dv.Interface(), src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if src == nil {
			return fmt.Errorf("converting dv.Kind(): %v", dv.Kind())
		}
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	case reflect.String:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		switch v := src.(type) {
		case string:
			dv.SetString(v)
			return nil
		case []byte:
			dv.SetString(string(v))
			return nil
		default:
			return nil
		}
	}
	return fmt.Errorf("unsupported dv.Kind(), storing driver.Value type %T", src)
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func convertTypePhase(dest, src interface{}) (error, bool) {
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return fmt.Errorf("destination pointer is nil, source: %s, dest: %v", s, d), false
			}
			*d = s
			return nil, false
		case *[]byte:
			if d == nil {
				return fmt.Errorf("destination pointer is nil, source: %s, dest: %v", s, d), false
			}
			*d = []byte(s)
			return nil, false
		default:
			ctlog.Debugf("source pointer is string, dest: %v, source: %v", d, s)
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return fmt.Errorf("destination pointer is nil, source: %s, dest: %v", s, d), false
			}
			*d = string(s)
			return nil, false
		case *interface{}:
			if d == nil {
				return fmt.Errorf("destination pointer is nil, source: %s, dest: %v", s, d), false
			}
			*d = cloneBytes(s)
			return nil, false
		case *[]byte:
			if d == nil {
				return fmt.Errorf("destination pointer is nil, source: %s, dest: %v", s, d), false
			}
			*d = cloneBytes(s)
			return nil, false
		default:
			ctlog.Debugf("source pointer is []byte, dest: %v, source: %v", d, s)
		}
	case nil:
		switch d := dest.(type) {
		case *interface{}:
			if d == nil {
				return fmt.Errorf("destination pointer is nil, source: %s, dest: %v", s, d), false
			}
			*d = nil
			return nil, false
		case *[]byte:
			if d == nil {
				return fmt.Errorf("destination pointer is nil, source: %s, dest: %v", s, d), false
			}
			*d = nil
			return nil, false
		default:
			ctlog.Debugf("source pointer is nil, dest: %v, source: %v", d, s)
		}
	default:
		ctlog.Debugf("source pointer is : %v", s)
	}
	return nil, true
}
