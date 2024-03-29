package zgen

import (
	"github.com/znxlc/zerror"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Valuer - generic interface that implements Value() without error
type Valuer interface {
	Value() any
}

// ValueError - generic interface that implements Value() with error
type ValueError interface {
	Value() (any, error)
}

// Stringable - generic interface that implements String()
type Stringable interface {
	String() string
}

// ToStringable - generic interface that implements ToString()
type ToStringable interface {
	ToString() string
}

// Scanner - generic interface that implements Scan(any) error
type Scanner interface {
	Scan(value any) error
}

// Int - tries to convert any to int (conversion loss may occur)
func Int(src any) (dst int, err zerror.Error) {
	if src == nil {
		return 0, nil
	}
	switch val := src.(type) {
	case int:
		return val, nil
	case int64:
		return int(val), nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case complex64:
		return int(real(val)), nil
	case complex128:
		return int(real(val)), nil
	case time.Duration:
		return int(val), nil
	case time.Time: //return the unix value
		return int(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
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
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int",
				"error":    err.Error(),
			})
		}
		return int(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return uint(val), nil
	case int64:
		return uint(val), nil
	case int8:
		return uint(val), nil
	case int16:
		return uint(val), nil
	case int32:
		return uint(val), nil
	case uint8:
		return uint(val), nil
	case uint16:
		return uint(val), nil
	case uint32:
		return uint(val), nil
	case uint64:
		return uint(val), nil
	case float32:
		return uint(val), nil
	case float64:
		return uint(val), nil
	case complex64:
		return uint(real(val)), nil
	case complex128:
		return uint(real(val)), nil
	case time.Duration:
		return uint(val), nil
	case time.Time: //return the unix value
		return uint(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
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
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint",
				"error":    err.Error(),
			})
		}
		return uint(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
	case time.Time: //return the unix value
		return val.Unix(), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "int64",
				"error":    err.Error(),
			})
		}
		return int64(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int64",
				"error":    err.Error(),
			})
		}
		return int64(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return int32(val), nil
	case int8:
		return int32(val), nil
	case int16:
		return int32(val), nil
	case int64:
		return int32(val), nil
	case uint:
		return int32(val), nil
	case uint8:
		return int32(val), nil
	case uint16:
		return int32(val), nil
	case uint32:
		return int32(val), nil
	case uint64:
		return int32(val), nil
	case float32:
		return int32(val), nil
	case float64:
		return int32(val), nil
	case complex64:
		return int32(real(val)), nil
	case complex128:
		return int32(real(val)), nil
	case time.Duration:
		return int32(val), nil
	case time.Time: //return the unix value, some conversion loss may occur
		return int32(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 32)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
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
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int32",
				"error":    err.Error(),
			})
		}
		return int32(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return int16(val), nil
	case int64:
		return int16(val), nil
	case int8:
		return int16(val), nil
	case int16:
		return val, nil
	case int32:
		return int16(val), nil
	case uint:
		return int16(val), nil
	case uint8:
		return int16(val), nil
	case uint16:
		return int16(val), nil
	case uint32:
		return int16(val), nil
	case uint64:
		return int16(val), nil
	case float32:
		return int16(val), nil
	case float64:
		return int16(val), nil
	case complex64:
		return int16(real(val)), nil
	case complex128:
		return int16(real(val)), nil
	case time.Duration:
		return int16(val), nil
	case time.Time: //return the unix value
		return int16(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
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
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int16",
				"error":    err.Error(),
			})
		}
		return int16(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return int8(val), nil
	case int64:
		return int8(val), nil
	case int8:
		return val, nil
	case int16:
		return int8(val), nil
	case int32:
		return int8(val), nil
	case uint:
		return int8(val), nil
	case uint8:
		return int8(val), nil
	case uint16:
		return int8(val), nil
	case uint32:
		return int8(val), nil
	case uint64:
		return int8(val), nil
	case float32:
		return int8(val), nil
	case float64:
		return int8(val), nil
	case complex64:
		return int8(real(val)), nil
	case complex128:
		return int8(real(val)), nil
	case time.Duration:
		return int8(val), nil
	case time.Time: //return the unix value
		return int8(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "int8",
				"error":    err.Error(),
			})
		}
		return int8(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "int8",
				"error":    err.Error(),
			})
		}
		return int8(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return uint64(val), nil
	case int8:
		return uint64(val), nil
	case int16:
		return uint64(val), nil
	case int32:
		return uint64(val), nil
	case int64:
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
		return uint64(val), nil
	case float64:
		return uint64(val), nil
	case complex64:
		return uint64(real(val)), nil
	case complex128:
		return uint64(real(val)), nil
	case time.Duration:
		return uint64(val), nil
	case time.Time: //return the unix value
		return uint64(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "uint64",
				"error":    err.Error(),
			})
		}
		return uint64(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint64",
				"error":    err.Error(),
			})
		}
		return uint64(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return uint32(val), nil
	case int8:
		return uint32(val), nil
	case int16:
		return uint32(val), nil
	case int32:
		return uint32(val), nil
	case int64:
		return uint32(val), nil
	case uint:
		return uint32(val), nil
	case uint8:
		return uint32(val), nil
	case uint16:
		return uint32(val), nil
	case uint64:
		return uint32(val), nil
	case float32:
		return uint32(val), nil
	case float64:
		return uint32(val), nil
	case complex64:
		return uint32(real(val)), nil
	case complex128:
		return uint32(real(val)), nil
	case time.Duration:
		return uint32(val), nil
	case time.Time: //return the unix value
		return uint32(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "uint32",
				"error":    err.Error(),
			})
		}
		return uint32(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint32",
				"error":    err.Error(),
			})
		}
		return uint32(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return uint16(val), nil
	case int8:
		return uint16(val), nil
	case int16:
		return uint16(val), nil
	case int32:
		return uint16(val), nil
	case int64:
		return uint16(val), nil
	case uint:
		return uint16(val), nil
	case uint8:
		return uint16(val), nil
	case uint32:
		return uint16(val), nil
	case uint64:
		return uint16(val), nil
	case float32:
		return uint16(val), nil
	case float64:
		return uint16(val), nil
	case complex64:
		return uint16(real(val)), nil
	case complex128:
		return uint16(real(val)), nil
	case time.Duration:
		return uint16(val), nil
	case time.Time: //return the unix value
		return uint16(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "uint16",
				"error":    err.Error(),
			})
		}
		return uint16(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint16",
				"error":    err.Error(),
			})
		}
		return uint16(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return uint8(val), nil
	case int8:
		return uint8(val), nil
	case int16:
		return uint8(val), nil
	case int32:
		return uint8(val), nil
	case int64:
		return uint8(val), nil
	case uint:
		return uint8(val), nil
	case uint32:
		return uint8(val), nil
	case uint16:
		return uint8(val), nil
	case uint64:
		return uint8(val), nil
	case float32:
		return uint8(val), nil
	case float64:
		return uint8(val), nil
	case complex64:
		return uint8(real(val)), nil
	case complex128:
		return uint8(real(val)), nil
	case time.Duration:
		return uint8(val), nil
	case time.Time: //return the unix value
		return uint8(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "uint8",
				"error":    err.Error(),
			})
		}
		return uint8(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "uint8",
				"error":    err.Error(),
			})
		}
		return uint8(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return float64(val), nil
	case float32:
		return float64(val), nil
	case complex64:
		return float64(real(val)), nil
	case complex128:
		return real(val), nil
	case time.Duration:
		return float64(val), nil
	case time.Time: //return the unix value
		return float64(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatValue, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "float64",
				"error":    err.Error(),
			})
		}
		return floatValue, nil
	case string:
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "float64",
				"error":    err.Error(),
			})
		}
		return floatValue, nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return float32(val), nil
	case float64:
		return float32(val), nil
	case complex64:
		return real(val), nil
	case complex128:
		return float32(real(val)), nil
	case time.Duration:
		return float32(val), nil
	case time.Time: //return the unix value
		return float32(val.Unix()), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "float32",
				"error":    err.Error(),
			})
		}
		return float32(floatVal), nil
	case string:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "float32",
				"error":    err.Error(),
			})
		}
		return float32(floatVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
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
		return complex(float32(val), float32(0)), nil //there may be some conversion loss here
	case uint:
		return complex(float32(val), float32(0)), nil
	case uint8:
		return complex(float32(val), float32(0)), nil
	case uint16:
		return complex(float32(val), float32(0)), nil
	case uint32:
		return complex(float32(val), float32(0)), nil
	case uint64:
		return complex(float32(val), float32(0)), nil //there may be some conversion loss here
	case float32:
		return complex(val, float32(0)), nil
	case float64:
		return complex(float32(val), float32(0)), nil //there may be some conversion loss here
	case complex128:
		return complex64(val), nil //there may be some conversion loss here
	case time.Duration:
		return complex(float32(val), float32(0)), nil
	case time.Time: //return the unix value
		return complex(float32(val.Unix()), float32(0)), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		complexVal, err := strconv.ParseComplex(string(val), 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "complex64",
				"error":    err.Error(),
			})
		}
		return complex64(complexVal), nil
	case string:
		complexVal, err := strconv.ParseComplex(val, 64)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "complex64",
				"error":    err.Error(),
			})
		}
		return complex64(complexVal), nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
		return complex(float64(val), float64(0)), nil
	case float32:
		return complex(float64(val), float64(0)), nil
	case float64:
		return complex(val, float64(0)), nil
	case complex64:
		return complex128(val), nil //there may be some conversion loss here
	case time.Duration:
		return complex(float64(val), float64(0)), nil
	case time.Time: //return the unix value
		return complex(float64(val.Unix()), float64(0)), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		complexVal, err := strconv.ParseComplex(string(val), 128)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "[]byte",
				"dst_type": "complex128",
				"error":    err.Error(),
			})
		}
		return complexVal, nil
	case string:
		complexVal, err := strconv.ParseComplex(val, 128)
		if err != nil {
			return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      val,
				"src_type": "string",
				"dst_type": "complex128",
				"error":    err.Error(),
			})
		}
		return complexVal, nil
	default:
		return 0, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
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
	case time.Time: //return the unix value
		return decimal.NewFromInt(int64(val.Unix())), nil
	case bool:
		if val {
			return decimal.NewFromInt(1), nil
		}
		return decimal.NewFromInt(0), nil
	case []byte:
		floatval, _ := strconv.ParseFloat(string(val), 64)
		return decimal.NewFromFloat(floatval), nil
	case string:
		floatval, _ := strconv.ParseFloat(val, 64)
		return decimal.NewFromFloat(floatval), nil
	default:
		return decimal.NewFromInt(0), zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
			"dst_type": "decimal",
		})
	}
}

