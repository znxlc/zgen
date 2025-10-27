package zgen

// IsBool will return true if the value represents a bool
func IsBool(value any) bool {
  switch value.(type) {
  case bool:
    return true
  }
  return false
}

// IsNumber will return true if the value represents a number
func IsNumber(value any) bool {
  switch value.(type) {
  case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
    return true
  }
  return false
}
