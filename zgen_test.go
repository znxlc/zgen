package zgen

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CloneCircularTestStruct1 struct {
	Value     int
	Circular2 *CloneCircularTestStruct2
}
type CloneCircularTestStruct2 struct {
	Value     int
	Circular1 *CloneCircularTestStruct1
}

func TestUnit_UnpackBaseElement(t *testing.T) {
	tests := []struct {
		name         string
		input        any
		keepPointers bool
		expect       any
		expectType   string
		hasError     bool
	}{
		// ===== keepPointers = false (unpack everything) =====

		// Basic types without wrapping
		{name: "nil input", input: nil, keepPointers: false, expect: nil},
		{name: "int value", input: 42, keepPointers: false, expect: 42, expectType: "int"},
		{name: "string value", input: "test", keepPointers: false, expect: "test", expectType: "string"},
		{name: "bool value", input: true, keepPointers: false, expect: true, expectType: "bool"},
		{name: "float64 value", input: 3.14, keepPointers: false, expect: 3.14, expectType: "float64"},

		// Single interface wrapping
		{name: "interface{} int", input: interface{}(42), keepPointers: false, expect: 42, expectType: "int"},
		{name: "interface{} string", input: interface{}("test"), keepPointers: false, expect: "test", expectType: "string"},
		{name: "interface{} bool", input: interface{}(true), keepPointers: false, expect: true, expectType: "bool"},

		// Single pointer wrapping (should unpack)
		{name: "pointer to int", input: intPtr(42), keepPointers: false, expect: 42, expectType: "int"},
		{name: "pointer to string", input: stringPtr("test"), keepPointers: false, expect: "test", expectType: "string"},
		{name: "pointer to bool", input: boolPtr(true), keepPointers: false, expect: true, expectType: "bool"},

		// Nested interface wrapping
		{name: "double interface int", input: interface{}(interface{}(42)), keepPointers: false, expect: 42, expectType: "int"},
		{name: "triple interface string", input: interface{}(interface{}(interface{}("test"))), keepPointers: false, expect: "test", expectType: "string"},

		// Mixed interface and pointer wrapping (should unpack)
		{name: "interface pointer int", input: interface{}(intPtr(42)), keepPointers: false, expect: 42, expectType: "int"},
		{name: "pointer interface string", input: anyPtr("test"), keepPointers: false, expect: "test", expectType: "string"},
		{name: "interface pointer interface int", input: interface{}(anyPtr(42)), keepPointers: false, expect: 42, expectType: "int"},

		// Complex nested scenarios (should unpack)
		{name: "complex nested 1", input: anyPtr(anyPtr(42)), keepPointers: false, expect: 42, expectType: "int"},
		{name: "complex nested 2", input: interface{}(anyPtr(anyPtr("test"))), keepPointers: false, expect: "test", expectType: "string"},

		// &interface{} scenarios (should unpack)
		{name: "&interface{} int", input: anyPtr(42), keepPointers: false, expect: 42, expectType: "int"},
		{name: "&interface{} string", input: anyPtr("test"), keepPointers: false, expect: "test", expectType: "string"},
		{name: "&interface{} bool", input: anyPtr(true), keepPointers: false, expect: true, expectType: "bool"},

		// Nested &interface{} scenarios (should unpack)
		{name: "nested &interface{} 1", input: anyPtr(anyPtr(42)), keepPointers: false, expect: 42, expectType: "int"},
		{name: "nested &interface{} 2", input: anyPtr(anyPtr(anyPtr("test"))), keepPointers: false, expect: "test", expectType: "string"},

		// Mixed &interface{} with regular interfaces (should unpack)
		{name: "interface{} &interface{} int", input: interface{}(anyPtr(42)), keepPointers: false, expect: 42, expectType: "int"},
		{name: "&interface{} interface{} string", input: anyPtr(interface{}("test")), keepPointers: false, expect: "test", expectType: "string"},

		// Complex mixed scenarios with &interface{} (should unpack)
		{name: "complex mixed 1", input: anyPtr(interface{}(anyPtr(42))), keepPointers: false, expect: 42, expectType: "int"},
		{name: "complex mixed 2", input: interface{}(anyPtr(interface{}(anyPtr("test")))), keepPointers: false, expect: "test", expectType: "string"},

		// Nil scenarios
		{name: "nil pointer", input: (*int)(nil), keepPointers: false, expect: nil},
		{name: "nil interface", input: interface{}(nil), keepPointers: false, expect: nil},
		{name: "interface nil pointer", input: interface{}((*int)(nil)), keepPointers: false, expect: nil},
		{name: "&interface{} nil", input: anyPtr(nil), keepPointers: false, expect: nil},
		{name: "&interface{} nil pointer", input: anyPtr((*any)(nil)), keepPointers: false, expect: nil},

		// Complex types (should unpack)
		{name: "struct value", input: struct{ Name string }{Name: "test"}, keepPointers: false, expect: struct{ Name string }{Name: "test"}, expectType: "struct { Name string }"},
		{name: "slice value", input: []int{1, 2, 3}, keepPointers: false, expect: []int{1, 2, 3}, expectType: "[]int"},
		{name: "map value", input: map[string]int{"a": 1}, keepPointers: false, expect: map[string]int{"a": 1}, expectType: "map[string]int"},

		// Pointer to complex types (should unpack)
		{name: "pointer to struct", input: &struct{ Name string }{Name: "test"}, keepPointers: false, expect: struct{ Name string }{Name: "test"}, expectType: "struct { Name string }"},
		{name: "pointer to slice", input: &[]int{1, 2, 3}, keepPointers: false, expect: []int{1, 2, 3}, expectType: "[]int"},
		{name: "pointer to map", input: &map[string]int{"a": 1}, keepPointers: false, expect: map[string]int{"a": 1}, expectType: "map[string]int"},

		// Interface wrapping complex types (should unpack)
		{name: "interface struct", input: interface{}(struct{ Name string }{Name: "test"}), keepPointers: false, expect: struct{ Name string }{Name: "test"}, expectType: "struct { Name string }"},
		{name: "interface slice", input: interface{}([]int{1, 2, 3}), keepPointers: false, expect: []int{1, 2, 3}, expectType: "[]int"},
		{name: "interface map", input: interface{}(map[string]int{"a": 1}), keepPointers: false, expect: map[string]int{"a": 1}, expectType: "map[string]int"},

		// &interface{} with complex types (should unpack)
		{name: "&interface{} struct", input: anyPtr(struct{ Name string }{Name: "test"}), keepPointers: false, expect: struct{ Name string }{Name: "test"}, expectType: "struct { Name string }"},
		{name: "&interface{} slice", input: anyPtr([]int{1, 2, 3}), keepPointers: false, expect: []int{1, 2, 3}, expectType: "[]int"},
		{name: "&interface{} map", input: anyPtr(map[string]int{"a": 1}), keepPointers: false, expect: map[string]int{"a": 1}, expectType: "map[string]int"},

		// ===== keepPointers = true (preserve pointers) =====

		// Basic pointer types (should preserve)
		{name: "pointer to int keep", input: intPtr(42), keepPointers: true, expect: intPtr(42), expectType: "*int"},
		{name: "pointer to string keep", input: stringPtr("test"), keepPointers: true, expect: stringPtr("test"), expectType: "*string"},
		{name: "pointer to bool keep", input: boolPtr(true), keepPointers: true, expect: boolPtr(true), expectType: "*bool"},
		{name: "pointer to any keep", input: anyPtr(true), keepPointers: true, expect: anyPtr(true), expectType: "*interface {}"},

		// Complex nested scenarios (should preserve outermost pointer)
		{name: "complex nested keep pointers", input: anyPtr(anyPtr(42)), keepPointers: true, expect: anyPtr(anyPtr(42)), expectType: "*interface {}"},

		// &interface{} scenarios (should preserve)
		{name: "&interface{} int keep", input: anyPtr(42), keepPointers: true, expect: anyPtr(42), expectType: "*interface {}"},

		// Nested &interface{} scenarios (should preserve)
		{name: "nested &interface{} keep", input: anyPtr(anyPtr(42)), keepPointers: true, expect: anyPtr(anyPtr(42)), expectType: "*interface {}"},

		// Mixed &interface{} with regular interfaces (should preserve)
		{name: "interface{} &interface{} keep", input: interface{}(anyPtr(42)), keepPointers: true, expect: anyPtr(42), expectType: "*interface {}"},

		// Complex mixed scenarios with &interface{} (should preserve)
		{name: "complex mixed keep", input: anyPtr(interface{}(anyPtr(42))), keepPointers: true, expect: anyPtr(interface{}(anyPtr(42))), expectType: "*interface {}"},

		// Pointer to complex types (should preserve)
		{name: "pointer to struct keep", input: &struct{ Name string }{Name: "test"}, keepPointers: true, expect: &struct{ Name string }{Name: "test"}, expectType: "*struct { Name string }"},
		{name: "pointer to slice keep", input: &[]int{1, 2, 3}, keepPointers: true, expect: &[]int{1, 2, 3}, expectType: "*[]int"},
		{name: "pointer to map keep", input: &map[string]int{"a": 1}, keepPointers: true, expect: &map[string]int{"a": 1}, expectType: "*map[string]int"},

		// &interface{} with complex types (should preserve)
		{name: "&interface{} struct keep", input: anyPtr(struct{ Name string }{Name: "test"}), keepPointers: true, expect: anyPtr(struct{ Name string }{Name: "test"}), expectType: "*interface {}"},
		{name: "&interface{} slice keep", input: anyPtr([]int{1, 2, 3}), keepPointers: true, expect: anyPtr([]int{1, 2, 3}), expectType: "*interface {}"},
		{name: "&interface{} map keep", input: anyPtr(map[string]int{"a": 1}), keepPointers: true, expect: anyPtr(map[string]int{"a": 1}), expectType: "*interface {}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnpackBaseElement(tt.input, tt.keepPointers)

			if tt.expect != nil {
				assert.Equal(t, tt.expect, result, "Expected value mismatch for input: %v", tt.input)
			}

			if tt.expectType != "" {
				assert.Equal(t, tt.expectType, reflect.TypeOf(result).String(), "Expected type mismatch for input: %v", tt.input)
			}

			if tt.expect == nil {
				assert.Nil(t, result, "Expected nil result for input: %v", tt.input)
			}
		})
	}
}

