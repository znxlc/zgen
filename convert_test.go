package zgen

import (
  "fmt"
  "math"
  "testing"
  "time"

  "github.com/gofrs/uuid"
  "github.com/shopspring/decimal"
  "github.com/stretchr/testify/assert"
)

func Test_Bool(t *testing.T) {
  if res, err := Bool(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(0); err != nil {
    t.Error(err)
  } else {
    assert.False(t, res)
  }
  if res, err := Bool(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }

  if res, err := Bool(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }

  if res, err := Bool(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }

  if res, err := Bool(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(float64(0)); err != nil {
    t.Error(err)
  } else {
    assert.False(t, res)
  }

  if res, err := Bool(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool(complex128(0)); err != nil {
    t.Error(err)
  } else {
    assert.False(t, res)
  }

  if res, err := Bool([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool([]byte("false")); err != nil {
    t.Error(err)
  } else {
    assert.False(t, res)
  }
  if res, err := Bool("5"); err != nil {
    t.Error(err)
  } else {
    assert.True(t, res)
  }
  if res, err := Bool("false"); err != nil {
    t.Error(err)
  } else {
    assert.False(t, res)
  }

  if _, err := Bool([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Uint(t *testing.T) {
  if res, err := Uint(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }

  if res, err := Uint(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }

  if res, err := Uint(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(0), res)
  }

  if res, err := Uint(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }

  if res, err := Uint(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }
  if res, err := Uint([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(10), res)
  }
  if res, err := Uint("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(5), res)
  }

  if res, err := Uint(9223372036854775807); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint(9223372036854775807), res)
  }

  if _, err := Uint([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Int(t *testing.T) {
  if res, err := Int(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }

  if res, err := Int(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }

  if res, err := Int(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }

  if res, err := Int(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }

  if res, err := Int(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }
  if res, err := Int(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(0), res)
  }

  if res, err := Int([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(10), res)
  }
  if res, err := Int("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(5), res)
  }

  if res, err := Int(9223372036854775807); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int(9223372036854775807), res)
  }

  if _, err := Int([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Int8(t *testing.T) {
  if res, err := Int8(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }

  if res, err := Int8(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(0), res)
  }

  // Test complex number overflow
  _, err := Int8(complex64(math.MaxInt8 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))
  _, err = Int8(complex64(math.MinInt8 - 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int8(complex128(math.MaxInt8 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int8(complex128(math.MinInt8 - 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if res, err := Int8(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }

  if res, err := Int8(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }
  if res, err := Int8([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(10), res)
  }
  if res, err := Int8("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int8(5), res)
  }

  // this conversion will lose the data because we are converting an int64 to int8
  if res, err := Int8(math.MaxInt16); err != nil {
    assert.Equal(t, int8(0), res)
    assert.Equal(t, ErrorConvertorNumberOverflow, err.Get().Code())
  } else {
    assert.Equal(t, int8(-1), res)
  }

  if _, err := Int8([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Int16(t *testing.T) {
  if res, err := Int16(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }

  if res, err := Int16(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }

  if res, err := Int16(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }

  if res, err := Int16(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }
  if res, err := Int16(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(0), res)
  }

  // Test complex number overflow
  _, err := Int16(complex64(math.MaxInt16 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(complex64(math.MinInt16 - 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(complex128(math.MaxInt16 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(complex128(math.MinInt16 - 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if res, err := Int16([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(10), res)
  }
  if res, err := Int16("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int16(5), res)
  }

  // Test overflow cases for Int16
  _, err = Int16(int64(math.MaxInt16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(int64(math.MinInt16 - 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(uint16(math.MaxInt16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(uint32(math.MaxInt16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(uint64(math.MaxInt16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(float32(math.MaxInt16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(float32(math.MinInt16 - 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(float64(math.MaxInt16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(float64(math.MinInt16 - 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int16(math.NaN())
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if _, err := Int16([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Int32(t *testing.T) {
  if res, err := Int32(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }

  if res, err := Int32(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(0), res)
  }

  _, err := Int32(complex128(math.MaxInt32 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int32(complex128(math.MinInt32 - 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if res, err := Int32(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }

  if res, err := Int32(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }
  if res, err := Int32([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(10), res)
  }
  if res, err := Int32("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int32(5), res)
  }

  // Test overflow cases
  _, err = Int32(int64(math.MaxInt32 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int32(int64(math.MinInt32 - 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int32(uint64(math.MaxInt32 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int32(math.MaxFloat64)
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Int32(math.NaN())
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if _, err := Int32([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Int64(t *testing.T) {
  if res, err := Int64(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }

  if res, err := Int64(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(0), res)
  }

  if res, err := Int64(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }

  if res, err := Int64(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  if res, err := Int64([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(10), res)
  }
  if res, err := Int64("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(5), res)
  }
  // this conversion should work
  if res, err := Int64(9223372036854775807); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, int64(9223372036854775807), res)
  }

  if _, err := Int64([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Uint8(t *testing.T) {
  if res, err := Uint8(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }

  if res, err := Uint8(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(0), res)
  }

  // Test complex number overflow
  _, err := Uint8(complex64(math.MaxUint8 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(complex64(-1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(complex128(math.MaxUint8 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(complex128(-1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if res, err := Uint8(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }

  if res, err := Uint8(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }
  if res, err := Uint8([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(10), res)
  }
  if res, err := Uint8("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint8(5), res)
  }

  // Test overflow cases for Uint8
  _, err = Uint8(-1)
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(int8(-1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(int16(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(int32(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(int64(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(uint16(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(uint32(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(uint64(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(float32(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(float64(math.MaxUint8 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint8(math.NaN())
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if _, err := Uint8([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Uint16(t *testing.T) {
  if res, err := Uint16(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }

  if res, err := Uint16(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(0), res)
  }

  // Test complex number overflow
  _, err := Uint16(complex64(math.MaxUint16 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(complex64(-1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(complex128(math.MaxUint16 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(complex128(-1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if res, err := Uint16(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }

  if res, err := Uint16(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }
  if res, err := Uint16([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(10), res)
  }
  if res, err := Uint16("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint16(5), res)
  }

  // Test overflow cases for Uint16
  _, err = Uint16(-1)
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(int16(-1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(int32(math.MaxUint16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(int64(math.MaxUint16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(uint32(math.MaxUint16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(uint64(math.MaxUint16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(float32(math.MaxUint16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(float64(math.MaxUint16 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint16(math.NaN())
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if _, err := Uint16([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Uint32(t *testing.T) {
  if res, err := Uint32(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }

  if res, err := Uint32(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(0), res)
  }

  _, err := Uint32(complex64(-1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint32(complex128(math.MaxUint32 + 1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint32(complex128(-1 + 0i))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if res, err := Uint32(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }

  if res, err := Uint32(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }
  if res, err := Uint32([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(10), res)
  }
  if res, err := Uint32("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint32(5), res)
  }

  // Test overflow cases for Uint32
  _, err = Uint32(-1)
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint32(int32(-1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint32(int64(math.MaxUint32 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint32(uint64(math.MaxUint32 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint32(float64(math.MaxUint32 + 1))
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  _, err = Uint32(math.NaN())
  assert.Error(t, err)
  assert.True(t, err.Has(ErrorConvertorNumberOverflow))

  if _, err := Uint32([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Uint64(t *testing.T) {
  if res, err := Uint64(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }

  if res, err := Uint64(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(0), res)
  }
  if res, err := Uint64(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }

  if res, err := Uint64(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }
  if res, err := Uint64([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(10), res)
  }
  if res, err := Uint64("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, uint64(5), res)
  }

  if _, err := Uint64([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Float32(t *testing.T) {
  if res, err := Float32(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }

  if res, err := Float32(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(0), res)
  }
  if res, err := Float32(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }

  if res, err := Float32(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32(1.5134e+02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(151.34), res)
  }
  if res, err := Float32(15134e-02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(151.34), res)
  }
  if res, err := Float32([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(10), res)
  }
  if res, err := Float32("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(5), res)
  }
  if res, err := Float32("15134e-02"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float32(151.34), res)
  }

  if _, err := Float32([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Float64(t *testing.T) {
  if res, err := Float64(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }

  if res, err := Float64(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(0), res)
  }
  if res, err := Float64(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }

  if res, err := Float64(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64(1.5134e+02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(151.34), res)
  }
  if res, err := Float64(15134e-02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(151.34), res)
  }
  if res, err := Float64([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(10), res)
  }
  if res, err := Float64("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(5), res)
  }
  if res, err := Float64("15134e-02"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, float64(151.34), res)
  }

  if _, err := Float64([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Decimal(t *testing.T) {
  if res, err := Decimal(int(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(int8(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(int16(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(int32(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(int64(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }

  if res, err := Decimal(uint(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(uint8(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(uint16(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(uint32(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(uint64(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(0 + 5i); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(0)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromInt(5)
    assert.Equal(t, dec, res)
  }

  if res, err := Decimal(float32(5.5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromFloat(5.5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(float64(5.5)); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromFloat(5.5)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal(1.5134e+02); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromFloat(1.5134e+02)
    assert.Equal(t, dec, res)
  }

  if res, err := Decimal(1.5134e-02); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromFloat(1.5134e-02)
    resString, _ := String(res)
    assert.Equal(t, dec.String(), resString)
  }

  if res, err := Decimal([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    dec := decimal.NewFromFloat(10)
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal("5.5"); err != nil {
    t.Error(err)
  } else {
    dec, _ := decimal.NewFromString("5.5")
    assert.Equal(t, dec, res)
  }
  if res, err := Decimal("15134e-02"); err != nil {
    t.Error(err)
  } else {
    dec, _ := decimal.NewFromString("15134e-02")
    assert.Equal(t, dec, res)
  }

  if _, err := Decimal([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Complex64(t *testing.T) {
  if res, err := Complex64(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }

  if res, err := Complex64(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(complex64(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5+12i), res)
  }
  if res, err := Complex64(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(0+5i), res)
  }
  if res, err := Complex64(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }

  if res, err := Complex64(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5), res)
  }
  if res, err := Complex64(float64(-5.14)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(-5.14), res)
  }

  if res, err := Complex64(complex64(5 - 11.5i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5-11.5i), res)
  }
  if res, err := Complex64(1.5134e+02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(151.34), res)
  }
  if res, err := Complex64(15134e-02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(151.34), res)
  }
  if res, err := Complex64([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(10), res)
  }
  if res, err := Complex64("5+6.1i"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(5+6.1i), res)
  }
  if res, err := Complex64("15134e-02"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex64(151.34), res)
  }

  if _, err := Complex64([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_Complex128(t *testing.T) {
  if res, err := Complex128(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }

  if res, err := Complex128(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(complex128(5 + 12i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5+12i), res)
  }
  if res, err := Complex128(0 + 5i); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(0+5i), res)
  }
  if res, err := Complex128(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }

  if res, err := Complex128(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5), res)
  }
  if res, err := Complex128(-5.14); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(-5.14), res)
  }

  if res, err := Complex128(complex64(5 - 11.5i)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, 5-11.5i, res)
  }
  if res, err := Complex128(1.5134e+02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(151.34), res)
  }
  if res, err := Complex128(15134e-02); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(151.34), res)
  }
  if res, err := Complex128([]byte{49, 48, 46, 48, 48, 48, 48, 48, 48, 48, 48}); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(10), res)
  }
  if res, err := Complex128("5+6.1i"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(5+6.1i), res)
  }
  if res, err := Complex128("15134e-02"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, complex128(151.34), res)
  }

  if _, err := Complex128([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_String(t *testing.T) {
  if res, err := String(int(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(int8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(int16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(int32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(int64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }

  if res, err := String(uint(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(uint8(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(uint16(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(uint32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(uint64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }

  if res, err := String(complex(-5.14, 4.5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "(-5.14+4.5i)", res)
  }

  if res, err := String(complex64(complex(5, -4.5))); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "(5-4.5i)", res)
  }

  if res, err := String(complex(5, 0)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "(5+0i)", res)
  }

  if res, err := String(complex(0, -3.14)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "(0-3.14i)", res)
  }

  if res, err := String(complex(0, 0)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "(0+0i)", res)
  }

  if res, err := String(time.Duration(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }

  if res, err := String(float32(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(float64(5)); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String(0.000001); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "0.000001", res)
  }
  if res, err := String("5"); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  if res, err := String([]byte{53}); err != nil { // 53 is the ascii for 5
    t.Error(err)
  } else {
    assert.Equal(t, "5", res)
  }
  testUUID, _ := uuid.NewV4()
  if res, err := String(testUUID); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, testUUID.String(), res)
  }

  if _, err := String([]string{}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }
}

func Test_MapStringAny(t *testing.T) {
  // Type not supported error
  if _, err := MapStringAny("value"); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Contains(t, err.Error(), ErrorConvertorTypeNotSupported)
  }

  // Key type not supported
  type TestStruct struct { // struct that has no Stringable function associated so the conversion will fail
    Key string
  }
  if _, err := MapStringAny(map[TestStruct]string{TestStruct{"key"}: "value"}); err == nil {
    t.Error("Should trigger an error")
  } else {
    assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
  }

  mii := map[int]int{1: 1, 2: 2, 3: 3}
  miiconverted := map[string]interface{}{"1": 1, "2": 2, "3": 3}
  if convResult, err := MapStringAny(mii); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, miiconverted, convResult)
  }

  msi := map[string]interface{}{"1": 1, "2": 2, "3": 3}
  if convResult, err := MapStringAny(msi); err != nil {
    t.Error(err)
  } else {
    assert.Equal(t, miiconverted, convResult)
  }

  // if convResult, err := MapStringAny(&msi); err != nil {
  //  fmt.Printf("err: %#v\n", err.GetList()[0].Get())
  //  t.Error(err)
  // } else {
  //  assert.Equal(t, msi, convResult)
  // }

}

func Test_Slice(t *testing.T) {
  si := []int{1, 2, 3}

  res, err := SliceAny(si)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, si[0], res[0])
  ai := [3]int{1, 2, 3}

  res2, err := SliceAny(si)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, ai[0], res2[0])
}

func Test_SliceSimpleType(t *testing.T) {
  si := "first"

  expectedOutput := []interface{}{"first"}
  res, err := SliceAny(si)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, expectedOutput, res)
}

func Test_SliceError(t *testing.T) {
  si := map[string]interface{}{}

  _, err := SliceAny(si)
  if err == nil {
    t.Error("Should Error")
  }

  assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
}

func Test_SliceString(t *testing.T) {
  si := []int{1, 2, 3}

  res, err := SliceString(si)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, "1", res[0])
}

func Test_SliceStringSimpleType(t *testing.T) {
  si := "first"

  expectedOutput := []string{"first"}
  res, err := SliceString(si)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, expectedOutput, res)
}

func Test_SliceStringError(t *testing.T) {
  si := map[string]interface{}{}

  _, err := SliceString(si)
  if err == nil {
    t.Error("Should Error")
  }

  assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())

  si2 := []interface{}{"1", []string{"x", "y"}}

  _, err = SliceString(si2)
  if err == nil {
    t.Error("Should Error")
  }

  assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
}

func Test_SliceInt(t *testing.T) {
  si := []interface{}{"1", "2", "3"}

  res, err := SliceInt(si)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, 1, res[0])
}

func Test_SliceIntError(t *testing.T) {
  si := "invalid"

  _, err := SliceInt(si)
  if err == nil {
    t.Error("Should Error")
  }

  assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())

  si2 := []interface{}{"1", []string{"x", "y"}}

  _, err = SliceInt(si2)
  if err == nil {
    t.Error("Should Error")
  }

  assert.Equal(t, ErrorConvertorTypeNotSupported, err.Error())
}

// func Test_SliceMapStringAny(t *testing.T) {
//  type TestStruct struct {
//    Key   string `json:"Key"`
//    Value string `json:"Value"`
//  }
//  si := []interface{}{
//    map[string]interface{}{
//      "Key":   "testMapKey",
//      "Value": "testMapValue",
//    },
//    TestStruct{
//      Key:   "testStructKey",
//      Value: "testStructValue",
//    }}
//
//  res, err := SliceMapStringAny(si)
//  if err != nil {
//    t.Error(err)
//  }
//
//  assert.Equal(t, "testMapKey", res[0]["Key"])
//  assert.Equal(t, "testStructKey", res[1]["Key"])
// }

func TestUnit_Time(t *testing.T) {
  currTime := time.Now().UTC()
  timeRFC1123 := currTime.Format(time.RFC1123)
  timeIsoDate := currTime.Format(TimeFormatISODate)
  timeIsoDateTime := currTime.Format(TimeFormatISO)
  timeIsoDateTimeSTZ := currTime.Format(TimeFormatISOSTZ)
  timeIsoDateTimeTZ := currTime.Format(TimeFormatISOTZ)
  timeRFC3339 := currTime.Format(time.RFC3339)
  timeRFC3339Nano := currTime.Format(time.RFC3339Nano)
  unixTime := currTime.Unix()
  unixNano := currTime.Nanosecond()

  res, zerr := Time(currTime)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res, currTime)

  res, zerr = Time(timeIsoDate)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Format(TimeFormatISODate), timeIsoDate)

  res, zerr = Time(timeIsoDateTime)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Format(TimeFormatISO), timeIsoDateTime)

  res, zerr = Time(timeRFC1123)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Format(TimeFormatISO), timeIsoDateTime)

  res, zerr = Time(timeIsoDateTimeSTZ)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Format(TimeFormatISOSTZ), timeIsoDateTimeSTZ)

  res, zerr = Time(timeIsoDateTimeTZ)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Format(TimeFormatISOTZ), timeIsoDateTimeTZ)

  res, zerr = Time(timeRFC3339)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Format(time.RFC3339), timeRFC3339)

  res, zerr = Time(timeRFC3339Nano)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Format(time.RFC3339Nano), timeRFC3339Nano)

  res, zerr = Time(unixTime)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.Unix(), unixTime)

  res, zerr = Time(unixTime, unixNano)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.UTC(), currTime)

  dec, err := decimal.NewFromString(fmt.Sprintf("%d.%d", unixTime, unixNano))
  if err != nil {
    t.Error(zerr)
  }

  floatTime, _ := dec.Float64() // this will not be exact no matter what we do
  res, zerr = Time(floatTime)
  if zerr != nil {
    t.Error(zerr)
  }
  assert.Equal(t, res.UTC().Format(time.RFC3339), currTime.Format(time.RFC3339)) // testing only up to seconds
}
