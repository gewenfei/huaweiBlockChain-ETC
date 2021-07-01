package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"git.huawei.com/poissonsearch/wienerchain/contract/sdk"
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/logger"
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/smstub"
	"git.huawei.com/poissonsearch/wienerchain/proto/common"
)

type etcManager struct {
}

/**
 * @Description:数字证书结构体
 */
type DigitalCertificate struct {
	ID int64					// 证书序号
	PublicKey []byte			// 所属人公钥
	BasicInfo []byte			// 所属人基本信息Hash
	ChainID  string				// 证书所在链
	Expires int64				// 证书截止时间
	State int8					// 当前证书签名状态 0:不存在 1:未开始 2:交管局系统已签名 3:银行财务系统已签名 4:省中心已签名
	Burn bool					// 注销
	SettlementSig []byte		// 省结算中心签名
	BankSig		  []byte		// 银行财务系统签名
	TrafficManagementSig []byte	// 交管局系统签名
}

type keyHistory struct {
	Value string
	// TxHash []byte
	BlockNum uint64
	TxNum int32
	IsDeleted bool
	// Timestamp int64
}

var log = logger.GetDefaultLogger()

/**
 * @Description:初始化函数
 * @receiver e
 * @param stub	初始化参数: 省结算中心、银行财务系统、ETC收费系统认证密钥(一般与签名私钥不同)
 * @return common.InvocationResponse
 */
func (e etcManager) Init(stub sdk.ContractStub) common.InvocationResponse {
	fmt.Printf("请分别输入:省结算中心、银行财务系统、ETC收费系统的认证密钥\n")
	args := stub.Parameters()
	const numOfArgs = 3
	if len(args) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	// 初始化密钥
	settlementKey, bankKey, trafficManagementKey := args[0], args[1], args[2]
	err := stub.PutKV("settlementKey", settlementKey)
	if err != nil {
		return sdk.Error("存储settlementKey失败")
	}
	err = stub.PutKV("bankKey", bankKey)
	if err != nil {
		return sdk.Error("存储bankKey失败")
	}
	err = stub.PutKV("trafficManagementKey", trafficManagementKey)
	if err != nil {
		return sdk.Error("存储chargeSystemKey失败")
	}
	// 初始化证书编号
	err = stub.PutKV("totalCertificate", []byte("0"))
	if err != nil {
		return sdk.Error("存储totalCertificate失败")
	}
	return sdk.Success(nil)
}

/**
 * @Description:验证权限
 * @param stub
 * @param kind 权限机关key类别
 * @param key  key值
 * @return bool	成功/失败
 */
func VerifyAuth(stub sdk.ContractStub, kind []byte, key []byte) bool {
	// 获取对应认证密钥
	val, err := stub.GetKV(string(kind))
	if err != nil {
		log.Error("获取" + string(kind) + "失败")
		return false
	}
	return bytes.Compare(val, key) == 0
}

/**
 * @Description:请求路由函数
 * @receiver e
 * @param stub
 * @return common.InvocationResponse
 */
func (e etcManager) Invoke(stub sdk.ContractStub) common.InvocationResponse {
	funcName := stub.FuncName()
	args := stub.Parameters()
	if funcName == "getMyCertificate" {
		return getMyCertificate(stub, args)
	}else {
		// 验证权限
		kind, key := args[0], args[1]
		if !VerifyAuth(stub, kind, key) {
			str := fmt.Sprintf("请求函数权限不足, 请求的函数是:%s ", funcName)
			return sdk.Error(str)
		}
		// 路由
		switch funcName {
		case "storeTotalCertificateNum":
			// 存储证书总数
			return storeTotalCertificateNum(stub, kind, args[2:])
		case "getNowCertificateNum":
			// 获取证书当前总数
			return getNowCertificateNum(stub)
		case "storeBasicInfoHash":
			// 存储车主基本信息Hash
			return storeBasicInfoHash(stub, kind, args[2:])
		case "getBasicInfoHash":
			// 获取车主基本信息Hash
			return getBasicInfoHash(stub, args[2:])
		case "generateCertificate":
			// 生成初始化的证书
			return generateCertificate(stub, kind, args[2:])
		case "getInitCertificate":
			// 获取初始化的证书（未签名）
			return getInitCertificate(stub, args[2:])
		case "attachSigning":
			// 各部门附加签名
			return attachSigning(stub, kind, args[2:])
		case "getRangeCertificates":
			// 获取所有的已认证证书
			return getRangeCertificates(stub, args[2:])
			// 初始化某用户ETC
		case "initFunds":
			return initFunds(stub, kind, args[2:])
		// 质押资金
		case "depositFunds":
			return depositFunds(stub, kind, args[2:])
		// 抵扣资金
		case "deductFunds":
			return deductFunds(stub, kind, args[2:])
		// 转移资金
		case "transferFunds":
			return transferFunds(stub, kind, args[2:])
		// 查询资金余额
		case "getFunds":
			return getFunds(stub, args[2:])
		// 查询资金余额
		case "getFundsHistory":
			return getFundsHistory(stub, args[2:])
		}
	}
	str := fmt.Sprintf("请求的函数不支持, 请求的函数是:%s ", funcName)
	return sdk.Error(str)

}

