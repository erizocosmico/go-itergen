package generator

import . "gopkg.in/check.v1"

type GeneratorSuite struct{}

var _ = Suite(&GeneratorSuite{})

func (s *GeneratorSuite) TestParseType(c *C) {
	tcs := []struct {
		raw string
		typ string
		pkg string
	}{
		{"int64", "int64", ""},
		{"github.com/mvader/go-itergen/generator:generator.Generator", "generator.Generator", "github.com/mvader/go-itergen/generator"},
		{"os:*os.File", "*os.File", "os"},
	}
	for _, tc := range tcs {
		g := &Generator{RawType: tc.raw}
		g.parseType()
		c.Assert(g.Type.Type, Equals, tc.typ)
		c.Assert(g.Type.Package, Equals, tc.pkg)
	}
}

func (s *GeneratorSuite) TestFileName(c *C) {
	tcs := []struct {
		typ    string
		output string
	}{
		{"*os.File", "osfile_iter.go"},
		{"int64", "int64_iter.go"},
	}

	for _, tc := range tcs {
		g := &Generator{}
		g.Type.Type = tc.typ
		c.Assert(g.fileName(), Equals, tc.output)
	}
}
