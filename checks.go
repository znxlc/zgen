package zgen

import "strings"

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
	if value == nil {
		return false
	}
	switch val := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case bool:
		return false
	case string:
		value = strings.TrimSpace(val)
	}
	if _, er := Float64(value); er == nil { // if we can convert to float64, then it's a number
		return true
	}
	if _, er := Int(value); er == nil { // if we can convert to int, then it's a number
		return true
	}
	return false
}
