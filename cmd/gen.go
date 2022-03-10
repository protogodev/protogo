package protogocmd

import (
	"path/filepath"

	"github.com/alecthomas/kong"

	"github.com/protogodev/protogo/generator"
	"github.com/protogodev/protogo/parser"
	"github.com/protogodev/protogo/parser/ifacetool"
)

type Generator interface {
	Generate(data *ifacetool.Data) (*generator.File, error)
}

type Gen struct {
	Generator

	SrcFilename   string `arg:"" name:"source-file" help:"source-file"`
	InterfaceName string `arg:"" name:"interface-name" help:"interface-name"`
}

func NewGen(generator Generator) *Gen {
	return &Gen{Generator: generator}
}

func (g *Gen) Run(ctx *kong.Context) error {
	srcFilename, err := filepath.Abs(g.SrcFilename)
	if err != nil {
		return err
	}

	// Non-empty pkgName makes all type names used in the interface full-qualified.
	pkgName := "x"
	data, err := parser.ParseInterface(pkgName, srcFilename, g.InterfaceName)
	if err != nil {
		return err
	}

	file, err := g.Generator.Generate(data)
	if err != nil {
		return err
	}

	return file.Write()
}
