package zgen

import (
  "flag"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "net/http"
  "os"
  "testing"
  "time"
)

var (
  configFolder = flag.String("config", "../../../config", "Path to config folder")
)

func TestMain(m *testing.M) {
  flag.Parse()
  os.Exit(m.Run())
}

type SimpleStruct struct {
  Foo string `json:"foo"`
  Bar []int  `json:"bar"`
}

type DoublePointerStruct struct {
  Data      map[string]interface{} `json:"data"`
  ANumber   int                    `json:"a_number"`
  DoublePtr **SimpleStruct         `json:"double_ptr"`
}

type GOTestStruct struct {
  ID      int                `tst:"tst_id"`
  Key     string             `tst:"tst_key,omitempty"`
  Value   string             `tst:"tst_value,omitempty"`
  Str     GOTestNestedStruct `tst:"str"`
  NoNest  interface{}        `tst:"nonest,omitnested"`
  Nested  interface{}        `tst:"nested"`
  NB      NullBool           `tst:"nb"`
  NU      uint64             `tst:"nu"`
  NI      int64              `tst:"ni"`
  NF      float64            `tst:"nf,omitempty"`
  NS      string             `tst:"ns"`
  NT      NullTime           `tst:"nt"`
  RWField NullTime           // this field has no tag so we can test some of the scan features
}

type GOTestStruct2 struct {
  ID     int                           `tst:"tst_id"`
  Key    string                        `tst:"tst_key,omitempty"`
  Value  string                        `tst:"tst_value,omitempty"`
  Str    GOTestNestedStruct            `tst:"str"`
  NoNest interface{}                   `tst:"nonest,omitnested"`
  Nested interface{}                   `tst:"nested"`
  Map    map[string]GOTestNestedStruct `tst:"map"`
  Slice  []GOTestNestedStruct          `tst:"slice"`
  NB     NullBool                      `tst:"nb"`
  NU     uint64                        `tst:"nu"`
  NI     int64                         `tst:"ni"`
  NF     float64                       `tst:"nf,omitempty"`
  NS     string                        `tst:"ns"`
  NT     NullTime                      `tst:"nt"`
  Time   time.Time                     `tst:"time"`
  RWNT   NullTime                      // element without tag that will be shown only if NameReturnWhenNoTagsFound: true
}

type GOTestNestedStruct struct {
  ID    int        `tst:"tst_id"`
  Key   string     `tst:"tst_key,omitempty"`
  Value string     `tst:"tst_value"`
  Str   testStruct `tst:"tst_struct"`
}

type testStruct struct {
  KeyString string                 `tst:"key_string" json:"key_string"`
  KeyMap    map[string]interface{} `tst:"key_map" json:"key_map"`
  KeySlice  []interface{}          `tst:"key_slice" json:"key_slice"`
}

func (x testStruct) GetKeyString() string {
  return x.KeyString
}
func (x testStruct) NonGetKeyString(key string) string {
  return key
}

var ExpectedStruct = GOTestStruct{
  ID:    13,
  Key:   "GTSKey",
  Value: "GTSValue",
}

var testSample = map[string]interface{}{
  "stringKey": "stringValue",
  "intKey":    10,
  "sliceKey": []interface{}{
    "sliceValue1",
    map[string]interface{}{
      "slice2Key": "slice2Value",
      "slice2Int": 12,
      "slice2Ref": &GOTestStruct{
        ID:    13,
        Key:   "GTSKey",
        Value: "GTSValue",
      },
      "slice2Struct": testStruct{
        KeyString: "structString",
        KeyMap: map[string]interface{}{
          "slice2structMap1": "test",
        },
        KeySlice: []interface{}{1, 2, 3, 4, 5},
      },
    },
  },
  "GTS": GOTestStruct{
    ID:    13,
    Key:   "GTSKey",
    Value: "GTSValue"},
  "boolKey": true,
}

