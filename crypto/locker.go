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
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"

	"golang.org/x/crypto/chacha20"
)

// Locker implements an AE which encrypts the incoming data with the given
// key using Xchacha20 along with a MAC and writes it to the provided
// io.Writer. The written data can be unencrypted and authenticated easily.
type Locker struct {
	// user inputs
	Password    []byte
	Destination io.Writer

	// locker state
	hash   hash.Hash
	cipher *chacha20.Cipher
	closed bool
}

// Write allows locker to adhere to the io.Writer interface. The provided
// slice is encrypted with Xchacha20 and written to the destination.
func (l *Locker) Write(b []byte) (int, error) {
	// (*Locker).Close has already been called
	if l.closed {
		panic("(*Locker).Write(): locker is closed")
	}

	// first call to write, initialize cipher
	if l.cipher == nil {
		// generate salt and nonce
		rand, err := randBytes(randLen)
		if err != nil {
			return 0, err
		}

		// write the salt and nonce to the destination
		if n, err := l.Destination.Write(rand); err != nil {
			return n, err
		}

		salt := rand[:saltLen] // extract salt from bytes
		key := DeriveKey(l.Password, salt)

		nonce := rand[saltLen:] // extract nonce from bytes
		l.cipher, _ = chacha20.NewUnauthenticatedCipher(key, nonce)

		l.hash = hmac.New(sha256.New, key)
	}

	l.hash.Write(b)               // update hash value
	l.cipher.XORKeyStream(b, b)   // encrypt the data
	return l.Destination.Write(b) // write ciphertext
}

// Close signals that all the data that needs to be encrypted has been
// provided to the locker. The locker writes the MAC to the destination.
func (l *Locker) Close() error {
	_, err := l.Write(l.hash.Sum(nil))
	l.closed = true // panic on write calls
	return err
}

// randBytes generates n cryptographically secure random bytes.
func randBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}
