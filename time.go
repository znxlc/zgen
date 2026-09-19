package zgen

import (
	"reflect"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/znxlc/zerror"
)

// Time - tries to convert any to time.Time
// params:
//
//		src
//		   1 - string                - will try parse the string using all formats in TimeLayoutMap
//		   1,2 - number              - (int, uint, float types) - will be converted to integers and assumes unix time and/or unixnano time
//		   7 numbers                 - time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
//		   7 numbers + location      - time.Date(year, time.Month(month), day, hour, min, sec, nsec, location), uses time.UTC if location is nil
//		   time.Time                 - returns value as is
//	       time.Duration             - returns
//		   NullTime / pq.NullTime    - returns the underlying time when valid, otherwise returns an error
//		   other                     - will return an error
func Time(args ...any) (dst time.Time, err zerror.Error) {
	result := time.Time{}
	if len(args) == 0 {
		return result, nil
	}

	// Handle single argument
	if len(args) == 1 {
		if args[0] == nil {
			return result, nil
		}

		// Handle time.Time directly
		if timeVal, ok := args[0].(time.Time); ok {
			return timeVal, nil
		}

		// Handle time.Duration
		if timeVal, ok := args[0].(time.Duration); ok {
			return time.Unix(0, timeVal.Nanoseconds()), nil
		}

		// Handle NullTime / pq.NullTime (and pointers) - returns the underlying
		// time when valid, otherwise returns a conversion error.
		switch nt := args[0].(type) {
		case NullTime:
			if nt.Valid {
				return nt.Time, nil
			}
			return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      args[0],
				"src_type": "nil or invalid time",
				"dst_type": "time",
			})
		case *NullTime:
			if nt != nil && nt.Valid {
				return nt.Time, nil
			}
			return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      args[0],
				"src_type": "nil or invalid time",
				"dst_type": "time",
			})
		case pq.NullTime:
			if nt.Valid {
				return nt.Time, nil
			}
			return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      args[0],
				"src_type": "nil or invalid time",
				"dst_type": "time",
			})
		case *pq.NullTime:
			if nt != nil && nt.Valid {
				return nt.Time, nil
			}
			return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      args[0],
				"src_type": "nil or invalid time",
				"dst_type": "time",
			})
		}

		// Handle []byte by converting to string
		if timeVal, ok := args[0].([]byte); ok {
			return Time(string(timeVal))
		}

		// Handle other types based on reflection
		elemKind := reflect.TypeOf(args[0]).Kind()
		switch elemKind {
		case reflect.Slice:
			return handleSliceForTime(args[0], elemKind)
		case reflect.String:
			return handleStringForTime(args[0], elemKind)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return handleIntForTime(args[0], elemKind)
		case reflect.Float32, reflect.Float64:
			return handleFloatForTime(args[0], elemKind)
		default:
			return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
				"src":      args[0],
				"src_type": elemKind.String(),
				"dst_type": "time",
			})
		}
	}

	// Handle multiple arguments
	if len(args) == 2 {
		return handleTwoArgsForTime(args)
	}
	if len(args) == 7 {
		return handleSevenArgsForTime(args)
	}
	if len(args) == 8 {
		return handleEightArgsForTime(args)
	}

	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
		"src":      args,
		"src_type": reflect.TypeOf(args).String(),
		"dst_type": "time",
	})
}

// handleSliceForTime handles slice arguments for Time conversion
func handleSliceForTime(arg any, elemKind reflect.Kind) (time.Time, zerror.Error) {
	result := time.Time{}
	sliceParam, zer := SliceAny(arg)
	if zer != nil {
		err := zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      arg,
			"src_type": elemKind.String(),
			"dst_type": "time",
		})
		err.Add(zer.GetList())
		return result, err
	}
	return Time(sliceParam...)
}

