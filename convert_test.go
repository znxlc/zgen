package zgen

import (
  "fmt"
  "github.com/znxlc/zerror"
  "math"
  "testing"
  "time"

  "github.com/gofrs/uuid"
  "github.com/shopspring/decimal"
  "github.com/stretchr/testify/assert"
)

func TestUnit_Bool(t *testing.T) {
  tests := []struct {
    name      string
    input     any
    expect    bool
    hasError  bool
    errorCode string
  }{
    // Integer tests
    {name: "int true", input: int(5), expect: true},
    {name: "int false", input: 0, expect: false},
    {name: "int8 true", input: int8(5), expect: true},
    {name: "int16 true", input: int16(5), expect: true},
    {name: "int64 true", input: int64(5), expect: true},

    // Unsigned integer tests
    {name: "uint true", input: uint(5), expect: true},
    {name: "uint8 true", input: uint8(5), expect: true},
    {name: "uint16 true", input: uint16(5), expect: true},
    {name: "uint64 true", input: uint64(5), expect: true},

    // Time duration
    {name: "time.Duration true", input: time.Duration(5), expect: true},

    // Float tests
    {name: "float32 true", input: float32(5), expect: true},
    {name: "float64 true", input: float64(5), expect: true},
    {name: "float64 false", input: float64(0), expect: false},

    // Complex tests
    {name: "complex64 true", input: complex64(5 + 12i), expect: true},
    {name: "complex128 false", input: complex128(0), expect: false},

    // Byte slice and string tests
    {name: "[]byte true", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: true},
    {name: "[]byte false", input: []byte("false"), expect: false},
    {name: "string true", input: "5", expect: false},
    {name: "string false", input: "false", expect: false},

    // Error case
    {name: "unsupported type", input: []string{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      res, err := Bool(tt.input)
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

      if tt.expect {
        assert.True(t, res, "Expected true but got false for input: %v", tt.input)
      } else {
        assert.False(t, res, "Expected false but got true for input: %v", tt.input)
      }
    })
  }
}

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
          fmt.Printf("res: %v \n", res)
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
      fmt.Printf("Converting %v to %v , expecting %v\n", tt.input, res, tt.expect)
      if tt.hasError {
        if err == nil {
          fmt.Printf("error %#v\n", err.Get().Args())
          t.Error("Expected an error but got none")
        } else if tt.errorCode != "" {
          assert.True(t, err.Has(tt.errorCode), "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
        }
        return
      }

      if err != nil {
        fmt.Printf("error %#v\n", err.Get().Args())
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

func TestUnit_String(t *testing.T) {
  testUUID, _ := uuid.NewV4()

  tests := []struct {
    name      string
    input     any
    expect    string
    hasError  bool
    errorCode string
  }{
    // Integer tests
    {name: "int", input: int(5), expect: "5"},
    {name: "int8", input: int8(5), expect: "5"},
    {name: "int16", input: int16(5), expect: "5"},
    {name: "int32", input: int32(5), expect: "5"},
    {name: "int64", input: int64(5), expect: "5"},
    {name: "int negative", input: -5, expect: "-5"},

    // Unsigned integer tests
    {name: "uint", input: uint(5), expect: "5"},
    {name: "uint8", input: uint8(5), expect: "5"},
    {name: "uint16", input: uint16(5), expect: "5"},
    {name: "uint32", input: uint32(5), expect: "5"},
    {name: "uint64", input: uint64(5), expect: "5"},

    // Complex tests
    {name: "complex128", input: complex(-5.14, 4.5), expect: "(-5.14+4.5i)"},
    {name: "complex64", input: complex64(complex(5, -4.5)), expect: "(5-4.5i)"},
    {name: "complex zero imaginary", input: complex(5, 0), expect: "(5+0i)"},
    {name: "complex zero real", input: complex(0, -3.14), expect: "(0-3.14i)"},
    {name: "complex zero both", input: complex(0, 0), expect: "(0+0i)"},

    // time.Duration tests
    {name: "time.Duration positive", input: time.Duration(5), expect: "5ns"},
    {name: "time.Duration zero", input: time.Duration(0), expect: "0s"},
    {name: "time.Duration negative", input: -time.Second, expect: "-1s"},

    // Float tests
    {name: "float32", input: float32(5), expect: "5"},
    {name: "float64", input: 5.0, expect: "5"},
    {name: "float64 with decimal", input: 5.12345, expect: "5.12345"},
    {name: "float64 small decimal", input: 0.000001, expect: "0.000001"},
    {name: "float64 scientific positive", input: 1.5134e+02, expect: "151.34"},
    {name: "float64 scientific negative", input: 15134e-02, expect: "151.34"},
    {name: "float64 NaN", input: math.NaN(), expect: "NaN"},
    {name: "float64 positive infinity", input: math.Inf(1), expect: "+Inf"},
    {name: "float64 negative infinity", input: math.Inf(-1), expect: "-Inf"},

    // String and byte slice tests
    {name: "string", input: "test", expect: "test"},
    {name: "string number", input: "5", expect: "5"},
    {name: "[]byte", input: []byte{116, 101, 115, 116}, expect: "test"}, // "test" in ASCII
    {name: "[]byte number", input: []byte{53}, expect: "5"},             // 53 is ASCII for '5'
    {name: "[]byte with decimal", input: []byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}, expect: "10.00000000"},

    // UUID test
    {name: "UUID", input: testUUID, expect: testUUID.String()},

    // Error cases
    {name: "nil value", input: nil, expect: ""},
    {name: "unsupported type - slice", input: []string{"1", "2"}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
    {name: "unsupported type - map", input: map[string]int{"a": 1}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
    {name: "unsupported type - struct", input: struct{}{}, hasError: true, errorCode: ErrorConvertorTypeNotSupported},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      res, err := String(tt.input)
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

func TestUnit_MapStringAny(t *testing.T) {
  // Define test structs
  type TestStruct struct {
    Key string
  }

  // Define test cases
  tests := []struct {
    name        string
    input       any
    expected    map[string]any
    hasError    bool
    errorCode   string
    errorString string // For partial error string matching
  }{
    {
      name:        "non-map type",
      input:       "value",
      hasError:    true,
      errorString: ErrorConvertorTypeNotSupported,
    },
    {
      name:      "unsupported key type",
      input:     map[TestStruct]string{{Key: "key"}: "value"},
      hasError:  true,
      errorCode: ErrorConvertorTypeNotSupported,
    },
    {
      name:     "map[int]int",
      input:    map[int]int{1: 1, 2: 2, 3: 3},
      expected: map[string]any{"1": 1, "2": 2, "3": 3},
      hasError: false,
    },
    {
      name:     "map[string]any",
      input:    map[string]any{"1": 1, "2": 2, "3": 3},
      expected: map[string]any{"1": 1, "2": 2, "3": 3},
      hasError: false,
    },
    // Uncomment and modify if needed when pointer to map is supported
    // {
    //   name:     "pointer to map",
    //   input:    &map[string]any{"1": 1, "2": 2, "3": 3},
    //   expected: map[string]any{"1": 1, "2": 2, "3": 3},
    //   hasError: false,
    // },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      result, err := MapStringAny(tt.input)

      if tt.hasError {
        if err == nil {
          t.Error("Expected an error but got none")
          return
        }

        if tt.errorCode != "" {
          assert.True(t, err.Has(tt.errorCode),
            "Expected error code '%s' but got '%s'", tt.errorCode, err.Error())
        }

        if tt.errorString != "" {
          assert.Contains(t, err.Error(), tt.errorString,
            "Expected error to contain '%s' but got '%s'", tt.errorString, err.Error())
        }
        return
      }

      if err != nil {
        t.Errorf("Unexpected error: %v", err)
        return
      }

      assert.Equal(t, tt.expected, result, "Unexpected result for input: %v", tt.input)
    })
  }
}

func TestUnit_Slice(t *testing.T) {
  tests := []struct {
    name     string
    input    any
    expected any
    hasError bool
  }{
    {
      name:     "slice of int",
      input:    []int{1, 2, 3},
      expected: []any{1, 2, 3},
      hasError: false,
    },
    {
      name:     "array of int",
      input:    [3]int{1, 2, 3},
      expected: []any{1, 2, 3},
      hasError: false,
    },
    {
      name:     "slice of string",
      input:    []string{"a", "b", "c"},
      expected: []any{"a", "b", "c"},
      hasError: false,
    },
    {
      name:     "single string",
      input:    "test",
      expected: []any{"test"},
      hasError: false,
    },
    {
      name:     "single int",
      input:    42,
      expected: []any{42},
      hasError: false,
    },
    {
      name:     "unsupported type - map",
      input:    map[string]int{"a": 1},
      hasError: true,
    },
    {
      name:     "nil input",
      input:    nil,
      expected: []any{},
      hasError: false,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      result, err := SliceAny(tt.input)

      if tt.hasError {
        if err == nil {
          t.Error("Expected an error but got none")
        }
        return
      }

      if err != nil {
        t.Errorf("Unexpected error: %v", err)
        return
      }

      assert.Equal(t, tt.expected, result, "Unexpected result for input: %v", tt.input)
    })
  }
}

func TestUnit_SliceSimpleType(t *testing.T) {
  si := "first"

  expectedOutput := []any{"first"}
  res, err := SliceAny(si)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, expectedOutput, res)
}

func TestUnit_SliceError(t *testing.T) {
  si := map[string]any{}

  _, err := SliceAny(si)
  if err == nil {
    t.Error("Should Error")
  }

  assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
}

func TestUnit_SliceString(t *testing.T) {
  tests := []struct {
    name     string
    input    any
    expected []string
    hasError bool
  }{
    {
      name:     "slice of int",
      input:    []int{1, 2, 3},
      expected: []string{"1", "2", "3"},
      hasError: false,
    },
    {
      name:     "array of int",
      input:    [3]int{1, 2, 3},
      expected: []string{"1", "2", "3"},
      hasError: false,
    },
    {
      name:     "slice of string",
      input:    []string{"a", "b", "c"},
      expected: []string{"a", "b", "c"},
      hasError: false,
    },
    {
      name:     "single string",
      input:    "test",
      expected: []string{"test"},
      hasError: false,
    },
    {
      name:     "single int",
      input:    42,
      expected: []string{"42"},
      hasError: false,
    },
    {
      name:     "unsupported type - map",
      input:    map[string]int{"a": 1},
      hasError: true,
    },
    {
      name:     "nested slice",
      input:    []any{"1", []string{"x", "y"}},
      hasError: true,
    },
    {
      name:     "nil input",
      input:    nil,
      expected: []string{},
      hasError: false,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      result, err := SliceString(tt.input)

      if tt.hasError {
        if err == nil {
          t.Error("Expected an error but got none")
        } else {
          assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
        }
        return
      }

      if err != nil {
        t.Errorf("Unexpected error: %v", err)
        return
      }

      assert.Equal(t, tt.expected, result, "Unexpected result for input: %v", tt.input)
    })
  }
}

