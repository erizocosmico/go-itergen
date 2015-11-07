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
  "os"
  "foo"
  "github.com/foo/bar"
)
`

var generatedType = `type OsFileIter []*os.File

func NewOsFileIter(items ...*os.File) OsFileIter {
  return OsFileIter(items)
}
`

var generatedMap = `
type Float64IterMapResult []interface{}

func (i Float64Iter) Map(fn func(float64) interface{}) Float64IterMapResult {
  var result []interface{}
  for _, item := range i {
    result = append(result, fn(item))
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
