package examples

import (
	"testing"

	. "gopkg.in/check.v1"
)

var _ = Suite(&IterSuite{})

func Test(t *testing.T) { TestingT(t) }

type IterSuite struct{}

func (s *IterSuite) TestMap(c *C) {
	expected := Float64IterMapResult{1, 2, 3}
	iter := NewFloat64Iter(1.2, 2.3, 3.4)

	result := iter.Map(fToInt)

	c.Assert(result, DeepEquals, expected)
}

func (s *IterSuite) TestMapResult(c *C) {
	expected := []int{1, 2, 3}
	iter := NewFloat64Iter(1.2, 2.3, 3.4)

	result, err := iter.Map(func(n int, f float64) interface{} {
		return int(f)
	}).ToInt()

	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, expected)
}

func (s *IterSuite) TestFilter(c *C) {
	expected := NewFloat64Iter(1., 3.)
	iter := NewFloat64Iter(1., 2., 3., 4.)

	result := iter.Filter(isOdd)

	c.Assert(result, DeepEquals, expected)
}

func isOdd(f float64) bool {
	return int(f)%2 == 1
}

func (s *IterSuite) TestAll(c *C) {
	a := NewFloat64Iter(1., 3., 5.)
	b := NewFloat64Iter(.2, .3, .5)

	c.Assert(a.All(isOdd), Equals, true)
	c.Assert(b.All(isOdd), Equals, false)
}

func (s *IterSuite) TestSome(c *C) {
	a := NewFloat64Iter(6., 2., 4.)
	b := NewFloat64Iter(2., 3., 5.)

	c.Assert(a.Some(isOdd), Equals, false)
	c.Assert(b.Some(isOdd), Equals, true)
}

func (s *IterSuite) TestForEach(c *C) {
	var arr []float64
	NewFloat64Iter(1., 2., 3.).ForEach(func(i int, f float64) {
		arr = append(arr, f)
	})
	c.Assert(arr, HasLen, 3)
}

func (s *IterSuite) TestConcat(c *C) {
	result := NewFloat64Iter(1., 2.).Concat(NewFloat64Iter(3., 4.))
	c.Assert(result, DeepEquals, NewFloat64Iter(1., 2., 3., 4.))
}

func (s *IterSuite) TestFind(c *C) {
	n, idx := NewFloat64Iter(2., 3., 4., 5.).Find(isOdd)
	c.Assert(int(n), Equals, 3)
	c.Assert(idx, Equals, 1)
}

func (s *IterSuite) TestReverse(c *C) {
	iter := NewFloat64Iter(1., 2., 3.)
	c.Assert(iter.Reverse(), DeepEquals, NewFloat64Iter(3., 2., 1.))
}

func (s *IterSuite) TestReduce(c *C) {
	iter := NewFloat64Iter(1., 2., 3.)
	n := iter.ReduceInt(func(current float64, acc int, idx int) int {
		return acc + int(current)
	}, 0)
	c.Assert(n, Equals, 6)
}

func (s *IterSuite) TestSplice(c *C) {
	iter := NewFloat64Iter(1., 2., 3.)

	c.Assert(iter.Splice(-1, 1), DeepEquals, iter)
	c.Assert(iter.Splice(0, 1), DeepEquals, NewFloat64Iter(2., 3.))
	c.Assert(iter.Splice(0, 4), DeepEquals, NewFloat64Iter())
	c.Assert(iter.Splice(1, 4), DeepEquals, NewFloat64Iter(1.))
}

func fToInt(n int, f float64) interface{} {
	return int(f)
}