func TestUnit_IsZeroValue(t *testing.T) {
	var i int
	assert.Equal(t, true, IsZeroValue(i))
	i = 5
	assert.Equal(t, false, IsZeroValue(i))

	var str string
	assert.Equal(t, true, IsZeroValue(str))
	str = "a"
	assert.Equal(t, false, IsZeroValue(str))

	var ar [2]int
	assert.Equal(t, true, IsZeroValue(ar))
	ar[1] = 5
	assert.Equal(t, false, IsZeroValue(ar))

	var mp map[string]string
	assert.Equal(t, true, IsZeroValue(mp))
	mp = map[string]string{"a": "b"}
	assert.Equal(t, false, IsZeroValue(mp))

	var Cf struct {
		IntVal  int
		BoolVal bool
	}

	assert.Equal(t, false, IsZeroValue(&Cf))
	assert.Equal(t, true, IsZeroValue(Cf))
	Cf.BoolVal = true
	assert.Equal(t, false, IsZeroValue(Cf))
}

func TestUnit_CloneMap(t *testing.T) {
	src := map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}
	clone, _ := MapStringAny(Clone(src))

	assert.Equal(t, src["key1"], clone["key1"])
	src["key1"] = "othervalue" // changing value, if clone is not a real clone this should change
	assert.NotEqual(t, src["key1"], clone["key1"])
}

