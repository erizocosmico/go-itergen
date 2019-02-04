package main

import (
	"fmt"
	"os"

	"github.com/erizocosmico/go-itergen"
	_ "github.com/erizocosmico/go-itergen/statik"
	"github.com/jessevdk/go-flags"
)

func main() {
	cmd := new(generator.Generator)
	parser := flags.NewParser(cmd, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}

	err = cmd.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
