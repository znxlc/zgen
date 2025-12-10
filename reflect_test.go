package zgen

import (
  "database/sql/driver"
  "fmt"
  "reflect"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"
  "github.com/znxlc/zerror"
)

// CustomType implements the Scanner interface for testing
type CustomType struct {
  Val string
}

// Scan implements the Scanner interface for CustomType
func (c *CustomType) Scan(src interface{}) error {
  if src == nil {
    c.Val = ""
    return nil
  }
  if s, ok := src.(string); ok {
    c.Val = "scanned: " + s
    return nil
  }
  return zerror.New("invalid type", nil)
}

// CustomValuer implements the driver.Valuer interface for testing
type CustomValuer struct {
  Val string
}

// Value implements the driver.Valuer interface for CustomValuer
func (c CustomValuer) Value() (driver.Value, error) {
  return "valuer: " + c.Val, nil
}

// ScannerValuer implements both Scanner and driver.Valuer interfaces for testing
type ScannerValuer struct {
  Val string
}

// Scan implements the Scanner interface for ScannerValuer
func (s *ScannerValuer) Scan(src interface{}) error {
  if src == nil {
    s.Val = ""
    return nil
  }
  s.Val = "scanned: " + src.(string)
  return nil
}

// Value implements the driver.Valuer interface for ScannerValuer
func (s ScannerValuer) Value() (driver.Value, error) {
  return "valuer: " + s.Val, nil
}

// TestStruct is used for testing field setting functionality
type TestStruct struct {
  IntField    int
  StringField string
  TimeField   time.Time
  PtrField    *int
  SliceField  []string
  MapField    map[string]int
  CustomField CustomType
}

// TestStructWithZero is used for testing zero value handling
type TestStructWithZero struct {
  IntField    int
  StringField string
  SliceField  []string
}

func TestIsNil(t *testing.T) {
  tests := []struct {
    name     string
    input    interface{}
    expected bool
  }{
    // Nil values
    {name: "nil interface", input: interface{}(nil), expected: true},
    {"nil pointer", (*int)(nil), true},
    {"nil slice", ([]int)(nil), true},
    {"nil map", (map[string]int)(nil), true},
    {"nil channel", (chan int)(nil), true},
    {"nil function", (func())(nil), true},

    // Non-nil values
    {"non-nil int", 42, false},
    {"non-nil string", "test", false},
    {"non-nil struct", struct{}{}, false},
    {"non-nil array", [3]int{1, 2, 3}, false},

    // Pointers to zero values
    {"pointer to zero int", new(int), false},
    {"pointer to empty struct", &struct{}{}, false},

    // Initialized collections
    {"empty slice", []int{}, false},
    {"empty map", map[string]int{}, false},
    {"initialized channel", make(chan int), false},
    {"initialized function", func() {}, false},

    // Custom types
    {"custom type with nil value", (*CustomType)(nil), true},
    {"initialized custom type", &CustomType{}, false},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      result := IsNil(reflect.ValueOf(tt.input))
      assert.Equal(t, tt.expected, result, "IsNil() = %v, want %v", result, tt.expected)
    })
  }
}