/**
 * @Description:存储证书总量
 * @param stub
 * @param kind
 * @param arg
 * @return common.InvocationResponse
 */
func storeTotalCertificateNum(stub sdk.ContractStub, kind []byte, arg [][]byte) common.InvocationResponse {
	if string(kind) != "settlementKey" {
		return sdk.Error("权限不足")
	}
	const numOfArgs = 1
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	c_num := arg[0]
	err := stub.PutKV("totalCertificate", c_num)
	if err != nil {
		return sdk.Error("存储totalCertificate失败")
	}
	return sdk.Success(nil)
}

/**
 * @Description:获取当前证书总量
 * @param stub
 * @return common.InvocationResponse
 */
func getNowCertificateNum(stub sdk.ContractStub) common.InvocationResponse  {
	numBytes, err := stub.GetKV("totalCertificate")
	if err != nil {
		return sdk.Error("读取totalCertificate失败")
	}
	if numBytes == nil {
		return sdk.Error("证书数量未初始化")
	}
	return sdk.Success(numBytes)
}

/**
 * @Description:存储车主基本信息Hash
 * @param stub
 * @param kind
 * @param arg
 * @return common.InvocationResponse
 */
func storeBasicInfoHash(stub sdk.ContractStub, kind []byte, arg [][]byte) common.InvocationResponse {
	if string(kind) != "settlementKey" {
		return sdk.Error("权限不足")
	}
	const numOfArgs = 2
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	publicKey, infoHash := arg[0], arg[1]
	err := stub.PutKV(string(publicKey), infoHash)
	if err != nil {
		return sdk.Error("存储infoHash失败")
	}
	return sdk.Success(nil)
}

/**
 * @Description:查询当前车主的InfoHash
 * @param stub
 * @param kind
 * @param arg
 * @return common.InvocationResponse
 */
func getBasicInfoHash(stub sdk.ContractStub, arg [][]byte) common.InvocationResponse  {
	const numOfArgs = 1
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	publicKey := arg[0]
	infoHash, err := stub.GetKV(string(publicKey))
	if err != nil {
		return sdk.Error("读取infoHash失败")
	}
	return sdk.Success(infoHash)
}


/**
 * @Description:生成证书
 * @param stub
 * @return common.InvocationResponse
 */
func generateCertificate(stub sdk.ContractStub, kind []byte, arg [][]byte) common.InvocationResponse  {
	// 1. 判断权限
	if string(kind) != "settlementKey" {
		return sdk.Error("权限不足")
	}
	// 2. 判断参数个数
	const numOfArgs = 2
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	// 3. 解析参数
	publicKey, expiresBytes := arg[0], arg[1]
	expires, _ := strconv.Atoi(string(expiresBytes))
	// 4. 获取基本信息Hash
	infoHash, err := stub.GetKV(string(publicKey))
	if err != nil {
		return sdk.Error("generateCertificate GetKV error")
	}
	if infoHash == nil {
		return sdk.Error("请先上传基本信息Hash")
	}
	// 5. 生成证书
	idBytes, _ := stub.GetKV("totalCertificate")
	id, _ := strconv.Atoi(string(idBytes))
	newCertificate := &DigitalCertificate{
		ID:                 int64(id),
		PublicKey:          publicKey,
		BasicInfo:          infoHash,
		ChainID:			  stub.ChainID(),
		Expires:            int64(expires),
		State:              1,
		Burn:               false,
	}
	// 6. 总量+1
	err = stub.PutKV("totalCertificate", []byte(strconv.Itoa(id+1)))
	if err != nil {
		return sdk.Error("读取totalCertificate错误")
	}
	// 7. 持久化存储
	err = stub.PutKVCommon(strconv.Itoa(id), newCertificate)
	if err != nil {
		return sdk.Error("GenerateCertificate PutKVCommon error")
	}
	return sdk.Success(idBytes)
}