// String - tries to convert any to string
func String(src any) (dst string, err zerror.Error) {
	if src == nil {
		return "", nil
	}
	switch val := src.(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case []any: //will try to convert to []byte
		sliceByte, err := SliceByte(val)
		if err != nil {
			return "", err
		}
		return string(sliceByte), nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil
	case int:
		return strconv.FormatInt(int64(val), 10), nil
	case int8:
		return strconv.FormatInt(int64(val), 10), nil
	case int16:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case time.Duration:
		return strconv.FormatUint(uint64(val), 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'g', -1, 64), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case complex64:
		return strconv.FormatComplex(complex128(val), 'g', -1, 64), nil
	case complex128:
		return strconv.FormatComplex(val, 'g', -1, 128), nil
	case time.Time: //return the standard ISO STZ format
		return val.Format(TimeFormatISOSTZ), nil
	case Stringable:
		return val.String(), nil
	case ToStringable:
		return val.ToString(), nil
	default:
		return "", zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
			"dst_type": "string",
		})

	}
}

// Bool - tries to convert any to bool
func Bool(src any) (dst bool, err zerror.Error) {
	if src == nil {
		return false, nil
	}
	switch val := src.(type) {
	case bool:
		return val, nil
	case uint64, uint32, uint16, uint8, uint, int, int8, int16, int32, int64, float64, float32, time.Duration:
		valF64, _ := Float64(val)
		return valF64 != float64(0), nil
	case complex64:
		return real(val) != float32(0) || imag(val) != float32(0), nil
	case complex128:
		return real(val) != float64(0) || imag(val) != float64(0), nil
	case []byte:
		stringVal := strings.ToLower(string(val))
		if stringVal == "" || stringVal == "false" { //element is empty or "false"
			return false, nil
		}
		return true, nil
	case string:
		stringVal := strings.ToLower(val)
		if stringVal == "" || stringVal == "false" { //element is empty or "false"
			return false, nil
		}
		return true, nil
	default:
		return false, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": "default",
			"dst_type": "bool",
		})
	}
}

