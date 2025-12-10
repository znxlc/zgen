package zgen

import (
  "database/sql/driver"
  "github.com/znxlc/zerror"
  "reflect"
  "time"
)

// SetFieldValueByType sets a field value based on the field type with automatic type conversion.
//
// This function handles setting values of various types with appropriate type conversion.
// It supports all basic Go types, pointers, interfaces, and custom types that implement
// the Scanner or driver.Valuer interfaces. For time.Time, it uses the Time() conversion function.
//
// The function makes the following assumptions:
//   - The destination field is valid and settable
//   - Non-settable elements are ignored (handled by the caller)
//   - Nil values without type information are not accepted
//
// Parameters:
//   - currentParseSettings: Configuration for the parser
//   - dstFieldReflectValue: The reflection value of the destination field to set
//   - srcValue: The source value to set (will be converted to the destination type if needed)
//
// Returns:
//   - zerror.Error: Returns an error if the conversion fails, nil otherwise
//
// Supported types:
//   - Basic types (bool, int, float, string, etc.)
//   - Complex numbers (complex64, complex128)
//   - Pointers to any supported type
//   - Interfaces (empty interface{} and interface with methods)
//   - time.Time and *time.Time
//   - Types implementing Scanner or driver.Valuer interfaces
//   - Slices, maps, and channels (with some limitations)
//
// Examples:
//
// 1. Basic type conversion:
//    var i int
//    field := reflect.ValueOf(&i).Elem()
//    err := SetFieldValueByType(Config{}, field, "42")
//    // i == 42
//
// 2. Setting a time.Time field from string:
//    var t time.Time
//    field := reflect.ValueOf(&t).Elem()
//    err := SetFieldValueByType(Config{}, field, "2023-01-02T15:04:05Z")
//
// 3. Setting a pointer field:
//    var s *string
//    field := reflect.ValueOf(&s).Elem()
//    err := SetFieldValueByType(Config{}, field, "pointer value")
//    // *s == "pointer value"
//
// 4. Using with custom types that implement Scanner:
//    type MyInt int
//    func (m *MyInt) Scan(src interface{}) error { /* ... */ }
//    
//    var m MyInt
//    field := reflect.ValueOf(&m).Elem()
//    err := SetFieldValueByType(Config{}, field, 42)
func SetFieldValueByType(currentParseSettings ParserConfig, dstFieldReflectValue reflect.Value, srcValue any) (err zerror.Error) {
  if dstFieldReflectValue.Kind() == reflect.Ptr && !dstFieldReflectValue.CanSet() {
    dstUnpacked := UnpackBaseElement(dstFieldReflectValue.Interface(), true) // this will result in an interface containing the actual element
    dstFieldReflectValue = reflect.ValueOf(dstUnpacked)
  }
  if !dstFieldReflectValue.IsValid() {
    return zerror.New(ErrorZGENInvalidField, map[string]any{
      "field": dstFieldReflectValue,
    })
  } else if dstFieldReflectValue.CanSet() { // we ignore fields that can not be set (like not exported ones) without throwing an error
    srcReflectValue := reflect.ValueOf(srcValue)
    if srcValue == nil {
      if dstFieldReflectValue.IsNil() { // already nil, nothing to be done
        return nil
      }
      switch dstFieldReflectValue.Kind() {
      case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
        dstFieldReflectValue.Set(srcReflectValue)
        return nil
      case reflect.Interface:
        dstFieldReflectValue.Set(reflect.New(reflect.TypeOf(dstFieldReflectValue.Type())).Elem())
        return nil
      }
    }

    if dstFieldReflectValue.Type() == srcReflectValue.Type() { // if they are the same type we set it directly
      dstFieldReflectValue.Set(srcReflectValue)
      return nil
    }
    // different types so we start converting
    switch dstFieldReflectValue.Kind() {
    case reflect.Ptr: // recurse inside the element
      if dstFieldReflectValue.IsNil() && dstFieldReflectValue.CanSet() {
        dstFieldReflectValue.Set(reflect.New(dstFieldReflectValue.Type().Elem()))
      }
      dstFieldValueReflectValue := dstFieldReflectValue.Elem() // the value the pointer is pointing at
      return SetFieldValueByType(currentParseSettings, dstFieldValueReflectValue, srcValue)
    case reflect.Interface: // we have an interface, we set it directly
      if !dstFieldReflectValue.IsNil() { // this interface already contains data so we try to scan to that element
        return SetFieldValueByType(currentParseSettings, dstFieldReflectValue.Elem(), srcValue)
      }
      // empty interface so we just try to add the srcVal data to it as is
      // checking the element for type
      // - interface {}(can hold any element type) - we can set anything inside it
      // - different type = interface{<method definitions here>} (the one that holds method definitions and matches against structs that implement them) - it is impossible to set data to it because it can be any struct that implements those methods
      // ...go should really not reuse types in this manner
      if dstFieldReflectValue.Type().String() == "interface {}" { // only setting data for interface {} type
        dstFieldReflectValue.Set(srcReflectValue)
      }
    case reflect.Struct: // we have a struct field
      if dstFieldScanner, ok := dstFieldReflectValue.Addr().Interface().(Scanner); ok { // we got a scanner
        if srcReflectValue.Kind() == reflect.Struct {
          // testing Valuer variants, unsuccessful scan(err != nil) will be ignored and we try next method
          if vlr, ok := srcReflectValue.Interface().(Valuer); ok { // the srcValue is a Valuer struct
            er := dstFieldScanner.Scan(vlr.Value())
            if er == nil { // scan successful
              return nil
            }
          } else if vlrErr, ok := srcReflectValue.Interface().(ValueError); ok { // the srcValue is a ValueError struct
            vData, er := vlrErr.Value()
            if er == nil {
              er = dstFieldScanner.Scan(vData)
              if er == nil { // scan successful
                return nil
              }
            }
          } else if drvVlr, ok := srcReflectValue.Interface().(driver.Valuer); ok { // the srcValue is a sql driver.Valuer struct
            vData, err := drvVlr.Value()
            if err == nil {
              err = dstFieldScanner.Scan(vData)
              if err == nil { // scan successful
                return nil
              }
            }
          }
        }
        // try direct scan for any other type
        er := dstFieldScanner.Scan(srcReflectValue.Interface())
        if er != nil {
          return zerror.New(ErrorZGENScannerFailed, er)
        }
      } else if _, ok := dstFieldReflectValue.Interface().(time.Time); ok { // we got a time.Time element
        timeSrcValue, err := Time(srcValue)
        if err != nil {
          return err
        }
        dstFieldReflectValue.Set(reflect.ValueOf(timeSrcValue))
      } else if _, ok := dstFieldReflectValue.Interface().(*time.Time); ok { // we got a *time.Time element
        timeSrcValue, er := Time(srcValue)
        if er != nil {
          return er
        }
        dstFieldReflectValue.Elem().Set(reflect.ValueOf(timeSrcValue))
      }
      if srcReflectValue.Kind() == reflect.Struct || srcReflectValue.Kind() == reflect.Map { // we recurse ToStruct
        fv := reflect.New(dstFieldReflectValue.Type()) // create a new element of the same type (pointer)
        err = ToStruct(fv.Interface(), currentParseSettings, srcValue)
        if err != nil {
          return err
        }
        dstFieldReflectValue.Set(fv.Elem()) // passing the srcValue of the pointer
      }
    case reflect.Map: // field is a map, srcValue must be a map (due to the ToMap any other struct should become a map, else the scan will return an error)
      switch srcReflectValue.Kind() {
      case reflect.Map:
        if dstFieldReflectValue.IsNil() { // create a new map if nil, else we will not be able to assign values to it
          newMap := reflect.MakeMap(dstFieldReflectValue.Type())
          dstFieldReflectValue.Set(newMap)
        }
        for _, mapKey := range srcReflectValue.MapKeys() {
          fvKey := reflect.New(dstFieldReflectValue.Type().Key())
          fvVal := reflect.New(dstFieldReflectValue.Type().Elem())
          err = SetFieldValueByType(currentParseSettings, fvKey.Elem(), mapKey.Interface())
          if err != nil {
            return err
          }
          err = SetFieldValueByType(currentParseSettings, fvVal.Elem(), srcReflectValue.MapIndex(mapKey).Interface())
          if err != nil {
            return err
          }
          dstFieldReflectValue.SetMapIndex(fvKey.Elem(), fvVal.Elem())
        }
      default: // no other types are supported to cast to a map
        err = zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
          "src_type": srcReflectValue.Type().String(),
          "dst_type": dstFieldReflectValue.Type().String(),
        })
        return err
      }
    case reflect.Slice:
      if dstFieldReflectValue.Len() > 0 { // we got an existing slice and it needs to be replaced to make sure we do not keep the old data
        newSlice := reflect.New(dstFieldReflectValue.Type())
        dstFieldReflectValue.Set(newSlice.Elem())
      }
      switch srcReflectValue.Kind() {
      case reflect.Slice, reflect.Array:
        for idx := 0; idx < srcReflectValue.Len(); idx++ {
          fvVal := reflect.New(dstFieldReflectValue.Type().Elem())
          err = SetFieldValueByType(currentParseSettings, fvVal.Elem(), srcReflectValue.Index(idx).Interface())
          if err != nil {
            return err
          }
          if idx < dstFieldReflectValue.Len() { // replace existing element with new value
            dstFieldReflectValue.Index(idx).Set(fvVal.Elem())
          } else {
            dstFieldReflectValue.Set(reflect.Append(dstFieldReflectValue, fvVal.Elem()))
          }
        }
      default: // no other types are supported to cast to a slice
        err = zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
          "src_type": srcReflectValue.Type().String(),
          "dst_type": dstFieldReflectValue.Type().String(),
        })
        return err
      }
    case reflect.Array:
      switch srcReflectValue.Kind() {
      case reflect.Slice, reflect.Array:
        maxLen := dstFieldReflectValue.Cap()
        if srcReflectValue.Len() < maxLen {
          maxLen = srcReflectValue.Len()
        }
        for idx := 0; idx < maxLen; idx++ {
          fvVal := reflect.New(dstFieldReflectValue.Type().Elem())
          err = SetFieldValueByType(currentParseSettings, fvVal.Elem(), srcReflectValue.Index(idx).Interface())
          if err != nil {
            return err
          }
          dstFieldReflectValue.Index(idx).Set(fvVal.Elem())
        }
      default: // no other types are supported to cast to array
        err = zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
          "src_type": srcReflectValue.Type().String(),
          "dst_type": dstFieldReflectValue.Type().String(),
        })
        return err
      }
    case reflect.Uint:
      retVal, err := Uint(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Uint8:
      retVal, err := Uint8(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Uint16:
      retVal, err := Uint16(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Uint32:
      retVal, err := Uint32(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Uint64:
      retVal, err := Uint64(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Int:
      retVal, err := Int(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Int8:
      retVal, err := Int8(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Int16:
      retVal, err := Int16(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Int32:
      retVal, err := Int32(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Int64:
      retVal, err := Int64(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Float32:
      retVal, err := Float32(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Float64:
      retVal, err := Float64(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Complex64:
      retVal, err := Complex64(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.Complex128:
      retVal, err := Complex128(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    case reflect.String:
      retVal, err := String(srcReflectValue.Interface())
      if err != nil {
        return err
      }
      dstFieldReflectValue.Set(reflect.ValueOf(retVal).Convert(dstFieldReflectValue.Type()))
    default:
      // TODO add ignoreErrors?
      err = zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
        "src_type": srcReflectValue.Type().String(),
        "dst_type": dstFieldReflectValue.Type().String(),
      })
      return err
    }
  }
  return nil
}

// IsNil safely checks if a reflect.Value is nil without panicking.
// 
// This function is a safer alternative to reflect.Value.IsNil() as it won't panic
// when called on non-nilable types. It handles all nil-checkable types:
//   - Pointers
//   - Slices
//   - Maps
//   - Interfaces
//   - Channels
//   - Functions
// 
// For other types, it always returns false since they cannot be nil in Go.
// 
// Parameters:
//   - element: The reflect.Value to check for nil
// 
// Returns:
//   - bool: true if the value is nil or the zero value of its type, false otherwise
// 
// Examples:
// 
// 1. Basic nil checks:
//    var s *string
//    var i int
//    var m map[string]int
//    var slice []int
//    var iface interface{} = nil
//    
//    IsNil(reflect.ValueOf(s))      // true - nil pointer
//    IsNil(reflect.ValueOf(&i))     // false - non-nil pointer to int
//    IsNil(reflect.ValueOf(m))      // true - nil map
//    IsNil(reflect.ValueOf(slice))  // true - nil slice
//    IsNil(reflect.ValueOf(iface))  // true - nil interface
//    IsNil(reflect.ValueOf(42))     // false - non-pointer type
//
// 2. Working with structs and pointers to structs:
//    type Person struct{ Name string }
//    var p1 *Person
//    p2 := &Person{Name: "Alice"}
//    var p3 Person
//    
//    IsNil(reflect.ValueOf(p1))  // true - nil pointer to struct
//    IsNil(reflect.ValueOf(p2))  // false - non-nil pointer to struct
//    IsNil(reflect.ValueOf(p3))  // false - non-pointer struct (can't be nil in Go)
//
// 3. Checking collections:
//    var m1 map[string]int
//    m2 := make(map[string]int)
//    var s1 []int
//    s2 := []int{}
//    
//    IsNil(reflect.ValueOf(m1))  // true - nil map
//    IsNil(reflect.ValueOf(m2))  // false - initialized empty map
//    IsNil(reflect.ValueOf(s1))  // true - nil slice
//    IsNil(reflect.ValueOf(s2))  // false - initialized empty slice
//
// 4. Safe nil check in a function with type assertion:
//    func processValue(v interface{}) error {
//        rv := reflect.ValueOf(v)
//        if IsNil(rv) {
//            return fmt.Errorf("value cannot be nil")
//        }
//        
//        // Safe to use v as it's not nil
//        // Additional processing...
//        return nil
//    }
//
// 5. Checking function and channel types:
//    var f func()
//    var ch chan int
//    
//    IsNil(reflect.ValueOf(f))   // true - nil function
//    IsNil(reflect.ValueOf(ch))  // true - nil channel
//    
//    f = func() { fmt.Println("Hello") }
//    ch = make(chan int)
//    
//    IsNil(reflect.ValueOf(f))   // false - initialized function
//    IsNil(reflect.ValueOf(ch))  // false - initialized channel
//
// 6. Working with interfaces and type assertions:
//    var i interface{}
//    IsNil(reflect.ValueOf(i))  // true - nil interface
//    
//    i = "some string"
//    IsNil(reflect.ValueOf(i))  // false - non-nil interface containing a string
//    
//    var ptr *int
//    i = ptr
//    IsNil(reflect.ValueOf(i))  // true - interface containing a nil pointer
//    
//    val := 42
//    i = &val
//    IsNil(reflect.ValueOf(i))  // false - interface containing a non-nil pointer
func IsNil(element reflect.Value) bool {
  switch element.Kind() {
  case reflect.Invalid:
    return true
  case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan, reflect.Func:
    return element.IsNil()
  default:
    return false
  }
}
