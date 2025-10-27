package zgen

import (
  "database/sql/driver"
  "github.com/znxlc/zerror"
  "reflect"
  "strings"
)

const (
  TagOmitEmpty  = "omitempty"  // ignores zero value elements
  TagOmitNested = "omitnested" // will not recurse into nested elements leaving them as they are
)

var (
  DefaultParserConfig ParserConfig = ParserConfig{
    EvaluateMethods: false,                  // ignore methods
    KeepPointers:    true,                   // keep pointer values
    Mode:            ParserModeTagsOnly,     // add only tags to the map
    OmitEmpty:       false,                  // do not omit empty fields
    Tags:            []string{"db", "json"}, // by default it parses db and json tags
  }
)

// ToMap - creates a map containing all the keys and properties found in the args. Error is always nil.
// Params:
//   destMap *map[string]any      - destination map pointer, the result will be merged to this map
//   args:                                - any number of arguments as follows
//      config Config                     - configuration settings (if not specified DefaultGenericObjectConfig will be used instead)
//      map map[stringable]any    - each element will be evaluated and added to the result map according to the config
//      struct any struct                 - all the struct elements will be evaluated and a key with the name/tag will be added to the map according with the config settings
//      other                             - any other element type will be ignored
// Note:
//   if config.ParseOnlyTags = false the map will contain the fieldNames as they are (omitempty, omitnested only works with tag fields)
func ToMap(destMap any, args ...any) (err zerror.Error) {
  var returnResponse = map[string]any{}

  currentParseSettings := DefaultParserConfig
  if len(args) > 0 {
    argSlice := make([]any, 0)
    for _, element := range args {
      if settings, ok := element.(ParserConfig); ok {
        currentParseSettings = settings
      } else {
        argSlice = append(argSlice, element)
      }
    }

    for _, element := range argSlice { // parsing data arguments
      element = UnpackBaseElement(element, false)
      elementKind := reflect.TypeOf(element).Kind()
      elemVal := reflect.ValueOf(element)

      if elementKind == reflect.Struct {
        for i := 0; i < elemVal.NumField(); i++ {
          // add field Name to the map if configured
          fieldName := elemVal.Type().Field(i).Name
          // TODO replace with elemVal.Field(i).PkgPath != "" or elemVal.Field(i).IsExported() (as available) when upgrading to a higher go version since this is not available in 1.12
          if fieldName[0:1] != strings.ToUpper(fieldName[0:1]) { // test if field name starts with small case = UNEXPORTED so we will ignore it since it will cause errors later on
            continue
          }

          mapKeyName := fieldName

          fieldValue := UnpackBaseElement(elemVal.Field(i).Interface(), currentParseSettings.KeepPointers)
          fieldKind := reflect.ValueOf(fieldValue).Kind()

          flagFoundTag := false // flag that shows if a tag was found
          // add field tag to the map if configured
          if currentParseSettings.Mode != ParserModeNameOnly {
            for _, tagName := range currentParseSettings.Tags {
              if tagKey, ok := elemVal.Type().Field(i).Tag.Lookup(tagName); ok {
                flagOmitEmpty := false
                flagOmitNested := false

                if tagKey == "-" || tagKey == "" { // skip empty or ignore tags
                  continue
                }

                tagElements := strings.Split(tagKey, ",") // separate elements in tag
                tagKey = tagElements[0]                   // this should exist since at this point we ruled out empty tags

                if len(tagElements) > 1 {
                  for _, te := range tagElements {
                    te = strings.Trim(te, " ") // trim the spaces
                    switch strings.ToLower(te) {
                    case TagOmitEmpty:
                      flagOmitEmpty = true
                    case TagOmitNested:
                      flagOmitNested = true
                    }
                  }
                }
                if flagOmitEmpty && IsZeroValue(fieldValue) {
                  continue
                }
                mapTagKey := tagKey

                if _, ok := returnResponse[mapTagKey]; !ok {
                  // add recursion
                  switch fieldKind {
                  case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
                    if !flagOmitNested {
                      fieldValue, err = s2mParseElement(currentParseSettings, fieldValue)
                      if err != nil {
                        return err
                      }
                    }
                  case reflect.Func: // ignore functions
                    continue
                  }
                  returnResponse[mapTagKey] = fieldValue
                  flagFoundTag = true
                }
              }
            }
          }
          // add field name to the map if configured
          if currentParseSettings.Mode == ParserModeNameOnly || currentParseSettings.Mode == ParserModeNameAndTags || (currentParseSettings.Mode == ParserModeNameIfNoTag && !flagFoundTag) { // add name as key in the map
            if _, ok := returnResponse[mapKeyName]; !ok {
              // add recursion
              switch fieldKind {
              case reflect.Struct, reflect.Map, reflect.Slice:
                fieldValue, err = s2mParseElement(currentParseSettings, fieldValue)
                if err != nil {
                  return err
                }
              case reflect.Func: // ignore functions
                continue
              }
              returnResponse[mapKeyName] = fieldValue
            }
          }
        }
      } else if elementKind == reflect.Map {
        for _, mapKey := range elemVal.MapKeys() {
          if _, ok := returnResponse[mapKey.String()]; !ok {
            fieldValue, err := s2mParseElement(currentParseSettings, elemVal.MapIndex(mapKey).Interface())
            if err != nil {
              return err
            }

            returnResponse[mapKey.String()] = fieldValue
          }
        }
      }
    }
  }

  destValue := reflect.ValueOf(destMap).Elem().Interface().(map[string]any)
  for key, value := range returnResponse {
    destValue[key] = value
  }
  return nil
}

