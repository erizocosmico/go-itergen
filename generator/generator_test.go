package generator

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type GeneratorSuite struct{}

var _ = Suite(&GeneratorSuite{})

func (s *GeneratorSuite) TestParseType(c *C) {
	tcs := []struct {
		raw  string
		typ  string
		pkg  string
		name string
	}{
		{"int64", "int64", "", "Int64"},
		{"github.com/mvader/go-itergen/generator:generator.Generator", "generator.Generator", "github.com/mvader/go-itergen/generator", "GeneratorGenerator"},
		{"os:*os.File", "*os.File", "os", "OsFile"},
	}
	for _, tc := range tcs {
		g := &Generator{RawType: tc.raw}
		g.parseType()
		c.Assert(g.Type.Type, Equals, tc.typ)
		c.Assert(g.Type.Package, Equals, tc.pkg)
		c.Assert(g.Type.Name, Equals, tc.name)
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

func (s *GeneratorSuite) TestGeneratePackage(c *C) {
	g := &Generator{Package: "foo"}
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generatePackage(buf), IsNil)
	c.Assert(buf.String(), Equals, "package foo\n\n")
}

func (s *GeneratorSuite) TestGenerateImports(c *C) {
	g := &Generator{}
	g.Type.Package = "os"
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateImports(buf), IsNil)
	c.Assert(buf.String(), Equals, "import \"os\"\n\n")
}

var generatedType = `type OsFileIter []*os.File

func NewOsFileIter(items ...*os.File) OsFileIter {
  return OsFileIter(items)
}
`

func (s *GeneratorSuite) TestGenerateType(c *C) {
	g := &Generator{}
	g.Type.Package = "os"
	g.Type.Type = "*os.File"
	g.Type.Name = "OsFile"
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateType(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedType)
}