// MapStringAny - tries to convert any to map[string]any
func MapStringAny(src any) (dst map[string]any, err zerror.Error) {
	result := map[string]any{}
	if src == nil {
		return result, nil
	}
	//fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.(map[string]any); ok {
		return srcVal, nil
	}

	//converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Map:
		for _, mapKey := range elemValue.MapKeys() {
			key, err := String(mapKey.Interface())
			if err != nil {
				return result, err
			}
			result[key] = elemValue.MapIndex(mapKey).Interface()
		}
		return result, nil
		//case reflect.Struct:
		//  err = ScanToMap(&result, Config{ParseTags: true, ParseOnlyTags: true, ParseTagList: []string{"db", "json"}, MapIgnoreObjectNames: true}, elemValue.Interface())
		//  if err != nil {
		//    return result, err
		//  }
		//  return result, nil
		//case reflect.Ptr:
		//  unpackedVal := Unpack(src, false) // removing pointer
		//  return MapStringAny(unpackedVal)
	}
	return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": elemKind.String(),
		"dst_type": "map[string]any",
	})
}

// SliceAny - tries to convert any to []any
func SliceAny(src any) (dst []any, err zerror.Error) {
	result := []any{}
	if src == nil {
		return result, nil
	}
	//fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]any); ok {
		return srcVal, nil
	}

	//converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Struct, reflect.Map, reflect.Func, reflect.Chan: //types not supported
	case reflect.Slice, reflect.Array:
		for i := 0; i < elemValue.Len(); i++ {
			result = append(result, elemValue.Index(i).Interface())
		}
		return result, nil
	default: //simple type, we convert to string and add it as a slice element
		result = append(result, elemValue.Interface())
		return result, nil
	}
	return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": elemKind.String(),
		"dst_type": "[]any",
	})
}