// s2mParseElement - helper function for ToMap that looks at the structure of an element and returns its parsed value
//  Params:
//    config Config - configuration for parser
//    elementData any - the element that needs to be parsed
//                  if simple type it will return the element
//                  if map,struct or slice it will recurse through each element and tries to process each node
func s2mParseElement(config ParserConfig, elementData any) (returnResult any, err zerror.Error) {
  if elementData == nil {
    return nil, nil
  }
  elementData = UnpackBaseElement(elementData, config.KeepPointers) // deepdive into pointers
  eKind := reflect.TypeOf(elementData).Kind()

  switch eKind {
  case reflect.Struct:
    if config.EvaluateMethods {
      // check if struct implements Valuer so we return the value directly
      if fv, ok := elementData.(Valuer); ok { // simple Value() any
        returnResult = fv.Value()
        return returnResult, nil
      } else if fv, ok := elementData.(ValueError); ok { // simple Value() (any,error)
        val, er := fv.Value()
        if er != nil {
          err = zerror.New(ErrorZGENScannerEvaluate, map[string]any{
            "caller":       "s2mParseElement",
            "error":        er.Error(),
            "element_type": "ValueError",
          })
          return nil, err
        }
        return val, nil
      } else if fv, ok := elementData.(driver.Valuer); ok { // sql driver.Valuer
        val, er := fv.Value()
        if er != nil {
          err = zerror.New(ErrorZGENScannerEvaluate, map[string]any{
            "caller":       "s2mParseElement",
            "error":        er.Error(),
            "element_type": "driver.Valuer",
          })
          return nil, err
        }
        return val, nil
      }
    }
    // no valuer so we try to parse the element
    dstMap := map[string]any{}
    err = ToMap(&dstMap, config, elementData) // parse the element
    if err != nil {
      return nil, err
    }
    returnResult = dstMap
    if len(dstMap) == 0 { // && config.RawStructReturnIfNoKeysFound { // return the raw result if no element is in map
      returnResult = elementData
    }

  case reflect.Map:
    dstMap := map[string]any{}
    err = ToMap(&dstMap, config, elementData) // parse the element
    if err != nil {
      return nil, err
    }
    returnResult = dstMap
  case reflect.Slice:
    sliceData, err := SliceAny(elementData)
    sliceResponse := []any{}
    if err != nil {
      return nil, err
    }
    flagFoundNestedElement := false // flag to mark if slice contains struct, map or slice elements
    for _, sliceElement := range sliceData {
      switch reflect.ValueOf(sliceElement).Kind() {
      case reflect.Slice, reflect.Struct, reflect.Map, reflect.Ptr:
        res, err := s2mParseElement(config, sliceElement)
        if err != nil {
          return nil, err
        }
        sliceResponse = append(sliceResponse, res)
        flagFoundNestedElement = true
      default:
        sliceResponse = append(sliceResponse, sliceElement)
      }
    }
    returnResult = elementData // assuming slice is made of simple elements so we keep its value(very important for []byte or []rune elements to remain the same else parsing may be broken
    if flagFoundNestedElement {
      returnResult = sliceResponse
    }
  case reflect.Array:
    sliceData, err := SliceAny(elementData)
    sliceResponse := []any{}
    if err != nil {
      return nil, err
    }
    flagFoundNestedElement := false // flag to mark if slice contains struct, map or slice elements
    for _, sliceElement := range sliceData {
      switch reflect.ValueOf(sliceElement).Kind() {
      case reflect.Slice, reflect.Array, reflect.Struct, reflect.Map:
        res, err := s2mParseElement(config, sliceElement)
        if err != nil {
          return nil, err
        }
        sliceResponse = append(sliceResponse, res)
        flagFoundNestedElement = true
      default:
        sliceResponse = append(sliceResponse, sliceElement)
      }
    }
    returnResult = elementData // assuming slice is made of simple elements so we keep its value(very important for []byte or []rune elements to remain the same else parsing may be broken
    if flagFoundNestedElement {
      returnResult = sliceResponse
    }
  default:
    returnResult = elementData
  }

  return returnResult, nil
}

