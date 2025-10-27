package zgen

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

type CloneCircularTestStruct1 struct {
  Value     int
  Circular2 *CloneCircularTestStruct2
}
type CloneCircularTestStruct2 struct {
  Value     int
  Circular1 *CloneCircularTestStruct1
}

func TestUnit_UnpackToBase(t *testing.T) {
  intTest := 1

  response := UnpackBaseElement(intTest, false)
  assert.Equal(t, intTest, response)

  // pack it inside 2 interfaces and a pointer
  packedIntTest := interface{}(interface{}(intTest))
  response = UnpackBaseElement(&packedIntTest, false)
  assert.Equal(t, intTest, response)

  stringTest := "Test"
  response = UnpackBaseElement(stringTest, false)
  assert.Equal(t, stringTest, response)

  packedStringTest := interface{}(interface{}(&stringTest))
  response = UnpackBaseElement(packedStringTest, false)
  assert.Equal(t, stringTest, response)

  var interfaceTest interface{}
  interfaceTest = "Test"
  response = UnpackBaseElement(interfaceTest, true)
  assert.Equal(t, interfaceTest, response)

  interfaceTest = interface{}(interface{}(map[string]interface{}{
    "key1": "value1",
    "key2": 2}))
  response = UnpackBaseElement(interfaceTest, true)
  assert.Equal(t, interfaceTest, response)

  // simulating a multiple passthru with pointers and all
  ift := interface{}(&interfaceTest)
  if _, ok := ift.(map[string]interface{}); ok {
    t.Error("Map Assertion Error")
  }
  response = UnpackBaseElement(ift, true)
  assert.Equal(t, interfaceTest, response)

  sft := interface{}(&packedStringTest)
  if _, ok := ift.(string); ok {
    t.Error("String Assertion")
  }
  response = UnpackBaseElement(sft, true) // will return a *string
  if _, ok := response.(*string); !ok {
    t.Error("String Assertion 2")
  }

  response = UnpackBaseElement(sft, false) // will unpack all the way including *string to string
  assert.Equal(t, stringTest, response)

}

func TestUnit_IsZeroValue(t *testing.T) {
  var i int
  assert.Equal(t, true, IsZeroValue(i))
  i = 5
  assert.Equal(t, false, IsZeroValue(i))

  var str string
  assert.Equal(t, true, IsZeroValue(str))
  str = "a"
  assert.Equal(t, false, IsZeroValue(str))

  var ar [2]int
  assert.Equal(t, true, IsZeroValue(ar))
  ar[1] = 5
  assert.Equal(t, false, IsZeroValue(ar))

  var mp map[string]string
  assert.Equal(t, true, IsZeroValue(mp))
  mp = map[string]string{"a": "b"}
  assert.Equal(t, false, IsZeroValue(mp))

  var Cf struct {
    IntVal  int
    BoolVal bool
  }

  assert.Equal(t, false, IsZeroValue(&Cf))
  assert.Equal(t, true, IsZeroValue(Cf))
  Cf.BoolVal = true
  assert.Equal(t, false, IsZeroValue(Cf))
}

//
func TestUnit_CloneMap(t *testing.T) {
  src := map[string]interface{}{
    "key1": "val1",
    "key2": "val2",
  }
  clone, _ := MapStringAny(Clone(src))

  assert.Equal(t, src["key1"], clone["key1"])
  src["key1"] = "othervalue" // changing value, if clone is not a real clone this should change
  assert.NotEqual(t, src["key1"], clone["key1"])
}

func TestUnit_CloneSlice(t *testing.T) {
  src := []interface{}{
    "val1",
    "val2",
  }
  clone, _ := SliceAny(Clone(src))

  assert.Equal(t, src[0], clone[0])
  src[0] = "othervalue" // changing value, if clone is not a real clone this should change
  assert.NotEqual(t, src[0], clone[0])
}

func TestUnit_CloneStruct(t *testing.T) {
  type CloneTestStruct struct {
    Value int
  }
  src := CloneTestStruct{Value: 1}

  clone, ok := Clone(src).(CloneTestStruct)
  if !ok {
    t.Error("Clone struct failed")
  }

  assert.Equal(t, src.Value, clone.Value)
  src.Value = 2 // changing value, if clone is not a real clone this should change
  assert.NotEqual(t, src.Value, clone.Value)
}

func TestUnit_CloneCircular(t *testing.T) {
  str2 := CloneCircularTestStruct2{
    Value: 1,
  }
  str1 := CloneCircularTestStruct1{
    Value:     1,
    Circular2: &str2,
  }
  str2.Circular1 = &str1

  clone, ok := Clone(str1).(CloneCircularTestStruct1) // this will trigger max iteration count and should stop there
  if !ok {
    t.Error("Clone struct failed")
  }

  assert.Equal(t, str1.Value, clone.Value)
  assert.Equal(t, str1.Circular2.Value, clone.Circular2.Value)
  str1.Value = 2 // changing value, if clone is not a real clone this should change
  str2.Value = 2
  assert.NotEqual(t, str1.Value, clone.Value)
  assert.NotEqual(t, str1.Circular2.Value, clone.Circular2.Value)
}

