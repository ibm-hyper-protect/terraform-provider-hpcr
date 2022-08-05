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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"

	"github.com/terraform-provider-hpcr/common"
	RA "github.com/terraform-provider-hpcr/fp/array"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	"golang.org/x/crypto/pbkdf2"
)

var (
	parseCertificateE   = E.Eitherize1(x509.ParseCertificate)
	parsePKIXPublicKeyE = E.Eitherize1(x509.ParsePKIXPublicKey)
	toRsaPublicKey      = common.ToTypeE[*rsa.PublicKey]
	randomSaltE         = cryptoRandomE(saltlen)
	aesCipherE          = E.Eitherize1(aes.NewCipher)
	salted              = []byte("Salted__")

	// certToRsaKey decodes a certificate into a public key
	certToRsaKey = F.Flow3(
		pemDecodeE,
		E.Chain(parseCertificateE),
		E.Chain(rsaFromCertificate),
	)

	// pubToRsaKey decodes a public key to rsa format
	pubToRsaKey = F.Flow3(
		pemDecodeE,
		E.Chain(parsePKIXPublicKeyE),
		E.Chain(toRsaPublicKey),
	)
)

// cryptoRandomE returns a random sequence of bytes with the given length
func cryptoRandomE(n int) func() E.Either[error, []byte] {
	return func() E.Either[error, []byte] {
		return E.TryCatchError(func() ([]byte, error) {
			buf := make([]byte, n)
			_, err := rand.Read(buf)
			return buf, err
		})
	}
}

// CryptoRandomPassword creates a random password of given length using characters from the base64 alphabet only
func CryptoRandomPassword(count int) func() E.Either[error, []byte] {
	slice := B.Slice(0, count)
	rnd := cryptoRandomE(count)
	return func() E.Either[error, []byte] {
		return F.Pipe3(
			rnd(),
			E.Map[error](common.Base64Encode),
			common.MapStgToBytesE,
			E.Map[error](slice),
		)
	}
}

// pemDecode will find the next PEM formatted block (certificate, private key etc) in the input
func pemDecodeE(data []byte) E.Either[error, []byte] {
	block, _ := pem.Decode(data)
	return F.Pipe1(
		E.FromNillable[error, pem.Block](fmt.Errorf("enable to decode block from PEM"))(block),
		E.Map[error](func(b *pem.Block) []byte {
			return b.Bytes
		}),
	)
}

// encryptPKCS1v15 creates a function that encrypts a piece of text using a public key
func encryptPKCS1v15(pub *rsa.PublicKey) func([]byte) E.Either[error, []byte] {
	return func(origData []byte) E.Either[error, []byte] {
		return E.TryCatchError(func() ([]byte, error) {
			return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
		})
	}
}

// cryptoAsymmetricEncrypt creates a function that encrypts a piece of text using a public key
func cryptoAsymmetricEncrypt(decKey func([]byte) E.Either[error, *rsa.PublicKey]) func(publicKey []byte) func([]byte) E.Either[error, string] {
	// prepare the encryption callback
	enc := F.Flow2(
		decKey,
		E.Map[error](encryptPKCS1v15),
	)
	return func(publicKey []byte) func([]byte) E.Either[error, string] {
		// decode the input to an RSA public key
		encE := F.Pipe1(
			publicKey,
			enc,
		)
		// returns the encryption function
		return func(data []byte) E.Either[error, string] {
			return F.Pipe2(
				encE,
				E.Chain(I.Ap[[]byte, E.Either[error, []byte]](data)),
				E.Map[error](common.Base64Encode),
			)
		}
	}
}

// // CryptoAsymmetricEncryptPub encrypts a piece of text using a public key
var CryptoAsymmetricEncryptPub = cryptoAsymmetricEncrypt(pubToRsaKey)

// CryptoAsymmetricEncryptCert encrypts a piece of text using a certificate
var CryptoAsymmetricEncryptCert = cryptoAsymmetricEncrypt(certToRsaKey)

// cbcEncrypt creates a new encrypter and then encrypts a plaintext into a cyphertext
func cbcEncrypt(b cipher.Block, iv []byte) func([]byte) []byte {
	return func(src []byte) []byte {
		ciphertext := make([]byte, len(src))
		cipher.NewCBCEncrypter(b, iv).CryptBlocks(ciphertext, src)
		return ciphertext
	}
}

// CryptoSymmetricEncrypt encrypts a set of bytes using a password
func CryptoSymmetricEncrypt(srcPlainbBytes []byte) func([]byte) E.Either[error, string] {
	// Pad plaintext to a multiple of BlockSize with random padding.
	bytesToPad := aes.BlockSize - (len(srcPlainbBytes) % aes.BlockSize)
	// pad the byte array
	paddedPlainBytes := B.Monoid.Concat(srcPlainbBytes, RA.Replicate(bytesToPad, byte(bytesToPad)))
	// length of plain text
	lenPlainBytes := len(paddedPlainBytes)
	// prepare the length buffer
	origSizeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(origSizeBuffer, uint64(lenPlainBytes))

	return func(password []byte) E.Either[error, string] {
		// the salt
		saltE := randomSaltE()
		// derive a key
		keyE := F.Pipe1(
			saltE,
			E.Map[error](func(salt []byte) []byte {
				return pbkdf2.Key(password, salt, iterations, keylen+aes.BlockSize, sha256.New)
			}),
		)
		// the initialization vector
		ivE := F.Pipe1(
			keyE,
			E.Map[error](B.Slice(keylen, keylen+aes.BlockSize)),
		)
		// the block
		blockE := F.Pipe2(
			keyE,
			E.Map[error](B.Slice(0, keylen)),
			E.Chain(aesCipherE),
		)
		// derive the encrypter
		ciphertextE := E.Sequence2(func(b cipher.Block, iv []byte) E.Either[error, func([]byte) []byte] {
			return F.Pipe2(
				cbcEncrypt(b, iv),
				I.Ap[[]byte, []byte](paddedPlainBytes),
				E.Of[error, []byte],
			)
		})(blockE, ivE)
		// derive the final bytes
		return E.Sequence2(func(salt, ciphertext []byte) E.Either[error, string] {
			return F.Pipe1(
				B.ConcatAll(salted, salt, ciphertext),
				common.Base64EncodeE,
			)
		})(saltE, ciphertextE)
	}
}

func rsaFromCertificate(cert *x509.Certificate) E.Either[error, *rsa.PublicKey] {
	return toRsaPublicKey(cert.PublicKey)
}

func rawFromCertificate(cert *x509.Certificate) []byte {
	return cert.Raw
}

// CryptoEncryptBasic implements basic encryption using golang crypto libraries given the certificate
func CryptoEncryptBasic(cert []byte) func([]byte) E.Either[error, string] {
	return EncryptBasic(CryptoRandomPassword(keylen), CryptoAsymmetricEncryptCert(cert), CryptoSymmetricEncrypt)
}

func shaToBytes(sha [32]byte) []byte {
	return sha[:]
}

// CryptoCertFingerprint computes the fingerprint of a certificate using the crypto library
var CryptoCertFingerprint = F.Flow5(
	pemDecodeE,
	E.Chain(parseCertificateE),
	E.Map[error](rawFromCertificate),
	E.Map[error](sha256.Sum256),
	E.Map[error](shaToBytes),
)
