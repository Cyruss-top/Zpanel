package main

import (
	"github.com/zex/zpanel/internal/cli"
)

// version 由 Makefile -ldflags 注入
var version = "0.5.0-dev"

func main() {
	cli.Execute(version)
}