func TestUnit_SliceInt(t *testing.T) {
  tests := []struct {
    name     string
    input    any
    expected []int
    hasError bool
  }{
    {
      name:     "slice of string numbers",
      input:    []any{"1", "2", "3"},
      expected: []int{1, 2, 3},
      hasError: false,
    },
    {
      name:     "slice of mixed numeric types",
      input:    []any{1, int8(2), int16(3), float32(4), int64(5), "6"},
      expected: []int{1, 2, 3, 4, 5, 6},
      hasError: false,
    },
    {
      name:     "array of int",
      input:    [3]int{1, 2, 3},
      expected: []int{1, 2, 3},
      hasError: false,
    },
    {
      name:     "single int",
      input:    42,
      expected: []int{42},
      hasError: false,
    },
    {
      name:     "single string number",
      input:    "42",
      expected: []int{42},
      hasError: false,
    },
    {
      name:     "unsupported type - string",
      input:    "invalid",
      hasError: true,
    },
    {
      name:     "unsupported type - map",
      input:    map[string]int{"a": 1},
      hasError: true,
    },
    {
      name:     "nested slice",
      input:    []any{"1", []string{"x", "y"}},
      hasError: true,
    },
    {
      name:     "nil input",
      input:    nil,
      expected: []int{},
      hasError: false,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      result, err := SliceInt(tt.input)

      if tt.hasError {
        if err == nil {
          t.Error("Expected an error but got none")
        } else {
          assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
        }
        return
      }

      if err != nil {
        fmt.Printf("got error: %#v", err.Get())
        t.Errorf("Unexpected error: %v", err)
        return
      }

      assert.Equal(t, tt.expected, result, "Unexpected result for input: %v", tt.input)
    })
  }
}

