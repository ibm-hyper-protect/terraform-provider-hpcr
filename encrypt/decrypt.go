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
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
)

// Decryption captures the crypto functions required to implement the source providers
type Decryption struct {
	// DecryptBasic implements basic decryption given the private key
	DecryptBasic func(privKey []byte) func(string) E.Either[error, []byte]
}

// openSSLDecryption returns the decryption environment using OpenSSL
func openSSLDecryption() Decryption {
	return Decryption{
		DecryptBasic: OpenSSLDecryptBasic,
	}
}

// cryptoDecryption returns the decryption environment using golang crypto
func cryptoDecryption() Decryption {
	return Decryption{
		DecryptBasic: CryptoDecryptBasic,
	}
}

// // DefaultDecryption detects the decryption environment
func DefaultDecryption() Decryption {
	return F.Pipe1(
		validOpenSSL(),
		E.Fold(F.Ignore1of1[error](cryptoDecryption), F.Ignore1of1[string](openSSLDecryption)),
	)
}