// SliceByte - tries to convert any to []byte
func SliceByte(src any) (dst []byte, err zerror.Error) {
	result := []byte{}
	if src == nil {
		return result, nil
	}
	//fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]byte); ok {
		return srcVal, nil
	}

	//converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Slice, reflect.Array:
		for i := 0; i < elemValue.Len(); i++ {
			resByte, err := Uint8(elemValue.Index(i).Interface())
			if err != nil {
				return result, err
			}
			result = append(result, resByte)
		}
		return result, nil
	case reflect.String, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		stringVal, _ := String(src)
		return []byte(stringVal), nil
	}
	return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": elemKind.String(),
		"dst_type": "[]byte",
	})

}

// SliceString - tries to convert any to []string
func SliceString(src any) (dst []string, err zerror.Error) {
	result := []string{}
	if src == nil {
		return result, nil
	}
	//fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]string); ok {
		return srcVal, nil
	}

	//converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Struct, reflect.Map, reflect.Func, reflect.Chan:
	case reflect.Slice, reflect.Array:
		for i := 0; i < elemValue.Len(); i++ {
			resString, err := String(elemValue.Index(i).Interface())
			if err != nil {
				return result, err
			}
			result = append(result, resString)
		}
		return result, nil
	default: //simple type, we convert to string and add it as a slice element
		resString, err := String(elemValue.Interface())
		if err != nil {
			return result, err
		}
		result = append(result, resString)
		return result, nil
	}

	return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": elemKind.String(),
		"dst_type": "[]string",
	})

}

// SliceInt - tries to convert any to []int
func SliceInt(src any) (dst []int, err zerror.Error) {
	result := []int{}
	if src == nil {
		return result, nil
	}
	//fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]int); ok {
		return srcVal, nil
	}

	//converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	if elemKind == reflect.Slice || elemKind == reflect.Array {
		for i := 0; i < elemValue.Len(); i++ {
			resInt, err := Int(elemValue.Index(i).Interface())
			if err != nil {
				return result, err
			}
			result = append(result, resInt)
		}
		return result, nil
	}
	return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": elemKind.String(),
		"dst_type": "[]int",
	})
}

// SliceMapStringAny - tries to convert any to []map[string]any
func SliceMapStringAny(src any) (dst []map[string]any, err zerror.Error) {
	result := []map[string]any{}
	if src == nil {
		return result, nil
	}
	//fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]map[string]any); ok {
		return srcVal, nil
	}

	//converting other types
	srcValue := reflect.ValueOf(src)
	srcKind := reflect.TypeOf(src).Kind()

	if srcKind == reflect.Slice || srcKind == reflect.Array {
		for i := 0; i < srcValue.Len(); i++ {
			elem := srcValue.Index(i).Interface()
			elemKind := reflect.TypeOf(elem).Kind()
			switch elemKind {
			case reflect.Map, reflect.Struct:
				res, err := MapStringAny(elem)
				if err != nil {
					return result, err
				}
				result = append(result, res)
			}
		}
		return result, nil
	}
	return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": srcKind.String(),
		"dst_type": "[]map[string]any",
	})

}

