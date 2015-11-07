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
		t := (&Generator{}).parseType(tc.raw)
		c.Assert(t.Type, Equals, tc.typ)
		c.Assert(t.Package, Equals, tc.pkg)
		c.Assert(t.Name, Equals, tc.name)
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
	g1 := &Generator{}
	g2 := &Generator{}
	g2.Type.Package = "os"
	g3 := &Generator{}
	g3.Type.Package = "os"
	g3.MapResults = []TypeDef{
		TypeDef{Package: "foo"},
		TypeDef{Package: "github.com/foo/bar"},
	}

	tc := []struct {
		g      *Generator
		result string
	}{
		{g1, generatedImport1},
		{g2, generatedImport2},
		{g3, generatedImport3},
	}

	for _, t := range tc {
		buf := bytes.NewBuffer(nil)
		c.Assert(t.g.generateImports(buf), IsNil)
		c.Assert(buf.String(), Equals, t.result)
	}
}

func (s *GeneratorSuite) TestGenerateType(c *C) {
	g := &Generator{}
	g.Type.Package = "os"
	g.Type.Type = "*os.File"
	g.Type.Name = "OsFile"
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateType(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedType)
}

func (s *GeneratorSuite) TestGenerateMap(c *C) {
	g := &Generator{
		RawType: "float64",
		Map: []string{
			"int",
			"string",
		},
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateMap(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedMap)
}

func (s *GeneratorSuite) TestGenerateMapResults(c *C) {
	g := &Generator{
		RawType: "float64",
		Map: []string{
			"int",
			"string",
		},
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateMapResults(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedMapResults)
}

func (s *GeneratorSuite) TestGenerateFilters(c *C) {
	g := &Generator{
		RawType: "float64",
		Filter:  true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateFilter(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedFilter)
}

func (s *GeneratorSuite) TestGenerateSome(c *C) {
	g := &Generator{
		RawType: "float64",
		Some:    true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateSome(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedSome)
}

func (s *GeneratorSuite) TestGenerateAll(c *C) {
	g := &Generator{
		RawType: "float64",
		All:     true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateAll(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedAll)
}

func (s *GeneratorSuite) TestGenerateConcat(c *C) {
	g := &Generator{
		RawType: "float64",
		Concat:  true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateConcat(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedConcat)
}

func (s *GeneratorSuite) TestGenerateFind(c *C) {
	g := &Generator{
		RawType: "float64",
		Find:    true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateFind(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedFind)
}

func (s *GeneratorSuite) TestGenerateForEach(c *C) {
	g := &Generator{
		RawType: "float64",
		ForEach: true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateForEach(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedForEach)
}

func (s *GeneratorSuite) TestGenerateReverse(c *C) {
	g := &Generator{
		RawType: "float64",
		Reverse: true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateReverse(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedReverse)
}

func (s *GeneratorSuite) TestGenerateSplice(c *C) {
	g := &Generator{
		RawType: "float64",
		Splice:  true,
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateSplice(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedSplice)
}

func (s *GeneratorSuite) TestGenerateReducers(c *C) {
	g := &Generator{
		RawType: "float64",
		Reduce: []string{
			"int",
			"string",
		},
	}
	g.parseTypes()
	buf := bytes.NewBuffer(nil)
	c.Assert(g.generateReduces(buf), IsNil)
	c.Assert(buf.String(), Equals, generatedReducers)
}
