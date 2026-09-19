package zgen

import (
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/znxlc/zerror"

	"github.com/shopspring/decimal"
)

// Int - tries to convert any to int (conversion loss may occur)
func Int(src any) (dst int, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case int:
		return val, nil
	case int64:
		if val > int64(math.MaxInt) || val < int64(math.MinInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(val), nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case uint:
		if uint64(val) > uint64(math.MaxInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		if uint64(val) > uint64(math.MaxInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint32",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(val), nil
	case uint64:
		if val > uint64(math.MaxInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(val), nil
	case float32:
		if math.IsNaN(float64(val)) || math.IsInf(float64(val), 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "int",
				"value":     val,
			})
		}
		if val > float32(math.MaxInt) || val < float32(math.MinInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(val), nil
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "int",
				"value":     val,
			})
		}
		if val > float64(math.MaxInt) || val < float64(math.MinInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "int",
				"value":     val,
			})
		}
		if float64(realVal) > float64(math.MaxInt) || float64(realVal) < float64(math.MinInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "int",
				"value":     val,
			})
		}
		if realVal > float64(math.MaxInt) || realVal < float64(math.MinInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(realVal), nil
	case time.Duration:
		ival := int64(val)
		if ival > int64(math.MaxInt) || ival < int64(math.MinInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(ival), nil
	case time.Time: // return the unix value
		return int(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "int",
				"error":    err.Error(),
			})
		}
		return int(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int",
				"error":    err.Error(),
			})
		}
		return int(floatVal), nil
	case decimal.Decimal:
		bi := val.BigInt()
		if !bi.IsInt64() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int",
				"value":     val,
			})
		}
		i64 := bi.Int64()
		if i64 > int64(math.MaxInt) || i64 < int64(math.MinInt) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int",
				"value":     val,
			})
		}
		return int(i64), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "int",
		})
	}
}

// Uint - tries to convert any to uint (conversion loss may occur)
func Uint(src any) (dst uint, err zerror.Error) {
	if src == nil {
		return 0, nil
	}

	switch val := src.(type) {
	case uint:
		return val, nil
	case int:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case int64:
		if val < 0 || uint64(val) > uint64(math.MaxUint) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case int8:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int8",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case int16:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int16",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case int32:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int32",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case uint8:
		return uint(val), nil
	case uint16:
		return uint(val), nil
	case uint32:
		return uint(val), nil
	case uint64:
		if val > uint64(math.MaxUint) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case float32:
		if math.IsNaN(float64(val)) || math.IsInf(float64(val), 0) || val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "uint",
				"value":     val,
			})
		}
		if val > float32(math.MaxUint) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) || val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "uint",
				"value":     val,
			})
		}
		if val > float64(math.MaxUint) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint",
				"value":     val,
			})
		}
		if float64(realVal) > float64(^uint(0)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint",
				"value":     val,
			})
		}
		if realVal > float64(^uint(0)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(realVal), nil
	case time.Duration:
		ival := int64(val)
		if ival < 0 || uint64(val) > math.MaxUint {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(ival), nil
	case time.Time: // return the unix value
		return uint(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "uint",
				"error":    err.Error(),
			})
		}
		return uint(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint",
				"error":    err.Error(),
			})
		}
		if floatVal < 0 || floatVal > math.MaxUint {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(floatVal), nil
	case decimal.Decimal:
		if val.IsNegative() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint",
				"value":     val,
			})
		}
		bi := val.BigInt()
		if !bi.IsUint64() || bi.Uint64() > uint64(math.MaxUint) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint",
				"value":     val,
			})
		}
		return uint(bi.Uint64()), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "uint",
		})
	}
}

// Int64 - tries to convert any to int64
func Int64(src any) (dst int64, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case int64:
		return val, nil
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case complex64:
		return int64(real(val)), nil
	case complex128:
		return int64(real(val)), nil
	case time.Duration:
		return int64(val), nil
	case time.Time: // return the unix value
		return val.Unix(), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		intVal, er := strconv.ParseInt(string(val), 10, 64)
		if er != nil {
			floatVal, er := strconv.ParseFloat(string(val), 64)
			if er != nil {
				return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
					"src":      val,
					"src_type": "[]byte",
					"dst_type": "int64",
					"error":    er.Error(),
				})
			}
			return int64(floatVal), nil
		}
		return intVal, nil
	case string:
		intVal, er := strconv.ParseInt(val, 10, 64)
		if er != nil {
			floatVal, er := strconv.ParseFloat(val, 64)
			if er != nil {
				return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
					"src":      val,
					"src_type": "string",
					"dst_type": "int64",
					"error":    er.Error(),
				})
			}
			return int64(floatVal), nil
		}
		return intVal, nil
	case decimal.Decimal:
		bi := val.BigInt()
		if !bi.IsInt64() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int64",
				"value":     val,
			})
		}
		return bi.Int64(), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "int64",
		})
	}
}

