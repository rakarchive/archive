package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"laptudirm.com/x/archive/internal/cmd"
)

func main() {
	root := cmd.Root()
	if cmd, err := root.ExecuteC(); err != nil {
		printError(cmd, err)
		os.Exit(1)
	}
}

func printError(cmd *cobra.Command, err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, cmd.UsageString())
}