// handleStringForTime handles string arguments for Time conversion
func handleStringForTime(arg any, elemKind reflect.Kind) (time.Time, zerror.Error) {
	result := time.Time{}
	timeStr, zer := String(arg)
	if zer != nil {
		err := zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      arg,
			"src_type": elemKind.String(),
			"dst_type": "time",
		})
		err.Add(zer.GetList())
		return result, err
	}

	// Check if it's a unix timestamp (numeric string)
	if IsNumber(timeStr) {
		if strings.Index(timeStr, ".") >= 0 { // unixnano time representation
			timeFragments := strings.Split(timeStr, ".")
			secStr := timeFragments[0]
			sec, er := Int64(secStr)
			if er != nil {
				return result, er
			}
			nsecStr := timeFragments[1]
			nsec, er := Int64(nsecStr)
			if er != nil {
				return result, er
			}
			nsec *= 1e9
			result = time.Unix(sec, nsec)
			return result, nil
		} else { // unix time representation
			sec, er := Int64(timeStr)
			if er != nil {
				return result, er
			}
			result = time.Unix(sec, 0)
			return result, nil
		}
	}

	// Try all layouts from TimeLayoutMap
	for _, layout := range TimeLayoutMap {
		result, er := time.Parse(layout, timeStr)
		if er == nil {
			return result, nil
		}
	}

	return result, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
		"src":      timeStr,
		"src_type": elemKind.String(),
		"dst_type": "time",
	})
}

// handleIntForTime handles integer arguments for Time conversion
func handleIntForTime(arg any, elemKind reflect.Kind) (time.Time, zerror.Error) {
	result := time.Time{}
	unixTime, zer := Int64(arg)
	if zer != nil {
		err := zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      arg,
			"src_type": elemKind.String(),
			"dst_type": "time",
		})
		err.Add(zer.GetList())
		return result, err
	}
	result = time.Unix(unixTime, 0)
	return result, nil
}

// handleFloatForTime handles float arguments for Time conversion
func handleFloatForTime(arg any, elemKind reflect.Kind) (time.Time, zerror.Error) {
	result := time.Time{}
	floatTime, zer := Float64(arg)
	if zer != nil {
		err := zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      arg,
			"src_type": elemKind.String(),
			"dst_type": "time",
		})
		err.Add(zer.GetList())
		return result, err
	}
	unixTime, zer := Int64(arg)
	if zer != nil {
		err := zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      arg,
			"src_type": elemKind.String(),
			"dst_type": "time",
		})
		err.Add(zer.GetList())
		return result, err
	}
	unixNano := int64((floatTime - float64(unixTime)) * 1e9)
	result = time.Unix(unixTime, unixNano)
	return result, nil
}

// handleTwoArgsForTime handles two arguments for Time conversion (unixTime, unixNano)
func handleTwoArgsForTime(args []any) (time.Time, zerror.Error) {
	result := time.Time{}
	unixTime, err := Int64(args[0])
	if err != nil {
		return result, err
	}
	unixNano, err := Int64(args[1])
	if err != nil {
		return result, err
	}
	result = time.Unix(unixTime, unixNano)
	return result, nil
}

// handleSevenArgsForTime handles seven arguments for Time conversion (year, month, day, hour, min, sec, nsec)
func handleSevenArgsForTime(args []any) (time.Time, zerror.Error) {
	result := time.Time{}
	year, err := Int(args[0])
	if err != nil {
		return result, err
	}
	month, err := Int(args[1])
	if err != nil {
		return result, err
	}
	day, err := Int(args[2])
	if err != nil {
		return result, err
	}
	hour, err := Int(args[3])
	if err != nil {
		return result, err
	}
	min, err := Int(args[4])
	if err != nil {
		return result, err
	}
	sec, err := Int(args[5])
	if err != nil {
		return result, err
	}
	nsec, err := Int(args[6])
	if err != nil {
		return result, err
	}
	result = time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
	return result, nil
}

// handleEightArgsForTime handles eight arguments for Time conversion (year, month, day, hour, min, sec, nsec, location)
func handleEightArgsForTime(args []any) (time.Time, zerror.Error) {
	result := time.Time{}
	year, err := Int(args[0])
	if err != nil {
		return result, err
	}
	month, err := Int(args[1])
	if err != nil {
		return result, err
	}
	day, err := Int(args[2])
	if err != nil {
		return result, err
	}
	hour, err := Int(args[3])
	if err != nil {
		return result, err
	}
	min, err := Int(args[4])
	if err != nil {
		return result, err
	}
	sec, err := Int(args[5])
	if err != nil {
		return result, err
	}
	nsec, err := Int(args[6])
	if err != nil {
		return result, err
	}

	locationItf := args[7]
	location, ok := locationItf.(*time.Location)
	if !ok {
		location = time.UTC
	}
	result = time.Date(year, time.Month(month), day, hour, min, sec, nsec, location)
	return result, nil
}

