package zgen

import (
  "github.com/znxlc/zerror"
  "reflect"
  "time"
)

// This package provides a set of functions that handle data extraction, conversion and manipulation

// UnpackBaseElement - removes wrap layers and returns dereferenced base element
// keepPointers=true will unpack only interface elements and keep pointer values intact
// example:
// interface{}(*interface{}(*interface{}(*string))) => string (keepPointers=false)
//	                                                => *string (keepPointers=true)
func UnpackBaseElement(response any, keepPointers bool) any {
  if response == nil {
    return nil
  }
  elemValue := reflect.ValueOf(response)
  result := elemValue
  if elemValue.Kind() == reflect.Ptr || elemValue.Kind() == reflect.Interface {
    dfElem := elemValue.Elem()

    if !dfElem.IsValid() { // element pointed to invalid value, returning nil to prevent a panic
      return nil
    }
    if dfElem.Kind() == reflect.Interface && (!dfElem.IsNil() || !keepPointers) { // continue unpacking if we get a non nil interface
      return UnpackBaseElement(dfElem.Interface(), keepPointers)
    } else if !keepPointers {
      result = dfElem
    }
  }
  res := result.Interface()

  return res
}

// Clone creates a deep copy of whatever is passed to it in an interface{}.
func Clone(src any) any {
  if src == nil {
    return nil
  }

  // Make the interface a reflect.Value
  original := reflect.ValueOf(src)

  // Make a copy of the same type as the original.
  clone := reflect.New(original.Type()).Elem()

  // Recursively copy the original.
  deepCopyRecursive(original, clone)

  // Return the copy as an interface.
  return clone.Interface()
}

// DeepMerge - returns the merged value between 2 elements
//
//	Merge is performed recursively by traversing all the nodes in a map or appending to the existing slices
//	  If FlagDeepMergeOverwriteEnabled is set, the merge is performed only on level 1 of the map, any other type will be overwritten
//	Inputs
//
//	  element1 interface{}
//	    First element to merge
//	  element2 interface{}
//	    Second element to merge
//	  flag int
//	    flag value comprised of the sum of following flags
//	      FlagDeepMergeOverwriteDisabled = 0  // DeepMerge overwrite disabled
//	      FlagDeepMergePriorityFirst     = 1  // Set DeepMerge Priority on first element
//	      FlagDeepMergePrioritySecond    = 2  // Set DeepMerge Priority on second element
//	      FlagDeepMergeOverwriteEnabled  = 10 // DeepMerge overwrite enabled
//	Output
//	  result interface{}
//	    The result of the DeepMerge
//	  err error
//	    can contain type conversion errors
//	Example
//	  DeepMerge( element1, element2 ) // triggers a DeepMerge with overwrite disabled, priority on element2 (defaults)
//	  DeepMerge( element1, element2, FlagDeepMergePriorityFirst ) // triggers a DeepMerge with overwrite disabled, priority on element1
//	  DeepMerge( element1, element2, FlagDeepMergePriorityFirst + FlagDeepMergeOverwriteEnabled ) // triggers a DeepMerge with overwrite enabled, priority on element1
func DeepMerge(element1, element2 interface{}, args ...int) (result interface{}, err zerror.Error) {
  result = nil
  flagPriority := FlagDeepMergePrioritySecond     // second element overwrites the first
  flagOverwrite := FlagDeepMergeOverwriteDisabled // disable overwrite as default

  if len(args) >= 1 { // set priority based on element
    flag := args[0]                            // full flag value
    if flag >= FlagDeepMergeOverwriteEnabled { // check if overwrite is enabled
      flag = flag - FlagDeepMergeOverwriteEnabled
      flagOverwrite = FlagDeepMergeOverwriteEnabled
    }
    if flag == FlagDeepMergePriorityFirst || flag == FlagDeepMergePrioritySecond { // validate if priority flag is correct and set it
      flagPriority = flag
    }
  }

  elem1 := reflect.ValueOf(element1)
  elem2 := reflect.ValueOf(element2)

  if elem1.Kind() == elem2.Kind() {
    switch elem1.Kind() {
    case reflect.Map:
      return deepMergeMap(element1, element2, flagPriority, flagOverwrite)
    case reflect.Slice, reflect.Array:
      return deepMergeSlice(element1, element2, flagPriority, flagOverwrite)
    }
  } else if (elem1.Kind() == reflect.Slice || elem1.Kind() == reflect.Array) &&
    (elem2.Kind() == reflect.Slice || elem2.Kind() == reflect.Array) { // elements are sliceable
    return deepMergeSlice(element1, element2, flagPriority, flagOverwrite)
  }
  if flagPriority == FlagDeepMergePriorityFirst {
    result = element1
  } else {
    result = element2
  }

  return
}

