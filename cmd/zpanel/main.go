package main

import (
	"fmt"
	"os"
)

// version 由 Makefile -ldflags 注入；开发期默认 dev
var version = "0.2.0-dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "server" {
		fmt.Println("zpanel server: not implemented yet (next step)")
		os.Exit(1)
	}
	fmt.Printf("zpanel %s\n", version)
}
