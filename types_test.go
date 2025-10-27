package zgen

import (
  "encoding/json"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"
)

func TestNullBool_MarshalJSON(t *testing.T) {
  nbVar := NullBool{}
  nbVar.Scan(nil)

  jsonVal, err := nbVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "null", string(jsonVal))

  nbVar.Scan(true)

  jsonVal, err = nbVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "true", string(jsonVal))
}

func TestNullBool_UnmarshalJSON(t *testing.T) {
  nbVar := NullBool{}

  err := json.Unmarshal([]byte("null"), &nbVar)
  assert.NoError(t, err)
  assert.Equal(t, false, nbVar.Valid)

  err = json.Unmarshal([]byte("true"), &nbVar)
  assert.NoError(t, err)
  assert.Equal(t, true, nbVar.Bool)
  assert.Equal(t, true, nbVar.Valid)
}

func TestNullInt64_MarshalJSON(t *testing.T) {
  niVar := NullInt64{}

  niVar.Scan(nil)

  jsonVal, err := niVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "null", string(jsonVal))

  niVar.Scan(3)

  jsonVal, err = niVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "3", string(jsonVal))
}

func TestNullInt64_UnmarshalJSON(t *testing.T) {
  niVar := NullInt64{}

  err := json.Unmarshal([]byte("null"), &niVar)
  assert.NoError(t, err)
  assert.Equal(t, false, niVar.Valid)

  err = json.Unmarshal([]byte("3"), &niVar)
  assert.NoError(t, err)
  assert.Equal(t, int64(3), niVar.Int64)
  assert.Equal(t, true, niVar.Valid)
}

func TestNullFloat64_MarshalJSON(t *testing.T) {
  nfVar := NullFloat64{}

  nfVar.Scan(nil)

  jsonVal, err := nfVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "null", string(jsonVal))

  nfVar.Scan(3.14)

  jsonVal, err = nfVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "3.14", string(jsonVal))
}

func TestNullFloat64_UnmarshalJSON(t *testing.T) {
  nfVar := NullFloat64{}

  err := json.Unmarshal([]byte("null"), &nfVar)
  assert.NoError(t, err)
  assert.Equal(t, false, nfVar.Valid)

  err = json.Unmarshal([]byte("3.14"), &nfVar)
  assert.NoError(t, err)
  assert.Equal(t, 3.14, nfVar.Float64)
  assert.Equal(t, true, nfVar.Valid)
}

func TestNullString_MarshalJSON(t *testing.T) {
  nsVar := NullString{}

  nsVar.Scan(nil)

  jsonVal, err := nsVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "null", string(jsonVal))

  nsVar.Scan(`test`)

  jsonVal, err = nsVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, `"test"`, string(jsonVal))
}

func TestNullString_UnmarshalJSON(t *testing.T) {
  nsVar := NullString{}

  err := json.Unmarshal([]byte("null"), &nsVar)
  assert.NoError(t, err)
  assert.Equal(t, false, nsVar.Valid)

  err = json.Unmarshal([]byte(`"test"`), &nsVar)
  assert.NoError(t, err)
  assert.Equal(t, "test", nsVar.String)
  assert.Equal(t, true, nsVar.Valid)
}

func TestNullTime_MarshalJSON(t *testing.T) {
  ntVar := NullTime{}

  ntVar.Scan(nil)

  jsonVal, err := ntVar.MarshalJSON()
  assert.NoError(t, err)
  assert.Equal(t, "null", string(jsonVal))

  currTime := time.Now().UTC()
  ntVar.Scan(currTime)

  jsonCT, _ := json.Marshal(currTime)
  jsonVal, err = ntVar.MarshalJSON()

  assert.NoError(t, err)
  assert.Equal(t, jsonCT, jsonVal)
}

func TestNullTime_UnmarshalJSON(t *testing.T) {
  ntVar := NullTime{}

  err := json.Unmarshal([]byte("null"), &ntVar)
  assert.NoError(t, err)
  assert.Equal(t, false, ntVar.Valid)

  currTime := time.Now().UTC()
  jsonCT, _ := json.Marshal(currTime)

  err = json.Unmarshal([]byte(jsonCT), &ntVar)
  assert.NoError(t, err)
  assert.Equal(t, currTime.String(), ntVar.Time.String())
  assert.Equal(t, true, ntVar.Valid)
}

func TestNullTime_ScanString(t *testing.T) {
  ntVar := NullTime{}

  err := ntVar.Scan("2020-01-01 00:15:30")
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, "2020-01-01 00:15:30", ntVar.Time.Format("2006-01-02 15:04:05"))
}

func TestNullTime_ScanInt(t *testing.T) {
  ntVar := NullTime{}

  currTime := time.Now().UnixNano()

  err := ntVar.Scan(currTime)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, currTime, ntVar.Time.UnixNano())
}