func TestUnit_ScanToStructSimple(t *testing.T) {
  destStruct := GOTestNestedStruct{}

  err := ToStruct(&destStruct,
    ParserConfig{Mode: ParserModeTagsOnly, Tags: []string{"tst"}},
    map[string]interface{}{
      "tst_id":    1,
      "tst_value": "value",
      "tst_struct": map[string]interface{}{
        "key_string": "key_string_value",
      },
    })
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, 1, destStruct.ID)
  assert.Equal(t, "", destStruct.Key)
  assert.Equal(t, "value", destStruct.Value)
  assert.Equal(t, "key_string_value", destStruct.Str.KeyString)

}

func TestUnit_ScanToElementTime(t *testing.T) {
  dest := time.Time{}

  src := time.Now().UTC()

  err := ScanToElement(&dest,
    src,
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src, dest)
  dest = time.Time{}
  // scanning to a struct from a string, using a format that has nanoseconds to make sure element remains the same
  err = ScanToElement(&dest,
    src.Format(time.RFC3339Nano),
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src, dest)

  var nilDest time.Time
  // scanning to a struct from a string, using a format that has nanoseconds to make sure element remains the same
  err = ScanToElement(&nilDest,
    src.Format(time.RFC3339Nano),
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src, nilDest)

  destPtr := &time.Time{}
  // scanning to a struct from a string, using a format that has nanoseconds to make sure element remains the same
  err = ScanToElement(&destPtr,
    src.Format(time.RFC3339Nano),
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src, *destPtr)

  var nilDestPtr *time.Time
  // scanning to a struct from a string, using a format that has nanoseconds to make sure element remains the same
  err = ScanToElement(&nilDestPtr,
    src.Format(time.RFC3339Nano),
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src, *nilDestPtr)
}

func TestUnit_ScanToElement_SameStructDifferentType(t *testing.T) {
  // we use elements of same struct but different name types
  type HTTPCookie struct {
    Name       string        `json:"name" mapstructure:"name" structs:"name"`
    Value      string        `json:"value" mapstructure:"value" structs:"value"`
    Path       string        `json:"path,omitempty" mapstructure:"path,omitempty" structs:"path"`
    Domain     string        `json:"domain,omitempty" mapstructure:"domain,omitempty" structs:"domain"`
    Expires    time.Time     `json:"expires,omitempty" mapstructure:"expires,omitempty" structs:"expires"`
    RawExpires string        `json:"raw_expires" mapstructure:"raw_expires" structs:"raw_expires"` // RawExpires is a read-only, string version of Expires in GMT timezone and formatted as RFC1123 ("cookie time").
    MaxAge     int           `json:"max_age" mapstructure:"max_age" structs:"max_age"`
    Secure     bool          `json:"secure" mapstructure:"secure" structs:"secure"`
    HTTPOnly   bool          `json:"http_only" mapstructure:"http_only" structs:"http_only"`
    SameSite   http.SameSite `json:"same_site" mapstructure:"same_site" structs:"same_site"`
    Raw        string        `json:"raw" mapstructure:"raw" structs:"raw"`
    Unparsed   []string      `json:"unparsed" mapstructure:"unparsed" structs:"unparsed"`
  }

  type HTTPCookie2 struct {
    Name       string        `json:"name" mapstructure:"name" structs:"name"`
    Value      string        `json:"value" mapstructure:"value" structs:"value"`
    Path       string        `json:"path,omitempty" mapstructure:"path,omitempty" structs:"path"`
    Domain     string        `json:"domain,omitempty" mapstructure:"domain,omitempty" structs:"domain"`
    Expires    time.Time     `json:"expires,omitempty" mapstructure:"expires,omitempty" structs:"expires"`
    RawExpires string        `json:"raw_expires" mapstructure:"raw_expires" structs:"raw_expires"` // RawExpires is a read-only, string version of Expires in GMT timezone and formatted as RFC1123 ("cookie time").
    MaxAge     int           `json:"max_age" mapstructure:"max_age" structs:"max_age"`
    Secure     bool          `json:"secure" mapstructure:"secure" structs:"secure"`
    HTTPOnly   bool          `json:"http_only" mapstructure:"http_only" structs:"http_only"`
    SameSite   http.SameSite `json:"same_site" mapstructure:"same_site" structs:"same_site"`
    Raw        string        `json:"raw" mapstructure:"raw" structs:"raw"`
    Unparsed   []string      `json:"unparsed" mapstructure:"unparsed" structs:"unparsed"`
  }

  dest := HTTPCookie{}

  src := HTTPCookie2{
    Name:       "tst",
    Value:      "cookieVal",
    Path:       "/",
    Domain:     "test.com",
    Expires:    time.Now().Add(300).UTC(), // expiration is now+300sec
    RawExpires: "xxx",
    MaxAge:     1000,
    Secure:     true,
    HTTPOnly:   true,
    SameSite:   10,
    Raw:        "tst",
    Unparsed:   nil,
  }

  err := ToStruct(&dest,
    src,
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src.Name, dest.Name)
  assert.Equal(t, src.Value, dest.Value)
  assert.Equal(t, src.Path, dest.Path)
  assert.Equal(t, src.Domain, dest.Domain)
  assert.Equal(t, src.Expires, dest.Expires)
  assert.Equal(t, src.RawExpires, dest.RawExpires)
  assert.Equal(t, src.MaxAge, dest.MaxAge)
  assert.Equal(t, src.Secure, dest.Secure)
  assert.Equal(t, src.HTTPOnly, dest.HTTPOnly)
  assert.Equal(t, src.SameSite, dest.SameSite)
  assert.Equal(t, src.Raw, dest.Raw)
  assert.Equal(t, src.Unparsed, dest.Unparsed)
}

