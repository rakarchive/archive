package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"laptudirm.com/x/archive/pkg/archive"
)

func sealCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seal <directory>",
		Short: "create and seal a new archive",
		Args:  cobra.ExactArgs(1),
		RunE:  seal,
	}
}

func seal(cmd *cobra.Command, args []string) error {
	dir := args[0]

	if !archive.IsDirectory(dir) {
		return fmt.Errorf("seal: '%s' is not a directory", dir)
	}

	fmt.Print("Enter a password for the archive: ")
	pass1, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return err
	}

	fmt.Print("Confirm the password for the archive:")
	pass2, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n\n")
	if err != nil {
		return err
	}

	if string(pass1) != string(pass2) {
		return errors.New("seal: passwords do not match")
	}

	a := &archive.Archive{
		Name: filepath.Base(dir),

		Pass: pass1,

		SrcDirectory: dir,
		ArcDirectory: filepath.Dir(dir),
	}
	return a.Create()
}
