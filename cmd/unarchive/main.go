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
		fmt.Fprintln(os.Stderr, "usage: unarchive [name]")
	}

	archive, err := filepath.Abs(os.Args[1] + ".march")
	if err != nil {
		return err
	}

	fmt.Print("Enter a password for the archive: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n\n")
	if err != nil {
		return err
	}
	key := crypto.DeriveKey(pass)

	w, err := os.Open(archive)
	if err != nil {
		return err
	}

	fmt.Print("decrypting files... ")
	u, err := crypto.NewUnlocker(w, key)
	if err != nil {
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

	fmt.Printf("Successfully unarchived '%s'.", filepath.Base(archive))
	return nil
}
