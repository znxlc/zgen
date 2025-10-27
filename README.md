# zgen

[![Go Reference](https://pkg.go.dev/badge/github.com/znxlc/zgen.svg)](https://pkg.go.dev/github.com/znxlc/zgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/znxlc/zgen)](https://goreportcard.com/report/github.com/znxlc/zgen)

A Go utility library providing type conversion and data manipulation functions with almost zero dependencies (except for `github.com/shopspring/decimal` for decimal operations).

## Features

- **Type Conversion**: Safely convert between all basic Go types (int, float, string, bool, etc.)
- **Deep Copying**: Create deep copies of complex data structures
- **Merging**: Deep merge maps and slices with configurable priority
- **Zero Value Checks**: Easily check if a value is the zero value for its type
- **Pointer Unpacking**: Safely extract values from nested pointers and interfaces

## Installation

```bash
go get github.com/znxlc/zgen
```

## Usage

### Type Conversion

Convert between different types easily:

```go
import "github.com/znxlc/zgen"

// String to int
num, err := zgen.Int("42") // 42, nil

// Float to string
str, err := zgen.String(3.14) // "3.14", nil

```

### Deep Copy

Create deep copies of complex data structures:

```go
original := map[string]any{
    "nested": []int{1, 2, 3},
    "value":  "test",
}

copied := zgen.Clone(original).(map[string]any)
```

### Deep Merge

Merge maps and slices with configurable priority:

```go
map1 := map[string]any{"a": 1, "b": 2}
map2 := map[string]any{"b": 3, "c": 4}

// Merge with priority on second element
merged, err := zgen.DeepMerge(map1, map2)
// Result: map[a:1 b:3 c:4]

// Merge with priority on first element
merged, err := zgen.DeepMerge(map1, map2, zgen.FlagDeepMergePriorityFirst)
// Result: map[a:1 b:2 c:4]
```

## Available Functions

### Type Conversion
- `Int()`, `Int8()`, `Int16()`, `Int32()`, `Int64()`
- `Uint()`, `Uint8()`, `Uint16()`, `Uint32()`, `Uint64()`
- `Float32()`, `Float64()`
- `Complex64()`, `Complex128()`
- `String()`
- `Bool()`
- `Time()`
- `Decimal()`
- `MapStringAny()`
- `SliceAny()`, `SliceByte()`, `SliceString()`, `SliceInt()`, `SliceMapStringAny()`

### Data Manipulation
- `Clone(any) any` - Create a deep copy of any value
- `DeepMerge(any, any, ...int) (any, error)` - Deep merge two values
- `IsZeroValue(any) bool` - Check if a value is the zero value for its type
- `UnpackBaseElement(any, bool) any` - Safely extract values from nested pointers and interfaces

## Error Handling

All conversion functions return a `zerror.Error` type which provides detailed error information. Check errors using the standard Go pattern:

```go
value, err := zgen.Int("not a number")
if err != nil {
    // Handle error
    log.Printf("Conversion failed: %v", err)
    return
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