func TestUnit_CloneSimple(t *testing.T) {
  src := "x"
  clone, _ := String(Clone(src))

  assert.Equal(t, src, clone) // all simple types are passed by value so we do not need to test changed values
}

func TestUnit_DeepMergeSimple(t *testing.T) {
  elem1 := "x"
  elem2 := 12

  res, err := DeepMerge(elem1, elem2, 1)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, "x", res)

  res, err = DeepMerge(elem1, elem2, 2)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, 12, res)
}

func TestUnit_DeepMergeSlice(t *testing.T) {
  elem1 := []interface{}{1, "2", 3.14}
  elem2 := [3]interface{}{4, "5", 6.7}

  res, err := DeepMerge(elem1, elem2, FlagDeepMergePriorityFirst)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, []interface{}{4, "5", 6.7, 1, "2", 3.14}, res)

  res, err = DeepMerge(elem1, elem2, FlagDeepMergePrioritySecond)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, []interface{}{1, "2", 3.14, 4, "5", 6.7}, res)

  res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePriorityFirst)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, []interface{}{1, "2", 3.14}, res)

  res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePrioritySecond)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, []interface{}{4, "5", 6.7}, res)

}

func TestUnit_DeepMergeMap(t *testing.T) {
  elem1 := map[string]interface{}{
    "simpleInt":    12,
    "simpleString": "a string",
    "simpleMap": map[string]interface{}{
      "x1": 1,
      "x2": "x2Val1",
      "x3": []interface{}{1, "2", 3.14},
      "x4": map[string]interface{}{
        "x4v1": "1",
        "x4v2": 2,
        "x4v3": []interface{}{1, "2", 3.14},
      },
    },
  }
  elem2 := map[string]interface{}{
    "simpleIntg":   13,
    "simpleString": "another string",
    "simpleMap": map[string]interface{}{
      "x1.1": 1,
      "x2":   "x2Val2",
      "x3":   []interface{}{4, "5", 6.7},
      "x4": map[string]interface{}{
        "x4vx": "2",
        "x4v2": 3,
        "x4v3": [3]interface{}{4, "5", 6.7}, // array here, will be converted to slice
      },
    },
  }
  elemResult1 := map[string]interface{}{
    "simpleInt":    12,
    "simpleIntg":   13,
    "simpleString": "a string",
    "simpleMap": map[string]interface{}{
      "x1":   1,
      "x1.1": 1,
      "x2":   "x2Val1",
      "x3":   []interface{}{4, "5", 6.7, 1, "2", 3.14},
      "x4": map[string]interface{}{
        "x4v1": "1",
        "x4vx": "2",
        "x4v2": 2,
        "x4v3": []interface{}{4, "5", 6.7, 1, "2", 3.14},
      },
    },
  }

  elemResult2 := map[string]interface{}{
    "simpleInt":    12,
    "simpleIntg":   13,
    "simpleString": "another string",
    "simpleMap": map[string]interface{}{
      "x1":   1,
      "x1.1": 1,
      "x2":   "x2Val2",
      "x3":   []interface{}{1, "2", 3.14, 4, "5", 6.7},
      "x4": map[string]interface{}{
        "x4v1": "1",
        "x4vx": "2",
        "x4v2": 3,
        "x4v3": []interface{}{1, "2", 3.14, 4, "5", 6.7},
      },
    },
  }

  elemResult3 := map[string]interface{}{
    "simpleInt":    12,
    "simpleIntg":   13,
    "simpleString": "a string",
    "simpleMap": map[string]interface{}{
      "x1": 1,
      "x2": "x2Val1",
      "x3": []interface{}{1, "2", 3.14},
      "x4": map[string]interface{}{
        "x4v1": "1",
        "x4v2": 2,
        "x4v3": []interface{}{1, "2", 3.14},
      },
    },
  }

  elemResult4 := map[string]interface{}{
    "simpleInt":    12,
    "simpleIntg":   13,
    "simpleString": "another string",
    "simpleMap": map[string]interface{}{
      "x1.1": 1,
      "x2":   "x2Val2",
      "x3":   []interface{}{4, "5", 6.7},
      "x4": map[string]interface{}{
        "x4vx": "2",
        "x4v2": 3,
        "x4v3": [3]interface{}{4, "5", 6.7}, // array here, will be converted to slice
      },
    },
  }

  res, err := DeepMerge(elem1, elem2, FlagDeepMergePriorityFirst)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, elemResult1, res)

  res, err = DeepMerge(elem1, elem2, FlagDeepMergePrioritySecond)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, elemResult2, res)

  res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePriorityFirst)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, elemResult3, res)

  res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled+FlagDeepMergePrioritySecond)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, elemResult4, res)

  res, err = DeepMerge(elem1, elem2, FlagDeepMergeOverwriteEnabled)
  if err != nil {
    t.Error(err)
  }
  assert.Equal(t, elemResult4, res)
}