func TestUnit_ScanToElement_StructFromMap(t *testing.T) {
  // elements are compatible (Same fields of compatible types) so this will work
  type HTTPCookie struct {
    Name       string        `json:"name" mapstructure:"name" structs:"name"`
    Value      string        `json:"value" mapstructure:"value" structs:"value"`
    Path       string        `json:"path,omitempty" mapstructure:"path,omitempty" structs:"path"`
    Domain     string        `json:"domain,omitempty" mapstructure:"domain,omitempty" structs:"domain"`
    Expires    time.Time     `json:"expires,omitempty" mapstructure:"expires,omitempty" structs:"expires"`
    RawExpires string        `json:"raw_expires" mapstructure:"raw_expires" structs:"raw_expires"` // RawExpires is a read-only, string version of Expires in GMT timezone and formatted as RFC1123 ("cookie time").
    MaxAge     int           `json:"max_age" mapstructure:"max_age" structs:"max_age"`
    Secure     bool          `json:"secure" mapstructure:"secure" structs:"secure"`
    HTTPOnly   bool          `json:"http_only" mapstructure:"http_only" structs:"http_only"`
    SameSite   http.SameSite `json:"same_site" mapstructure:"same_site" structs:"same_site"`
    Raw        string        `json:"raw" mapstructure:"raw" structs:"raw"`
    Unparsed   []string      `json:"unparsed" mapstructure:"unparsed" structs:"unparsed"`
  }

  dest := HTTPCookie{}

  src := map[string]interface{}{
    "name":        "tst",
    "value":       "cookieVal",
    "path":        "/",
    "domain":      "test.com",
    "expires":     time.Now().Add(300).UTC(), // expiration is now+300sec
    "raw_expires": "xxx",
    "max_age":     1000,
    "secure":      true,
    "http_only":   true,
    "same_site":   10,
    "raw":         "tst",
    "unparsed":    []string{"tst1", "tst2"},
  }

  err := ScanToElement(&dest,
    src,
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src["name"], dest.Name)
  assert.Equal(t, src["value"], dest.Value)
  assert.Equal(t, src["path"], dest.Path)
  assert.Equal(t, src["domain"], dest.Domain)
  assert.Equal(t, src["expires"], dest.Expires)
  assert.Equal(t, src["raw_expires"], dest.RawExpires)
  assert.Equal(t, src["max_age"], dest.MaxAge)
  assert.Equal(t, src["secure"], dest.Secure)
  assert.Equal(t, src["http_only"], dest.HTTPOnly)
  assert.Equal(t, http.SameSite(src["same_site"].(int)), dest.SameSite)
  assert.Equal(t, src["raw"], dest.Raw)
  assert.Equal(t, src["unparsed"], dest.Unparsed)
}

