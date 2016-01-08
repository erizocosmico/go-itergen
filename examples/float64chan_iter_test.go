package examples

import (
	"sort"
	"sync"

	. "gopkg.in/check.v1"
)

var _ = Suite(&ChanSuite{})

type ChanSuite struct{}

func (s *ChanSuite) TestMap(c *C) {
	var result []interface{}
	var i = make(Float64ChanIter)

	ch := i.Map(fToInt)

	go func() {
		i <- 1.2
		i <- 2.5
		i <- 3.4
		close(i)
	}()

	for v := range ch {
		result = append(result, v)
	}

	c.Assert(result, DeepEquals, []interface{}{1, 2, 3})
}

func (s *ChanSuite) TestMapTo(c *C) {
	var result []int
	var i = make(Float64ChanIter)

	out, err := i.Map(fToInt).ToInt()

	go func() {
		i <- 1.2
		i <- 2.5
		i <- 3.4
		close(i)
	}()

	for v := range out {
		result = append(result, v)
	}

	_, ok := <-err
	c.Assert(ok, Equals, false)
	c.Assert(result, DeepEquals, []int{1, 2, 3})
}

func (s *ChanSuite) TestFilter(c *C) {
	var result []float64
	var i = make(Float64ChanIter)

	out := i.Filter(isOdd)

	go func() {
		i <- 1.2
		i <- 2.5
		i <- 3.4
		close(i)
	}()

	for v := range out {
		result = append(result, v)
	}

	c.Assert(result, DeepEquals, []float64{1.2, 3.4})
}

func (s *ChanSuite) TestForEach(c *C) {
	var result []float64
	var wg sync.WaitGroup
	wg.Add(3)
	var i = make(Float64ChanIter)

	go func() {
		i <- 1.2
		i <- 2.5
		i <- 3.4
		close(i)
	}()

	i.ForEach(func(i int, f float64) {
		result = append(result, f)
		wg.Done()
	})

	wg.Wait()
	c.Assert(result, DeepEquals, []float64{1.2, 2.5, 3.4})
}

func (s *ChanSuite) TestConcat(c *C) {
	var result []int
	var a = make(Float64ChanIter)
	var b = make(Float64ChanIter)
	var d = make(Float64ChanIter)

	ch := a.Concat(b, d)

	go func() {
		a <- 1.2
		b <- 2.5
		a <- 3.4
		b <- 4.6
		a <- 5.8
		d <- 6.0
		close(a)
		close(b)
		close(d)
	}()

	for v := range ch {
		result = append(result, int(v))
	}

	sort.Ints(result)

	c.Assert(result, DeepEquals, []int{1, 2, 3, 4, 5, 6})
}

func (s *ChanSuite) TestReduce(c *C) {
	var i = make(Float64ChanIter)

	out := i.ReduceInt(func(current float64, acc int, index int) int {
		return acc + int(current)
	}, 0)

	go func() {
		i <- 1.2
		i <- 2.5
		i <- 3.4
		close(i)
	}()

	v := <-out
	c.Assert(v, Equals, 6)
}

func (s *ChanSuite) TestArray(c *C) {
	var i = make(Float64ChanIter)
	var done = make(chan struct{})
	var result []float64

	go func() {
		result = i.Array(done)
	}()

	go func() {
		i <- 1.2
		i <- 2.5
		i <- 3.4
		close(i)
	}()

	<-done
	close(done)
	c.Assert(result, DeepEquals, []float64{1.2, 2.5, 3.4})
}
