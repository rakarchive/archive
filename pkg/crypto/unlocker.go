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
	"errors"
	"io"
)

// NewUnlocker creates a new *unlocker with the given source and key.
func NewUnlocker(src io.Reader, key []byte) (*unlocker, error) {
	var u unlocker
	u.src = src

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	u.aead, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ct, err := io.ReadAll(u.src)
	if err != nil {
		return nil, err
	}

	nonceSize := u.aead.NonceSize()
	if len(ct) < nonceSize {
		return nil, errors.New("unlocker: nonce not found")
	}

	nonce, ct := ct[:nonceSize], ct[nonceSize:]
	clt, err := u.aead.Open(nil, nonce, ct, nil)

	u.Size = int64(len(clt))
	u.ReaderAt = bytes.NewReader(clt)

	return &u, err
}

// unlocker implements an io.ReaderAt which decrypts the data provided to
// it using the given key and returns the cleartext on Read calls.
type unlocker struct {
	io.ReaderAt
	src io.Reader

	Size int64

	aead cipher.AEAD
}
