package main

import (
	"os"

	"github.com/feloy/tesh/pkg/cmd"
	"github.com/spf13/pflag"
)

func main() {
	flags := pflag.NewFlagSet("tesh", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewTesh()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
