package protogocmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	Plugins     []string `name:"plugin" help:"Plugins to build. Format: <module[=replacement]>"`
	SkipCleanup int      `env:"PROTOGO_SKIP_CLEANUP" default:"0" help:"Whether to leave build artifacts on disk after exiting."`
}

func (b *Build) Run(ctx *kong.Context) error {
	mods, err := toModules(b.Plugins)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir(wd, "build-")
	if err != nil {
		return err
	}
	if b.SkipCleanup == 0 {
		defer os.RemoveAll(dir)
	}

	mainDir := filepath.Join(dir, "protogo")
	if err := os.Mkdir(mainDir, os.ModePerm); err != nil {
		return err
	}

	mainFile := filepath.Join(mainDir, "main.go")
	if err := b.genMain(mainFile, mods.Paths()); err != nil {
		return err
	}

	var cmds commands
	cmds.Add("go", "mod", "init")
	for _, replace := range mods.ReplacementDirectives() {
		cmds.Add("go", "mod", "edit", "-replace", replace)
	}
	cmds.Add("go", "mod", "tidy")
	cmds.Add("go", "build")
	cmds.Add("cp", "protogo", wd)
	return cmds.Run(mainDir)
}

func (b *Build) genMain(filename string, paths []string) error {
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
	file, err := generator.Generate(tmpl, paths, generator.Options{
		Formatted:      true,
		TargetFileName: filename,
	})
	if err != nil {
		return err
	}
	return file.Write()
}

type command struct {
	name string
	args []string
}

type commands []command

func (cs *commands) Add(name string, args ...string) {
	*cs = append(*cs, command{
		name: name,
		args: args,
	})
}

func (cs commands) Run(dir string) error {
	for _, c := range cs {
		cmd := exec.Command(c.name, c.args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("err: %s", out)
		}
	}
	return nil
}

type module struct {
	path        string
	replacement string
}

type modules []module

func toModules(plugins []string) (modules, error) {
	var ms modules
	for _, p := range plugins {
		parts := strings.SplitN(p, "=", 2)
		m := module{path: parts[0]}

		if len(parts) == 2 {
			replacement, err := filepath.Abs(parts[1])
			if err != nil {
				return nil, err
			}
			m.replacement = replacement
		}

		ms = append(ms, m)
	}
	return ms, nil
}

func (ms modules) Paths() (p []string) {
	for _, m := range ms {
		p = append(p, m.path)
	}
	return
}

func (ms modules) ReplacementDirectives() (r []string) {
	for _, m := range ms {
		if m.replacement != "" {
			directive := m.path + "=" + m.replacement
			r = append(r, directive)
		}
	}
	return
}
