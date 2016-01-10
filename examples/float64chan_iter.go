package examples

import (
	"errors"
	"sync"
)

type Float64ChanIter chan float64

type Float64ChanMapResult <-chan interface{}

func (i Float64ChanIter) Map(fn func(int, float64) interface{}) Float64ChanMapResult {
	out := make(chan interface{})

	go func() {
		var idx int
		for v := range i {
			out <- fn(idx, v)
			idx++
		}
		close(out)
	}()

	return out
}

var ErrFloat64ChanToFloat64 = errors.New("cannot convert Float64ChanMapResult to chan float64")

func (r Float64ChanMapResult) Iter() (Float64ChanIter, chan error) {
	out := make(chan float64)
	err := make(chan error)

	go func() {
		for v := range r {
			if _, ok := v.(float64); !ok {
				err <- ErrFloat64ChanToFloat64
				break
			}
			out <- v.(float64)
		}
		close(err)
		close(out)
	}()

	return out, err
}

var ErrFloat64ChanToInt = errors.New("cannot convert Float64ChanMapResult to chan int")

func (r Float64ChanMapResult) ToInt() (chan int, chan error) {
	out := make(chan int)
	err := make(chan error)

	go func() {
		for v := range r {
			if _, ok := v.(int); !ok {
				err <- ErrFloat64ChanToInt
				break
			}
			out <- v.(int)
		}
		close(err)
		close(out)
	}()

	return out, err
}

func (i Float64ChanIter) Filter(fn func(float64) bool) Float64ChanIter {
	out := make(chan float64)

	go func() {
		for v := range i {
			if fn(v) {
				out <- v
			}
		}
		close(out)
	}()

	return out
}

func (i Float64ChanIter) ForEach(fn func(int, float64)) {
	var n int
	go func() {
		for v := range i {
			fn(n, v)
			n++
		}
	}()
}

func (i Float64ChanIter) Concat(args ...Float64ChanIter) Float64ChanIter {
	var (
		out   = make(chan float64)
		wg    sync.WaitGroup
		chans = []Float64ChanIter{i}
	)

	for _, a := range args {
		chans = append(chans, a)
	}

	for _, c := range chans {
		wg.Add(1)
		go func(in Float64ChanIter) {
			for v := range in {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (i Float64ChanIter) ReduceInt(fn func(current float64, acc int, index int) int, initial int) chan int {
	out := make(chan int)
	result := initial

	go func() {
		var idx int
		for item := range i {
			result = fn(item, result, idx)
			idx++
		}

		out <- result
		close(out)
	}()

	return out
}

func (i Float64ChanIter) Array(done chan struct{}) []float64 {
	var (
		result []float64
		wg     sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		for v := range i {
			result = append(result, v)
		}
		wg.Done()
	}()

	defer func() {
		done <- struct{}{}
	}()

	wg.Wait()
	return result
}