// TimeFromFormattedString - parses a string to time.Time using a specified format
// params:
//
//	format - the time format to use for parsing (can be human-readable format that GetTimeLayout supports)
//	value  - the string value to parse
//
// Returns an error if the value cannot be parsed with the specified format.
// For parsing unix timestamps or other numeric/date component inputs, use the Time function instead.
func TimeFromFormattedString(format string, value string) (dst time.Time, err zerror.Error) {
	if value == "" {
		return time.Time{}, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      value,
			"src_type": "string",
			"dst_type": "time",
			"format":   format,
		})
	}

	layout := GetTimeLayout(format)
	result, er := time.Parse(layout, value)
	if er != nil {
		return time.Time{}, zerror.New(ErrorConvertorTypeNotSupported, map[string]any{
			"src":      value,
			"src_type": "string",
			"dst_type": "time",
			"format":   format,
			"layout":   layout,
			"error":    er.Error(),
		})
	}

	return result, nil
}

// TimeToFormattedString - formats a time.Time into a string using a specified format
// params:
//
//	format - the time format to use for formatting (can be human-readable format that GetTimeLayout supports)
//	value  - the time.Time value to format
//
// This is the inverse of TimeFromFormattedString.
func TimeToFormattedString(format string, value time.Time) (dst string) {
	layout := GetTimeLayout(format)
	return value.Format(layout)
}

// DateUTC parses a date-only value and normalizes it to 12:00:00 UTC.
// When format is provided the value is parsed with TimeFromFormattedString,
// otherwise the flexible Time parser is used. The calendar date is taken from the
// parsed value's own location (not converted to UTC first) and fixed to noon UTC,
// keeping the calendar date stable across timezone conversions.
func DateUTC(value any, format string) (dst time.Time, err zerror.Error) {
	var t time.Time
	if format != "" {
		var s string
		s, err = String(value)
		if err != nil {
			return time.Time{}, err
		}
		t, err = TimeFromFormattedString(format, s)
	} else {
		t, err = Time(value)
	}
	if err != nil {
		return time.Time{}, err
	}

	year, month, day := t.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, time.UTC), nil
}

// GetTimeLayout returns a go time layout from a human readable format.
// Supports common date formats from the import field configuration.
//
// Supported formats:
//   - Date only: mm/dd/yy, mm/dd/yyyy, dd/mm/yy, dd/mm/yyyy, yyyy-mm-dd, yyyy/mm/dd,
//     mm-dd-yy, mm-dd-yyyy, dd-mm-yy, dd-mm-yyyy
//   - ISO 8601: yyyy-mm-ddThh:mm:ssZ, yyyy-mm-ddThh:mm:ss.sssZ, yyyy-mm-ddThh:mm:ss+hh:mm,
//     yyyy-mm-ddThh:mm:ss-hh:mm
//   - RFC formats: RFC1123, RFC1123Z, RFC822, RFC822Z, RFC850, RFC3339, RFC3339Nano
//   - Time only: hh:mm:ss, hh:mm, hh:mm:ss.sss, hh:mm:ss.ssssss, hh:mm:ss.sssssssss
//   - Kitchen: hh:mmPM, hh:mm AM/PM
//   - Stamp: Jan _2 15:04:05, Jan _2 15:04:05.000, Jan _2 15:04:05.000000,
//     Jan _2 15:04:05.000000000
//
// If the format is not recognized, it is returned as-is.
func GetTimeLayout(format string) string {
	if layout, ok := TimeLayoutMap[format]; ok {
		return layout
	}

	// Return the format as-is if not recognized
	return format
}

// ToNullTime - tries to convert any to a NullTime using the same rules as Time.
// A nil input returns a null (invalid) NullTime with no error. On a successful
// conversion it returns a valid NullTime. If the conversion fails, it returns an
// invalid (null) NullTime together with the conversion error.
func ToNullTime(args ...any) (dst NullTime, err zerror.Error) {
	// A nil input maps to a null NullTime rather than a valid zero time.
	if len(args) == 1 && args[0] == nil {
		return NullTime{}, nil
	}

	nullTime := NullTime{}
	t, err := Time(args...)
	if err != nil {
		nullTime.Valid = false
		return nullTime, err
	}

	dst.Time = t
	dst.Valid = true

	return dst, nil
}