// Int32 - tries to convert any to int32, data may be lost in the conversion so use at your own risk
func Int32(src any) (dst int32, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case int32:
		return val, nil
	case int:
		dst = int32(val)
		if int64(val) > math.MaxInt32 || int64(val) < math.MinInt32 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "int32",
				"value":     val,
			})
			dst = 0
		}
		return
	case int8:
		return int32(val), nil
	case int16:
		return int32(val), nil
	case int64:
		dst = int32(val)
		if val > math.MaxInt32 || val < math.MinInt32 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "int32",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint:
		dst = int32(val)
		if val > math.MaxInt32 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint",
				"to_type":   "int32",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint8:
		return int32(val), nil
	case uint16:
		return int32(val), nil
	case uint32:
		dst = int32(val)
		if val > math.MaxInt32 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint32",
				"to_type":   "int32",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint64:
		dst = int32(val)
		if val > math.MaxInt32 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "int32",
				"value":     val,
			})
			dst = 0
		}
		return
	case float32:
		if val > float32(math.MaxInt32) || val < float32(math.MinInt32) || math.IsNaN(float64(val)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "int32",
				"value":     val,
			})
		}
		return int32(val), nil
	case float64:
		if val > float64(math.MaxInt32) || val < float64(math.MinInt32) || math.IsNaN(val) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "int32",
				"value":     val,
			})
		}
		return int32(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "int32",
				"value":     val,
			})
		}
		if realVal > math.MaxInt32 || realVal < math.MinInt32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "int32",
				"value":     val,
			})
		}
		return int32(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "int32",
				"value":     val,
			})
		}
		if realVal > math.MaxInt32 || realVal < math.MinInt32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "int32",
				"value":     val,
			})
		}
		return int32(realVal), nil
	case time.Duration:
		ival := int64(val)
		if ival > int64(math.MaxInt32) || ival < int64(math.MinInt32) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "int32",
				"value":     val,
			})
		}
		return int32(ival), nil
	case time.Time: // return the unix value, some conversion loss may occur
		return int32(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 32)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "int32",
				"error":    err.Error(),
			})
		}
		return int32(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int32",
				"error":    err.Error(),
			})
		}
		if floatVal > float64(math.MaxInt32) || floatVal < float64(math.MinInt32) || math.IsNaN(floatVal) || math.IsInf(floatVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "int32",
				"value":     val,
			})
		}
		return int32(floatVal), nil
	case decimal.Decimal:
		bi := val.BigInt()
		if !bi.IsInt64() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int32",
				"value":     val,
			})
		}
		i64 := bi.Int64()
		if i64 > math.MaxInt32 || i64 < math.MinInt32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int32",
				"value":     val,
			})
		}
		return int32(i64), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "int32",
		})
	}
}

