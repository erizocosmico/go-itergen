package generator

// Generator generates functions for iterable types based on the options received
type Generator struct {
	Type   string   `short:"t" long:"type" description:"type to generate the code for"`
	Map    []string `long:"map" description:"generate Map function with transformer for given type"`
	Filter bool     `long:"filter" description:"generate Filter function"`
	All    bool     `long:"all" description:"generate All function"`
	Some   bool     `long:"some" description:"generate Some function"`
}

// Generate writes the generated code to the correspondant file and returns an error if something failed
func (g *Generator) Generate() error {
	return nil
}
