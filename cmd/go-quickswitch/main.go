package main

import (
	cli "github.com/HellstromIT/go-quickswitch/internal/quickswitch"
)

var version = "dev"

func main() {
	cli.Cli(version)
}