func TestUnit_CloneSlice(t *testing.T) {
	src := []interface{}{
		"val1",
		"val2",
	}
	clone, _ := SliceAny(Clone(src))

	assert.Equal(t, src[0], clone[0])
	src[0] = "othervalue" // changing value, if clone is not a real clone this should change
	assert.NotEqual(t, src[0], clone[0])
}

func TestUnit_CloneStruct(t *testing.T) {
	type CloneTestStruct struct {
		Value int
	}
	src := CloneTestStruct{Value: 1}

	clone, ok := Clone(src).(CloneTestStruct)
	if !ok {
		t.Error("Clone struct failed")
	}

	assert.Equal(t, src.Value, clone.Value)
	src.Value = 2 // changing value, if clone is not a real clone this should change
	assert.NotEqual(t, src.Value, clone.Value)
}

func TestUnit_CloneCircular(t *testing.T) {
	str2 := CloneCircularTestStruct2{
		Value: 1,
	}
	str1 := CloneCircularTestStruct1{
		Value:     1,
		Circular2: &str2,
	}
	str2.Circular1 = &str1

	clone, ok := Clone(str1).(CloneCircularTestStruct1) // this will trigger max iteration count and should stop there
	if !ok {
		t.Error("Clone struct failed")
	}

	assert.Equal(t, str1.Value, clone.Value)
	assert.Equal(t, str1.Circular2.Value, clone.Circular2.Value)
	str1.Value = 2 // changing value, if clone is not a real clone this should change
	str2.Value = 2
	assert.NotEqual(t, str1.Value, clone.Value)
	assert.NotEqual(t, str1.Circular2.Value, clone.Circular2.Value)
}

func TestUnit_CloneSimple(t *testing.T) {
	src := "x"
	clone, _ := String(Clone(src))

	assert.Equal(t, src, clone) // all simple types are passed by value so we do not need to test changed values
}

