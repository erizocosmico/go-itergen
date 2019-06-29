# go-itergen [![Build Status](https://travis-ci.org/erizocosmico/go-itergen.svg)](https://travis-ci.org/erizocosmico/go-itergen)

**go-itergen** addresses a common and big problem of go: no maps, no filters, no nothing. If you come from a functional background that could be really frustrating. Since generics are not going to be around for a while code generation is the only way we have to achieve such things.

This is a naive attempt to make easier deal with this kind of operations over iterable types without having to write every single time the same code over and over.

It also has a very nice feature: it works with channels, not just slices. So, you can map, filter, concat and reduce channels. And even convert them to arrays. All without always having to write the same boilerplate code. TL;DR: channel manipulation made easy.

## Install

```
go get -u github.com/erizocosmico/go-itergen/cmd/...
```

## Available operations

go-itergen generates the following functions for a slice type:
* **Map:** apply a function to every element and return a slice/channel with the modifications. It actually returns a XXXIterMapResult, which will have a set of operations to convert the `interface{}` result to other types.
* **Filter:** apply a function and will return a slice/channel with all the elements whose result was true.
* **All (only for slices):** will return true if all the elements return true after applying the given function.
* **Some (only for slices):** will return true if any of the elements return true after applying the given function.
* **Concat:** 
  * In the case of a slice, will return a new array with the contents of `a` and `b` given `a.Concat(b)`.
  * In the case of a channel, will return a new channel multiplexing all the other given channels. E.g: `a.Concat(b, c, d)` will return a channel with all the items sent to `a`, `b`, `c` and `d`.
* **Find (only for slices):** will return the item and the index of the first occurrence that returns true after applying the given function.
* **ForEach:** will execute a function for every item in the slice/channel.
* **Reverse (only for slices):** will return the slice in reversed order.
* **Splice (only for slices):** will return a new slice with a number of items removed after the given start.
* **Reduce:** applies a function against an accumulator and each value of the slice/channel (from first to last) to reduce it to a single value of the given type.
* **Array (only for channels):** converts the channel into an array. The operation blocks, but can be done in a goroutine and you will be notified via the `done` parameter.

You can choose which operations you want for your type, that is, if you don't need `Map` or another function it won't be generated.

## Generate code

You just have to add that to a file in the package you want the code to be generated in.

```go
//go:generate go-itergen -t "float64" --pkg="mypkg" --map="string" --map="int" --filter --all --some --foreach --concat --find --reverse --splice --reduce="string" --reduce="int"
```

Or execute the binary:

```bash
go-itergen -t "float64" --pkg="mypkg" --map="string" --map="int" --filter --all --some --foreach --concat --find --reverse --splice --reduce="string" --reduce="int"
```

#### Types from external packages

If you want to generate an iterable type for an external package, you can do that with `:`.
For example, let's say we want a `os.File` iterable. We would do it like this:
```
go-itergen -t "os:os.File" --pkg="mypkg" --map="string"
```

Note that what goes before the `:` is what would go in the import, and that you have to type the full type name. This is because of how Go packages work. One can not guarantee that `github.com/foo/go-bar` will be `go-bar.Foo`. It may be `bar.Foo`. And thus, it has to be specified.

Another example:
```
go-itergen -t "golang.org/x/net/context:context.Context" --pkg="mypkg" --map="github.com/foo/ctx:ctx.MyCtx"
```

Take a look at the map. We can also specify the external packages in `map` and `reduce` arguments.

## Example

For examples of generated code see the `examples` folder. Contains a file with a `chan float64` iterable and another with a `float64` slice iterable.

**Usage example:**

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