func Test_SliceMapStringAny(t *testing.T) {
  type TestStruct struct {
    Key   string `json:"Key"`
    Value string `json:"Value"`
  }
  type TestStructUnexported struct {
    key   string
    value string
  }
  type TestStructWithTags struct {
    Key   string `json:"key_name"`
    Value string `json:"value_name"`
  }

  tests := []struct {
    name     string
    input    any
    expected []map[string]any
    hasError bool
  }{
    {
      name: "slice of maps and structs",
      input: []any{
        map[string]any{
          "Key":   "testMapKey",
          "Value": "testMapValue",
        },
        TestStruct{
          Key:   "testStructKey",
          Value: "testStructValue",
        },
      },
      expected: []map[string]any{
        {
          "Key":   "testMapKey",
          "Value": "testMapValue",
        },
        {
          "Key":   "testStructKey",
          "Value": "testStructValue",
        },
      },
      hasError: false,
    },
    {
      name: "single map",
      input: map[string]any{
        "Key":   "singleMapKey",
        "Value": "singleMapValue",
      },
      expected: []map[string]any{
        {
          "Key":   "singleMapKey",
          "Value": "singleMapValue",
        },
      },
      hasError: false,
    },
    {
      name: "single struct",
      input: TestStruct{
        Key:   "singleStructKey",
        Value: "singleStructValue",
      },
      expected: []map[string]any{
        {
          "Key":   "singleStructKey",
          "Value": "singleStructValue",
        },
      },
      hasError: false,
    },
    {
      name: "struct with json tags",
      input: TestStructWithTags{
        Key:   "taggedKey",
        Value: "taggedValue",
      },
      expected: []map[string]any{
        {
          "key_name":   "taggedKey",
          "value_name": "taggedValue",
        },
      },
      hasError: false,
    },
    {
      name:     "nil input",
      input:    nil,
      expected: []map[string]any{},
      hasError: false,
    },
    {
      name:     "unsupported type - string",
      input:    "not a map or struct",
      expected: nil,
      hasError: true,
    },
    {
      name:     "unsupported type - slice of strings",
      input:    []string{"a", "b"},
      expected: nil,
      hasError: true,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      result, err := SliceMapStringAny(tt.input)

      if tt.hasError {
        if err == nil {
          t.Error("Expected an error but got none")
        }
        return
      }

      if err != nil {
        t.Errorf("Unexpected error: %v", err)
        return
      }

      assert.Equal(t, len(tt.expected), len(result), "Result length mismatch")

      for i, expectedItem := range tt.expected {
        for key, expectedValue := range expectedItem {
          assert.Equal(t, expectedValue, result[i][key],
            "Mismatch for key '%s' at index %d", key, i)
        }
      }
    })
  }
}
func TestUnit_Time(t *testing.T) {
  // Base time and its components
  currTime := time.Now().UTC()
  year, month, day := currTime.Date()
  hour, minute, sec := currTime.Clock()
  nsec := currTime.Nanosecond()
  unixTime := currTime.Unix()
  floatTime := float64(unixTime) + float64(nsec)/1e9

  // Time formats
  timeRFC1123 := currTime.Format(time.RFC1123)
  timeRFC850 := currTime.Format(time.RFC850)
  timeRFC1123Z := currTime.Format(time.RFC1123Z)
  timeRFC3339 := currTime.Format(time.RFC3339)
  timeRFC3339Nano := currTime.Format(time.RFC3339Nano)
  timeIsoDate := currTime.Format(TimeFormatISODate)
  timeIsoDateTime := currTime.Format(TimeFormatISO)
  timeIsoDateTimeSTZ := currTime.Format(TimeFormatISOSTZ)
  timeIsoDateTimeTZ := currTime.Format(TimeFormatISOTZ)

  // Test cases
  tests := []struct {
    name     string
    input    any
    expected time.Time
    hasError bool
  }{
    // Basic cases
    {name: "nil input", input: nil, expected: time.Time{}},
    {name: "time.Time input", input: currTime, expected: currTime},

    // String formats
    {name: "RFC1123 format", input: timeRFC1123, expected: currTime.Truncate(time.Second)},
    {name: "RFC850 format", input: timeRFC850, expected: currTime.Truncate(time.Second)},
    {name: "RFC1123Z format", input: timeRFC1123Z, expected: currTime.Truncate(time.Second)},
    {name: "RFC3339 format", input: timeRFC3339, expected: currTime.Truncate(time.Second)},
    {name: "RFC3339Nano format", input: timeRFC3339Nano, expected: currTime},
    {name: "ISO date format", input: timeIsoDate, expected: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)},
    {name: "ISO datetime format", input: timeIsoDateTime, expected: time.Date(year, month, day, hour, minute, sec, 0, time.UTC)},
    {name: "ISO datetime with space timezone", input: timeIsoDateTimeSTZ, expected: currTime.Truncate(time.Second)},
    {name: "ISO datetime with Z timezone", input: timeIsoDateTimeTZ, expected: currTime.Truncate(time.Second)},

    // Numeric inputs
    {name: "unix timestamp (seconds)", input: unixTime, expected: time.Unix(unixTime, 0).UTC()},
    {name: "unix timestamp with nanoseconds", input: []any{unixTime, nsec}, expected: time.Unix(unixTime, int64(nsec)).UTC()},
    {name: "unix timestamp as float64", input: floatTime, expected: time.Unix(int64(floatTime), int64((floatTime-float64(int64(floatTime)))*1e9)).UTC()},

    // Component inputs
    {name: "7 date components", input: []any{year, int(month), day, hour, minute, sec, nsec}, expected: currTime},
    {name: "8 date components", input: []any{year, int(month), day, hour, minute, sec, nsec, time.UTC}, expected: currTime},
    {name: "time.Duration", input: time.Duration(nsec), expected: time.Unix(0, int64(nsec)).UTC()},

    // Edge cases
    {name: "empty string", input: "", hasError: true},
    {name: "invalid date string", input: "not a date", hasError: true},
    {name: "unsupported type", input: map[string]string{"key": "value"}, hasError: true},
    {name: "too many components", input: []any{1, 2, 3, 4, 5, 6, 7, 8, 9}, hasError: true},
    {name: "extra month", input: []any{2023, 13, 1, 0, 0, 0, 0}, expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},

    // []byte input
    {name: "[]byte input", input: []byte(timeRFC3339), expected: currTime.Truncate(time.Second)},
  }

  // Run test cases
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      var result time.Time
      var zerr zerror.Error

      // Handle different input types
      switch v := tt.input.(type) {
      case []any:
        result, zerr = Time(v...)
      case nil:
        result, zerr = Time(nil)
      default:
        result, zerr = Time(v)
      }

      if tt.hasError {
        if zerr == nil {
          fmt.Printf("result: %#v\n", result)

          t.Error("Expected error but got none")
        }
        return
      }

      if zerr != nil {
        t.Errorf("Unexpected error: %v", zerr)
        return
      }

      // For time comparisons, handle different precisions
      if !result.Equal(tt.expected) {
        // If nanosecond precision differs but seconds are the same, it might be a formatting issue
        if result.Unix() != tt.expected.Unix() ||
          (result.Nanosecond()/1000) != (tt.expected.Nanosecond()/1000) {
          t.Errorf("Expected %v, got %v", tt.expected, result)
        }
      }
    })
  }

  // Additional test cases for specific scenarios
  t.Run("with time.Duration", func(t *testing.T) {
    d := time.Hour*2 + time.Minute*30
    res, zerr := Time(d)
    if zerr != nil {
      t.Fatalf("Unexpected error: %v", zerr)
    }
    expected := time.Unix(0, d.Nanoseconds()).UTC()
    if !res.Equal(expected) {
      t.Errorf("Expected %v, got %v", expected, res)
    }
  })

  t.Run("with invalid components", func(t *testing.T) {
    _, zerr := Time(2023, 2, 30) // Invalid date
    if zerr == nil {
      t.Error("Expected error for invalid date, got none")
    }
  })
}