// Int16 - tries to convert any to int16 (conversion loss may occur)
func Int16(src any) (dst int16, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case int:
		dst = int16(val)
		if int64(val) > math.MaxInt16 || int64(val) < math.MinInt16 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "int16",
				"value":     val,
			})
			dst = 0
		}
		return
	case int64:
		dst = int16(val)
		if val > math.MaxInt16 || val < math.MinInt16 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "int16",
				"value":     val,
			})
			dst = 0
		}
		return
	case int8:
		return int16(val), nil
	case int16:
		return val, nil
	case int32:
		dst = int16(val)
		if val > math.MaxInt16 || val < math.MinInt16 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int32",
				"to_type":   "int16",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint:
		dst = int16(val)
		if val > math.MaxInt16 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint",
				"to_type":   "int16",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint8:
		return int16(val), nil
	case uint16:
		dst = int16(val)
		if val > math.MaxInt16 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint16",
				"to_type":   "int16",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint32:
		dst = int16(val)
		if val > math.MaxInt16 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint32",
				"to_type":   "int16",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint64:
		dst = int16(val)
		if val > math.MaxInt16 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "int16",
				"value":     val,
			})
			dst = 0
		}
		return
	case float32:
		if val > float32(math.MaxInt16) || val < float32(math.MinInt16) || math.IsNaN(float64(val)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "int16",
				"value":     val,
			})
		}
		return int16(val), nil
	case float64:
		if val > float64(math.MaxInt16) || val < float64(math.MinInt16) || math.IsNaN(val) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "int16",
				"value":     val,
			})
		}
		return int16(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "int16",
				"value":     val,
			})
		}
		if realVal > math.MaxInt16 || realVal < math.MinInt16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "int16",
				"value":     val,
			})
		}
		return int16(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "int16",
				"value":     val,
			})
		}
		if realVal > math.MaxInt16 || realVal < math.MinInt16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "int16",
				"value":     val,
			})
		}
		return int16(realVal), nil
	case time.Duration:
		ival := int64(val)
		if ival > int64(math.MaxInt16) || ival < int64(math.MinInt16) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "int16",
				"value":     val,
			})
		}
		return int16(ival), nil
	case time.Time: // return the unix value
		return int16(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "int16",
				"error":    err.Error(),
			})
		}
		return int16(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int16",
				"error":    err.Error(),
			})
		}
		if floatVal > float64(math.MaxInt16) || floatVal < float64(math.MinInt16) || math.IsNaN(floatVal) || math.IsInf(floatVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "int16",
				"value":     val,
			})
		}
		return int16(floatVal), nil
	case decimal.Decimal:
		bi := val.BigInt()
		if !bi.IsInt64() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int16",
				"value":     val,
			})
		}
		i64 := bi.Int64()
		if i64 > math.MaxInt16 || i64 < math.MinInt16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int16",
				"value":     val,
			})
		}
		return int16(i64), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "int16",
		})
	}
}

// Int8 - tries to convert any to int8 (conversion loss may occur)
func Int8(src any) (dst int8, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case int:
		dst = int8(val)
		if val > math.MaxInt8 || val < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case int64:
		dst = int8(val)
		if val > math.MaxInt8 || val < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case int8:
		return val, nil
	case int16:
		dst = int8(val)
		if val > math.MaxInt8 || val < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case int32:
		dst = int8(val)
		if val > math.MaxInt8 || val < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int32",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint:
		dst = int8(val)
		if val > math.MaxInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint8:
		dst = int8(val)
		if val > math.MaxInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint8",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint16:
		dst = int8(val)
		if val > math.MaxInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint16",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint32:
		dst = int8(val)
		if val > math.MaxInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint32",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case uint64:
		dst = int8(val)
		if val > math.MaxInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case float32:
		dst = int8(val)
		if val > math.MaxInt8 || val < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case float64:
		dst = int8(val)
		if val > math.MaxInt8 || val < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case complex64:
		rVal := real(val)
		dst = int8(rVal)
		if rVal > math.MaxInt8 || rVal < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case complex128:
		rVal := real(val)
		dst = int8(rVal)
		if rVal > math.MaxInt8 || rVal < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case time.Duration:
		ival := int64(val)
		if ival > int64(math.MaxInt8) || ival < int64(math.MinInt8) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "int8",
				"value":     val,
			})
		}
		return int8(ival), nil
	case time.Time: // return the unix value
		rVal := val.Unix()
		dst = int8(rVal)
		if rVal > math.MaxInt8 || rVal < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return

	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, er := strconv.ParseFloat(string(val), 64)
		if er != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "int8",
				"error":    er.Error(),
			})
		}
		dst = int8(floatVal)
		if floatVal > math.MaxInt8 || floatVal < math.MinInt8 {
			err = zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "int8",
				"value":     val,
			})
			dst = 0
		}
		return
	case string:
		floatVal, er := strconv.ParseFloat(val, 64)
		if er != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int8",
				"error":    er.Error(),
			})
		}
		if floatVal > float64(math.MaxInt8) || floatVal < float64(math.MinInt8) || math.IsNaN(floatVal) || math.IsInf(floatVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "int8",
				"value":     val,
			})
		}
		return int8(floatVal), nil
	case decimal.Decimal:
		bi := val.BigInt()
		if !bi.IsInt64() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int8",
				"value":     val,
			})
		}
		i64 := bi.Int64()
		if i64 > math.MaxInt8 || i64 < math.MinInt8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "int8",
				"value":     val,
			})
		}
		return int8(i64), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "int8",
		})
	}
}

