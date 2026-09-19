package zgen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/znxlc/zerror"
)

func TestUnit_Time(t *testing.T) {
	// Base time and its components
	currTime := time.Now().UTC()
	year, month, day := currTime.Date()
	hour, minute, sec := currTime.Clock()
	nsec := currTime.Nanosecond()
	unixTime := currTime.Unix()
	floatTime := float64(unixTime) + float64(nsec)/1e9

	// Time formats
	timeRFC1123 := currTime.Format(time.RFC1123)
	timeRFC850 := currTime.Format(time.RFC850)
	timeRFC1123Z := currTime.Format(time.RFC1123Z)
	timeRFC3339 := currTime.Format(time.RFC3339)
	timeRFC3339Nano := currTime.Format(time.RFC3339Nano)
	timeIsoDate := currTime.Format(TimeFormatISODate)
	timeIsoDateTime := currTime.Format(TimeFormatISO)
	timeIsoDateTimeSTZ := currTime.Format(TimeFormatISOSTZ)
	timeIsoDateTimeTZ := currTime.Format(TimeFormatISOTZ)

	// Test cases
	tests := []struct {
		name     string
		input    any
		expected time.Time
		hasError bool
	}{
		// Basic cases
		{name: "nil input", input: nil, expected: time.Time{}},
		{name: "time.Time input", input: currTime, expected: currTime},

		// String formats
		{name: "RFC1123 format", input: timeRFC1123, expected: currTime.Truncate(time.Second)},
		{name: "RFC850 format", input: timeRFC850, expected: currTime.Truncate(time.Second)},
		{name: "RFC1123Z format", input: timeRFC1123Z, expected: currTime.Truncate(time.Second)},
		{name: "RFC3339 format", input: timeRFC3339, expected: currTime.Truncate(time.Second)},
		{name: "RFC3339Nano format", input: timeRFC3339Nano, expected: currTime},
		{name: "ISO date format", input: timeIsoDate, expected: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)},
		{name: "ISO datetime format", input: timeIsoDateTime, expected: time.Date(year, month, day, hour, minute, sec, 0, time.UTC)},
		{name: "ISO datetime with space timezone", input: timeIsoDateTimeSTZ, expected: currTime.Truncate(time.Second)},
		{name: "ISO datetime with Z timezone", input: timeIsoDateTimeTZ, expected: currTime.Truncate(time.Second)},
		{name: "string unix timestamp (seconds)", input: "1672531200", expected: time.Unix(1672531200, 0).UTC()},
		{name: "string unix timestamp as float64", input: "1672531200.1", expected: time.Unix(1672531200, 1e9).UTC()},

		// Numeric inputs
		{name: "unix timestamp (seconds)", input: unixTime, expected: time.Unix(unixTime, 0).UTC()},
		{name: "unix timestamp with nanoseconds", input: []any{unixTime, nsec}, expected: time.Unix(unixTime, int64(nsec)).UTC()},
		{name: "unix timestamp as float64", input: floatTime, expected: time.Unix(int64(floatTime), int64((floatTime-float64(int64(floatTime)))*1e9)).UTC()},

		// Component inputs
		{name: "7 date components", input: []any{year, int(month), day, hour, minute, sec, nsec}, expected: currTime},
		{name: "8 date components", input: []any{year, int(month), day, hour, minute, sec, nsec, time.UTC}, expected: currTime},
		{name: "time.Duration", input: time.Duration(nsec), expected: time.Unix(0, int64(nsec)).UTC()},

		// Edge cases
		{name: "empty string", input: "", hasError: true},
		{name: "invalid date string", input: "not a date", hasError: true},
		{name: "unsupported type", input: map[string]string{"key": "value"}, hasError: true},
		{name: "too many components", input: []any{1, 2, 3, 4, 5, 6, 7, 8, 9}, hasError: true},
		{name: "extra month", input: []any{2023, 13, 1, 0, 0, 0, 0}, expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},

		// []byte input
		{name: "[]byte input", input: []byte(timeRFC3339), expected: currTime.Truncate(time.Second)},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result time.Time
			var zerr zerror.Error

			// Handle different input types
			switch v := tt.input.(type) {
			case []any:
				result, zerr = Time(v...)
			case nil:
				result, zerr = Time(nil)
			default:
				result, zerr = Time(v)
			}

			if tt.hasError {
				if zerr == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if zerr != nil {
				t.Errorf("Unexpected error: %v", zerr)
				return
			}

			// For time comparisons, handle different precisions
			if !result.Equal(tt.expected) {
				// If nanosecond precision differs but seconds are the same, it might be a formatting issue
				if result.Unix() != tt.expected.Unix() ||
					(result.Nanosecond()/1000) != (tt.expected.Nanosecond()/1000) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}

	// Additional test cases for specific scenarios
	t.Run("with time.Duration", func(t *testing.T) {
		d := time.Hour*2 + time.Minute*30
		res, zerr := Time(d)
		if zerr != nil {
			t.Fatalf("Unexpected error: %v", zerr)
		}
		expected := time.Unix(0, d.Nanoseconds()).UTC()
		if !res.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, res)
		}
	})

	t.Run("with invalid components", func(t *testing.T) {
		_, zerr := Time(2023, 2, 30) // Invalid date
		if zerr == nil {
			t.Error("Expected error for invalid date, got none")
		}
	})
}