func TestUnit_ScanToStruct_StructFromMap(t *testing.T) {
  type HTTPCookie struct {
    Name       string        `json:"name" mapstructure:"name" structs:"name"`
    Value      string        `json:"value" mapstructure:"value" structs:"value"`
    Path       string        `json:"path,omitempty" mapstructure:"path,omitempty" structs:"path"`
    Domain     string        `json:"domain,omitempty" mapstructure:"domain,omitempty" structs:"domain"`
    Expires    time.Time     `json:"expires,omitempty" mapstructure:"expires,omitempty" structs:"expires"`
    RawExpires string        `json:"raw_expires" mapstructure:"raw_expires" structs:"raw_expires"` // RawExpires is a read-only, string version of Expires in GMT timezone and formatted as RFC1123 ("cookie time").
    MaxAge     int           `json:"max_age" mapstructure:"max_age" structs:"max_age"`
    Secure     bool          `json:"secure" mapstructure:"secure" structs:"secure"`
    HTTPOnly   bool          `json:"http_only" mapstructure:"http_only" structs:"http_only"`
    SameSite   http.SameSite `json:"same_site" mapstructure:"same_site" structs:"same_site"`
    Raw        string        `json:"raw" mapstructure:"raw" structs:"raw"`
    Unparsed   []string      `json:"unparsed" mapstructure:"unparsed" structs:"unparsed"`
  }

  dest := HTTPCookie{}

  src := map[string]interface{}{
    "name":        "tst",
    "value":       "cookieVal",
    "path":        "/",
    "domain":      "test.com",
    "expires":     time.Now().Add(300).UTC(), // expiration is now+300sec
    "raw_expires": "xxx",
    "max_age":     1000,
    "secure":      true,
    "http_only":   true,
    "same_site":   10,
    "raw":         "tst",
    "unparsed":    []string{"tst1", "tst2"},
  }

  err := ToStruct(&dest,
    src,
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src["name"], dest.Name)
  assert.Equal(t, src["value"], dest.Value)
  assert.Equal(t, src["path"], dest.Path)
  assert.Equal(t, src["domain"], dest.Domain)
  assert.Equal(t, src["expires"], dest.Expires)
  assert.Equal(t, src["raw_expires"], dest.RawExpires)
  assert.Equal(t, src["max_age"], dest.MaxAge)
  assert.Equal(t, src["secure"], dest.Secure)
  assert.Equal(t, src["http_only"], dest.HTTPOnly)
  assert.Equal(t, http.SameSite(src["same_site"].(int)), dest.SameSite)
  assert.Equal(t, src["raw"], dest.Raw)
  assert.Equal(t, src["unparsed"], dest.Unparsed)
}

