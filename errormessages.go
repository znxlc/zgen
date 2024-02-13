package zgen

import (
  "github.com/znxlc/zerror/errormessage"
)

//Error constants
const (
  ErrorZGENConvertorTypeNotSupported = "ERROR_ZGEN_CONVERTOR_TYPE_NOT_SUPPORTED"
)

//ErrorMap - main error definition map
var ErrorMap = map[string]errormessage.Message{
  //zerror Errors
  ErrorZGENConvertorTypeNotSupported: {
    Code: ErrorZGENConvertorTypeNotSupported,
    Msg:  "ZGEN Conversion Error, Type not supported",
  },
}

func init() {
  for _, value := range ErrorMap {
    errormessage.RegisterErrors(value)
  }
}
