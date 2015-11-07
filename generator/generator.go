package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"path/filepath"
	"strings"
)

// Generator generates functions for iterable types based on the options received
type Generator struct {
	RawType string   `short:"t" long:"type" description:"type to generate the code for" required:"true"`
	Package string   `long:"pkg" description:"package of the resultant file" required:"true"`
	Map     []string `long:"map" description:"generate Map function with transformer for given type"`
	Filter  bool     `long:"filter" description:"generate Filter function"`
	All     bool     `long:"all" description:"generate All function"`
	Some    bool     `long:"some" description:"generate Some function"`

	Type struct {
		Package string
		Type    string
		Name    string
	}
}

type generatorFunc func(io.Writer) error

func (g *Generator) parseType() {
	g.Type.Package, g.Type.Type = g.parseRawType(g.RawType)
	g.Type.Name = g.getTypeName(g.Type.Type)
}

func (g *Generator) parseRawType(raw string) (string, string) {
	var (
		pkg string
		typ string
	)

	typeParts := strings.Split(raw, ":")
	if len(typeParts) == 1 {
		typ = g.RawType
	} else {
		pkg = typeParts[0]
		typ = typeParts[1]
	}

	return pkg, typ
}

func (g *Generator) getTypeName(t string) string {
	if strings.HasPrefix(t, "*") {
		t = t[1:]
	}

	if strings.Contains(t, ".") {
		tParts := strings.Split(t, ".")
		return strings.Title(tParts[0]) + strings.Title(tParts[1])
	}

	return strings.Title(t)
}

func (g *Generator) generatePackage(w io.Writer) error {
	pkg := fmt.Sprintf("package %s\n\n", g.Package)
	_, err := w.Write([]byte(pkg))
	return err
}

func (g *Generator) generateImports(w io.Writer) error {
	if g.Type.Package != "" {
		imp := fmt.Sprintf("import \"%s\"\n\n", g.Type.Package)
		_, err := w.Write([]byte(imp))
		return err
	}

	return nil
}

func (g *Generator) generateType(w io.Writer) error {
	return typeTpl.Execute(w, g.Type)
}

func (g *Generator) generateSome(w io.Writer) error {
	return nil
}

func (g *Generator) generateAll(w io.Writer) error {
	return nil
}

func (g *Generator) generateFilter(w io.Writer) error {
	return nil
}

func (g *Generator) generateMap(w io.Writer) error {
	return nil
}

func (g *Generator) generateMapResults(w io.Writer) error {
	return nil
}

func (g *Generator) generateCode() ([]byte, error) {
	generators := []generatorFunc{
		g.generatePackage,
		g.generateImports,
		g.generateType,
		g.generateMap,
		g.generateMapResults,
		g.generateFilter,
		g.generateAll,
		g.generateSome,
	}

	buf := bytes.NewBuffer(nil)
	for _, gen := range generators {
		err := gen(buf)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (g *Generator) fileName() string {
	file := fmt.Sprintf("%s_iter.go", fileify(g.Type.Type))
	return filepath.Join(".", file)
}

// Generate writes the generated code to the correspondant file and returns an error if something failed
func (g *Generator) Generate() error {
	g.parseType()
	code, err := g.generateCode()
	if err != nil {
		return err
	}

	code, err = format.Source(code)
	if err != nil {
		return err
	}

	return write(g.fileName(), code)
}