// Testing to make sure the existing data is overwritten properly with the new data
func TestUnit_ScanToStruct_PrefilledStructFromMap(t *testing.T) {
  type HTTPCookie struct {
    Name       string        `json:"name" mapstructure:"name" structs:"name"`
    Value      string        `json:"value" mapstructure:"value" structs:"value"`
    Path       string        `json:"path,omitempty" mapstructure:"path,omitempty" structs:"path"`
    Domain     string        `json:"domain,omitempty" mapstructure:"domain,omitempty" structs:"domain"`
    Expires    time.Time     `json:"expires,omitempty" mapstructure:"expires,omitempty" structs:"expires"`
    RawExpires string        `json:"raw_expires" mapstructure:"raw_expires" structs:"raw_expires"` // RawExpires is a read-only, string version of Expires in GMT timezone and formatted as RFC1123 ("cookie time").
    MaxAge     int           `json:"max_age" mapstructure:"max_age" structs:"max_age"`
    Secure     bool          `json:"secure" mapstructure:"secure" structs:"secure"`
    HTTPOnly   bool          `json:"http_only" mapstructure:"http_only" structs:"http_only"`
    SameSite   http.SameSite `json:"same_site" mapstructure:"same_site" structs:"same_site"`
    Raw        string        `json:"raw" mapstructure:"raw" structs:"raw"`
    Unparsed   []string      `json:"unparsed" mapstructure:"unparsed" structs:"unparsed"`
  }

  dest := HTTPCookie{
    Name:     "initial",
    Value:    "initial",
    Expires:  time.Now(),
    MaxAge:   0,
    SameSite: 1,
    Secure:   false,
    Unparsed: []string{"initial"},
  }

  src := map[string]interface{}{
    "name":        "tst",
    "value":       "cookieVal",
    "path":        "/",
    "domain":      "test.com",
    "expires":     time.Now().Add(300).UTC(), // expiration is now+300sec
    "raw_expires": "xxx",
    "max_age":     1000,
    "secure":      true,
    "http_only":   true,
    "same_site":   10,
    "raw":         "tst",
    "unparsed":    []string{"tst1", "tst2"},
  }

  err := ToStruct(&dest,
    src,
  )
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, src["name"], dest.Name)
  assert.Equal(t, src["value"], dest.Value)
  assert.Equal(t, src["path"], dest.Path)
  assert.Equal(t, src["domain"], dest.Domain)
  assert.Equal(t, src["expires"], dest.Expires)
  assert.Equal(t, src["raw_expires"], dest.RawExpires)
  assert.Equal(t, src["max_age"], dest.MaxAge)
  assert.Equal(t, src["secure"], dest.Secure)
  assert.Equal(t, src["http_only"], dest.HTTPOnly)
  assert.Equal(t, http.SameSite(src["same_site"].(int)), dest.SameSite)
  assert.Equal(t, src["raw"], dest.Raw)
  assert.Equal(t, src["unparsed"], dest.Unparsed)

}

func TestUnit_ScanToStructValidationErrors(t *testing.T) {
  destStruct := GOTestStruct{}

  err := ToStruct(destStruct) // error destStruct not a pointer
  assert.Equal(t, "ERROR_ZGEN_SCANNER_DST_STRUCTURE_INVALID", err.Error())

  err = ToStruct(&destStruct, "x") // error empty data to add
  assert.Equal(t, "ERROR_ZGEN_SCANNER_ARGUMENT_INVALID", err.Error())
}

func TestUnit_ScanToStruct(t *testing.T) {
  destStruct := GOTestStruct{}
  currTime := time.Now()
  nt := NullTime{}
  _ = nt.Scan(currTime)

  err := ToStruct(&destStruct,
    ParserConfig{Mode: ParserModeTagsOnly, Tags: []string{"tst"}},
    map[string]interface{}{
      "key1": "value1",
      "key2": "value2",
      "nf":   "3.14",
      "ns":   []byte("testNS")},
    map[string]interface{}{
      "key3": 1,
      "key4": 2,
      "ni":   int(30),
      "nb":   true,
      "nu":   "13"},
    GOTestStruct{
      ID:    10,
      Key:   "GTSKey",
      Value: "GTSValue",
      NT:    nt,
    })
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, 10, destStruct.ID)
  assert.Equal(t, "GTSKey", destStruct.Key)
  assert.Equal(t, "GTSValue", destStruct.Value)
  assert.Equal(t, int64(30), destStruct.NI)
  assert.Equal(t, uint64(13), destStruct.NU)
  assert.Equal(t, float64(3.14), destStruct.NF)
  assert.Equal(t, true, destStruct.NB.Bool)
  assert.Equal(t, "testNS", destStruct.NS)
  assert.Equal(t, currTime, destStruct.NT.Time)
}

