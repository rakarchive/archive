package archive

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"laptudirm.com/x/archive/crypto"
	"laptudirm.com/x/archive/zipper"
)

const Ext = ".march"

func IsArchive(path string) bool {
	info, err := os.Stat(path + Ext)
	return !errors.Is(err, os.ErrNotExist) && !info.IsDir()
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist) && info.IsDir()
}

type Archive struct {
	Name string

	Pass []byte
	Salt []byte

	SrcDirectory string
	ArcDirectory string
}

func (a *Archive) Create() error {
	dst := a.archivePath()

	arc, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer arc.Close()

	l := &crypto.Locker{
		Password:    a.Pass,
		Destination: arc,
	}

	fmt.Print("zipping files... ")
	if err = zipper.Zip(a.SrcDirectory, l); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("encrypting archive... ")
	if err := l.Close(); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("removing old files... ")
	if err = os.RemoveAll(a.SrcDirectory); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Printf("\nSuccessfully created archive '%s'.\n", a.Name)
	return nil
}

func (a *Archive) Open() error {
	archive := a.archivePath()

	arc, err := os.Open(archive)
	if err != nil {
		return err
	}

	fmt.Print("decrypting files... ")
	u := &crypto.Unlocker{
		Password: a.Pass,
		Source:   arc,
	}
	if err = u.Unlock(); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("unzipping files... ")
	if err := zipper.Unzip(u, u.Size, a.ArcDirectory); err != nil {
		return err
	}
	fmt.Println("done.")

	fmt.Print("removing archive... ")
	if err = os.Remove(archive); err != nil {
		return err
	}
	fmt.Print("done.\n\n")

	fmt.Printf("Successfully unarchived '%s'.\n", a.Name)
	return nil
}

func (a *Archive) archivePath() string {
	return filepath.Join(a.ArcDirectory, a.Name+Ext)
}
