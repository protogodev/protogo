package protogocmd

import (
	"log"

	"github.com/alecthomas/kong"
)

var CLI struct{}

func Main() {
	ctx := kong.Parse(&CLI, KongOptions()...)
	if err := ctx.Run(); err != nil {
		log.Fatal(err)
	}
}
