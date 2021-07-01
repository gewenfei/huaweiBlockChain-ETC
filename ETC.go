package main

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"strconv"
	// "strings"

	"git.huawei.com/poissonsearch/wienerchain/contract/sdk"
	// "git.huawei.com/poissonsearch/wienerchain/contract/sdk/logger"
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/smstub"
	"git.huawei.com/poissonsearch/wienerchain/proto/common"
)

type ETC struct {
}

// var log = logger.GetDefaultLogger()

func (e ETC) Init(stub sdk.ContractStub) common.InvocationResponse {
	fmt.Printf("Enter ETC init function\n")
	args := stub.Parameters()
	const numOfArgs = 2
	if len(args) < numOfArgs {
		return sdk.Error("Init parameter is not correct")
	}
	// 用户地址
	address := string(args[0])
	// 初始化余额
	value := args[1]
	err := stub.PutKV(address, value)
	if err != nil {
		return sdk.Error("Init put kv failed")
	}

	return sdk.Success([]byte("初始化成功！"))
}

func (e ETC) Invoke(stub sdk.ContractStub) common.InvocationResponse {
	funcName := stub.FuncName()
	args := stub.Parameters()

	switch funcName {
	// 初始化某用户ETC
	case "initFunds":
		return initFunds(stub, args)
	// 质押资金
	case "depositFunds":
		return depositFunds(stub, args)
	// 抵扣资金
	case "deductFunds":
		return deductFunds(stub, args)
	// 转移资金
	case "transferFunds":
		return transferFunds(stub, args)
	// 查询资金余额
	case "getFunds":
		return getFunds(stub, args)
	}
	str := fmt.Sprintf("Func name is not correct, the function name is %s ", funcName)
	return sdk.Error(str)
}

// 质押资金
func depositFunds(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
	const numOfArgs = 2
	if len(args) < numOfArgs {
		return sdk.Error("the number of args is not correct")
	}

	address := string(args[0])
	value,err := strconv.Atoi(string(args[1]))
	if err != nil {
		fmt.Printf("The ETC size is not int type\n")
		return sdk.Error("The ETC size is not int type")
	}

	oldValue, err := stubInterface.GetKV(address)
	if err != nil {
		fmt.Printf("getkv error, the err is :%s\n", err.Error())
		return sdk.Error("get kv failed")
	}
	intOldValue,err := strconv.Atoi(string(oldValue))
	if err != nil {
		fmt.Printf("The ETC size is not int type\n")
		return sdk.Error("The ETC size is not int type")
	}
	err = stubInterface.PutKV(address, []byte(strconv.Itoa(intOldValue + value)))
	if err != nil {
		return sdk.Error(err.Error())
	}
	info := fmt.Sprintf("质押成功 地址: %s 质押金额：%d 余额：%d", address, value, intOldValue + value)
	return sdk.Success([]byte(info))
}

// 抵扣资金
func deductFunds(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
	const numOfArgs = 2
	if len(args) < numOfArgs {
		return sdk.Error("the number of args is not correct")
	}

	address := string(args[0])
	value,err := strconv.Atoi(string(args[1]))
	if err != nil {
		fmt.Printf("The ETC size is not int type\n")
		return sdk.Error("The ETC size is not int type")
	}

	oldValue, err := stubInterface.GetKV(address)
	if err != nil {
		fmt.Printf("getkv error, the err is :%s\n", err.Error())
		return sdk.Error("get kv failed")
	}
	intOldValue,_ := strconv.Atoi(string(oldValue))
	if(intOldValue < value){
		return sdk.Error("余额不足，抵扣失败！")
	}
	err = stubInterface.PutKV(address, []byte(strconv.Itoa(intOldValue - value)))
	if err != nil {
		return sdk.Error(err.Error())
	}
	info := fmt.Sprintf("抵扣成功 地址: %s 抵扣金额：%d 余额：%d", address, value, intOldValue - value)
	return sdk.Success([]byte(info))
}

