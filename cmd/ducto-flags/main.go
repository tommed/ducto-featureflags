// File: ducto-featureflags/cmd/ducto-flags/main.go
package main

import (
	"github.com/tommed/ducto-featureflags/internal/cli"
	"os"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
