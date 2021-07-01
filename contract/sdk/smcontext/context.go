/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

// Package smcontext the smart contract context.
package smcontext

import (
	"fmt"
	"sync"

	"git.huawei.com/poissonsearch/wienerchain/contract/sdk"
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/logger"
)

// SmCtx is smart contract context.
type SmCtx struct {
	smartContract sdk.Contract
	stubMap       sync.Map
}

// SmCtxSQLInfo is SQL contract context.
type SmCtxSQLInfo struct {
	smSQLContract sdk.SQLContract
	stubMap       sync.Map
}

var log = logger.GetDefaultLogger()

// SmContext is leveldb context
var smContext *SmCtx

// SmSQLContext is sql context
var smSQLContext *SmCtxSQLInfo

var lockContext sync.Mutex

// GetSmContext generate and return single smcontext.
func GetSmContext() sdk.Context {
	lockContext.Lock()
	defer lockContext.Unlock()
	if smContext != nil {
		return smContext
	}
	smContext = &SmCtx{}
	return smContext
}

// SetContractInfo set the user contract in the context.
func (sc *SmCtx) SetContractInfo(contract sdk.Contract) {
	sc.smartContract = contract
}

// GetContractInfo get the user contract.
func (sc *SmCtx) GetContractInfo() sdk.Contract {
	return sc.smartContract
}

// DelStub delete the stub for certain tx id.
func (sc *SmCtx) DelStub(stub sdk.ContractStub) {
	sc.stubMap.Delete(stub.TxID())
}

// AddStub add stub for certain tx id.
func (sc *SmCtx) AddStub(stub sdk.ContractStub) error {
	_, ok := sc.stubMap.Load(stub.TxID())
	if !ok {
		sc.stubMap.Store(stub.TxID(), stub)
		return nil
	}
	log.Infof("The smart contract stub for txid: %s already exist", stub.TxID())
	return fmt.Errorf("the stub for txid: %s already exist", stub.TxID())
}

// GetStub get stub for certain tx id
func (sc *SmCtx) GetStub(txID string) sdk.ContractStub {
	stub, ok := sc.stubMap.Load(txID)
	if !ok {
		return nil
	}
	trueStub, ok := stub.(sdk.ContractStub)
	if !ok {
		return nil
	}
	return trueStub
}

// GetSmSQLContext Generate and return single SmSQLContext
func GetSmSQLContext() sdk.SQLContext {
	lockContext.Lock()
	defer lockContext.Unlock()
	if smSQLContext != nil {
		return smSQLContext
	}
	smSQLContext = &SmCtxSQLInfo{}
	return smSQLContext
}

// SetContractInfo is Set the user contract in the context
func (sc *SmCtxSQLInfo) SetContractInfo(contract sdk.SQLContract) {
	sc.smSQLContract = contract
}

// GetContractInfo Get the user contract
func (sc *SmCtxSQLInfo) GetContractInfo() sdk.SQLContract {
	return sc.smSQLContract
}

// DelStub is the stub for certain tx id
func (sc *SmCtxSQLInfo) DelStub(stub sdk.SQLContractStub) {
	sc.stubMap.Delete(stub.TxID())
}

// AddStub is Add stub for certain tx id
func (sc *SmCtxSQLInfo) AddStub(stub sdk.SQLContractStub) error {
	_, ok := sc.stubMap.Load(stub.TxID())
	if !ok {
		sc.stubMap.Store(stub.TxID(), stub)
	} else {
		log.Infof("The stub for txid: %s already exist", stub.TxID())
		err := fmt.Errorf("the stub for txid: %s already exist", stub.TxID())
		return err
	}

	return nil
}

// GetStub is Get stub for certain tx ID
func (sc *SmCtxSQLInfo) GetStub(txID string) sdk.SQLContractStub {
	stub, ok := sc.stubMap.Load(txID)
	if !ok {
		return nil
	}
	trueStub, ok := stub.(sdk.SQLContractStub)
	if !ok {
		return nil
	}
	return trueStub
}
