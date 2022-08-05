// Copyright 2022 IBM Corp.
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
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

type Encryption struct {
	// EncryptBasic implements basic encryption using openSSL given the certificate
	EncryptBasic func(cert []byte) func([]byte) E.Either[error, string]
	// CertFingerprint computes the fingerprint of a certificate
	CertFingerprint (func([]byte) E.Either[error, []byte])
}

// openSSLEncryption returns the encryption environment using OpenSSL
func openSSLEncryption() Encryption {
	return Encryption{
		EncryptBasic:    OpenSSLEncryptBasic,
		CertFingerprint: OpenSSLCertFingerprint,
	}
}

// cryptoEncryption returns the encryption environment using golang crypto
func cryptoEncryption() Encryption {
	return Encryption{
		EncryptBasic:    CryptoEncryptBasic,
		CertFingerprint: CryptoCertFingerprint,
	}
}

// DefaultEncryption detects the encryption environment
func DefaultEncryption() Encryption {
	return F.Pipe1(
		validOpenSSL(),
		E.Fold(F.Ignore1[error](cryptoEncryption), F.Ignore1[string](openSSLEncryption)),
	)
}
