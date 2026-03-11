package main

import (
	"os"

	"github.com/misonikomipan/homebox-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