// Uint64 - tries to convert any to uint64
func Uint64(src any) (dst uint64, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case uint64:
		return val, nil
	case int:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case int8:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int8",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case int16:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int16",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case int32:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int32",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case int64:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case uint:
		return uint64(val), nil
	case uint8:
		return uint64(val), nil
	case uint16:
		return uint64(val), nil
	case uint32:
		return uint64(val), nil
	case float32:
		if math.IsNaN(float64(val)) || math.IsInf(float64(val), 0) || val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		if val > float32(math.MaxUint64) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) || val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		if val > float64(^uint64(0)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		if realVal > float32(^uint64(0)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		if realVal > float64(^uint64(0)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(realVal), nil
	case time.Duration:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return uint64(val), nil
	case time.Time: // return the unix value
		return uint64(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		// First try parsing as integer
		intVal, err := strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			// If not an integer, try parsing as float
			floatVal, ferr := strconv.ParseFloat(string(val), 64)
			if ferr != nil {
				return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
					"src":      val,
					"src_type": "[]byte",
					"dst_type": "uint64",
					"error":    err.Error(),
				})
			}
			if floatVal < 0 || floatVal > float64(^uint64(0)) || math.IsNaN(floatVal) || math.IsInf(floatVal, 0) {
				return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
					"from_type": "[]byte",
					"to_type":   "uint64",
					"value":     val,
				})
			}
			return uint64(floatVal), nil
		}
		return intVal, nil
	case string:
		// First try parsing as integer
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			// If not an integer, try parsing as float
			floatVal, ferr := strconv.ParseFloat(val, 64)
			if ferr != nil {
				return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
					"src":      val,
					"src_type": "string",
					"dst_type": "uint64",
					"error":    err.Error(),
				})
			}
			if floatVal < 0 || floatVal > float64(^uint64(0)) || math.IsNaN(floatVal) || math.IsInf(floatVal, 0) {
				return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
					"from_type": "string",
					"to_type":   "uint64",
					"value":     val,
				})
			}
			return uint64(floatVal), nil
		}
		return intVal, nil
	case decimal.Decimal:
		if val.IsNegative() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		bi := val.BigInt()
		if !bi.IsUint64() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint64",
				"value":     val,
			})
		}
		return bi.Uint64(), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "uint64",
		})
	}
}

// Uint32 - tries to convert any to uint32 (conversion loss may occur)
func Uint32(src any) (dst uint32, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case uint32:
		return val, nil
	case int:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case int8:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int8",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case int16:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int16",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case int32:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int32",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case int64:
		if val < 0 || val > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case uint:
		if val > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case uint8:
		return uint32(val), nil
	case uint16:
		return uint32(val), nil
	case uint64:
		if val > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case float32:
		if val < 0 || val > math.MaxUint32 || math.IsNaN(float64(val)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case float64:
		if val < 0 || val > math.MaxUint32 || math.IsNaN(val) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		if realVal > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		if realVal > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(realVal), nil
	case time.Duration:
		ival := int64(val)
		if ival < 0 || ival > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(ival), nil
	case time.Time: // return the unix value
		unixTime := val.Unix()
		if unixTime < 0 || unixTime > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Time",
				"to_type":   "uint32",
				"value":     unixTime,
			})
		}
		return uint32(unixTime), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		// First try parsing as uint64
		if uintVal, err := strconv.ParseUint(string(val), 10, 32); err == nil {
			return uint32(uintVal), nil
		}

		// If that fails, try parsing as float
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "uint32",
				"error":    err.Error(),
			})
		}

		// Check for overflow, NaN, and negative values
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal < 0 || floatVal > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "uint32",
				"value":     string(val),
			})
		}

		return uint32(floatVal), nil
	case string:
		// First try parsing as uint64
		if uintVal, err := strconv.ParseUint(val, 10, 32); err == nil {
			return uint32(uintVal), nil
		}

		// If that fails, try parsing as float
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint32",
				"error":    err.Error(),
			})
		}

		// Check for overflow, NaN, and negative values
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal < 0 || floatVal > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "uint32",
				"value":     val,
			})
		}

		return uint32(floatVal), nil
	case decimal.Decimal:
		if val.IsNegative() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		bi := val.BigInt()
		if !bi.IsUint64() || bi.Uint64() > math.MaxUint32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint32",
				"value":     val,
			})
		}
		return uint32(bi.Uint64()), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "uint32",
		})
	}
}

