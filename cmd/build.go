package protogocmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alecthomas/kong"

	"github.com/protogodev/protogo/generator"
)

func init() {
	MustRegister(&Plugin{
		Name: "build",
		Help: "Build protogo CLI with plugins",
		Cmd:  &Build{},
	})
}

type Build struct {
	Plugins []string `name:"plugin" help:"Plugins to build."`
}

func (b *Build) Run(ctx *kong.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir, err := ioutil.TempDir(wd, "build-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	mainDir := filepath.Join(dir, "protogo")
	if err := os.Mkdir(mainDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	mainFile := filepath.Join(mainDir, "main.go")
	if err := b.genMain(mainFile); err != nil {
		log.Fatal(err)
	}

	run := func(name string, arg ...string) {
		cmd := exec.Command(name, arg...)
		cmd.Dir = mainDir
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("err: %s", out)
		}
	}

	run("go", "mod", "init")
	run("go", "mod", "tidy")
	run("go", "build")
	run("cp", "protogo", wd)

	return nil
}

func (b *Build) genMain(filename string) error {
	tmpl := `package main

import (
	protogocmd "github.com/protogodev/protogo/cmd"

	{{range . -}}
	_ "{{.}}"
	{{end}}
)

func main() {
	protogocmd.Main()
}
`
	file, err := generator.Generate(tmpl, b.Plugins, generator.Options{
		Formatted:      true,
		TargetFileName: filename,
	})
	if err != nil {
		return err
	}
	return file.Write()
}
