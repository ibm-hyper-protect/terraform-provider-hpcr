// Copyright 2023 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package encrypt

import (
	"fmt"
	"math/rand"
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func testSymmetricDecrypt(
	t *testing.T,
	srcLen int,
	encrypt func(srcPlainBytes []byte) func([]byte) E.Either[error, string],
	decrypt func(token string) func([]byte) E.Either[error, []byte]) E.Either[error, bool] {

	// generate some random data
	dataGen := cryptoRandomE(srcLen)
	dataE := dataGen()

	// generate some password
	pwdGen := CryptoRandomPassword(keylen)
	pwdE := pwdGen()

	// encrypt the data record
	encE := F.Pipe3(
		dataE,
		E.Map[error](encrypt),
		E.Ap[E.Either[error, string]](pwdE),
		E.Flatten[error, string],
	)

	// decrypt the data record
	decE := F.Pipe3(
		encE,
		E.Map[error](decrypt),
		E.Ap[E.Either[error, []byte]](pwdE),
		E.Flatten[error, []byte],
	)

	// verify that the encrypted and the decrypted data match
	return E.Sequence2(func(exp, actual []byte) E.Either[error, bool] {
		return E.Of[error](assert.Equal(t, exp, actual))
	})(dataE, decE)
}

func testAsymmetricDecrypt(
	t *testing.T,
	srcLen int,
	encrypt func([]byte) func([]byte) E.Either[error, string],
	decrypt func([]byte) func(string) E.Either[error, []byte]) E.Either[error, bool] {

	// generate some random data
	dataGen := cryptoRandomE(srcLen)
	dataE := dataGen()

	// generate a key pair
	privKeyE := CryptoPrivateKey()
	pubKeyE := F.Pipe1(
		privKeyE,
		E.Chain(CryptoPublicKey),
	)

	// encrypt the data record
	encE := F.Pipe3(
		pubKeyE,
		E.Map[error](encrypt),
		E.Ap[E.Either[error, string]](dataE),
		E.Flatten[error, string],
	)

	// decrypt the data record
	decE := F.Pipe3(
		privKeyE,
		E.Map[error](decrypt),
		E.Ap[E.Either[error, []byte]](encE),
		E.Flatten[error, []byte],
	)

	// verify that the encrypted and the decrypted data match
	return E.Sequence2(func(exp, actual []byte) E.Either[error, bool] {
		return E.Of[error](assert.Equal(t, exp, actual))
	})(dataE, decE)
}

type (
	// definition of an encryption
	SymmEncDecItem struct {
		Encrypt func(srcPlainBytes []byte) func([]byte) E.Either[error, string]
		Decrypt func(token string) func([]byte) E.Either[error, []byte]
	}

	// definition of an encryption
	AsymmEncDecItem struct {
		Encrypt func([]byte) func([]byte) E.Either[error, string]
		Decrypt func([]byte) func(string) E.Either[error, []byte]
	}
)

var (
	// test matrix for symmetric encryption combinations
	SymmEncDecMatrix = []SymmEncDecItem{
		{Encrypt: OpenSSLSymmetricEncrypt, Decrypt: OpenSSLSymmetricDecrypt},
		{Encrypt: CryptoSymmetricEncrypt, Decrypt: OpenSSLSymmetricDecrypt},
		{Encrypt: OpenSSLSymmetricEncrypt, Decrypt: CryptoSymmetricDecrypt},
		{Encrypt: CryptoSymmetricEncrypt, Decrypt: CryptoSymmetricDecrypt},
	}

	// test matrix for asymmetric encryption combinations
	AsymmEncDecMatrix = []AsymmEncDecItem{
		{Encrypt: OpenSSLAsymmetricEncryptPub, Decrypt: OpenSSLAsymmetricDecrypt},
		{Encrypt: CryptoAsymmetricEncryptPub, Decrypt: OpenSSLAsymmetricDecrypt},
		{Encrypt: OpenSSLAsymmetricEncryptPub, Decrypt: CryptoAsymmetricDecrypt},
		{Encrypt: CryptoAsymmetricEncryptPub, Decrypt: CryptoAsymmetricDecrypt},
	}
)

// TestSymmetricDecrypt checks if the symmetric decryption works
func TestSymmetricDecrypt(t *testing.T) {

	for i := 0; i < 10; i++ {

		len := rand.Intn(10000) + 1
		for idx, item := range SymmEncDecMatrix {

			t.Run(fmt.Sprintf("Message Size [%d], Combination [%d]", len, idx), func(t *testing.T) {
				res := testSymmetricDecrypt(t, len, item.Encrypt, item.Decrypt)
				assert.Equal(t, E.Of[error](true), res)
			})
		}

	}
}

func TestAsymmetricDecrypt(t *testing.T) {

	for i := 0; i < 2; i++ {

		len := rand.Intn(255) + 1
		for idx, item := range AsymmEncDecMatrix {

			t.Run(fmt.Sprintf("Message Size [%d], Combination [%d]", len, idx), func(t *testing.T) {
				res := testAsymmetricDecrypt(t, len, item.Encrypt, item.Decrypt)
				assert.Equal(t, E.Of[error](true), res)
			})
		}

	}
}
