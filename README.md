# go-itergen [![Build Status](https://travis-ci.org/mvader/go-itergen.svg)](https://travis-ci.org/mvader/go-itergen)

**go-itergen** addresses a common and big problem of go: no maps, no filters, no nothing. If you come from a functional background that could be really frustrating. Since generics are not going to be around for a while code generation is the only way we have to achieve such things.

This is a naive attempt to make easier deal with this kind of operations over iterable types without having to write every single time the same code over and over.

## Available operations

go-itergen generates the following functions for a slice type:
* **Map:** apply a function to every element and return a slice with the modifications. It actually returns a XXXIterMapResult, which will have a set of operations to convert the `interface{}` result to other types.
* **Filter:** apply a function and will return a slice with all the elements whose result was true.
* **All:** will return true if all the elements return true after applying the given function.
* **Some:** will return true if any of the elements return true after applying the given function.
* **Concat:** will return a new slice with the contents of the current appending the given slice.
* **Find:** will return the item and the index of the first occurrence that returns true after applying the given function.
* **ForEach:** will execute a function for every item.
* **Reverse:** will return the slice in reversed order.
* **Splice:** will return a new slice with a number of items removed after the given start.
* **Reduce:** applies a function against an accumulator and each value of the slice (from left to right) to reduce it to a single value of the given type.

You can choose which operations you want for your type, that is, if you don't need `Map` or another function it won't be generated.

## Generate code

You just have to add that to a file in the package you want the code to be generated in.

```go
//go:generate go-itergen -t "float64" --pkg="mypkg" --map="string" --map="int" --filter --all --some --foreach --concat --find --reverse --splice --reduce="string" --reduce="int"
```

## Example

For example, this is what the generated code will look like for the previous generate directive.

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

func (i Float64Iter) Filter(fn func(float64) bool) Float64Iter {
	var result []float64
	for _, item := range i {
		if fn(item) {
			result = append(result, item)
		}
	}
	return Float64Iter(result)
}

func (i Float64Iter) All(fn func(float64) bool) bool {
	for _, item := range i {
		if !fn(item) {
			return false
		}
	}
	return true
}

func (i Float64Iter) Some(fn func(float64) bool) bool {
	for _, item := range i {
		if fn(item) {
			return true
		}
	}
	return false
}

func (i Float64Iter) ForEach(fn func(int, float64) interface{}) {
	for n, item := range i {
		fn(n, item)
	}
}

func (i Float64Iter) Concat(i2 Float64Iter) Float64Iter {
	return append(i, i2...)
}

func (i Float64Iter) Find(fn func(float64) bool) (float64, int) {
	var zero float64
	for i, item := range i {
		if fn(item) {
			return item, i
		}
	}
	return zero, -1
}

func (i Float64Iter) Reverse() Float64Iter {
	var result []float64
	for j := len(i) - 1; j >= 0; j-- {
		result = append(result, i[j])
	}
	return result
}

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

func (i Float64Iter) ReduceInt(fn func(current float64, acc int, index int) int, initial int) int {
	var result = initial
	for idx, item := range i {
		initial = fn(item, result, idx)
	}
	return result
}

func (i Float64Iter) ReduceString(fn func(current float64, acc string, index int) string, initial string) string {
	var result = initial
	for idx, item := range i {
		initial = fn(item, result, idx)
	}
	return result
}

```

And would be used like:

```go
func main() {
  rounded, err := NewFloat64Iter(1.2, 2.4, 3.5, 5.6).Filter(func(n float64) bool {
		return n > 2.0
	}).Map(func(int i, n float64) interface{} {
    return int(n)
  }).ToInt()
  fmt.Println(rounded) // [3 5]
}
```