func TestGetTimeLayout(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Date only formats
		{name: "mm/dd/yy", input: "mm/dd/yy", expected: "01/02/06"},
		{name: "mm/dd/yyyy", input: "mm/dd/yyyy", expected: "01/02/2006"},
		{name: "dd/mm/yy", input: "dd/mm/yy", expected: "02/01/06"},
		{name: "dd/mm/yyyy", input: "dd/mm/yyyy", expected: "02/01/2006"},
		{name: "yyyy-mm-dd", input: "yyyy-mm-dd", expected: "2006-01-02"},
		{name: "yyyy/mm/dd", input: "yyyy/mm/dd", expected: "2006/01/02"},
		{name: "mm-dd-yy", input: "mm-dd-yy", expected: "01-02-06"},
		{name: "mm-dd-yyyy", input: "mm-dd-yyyy", expected: "01-02-2006"},
		{name: "dd-mm-yy", input: "dd-mm-yy", expected: "02-01-06"},
		{name: "dd-mm-yyyy", input: "dd-mm-yyyy", expected: "02-01-2006"},
		{name: "yyyymmdd", input: "yyyymmdd", expected: "20060102"},
		{name: "yy-mm-dd", input: "yy-mm-dd", expected: "06-01-02"},
		{name: "yy/mm/dd", input: "yy/mm/dd", expected: "06/01/02"},

		// ISO 8601 formats
		{name: "yyyy-mm-ddThh:mm:ssZ", input: "yyyy-mm-ddThh:mm:ssZ", expected: time.RFC3339},
		{name: "yyyy-mm-ddThh:mm:ss.sssZ", input: "yyyy-mm-ddThh:mm:ss.sssZ", expected: "2006-01-02T15:04:05.000Z"},
		{name: "yyyy-mm-ddThh:mm:ss+hh:mm", input: "yyyy-mm-ddThh:mm:ss+hh:mm", expected: "2006-01-02T15:04:05-07:00"},
		{name: "yyyy-mm-ddThh:mm:ss-hh:mm", input: "yyyy-mm-ddThh:mm:ss-hh:mm", expected: "2006-01-02T15:04:05-07:00"},
		{name: "yyyy-mm-ddThh:mm:ssZ07:00", input: "yyyy-mm-ddThh:mm:ssZ07:00", expected: time.RFC3339},
		{name: "yyyy-mm-ddThh:mm:ss.sssZ07:00", input: "yyyy-mm-ddThh:mm:ss.sssZ07:00", expected: time.RFC3339Nano},

		// RFC formats
		{name: "RFC1123", input: "ddd, dd mmm yyyy hh:mm:ss GMT", expected: time.RFC1123},
		{name: "RFC1123Z", input: "ddd, dd mmm yyyy hh:mm:ss +0000", expected: time.RFC1123Z},
		{name: "RFC822", input: "ddd, dd mmm yy hh:mm:ss -0700", expected: time.RFC822},
		{name: "RFC822Z", input: "ddd, dd mmm yy hh:mm:ss -0000", expected: time.RFC822Z},
		{name: "RFC850", input: "dddd, dd-mmm-yy hh:mm:ss GMT", expected: time.RFC850},
		{name: "rfc1123", input: "rfc1123", expected: time.RFC1123},
		{name: "rfc1123z", input: "rfc1123z", expected: time.RFC1123Z},
		{name: "rfc822", input: "rfc822", expected: time.RFC822},
		{name: "rfc822z", input: "rfc822z", expected: time.RFC822Z},
		{name: "rfc850", input: "rfc850", expected: time.RFC850},
		{name: "rfc3339", input: "rfc3339", expected: time.RFC3339},
		{name: "rfc3339nano", input: "rfc3339nano", expected: time.RFC3339Nano},

		// Time only formats
		{name: "hh:mm:ss", input: "hh:mm:ss", expected: "15:04:05"},
		{name: "hh:mm", input: "hh:mm", expected: "15:04"},
		{name: "hh:mm:ss.sss", input: "hh:mm:ss.sss", expected: "15:04:05.000"},
		{name: "hh:mm:ss.ssssss", input: "hh:mm:ss.ssssss", expected: "15:04:05.000000"},
		{name: "hh:mm:ss.sssssssss", input: "hh:mm:ss.sssssssss", expected: "15:04:05.000000000"},

		// Kitchen format
		{name: "hh:mmPM", input: "hh:mmPM", expected: "3:04PM"},
		{name: "hh:mm am/pm", input: "hh:mm am/pm", expected: "3:04PM"},

		// Stamp formats
		{name: "Stamp", input: "Jan _2 15:04:05", expected: time.Stamp},
		{name: "StampMilli", input: "Jan _2 15:04:05.000", expected: time.StampMilli},
		{name: "StampMicro", input: "Jan _2 15:04:05.000000", expected: time.StampMicro},
		{name: "StampNano", input: "Jan _2 15:04:05.000000000", expected: time.StampNano},

		// Date and time combined
		{name: "yyyy-mm-dd hh:mm:ss", input: "yyyy-mm-dd hh:mm:ss", expected: "2006-01-02 15:04:05"},
		{name: "yyyy-mm-dd hh:mm:ss.sss", input: "yyyy-mm-dd hh:mm:ss.sss", expected: "2006-01-02 15:04:05.000"},
		{name: "mm/dd/yyyy hh:mm:ss", input: "mm/dd/yyyy hh:mm:ss", expected: "01/02/2006 15:04:05"},
		{name: "dd/mm/yyyy hh:mm:ss", input: "dd/mm/yyyy hh:mm:ss", expected: "02/01/2006 15:04:05"},

		// Unknown format - should return as-is
		{name: "unknown format", input: "unknown-format", expected: "unknown-format"},
		{name: "empty string", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTimeLayout(tt.input)
			assert.Equal(t, tt.expected, result, "GetTimeLayout(%q) = %q, want %q", tt.input, result, tt.expected)
		})
	}
}

