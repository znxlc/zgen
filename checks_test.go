package zgen

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUnit_IsBool(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		// Boolean inputs
		{name: "true boolean", input: true, expected: true},
		{name: "false boolean", input: false, expected: true},

		// String inputs
		{name: "string true", input: "true", expected: false},
		{name: "string false", input: "false", expected: false},
		{name: "empty string", input: "", expected: false},

		// Numeric inputs
		{name: "int zero", input: 0, expected: false},
		{name: "int one", input: 1, expected: false},
		{name: "int negative", input: -1, expected: false},
		{name: "float zero", input: 0.0, expected: false},
		{name: "float one", input: 1.0, expected: false},

		// Other types
		{name: "nil", input: nil, expected: false},
		{name: "slice", input: []int{1, 2, 3}, expected: false},
		{name: "map", input: map[string]int{"key": 1}, expected: false},
		{name: "time.Time", input: time.Now(), expected: false},
		{name: "decimal.Decimal", input: decimal.NewFromInt(123), expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBool(tt.input)
			assert.Equal(t, tt.expected, result, "IsBool(%v) should return %v", tt.input, tt.expected)
		})
	}
}

func TestUnit_IsNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		// Direct numeric types
		{name: "int", input: 42, expected: true},
		{name: "int8", input: int8(42), expected: true},
		{name: "int16", input: int16(42), expected: true},
		{name: "int32", input: int32(42), expected: true},
		{name: "int64", input: int64(42), expected: true},
		{name: "uint", input: uint(42), expected: true},
		{name: "uint8", input: uint8(42), expected: true},
		{name: "uint16", input: uint16(42), expected: true},
		{name: "uint32", input: uint32(42), expected: true},
		{name: "uint64", input: uint64(42), expected: true},
		{name: "float32", input: float32(42.5), expected: true},
		{name: "float64", input: 42.5, expected: true},

		// Numeric edge cases
		{name: "int zero", input: 0, expected: true},
		{name: "int negative", input: -42, expected: true},
		{name: "float zero", input: 0.0, expected: true},
		{name: "float negative", input: -42.5, expected: true},
		{name: "max int64", input: int64(9223372036854775807), expected: true},
		{name: "min int64", input: int64(-9223372036854775808), expected: true},

		// String representations that can be converted
		{name: "string int", input: "42", expected: true},
		{name: "string float", input: "42.5", expected: true},
		{name: "string negative int", input: "-42", expected: true},
		{name: "string negative float", input: "-42.5", expected: true},
		{name: "string zero", input: "0", expected: true},
		{name: "string scientific notation", input: "1.23e-4", expected: true},

		// Non-numeric types
		{name: "boolean true", input: true, expected: false},
		{name: "boolean false", input: false, expected: false},
		{name: "empty string", input: "", expected: false},
		{name: "non-numeric string", input: "hello", expected: false},
		{name: "nil", input: nil, expected: false},
		{name: "slice", input: []int{1, 2, 3}, expected: false},
		{name: "map", input: map[string]int{"key": 1}, expected: false},
		{name: "time.Time", input: time.Now(), expected: true},

		// Edge cases for string conversion
		{name: "string with spaces", input: " 42 ", expected: true},
		{name: "string with leading zeros", input: "042", expected: true},
		{name: "hex string", input: "0x42", expected: false},      // Float64 conversion will fail
		{name: "binary string", input: "0b1010", expected: false}, // Float64 conversion will fail

		// Decimal type
		{name: "decimal.Decimal", input: decimal.NewFromInt(123), expected: false}, // Not in direct types, Float64 conversion fails
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumber(tt.input)
			assert.Equal(t, tt.expected, result, "IsNumber(%v) should return %v", tt.input, tt.expected)
		})
	}
}

func TestUnit_IsNumber_ConversionFallback(t *testing.T) {
	// Test the fallback conversion logic specifically
	t.Run("string to float64 conversion", func(t *testing.T) {
		// These should be true due to Float64 conversion fallback
		assert.True(t, IsNumber("123.456"))
		assert.True(t, IsNumber("-123.456"))
		assert.True(t, IsNumber("0"))
		assert.True(t, IsNumber("1e10"))
		assert.True(t, IsNumber("-1e-10"))
	})

	t.Run("string to int conversion", func(t *testing.T) {
		// These should be true due to Int conversion fallback
		assert.True(t, IsNumber("123"))
		assert.True(t, IsNumber("-123"))
		assert.True(t, IsNumber("0"))
	})

	t.Run("invalid conversions", func(t *testing.T) {
		// These should be false as both conversions fail
		assert.False(t, IsNumber("not a number"))
		assert.False(t, IsNumber("123abc"))
		assert.False(t, IsNumber("1.2.3"))
		assert.False(t, IsNumber(""))
	})
}
