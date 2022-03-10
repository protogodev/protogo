package protogocmd

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Plugin struct {
	Name  string
	Help  string
	Group string
	Cmd   interface{}
	Tags  []string
}

type Registry map[string]*Plugin

func (r Registry) Register(plugin *Plugin) error {
	if _, ok := r[plugin.Name]; ok {
		return fmt.Errorf("cmd %q is already registered", plugin.Name)
	}
	r[plugin.Name] = plugin
	return nil
}

func (r Registry) MustRegister(cmd *Plugin) {
	if err := r.Register(cmd); err != nil {
		panic(err)
	}
}

func (r Registry) KongOptions() []kong.Option {
	var opts []kong.Option
	for _, cmd := range r {
		opts = append(opts, kong.DynamicCommand(
			cmd.Name,
			cmd.Help,
			cmd.Group,
			cmd.Cmd,
			cmd.Tags...,
		))
	}
	return opts
}

func MustRegister(cmd *Plugin) {
	globalRegistry.MustRegister(cmd)
}

func KongOptions() []kong.Option {
	return globalRegistry.KongOptions()
}

var globalRegistry = Registry{}
