package cmd

import (
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"laptudirm.com/x/archive/pkg/archive"
)

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open <archive>",
		Short: "open the given sealed archive",
		Args:  cobra.ExactArgs(1),
		RunE:  open,
	}
}

func open(cmd *cobra.Command, args []string) error {
	arc := args[0]

	if !archive.IsArchive(arc) {
		return fmt.Errorf("open: '%s' is not an archive", arc)
	}

	fmt.Print("Enter a password for the archive: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n\n")
	if err != nil {
		return err
	}

	a := &archive.Archive{
		Name: filepath.Base(arc),

		Pass: pass,

		SrcDirectory: arc,
		ArcDirectory: filepath.Dir(arc),
	}
	return a.Open()
}
