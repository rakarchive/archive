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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Locker implements an io.Writer which encrypts the incoming data with the
// given key and writes it to the provided io.Writer.
type Locker struct {
	Password    []byte
	Destination io.Writer

	data []byte
}

func (l *Locker) Write(b []byte) (int, error) {
	l.data = append(l.data, b...)
	return len(b), nil
}

// Close signals that all the data that needs to be encrypted has been
// provided to the locker. It encrypts the collected data and writes it
// to the destination.
func (l *Locker) Close() error {
	salt, err := randBytes(8)
	if err != nil {
		return err
	}

	key := DeriveKey(l.Password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce, err := randBytes(aead.NonceSize())
	if err != nil {
		return err
	}

	data := append(salt, nonce...)
	enc := aead.Seal(data, nonce, l.data, nil)

	_, err = l.Destination.Write(enc)
	return err
}

func randBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}
