// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
	"laptudirm.com/x/archive/pkg/crypto"
	"laptudirm.com/x/archive/pkg/zipper"
)

func main() {
	if err := mainF(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mainF() error {
	if len(os.Args) != 2 {
		return errors.New("usage: archive [directory]")
	}

	dir, err := filepath.Abs(os.Args[1])
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
	key := crypto.DeriveKey(pass)

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	l, err := crypto.NewLocker(w, key)
	if err != nil {
		return err
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
