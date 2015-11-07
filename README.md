# go-itergen

**go-itergen** addresses a common and big problem of go: no maps, no filters, no nothing. If you come from a functional background that could be really frustrating. Since generics are not going to be around for a while code generation is the only way we have to achieve such things.

This is a naive attempt to make easier deal with this kind of operations over iterable types without having to write every single time the same code over and over.

## Available operations

go-itergen generates the following functions for an array type:
* [ ] Map: apply a function to every element and return an array with the modifications. It actually returns a XXXIterMapResult, which will have a set of operations to convert the `interface{}` result to other types.
* [ ] Filter: apply a function and will return an array with all the elements whose result was true.
* [ ] All: will return true if all the elements return true after applying the given function.
* [ ] Some: will return true if any of the elements return true after applying the given function.

More will come after these four, but this is the basic functionality that is going to be provided.

You can choose which operations you want for your type, that is, if you don't need `Map` it won't be generated.

## Generate code

You just have to add that to a file in the package you want the code to be generated in.

```go
//go:generate go-itergen -t "float64" --pkg="mypkg" --map="int" --map="string"
```

## Example

For example, this is what the generated code will look like for a Map operation.

```go
package mypkg

import (
	"errors"
)

type Float64Iter []float64

func NewFloat64Iter(items ...float64) Float64Iter {
	return Float64Iter(items)
}

type Float64IterMapResult []interface{}

func (i Float64Iter) Map(fn func(float64) interface{}) Float64IterMapResult {
	var result []interface{}
	for _, item := range i {
		result = append(result, fn(item))
	}
	return result
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

```

And would be used like:

```go
func main() {
  rounded, err := NewFloat64Iter(1.2, 2.4, 3.5, 5.6).Map(func(n float64) interface{} {
    return int(n)
  }).ToInt()
  fmt.Println(rounded)
}
```

## TODO

* [x] Write generator
* [x] Generate `XXXIter` type
* [x] Generate `Map` function
* [x] Generate `XXXIterMapResult` transformers
* [ ] Generate `Filter` function
* [ ] Generate `All` function
* [ ] Generate `Some` function
