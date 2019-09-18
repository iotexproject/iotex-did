package main

import (
	"os"

	"github.com/iotexproject/iotex-DID/cmd"
)

// main runs the did command
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
