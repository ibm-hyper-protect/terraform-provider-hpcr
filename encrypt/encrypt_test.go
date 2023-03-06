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
	"fmt"
	"testing"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/data"
	D "github.com/ibm-hyper-protect/terraform-provider-hpcr/data"
	RA "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/array"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	S "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/string"
	"github.com/stretchr/testify/assert"
)

type Encrypter = func([]byte) E.Either[error, string]
type Decrypter = func(string) E.Either[error, []byte]

var (
	// keypair for testing
	privKey = OpenSSLPrivateKey()
	pubKey  = F.Pipe1(
		privKey,
		E.Chain(OpenSSLPublicKey),
	)

	// the encryption function based on the keys
	openSSLEncryptBasic = createEncryptBasic(func(pubKey []byte) Encrypter {
		return EncryptBasic(OpenSSLRandomPassword(keylen), AsymmetricEncryptPub(pubKey), SymmetricEncrypt)
	})

	// the encryption function based on the keys
	cryptoEncryptBasic = createEncryptBasic(func(pubKey []byte) Encrypter {
		return EncryptBasic(CryptoRandomPassword(keylen), CryptoAsymmetricEncryptPub(pubKey), CryptoSymmetricEncrypt)
	})

	// the decryption function based on the keys
	openSSLDecryptBasic = createDecryptBasic(OpenSSLDecryptBasic)
)

func createEncryptBasic(f func([]byte) Encrypter) E.Either[error, Encrypter] {
	return F.Pipe1(
		pubKey,
		E.Map[error](f),
	)
}

func createDecryptBasic(f func([]byte) Decrypter) E.Either[error, Decrypter] {
	return F.Pipe1(
		privKey,
		E.Map[error](f),
	)
}

func encryptBasic(encE E.Either[error, Encrypter], decE E.Either[error, Decrypter]) func(t *testing.T) {
	// some random test data
	randomData := OpenSSLRandomPassword(1023)

	textE := randomData()
	// encrypt the text
	encTextE := F.Pipe2(
		encE,
		E.Ap[error, []byte, E.Either[error, string]](textE),
		E.Flatten[error, string],
	)
	// decrypt
	decTextE := F.Pipe2(
		decE,
		E.Ap[error, string, E.Either[error, []byte]](encTextE),
		E.Flatten[error, []byte],
	)

	return func(t *testing.T) {
		// compare
		resE := F.Pipe2(
			[]E.Either[error, []byte]{textE, decTextE},
			E.SequenceArray[error, []byte](),
			E.Map[error](func(data [][]byte) bool {
				return assert.Equal(t, data[0], data[1])
			}),
		)

		assert.Equal(t, E.Of[error](true), resE)
	}

}

func TestDefaultEncryption(t *testing.T) {
	// detect the default encryption environment
	env := DefaultEncryption()
	assert.NotNil(t, env.EncryptBasic)
}

func TestDefaultEncryptionFallback(t *testing.T) {
	somepath := "/somepath/openssl.exe"
	t.Setenv(KeyEnvOpenSSL, somepath)
	// detect the default encryption environment
	env := DefaultEncryption()
	assert.NotNil(t, env.EncryptBasic)
}

func TestOpenSSLEncryptBasic(t *testing.T) {
	enc := encryptBasic(openSSLEncryptBasic, openSSLDecryptBasic)
	enc(t)
}

func TestCryptoEncryptBasic(t *testing.T) {
	enc := encryptBasic(cryptoEncryptBasic, openSSLDecryptBasic)
	enc(t)
}

func TestSplitToken(t *testing.T) {
	goodTokens := []string{
		`hyper-protect-basic.UMs93kGaZrzYa6oeoYk8CyaCnsTtRPVdyT+zWBRKKaQD9H71G8bN3PQzbWVx/N84OeyorvERI9RVnpuWwlvnhXj5mu7KZdMXrPoLzW13/zB9HaKYLh64yV3fBsZbGkhlyyjW5n/dcoJ7zbAF5ZRe4m2unpsDUne2cLs27s1FD08oj7iWw/BrzNqqcyOayQnH1WUtHN2OhR4T3k+qSdj3XtnD6t+dsrxg9XFue0zciNQqxDfayBPiUWGpmtOKF2sc+Dp4cq9bV8SsF1crs3dXBsWc21Zl7nVcwt3bmQET++rBdgwI9TZDMa7gjB9Iu/JbjgbPHuBdIycWJMfIH4mseAH6r+HFg5Wq2t/s3FrWg5qdkwCWjzT3r5OoMOafiG06U0SFp29mND1t0kVypf3nEQJQjb6+WoIGcDvKzvUMz5NcRFi8zubziXg0wAJoSZWFL+/gXiDyg9ZbfR8/Ukx52CVLTYGW/IATChfIw51c57b2EddKT3aS/ZksZpyLfLdiLRxLn6X/lEmVGCUojAhmgiFQZzEjeREAV9HMNRnymiyq+qtK+zSMsfZMMdhesHalaRqK9ORqUgBaYII+AG7sWC1xS0FD5LNtN739SjY18/NAY0OznQWI8Yvfu0BoMRSVNIrZl4QWYHdmNHywSfkktc/Bk6qlkgTy392RbfgbcPw=.U2FsdGVkX1/DbyZBRupGSoukxfU91ywFu5HTUsqs8+LLU+MkGP3PJY1XxwaioHoq`,
	}
	goodE := F.Pipe2(
		goodTokens,
		RA.Map(splitToken),
		E.SequenceArray[error, SplitToken](),
	)

	fmt.Println(goodE)
}

func TestCertificate(t *testing.T) {

	assert.NoError(t, F.Pipe3(
		D.DefaultCertificate,
		S.ToBytes,
		CertSerial,
		E.ToError[[]byte],
	))
}

func TestCertFingerprint(t *testing.T) {
	// fingerprint from openSSL
	fpOpenSSL := F.Pipe2(
		data.DefaultCertificate,
		S.ToBytes,
		OpenSSLCertFingerprint,
	)
	// fingerprint directly from crypto
	fpCrypto := F.Pipe2(
		data.DefaultCertificate,
		S.ToBytes,
		CryptoCertFingerprint,
	)
	// make sure they match
	assert.Equal(t, fpOpenSSL, fpCrypto)
}

func TestPrivKeyFingerprints(t *testing.T) {
	// fingerprint from openSSL
	fpOpenSSL := F.Pipe1(
		privKey,
		E.Chain(OpenSSLPrivKeyFingerprint),
	)
	// fingerprint directly from crypto
	fpCrypto := F.Pipe1(
		privKey,
		E.Chain(CryptoPrivKeyFingerprint),
	)
	// make sure they match
	assert.Equal(t, fpOpenSSL, fpCrypto)
}
