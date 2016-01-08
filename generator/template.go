package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

var tpls = map[string]*template.Template{
	"type":        loadTemplate("type"),
	"imports":     loadTemplate("imports"),
	"map":         loadTemplate("map"),
	"map_results": loadTemplate("map_results"),
	"filter":      loadTemplate("filter"),
	"some":        loadTemplate("some"),
	"all":         loadTemplate("all"),
	"foreach":     loadTemplate("foreach"),
	"concat":      loadTemplate("concat"),
	"find":        loadTemplate("find"),
	"reverse":     loadTemplate("reverse"),
	"splice":      loadTemplate("splice"),
	"reduce":      loadTemplate("reduce"),

	"chan_type":        loadTemplate("chan_type"),
	"chan_concat":      loadTemplate("chan_concat"),
	"chan_filter":      loadTemplate("chan_filter"),
	"chan_map":         loadTemplate("chan_map"),
	"chan_map_results": loadTemplate("chan_map_results"),
	"chan_imports":     loadTemplate("imports"),
	"chan_foreach":     loadTemplate("chan_foreach"),
	"chan_reduce":      loadTemplate("chan_reduce"),
	"chan_array":       loadTemplate("chan_array"),
}

const (
	typeTpl       = "type"
	importsTpl    = "imports"
	mapTpl        = "map"
	mapResultsTpl = "map_results"
	filterTpl     = "filter"
	someTpl       = "some"
	allTpl        = "all"
	forEachTpl    = "foreach"
	concatTpl     = "concat"
	findTpl       = "find"
	reverseTpl    = "reverse"
	spliceTpl     = "splice"
	reduceTpl     = "reduce"
	arrayTpl      = "array"
)

func loadTemplateText(name string) string {
	f := filepath.Join(os.Getenv("GOPATH"), "src/github.com/mvader/go-itergen/generator/templates", fmt.Sprintf("%s.tgo", name))

	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func loadTemplate(name string) *template.Template {
	text := loadTemplateText(name)
	return template.Must(template.New(name).Parse(text))
}

func getTemplate(name string, isChan bool) (*template.Template, error) {
	if isChan {
		name = "chan_" + name
	}

	tpl, ok := tpls[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}

	return tpl, nil
}