func TestSetFieldValueByType(t *testing.T) {
  now := time.Now().UTC()
  intVal := 42

  // Custom types for testing
  customValuer := CustomValuer{Val: "custom valuer"}
  scannerValuer := ScannerValuer{Val: "scanner valuer"}

  tests := []struct {
    name        string
    fieldName   string
    value       interface{}
    expectError bool
    expected    interface{}
  }{
    // Integer types
    {
      name:        "set int field from string",
      fieldName:   "IntField",
      value:       "42",
      expectError: false,
      expected:    42,
    },
    {
      name:        "set int field from int",
      fieldName:   "IntField",
      value:       42,
      expectError: false,
      expected:    42,
    },
    {
      name:        "set int field from int8",
      fieldName:   "IntField",
      value:       int8(42),
      expectError: false,
      expected:    42,
    },
    {
      name:        "set int field from uint",
      fieldName:   "IntField",
      value:       uint(42),
      expectError: false,
      expected:    42,
    },
    {
      name:        "set int field from float",
      fieldName:   "IntField",
      value:       42.0,
      expectError: false,
      expected:    42,
    },

    // String type
    {
      name:        "set string field",
      fieldName:   "StringField",
      value:       "test string",
      expectError: false,
      expected:    "test string",
    },
    {
      name:        "set string field from int",
      fieldName:   "StringField",
      value:       42,
      expectError: false,
      expected:    "42",
    },
    {
      name:        "set string field from bool",
      fieldName:   "StringField",
      value:       true,
      expectError: false,
      expected:    "true",
    },

    // Time type
    {
      name:        "set time field from string",
      fieldName:   "TimeField",
      value:       now.Format(time.RFC3339),
      expectError: false,
      expected:    now.Truncate(time.Second), // Truncate to ignore nanoseconds for comparison
    },
    {
      name:        "set time field from time.Time",
      fieldName:   "TimeField",
      value:       now,
      expectError: false,
      expected:    now.Truncate(time.Second),
    },

    // Pointer type
    {
      name:        "set pointer field",
      fieldName:   "PtrField",
      value:       &intVal,
      expectError: false,
      expected:    &intVal,
    },
    {
      name:        "set pointer field from value",
      fieldName:   "PtrField",
      value:       84,
      expectError: false,
      expected:    func() *int { v := 84; return &v }(),
    },

    // Slice type
    {
      name:        "set slice field",
      fieldName:   "SliceField",
      value:       []string{"one", "two", "three"},
      expectError: false,
      expected:    []string{"one", "two", "three"},
    },
    {
      name:        "set slice field from string",
      fieldName:   "SliceField",
      value:       "one,two,three",
      expectError: true,
    },

    // Map type
    {
      name:        "set map field",
      fieldName:   "MapField",
      value:       map[string]int{"one": 1, "two": 2},
      expectError: false,
      expected:    map[string]int{"one": 1, "two": 2},
    },

    // Custom types
    {
      name:        "set custom type with scanner",
      fieldName:   "CustomField",
      value:       "test value",
      expectError: false,
      expected:    CustomType{Val: "scanned: test value"},
    },
    {
      name:        "set custom valuer",
      fieldName:   "StringField",
      value:       customValuer,
      expectError: true,
    },
    {
      name:        "set scanner valuer",
      fieldName:   "StringField",
      value:       scannerValuer,
      expectError: true,
    },

    // Error cases
    {
      name:        "invalid type conversion",
      fieldName:   "IntField",
      value:       "not a number",
      expectError: true,
    },
    {
      name:        "invalid field name",
      fieldName:   "NonexistentField",
      value:       "value",
      expectError: true,
    },
    {
      name:        "invalid time format",
      fieldName:   "TimeField",
      value:       "not a time",
      expectError: true,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      s := &TestStruct{}
      field := reflect.ValueOf(s).Elem().FieldByName(tt.fieldName)
      fmt.Printf("field: %v\n", field)
      err := SetFieldValueByType(DefaultParserConfig, field, tt.value)

      fmt.Printf("result: %v, %#v\n", err, s)

      if tt.expectError {
        fmt.Printf("Element: %#v\n", s)
        assert.Error(t, err, "Expected error for %s", tt.name)
        return
      }

      assert.NoError(t, err, "Unexpected error for %s: %v", tt.name, err)

      // Special handling for time.Time to avoid precision issues
      if tt.fieldName == "TimeField" {
        assert.WithinDuration(t, tt.expected.(time.Time), s.TimeField, time.Second)
        return
      }

      // Get the field value for comparison
      fieldValue := reflect.ValueOf(s).Elem().FieldByName(tt.fieldName).Interface()

      // Handle pointer comparison
      if reflect.TypeOf(tt.expected).Kind() == reflect.Ptr {
        assert.Equal(t, tt.expected, fieldValue, "Field %s value mismatch", tt.fieldName)
      } else {
        // For non-pointer values, use DeepEqual for better comparison of slices/maps
        assert.True(t, reflect.DeepEqual(tt.expected, fieldValue),
          "Field %s value mismatch: expected %v, got %v",
          tt.fieldName, tt.expected, fieldValue)
      }
    })
  }
}

// Test setting pointer fields with automatic allocation
func TestSetFieldValueByType_PointerFields(t *testing.T) {
  initialValue := 10

  tests := []struct {
    name        string
    initial     *TestStruct
    fieldName   string
    value       interface{}
    expectError bool
    expected    int
  }{
    {
      name:        "set nil pointer field with value",
      initial:     &TestStruct{},
      fieldName:   "PtrField",
      value:       42,
      expectError: false,
      expected:    42,
    },
    {
      name:        "set existing pointer field with new value",
      initial:     &TestStruct{PtrField: &initialValue},
      fieldName:   "PtrField",
      value:       84,
      expectError: false,
      expected:    84,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      s := tt.initial
      field := reflect.ValueOf(s).Elem().FieldByName(tt.fieldName)

      err := SetFieldValueByType(ParserConfig{}, field, tt.value)

      if tt.expectError {
        assert.Error(t, err, "Expected an error but got none")
      } else {
        assert.NoError(t, err, "Unexpected error")
      }

      assert.NotNil(t, s.PtrField, "PtrField should not be nil after setting value")
      assert.Equal(t, tt.expected, *s.PtrField, "PtrField has unexpected value")
    })
  }
}

// Test setting zero values
func TestSetFieldValueByType_ZeroValues(t *testing.T) {
  tests := []struct {
    name      string
    fieldName string
    initial   interface{}
    value     interface{}
    expected  interface{}
  }{
    {
      name:      "set zero int",
      fieldName: "IntField",
      initial:   42,
      value:     0,
      expected:  0,
    },
    {
      name:      "set empty string",
      fieldName: "StringField",
      initial:   "initial",
      value:     "",
      expected:  "",
    },
    {
      name:      "set empty slice",
      fieldName: "SliceField",
      initial:   []string{"one", "two"},
      value:     []string{},
      expected:  []string{},
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // Create test struct and get the field
      s := &TestStructWithZero{}
      field := reflect.ValueOf(s).Elem().FieldByName(tt.fieldName)

      // Set initial value
      field.Set(reflect.ValueOf(tt.initial))

      // Set zero value
      err := SetFieldValueByType(ParserConfig{}, field, tt.value)
      assert.NoError(t, err)

      // Get the field value for verification
      fieldValue := reflect.ValueOf(s).Elem().FieldByName(tt.fieldName).Interface()

      // Verify the value was set correctly
      if slice, ok := fieldValue.([]string); ok {
        assert.Empty(t, slice, "Expected empty slice")
      } else {
        assert.Equal(t, tt.expected, fieldValue)
      }
    })
  }
}