// func TestUnit_Time(t *testing.T) {
//   currTime := time.Now().UTC()
//   timeRFC1123 := currTime.Format(time.RFC1123)
//   timeIsoDate := currTime.Format(TimeFormatISODate)
//   timeIsoDateTime := currTime.Format(TimeFormatISO)
//   timeIsoDateTimeSTZ := currTime.Format(TimeFormatISOSTZ)
//   timeIsoDateTimeTZ := currTime.Format(TimeFormatISOTZ)
//   timeRFC3339 := currTime.Format(time.RFC3339)
//   timeRFC3339Nano := currTime.Format(time.RFC3339Nano)
//   unixTime := currTime.Unix()
//   unixNano := currTime.Nanosecond()
//
//   res, zerr := Time(currTime)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res, currTime)
//
//   res, zerr = Time(timeIsoDate)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Format(TimeFormatISODate), timeIsoDate)
//
//   res, zerr = Time(timeIsoDateTime)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Format(TimeFormatISO), timeIsoDateTime)
//
//   res, zerr = Time(timeRFC1123)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Format(TimeFormatISO), timeIsoDateTime)
//
//   res, zerr = Time(timeIsoDateTimeSTZ)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Format(TimeFormatISOSTZ), timeIsoDateTimeSTZ)
//
//   res, zerr = Time(timeIsoDateTimeTZ)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Format(TimeFormatISOTZ), timeIsoDateTimeTZ)
//
//   res, zerr = Time(timeRFC3339)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Format(time.RFC3339), timeRFC3339)
//
//   res, zerr = Time(timeRFC3339Nano)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Format(time.RFC3339Nano), timeRFC3339Nano)
//
//   res, zerr = Time(unixTime)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.Unix(), unixTime)
//
//   res, zerr = Time(unixTime, unixNano)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.UTC(), currTime)
//
//   dec, err := decimal.NewFromString(fmt.Sprintf("%d.%d", unixTime, unixNano))
//   if err != nil {
//     t.Error(zerr)
//   }
//
//   floatTime, _ := dec.Float64() // this will not be exact no matter what we do
//   res, zerr = Time(floatTime)
//   if zerr != nil {
//     t.Error(zerr)
//   }
//   assert.Equal(t, res.UTC().Format(time.RFC3339), currTime.Format(time.RFC3339)) // testing only up to seconds
// }
