package zgen

import (
  "database/sql"
  "database/sql/driver"
  "encoding/json"
  "errors"
  "time"

  "github.com/araddon/dateparse"
  "github.com/lib/pq"
)

// Number is a constraint interface that defines all numeric types supported by the package.
// This includes both signed and unsigned integers of various sizes, as well as
// floating-point numbers. The interface is used for type constraints in generic
// functions and validation operations.
type Number = interface {
  ~int | ~int8 | ~int16 | ~int32 | ~int64 |
  ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
  ~float32 | ~float64
}

// Valuer - generic interface that implements Value() without error
type Valuer interface {
  Value() any
}

// ValueError - generic interface that implements Value() with error
type ValueError interface {
  Value() (any, error)
}

// Stringable - generic interface that implements String()
type Stringable interface {
  String() string
}

// ToStringable - generic interface that implements ToString()
type ToStringable interface {
  ToString() string
}

// Scanner - generic interface that implements Scan(any) error
type Scanner interface {
  Scan(value any) error
}

// DeepCopier is an interface for delegating copy process to type
type DeepCopier interface {
  DeepCopy() any
}

// These types create or extend frequently used types so that they allow json and sql marshalling

// NullBool - nullable boolean extension
type NullBool struct {
  sql.NullBool
}

// UnmarshalJSON - extension to make element compatible with json.Unmarshal
func (nb *NullBool) UnmarshalJSON(data []byte) error {
  if string(data) == "null" || string(data) == "nil" {
    nb.Valid = false
    return nil
  }

  var temp bool
  if err := json.Unmarshal(data, &temp); err != nil {
    return err
  }
  nb.Bool = temp
  nb.Valid = true

  return nil
}

// MarshalJSON - extension to make element compatible with json.Marshal
func (nb *NullBool) MarshalJSON() ([]byte, error) {
  if nb.Valid == false {
    return json.Marshal(nil)
  }

  return json.Marshal(nb.Bool)
}

// NullInt64 - nullable int64 extension
type NullInt64 struct {
  sql.NullInt64
}

// UnmarshalJSON - extension to make element compatible with json.Unmarshal
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
  if string(data) == "null" || string(data) == "nil" {
    ni.Valid = false
    return nil
  }

  var temp int64
  if err := json.Unmarshal(data, &temp); err != nil {
    return err
  }
  ni.Int64 = temp
  ni.Valid = true

  return nil
}

// MarshalJSON - extension to make element compatible with json.Marshal
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
  if ni.Valid == false {
    return json.Marshal(nil)
  }

  return json.Marshal(ni.Int64)
}

// NullFloat64 - nullable float64 extension
type NullFloat64 struct {
  sql.NullFloat64
}

// UnmarshalJSON - extension to make element compatible with json.Unmarshal
func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
  if string(data) == "null" || string(data) == "nil" {
    nf.Valid = false
    return nil
  }

  var temp float64
  if err := json.Unmarshal(data, &temp); err != nil {
    return err
  }
  nf.Float64 = temp
  nf.Valid = true

  return nil
}

// MarshalJSON - extension to make element compatible with json.Marshal
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
  if nf.Valid == false {
    return json.Marshal(nil)
  }

  return json.Marshal(nf.Float64)
}

// NullString - nullable string extension
type NullString struct {
  sql.NullString
}

// UnmarshalJSON - extension to make element compatible with json.Unmarshal
func (ns *NullString) UnmarshalJSON(data []byte) error {
  if string(data) == "null" || string(data) == "nil" {
    ns.Valid = false
    return nil
  }

  var temp string
  if err := json.Unmarshal(data, &temp); err != nil {
    return err
  }
  ns.String = temp
  ns.Valid = true

  return nil
}

// MarshalJSON - extension to make element compatible with json.Marshal
func (ns *NullString) MarshalJSON() ([]byte, error) {
  if ns.Valid == false {
    return json.Marshal(nil)
  }

  return json.Marshal(ns.String)
}

// NullTime - nullable time extension
type NullTime struct {
  pq.NullTime
}

// UnmarshalJSON - extension to make element compatible with json.Unmarshal
func (nt *NullTime) UnmarshalJSON(data []byte) error {
  if string(data) == "null" || string(data) == "nil" {
    nt.Valid = false
    return nil
  }

  var temp time.Time
  if err := json.Unmarshal(data, &temp); err != nil {
    return err
  }
  nt.Time = temp
  nt.Valid = true

  return nil
}

// MarshalJSON - extension to make element compatible with json.Marshal
func (nt *NullTime) MarshalJSON() ([]byte, error) {
  if nt.Valid == false {
    return json.Marshal(nil)
  }

  return json.Marshal(nt.Time)
}

// Scan - Scan override to support parsing from strings
func (nt *NullTime) Scan(value interface{}) (err error) {
  if value == nil {
    nt.Valid = false
  } else {
    switch val := value.(type) {
    case NullTime:
      nt.Valid = val.Valid
      nt.Time = val.Time
    case *NullTime:
      nt.Valid = val.Valid
      nt.Time = val.Time
    case time.Time:
      nt.Time = val
      nt.Valid = true
    case string:
      nt.Valid = true
      nt.Time, err = dateparse.ParseAny(val)
      if err != nil {
        nt.Valid = false
      }
    case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
      str, _ := String(val)
      nt.Valid = true
      nt.Time, err = dateparse.ParseAny(str)
      if err != nil {
        nt.Valid = false
      }
    default:
      return errors.New(ErrorConvertorTypeNotSupported)
    }
  }

  return nil
}

// DBJSONField is a helper type for storing json data in the db
type DBJSONField map[string]any

// Value is a helper to transform the field into json
func (sett DBJSONField) Value() (driver.Value, error) {
  return json.Marshal(sett)
}

// Scan will read the db field into the data structure
func (sett *DBJSONField) Scan(value any) error {
  b, ok := value.([]byte)
  if !ok {
    return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &sett)
}
