package main

import (
	"github.com/HellstromIT/go-quickswitch/cmd/go-quickswitch/internal/quickswitch"
)

var version = "dev"

func main() {
	quickswitch.Cli(version)
}