// Uint16 - tries to convert any to uint16 (conversion loss may occur)
func Uint16(src any) (dst uint16, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case uint16:
		return val, nil
	case int:
		if val < 0 || val > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case int8:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int8",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case int16:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int16",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case int32:
		if val < 0 || val > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int32",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case int64:
		if val < 0 || val > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case uint:
		if val > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case uint8:
		return uint16(val), nil
	case uint32:
		if val > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint32",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case uint64:
		if val > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case float32:
		if val < 0 || val > math.MaxUint16 || math.IsNaN(float64(val)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case float64:
		if val < 0 || val > math.MaxUint16 || math.IsNaN(val) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		if realVal > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		if realVal > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(realVal), nil
	case time.Duration:
		ival := int64(val)
		if ival < 0 || ival > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(ival), nil
	case time.Time: // return the unix value
		unixTime := val.Unix()
		if unixTime < 0 || unixTime > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Time",
				"to_type":   "uint16",
				"value":     unixTime,
			})
		}
		return uint16(unixTime), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		// First try parsing as uint64
		if uintVal, err := strconv.ParseUint(string(val), 10, 16); err == nil {
			return uint16(uintVal), nil
		}

		// If that fails, try parsing as float
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "uint16",
				"error":    err.Error(),
			})
		}

		// Check for overflow, NaN, and negative values
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal < 0 || floatVal > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "uint16",
				"value":     string(val),
			})
		}

		return uint16(floatVal), nil
	case string:
		// First try parsing as uint64
		if uintVal, err := strconv.ParseUint(val, 10, 16); err == nil {
			return uint16(uintVal), nil
		}

		// If that fails, try parsing as float
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint16",
				"error":    err.Error(),
			})
		}

		// Check for overflow, NaN, and negative values
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal < 0 || floatVal > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "uint16",
				"value":     val,
			})
		}

		return uint16(floatVal), nil
	case decimal.Decimal:
		if val.IsNegative() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		bi := val.BigInt()
		if !bi.IsUint64() || bi.Uint64() > math.MaxUint16 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint16",
				"value":     val,
			})
		}
		return uint16(bi.Uint64()), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "uint16",
		})
	}
}

// Uint8 - tries to convert any to uint8 (conversion loss may occur)
func Uint8(src any) (dst uint8, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case uint8:
		return val, nil
	case int:
		if val < 0 || val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case int8:
		if val < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int8",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case int16:
		if val < 0 || val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int16",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case int32:
		if val < 0 || val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int32",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case int64:
		if val < 0 || val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case uint:
		if val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case uint16:
		if val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint16",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case uint32:
		if val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint32",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case uint64:
		if val > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case float32:
		if val < 0 || val > math.MaxUint8 || math.IsNaN(float64(val)) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case float64:
		if val < 0 || val > math.MaxUint8 || math.IsNaN(val) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		if realVal > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || realVal < 0 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		if realVal > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(realVal), nil
	case time.Duration:
		ival := int64(val)
		if ival < 0 || ival > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Duration",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(ival), nil
	case time.Time: // return the unix value
		unixTime := val.Unix()
		if unixTime < 0 || unixTime > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "time.Time",
				"to_type":   "uint8",
				"value":     unixTime,
			})
		}
		return uint8(unixTime), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		// First try parsing as uint64
		if uintVal, err := strconv.ParseUint(string(val), 10, 8); err == nil {
			return uint8(uintVal), nil
		}

		// If that fails, try parsing as float
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "uint8",
				"error":    err.Error(),
			})
		}

		// Check for overflow, NaN, and negative values
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal < 0 || floatVal > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "uint8",
				"value":     string(val),
			})
		}

		return uint8(floatVal), nil
	case string:
		// First try parsing as uint64
		if uintVal, err := strconv.ParseUint(val, 10, 8); err == nil {
			return uint8(uintVal), nil
		}

		// If that fails, try parsing as float
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint8",
				"error":    err.Error(),
			})
		}

		// Check for overflow, NaN, and negative values
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal < 0 || floatVal > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "uint8",
				"value":     val,
			})
		}

		return uint8(floatVal), nil
	case decimal.Decimal:
		if val.IsNegative() {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		bi := val.BigInt()
		if !bi.IsUint64() || bi.Uint64() > math.MaxUint8 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "uint8",
				"value":     val,
			})
		}
		return uint8(bi.Uint64()), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "uint8",
		})
	}
}