/**
 * @Description:获取初始化证书
 * @param stub
 * @param kind
 * @param arg 1:证书id
 * @return common.InvocationResponse
 */
func getInitCertificate(stub sdk.ContractStub, arg [][]byte) common.InvocationResponse {
	const numOfArgs = 1
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	idBytes := arg[0]
	certificate, err := stub.GetKV(string(idBytes))
	if err != nil {
		return sdk.Error("getInitCertificate GetKV error")
	}
	if certificate == nil {
		return sdk.Error("不存在此证书")
	}
	// 解码
	dc, err := Unmarshal(certificate)
	if err != nil {
		return sdk.Error("getInitCertificate Unmarshal error")
	}
	s := fmt.Sprintf("DigitalCertificate的序号是%d,所属人公钥是%s,所属人基本信息Hash是%s,证书截止时间是%d,当前证书状态是%d", dc.ID, dc.PublicKey, dc.BasicInfo, dc.Expires, dc.State)
	log.Info(s)
	dcJSON, err := dc.Marshal()
	if err != nil {
		return sdk.Error("getMyCertificate Marshal error")
	}
	return sdk.Success(dcJSON)
}

/**
 * @Description:证书附加签名
 * @param stub
 * @param kind
 * @param arg  1:证书ID 2:签名
 * @return common.InvocationResponse
 */
func attachSigning(stub sdk.ContractStub, kind []byte, arg [][]byte) common.InvocationResponse  {
	const numOfArgs = 2
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	idBytes, signature := arg[0], arg[1]
	certificate, err := stub.GetKV(string(idBytes))
	if err != nil {
		return sdk.Error("settlementSigning GetKV error")
	}
	if certificate == nil {
		return sdk.Error("不存在此证书")
	}
	dc, err := Unmarshal(certificate)
	if err != nil {
		return sdk.Error("settlementSigning Unmarshal error")
	}
	// 设置签名
	switch string(kind) {
	case "settlementKey":
		if dc.State < 3 {
			return sdk.Error("请先在交管局、银行部门核实签名")
		}
		dc.SettlementSig = signature
		dc.State = 4
	case "bankKey":
		if dc.State < 2 {
			return sdk.Error("请先在交管局核实签名")
		}
		dc.BankSig = signature
		dc.State = 3
	case "trafficManagementKey":
		if dc.State < 1 {
			return sdk.Error("请先创建证书")
		}
		dc.TrafficManagementSig = signature
		dc.State = 2
	}
	// 存储
	err = stub.PutKVCommon(string(idBytes), dc)
	if err != nil {
		return sdk.Error("settlementSigning PutKVCommon error")
	}
	return sdk.Success(nil)
}

/**
 * @Description:用户调用，获取自己的证书
 * @param stub
 * @param arg 1:id 0:不存在 1:未开始 2:交管局系统已签名 3:银行财务系统已签名 4:省中心已签名
 * @return common.InvocationResponse
 */
func getMyCertificate(stub sdk.ContractStub, arg [][]byte) common.InvocationResponse {
	const numOfArgs = 1
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	idBytes := arg[0]
	certificate, err := stub.GetKV(string(idBytes))
	if err != nil {
		return sdk.Error("settlementSigning GetKV error")
	}
	if certificate == nil {
		return sdk.Error("查找的证书不存在!")
	}
	dc, err := Unmarshal(certificate)
	if err != nil {
		return sdk.Error("settlementSigning Unmarshal error")
	}
	switch dc.State {
	case 1:
		return sdk.Error("您的证书正在签署")
	case 2:
		return sdk.Error("交管局已签名核实，请耐心等待")
	case 3:
		return sdk.Error("银行财务系统已签名，请耐心等待")
	case 4:
		// 返回证书
		s := fmt.Sprintf("DigitalCertificate的序号是%d,所属人公钥是%s,所属人基本信息Hash是%s,证书截止时间是%d,当前证书状态是%d,省中心签名是%s,交管局签名是%s,银行签名是%s", dc.ID, dc.PublicKey, dc.BasicInfo, dc.Expires, dc.State, dc.SettlementSig, dc.TrafficManagementSig ,dc.BankSig)
		log.Info(s)
		dcJSON, err := dc.Marshal()
		if err != nil {
			return sdk.Error("getMyCertificate Marshal error")
		}
		return sdk.Success(dcJSON)
	}
	return sdk.Success(nil)
}