func TestGetTimeLayout_ActualParsing(t *testing.T) {
	// Test that the layouts actually work with time.Parse
	tests := []struct {
		name   string
		format string
		input  string
	}{
		{name: "mm/dd/yyyy", format: "mm/dd/yyyy", input: "01/15/2024"},
		{name: "dd/mm/yyyy", format: "dd/mm/yyyy", input: "15/01/2024"},
		{name: "yyyy-mm-dd", format: "yyyy-mm-dd", input: "2024-01-15"},
		{name: "yyyy-mm-dd hh:mm:ss", format: "yyyy-mm-dd hh:mm:ss", input: "2024-01-15 14:30:45"},
		{name: "hh:mm:ss", format: "hh:mm:ss", input: "14:30:45"},
		{name: "rfc3339", format: "rfc3339", input: "2024-01-15T14:30:45Z"},
		{name: "rfc1123", format: "rfc1123", input: "Mon, 15 Jan 2024 14:30:45 GMT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := GetTimeLayout(tt.format)
			_, err := time.Parse(layout, tt.input)
			assert.NoError(t, err, "time.Parse with layout %q and input %q should not error", layout, tt.input)
		})
	}
}

func TestUnit_TimeFromFormattedString(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		format   string
		input    string
		expected time.Time
		hasError bool
	}{
		// Date only formats
		{name: "yyyy-mm-dd format", format: "yyyy-mm-dd", input: "2024-01-15", expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{name: "mm/dd/yyyy format", format: "mm/dd/yyyy", input: "01/15/2024", expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{name: "dd/mm/yyyy format", format: "dd/mm/yyyy", input: "15/01/2024", expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{name: "yyyymmdd format", format: "yyyymmdd", input: "20240115", expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{name: "mm-dd-yyyy format", format: "mm-dd-yyyy", input: "01-15-2024", expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{name: "dd-mm-yyyy format", format: "dd-mm-yyyy", input: "15-01-2024", expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{name: "yyyy/mm/dd format", format: "yyyy/mm/dd", input: "2024/01/15", expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},

		// Date and time formats
		{name: "yyyy-mm-dd hh:mm:ss format", format: "yyyy-mm-dd hh:mm:ss", input: "2024-01-15 14:30:45", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},
		{name: "yyyy-mm-dd hh:mm:ss.sss format", format: "yyyy-mm-dd hh:mm:ss.sss", input: "2024-01-15 14:30:45.123", expected: time.Date(2024, 1, 15, 14, 30, 45, 123000000, time.UTC)},
		{name: "mm/dd/yyyy hh:mm:ss format", format: "mm/dd/yyyy hh:mm:ss", input: "01/15/2024 14:30:45", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},
		{name: "dd/mm/yyyy hh:mm:ss format", format: "dd/mm/yyyy hh:mm:ss", input: "15/01/2024 14:30:45", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},

		// RFC formats
		{name: "rfc3339 format", format: "rfc3339", input: "2024-01-15T14:30:45Z", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},
		{name: "rfc1123 format", format: "rfc1123", input: "Mon, 15 Jan 2024 14:30:45 GMT", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},
		{name: "rfc1123z format", format: "rfc1123z", input: "Mon, 15 Jan 2024 14:30:45 +0000", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},
		{name: "rfc850 format", format: "rfc850", input: "Monday, 15-Jan-24 14:30:45 GMT", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},

		// Time only formats
		{name: "hh:mm:ss format", format: "hh:mm:ss", input: "14:30:45", expected: time.Date(0, 1, 1, 14, 30, 45, 0, time.UTC)},
		{name: "hh:mm format", format: "hh:mm", input: "14:30", expected: time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC)},

		// ISO 8601 formats
		{name: "yyyy-mm-ddThh:mm:ssZ format", format: "yyyy-mm-ddThh:mm:ssZ", input: "2024-01-15T14:30:45Z", expected: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)},
		{name: "yyyy-mm-ddThh:mm:ss.sssZ format", format: "yyyy-mm-ddThh:mm:ss.sssZ", input: "2024-01-15T14:30:45.123Z", expected: time.Date(2024, 1, 15, 14, 30, 45, 123000000, time.UTC)},

		// Edge cases
		{name: "empty string", format: "yyyy-mm-dd", input: "", hasError: true},
		{name: "invalid date string for format", format: "yyyy-mm-dd", input: "not a date", hasError: true},
		{name: "wrong format for input", format: "mm/dd/yyyy", input: "2024-01-15", hasError: true},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimeFromFormattedString(tt.format, tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// For time comparisons, handle different precisions
			if !result.Equal(tt.expected) {
				// If nanosecond precision differs but seconds are the same, it might be a formatting issue
				if result.Unix() != tt.expected.Unix() ||
					(result.Nanosecond()/1000) != (tt.expected.Nanosecond()/1000) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnit_TimeToFormattedString(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		input    time.Time
		expected string
	}{
		{name: "yyyy-mm-dd format", format: "yyyy-mm-dd", input: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), expected: "2024-01-15"},
		{name: "mm/dd/yyyy format", format: "mm/dd/yyyy", input: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), expected: "01/15/2024"},
		{name: "yyyy-mm-dd hh:mm:ss format", format: "yyyy-mm-dd hh:mm:ss", input: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC), expected: "2024-01-15 14:30:45"},
		{name: "rfc3339 format", format: "rfc3339", input: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC), expected: "2024-01-15T14:30:45Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TimeToFormattedString(tt.format, tt.input))
		})
	}
}

func TestUnit_DateUTC(t *testing.T) {
	noon := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, 12, 0, 0, 0, time.UTC)
	}

	tests := []struct {
		name     string
		value    any
		format   string
		expected time.Time
		hasError bool
	}{
		{name: "date-only string, no format", value: "2025-01-01", expected: noon(2025, 1, 1)},
		{name: "datetime string strips time to noon", value: "2025-01-01 14:30:45", expected: noon(2025, 1, 1)},
		{name: "rfc3339 string, no format", value: "2025-01-01T14:30:45Z", expected: noon(2025, 1, 1)},
		{name: "time.Time input, no format", value: time.Date(2025, 1, 1, 8, 0, 0, 0, time.UTC), expected: noon(2025, 1, 1)},
		{name: "with format dd/mm/yyyy", value: "01/02/2025", format: "dd/mm/yyyy", expected: noon(2025, 2, 1)},
		{name: "with format mm/dd/yyyy", value: "02/01/2025", format: "mm/dd/yyyy", expected: noon(2025, 2, 1)},
		{name: "empty string with format errors", value: "", format: "yyyy-mm-dd", hasError: true},
		{name: "value does not match format", value: "not-a-date", format: "yyyy-mm-dd", hasError: true},
		{name: "unsupported value errors", value: map[string]any{}, hasError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DateUTC(tt.value, tt.format)

			if tt.hasError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.True(t, tt.expected.Equal(result), "expected %v, got %v", tt.expected, result)
			assert.Equal(t, time.UTC, result.Location())
			assert.Equal(t, 12, result.Hour())
		})
	}
}

