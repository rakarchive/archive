package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"laptudirm.com/x/archive/pkg/crypto"
	"laptudirm.com/x/archive/pkg/zipper"
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
	dir, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	dst := filepath.Base(dir) + ".march"

	fmt.Print("Enter a password for the archive: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n\n")
	if err != nil {
		return err
	}

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	l := &crypto.Locker{
		Password:    pass,
		Destination: w,
	}

	fmt.Print("zipping files... ")
	if err = zipper.Zip(dir, l); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("encrypting archive... ")
	if err := l.Close(); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("removing old files... ")
	if err = os.RemoveAll(dir); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Println("\nSuccessfully created archive.")
	return nil
}