// 转移资金
func transferFunds(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
	const numOfArgs = 3
	if len(args) < numOfArgs {
		return sdk.Error("the number of args is not correct")
	}

	addressFrom := string(args[0])
	addressTo := string(args[1])

	// 地址存在检查
	// errCheck := checkAccountExist(stubInterface, addressFrom)
	// if errCheck.Status != 0 {
	// 	return errCheck
	// }
	// errCheck = checkAccountExist(stubInterface, addressTo)
	// if errCheck.Status != 0 {
	// 	return errCheck
	// }

	value,err := strconv.Atoi(string(args[2]))
	if err != nil {
		fmt.Printf("The ETC size is not int type\n")
		return sdk.Error("The ETC size is not int type")
	}

	oldValueFrom, err := stubInterface.GetKV(addressFrom)
	if err != nil {
		fmt.Printf("getkv error, the err is :%s\n", err.Error())
		return sdk.Error("get kv failed")
	}
	oldValueTo, err := stubInterface.GetKV(addressTo)
	if err != nil {
		fmt.Printf("getkv error, the err is :%s\n", err.Error())
		return sdk.Error("get kv failed")
	}
	intOldValueFrom,_ := strconv.Atoi(string(oldValueFrom))
	intOldValueTo,_ := strconv.Atoi(string(oldValueTo))
	if(intOldValueFrom < value){
		info := fmt.Sprintf("余额不足，转账失败！ 地址：%s，余额：%d，需要额度：%d ", addressFrom, intOldValueFrom, value)
		return sdk.Error(info)
	}
	err = stubInterface.PutKV(addressFrom, []byte(strconv.Itoa(intOldValueFrom - value)))
	if err != nil {
		return sdk.Error(err.Error())
	}
	err = stubInterface.PutKV(addressTo, []byte(strconv.Itoa(intOldValueTo + value)))
	if err != nil {
		return sdk.Error(err.Error())
	}

	info := fmt.Sprintf("转账成功 From地址: %s 抵扣金额：%d 余额：%d To地址：%s 余额：%d", addressFrom, value, intOldValueFrom - value, addressTo, intOldValueTo + value)
	return sdk.Success([]byte(info))
}

// 查询资金余额
func getFunds(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
	if len(args) < 1 {
		fmt.Println("The args number is not correct")
		return sdk.Error("The args is not correct.")
	}

	address := string(args[0])
	value, err := stubInterface.GetKV(address)
	if err != nil {
		errInfo := fmt.Sprintf("Get the key: %s failed", address)
		return sdk.Error(errInfo)
	}

	info := fmt.Sprintf("查询成功 地址: %s 余额：%s ", address, value)
	return sdk.Success([]byte(info))
}

// 初始化ETC
func initFunds(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
	const numOfArgs = 2
	if len(args) < numOfArgs {
		fmt.Println("The args number is not correct")
		return sdk.Error("The args is not correct.")
	}
	address := string(args[0])

	value, err := stubInterface.GetKV(address)
	if err != nil {
		errInfo := fmt.Sprintf("Get the ETC info err:%s", err.Error())
		return sdk.Error(errInfo)
	}
	if value != nil {
		sdk.Error("The key to be add is already exist")
	}
	value = args[1]

	err = stubInterface.PutKV(address, value)
	if err != nil {
		return sdk.Error(err.Error())
	}
	info := fmt.Sprintf("初始化成功 地址: %s 余额：%s ", address, value)
	return sdk.Success([]byte(info))
}

func checkAccountExist(stubInterface sdk.ContractStub, address string) common.InvocationResponse{
	value, err := stubInterface.GetKV(address)
	if err != nil {
		errInfo := fmt.Sprintf("Get the TEC info err:%s", err.Error())
		return sdk.Error(errInfo)
	}
	if value != nil {
		errInfo := fmt.Sprintf("地址：%s 不存在！", address)
		return sdk.Error(errInfo)
	}
	return sdk.Success(nil)
}

func main() {
	smstub.Start(&ETC{})
}