func TestUnit_DateUTC_TimezoneStable(t *testing.T) {
	// converting the noon-UTC result to a non-UTC zone must keep the calendar date
	result, err := DateUTC("2025-01-01", "")
	assert.Nil(t, err)

	loc, lerr := time.LoadLocation("America/New_York")
	assert.NoError(t, lerr)
	local := result.In(loc)
	assert.Equal(t, 2025, local.Year())
	assert.Equal(t, time.January, local.Month())
	assert.Equal(t, 1, local.Day())
}

func TestUnit_TimeFormattedString_RoundTrip(t *testing.T) {
	const format = "yyyy-mm-dd hh:mm:ss"
	original := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)

	formatted := TimeToFormattedString(format, original)
	parsed, err := TimeFromFormattedString(format, formatted)

	assert.Nil(t, err)
	assert.True(t, original.Equal(parsed))
}

func TestUnit_ToNullTime(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		src := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

		nt, err := ToNullTime(src)

		assert.Nil(t, err)
		assert.True(t, nt.Valid)
		assert.True(t, nt.Time.Equal(src))
	})

	t.Run("valid string", func(t *testing.T) {
		nt, err := ToNullTime("2020-01-02T03:04:05Z")

		assert.Nil(t, err)
		assert.True(t, nt.Valid)
		assert.Equal(t, "2020-01-02T03:04:05Z", nt.Time.UTC().Format(time.RFC3339))
	})

	t.Run("nil yields null", func(t *testing.T) {
		nt, err := ToNullTime(nil)

		assert.Nil(t, err)
		assert.False(t, nt.Valid)
	})

	t.Run("no args yields valid zero time", func(t *testing.T) {
		nt, err := ToNullTime()

		assert.Nil(t, err)
		assert.True(t, nt.Valid)
		assert.True(t, nt.Time.IsZero())
	})

	t.Run("invalid conversion yields null with error", func(t *testing.T) {
		nt, err := ToNullTime(map[string]any{})

		assert.NotNil(t, err)
		assert.False(t, nt.Valid)
	})
}