func TestUnit_ScanToStructMapsAndSlices(t *testing.T) {
  currTime := time.Now()
  nt := NullTime{}
  _ = nt.Scan(currTime)

  mapSI := map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
    "key3": 1,
    "key4": 2,
    "map": map[string]interface{}{
      "val": map[string]interface{}{
        "tst_id":  10,
        "tst_key": "mapKey",
        "tst_struct": map[string]interface{}{
          "key_map": map[string]interface{}{
            "m1": "1",
            "m2": "2",
          },
          "key_slice":  []interface{}(nil),
          "key_string": "aa",
        },
        "tst_value": "mapValue",
      },
    },
    "slice": []map[string]interface{}{ // keys by stuct field name
      {
        "ID":  15,
        "Key": "s1Key",
        "Str": map[string]interface{}{
          "KeyMap": map[string]interface{}{
            "11": "1",
            "21": "2",
          },
          "KeySlice":  []interface{}{"1", "2"},
          "KeyString": "zz",
        },
        "Value": "s1Value",
      },
      {
        "tst_id":  16,
        "tst_key": "s2Key",
        "tst_struct": map[string]interface{}{
          "key_map": map[string]interface{}{
            "12": "1",
            "22": "2",
          },
          "key_slice":  []interface{}(nil),
          "key_string": "zz2",
        },
        "tst_value": "s2Value",
      },
    },
    "nb":     interface{}(nil),
    "nested": []uint8{0x7b, 0x22, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x3a, 0x22, 0x43, 0x41, 0x22, 0x2c, 0x22, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x3a, 0x22, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x40, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x63, 0x6f, 0x6d, 0x22, 0x7d},
    "ni":     0,
    "nonest": interface{}(nil),
    "ns":     "",
    "nt":     nt,
    "nu":     0x0,
    "str": map[string]interface{}{
      "tst_id": 0,
      "tst_struct": map[string]interface{}{
        "key_map":    map[string]interface{}{},
        "key_slice":  []interface{}(nil),
        "key_string": "",
      },
      "tst_value": "",
    },
    "tst_id":    10,
    "tst_key":   "GTSKey",
    "tst_value": "GTSValue",
  }

  dstStruct := GOTestStruct2{}

  err := ToStruct(
    &dstStruct,
    ParserConfig{Mode: ParserModeTagsOnly, Tags: []string{"tst"}},
    mapSI)

  if err != nil {
    t.Error(err)
  }

  // testing more complex scanning
  assert.Equal(t, 2, len(dstStruct.Slice))
  assert.Equal(t, 15, dstStruct.Slice[0].ID)
  assert.Equal(t, 16, dstStruct.Slice[1].ID)
  assert.Equal(t, "1", dstStruct.Slice[0].Str.KeySlice[0])
  assert.Equal(t, mapSI["nested"], dstStruct.Nested)
  assert.Equal(t, "GTSKey", dstStruct.Key)
  assert.Equal(t, uint64(0), dstStruct.NU)
  assert.Equal(t, currTime, dstStruct.NT.Time)
  assert.Equal(t, true, dstStruct.NT.Valid)
  assert.Equal(t, 10, dstStruct.Map["val"].ID)
}