func TestUnit_DeepMergeSimple(t *testing.T) {
	elem1 := "x"
	elem2 := 12

	res, err := DeepMerge(elem1, elem2, 1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "x", res)

	res, err = DeepMerge(elem1, elem2, 2)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 12, res)
}

func TestUnit_DeepMergeSlice(t *testing.T) {
	elem1 := []interface{}{1, "2", 3.14}
	elem2 := [3]interface{}{4, "5", 6.7}

	res, err := DeepMerge(elem1, elem2, FlagDeepMergePriorityFirst)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []interface{}{4, "5", 6.7, 1, "2", 3.14}, res)

	res, err = DeepMerge(elem1, elem2, FlagDeepMergePrioritySecond)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []interface{}{1, "2", 3.14, 4, "5", 6.7}, res)

	res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePriorityFirst)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []interface{}{1, "2", 3.14}, res)

	res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePrioritySecond)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []interface{}{4, "5", 6.7}, res)

}

func TestUnit_DeepMergeMap(t *testing.T) {
	elem1 := map[string]interface{}{
		"simpleInt":    12,
		"simpleString": "a string",
		"simpleMap": map[string]interface{}{
			"x1": 1,
			"x2": "x2Val1",
			"x3": []interface{}{1, "2", 3.14},
			"x4": map[string]interface{}{
				"x4v1": "1",
				"x4v2": 2,
				"x4v3": []interface{}{1, "2", 3.14},
			},
		},
	}
	elem2 := map[string]interface{}{
		"simpleIntg":   13,
		"simpleString": "another string",
		"simpleMap": map[string]interface{}{
			"x1.1": 1,
			"x2":   "x2Val2",
			"x3":   []interface{}{4, "5", 6.7},
			"x4": map[string]interface{}{
				"x4vx": "2",
				"x4v2": 3,
				"x4v3": [3]interface{}{4, "5", 6.7}, // array here, will be converted to slice
			},
		},
	}
	elemResult1 := map[string]interface{}{
		"simpleInt":    12,
		"simpleIntg":   13,
		"simpleString": "a string",
		"simpleMap": map[string]interface{}{
			"x1":   1,
			"x1.1": 1,
			"x2":   "x2Val1",
			"x3":   []interface{}{4, "5", 6.7, 1, "2", 3.14},
			"x4": map[string]interface{}{
				"x4v1": "1",
				"x4vx": "2",
				"x4v2": 2,
				"x4v3": []interface{}{4, "5", 6.7, 1, "2", 3.14},
			},
		},
	}

	elemResult2 := map[string]interface{}{
		"simpleInt":    12,
		"simpleIntg":   13,
		"simpleString": "another string",
		"simpleMap": map[string]interface{}{
			"x1":   1,
			"x1.1": 1,
			"x2":   "x2Val2",
			"x3":   []interface{}{1, "2", 3.14, 4, "5", 6.7},
			"x4": map[string]interface{}{
				"x4v1": "1",
				"x4vx": "2",
				"x4v2": 3,
				"x4v3": []interface{}{1, "2", 3.14, 4, "5", 6.7},
			},
		},
	}

	elemResult3 := map[string]interface{}{
		"simpleInt":    12,
		"simpleIntg":   13,
		"simpleString": "a string",
		"simpleMap": map[string]interface{}{
			"x1": 1,
			"x2": "x2Val1",
			"x3": []interface{}{1, "2", 3.14},
			"x4": map[string]interface{}{
				"x4v1": "1",
				"x4v2": 2,
				"x4v3": []interface{}{1, "2", 3.14},
			},
		},
	}

	elemResult4 := map[string]interface{}{
		"simpleInt":    12,
		"simpleIntg":   13,
		"simpleString": "another string",
		"simpleMap": map[string]interface{}{
			"x1.1": 1,
			"x2":   "x2Val2",
			"x3":   []interface{}{4, "5", 6.7},
			"x4": map[string]interface{}{
				"x4vx": "2",
				"x4v2": 3,
				"x4v3": [3]interface{}{4, "5", 6.7}, // array here, will be converted to slice
			},
		},
	}

	res, err := DeepMerge(elem1, elem2, FlagDeepMergePriorityFirst)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, elemResult1, res)

	res, err = DeepMerge(elem1, elem2, FlagDeepMergePrioritySecond)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, elemResult2, res)

	res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePriorityFirst)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, elemResult3, res)

	res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePrioritySecond)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, elemResult4, res)

	res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, elemResult4, res)
}
