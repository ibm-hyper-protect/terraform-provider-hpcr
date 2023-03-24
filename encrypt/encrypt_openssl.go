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
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
)

// OpenSSLEncryptBasic implements basic encryption using openSSL given the certificate
func OpenSSLEncryptBasic(cert []byte) func([]byte) E.Either[error, string] {
	return EncryptBasic(OpenSSLRandomPassword(keylen), OpenSSLAsymmetricEncryptCert(cert), OpenSSLSymmetricEncrypt)
}

// OpenSSLDecryptBasic implements basic decryption using openSSL given the private key
func OpenSSLDecryptBasic(privKey []byte) func(string) E.Either[error, []byte] {
	return DecryptBasic(OpenSSLAsymmetricDecrypt(privKey), OpenSSLSymmetricDecrypt)
}