func TestUnit_ScanToMapWithEvaluate(t *testing.T) {
  mapData := map[string]interface{}{}
  currTime := time.Now()
  nt := NullTime{}
  _ = nt.Scan(currTime)

  err := ToMap(&mapData,
    ParserConfig{Mode: ParserModeNameAndTags, Tags: []string{"tst"}, EvaluateMethods: true},
    map[string]interface{}{
      "key1": "value1",
      "key2": "value2",
    },
    &map[string]int{
      "key3": 1,
      "key4": 2,
    },
    GOTestStruct{
      ID:      10,
      Key:     "GTSKey",
      Value:   "GTSValue",
      NT:      nt,
      RWField: nt,
    })
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, "value1", mapData["key1"])
  assert.Equal(t, "value2", mapData["key2"])
  assert.Equal(t, 1, mapData["key3"])
  assert.Equal(t, 2, mapData["key4"])
  assert.Equal(t, 10, mapData["ID"])
  assert.Equal(t, 10, mapData["tst_id"])
  assert.Equal(t, "GTSKey", mapData["Key"])
  assert.Equal(t, "GTSValue", mapData["Value"])
  assert.Equal(t, currTime, mapData["nt"])
  assert.Equal(t, currTime, mapData["RWField"])
}

func TestUnit_ScanToMap_UnpackToInterface_DoublePointerToStruct(t *testing.T) {
  q := &SimpleStruct{
    Foo: "This space for rent.",
    Bar: []int{1, 2, 3},
  }
  b := DoublePointerStruct{
    Data: map[string]interface{}{
      "foo": "bar",
    },
    ANumber:   1,
    DoublePtr: &q,
  }
  expected := map[string]interface{}{
    "data": map[string]interface{}{
      "foo": "bar",
    },
    "a_number":   1,
    "double_ptr": &q,
  }
  actual := make(map[string]interface{})
  err := ToMap(&actual, b)
  require.Nil(t, err)
  require.Equal(t, expected, actual)

  // Scan the map back into the struct
  scanned := DoublePointerStruct{}
  err = ToStruct(&scanned, actual)
  require.Nil(t, err)
  require.Equal(t, b, scanned)
}

func TestUnit_ScanToMap(t *testing.T) {
  mapData := map[string]interface{}{}

  currTime := time.Now()
  nt := NullTime{}
  _ = nt.Scan(currTime)

  nb := NullBool{}

  NoNestStruct := GOTestNestedStruct{
    ID:    1,
    Key:   "somekey",
    Value: "someValue",
  }

  NestedStruct := GOTestNestedStruct{
    ID:    3,
    Key:   "slicekey",
    Value: "sliceValue",
  }

  config := DefaultParserConfig
  config.Tags = []string{"tst"}
  err := ToMap(&mapData,
    GOTestStruct{
      ID:    10,
      Key:   "GTSKey",
      Value: "GTSValue",
      NB:    nb,
      NoNest: map[string]interface{}{
        "1":   "1",
        "str": NoNestStruct,
      },
      Nested: []interface{}{
        NestedStruct,
        "nestedString",
      },
      Str:     NoNestStruct,
      RWField: nt, // this will be ignored since it has no tag
    }, config)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, 10, mapData["tst_id"])
  assert.Equal(t, nb, mapData["nb"])
  assert.Equal(t, nil, mapData["RWField"])
  nestedSlice, err := SliceAny(mapData["nested"])
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, 3, nestedSlice[0].(map[string]interface{})["tst_id"]) // nested slice successfully scanned

  NoNMap, err := MapStringAny(mapData["nonest"])
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, NoNestStruct, NoNMap["str"]) // NonNested struct
}

