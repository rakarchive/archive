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

// NewLocker creates a new *locker with the given key and destination.
func NewLocker(dst io.Writer, key []byte) (*locker, error) {
	var l locker
	l.dst = dst

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	l.aead, err = cipher.NewGCM(block)
	return &l, err
}

// locker implements an io.Writer which encrypts the incoming data with the
// given key and writes it to the provided io.Writer.
type locker struct {
	dst io.Writer

	src  []byte
	aead cipher.AEAD
}

func (l *locker) Write(b []byte) (int, error) {
	l.src = append(l.src, b...)
	return len(b), nil
}

// Close signals that all the data that needs to be encrypted has been
// provided to the locker. It encrypts the collected data and writes it
// to the destination.
func (l *locker) Close() error {
	nonce := make([]byte, l.aead.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return err
	}

	enc := l.aead.Seal(nonce, nonce, l.src, nil)
	_, err = l.dst.Write(enc)
	return err
}
