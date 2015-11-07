package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"strings"
)

// Generator generates functions for iterable types based on the options received
type Generator struct {
	RawType string   `short:"t" long:"type" description:"type to generate the code for"`
	Map     []string `long:"map" description:"generate Map function with transformer for given type"`
	Filter  bool     `long:"filter" description:"generate Filter function"`
	All     bool     `long:"all" description:"generate All function"`
	Some    bool     `long:"some" description:"generate Some function"`

	Type struct {
		Package string
		Type    string
	}
}

type generatorFunc func() ([]byte, error)

func (g *Generator) parseType() {
	typeParts := strings.Split(g.RawType, ":")
	if len(typeParts) == 1 {
		g.Type.Type = g.RawType
	} else {
		g.Type.Package = typeParts[0]
		g.Type.Type = typeParts[1]
	}
}

func (g *Generator) generatePackage() ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateImports() ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateType() ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateSome() ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateAll() ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateFilter() ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateMap() ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateMapResults() ([]byte, error) {
	return nil, nil
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
		b, err := gen()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
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