/**
 * @Description:获取所有已完全签署的证书
 * @param stub
 * @param arg	starKey, endKey
 * @return common.InvocationResponse
 */
func getRangeCertificates(stub sdk.ContractStub, arg [][]byte) common.InvocationResponse {
	const numOfArgs = 2
	if len(arg) < numOfArgs {
		return sdk.Error("输入参数错误")
	}
	starKey, endKey := arg[0], arg[1]
	iterator, err := stub.GetIterator(string(starKey), string(endKey))
	if err != nil {
		return sdk.Error("getAllCertificates GetIterator error")
	}
	defer iterator.Close()
	rangeMap := make(map[string]string)
	var count = 0
	for {
		b := iterator.Next()
		if b {
			key := iterator.Key()
			value := string(iterator.Value())
			rangeMap[key] = value
			count++
			log.Debugf("The iterator read key is %s, value is %s, count is %d\n", key, value, count)
		} else {
			log.Debugf("The iterator break\n")
			break
		}
	}
	rangeMapBytes, err := json.Marshal(rangeMap)
	if err != nil {
		return sdk.Error(err.Error())
	}
	log.Debugf("rangeMap is %v\n", rangeMap)

	return sdk.Success(rangeMapBytes)
}

// type depositFundsResponse struct {
// 	message string
// 	address	string
// 	amount	int
// 	balance	int
// }

func depositFunds(stubInterface sdk.ContractStub, kind []byte, args [][]byte) common.InvocationResponse {
	// 1. 判断权限
	if string(kind) != "bankKey" {
		return sdk.Error("权限不足")
	}
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
	info := fmt.Sprintf("质押成功,地址:%s,质押金额:%d,余额:%d", address, value, intOldValue + value)
	return sdk.Success([]byte(info))
	// depositFundsResponse := &depositFundsResponse {message :"质押资金成功！", address :address, amount :value, balance :intOldValue + value}
	// responseData, err := depositFundsResponse.Marshal()
	// data := {
	// 	"address" :address,
	// 	"amount" :value,
	// 	"balance" :intOldValue + value
	// }
	// responseData, err := json.Marshal(depositFundsResponse)
	// fmt.Printf(string(responseData))
	// return sdk.Success(responseData)
}

// type deductFundsResponse struct {
// 	message string
// 	address	string
// 	amount	int
// 	balance	int
// }

// 抵扣资金
func deductFunds(stubInterface sdk.ContractStub, kind []byte, args [][]byte) common.InvocationResponse {
	// 1. 判断权限
	if string(kind) != "bankKey" {
		return sdk.Error("权限不足")
	}
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
	info := fmt.Sprintf("抵扣成功,地址:%s,抵扣金额:%d,余额:%d", address, value, intOldValue - value)
	return sdk.Success([]byte(info))
	// deductFundsResponse := &deductFundsResponse{message :"质押抵扣成功！", address :address, amount :value, balance :intOldValue - value}
	// responseData, err := json.Marshal(deductFundsResponse)
	// return sdk.Success(responseData)
}

// type transferFundsResponse struct {
// 	message string
// 	addressFrom	string
// 	amount	int
// 	balanceFrom	int
// 	addressTo	string
// 	balanceTo	int
// }

// 转移资金
func transferFunds(stubInterface sdk.ContractStub, kind []byte, args [][]byte) common.InvocationResponse {
	// 1. 判断权限
	if string(kind) != "bankKey" {
		return sdk.Error("权限不足")
	}
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
		info := fmt.Sprintf("余额不足，转账失败！ 地址:%s，余额:%d，需要额度:%d ", addressFrom, intOldValueFrom, value)
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

	info := fmt.Sprintf("转账成功,From地址:%s,抵扣金额:%d,余额:%d,To地址:%s,余额:%d", addressFrom, value, intOldValueFrom - value, addressTo, intOldValueTo + value)
	return sdk.Success([]byte(info))
	// transferFundsResponse := &transferFundsResponse{message:"转账成功！", addressFrom:addressFrom, amount:value, balanceFrom:intOldValueFrom - value, addressTo:addressTo, balanceTo:intOldValueTo + value}
	// responseData, err := json.Marshal(transferFundsResponse)
	// return sdk.Success(responseData)
}

// type getFundsResponse struct {
// 	message string
// 	address	string
// 	balance	int
// }