// Time - tries to convert any to time.Time
// params:
//
//	src
//	   1 - string     - will try parse the string and returns appropriate value (see code for supported RFC and ISO formats)
//	   1,2 - number   - (int, uint, float types) - will be converted to integers and assumes unix time and/or unixnano time
//	   7 numbers      - time.Date(year, time.Month(month), day, hour, min, sec, nsec)
//	   time.Time      - returns value as is
//	   other          - will return an error
func Time(args ...any) (dst time.Time, err zerror.Error) {
	result := time.Time{}
	if len(args) == 0 {
		return result, nil
	} else if len(args) == 1 {
		if args[0] == nil {
			return result, nil
		}
		if timeVal, ok := args[0].(time.Time); ok {
			return timeVal, nil
		}
		//checking other types
		elemKind := reflect.TypeOf(args[0]).Kind()
		switch elemKind {
		case reflect.Slice: // we spread the slice
			sliceParam, zer := SliceAny(args[0])
			if zer != nil {
				err = zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
					"src":      args[0],
					"src_type": elemKind.String(),
					"dst_type": "time",
				})
				err.Add(zer.GetList())
				return result, err
			}
			return Time(sliceParam...)
		case reflect.String: //date is sent as a string so we will try to parse it
			timeStr, zer := String(args[0])
			if zer != nil {
				err = zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
					"src":      args[0],
					"src_type": elemKind.String(),
					"dst_type": "time",
				})
				err.Add(zer.GetList())
			}
			//we try to match one of the supported time formats
			result, er := time.Parse(time.RFC822, timeStr)
			if er != nil {
				result, er = time.Parse(time.RFC850, timeStr)
				if er != nil {
					result, er = time.Parse(time.RFC1123, timeStr)
					if er != nil {
						result, er = time.Parse(time.RFC822Z, timeStr)
						if er != nil {
							result, er = time.Parse(time.RFC3339, timeStr)
							if er != nil {
								result, er = time.Parse(time.RFC1123Z, timeStr)
								if er != nil {
									result, er = time.Parse(time.RFC3339Nano, timeStr)
									if er != nil {
										result, er = time.Parse(TimeFormatISOTZ, timeStr) //ISO datetime format with Z timezone
										if er != nil {
											result, er = time.Parse(TimeFormatISOSTZ, timeStr) //ISO datetime format with spacer timezone
											if er != nil {
												result, er = time.Parse(TimeFormatISO, timeStr) //ISO datetime format with no timezone
												if er != nil {
													result, er = time.Parse(TimeFormatISODate, timeStr) //ISO date format
													if er != nil {
														return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
															"src":      timeStr,
															"src_type": elemKind.String(),
															"dst_type": "time",
															"error":    er.Error(),
														})
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return result, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			unixTime, zer := Int64(args[0])
			if zer != nil {
				err = zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
					"src":      args[0],
					"src_type": elemKind.String(),
					"dst_type": "time",
				})
				err.Add(zer.GetList())
				return result, err
			}
			result = time.Unix(unixTime, 0)
		case reflect.Float32, reflect.Float64: // time is in float format, int part is unixTime, decimals are unixNano, there will be some nanosecond errors because of some floating point operations
			floatTime, zer := Float64(args[0])
			if zer != nil {
				err = zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
					"src":      args[0],
					"src_type": elemKind.String(),
					"dst_type": "time",
				})
				err.Add(zer.GetList())
				return result, err
			}
			unixTime, zer := Int64(args[0])
			if zer != nil {
				err = zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
					"src":      args[0],
					"src_type": elemKind.String(),
					"dst_type": "time",
				})
				err.Add(zer.GetList())

				return result, err
			}
			unixNano := int64((floatTime - float64(unixTime)) * 1e9)
			result = time.Unix(unixTime, unixNano)
		default:
			err = zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
				"src":      args[0],
				"src_type": elemKind.String(),
				"dst_type": "time",
			})
			return result, err
		}
	} else if len(args) == 2 { //assuming 2 integers for unixTime and unixNano
		unixTime, err := Int64(args[0])
		if err != nil {
			return result, err
		}
		unixNano, err := Int64(args[1])
		if err != nil {
			return result, err
		}
		result = time.Unix(unixTime, unixNano)
	} else if len(args) == 7 { //assuming 7 integers, no location
		year, err := Int(args[0])
		if err != nil {
			return result, err
		}
		month, err := Int(args[1])
		if err != nil {
			return result, err
		}
		day, err := Int(args[2])
		if err != nil {
			return result, err
		}
		hour, err := Int(args[3])
		if err != nil {
			return result, err
		}
		min, err := Int(args[4])
		if err != nil {
			return result, err
		}
		sec, err := Int(args[5])
		if err != nil {
			return result, err
		}
		nsec, err := Int(args[6])
		if err != nil {
			return result, err
		}
		var location *time.Location //will remain nil
		result = time.Date(year, time.Month(month), day, hour, min, sec, nsec, location)
	} else {
		return result, zerror.New(ErrorZGENConvertorTypeNotSupported, map[string]any{
			"src":      args,
			"src_type": "default",
			"dst_type": "time",
		})
	}

	return result, nil
}
