package zgen

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
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
		{name: "float32", input: float32(3.14), expect: "3.14"},
		{name: "float64", input: 3.14, expect: "3.14"},
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
		{name: "[]any stringable", input: []any{1.11, 2, "ab", 3, true, []any{1, 2, 3}}, expect: "1.112ab3true123"},
		{name: "[]any non stringable", input: []any{true, "x", "-12.35", []any{1, 2, 3}, map[string]any{"a": 1}}, expect: "", hasError: true, errorCode: ErrorConvertorTypeNotSupported},

		// UUID test
		{name: "UUID", input: testUUID, expect: testUUID.String()},

		// Error cases
		{name: "nil value", input: nil, expect: ""},
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

func TestUnit_MapStringSliceByte(t *testing.T) {
	// Define test structs
	type TestStruct struct {
		Key string
	}

	// Define test cases
	tests := []struct {
		name        string
		input       any
		expected    map[string][]byte
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
			name:     "map[int]string",
			input:    map[int]string{1: "test1", 2: "test2", 3: "test3"},
			expected: map[string][]byte{"1": []byte("test1"), "2": []byte("test2"), "3": []byte("test3")},
			hasError: false,
		},
		{
			name:     "map[string]string",
			input:    map[string]string{"key1": "value1", "key2": "value2"},
			expected: map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			hasError: false,
		},
		{
			name:     "map[string][]byte",
			input:    map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			expected: map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			hasError: false,
		},
		{
			name:     "map[int][]byte",
			input:    map[int][]byte{1: []byte("value1"), 2: []byte("value2")},
			expected: map[string][]byte{"1": []byte("value1"), "2": []byte("value2")},
			hasError: false,
		},
		{
			name:     "map[string]any with strings",
			input:    map[string]any{"key1": "value1", "key2": "value2"},
			expected: map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			hasError: false,
		},
		{
			name:     "map[string]any with []byte",
			input:    map[string]any{"key1": []byte("value1"), "key2": []byte("value2")},
			expected: map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			hasError: false,
		},
		{
			name:     "map[string]any with mixed types",
			input:    map[string]any{"key1": "value1", "key2": []byte("value2"), "key3": 123},
			expected: map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2"), "key3": []byte("123")},
			hasError: false,
		},
		{
			name:     "map[string]any with boolean values",
			input:    map[string]any{"key1": true, "key2": false, "key3": true},
			expected: map[string][]byte{"key1": []byte("true"), "key2": []byte("false"), "key3": []byte("true")},
			hasError: false,
		},
		{
			name:      "map[string]any with unsupported value type",
			input:     map[string]any{"key1": struct{}{}},
			hasError:  true,
			errorCode: ErrorConvertorTypeNotSupported,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: map[string][]byte{},
			hasError: false,
		},
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: map[string][]byte{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MapStringSliceByte(tt.input)

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

func TestUnit_SliceByte(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		expected    []byte
		hasError    bool
		errorCode   string
		errorString string
	}{
		{
			name:     "bool true",
			input:    true,
			expected: []byte("true"),
			hasError: false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: []byte("false"),
			hasError: false,
		},
		{
			name:     "string",
			input:    "hello",
			expected: []byte("hello"),
			hasError: false,
		},
		{
			name:     "[]byte",
			input:    []byte("test"),
			expected: []byte("test"),
			hasError: false,
		},
		{
			name:     "int",
			input:    123,
			expected: []byte("123"),
			hasError: false,
		},
		{
			name:     "float",
			input:    45.67,
			expected: []byte("45.67"),
			hasError: false,
		},
		{
			name:     "nil",
			input:    nil,
			expected: []byte{},
			hasError: false,
		},
		{
			name:        "unsupported type",
			input:       struct{}{},
			hasError:    true,
			errorString: ErrorConvertorTypeNotSupported,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SliceByte(tt.input)

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
			name:     "nested slice stringable",
			input:    []any{"1", []string{"x", "y"}},
			expected: []string{"1", "xy"},
		},
		{
			name:     "nested slice not stringable",
			input:    []any{"1", map[string]any{"x": "y"}},
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

// Helper functions for creating pointers
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func anyPtr(a any) *any {
	return &a
}

// dec is a small helper to build a decimal.Decimal from its string form in tests.
func dec(s string) decimal.Decimal {
	return decimal.RequireFromString(s)
}
