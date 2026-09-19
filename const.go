package zgen

import "time"

// TimeLayoutMap maps human-readable time format strings to Go time layout constants.
// This map is exported for use in other packages that need to convert between
// human-readable formats and Go's time layout format.
//
// Supported format categories:
//   - Date only: mm/dd/yy, mm/dd/yyyy, dd/mm/yy, dd/mm/yyyy, yyyy-mm-dd, yyyy/mm/dd,
//     mm-dd-yy, mm-dd-yyyy, dd-mm-yy, dd-mm-yyyy, yyyymmdd, yy-mm-dd, yy/mm/dd
//   - ISO 8601: yyyy-mm-ddThh:mm:ssZ, yyyy-mm-ddThh:mm:ss.sssZ, yyyy-mm-ddThh:mm:ss+hh:mm,
//     yyyy-mm-ddThh:mm:ss-hh:mm, yyyy-mm-ddThh:mm:ssZ07:00, yyyy-mm-ddThh:mm:ss.sssZ07:00
//   - RFC formats: RFC1123, RFC1123Z, RFC822, RFC822Z, RFC850, RFC3339, RFC3339Nano
//   - Time only: hh:mm:ss, hh:mm, hh:mm:ss.sss, hh:mm:ss.ssssss, hh:mm:ss.sssssssss
//   - Kitchen: hh:mmPM, hh:mm AM/PM
//   - Stamp: Jan _2 15:04:05, Jan _2 15:04:05.000, Jan _2 15:04:05.000000,
//     Jan _2 15:04:05.000000000
//   - Date and time combined: yyyy-mm-dd hh:mm:ss, yyyy-mm-dd hh:mm:ss Z, yyyy-mm-dd hh:mm:ss.sss,
//     mm/dd/yyyy hh:mm:ss, dd/mm/yyyy hh:mm:ss
var TimeLayoutMap = map[string]string{
	// Date only formats
	"mm/dd/yy":     "01/02/06",
	"mm/dd/yyyy":   "01/02/2006",
	"dd/mm/yy":     "02/01/06",
	"dd/mm/yyyy":   "02/01/2006",
	"yyyy-mm-dd":   time.DateOnly,
	"yyyy/mm/dd":   "2006/01/02",
	"mm-dd-yy":     "01-02-06",
	"mm-dd-yyyy":   "01-02-2006",
	"dd-mm-yy":     "02-01-06",
	"dd-mm-yyyy":   "02-01-2006",
	"yyyymmdd":     "20060102",
	"yy-mm-dd":     "06-01-02",
	"yy/mm/dd":     "06/01/02",

	// ISO 8601 formats
	"yyyy-mm-ddThh:mm:ssZ":       time.RFC3339,
	"yyyy-mm-ddThh:mm:ss.sssZ":   "2006-01-02T15:04:05.000Z",
	"yyyy-mm-ddThh:mm:ss+hh:mm":  "2006-01-02T15:04:05-07:00",
	"yyyy-mm-ddThh:mm:ss-hh:mm":  "2006-01-02T15:04:05-07:00",
	"yyyy-mm-ddThh:mm:ssZ07:00":  time.RFC3339,
	"yyyy-mm-ddThh:mm:ss.sssZ07:00": time.RFC3339Nano,

	// RFC formats
	"ddd, dd mmm yyyy hh:mm:ss GMT":    time.RFC1123,
	"ddd, dd mmm yyyy hh:mm:ss +0000":  time.RFC1123Z,
	"ddd, dd mmm yy hh:mm:ss -0700":    time.RFC822,
	"ddd, dd mmm yy hh:mm:ss -0000":    time.RFC822Z,
	"dddd, dd-mmm-yy hh:mm:ss GMT":     time.RFC850,
	"rfc1123":   time.RFC1123,
	"rfc1123z":  time.RFC1123Z,
	"rfc822":    time.RFC822,
	"rfc822z":   time.RFC822Z,
	"rfc850":    time.RFC850,
	"rfc3339":   time.RFC3339,
	"rfc3339nano": time.RFC3339Nano,

	// Time only formats
	"hh:mm:ss":        "15:04:05",
	"hh:mm":           "15:04",
	"hh:mm:ss.sss":    "15:04:05.000",
	"hh:mm:ss.ssssss": "15:04:05.000000",
	"hh:mm:ss.sssssssss": "15:04:05.000000000",

	// Kitchen format
	"hh:mmPM":         "3:04PM",
	"hh:mm am/pm":     "3:04PM",

	// Stamp formats
	"Jan _2 15:04:05":            time.Stamp,
	"Jan _2 15:04:05.000":        time.StampMilli,
	"Jan _2 15:04:05.000000":     time.StampMicro,
	"Jan _2 15:04:05.000000000":  time.StampNano,

	// Date and time combined
	"yyyy-mm-dd hh:mm:ss Z":        TimeFormatISOSTZ,
	"yyyy-mm-dd hh:mm:ss":        time.DateTime,
	"yyyy-mm-dd hh:mm:ss.sss":    "2006-01-02 15:04:05.000",
	"mm/dd/yyyy hh:mm:ss":        "01/02/2006 15:04:05",
	"dd/mm/yyyy hh:mm:ss":        "02/01/2006 15:04:05",
}
