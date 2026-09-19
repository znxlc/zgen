package zgen

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/znxlc/zerror"

	"github.com/shopspring/decimal"
)

// String converts any supported type to its string representation.
//
// This function handles a wide variety of input types and returns appropriate
// string representations. For nil input, it returns an empty string with no error.
//
// Supported types and their conversions:
//   - string: returned as-is
//   - []byte: converted using string(val)
//   - bool: "true" or "false"
//   - int, int8, int16, int32, int64: decimal string representation
//   - uint, uint8, uint16, uint32, uint64: decimal string representation
//   - float32, float64: formatted as decimal without unnecessary trailing zeros
//   - complex64, complex128: formatted using Go's complex number format
//   - time.Duration: using the standard Duration.String() format
//   - time.Time: formatted using TimeFormatISOSTZ (ISO 8601 with timezone)
//   - decimal.Decimal: using the Decimal.String() method
//   - Stringable interface: using the String() method
//   - ToStringable interface: using the ToString() method
//   - Slice types: converted by joining string representations of elements
//
// For unsupported types, returns an error with ErrorConvertorTypeNotSupported.
//
// Parameters:
//   - src: the value to convert to string
//
// Returns:
//   - string: the string representation of src
//   - error: zerror.Error if conversion fails, nil otherwise
func String(src any) (dst string, err zerror.Error) {
	if src == nil {
		return "", nil
	}
	switch val := src.(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case []string:
		return strings.Join(val, ""), nil
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
		return val.String(), nil
		// return strconv.FormatInt(int64(val), 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case complex64:
		return strconv.FormatComplex(complex128(val), 'g', -1, 64), nil
	case complex128:
		return strconv.FormatComplex(val, 'g', -1, 128), nil
	case time.Time: // return the standard ISO STZ format
		return val.Format(TimeFormatISOSTZ), nil
	case decimal.Decimal:
		return val.String(), nil
	case Stringable:
		return val.String(), nil
	case ToStringable:
		return val.ToString(), nil
	default: // trying to process other types using reflection
		reflectVal := reflect.ValueOf(val)
		if reflectVal.Kind() == reflect.Slice {
			sliceString, err := SliceString(val)
			if err != nil {
				return "", err
			}
			return strings.Join(sliceString, ""), nil
		}
		return "", zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
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
		if stringVal == "" || stringVal == "false" { // element is empty or "false"
			return false, nil
		}
		return true, nil
	case string:
		tst, er := strconv.ParseBool(val)
		if er != nil {
			return false, nil
		}
		return tst, nil
	default:
		return false, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      val,
			"src_type": reflect.TypeOf(val).String(),
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
	// fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.(map[string]any); ok {
		return srcVal, nil
	}

	// converting other map types
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
	case reflect.Struct:
		err = ToMap(&result, DefaultParserConfig, elemValue.Interface())
		if err != nil {
			return result, err
		}
		return result, nil
	case reflect.Ptr:
		unpackedVal := UnpackBaseElement(src, false) // removing pointer
		return MapStringAny(unpackedVal)
	}
	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": elemKind.String(),
		"dst_type": "map[string]any",
	})
}

// MapStringSliceByte - tries to convert any to map[string][]byte
func MapStringSliceByte(src any) (dst map[string][]byte, err zerror.Error) {
	result := map[string][]byte{}
	if src == nil {
		return result, nil
	}
	// fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.(map[string][]byte); ok {
		return srcVal, nil
	}

	// convert to map[string]any
	mapStringAny, err := MapStringAny(src)
	if err != nil {
		return result, err
	}

	for key, value := range mapStringAny {
		sbValue, er := SliceByte(value)
		if er != nil {
			err = zerror.New(er)
			err.Add(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      src,
				"src_type": reflect.ValueOf(src).Kind().String(),
				"dst_type": "map[string][]byte",
			})
			return result, err
		}

		result[key] = sbValue
	}

	return result, nil

}

// SliceAny - tries to convert any to []any
func SliceAny(src any) (dst []any, err zerror.Error) {
	result := []any{}
	if src == nil {
		return result, nil
	}
	// fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]any); ok {
		return srcVal, nil
	}

	// converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Map, reflect.Struct, reflect.Chan, reflect.Func, reflect.Invalid:
	case reflect.Slice, reflect.Array:
		for i := 0; i < elemValue.Len(); i++ {
			result = append(result, elemValue.Index(i).Interface())
		}
		return result, nil
	default: // simple type, we convert to string and add it as a slice element
		result = append(result, elemValue.Interface())
		return result, nil
	}
	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
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
	// fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]byte); ok {
		return srcVal, nil
	}

	// converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Map, reflect.Struct, reflect.Chan, reflect.Func, reflect.Invalid:
	case reflect.Slice, reflect.Array:
		for i := 0; i < elemValue.Len(); i++ {
			resByte, err := Uint8(elemValue.Index(i).Interface())
			if err != nil {
				return result, err
			}
			result = append(result, resByte)
		}
		return result, nil
	case reflect.String, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Bool:
		stringVal, _ := String(src)
		return []byte(stringVal), nil
	}
	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
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
	// fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]string); ok {
		return srcVal, nil
	}

	// converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Map, reflect.Struct, reflect.Chan, reflect.Func, reflect.Invalid:
	case reflect.Slice, reflect.Array:
		for i := 0; i < elemValue.Len(); i++ {
			resString, err := String(elemValue.Index(i).Interface())
			if err != nil {
				return result, err
			}
			result = append(result, resString)
		}
		return result, nil
	default: // simple type, we convert to string and add it as a slice element
		resString, err := String(elemValue.Interface())
		if err != nil {
			return result, err
		}
		result = append(result, resString)
		return result, nil
	}

	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
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
	// fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]int); ok {
		return srcVal, nil
	}

	// converting other map types
	elemValue := reflect.ValueOf(src)
	elemKind := reflect.TypeOf(src).Kind()
	switch elemKind {
	case reflect.Map, reflect.Struct, reflect.Chan, reflect.Func, reflect.Invalid:
	case reflect.Slice, reflect.Array:
		for i := 0; i < elemValue.Len(); i++ {
			resInt, err := Int(elemValue.Index(i).Interface())
			if err != nil {
				return result, err
			}
			result = append(result, resInt)
		}
		return result, nil
	default:
		resInt, err := Int(elemValue.Interface())
		if err != nil {
			return result, err
		}
		result = append(result, resInt)
		return result, nil
	}
	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
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
	// fast check to see the src type is same as dst to avoid fancy reflect operations
	if srcVal, ok := src.([]map[string]any); ok {
		return srcVal, nil
	}

	// converting other types
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
			default:
				return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
					"src":      elem,
					"src_type": elemKind.String(),
					"dst_type": "map[string]any",
				})
			}
		}
		return result, nil
	} else if srcKind == reflect.Map || srcKind == reflect.Struct {
		res, err := MapStringAny(src)
		if err != nil {
			return result, err
		}
		result = append(result, res)
		return result, nil
	}
	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
		"src":      src,
		"src_type": srcKind.String(),
		"dst_type": "[]map[string]any",
	})

}
