package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

var (
	typeTpl       = loadTemplate("type")
	importsTpl    = loadTemplate("imports")
	mapTpl        = loadTemplate("map")
	mapResultsTpl = loadTemplate("map_results")
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