// 查询资金余额
func getFunds(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
	if len(args) < 1 {
		fmt.Println("The args number is not correct")
		return sdk.Error("The args is not correct.")
	}

	address := string(args[0])
	value, err := stubInterface.GetKV(address)
	if err != nil {
		errInfo := fmt.Sprintf("Get the key:%s failed", address)
		return sdk.Error(errInfo)
	}
	intValue,_ := strconv.Atoi(string(value))
	info := fmt.Sprintf("查询成功,地址:%s,余额:%d", address, intValue)
	return sdk.Success([]byte(info))
	// getFundsResponse := &getFundsResponse{message:"查询成功！", address:address, balance:intValue}
	// responseData, err := json.Marshal(getFundsResponse)
	// return sdk.Success(responseData)
}

// type initFundsResponse struct {
// 	message string
// 	address	string
// 	balance	int
// }

// 初始化ETC
func initFunds(stubInterface sdk.ContractStub, kind []byte, args [][]byte) common.InvocationResponse {
	// 1. 判断权限
	if string(kind) != "bankKey" {
		return sdk.Error("权限不足")
	}
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
	info := fmt.Sprintf("初始化成功,地址:%s,余额:%s ", address, value)
	return sdk.Success([]byte(info))
	// intValue,_ := strconv.Atoi(string(value))
	// initFundsResponse := &initFundsResponse{message:"初始化成功！", address:address, balance:intValue}
	// responseData, err := json.Marshal(initFundsResponse)
	// return sdk.Success(responseData)
}

// 获取交易记录
// func getFundsHistory(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
// 	if len(args) < 1 {
// 		fmt.Println("The args number is not correct")
// 		return sdk.Error("The args is not correct.")
// 	}

// 	address := string(args[0])
// 	iterator, err := stubInterface.GetKeyHistoryIterator(address)
// 	if err != nil {
// 		errInfo := fmt.Sprintf("Get the key:%s failed", address)
// 		return sdk.Error(errInfo)
// 	}
// 	defer iterator.Close()

// 	var historyArray []int 
// 	var count = 0
// 	for {
// 		b := iterator.Next()
// 		if b {
// 			key := iterator.Key()
// 			value := string(iterator.Value())
// 			intValue,err := strconv.Atoi(value)
// 			if err != nil {
// 				fmt.Printf("The ETC size is not int type\n")
// 				return sdk.Error("The ETC size is not int type")
// 			}
// 			historyArray = append(historyArray, intValue)
// 			count++
// 			log.Debugf("The iterator read key is %s, value is %s, count is %d\n", key, value, count)
// 		} else {
// 			log.Debugf("The iterator break\n")
// 			break
// 		}
// 	}
// 	historyBytes, err := json.Marshal(historyArray)
// 	if err != nil {
// 		return sdk.Error(err.Error())
// 	}

// 	return sdk.Success(historyBytes)
// }
func getFundsHistory(stubInterface sdk.ContractStub, args [][]byte) common.InvocationResponse {
    // stubInterface.Debugf("Enter getKeyHistory")
    const numOfArgs = 1
    if len(args) != numOfArgs {
        // stubInterface.Errorf("the args for getKeyHistory is not correct")
        return sdk.Error("the args for getKeyHistory is not correct")
    }

    key := string(args[0])

    iterator, err := stubInterface.GetKeyHistoryIterator(key)
    if err != nil {
        return sdk.Error("GetKeyHistoryIterator 出错！")
		// return sdk.Error(err.Error())
    }
    defer iterator.Close()

    var historyArray []keyHistory
    for {
        b := iterator.Next()
        if !b {
            // stubInterface.Debugf("the iterator break")
            break
        }

        var history keyHistory
        history.Value = string(iterator.Value())
        // history.TxHash = iterator.TxHash()
        history.BlockNum, history.TxNum = iterator.Version()
        history.IsDeleted = iterator.IsDeleted()
        // history.Timestamp = iterator.Timestamp()
        historyArray = append(historyArray, history)
    }

    historyMapBytes, err := json.Marshal(historyArray)
    if err != nil {
        return sdk.Error(err.Error())
    }
    // stub.Debugf("historyArray is %v", historyArray)

    return sdk.Success(historyMapBytes)
}



func (dc DigitalCertificate) Marshal() ([]byte, error) {
	return json.Marshal(dc)
}

// Unmarshal Unmarshal json data.
func Unmarshal(data []byte) (*DigitalCertificate, error) {
	var dc DigitalCertificate
	err := json.Unmarshal(data, &dc)
	if err != nil {
		return nil, err
	}
	return &dc, nil
}

func main() {
	smstub.Start(&etcManager{})
}