// deepMergeMap - will iterate through 2 MapSI with priority element overwriting the other
func deepMergeMap(element1, element2 any, flagPriority, flagOverwrite int) (result map[string]any, err zerror.Error) {
  result = map[string]interface{}{}

  e1Map, err := MapStringAny(element1)
  if err != nil {
    return result, err
  }
  e2Map, err := MapStringAny(element2)
  if err != nil {
    return result, err
  }

  for key, value := range e1Map {
    result[key] = value
  }

  for key, value := range e2Map {
    if existingValue, found := result[key]; found && flagOverwrite == FlagDeepMergeOverwriteDisabled { // we already have the key and overwrite is disabled, must merge it
      val, err := DeepMerge(existingValue, value, flagPriority+flagOverwrite)
      if err != nil {
        return result, err
      }
      result[key] = val
    } else {
      if !found || flagPriority == FlagDeepMergePrioritySecond {
        result[key] = value
      }
    }
  }

  return
}

// deepMergeSlice - appends slices one after another depending on priority
func deepMergeSlice(element1, element2 interface{}, flagPriority, flagOverwrite int) (result interface{}, err zerror.Error) {
  resultSlice := []interface{}{}
  e1SI, err := SliceAny(element1)
  if err != nil {
    return resultSlice, err
  }

  e2SI, err := SliceAny(element2)
  if err != nil {
    return resultSlice, err
  }

  if flagPriority == FlagDeepMergePriorityFirst { // element1 has priority so it will be added last
    if flagOverwrite == FlagDeepMergeOverwriteEnabled {
      return e1SI, nil
    }
    resultSlice = append(resultSlice, e2SI...)
    resultSlice = append(resultSlice, e1SI...)
  } else { // element2 has priority so it will be added last
    if flagOverwrite == FlagDeepMergeOverwriteEnabled {
      return e2SI, nil
    }
    resultSlice = append(resultSlice, e1SI...)
    resultSlice = append(resultSlice, e2SI...)
  }
  return resultSlice, nil
}

// deepCopyRecursive does the actual copying by value from original to cpy.
// It currently has limited support for what it can handle. Add as needed.
// WARNING - using this on a struct that contains circular links may cause a stack overflow
// params:
//    original reflect.Value - the original value we want to copy
//    cpy      reflect.Value - the element we want to copy the original value into
//    args     ...interface{}
//      0 - iterationCount int     - (optional) current iteration count (checked against the max number of iterations to prevent stack overflows when cloning a circular link element)
func deepCopyRecursive(original, cpy reflect.Value, args ...interface{}) {
  iterationCount := 0
  if len(args) > 0 { // leaving room for future improvements
    // iteration count will be first element, if not a number we will default it to 0
    itc, err := Int(args[0])
    if err != nil {
      itc = 0
    }
    iterationCount = itc + 1 // we increment the iteration counter
  }
  if iterationCount > DeepCopyMaxIterationCount {
    return
  }
  // check for implement deepcopy.Interface
  if original.CanInterface() {
    if copier, ok := original.Interface().(DeepCopier); ok {
      cpy.Set(reflect.ValueOf(copier.DeepCopy()))
      return
    }
  }

  // handle according to original's Kind
  switch original.Kind() {
  case reflect.Ptr:
    // Get the actual value being pointed to.
    originalValue := original.Elem()

    // if  it isn't valid, return.
    if !originalValue.IsValid() {
      return
    }
    cpy.Set(reflect.New(originalValue.Type()))
    deepCopyRecursive(originalValue, cpy.Elem(), iterationCount)

  case reflect.Interface:
    // If this is a nil, don't do anything
    if original.IsNil() {
      return
    }
    // Get the value for the interface, not the pointer.
    originalValue := original.Elem()

    // Get the value by calling Elem().
    copyValue := reflect.New(originalValue.Type()).Elem()
    deepCopyRecursive(originalValue, copyValue, iterationCount)
    cpy.Set(copyValue)

  case reflect.Struct: // only copies exported fields
    t, ok := original.Interface().(time.Time)
    if ok {
      cpy.Set(reflect.ValueOf(t))
      return
    }
    // Go through each field of the struct and copy it.
    for i := 0; i < original.NumField(); i++ {
      // PkgPath is checked, it is only set for unexported fields so we skip them
      if original.Type().Field(i).PkgPath != "" {
        continue
      }
      deepCopyRecursive(original.Field(i), cpy.Field(i), iterationCount)
    }

  case reflect.Slice:
    if original.IsNil() {
      return
    }
    // Make a new slice and copy each element.
    cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
    for i := 0; i < original.Len(); i++ {
      deepCopyRecursive(original.Index(i), cpy.Index(i), iterationCount)
    }

  case reflect.Map:
    if original.IsNil() {
      return
    }
    cpy.Set(reflect.MakeMap(original.Type()))
    for _, key := range original.MapKeys() {
      originalValue := original.MapIndex(key)
      copyValue := reflect.New(originalValue.Type()).Elem()
      deepCopyRecursive(originalValue, copyValue, iterationCount)
      copyKey := Clone(key.Interface())
      cpy.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
    }

  default:
    cpy.Set(original)
  }
}

// IsZeroValue - returns true if that element is of zero value
//	params:
//	  element any - the element we want to check
func IsZeroValue(element any) bool {
  if element == nil {
    return true
  }

  return reflect.ValueOf(element).IsZero()
}
