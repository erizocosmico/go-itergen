package generator

import (
	"io/ioutil"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type FileSuite struct{}

var _ = Suite(&FileSuite{})

func (s *FileSuite) TestFileify(c *C) {
	tcs := []struct {
		input  string
		output string
	}{
		{"int64", "int64"},
		{"interface{}", "interface"},
		{"*os.File", "osfile"},
	}

	for _, tc := range tcs {
		c.Assert(fileify(tc.input), Equals, tc.output)
	}
}

func (s *FileSuite) TestDeleteIfExists(c *C) {
	f := "test.txt"
	c.Assert(ioutil.WriteFile(f, []byte("hi"), 0644), IsNil)

	c.Assert(deleteIfExists(f), IsNil)
	_, err := os.Stat(f)
	c.Assert(os.IsNotExist(err), Equals, true)

	c.Assert(deleteIfExists(f), IsNil)
}

func (s *FileSuite) TestWrite(c *C) {
	f := "test.txt"
	c.Assert(ioutil.WriteFile(f, []byte("hi"), 0644), IsNil)

	c.Assert(write(f, []byte("Hello")), IsNil)
	b, err := ioutil.ReadFile(f)
	c.Assert(err, IsNil)
	c.Assert(string(b), Equals, "Hello")

	// Check that it does not fail when the file already exists
	c.Assert(write(f, []byte("Hello")), IsNil)

	c.Assert(deleteIfExists(f), IsNil)
}