// Float64 - tries to convert any to float64
func Float64(src any) (dst float64, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		// Check for potential precision loss when converting from int64 to float64
		// 1<<53 (2^53) is the largest integer that can be exactly represented in a float64 (IEEE 754 double-precision)
		if val > int64(1<<53) || val < -int64(1<<53) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "float64",
				"value":     val,
			})
		}
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		// Check for potential precision loss when converting from uint64 to float64
		// 1<<53 is the largest integer that can be exactly represented in a float64 (IEEE 754 double-precision)
		if val > uint64(1<<53) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "float64",
				"value":     val,
			})
		}
		return float64(val), nil
	case float32:
		if math.IsNaN(float64(val)) || math.IsInf(float64(val), 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float32",
				"to_type":   "float64",
				"value":     val,
			})
		}
		return float64(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "float64",
				"value":     val,
			})
		}
		return float64(realVal), nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "float64",
				"value":     val,
			})
		}
		return realVal, nil
	case time.Duration:
		return float64(val), nil
	case time.Time: // return the unix value
		return float64(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatValue, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "float64",
				"error":    err.Error(),
			})
		}
		if math.IsNaN(floatValue) || math.IsInf(floatValue, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "float64",
				"value":     string(val),
			})
		}
		return floatValue, nil
	case string:
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "float64",
				"error":    err.Error(),
			})
		}
		if math.IsNaN(floatValue) || math.IsInf(floatValue, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "float64",
				"value":     val,
			})
		}
		return floatValue, nil
	case decimal.Decimal:
		f, _ := val.Float64()
		if math.IsInf(f, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "float64",
				"value":     val,
			})
		}
		return f, nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "float64",
		})
	}
}

// Float32 - tries to convert any to float32(conversion loss may occur)
func Float32(src any) (dst float32, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case float32:
		return val, nil
	case int:
		return float32(val), nil
	case int8:
		return float32(val), nil
	case int16:
		return float32(val), nil
	case int32:
		return float32(val), nil
	case int64:
		// Check for potential precision loss when converting from int64 to float32
		// 1<<24 is the largest integer that can be exactly represented in a float32 (IEEE 754 single-precision)
		if val > int64(1<<24) || val < -int64(1<<24) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "float32",
				"value":     val,
			})
		}
		return float32(val), nil
	case uint:
		return float32(val), nil
	case uint8:
		return float32(val), nil
	case uint16:
		return float32(val), nil
	case uint32:
		return float32(val), nil
	case uint64:
		// Check for potential precision loss when converting from uint64 to float32
		// 1<<24 is the largest integer that can be exactly represented in a float32 (IEEE 754 single-precision)
		if val > uint64(1<<24) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "float32",
				"value":     val,
			})
		}
		return float32(val), nil
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) || val > math.MaxFloat32 || val < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "float32",
				"value":     val,
			})
		}
		return float32(val), nil
	case complex64:
		realVal := real(val)
		if math.IsNaN(float64(realVal)) || math.IsInf(float64(realVal), 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex64",
				"to_type":   "float32",
				"value":     val,
			})
		}
		return realVal, nil
	case complex128:
		realVal := real(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "float32",
				"value":     val,
			})
		}
		if realVal > math.MaxFloat32 || realVal < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "float32",
				"value":     val,
			})
		}
		return float32(realVal), nil
	case time.Duration:
		return float32(val), nil
	case time.Time: // return the unix value
		return float32(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 32)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "float32",
				"error":    err.Error(),
			})
		}
		// Check for overflow, NaN, and Inf after parsing
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal > math.MaxFloat32 || floatVal < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "float32",
				"value":     string(val),
			})
		}
		return float32(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "float32",
				"error":    err.Error(),
			})
		}
		// Check for overflow, NaN, and Inf after parsing
		if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) || floatVal > math.MaxFloat32 || floatVal < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "float32",
				"value":     val,
			})
		}
		return float32(floatVal), nil
	case decimal.Decimal:
		f, _ := val.Float64()
		if math.IsInf(f, 0) || f > math.MaxFloat32 || f < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "float32",
				"value":     val,
			})
		}
		return float32(f), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "string",
			"dst_type": "float32",
		})
	}
}

