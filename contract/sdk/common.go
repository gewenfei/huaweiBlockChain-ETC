/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

// Package sdk for smart contract.
package sdk

import (
	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/logger"
	"git.huawei.com/poissonsearch/wienerchain/proto/common"
)

// Log the log for smart contract.
var Log = logger.GetDefaultLogger()

// Success construct success response.
func Success(payload []byte) common.InvocationResponse {
	res := common.InvocationResponse{Status: common.Status_SUCCESS, Payload: payload}
	return res
}

// Error construct error response.
func Error(errMsg string) common.InvocationResponse {
	res := common.InvocationResponse{Status: common.Status_BAD_REQUEST, StatusInfo: errMsg}
	return res
}
