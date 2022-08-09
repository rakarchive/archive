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
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

// DeriveKey derives a 32 byte key from the given password, using 4096
// iterations of the pbkdf2 algorithm with sha256.
func DeriveKey(pass []byte) []byte {
	iter := 4096     // no of pbkdf2 iterations
	klen := 32       // length of key in bytes
	salt := []byte{} // no salt

	algo := sha256.New // hash algorithm for pbkdf

	return pbkdf2.Key(pass, salt, iter, klen, algo)
}
