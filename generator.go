package generator

//go:generate statik -src=templates

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// Generator generates functions for iterable types based on the options received
type Generator struct {
	RawType string   `short:"t" long:"type" description:"type to generate the code for" required:"true"`
	Package string   `long:"pkg" description:"package of the resultant file" required:"true"`
	Map     []string `long:"map" description:"generate Map function with transformer for given type"`
	Filter  bool     `long:"filter" description:"generate Filter function"`
	All     bool     `long:"all" description:"generate All function"`
	Some    bool     `long:"some" description:"generate Some function"`
	ForEach bool     `long:"foreach" description:"generate ForEach function"`
	Concat  bool     `long:"concat" description:"generate Concat function"`
	Find    bool     `long:"find" description:"generate Find function"`
	Reverse bool     `long:"reverse" description:"generate Reverse function"`
	Splice  bool     `long:"splice" description:"generate Splice function"`
	Reduce  []string `long:"reduce" description:"generate Reduce function for given type"`
	Array   bool     `long:"array" description:"generate Array function for channel type"`

	Type        TypeDef
	MapResults  []TypeDef
	ReduceTypes []TypeDef
}

// TypeDef is a type definition, with name, package and type
type TypeDef struct {
	Name    string
	Package string
	Type    string
	IsChan  bool
}

type generatorFunc func(io.Writer) error

func (g *Generator) parseTypes() error {
	td, err := g.parseType(g.RawType)
	if err != nil {
		return err
	}
	g.Type = td

	for _, m := range g.Map {
		td, err := g.parseType(m)
		if err != nil {
			return err
		}

		g.MapResults = append(g.MapResults, td)
	}

	for _, r := range g.Reduce {
		td, err := g.parseType(r)
		if err != nil {
			return err
		}

		g.ReduceTypes = append(g.ReduceTypes, td)
	}

	return nil
}

func (g *Generator) parseType(raw string) (TypeDef, error) {
	var (
		t   TypeDef
		err error
	)

	t.Package, t.Type, t.IsChan, err = g.parseRawType(raw)
	t.Name = g.getTypeName(t.Type)

	return t, err
}

func (g *Generator) parseRawType(raw string) (string, string, bool, error) {
	var (
		pkg    string
		typ    string
		isChan bool
		err    error
	)

	typeParts := strings.Split(raw, ":")
	if len(typeParts) == 1 {
		typ, isChan, err = parseType(raw)
	} else {
		pkg = typeParts[0]
		typ, isChan, err = parseType(typeParts[1])
	}

	return pkg, typ, isChan, err
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
	pkgs := map[string]struct{}{
		"errors": struct{}{},
	}

	if g.Type.Package != "" {
		pkgs[g.Type.Package] = struct{}{}
	}

	for _, mr := range g.MapResults {
		if mr.Package != "" {
			pkgs[mr.Package] = struct{}{}
		}
	}

	if g.Type.IsChan && g.Concat {
		pkgs["sync"] = struct{}{}
	}

	var packages []string
	for pkg := range pkgs {
		packages = append(packages, pkg)
	}

	tpl, err := g.getTpl(importsTpl)
	if err != nil {
		return err
	}

	sort.Strings(packages)
	return tpl.Execute(w, packages)
}

func (g *Generator) generateType(w io.Writer) error {
	tpl, err := g.getTpl(typeTpl)
	if err != nil {
		return err
	}
	return tpl.Execute(w, g.Type)
}

func (g *Generator) generateSome(w io.Writer) error {
	if g.Some {
		if g.Type.IsChan {
			return errors.New("chan type does not support some")
		}

		tpl, err := g.getTpl(someTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateAll(w io.Writer) error {
	if g.All {
		if g.Type.IsChan {
			return errors.New("chan type does not support all")
		}

		tpl, err := g.getTpl(allTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateFilter(w io.Writer) error {
	if g.Filter {
		tpl, err := g.getTpl(filterTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateMap(w io.Writer) error {
	if len(g.Map) > 0 {
		tpl, err := g.getTpl(mapTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateMapResults(w io.Writer) error {
	if len(g.Map) > 0 {
		data := struct {
			Name    string
			Results []TypeDef
		}{
			Name:    g.Type.Name,
			Results: g.MapResults,
		}

		tpl, err := g.getTpl(mapResultsTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, data)
	}

	return nil
}

func (g *Generator) generateForEach(w io.Writer) error {
	if g.ForEach {
		tpl, err := g.getTpl(forEachTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateConcat(w io.Writer) error {
	if g.Concat {
		tpl, err := g.getTpl(concatTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateFind(w io.Writer) error {
	if g.Find {
		if g.Type.IsChan {
			return errors.New("chan iter does not support find")
		}

		tpl, err := g.getTpl(findTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateReverse(w io.Writer) error {
	if g.Reverse {
		if g.Type.IsChan {
			return errors.New("chan iter does not support reverse")
		}

		tpl, err := g.getTpl(reverseTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateSplice(w io.Writer) error {
	if g.Splice {
		if g.Type.IsChan {
			return errors.New("a chan iter does not support splice")
		}

		tpl, err := g.getTpl(spliceTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}
	return nil
}

func (g *Generator) generateReduces(w io.Writer) error {
	if len(g.Reduce) > 0 {
		data := struct {
			Name     string
			Type     string
			Reducers []TypeDef
		}{
			Name:     g.Type.Name,
			Type:     g.Type.Type,
			Reducers: g.ReduceTypes,
		}

		tpl, err := g.getTpl(reduceTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, data)
	}

	return nil
}

func (g *Generator) generateArray(w io.Writer) error {
	if g.Array {
		if !g.Type.IsChan {
			return errors.New("array is not supported on non-channel types")
		}

		tpl, err := g.getTpl(arrayTpl)
		if err != nil {
			return err
		}
		return tpl.Execute(w, g.Type)
	}

	return nil
}

func (g *Generator) getTpl(tpl string) (*template.Template, error) {
	return getTemplate(tpl, g.Type.IsChan)
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
		g.generateForEach,
		g.generateConcat,
		g.generateFind,
		g.generateReverse,
		g.generateSplice,
		g.generateReduces,
		g.generateArray,
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

const (
	fileTpl     = "%s_iter.go"
	chanFileTpl = "%schan_iter.go"
)

func (g *Generator) fileName() string {
	tpl := fileTpl
	if g.Type.IsChan {
		tpl = chanFileTpl
	}

	file := fmt.Sprintf(tpl, fileify(g.Type.Type))
	return filepath.Join(".", file)
}

// Generate writes the generated code to the correspondant file and returns an error if something failed
func (g *Generator) Generate() error {
	err := g.parseTypes()
	if err != nil {
		return err
	}

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

func parseType(t string) (string, bool, error) {
	if strings.Contains(t, "<-") {
		return "", false, fmt.Errorf("invalid channel type given: %s", t)
	}

	tParts := strings.Split(t, " ")
	return tParts[len(tParts)-1], strings.HasPrefix(t, "chan"), nil
}
