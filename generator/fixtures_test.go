package generator

var generatedImport1 = `import (
  "errors"
)
`

var generatedImport2 = `import (
  "errors"
  "os"
)
`

var generatedImport3 = `import (
  "errors"
  "foo"
  "github.com/foo/bar"
  "os"
)
`

var generatedType = `type OsFileIter []*os.File

func NewOsFileIter(items ...*os.File) OsFileIter {
  return OsFileIter(items)
}
`

var generatedMap = `
type Float64IterMapResult []interface{}

func (i Float64Iter) Map(fn func(int, float64) interface{}) Float64IterMapResult {
  var result []interface{}
  for n, item := range i {
    result = append(result, fn(n, item))
  }
  return result
}

var ErrFloat64ToFloat64 = errors.New("cannot convert Float64IterMapResult to []float64")

func (r Float64IterMapResult) Iter() (Float64Iter, error) {
  var result []float64
  for _, i := range r {
    if _, ok := i.(float64); !ok {
      return nil, ErrFloat64ToFloat64
    }
    result = append(result, i.(float64))
  }
  return Float64Iter(result), nil
}
`

var generatedMapResults = `
var ErrFloat64ToInt = errors.New("cannot convert Float64IterMapResult to []int")

func (r Float64IterMapResult) ToInt() ([]int, error) {
  var result []int
  for _, i := range r {
    if _, ok := i.(int); !ok {
      return nil, ErrFloat64ToInt
    }
    result = append(result, i.(int))
  }
  return result, nil
}
var ErrFloat64ToString = errors.New("cannot convert Float64IterMapResult to []string")

func (r Float64IterMapResult) ToString() ([]string, error) {
  var result []string
  for _, i := range r {
    if _, ok := i.(string); !ok {
      return nil, ErrFloat64ToString
    }
    result = append(result, i.(string))
  }
  return result, nil
}
`

var generatedFilter = `
func (i Float64Iter) Filter(fn func(float64) bool) Float64Iter {
  var result []float64
  for _, item := range i {
    if fn(item) {
      result = append(result, item)
    }
  }
  return Float64Iter(result)
}
`

var generatedSome = `
func (i Float64Iter) Some(fn func(float64) bool) bool {
  for _, item := range i {
    if fn(item) {
      return true
    }
  }
  return false
}
`
var generatedAll = `
func (i Float64Iter) All(fn func(float64) bool) bool {
  for _, item := range i {
    if !fn(item) {
      return false
    }
  }
  return true
}
`

var generatedConcat = `
func (i Float64Iter) Concat(i2 Float64Iter) Float64Iter {
  return append(i, i2...)
}
`

var generatedFind = `
func (i Float64Iter) Find(fn func(float64) bool) (float64, int) {
  var zero float64
  for i, item := range i {
    if fn(item) {
      return item, i
    }
  }
  return zero, -1
}
`

var generatedForEach = `
func (i Float64Iter) ForEach(fn func(int, float64)) {
  for n, item := range i {
    fn(n, item)
  }
}
`

var generatedReverse = `
func (i Float64Iter) Reverse() Float64Iter {
  var result []float64
  for j := len(i)-1; j >= 0; j-- {
    result = append(result, i[j])
  }
  return result
}
`

var generatedSplice = `
// Splice removes numDelete items from the slice
// since start. If numDelete is -1 it will delete all
// items after start. If start is higher than the
// slice length or lower than 0 the whole slice
// will be returned.
func (i Float64Iter) Splice(start, numDelete int) Float64Iter {
  var result Float64Iter
  length := len(i)
  if start >= length-1 || start < 0 {
    return i
  }

  result = append(result, i[:start]...)
  if numDelete > -1 && numDelete+start < length {
    result = append(result, i[start+numDelete:]...)
  }

  return result
}
`

var generatedReducers = `

func (i Float64Iter) ReduceInt(fn func(current float64, acc int, index int) int, initial int) int {
  var result = initial
  for idx, item := range i {
    result = fn(item, result, idx)
  }
  return result
}

func (i Float64Iter) ReduceString(fn func(current float64, acc string, index int) string, initial string) string {
  var result = initial
  for idx, item := range i {
    result = fn(item, result, idx)
  }
  return result
}

`
