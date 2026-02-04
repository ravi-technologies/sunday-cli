// Package main is the entry point for the Sunday CLI application.
//
// The application is built with the API base URL injected at build time:
//
//	make build API_URL=https://api.sunday.example.com
//
// Run with --help to see available commands.
package main

import (
	"fmt"
	"os"

	"github.com/ravi-technologies/sunday-cli/internal/output"
	"github.com/ravi-technologies/sunday-cli/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		output.Current.PrintError(err)
		os.Exit(1)
	}
	fmt.Println()
}
