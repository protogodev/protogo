package protogocmd

import (
	"log"

	"github.com/alecthomas/kong"
)

var CLI struct{}

func Main() {
	ctx := kong.Parse(&CLI, KongOptions()...)
	log.Fatal(ctx.Run())
}
