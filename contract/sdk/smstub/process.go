/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

package smstub

import (
	"fmt"
	"runtime/debug"

	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/logger"
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/smcontext"
	"git.huawei.com/poissonsearch/wienerchain/proto/common"
	"git.huawei.com/poissonsearch/wienerchain/proto/contract"
	"github.com/golang/protobuf/proto"
)

var plog = logger.GetDefaultLogger()

const (
	recvChanBuf = 10
)

// ProcessMsgFromWnode process msg from wnode
func ProcessMsgFromWnode(conMsg *contract.ContractMsg, streamSendChan chan contract.ContractMsg) {
	switch conMsg.Type {
	case contract.ContractMsg_INIT, contract.ContractMsg_INVOKE,
		contract.ContractMsg_SQL_INIT, contract.ContractMsg_SQL_INVOKE:
		plog.Debugf("Contract msg type: %d", conMsg.Type)
		err := ProcessWnodeReq(conMsg, streamSendChan)
		if err != nil {
			plog.Errorf("Contract msg type: %d, process wnode request failed:%s", conMsg.Type, err.Error())
		}
	case contract.ContractMsg_PUT_KV, contract.ContractMsg_GET_KV,
		contract.ContractMsg_DEL_KV, contract.ContractMsg_GET_KV_BY_RANGE,
		contract.ContractMsg_SQL_EXEC, contract.ContractMsg_SQL_QUERY,
		contract.ContractMsg_GET_HISTORY_FOR_KEY:
		err := ProcessWnodeDBRes(conMsg)
		if err != nil {
			plog.Errorf("When process:%s, process wnode db response failed:%s", conMsg.Type, err.Error())
		}
	default:
		plog.Debugf("Receive unknown type message.")
	}
}

// ProcessWnodeDBRes process msg from wnode DB response.
func ProcessWnodeDBRes(conMsg *contract.ContractMsg) error {
	switch conMsg.Type {
	case contract.ContractMsg_PUT_KV, contract.ContractMsg_DEL_KV:
		stub := smcontext.GetSmContext().GetStub(conMsg.Txid)
		if stub == nil {
			return fmt.Errorf("the db response is not correct")
		}
		placeHolder := make([]byte, 0)
		te, ok := stub.(*TxExecution)
		if !ok {
			return fmt.Errorf("the stub is not TxExecution")
		}
		te.RecvChan <- placeHolder
	case contract.ContractMsg_GET_KV, contract.ContractMsg_GET_KV_BY_RANGE,
		contract.ContractMsg_GET_HISTORY_FOR_KEY:
		stub := smcontext.GetSmContext().GetStub(conMsg.Txid)
		if stub == nil {
			return fmt.Errorf("the stub response is not correct")
		}
		te, ok := stub.(*TxExecution)
		if !ok {
			return fmt.Errorf("the stub is not TxExecution")
		}
		te.RecvChan <- conMsg.Payload
	case contract.ContractMsg_SQL_EXEC:
		sqlStub := smcontext.GetSmSQLContext().GetStub(conMsg.Txid)
		if sqlStub == nil {
			return fmt.Errorf("the sql stub response is not correct")
		}
		placeHolder := make([]byte, 0)
		te, ok := sqlStub.(*TxExecution)
		if !ok {
			return fmt.Errorf("the stub is not TxExecution")
		}
		te.RecvChan <- placeHolder
	case contract.ContractMsg_SQL_QUERY:
		sqlStub := smcontext.GetSmSQLContext().GetStub(conMsg.Txid)
		if sqlStub == nil {
			return fmt.Errorf("the sql stub response is not correct")
		}
		te, ok := sqlStub.(*TxExecution)
		if !ok {
			return fmt.Errorf("the sql stub is not TxExecution")
		}
		te.RecvChan <- conMsg.Payload
	default:
		{
			return fmt.Errorf("unknown msg type")
		}
	}
	return nil
}

