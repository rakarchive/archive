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
	"crypto/sha256"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/chacha20"
)

// ErrAuth represents the error returned by (*Unlocker).Close if the
// verification of the message's authenticity fails.
var ErrAuth = errors.New("unlocker: message authentication failed")

// Unlocker implements an AE which decrypts and authenticates the
// provided data, and writes it to the given destination.
type Unlocker struct {
	// user input
	Password    []byte
	Destination io.Writer

	rand []byte

	// unlocker state
	hash    hash.Hash
	cipher  *chacha20.Cipher
	closed  bool
	final32 []byte
}

// Write decrypts the given data and writes it to the destination.
func (u *Unlocker) Write(b []byte) (int, error) {
	// (*Unlocker).Close has already been called
	if u.closed {
		panic("(*Unlocker).Write(): unlocker is closed")
	}

	// if all the random bytes(salt and nonce) have not been received,
	// keep collecting them
	if i := randLen - len(u.rand); i > 0 {
		var rand []byte

		if i >= len(b) {
			// collect all
			rand = b
			b = []byte{}
		} else {
			// collect some
			rand = b[:i]
			b = b[i:]
		}

		u.rand = append(u.rand, rand...)
	}

	// first call to write, initialize cipher
	if u.cipher == nil {
		// extract salt
		salt := u.rand[:saltLen]
		key := DeriveKey(u.Password, salt)

		// extract nonce
		nonce := u.rand[saltLen:]
		u.cipher, _ = chacha20.NewUnauthenticatedCipher(key, nonce)

		u.final32 = make([]byte, sha256.Size)
		u.hash = hmac.New(sha256.New, key)
	}

	// keep track of last 32 bytes
	buffer := append(u.final32, b...)
	i := len(buffer) - sha256.Size
	u.final32 = buffer[i:]
	b = buffer[:i]

	u.cipher.XORKeyStream(b, b)   // decrypt data
	u.hash.Write(b)               // update hash
	return u.Destination.Write(b) // write to destination
}

// Close authenticates the decrypted message. It returns ErrAuth if the
// authentication fails due to any reason.
func (u *Unlocker) Close() error {
	// decrypt hash
	u.cipher.XORKeyStream(u.final32, u.final32)

	// verify hash
	if !hmac.Equal(u.hash.Sum(nil), u.final32) {
		return ErrAuth
	}
	return nil
}
