package zgen

import (
  "github.com/znxlc/zerror/errormessage"
)

// Error constants
const (
  ErrorConvertorTypeNotSupported = "ERROR_ZGEN_CONVERTOR_TYPE_NOT_SUPPORTED"
  ErrorConvertorNumberOverflow   = "ERROR_ZGEN_CONVERTOR_NUMBER_OVERFLOW"

  // Scanner Errors
  ErrorZGENScannerEvaluate            = "ERROR_ZGEN_SCANNER_EVALUATE"
  ErrorZGENScannerDstStructureInvalid = "ERROR_ZGEN_SCANNER_DST_STRUCTURE_INVALID"
  ErrorZGENScannerArgumentInvalid     = "ERROR_ZGEN_SCANNER_ARGUMENT_INVALID"
  ErrorZGENScannerFailed              = "ERROR_ZGEN_SCANNER_FAILED"

  ErrorZGENInvalidField = "ERROR_ZGEN_INVALID_FIELD"
)

// ErrorMap - main error definition map
var ErrorMap = map[string]errormessage.Message{
  // zgen Errors
  ErrorConvertorTypeNotSupported: {
    Code: ErrorConvertorTypeNotSupported,
    Msg:  "ZGEN Conversion Error, Type not supported",
  },
  ErrorConvertorNumberOverflow: {
    Code: ErrorConvertorNumberOverflow,
    Msg:  "ZGEN Conversion Error, number overflow",
  },

  // Scanner errors
  ErrorZGENScannerEvaluate: {
    Code: ErrorZGENScannerEvaluate,
    Msg:  "zgen.Scanner: Unable to evaluate value",
  },
  ErrorZGENScannerDstStructureInvalid: {
    Code: ErrorZGENScannerDstStructureInvalid,
    Msg:  "zgen.Scanner: Destination structure is invalid",
  },
  ErrorZGENScannerArgumentInvalid: {
    Code: ErrorZGENScannerArgumentInvalid,
    Msg:  "zgen.Scanner: Argument to assign is Invalid",
  },
  ErrorZGENScannerFailed: {
    Code: ErrorZGENScannerFailed,
    Msg:  "zgen.Scanner: Failed to scan value",
  },

  ErrorZGENInvalidField: {
    Code: ErrorZGENInvalidField,
    Msg:  "zgen: Invalid field",
  },
}

func init() {
  for _, value := range ErrorMap {
    errormessage.RegisterErrors(value)
  }
}