// Complex64 - tries to convert any to complex64(conversion loss may occur because complex64 uses 2 float32 behind the scenes)
func Complex64(src any) (dst complex64, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case complex64:
		return val, nil
	case int:
		return complex(float32(val), float32(0)), nil
	case int8:
		return complex(float32(val), float32(0)), nil
	case int16:
		return complex(float32(val), float32(0)), nil
	case int32:
		return complex(float32(val), float32(0)), nil
	case int64:
		// Check for potential overflow when converting from int64 to float32
		if val > int64(1<<24) || val < -int64(1<<24) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "complex64",
				"value":     val,
			})
		}
		return complex(float32(val), float32(0)), nil
	case uint:
		return complex(float32(val), float32(0)), nil
	case uint8:
		return complex(float32(val), float32(0)), nil
	case uint16:
		return complex(float32(val), float32(0)), nil
	case uint32:
		return complex(float32(val), float32(0)), nil
	case uint64:
		// Check for potential overflow when converting from uint64 to float32
		if val > uint64(1<<24) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "complex64",
				"value":     val,
			})
		}
		return complex(float32(val), float32(0)), nil
	case float32:
		return complex(val, float32(0)), nil
	case float64:
		// Check for overflow, NaN, and Inf when converting from float64 to float32
		if math.IsNaN(val) || math.IsInf(val, 0) || val > math.MaxFloat32 || val < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "float64",
				"to_type":   "complex64",
				"value":     val,
			})
		}
		return complex(float32(val), float32(0)), nil
	case complex128:
		// Check for overflow, NaN, and Inf in both real and imaginary parts
		realVal := real(val)
		imagVal := imag(val)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || realVal > math.MaxFloat32 || realVal < -math.MaxFloat32 ||
			math.IsNaN(imagVal) || math.IsInf(imagVal, 0) || imagVal > math.MaxFloat32 || imagVal < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "complex128",
				"to_type":   "complex64",
				"value":     val,
			})
		}
		return complex64(val), nil
	case time.Duration:
		return complex(float32(val), float32(0)), nil
	case time.Time: // return the unix value
		return complex(float32(val.Unix()), float32(0)), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		complexVal, err := strconv.ParseComplex(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "complex64",
				"error":    err.Error(),
			})
		}
		// Check for overflow in the parsed complex value
		realVal := real(complexVal)
		imagVal := imag(complexVal)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || realVal > math.MaxFloat32 || realVal < -math.MaxFloat32 ||
			math.IsNaN(imagVal) || math.IsInf(imagVal, 0) || imagVal > math.MaxFloat32 || imagVal < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "complex64",
				"value":     string(val),
			})
		}
		return complex64(complexVal), nil
	case string:
		complexVal, err := strconv.ParseComplex(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "complex64",
				"error":    err.Error(),
			})
		}
		return complex64(complexVal), nil
	case decimal.Decimal:
		f, _ := val.Float64()
		if math.IsInf(f, 0) || f > math.MaxFloat32 || f < -math.MaxFloat32 {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "complex64",
				"value":     val,
			})
		}
		return complex(float32(f), float32(0)), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "complex64",
		})
	}
}