// ToStruct - fills the struct with all the keys and properties found in the args
//  Params:
//   destStruct *struct                   - destination struct pointer
//   args:                                - any number of arguments as follows
//      config Config                     - configuration settings (if not specified DefaultGenericObjectConfig will be used instead)
//      map map[stringable]any    - each element will be evaluated and added to the destination struct according to the config
//      struct any struct                 - all the struct elements will be evaluated and a key with the name/tag will be added to the destination struct according with the config settings
//      other                             - any other element type will be ignored
func ToStruct(destStruct any, args ...any) zerror.Error {
  var dataMap = map[string]any{}
  structKind := reflect.TypeOf(destStruct).Kind()
  if structKind != reflect.Ptr {
    return zerror.New(ErrorZGENScannerDstStructureInvalid, map[string]any{
      "caller":        "ToStruct",
      "type_expected": "pointer",
    })
  }

  unpackedDestStruct := UnpackBaseElement(destStruct, false)
  udsVal := reflect.ValueOf(unpackedDestStruct)

  if udsVal.Kind() != reflect.Struct {
    return zerror.New(ErrorZGENScannerDstStructureInvalid, map[string]any{
      "caller":        "ToStruct",
      "type_expected": "struct",
    })
  }

  elemVal := reflect.ValueOf(destStruct).Elem()
  if len(args) > 0 {
    newArgs := []any{}

    currentParseSettings := DefaultParserConfig
    mapParseSettings := currentParseSettings
    // make sure we do not miss non tagged fields
    mapParseSettings.Mode = ParserModeNameAndTags

    for _, element := range args {
      if settings, ok := element.(ParserConfig); ok { // parse settings have been sent
        currentParseSettings = settings
        mapParseSettings = currentParseSettings
      } else {
        newArgs = append(newArgs, element)
      }
    }

    newArgs = append(newArgs, mapParseSettings)
    err := ToMap(&dataMap, newArgs...)
    if err != nil {
      return err
    }
    if len(dataMap) == 0 { // no elements in the map
      return zerror.New(ErrorZGENScannerArgumentInvalid, map[string]any{
        "caller": "ToStruct",
        "error":  "Empty map resulted from initial argument",
      })
    }

    if elemVal.Kind() == reflect.Interface { // go past interface to the value behind it
      elemVal = elemVal.Elem()
    }
    for i := 0; i < elemVal.NumField(); i++ {
      fieldVal := elemVal.Field(i)
      fieldType := fieldVal.Type()
      fieldName := elemVal.Type().Field(i).Name
      if fieldVal.IsValid() && fieldVal.CanSet() {
        if value, found := dataMap[fieldName]; found { // search for map fields that match the name of the struct fields
          if value != nil {
            val := reflect.ValueOf(value)
            if fieldType == val.Type() { // map field has the same type as struct field
              fieldVal.Set(val)
            } else {
              err = SetFieldValueByType(currentParseSettings, fieldVal, value)
              if err != nil {
                return err
              }
            }
          }
        } else { // fieldName was not found in the map, we will try the tags
          if currentParseSettings.Mode != ParserModeNameOnly {
            for _, tagName := range currentParseSettings.Tags {
              if tagKey, ok := elemVal.Type().Field(i).Tag.Lookup(tagName); ok {
                if tagKey == "-" || tagKey == "" { // jump over empty tag
                  continue
                }
                tagElements := strings.Split(tagKey, ",")   // separate elements in tag
                tagKey = tagElements[0]                     // this should exist since at this point we ruled out empty tags
                if value, found := dataMap[tagKey]; found { // we found a value by tag
                  if value != nil {
                    val := reflect.ValueOf(value)
                    if fieldType == val.Type() { // struct field type matches map field type
                      fieldVal.Set(val)
                    } else { // we try to convert the type
                      err = SetFieldValueByType(currentParseSettings, fieldVal, value)
                      if err != nil {
                        return err
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  return nil
}

// ScanToElement - sets a value into a compatible or convertible destination element
// use ToMap or ToStruct before calling this method if types are not compatible
//  Params:
//    dstElement any - pointer to the destination element(can be any type)
//    value      any - any value compatible or convertible to dstElement type
//    parseSettings Config   - optional config
func ScanToElement(dstElement, value any, parseSettings ...ParserConfig) (err zerror.Error) {
  currentParserSettings := DefaultParserConfig
  if len(parseSettings) > 0 {
    currentParserSettings = parseSettings[0]
  }

  fieldVal := reflect.ValueOf(dstElement)

  return SetFieldValueByType(currentParserSettings, fieldVal.Elem(), value)
}

// ScanToTemplate - will return the value of the key (if key is a string) or will fill a map[string] with the keys found
func ScanToTemplate(destData, key any, args ...any) zerror.Error {
  var paramMap = map[string]any{}
  _ = ToMap(&paramMap, args...)

  destKeyKind := reflect.TypeOf(destData).Kind()
  if destKeyKind != reflect.Ptr {
    return zerror.New(ErrorZGENScannerDstStructureInvalid, map[string]any{
      "caller": "ScanToTemplate",
      "error":  "destination must be a pointer",
    })
  }

  keyKind := reflect.TypeOf(key).Kind()
  if keyKind == reflect.Map {
    destValue := reflect.ValueOf(destData).Elem().Interface().(map[string]any)
    for _, mapKey := range reflect.ValueOf(key).MapKeys() {
      if resp, ok := paramMap[mapKey.String()]; ok {
        destValue[mapKey.String()] = resp
      }
    }
    return nil
  } else if keyKind == reflect.Slice || keyKind == reflect.Struct { // TODO add slice/array support
    return zerror.New(ErrorZGENScannerDstStructureInvalid, map[string]any{
      "caller": "ScanToTemplate",
      "error":  "Invalid template, key must be string or map",
    })

  } else {
    if resp, ok := paramMap[reflect.ValueOf(key).String()]; ok {
      destValue := reflect.ValueOf(destData).Elem()
      if destValue.CanSet() {
        val := reflect.ValueOf(resp)
        if destValue.Type() == val.Type() {
          destValue.Set(val)
        } else {
          return zerror.New(ErrorZGENScannerArgumentInvalid, map[string]any{
            "caller": "ScanToTemplate",
            "error":  "Argument and destination type mismatch",
          })
        }
      }
      return nil
    }
  }
  return nil
}