func TestUnit_ScanToMapByTagWithValuer(t *testing.T) {
  mapData := map[string]interface{}{}

  currTime := time.Now()
  nt := NullTime{}
  _ = nt.Scan(currTime)

  nb := NullBool{}

  NoNestSruct := GOTestNestedStruct{
    ID:    1,
    Key:   "somekey",
    Value: "someValue",
  }

  NestedStruct := GOTestNestedStruct{
    ID:    3,
    Key:   "slicekey",
    Value: "sliceValue",
  }

  config := DefaultParserConfig
  config.Mode = ParserModeNameAndTags
  config.Tags = []string{"tst"}
  config.EvaluateMethods = true

  err := ToMap(&mapData, config,
    GOTestStruct{
      ID:    10,
      Key:   "GTSKey",
      Value: "GTSValue",
      NB:    nb, // will be evaluated to nil
      NoNest: map[string]interface{}{
        "1":   "1",
        "str": NoNestSruct,
      },
      Nested: []interface{}{
        NestedStruct,
        "nestedString",
      },
      Str:     NoNestSruct,
      RWField: nt, // will be evaluated to currTime
    })
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, 10, mapData["tst_id"])
  assert.Equal(t, nil, mapData["nb"])
  assert.Equal(t, currTime, mapData["RWField"])
  nestedSlice, err := SliceAny(mapData["nested"])
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, 3, nestedSlice[0].(map[string]interface{})["tst_id"]) // nested slice successfully scanned

  NoNMap, err := MapStringAny(mapData["nonest"])
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, NoNestSruct, NoNMap["str"]) // NonNested struct
}

func TestUnit_ScanToTemplate(t *testing.T) {
  mapData := map[string]interface{}{}

  err := ScanToTemplate(mapData, map[string]interface{}{"key1": nil, "genericobject.GOTestStruct.Key": nil}) // error destMap not a pointer
  assert.Equal(t, "ERROR_ZGEN_SCANNER_DST_STRUCTURE_INVALID", err.Error())
  assert.Equal(t, "destination must be a pointer", err.GetList()[0].Args()["error"])

  err = ScanToTemplate(&mapData, []string{"key1", "key2"}) // error
  assert.Equal(t, "ERROR_ZGEN_SCANNER_DST_STRUCTURE_INVALID", err.Error())
  assert.Equal(t, "Invalid template, key must be string or map", err.GetList()[0].Args()["error"])

  err = ScanToTemplate(&mapData, map[string]interface{}{"key1": nil, "Key": nil, "notFoundKey": nil},
    ParserConfig{Mode: ParserModeNameAndTags},
    map[string]string{
      "key1": "value1",
      "key2": "value2"},
    map[string]int{
      "key3": 1,
      "key4": 2},
    GOTestStruct{
      ID:    10,
      Key:   "GTSKey",
      Value: "GTSValue"})
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, "value1", mapData["key1"])
  assert.Equal(t, "GTSKey", mapData["Key"])
  assert.Equal(t, nil, mapData["notFoundKey"])

  stringData := ""
  err = ScanToTemplate(&stringData, "key1",
    ParserConfig{Mode: ParserModeNameAndTags},
    map[string]string{
      "key1": "value1",
      "key2": "value2"},
    map[string]int{
      "key3": 1,
      "key4": 2},
    GOTestStruct{
      ID:    10,
      Key:   "GTSKey",
      Value: "GTSValue"})
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, "value1", stringData)

  intData := 0
  err = ScanToTemplate(&intData, "key3",
    ParserConfig{Mode: ParserModeNameAndTags},
    map[string]string{
      "key1": "value1",
      "key2": "value2"},
    map[string]int{
      "key3": 1,
      "key4": 2},
    GOTestStruct{
      ID:    10,
      Key:   "GTSKey",
      Value: "GTSValue"})
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, 1, intData)

  // assigning a string to a int variable
  intData = 0
  err = ScanToTemplate(&intData, "key1",
    ParserConfig{Mode: ParserModeNameAndTags},
    map[string]string{
      "key1": "value1",
      "key2": "value2"},
    map[string]int{
      "key3": 1,
      "key4": 2},
    GOTestStruct{
      ID:    10,
      Key:   "GTSKey",
      Value: "GTSValue"})

  assert.Equal(t, "ERROR_ZGEN_SCANNER_ARGUMENT_INVALID", err.Error())

}
