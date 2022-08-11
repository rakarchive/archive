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

package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// Unlocker implements an io.ReaderAt which decrypts the data provided to
// it using the given key and returns the cleartext on Read calls.
type Unlocker struct {
	Password []byte
	Source   io.ReadCloser

	io.ReaderAt

	Size int64
}

func (u *Unlocker) Unlock() error {
	defer u.Source.Close()

	salt := make([]byte, 8)
	if _, err := u.Source.Read(salt); err != nil {
		return err
	}

	key := DeriveKey(u.Password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = u.Source.Read(nonce); err != nil {
		return err
	}

	ct, err := io.ReadAll(u.Source)
	if err != nil {
		return err
	}

	clt, err := aead.Open(nil, nonce, ct, nil)

	u.Size = int64(len(clt))
	u.ReaderAt = bytes.NewReader(clt)

	return err
}
