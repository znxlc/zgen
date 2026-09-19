package zgen

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/shopspring/decimal"
)

func TestUnit_Uint(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    uint
		hasError  bool
		errorCode string
	}{
		// Signed integer tests
		{name: "int positive", input: int(5), expect: 5},
		{name: "int negative", input: int(-5), expect: 0, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int zero", input: int(0), expect: 0},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int64 max", input: int64(math.MaxInt64), expect: math.MaxInt64},
		{name: "int negative", input: int(-5), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint zero", input: uint(0), expect: 0},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},

		// Float tests
		{name: "float32", input: float32(5.5), expect: 5},
		{name: "float64", input: 5.5, expect: 5},
		{name: "float64 large", input: 1.23e6, expect: 1230000},
		{name: "float64 negative", input: -5.5, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 overflow", input: 1e100, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex number tests
		{name: "complex64 real part", input: complex64(5 + 12i), expect: 5},
		{name: "complex64 zero imaginary", input: complex64(5 + 0i), expect: 5},
		{name: "complex128 real part", input: 5 + 12i, expect: 5},
		{name: "complex128 zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex128 negative real", input: -5 + 12i, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// String and byte slice tests
		{name: "string decimal", input: "42.1", expect: 42},
		{name: "string negative decimal", input: "-42.1", expect: 0, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte decimal", input: []byte("42"), expect: 42},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Bool tests
		{name: "bool true", input: true, expect: 1},
		{name: "bool false", input: false, expect: 0},

		// Nil test
		{name: "nil value", input: nil, expect: 0},

		// Unsupported types
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []int{1, 2, 3}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Uint(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode))
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Int(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    int
		hasError  bool
		errorCode string
	}{
		// Signed integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int64 max", input: int64(math.MaxInt64), expect: math.MaxInt},
		{name: "negative int", input: int(-5), expect: -5},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint64 max", input: uint64(math.MaxInt64), expect: math.MaxInt},
		{name: "uint64 overflow", input: uint64(math.MaxInt64 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 negative", input: -5.0, expect: -5},
		{name: "float64 overflow", input: math.MaxFloat64, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex negative real", input: -5 + 12i, expect: -5},
		{name: "complex overflow", input: complex128(math.MaxInt64 + 1 + 10i), expect: math.MinInt},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5 * time.Second), expect: 5000000000},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, expect: -1000000000},

		// String and byte slice tests
		{name: "string decimal", input: "42", expect: 42},
		{name: "string negative decimal", input: "-42", expect: -42},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte decimal", input: []byte("42"), expect: 42},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Edge cases
		{name: "max int64", input: int64(math.MaxInt64), expect: math.MaxInt64},
		{name: "min int64", input: int64(math.MinInt64), expect: math.MinInt64},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []int{1, 2, 3}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Int(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Int8(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    int8
		hasError  bool
		errorCode string
	}{
		// Signed integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int max", input: math.MaxInt8, expect: math.MaxInt8},
		{name: "int min", input: math.MinInt8, expect: math.MinInt8},
		{name: "int overflow", input: math.MaxInt8 + 1, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int underflow", input: math.MinInt8 - 1, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint max", input: uint8(math.MaxInt8), expect: math.MaxInt8},
		{name: "uint overflow", input: uint8(math.MaxInt8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 negative", input: -5.0, expect: -5},
		{name: "float64 overflow", input: math.MaxInt8 + 1.0, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 underflow", input: math.MinInt8 - 1.0, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex negative real", input: -5 + 12i, expect: -5},
		{name: "complex overflow", input: complex128(math.MaxInt8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: 5, expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -5, expect: -5},
		{name: "time.Duration overflow", input: time.Duration(math.MaxInt8+1) * time.Second, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// String and byte slice tests
		{name: "string decimal", input: "42", expect: 42},
		{name: "string negative decimal", input: "-42", expect: -42},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "string overflow", input: "128", hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "string underflow", input: "-129", hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "[]byte decimal", input: []byte("42"), expect: 42},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Edge cases
		{name: "max int8", input: int8(math.MaxInt8), expect: math.MaxInt8},
		{name: "min int8", input: int8(math.MinInt8), expect: math.MinInt8},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []int{1, 2, 3}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Int8(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
	// All test cases are now in the table above
}

func TestUnit_Int16(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    int16
		hasError  bool
		errorCode string
	}{
		// Signed integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int max", input: math.MaxInt16, expect: math.MaxInt16},
		{name: "int min", input: math.MinInt16, expect: math.MinInt16},
		{name: "int overflow", input: math.MaxInt16 + 1, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int underflow", input: math.MinInt16 - 1, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint max", input: uint16(math.MaxInt16), expect: math.MaxInt16},
		{name: "uint overflow", input: uint16(math.MaxInt16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 negative", input: -5.0, expect: -5},
		{name: "float64 overflow", input: math.MaxInt16 + 1.0, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 underflow", input: math.MinInt16 - 1.0, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex negative real", input: -5 + 12i, expect: -5},
		{name: "complex overflow", input: complex128(math.MaxInt16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: 5 * time.Microsecond, expect: 5000},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -5 * time.Microsecond, expect: -5000},
		{name: "time.Duration overflow", input: time.Duration(math.MaxInt16+1) * time.Millisecond, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// String and byte slice tests
		{name: "string decimal", input: "42", expect: 42},
		{name: "string negative decimal", input: "-42", expect: -42},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "string overflow", input: "132768", hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "string underflow", input: "-132769", hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "[]byte decimal", input: []byte("42"), expect: 42},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Edge cases
		{name: "max int16", input: int16(math.MaxInt16), expect: math.MaxInt16},
		{name: "min int16", input: int16(math.MinInt16), expect: math.MinInt16},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []int{1, 2, 3}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Int16(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
	// All test cases are now in the table above
}

func TestUnit_Int32(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    int32
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int64 max", input: int64(math.MaxInt32), expect: math.MaxInt32},
		{name: "int64 min", input: int64(math.MinInt32), expect: math.MinInt32},
		{name: "int64 overflow", input: int64(math.MaxInt32 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int64 underflow", input: int64(math.MinInt32 - 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint64 max", input: uint64(math.MaxInt32), expect: math.MaxInt32},
		{name: "uint64 overflow", input: uint64(math.MaxInt32 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 negative", input: -5.0, expect: -5},
		{name: "float64 overflow", input: math.MaxFloat64, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 underflow", input: math.SmallestNonzeroFloat64, expect: 0},
		{name: "float64 NaN", input: math.NaN(), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex negative real", input: -5 + 12i, expect: -5},
		{name: "complex overflow", input: complex128(math.MaxInt32 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex underflow", input: complex128(math.MinInt32 - 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, expect: -1000000000},

		// String and byte slice tests
		{name: "string decimal", input: "42", expect: 42},
		{name: "string negative decimal", input: "-42", expect: -42},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "string overflow", input: "2147483648", hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "string underflow", input: "-2147483649", hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "[]byte decimal", input: []byte("42"), expect: 42},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Int32(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Int64(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    int64
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int64 max", input: int64(math.MaxInt64), expect: math.MaxInt64},
		{name: "int64 min", input: int64(math.MinInt64), expect: math.MinInt64},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint64 max", input: uint64(math.MaxInt64), expect: math.MaxInt64},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 negative", input: -5.0, expect: -5},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex negative real", input: -5 + 12i, expect: -5},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, expect: -1000000000},

		// String and byte slice tests
		{name: "string decimal", input: "42", expect: 42},
		{name: "string negative decimal", input: "-42", expect: -42},
		{name: "string max int64", input: "9223372036854775807", expect: math.MaxInt64},
		{name: "string min int64", input: "-9223372036854775808", expect: math.MinInt64},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "string negative decimal", input: "-10.12", expect: -10},
		{name: "[]byte decimal", input: []byte("42"), expect: 42},
		{name: "[]byte with decimal", input: []byte("10.12"), expect: 10},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Int64(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Uint8(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    uint8
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -1, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int8 negative", input: int8(-1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int16 overflow", input: int16(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int32 overflow", input: int32(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int64 overflow", input: int64(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint16 overflow", input: uint16(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "uint32 overflow", input: uint32(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "uint64 overflow", input: uint64(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float32 overflow", input: float32(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 overflow", input: float64(math.MaxUint8 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 NaN", input: math.NaN(), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex64 overflow", input: complex64(math.MaxUint8 + 1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex64 negative", input: complex64(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex128 overflow", input: complex128(math.MaxUint8 + 1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex128 negative", input: complex128(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Uint8(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Uint16(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    uint16
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -1, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int16 negative", input: int16(-1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int32 overflow", input: int32(math.MaxUint16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int64 overflow", input: int64(math.MaxUint16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint32 overflow", input: uint32(math.MaxUint16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "uint64 overflow", input: uint64(math.MaxUint16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float32 overflow", input: float32(math.MaxUint16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 overflow", input: float64(math.MaxUint16 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 NaN", input: math.NaN(), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex64 overflow", input: complex64(math.MaxUint16 + 1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex64 negative", input: complex64(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex128 overflow", input: complex128(math.MaxUint16 + 1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex128 negative", input: complex128(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Uint16(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Uint32(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    uint32
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -1, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int32 negative", input: int32(-1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int64 overflow", input: int64(math.MaxUint32 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint64 overflow", input: uint64(math.MaxUint32 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 overflow", input: float64(math.MaxUint32 + 1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 NaN", input: math.NaN(), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex64 negative", input: complex64(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex128 overflow", input: complex128(math.MaxUint32 + 1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex128 negative", input: complex128(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Uint32(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Uint64(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    uint64
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -1, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "int64 negative", input: int64(-1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 negative", input: -5.0, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 NaN", input: math.NaN(), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex64 negative", input: complex64(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "complex128 negative", input: complex128(-1 + 0i), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "string negative", input: "-5", hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Uint64(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expect, res, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Float32(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    float32
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -5, expect: -5},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 scientific positive exponent", input: 1.5134e+02, expect: 151.34},
		{name: "float64 scientific negative exponent", input: 15134e-02, expect: 151.34},
		{name: "float64 max", input: math.MaxFloat64, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 min", input: -math.MaxFloat64, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 NaN", input: math.NaN(), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},
		{name: "complex overflow", input: complex128(math.MaxFloat32 * 2), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, expect: -1000000000},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "string scientific", input: "15134e-02", expect: 151.34},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Float32(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Use InDelta for float comparison to handle floating point imprecision
			assert.InDelta(t, float64(tt.expect), float64(res), 0.0001, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Float64(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    float64
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -5, expect: -5},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},
		{name: "uint64 max", input: uint64(math.MaxUint64), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// Float tests
		{name: "float32", input: float32(5.0), expect: 5},
		{name: "float64", input: 5.0, expect: 5},
		{name: "float64 scientific positive exponent", input: 1.5134e+02, expect: 151.34},
		{name: "float64 scientific negative exponent", input: 15134e-02, expect: 151.34},
		{name: "float64 max", input: math.MaxFloat64, expect: math.MaxFloat64},
		{name: "float64 smallest non-zero", input: math.SmallestNonzeroFloat64, expect: math.SmallestNonzeroFloat64},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5},
		{name: "complex128", input: 5 + 12i, expect: 5},
		{name: "complex zero imaginary", input: 0 + 5i, expect: 0},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, expect: -1000000000},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "string scientific", input: "15134e-02", expect: 151.34},
		{name: "string with decimal point", input: "10.00000000", expect: 10},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Float64(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Use InDelta for float comparison to handle floating point imprecision
			assert.InDelta(t, tt.expect, res, 0.0001, "Unexpected result for input: %v", tt.input)
		})
	}
}

func TestUnit_Decimal(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		setup     func() decimal.Decimal
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "int8", input: int8(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "int16", input: int16(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "int32", input: int32(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "int64", input: int64(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "int negative", input: -5, setup: func() decimal.Decimal { return decimal.NewFromInt(-5) }},

		// Unsigned integer tests
		{name: "uint", input: uint(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "uint8", input: uint8(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "uint16", input: uint16(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "uint32", input: uint32(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "uint64", input: uint64(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},

		// Complex tests
		{name: "complex64 real part", input: complex64(5 + 12i), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "complex128 real part", input: 5 + 12i, setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "complex zero real part", input: 0 + 5i, setup: func() decimal.Decimal { return decimal.NewFromInt(0) }},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), setup: func() decimal.Decimal { return decimal.NewFromInt(5) }},
		{name: "time.Duration zero", input: time.Duration(0), setup: func() decimal.Decimal { return decimal.NewFromInt(0) }},
		{name: "time.Duration negative", input: -time.Second, setup: func() decimal.Decimal { return decimal.NewFromInt(-1000000000) }},

		// Float tests
		{name: "float32", input: float32(5.5), setup: func() decimal.Decimal { return decimal.NewFromFloat(5.5) }},
		{name: "float64", input: 5.5, setup: func() decimal.Decimal { return decimal.NewFromFloat(5.5) }},
		{name: "float64 scientific positive", input: 1.5134e+02, setup: func() decimal.Decimal { return decimal.NewFromFloat(1.5134e+02) }},
		{name: "float64 scientific negative", input: 1.5134e-02, setup: func() decimal.Decimal { return decimal.NewFromFloat(1.5134e-02) }},
		{name: "float64 max", input: math.MaxFloat64, setup: func() decimal.Decimal { d, _ := decimal.NewFromString("1.7976931348623157e+308"); return d }},
		{name: "float64 smallest non-zero", input: math.SmallestNonzeroFloat64, setup: func() decimal.Decimal { d, _ := decimal.NewFromString("5e-324"); return d }},

		// String and byte slice tests
		{name: "string decimal", input: "5.5", setup: func() decimal.Decimal { d, _ := decimal.NewFromString("5.5"); return d }},
		{name: "string scientific", input: "15134e-02", setup: func() decimal.Decimal { d, _ := decimal.NewFromString("15134e-02"); return d }},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte decimal", input: []byte("5.5"), setup: func() decimal.Decimal { d, _ := decimal.NewFromString("5.5"); return d }},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, setup: func() decimal.Decimal { return decimal.NewFromInt(10) }},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, setup: func() decimal.Decimal { return decimal.Decimal{} }},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Decimal(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			expected := tt.setup()

			// For very small numbers, compare string representations to avoid floating-point precision issues
			if expected.Abs().LessThan(decimal.NewFromFloat(1e-10)) {
				resString, _ := String(res)
				assert.Equal(t, expected.String(), resString, "Unexpected result for input: %v", tt.input)
			} else {
				assert.True(t, expected.Equal(res), "Expected %v, got %v for input: %v", expected, res, tt.input)
			}
		})
	}
}

func TestUnit_Complex64(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    complex64
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -5, expect: -5},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},

		// Complex tests
		{name: "complex64", input: complex64(5 + 12i), expect: 5 + 12i},
		{name: "complex128", input: 5 + 12i, expect: 5 + 12i},
		{name: "complex zero real part", input: 0 + 5i, expect: 0 + 5i},
		{name: "complex negative imaginary", input: 5 - 11.5i, expect: 5 - 11.5i},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, expect: -1000000000},

		// Float tests
		{name: "float32", input: float32(5), expect: 5},
		{name: "float64", input: float64(-5.14), expect: -5.14},
		{name: "float64 scientific positive", input: 1.5134e+02, expect: 151.34},
		{name: "float64 scientific negative", input: 15134e-02, expect: 151.34},
		{name: "float64 max", input: math.MaxFloat64, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 min", input: -math.MaxFloat64, hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 NaN", input: math.NaN(), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 positive infinity", input: math.Inf(1), hasError: true, errorCode: ErrorConvertorNumberOverflow},
		{name: "float64 negative infinity", input: math.Inf(-1), hasError: true, errorCode: ErrorConvertorNumberOverflow},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "string complex", input: "5+6.1i", expect: 5 + 6.1i},
		{name: "string scientific", input: "15134e-02", expect: 151.34},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "string invalid complex", input: "5+6.1i+7i", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Complex64(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Compare real and imaginary parts separately to handle floating-point imprecision
			assert.InDelta(t, real(tt.expect), real(res), 0.0001, "Real part mismatch for input: %v", tt.input)
			assert.InDelta(t, imag(tt.expect), imag(res), 0.0001, "Imaginary part mismatch for input: %v", tt.input)
		})
	}
}

func TestUnit_Complex128(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expect    complex128
		hasError  bool
		errorCode string
	}{
		// Integer tests
		{name: "int", input: int(5), expect: 5},
		{name: "int8", input: int8(5), expect: 5},
		{name: "int16", input: int16(5), expect: 5},
		{name: "int32", input: int32(5), expect: 5},
		{name: "int64", input: int64(5), expect: 5},
		{name: "int negative", input: -5, expect: -5},

		// Unsigned integer tests
		{name: "uint", input: uint(5), expect: 5},
		{name: "uint8", input: uint8(5), expect: 5},
		{name: "uint16", input: uint16(5), expect: 5},
		{name: "uint32", input: uint32(5), expect: 5},
		{name: "uint64", input: uint64(5), expect: 5},

		// Complex tests
		{name: "complex128", input: 5 + 12i, expect: 5 + 12i},
		{name: "complex zero real part", input: 0 + 5i, expect: 0 + 5i},
		{name: "complex negative imaginary", input: 5 - 11.5i, expect: 5 - 11.5i},
		{name: "complex64 conversion", input: complex64(5 - 11.5i), expect: 5 - 11.5i},

		// time.Duration tests
		{name: "time.Duration positive", input: time.Duration(5), expect: 5},
		{name: "time.Duration zero", input: time.Duration(0), expect: 0},
		{name: "time.Duration negative", input: -time.Second, expect: -1000000000},

		// Float tests
		{name: "float32", input: float32(5), expect: 5},
		{name: "float64", input: -5.14, expect: -5.14},
		{name: "float64 scientific positive", input: 1.5134e+02, expect: 151.34},
		{name: "float64 scientific negative", input: 15134e-02, expect: 151.34},
		{name: "float64 max", input: math.MaxFloat64, expect: math.MaxFloat64},
		{name: "float64 min", input: -math.MaxFloat64, expect: -math.MaxFloat64},

		// String and byte slice tests
		{name: "string decimal", input: "5", expect: 5},
		{name: "string complex", input: "5+6.1i", expect: 5 + 6.1i},
		{name: "string scientific", input: "15134e-02", expect: 151.34},
		{name: "string invalid", input: "not a number", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "string invalid complex", input: "5+6.1i+7i", hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "[]byte decimal", input: []byte("5"), expect: 5},
		{name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: 10},
		{name: "[]byte invalid", input: []byte("not a number"), hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// Error cases
		{name: "nil value", input: nil, expect: 0},
		{name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
		{name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Complex128(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if tt.errorCode != "" {
					assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Compare real and imaginary parts separately to handle floating-point imprecision
			assert.InDelta(t, real(tt.expect), real(res), 0.0001, "Real part mismatch for input: %v", tt.input)
			assert.InDelta(t, imag(tt.expect), imag(res), 0.0001, "Imaginary part mismatch for input: %v", tt.input)
		})
	}
}

func TestUnit_DecimalToInt(t *testing.T) {
	tests := []struct {
		name      string
		input     decimal.Decimal
		expect    int
		hasError  bool
		errorCode string
	}{
		{name: "whole", input: dec("42"), expect: 42},
		{name: "truncate positive", input: dec("42.9"), expect: 42},
		{name: "truncate negative", input: dec("-42.9"), expect: -42},
		{name: "zero", input: dec("0"), expect: 0},
		{name: "overflow", input: dec("1e40"), hasError: true, errorCode: ErrorConvertorNumberOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Int(tt.input)
			if tt.hasError {
				assert.NotNil(t, err)
				if tt.errorCode != "" && err != nil {
					assert.True(t, err.Has(tt.errorCode))
				}
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, res)
		})
	}
}

func TestUnit_DecimalToUint(t *testing.T) {
	tests := []struct {
		name     string
		input    decimal.Decimal
		expect   uint
		hasError bool
	}{
		{name: "whole", input: dec("42"), expect: 42},
		{name: "truncate", input: dec("42.9"), expect: 42},
		{name: "negative", input: dec("-1"), hasError: true},
		{name: "overflow", input: dec("1e40"), hasError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Uint(tt.input)
			if tt.hasError {
				assert.NotNil(t, err)
				assert.True(t, err.Has(ErrorConvertorNumberOverflow))
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, res)
		})
	}
}

func TestUnit_DecimalToSignedInts(t *testing.T) {
	// Int64
	t.Run("int64 whole", func(t *testing.T) {
		res, err := Int64(dec("42"))
		assert.Nil(t, err)
		assert.Equal(t, int64(42), res)
	})
	t.Run("int64 overflow", func(t *testing.T) {
		_, err := Int64(dec("9223372036854775808")) // MaxInt64 + 1
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Int32
	t.Run("int32 max", func(t *testing.T) {
		res, err := Int32(dec("2147483647"))
		assert.Nil(t, err)
		assert.Equal(t, int32(math.MaxInt32), res)
	})
	t.Run("int32 overflow", func(t *testing.T) {
		_, err := Int32(dec("2147483648"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Int16
	t.Run("int16 max", func(t *testing.T) {
		res, err := Int16(dec("32767"))
		assert.Nil(t, err)
		assert.Equal(t, int16(math.MaxInt16), res)
	})
	t.Run("int16 overflow", func(t *testing.T) {
		_, err := Int16(dec("32768"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Int8
	t.Run("int8 min", func(t *testing.T) {
		res, err := Int8(dec("-128"))
		assert.Nil(t, err)
		assert.Equal(t, int8(math.MinInt8), res)
	})
	t.Run("int8 overflow", func(t *testing.T) {
		_, err := Int8(dec("128"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})
}

func TestUnit_DecimalToUnsignedInts(t *testing.T) {
	// Uint64
	t.Run("uint64 whole", func(t *testing.T) {
		res, err := Uint64(dec("42"))
		assert.Nil(t, err)
		assert.Equal(t, uint64(42), res)
	})
	t.Run("uint64 negative", func(t *testing.T) {
		_, err := Uint64(dec("-1"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})
	t.Run("uint64 overflow", func(t *testing.T) {
		_, err := Uint64(dec("18446744073709551616")) // MaxUint64 + 1
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Uint32
	t.Run("uint32 max", func(t *testing.T) {
		res, err := Uint32(dec("4294967295"))
		assert.Nil(t, err)
		assert.Equal(t, uint32(math.MaxUint32), res)
	})
	t.Run("uint32 overflow", func(t *testing.T) {
		_, err := Uint32(dec("4294967296"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Uint16
	t.Run("uint16 max", func(t *testing.T) {
		res, err := Uint16(dec("65535"))
		assert.Nil(t, err)
		assert.Equal(t, uint16(math.MaxUint16), res)
	})
	t.Run("uint16 overflow", func(t *testing.T) {
		_, err := Uint16(dec("65536"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Uint8
	t.Run("uint8 max", func(t *testing.T) {
		res, err := Uint8(dec("255"))
		assert.Nil(t, err)
		assert.Equal(t, uint8(math.MaxUint8), res)
	})
	t.Run("uint8 negative", func(t *testing.T) {
		_, err := Uint8(dec("-1"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})
}

func TestUnit_DecimalToFloatsAndComplex(t *testing.T) {
	// Float64
	t.Run("float64 value", func(t *testing.T) {
		res, err := Float64(dec("3.14"))
		assert.Nil(t, err)
		assert.InDelta(t, 3.14, res, 1e-9)
	})
	t.Run("float64 overflow", func(t *testing.T) {
		_, err := Float64(dec("1e400")) // exceeds float64 range -> +Inf
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Float32
	t.Run("float32 value", func(t *testing.T) {
		res, err := Float32(dec("3.14"))
		assert.Nil(t, err)
		assert.InDelta(t, 3.14, float64(res), 1e-4)
	})
	t.Run("float32 overflow", func(t *testing.T) {
		_, err := Float32(dec("1e40")) // exceeds float32 range
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Complex64
	t.Run("complex64 value", func(t *testing.T) {
		res, err := Complex64(dec("3.14"))
		assert.Nil(t, err)
		assert.InDelta(t, 3.14, float64(real(res)), 1e-4)
		assert.Equal(t, float32(0), imag(res))
	})
	t.Run("complex64 overflow", func(t *testing.T) {
		_, err := Complex64(dec("1e40"))
		assert.NotNil(t, err)
		assert.True(t, err.Has(ErrorConvertorNumberOverflow))
	})

	// Complex128
	t.Run("complex128 value", func(t *testing.T) {
		res, err := Complex128(dec("3.14"))
		assert.Nil(t, err)
		assert.InDelta(t, 3.14, real(res), 1e-9)
		assert.Equal(t, float64(0), imag(res))
	})
}
