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
