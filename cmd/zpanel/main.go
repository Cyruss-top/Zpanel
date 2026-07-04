package main

import (
	"fmt"
	"os"

	"github.com/zex/zpanel/internal/app"
)

// version 由 Makefile -ldflags 注入；开发期默认 dev
var version = "0.2.0-dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "server" {
		if err := app.RunServer(app.Options{Version: version}); err != nil {
			fmt.Fprintf(os.Stderr, "zpanel server: %v\n", err)
			os.Exit(1)
		}
		return
	}
	fmt.Printf("zpanel %s\n", version)
}
