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

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open <archive>",
		Short: "open the given sealed archive",
		Args:  cobra.ExactArgs(1),
		RunE:  open,
	}
}

func open(cmd *cobra.Command, args []string) error {
	archive, err := filepath.Abs(args[0] + ".march")
	if err != nil {
		return err
	}

	fmt.Print("Enter the password for the archive: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n\n")
	if err != nil {
		return err
	}

	r, err := os.Open(archive)
	if err != nil {
		return err
	}

	fmt.Print("decrypting files... ")
	u := &crypto.Unlocker{
		Password: pass,
		Source:   r,
	}
	if err = u.Unlock(); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("unzipping files... ")
	if err := zipper.Unzip(u, u.Size, filepath.Dir(archive)); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("removing archive... ")
	if err = os.Remove(archive); err != nil {
		return err
	}
	fmt.Print("done.\n\n")

	fmt.Printf("Successfully unarchived '%s'.\n", filepath.Base(archive))
	return nil
}
