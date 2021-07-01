/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

// Package utils is to construct or convert function.
package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"time"

	"git.huawei.com/poissonsearch/wienerchain/proto/common"
	"github.com/golang/protobuf/ptypes"
)

var keySeperator = []byte{0x00}
var maxEndKeyTail = []byte{0x01}

// ConstructKey is construct key interface .
func ConstructKey(prefix []byte, key []byte) []byte {
	keyPrefix := append(prefix, keySeperator...)
	return append(keyPrefix, key...)
}

// ConstructMaxEndKey is construct max end key interface .
func ConstructMaxEndKey(prefix []byte) []byte {
	return append(prefix, maxEndKeyTail...)
}

// TrimPrefix is trim prefix interface .
func TrimPrefix(key []byte) []byte {
	splits := bytes.SplitN(key, keySeperator, 2)
	if len(splits) > 1 {
		return splits[1]
	}
	return key
}

// ConstructStringKey is construct string key .
func ConstructStringKey(prefix string, key string) string {
	return fmt.Sprintf("%s_%s", prefix, key)
}

// GetSortedKeysFromMap is get sorted keys from map .
func GetSortedKeysFromMap(mapIter interface{}) []string {
	if mapIter == nil {
		return nil
	}

	var keys []string
	elems := reflect.ValueOf(mapIter).MapKeys()
	for _, elem := range elems {
		keys = append(keys, elem.String())
	}
	sort.Strings(keys)
	return keys
}

// ConvertToProtoPrimitives is get sorted keys from map .
func ConvertToProtoPrimitives(valueInterfaces ...interface{}) (*common.PrimitiveValues, error) {
	values := &common.PrimitiveValues{}
	for _, valueInterface := range valueInterfaces {
		value, err := ConvertToProtoPrimitive(valueInterface)
		if err != nil {
			return nil, err
		}
		values.Values = append(values.Values, value)
	}
	return values, nil
}

// ConvertToProtoPrimitive is Convert to Proto Primitive .
func ConvertToProtoPrimitive(valueInterface interface{}) (*common.PrimitiveValue, error) { /// nolint
	primitiveValue := &common.PrimitiveValue{}
	switch value := valueInterface.(type) {
	case float32:
		primitiveValue.Value = &common.PrimitiveValue_FloatValue{FloatValue: value}
	case float64:
		primitiveValue.Value = &common.PrimitiveValue_DoubleValue{DoubleValue: value}
	case int:
		primitiveValue.Value = &common.PrimitiveValue_Int32Value{Int32Value: int32(value)}
	case int8:
		primitiveValue.Value = &common.PrimitiveValue_Int32Value{Int32Value: int32(value)}
	case int16:
		primitiveValue.Value = &common.PrimitiveValue_Int32Value{Int32Value: int32(value)}
	case int32:
		primitiveValue.Value = &common.PrimitiveValue_Int32Value{Int32Value: value}
	case int64:
		primitiveValue.Value = &common.PrimitiveValue_Int64Value{Int64Value: value}
	case uint:
		primitiveValue.Value = &common.PrimitiveValue_Uint32Value{Uint32Value: uint32(value)}
	case uint8:
		primitiveValue.Value = &common.PrimitiveValue_Uint32Value{Uint32Value: uint32(value)}
	case uint16:
		primitiveValue.Value = &common.PrimitiveValue_Uint32Value{Uint32Value: uint32(value)}
	case uint32:
		primitiveValue.Value = &common.PrimitiveValue_Uint32Value{Uint32Value: value}
	case uint64:
		primitiveValue.Value = &common.PrimitiveValue_Uint64Value{Uint64Value: value}
	case bool:
		primitiveValue.Value = &common.PrimitiveValue_BoolValue{BoolValue: value}
	case string:
		primitiveValue.Value = &common.PrimitiveValue_StringValue{StringValue: value}
	case []byte:
		primitiveValue.Value = &common.PrimitiveValue_BytesValue{BytesValue: value}
	case time.Time:
		timeValue, err := ptypes.TimestampProto(value)
		if err != nil {
			return nil, fmt.Errorf("convert to timestamp proto error: %s", err)
		}
		primitiveValue.Value = &common.PrimitiveValue_TimeValue{TimeValue: timeValue}
	default:
		return nil, fmt.Errorf("unexpected type %T", value)
	}
	return primitiveValue, nil
}

// ConvertToGolangPrimitives is Convert to Golang Primitives .
func ConvertToGolangPrimitives(primitiveValues *common.PrimitiveValues) ([]interface{}, error) {
	if primitiveValues == nil {
		return nil, nil
	}
	var values []interface{}
	for _, primitiveValue := range primitiveValues.Values {
		value, err := ConvertToGolangPrimitive(primitiveValue)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

// ConvertToGolangPrimitive is Convert to Golang Primitive .
func ConvertToGolangPrimitive(primitiveValue *common.PrimitiveValue) (interface{}, error) {
	if primitiveValue == nil {
		return nil, nil
	}
	var result interface{}
	switch value := primitiveValue.Value.(type) {
	case *common.PrimitiveValue_FloatValue:
		result = value.FloatValue
	case *common.PrimitiveValue_DoubleValue:
		result = value.DoubleValue
	case *common.PrimitiveValue_Int32Value:
		result = value.Int32Value
	case *common.PrimitiveValue_Int64Value:
		result = value.Int64Value
	case *common.PrimitiveValue_Uint32Value:
		result = value.Uint32Value
	case *common.PrimitiveValue_Uint64Value:
		result = value.Uint64Value
	case *common.PrimitiveValue_BoolValue:
		result = value.BoolValue
	case *common.PrimitiveValue_StringValue:
		result = value.StringValue
	case *common.PrimitiveValue_BytesValue:
		result = value.BytesValue
	case *common.PrimitiveValue_TimeValue:
		timeValue, err := ptypes.Timestamp(value.TimeValue)
		if err != nil {
			return nil, fmt.Errorf("convert to native time error: %s", err)
		}
		result = timeValue
	default:
		return nil, fmt.Errorf("unexpected type %T", value)
	}
	return result, nil
}