// ProcessWnodeReq process request from wnode.
func ProcessWnodeReq(conMsg *contract.ContractMsg, streamSendChan chan contract.ContractMsg) error {
	plog.Debugf("ProcessWnodeReq msg type :%s", conMsg.Type)
	defer func() {
		err := recover()
		if err != nil {
			plog.Infof("ProcessWnodeReq failed, recover, the err is :%v. \n the stacktrace is %s\n", err, string(debug.Stack()))
		}
	}()

	stubOb := &TxExecution{}
	stubOb.ChainName = conMsg.Chainid
	stubOb.ContractID = conMsg.Contractid
	stubOb.TransactionID = conMsg.Txid
	stubOb.StreamSendChan = streamSendChan
	contractInvocation := &common.ContractInvocation{}
	if err := proto.Unmarshal(conMsg.Payload, contractInvocation); err != nil {
		return fmt.Errorf("unmarshal ContractInvocation failed")
	}
	stubOb.Invocation = contractInvocation
	//In one tx invocation, maybe there is multi thread to call wnode DB
	stubOb.RecvChan = make(chan []byte, recvChanBuf)

	if conMsg.Type == contract.ContractMsg_INIT ||
		conMsg.Type == contract.ContractMsg_INVOKE {
		if err := smcontext.GetSmContext().AddStub(stubOb); err != nil {
			resMsg := constructErrRes(conMsg.Txid, conMsg.Contractid, conMsg.Chainid, "The tx id already exist.")
			streamSendChan <- *resMsg
			return err
		}
	} else if conMsg.Type == contract.ContractMsg_SQL_INIT ||
		conMsg.Type == contract.ContractMsg_SQL_INVOKE {
		if err := smcontext.GetSmSQLContext().AddStub(stubOb); err != nil {
			resMsg := constructErrRes(conMsg.Txid, conMsg.Contractid, conMsg.Chainid, "The tx id already exist.")
			streamSendChan <- *resMsg
			return err
		}
	}

	plog.Debugf("ProcessWnodeReq msg type :%s\n", conMsg.Type)
	var response common.InvocationResponse
	switch conMsg.Type {
	case contract.ContractMsg_INVOKE:
		response = smcontext.GetSmContext().GetContractInfo().Invoke(stubOb)
	case contract.ContractMsg_INIT:
		response = smcontext.GetSmContext().GetContractInfo().Init(stubOb)
	case contract.ContractMsg_SQL_INIT:
		response = smcontext.GetSmSQLContext().GetContractInfo().InitSQL(stubOb)
	case contract.ContractMsg_SQL_INVOKE:
		response = smcontext.GetSmSQLContext().GetContractInfo().InvokeSQL(stubOb)
	}
	resMsg := constructNormalRes(conMsg, response)
	stubOb.StreamSendChan <- *resMsg
	smcontext.GetSmContext().DelStub(stubOb)
	return nil
}

func constructNormalRes(conMsg *contract.ContractMsg, invokeRes common.InvocationResponse) *contract.ContractMsg {
	resMsg := &contract.ContractMsg{}
	resMsg.Type = conMsg.Type
	resMsg.Txid = conMsg.Txid
	resMsg.Contractid = conMsg.Contractid
	resMsg.Chainid = conMsg.Chainid
	resPayload, err := proto.Marshal(&invokeRes)
	if err != nil {
		return nil
	}
	resMsg.Payload = resPayload
	return resMsg
}

func constructErrRes(txid string, contractID string, chainID string, errInfo string) *contract.ContractMsg {
	resMsg := &contract.ContractMsg{}
	resMsg.Type = contract.ContractMsg_INVOKE
	resMsg.Txid = txid
	resMsg.Contractid = contractID
	resMsg.Chainid = chainID
	invokeRes := &common.InvocationResponse{}
	invokeRes.Status = common.Status_BAD_REQUEST

	invokeRes.StatusInfo = errInfo
	payload, _ := proto.Marshal(invokeRes)
	resMsg.Payload = payload
	return resMsg
}
