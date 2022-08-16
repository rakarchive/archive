package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func completionCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "completion",
		Hidden: true,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stderr, cmd.UsageString())
		},
	}
}
