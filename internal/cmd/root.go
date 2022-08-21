package cmd

import "github.com/spf13/cobra"

func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "archive",
		Args: cobra.NoArgs,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// global flags
	cmd.PersistentFlags().BoolP("help", "h", false, "show help for command")

	// set custom help
	hc := helpCmd()
	cmd.SetHelpCommand(hc)
	cmd.SetHelpFunc(hc.Run)

	// add commands
	cmd.AddCommand(sealCmd())
	cmd.AddCommand(openCmd())
	cmd.AddCommand(completionCmd())

	return cmd
}
