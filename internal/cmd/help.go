package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func helpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help <command>",
		Short: "get help about any command",
		Args:  cobra.ArbitraryArgs,
		RunE:  help,
	}
}

func help(cmd *cobra.Command, args []string) error {
	root := cmd.Root()
	target, extra, err := root.Traverse(args)
	if err != nil {
		return err
	}

	if len(extra) > 0 {
		return fmt.Errorf("unknown command \"%v\" for \"%v\"", extra[0], getFullName(target))
	}

	if target.Long != "" {
		fmt.Println(target.Long)
	} else {
		fmt.Println(target.Short)
	}
	fmt.Println()

	if err = target.Usage(); err != nil {
		return err
	}

	fmt.Println()
	return nil
}

func getFullName(cmd *cobra.Command) string {
	name := cmd.Name()

	if cmd.HasParent() {
		name = getFullName(cmd.Parent()) + " " + name
	}

	return name
}