// Complex128 - tries to convert any to complex128
func Complex128(src any) (dst complex128, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case complex128:
		return val, nil
	case int:
		return complex(float64(val), float64(0)), nil
	case int8:
		return complex(float64(val), float64(0)), nil
	case int16:
		return complex(float64(val), float64(0)), nil
	case int32:
		return complex(float64(val), float64(0)), nil
	case int64:
		// Check for potential precision loss when converting from int64 to float64
		// 1<<53 is the largest integer that can be exactly represented in a float64
		if val > int64(1<<53) || val < -int64(1<<53) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "int64",
				"to_type":   "complex128",
				"value":     val,
			})
		}
		return complex(float64(val), float64(0)), nil
	case uint:
		return complex(float64(val), float64(0)), nil
	case uint8:
		return complex(float64(val), float64(0)), nil
	case uint16:
		return complex(float64(val), float64(0)), nil
	case uint32:
		return complex(float64(val), float64(0)), nil
	case uint64:
		// Check for potential precision loss when converting from uint64 to float64
		// 1<<53 is the largest integer that can be exactly represented in a float64
		if val > uint64(1<<53) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "uint64",
				"to_type":   "complex128",
				"value":     val,
			})
		}
		return complex(float64(val), float64(0)), nil
	case float32:
		return complex(float64(val), float64(0)), nil
	case float64:
		return complex(val, float64(0)), nil
	case complex64:
		// No overflow possible when converting from complex64 to complex128
		return complex128(val), nil // there may be some conversion loss here
	case time.Duration:
		return complex(float64(val), float64(0)), nil
	case time.Time: // return the unix value
		return complex(float64(val.Unix()), float64(0)), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		complexVal, err := strconv.ParseComplex(string(val), 128)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      string(val),
				"src_type": "[]byte",
				"dst_type": "complex128",
				"error":    err.Error(),
			})
		}
		// Check for NaN and Inf in the parsed complex value
		realVal := real(complexVal)
		imagVal := imag(complexVal)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || math.IsNaN(imagVal) || math.IsInf(imagVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "[]byte",
				"to_type":   "complex128",
				"value":     string(val),
			})
		}
		return complexVal, nil
	case string:
		complexVal, err := strconv.ParseComplex(val, 128)
		if err != nil {
			return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "complex128",
				"error":    err.Error(),
			})
		}
		// Check for NaN and Inf in the parsed complex value
		realVal := real(complexVal)
		imagVal := imag(complexVal)
		if math.IsNaN(realVal) || math.IsInf(realVal, 0) || math.IsNaN(imagVal) || math.IsInf(imagVal, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "string",
				"to_type":   "complex128",
				"value":     val,
			})
		}
		return complexVal, nil
	case decimal.Decimal:
		f, _ := val.Float64()
		if math.IsInf(f, 0) {
			return 0, zerror.New(ErrorConvertorNumberOverflow, map[string]any{
				"from_type": "decimal.Decimal",
				"to_type":   "complex128",
				"value":     val,
			})
		}
		return complex(f, float64(0)), nil
	default:
		return 0, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "complex128",
		})
	}
}

// Decimal - tries to convert any to decimal
func Decimal(src any) (dst decimal.Decimal, err zerror.Error) {
	if src == nil {
		return decimal.NewFromInt(0), nil
	}
	switch val := src.(type) {
	case decimal.Decimal:
		return val, nil
	case int:
		return decimal.NewFromInt(int64(val)), nil
	case int8:
		return decimal.NewFromInt(int64(val)), nil
	case int16:
		return decimal.NewFromInt(int64(val)), nil
	case int32:
		return decimal.NewFromInt(int64(val)), nil
	case int64:
		return decimal.NewFromInt(val), nil
	case uint:
		return decimal.NewFromInt(int64(val)), nil
	case uint8:
		return decimal.NewFromInt(int64(val)), nil
	case uint16:
		return decimal.NewFromInt(int64(val)), nil
	case uint32:
		return decimal.NewFromInt(int64(val)), nil
	case uint64:
		return decimal.NewFromInt(int64(val)), nil
	case float32:
		return decimal.NewFromFloat32(val), nil
	case float64:
		return decimal.NewFromFloat(val), nil
	case complex64:
		return decimal.NewFromFloat32(real(val)), nil
	case complex128:
		return decimal.NewFromFloat(real(val)), nil
	case time.Duration:
		return decimal.NewFromInt(int64(val)), nil
	case time.Time: // return the unix value
		return decimal.NewFromInt(int64(val.Unix())), nil
	case bool:
		if val {
			return decimal.NewFromInt(1), nil
		}
		return decimal.NewFromInt(0), nil
	case []byte:
		decval, er := decimal.NewFromString(string(val))
		if er != nil {
			err = zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"error":    er.Error(),
				"src":      val,
				"src_type": reflect.TypeOf(val).String(),
				"dst_type": "decimal",
			})
		}
		return decval, err
	case string:
		decval, er := decimal.NewFromString(val)
		if er != nil {
			err = zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"error":    er.Error(),
				"src":      val,
				"src_type": reflect.TypeOf(val).String(),
				"dst_type": "decimal",
			})
		}
		return decval, err
	default:
		return decimal.NewFromInt(0), zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
			"dst_type": "decimal",
		})
	}
}
