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
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	RA "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/array"
	B "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/bytes"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	I "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/identity"
	O "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/option"
	"golang.org/x/crypto/pbkdf2"
)

var (
	parseCertificateE     = E.Eitherize1(x509.ParseCertificate)
	parsePKIXPublicKeyE   = E.Eitherize1(x509.ParsePKIXPublicKey)
	parsePKCS1PrivateKeyE = E.Eitherize1(x509.ParsePKCS1PrivateKey)
	marshalPKIXPublicKeyE = E.Eitherize1(x509.MarshalPKIXPublicKey)
	toRsaPublicKey        = common.ToTypeE[*rsa.PublicKey]
	randomSaltE           = cryptoRandomE(saltlen)
	aesCipherE            = E.Eitherize1(aes.NewCipher)
	salted                = []byte("Salted__")

	// certToRsaKey decodes a certificate into a public key
	certToRsaKey = F.Flow3(
		pemDecodeE,
		E.Chain(parseCertificateE),
		E.Chain(rsaFromCertificate),
	)

	// privToRsaKey decodes a pkcs file into a private key
	privToRsaKey = F.Flow2(
		pemDecodeE,
		E.Chain(parsePKCS1PrivateKeyE),
	)

	// pubToRsaKey decodes a public key to rsa format
	pubToRsaKey = F.Flow3(
		pemDecodeE,
		E.Chain(parsePKIXPublicKeyE),
		E.Chain(toRsaPublicKey),
	)

	// CryptoCertFingerprint computes the fingerprint of a certificate using the crypto library
	CryptoCertFingerprint = F.Flow5(
		pemDecodeE,
		E.Chain(parseCertificateE),
		E.Map[error](rawFromCertificate),
		E.Map[error](sha256.Sum256),
		E.Map[error](shaToBytes),
	)

	// CryptoPrivKeyFingerprint computes the fingerprint of a private key using the crypto library
	CryptoPrivKeyFingerprint = F.Flow7(
		pemDecodeE,
		E.Chain(parsePKCS1PrivateKeyE),
		E.Map[error](privToPub),
		E.Map[error](pubToAny),
		E.Chain(marshalPKIXPublicKeyE),
		E.Map[error](sha256.Sum256),
		E.Map[error](shaToBytes),
	)

	// CryptoVerifyDigest verifies the signature of the input data against a signature
	CryptoVerifyDigest = F.Flow2(
		pubToRsaKey,
		E.Fold(errorValidator, verifyPKCS1v15),
	)

	// CryptoPublicKey extracts the public key from a private key
	CryptoPublicKey = F.Flow6(
		pemDecodeE,
		E.Chain(parsePKCS1PrivateKeyE),
		E.Map[error](privToPub),
		E.Map[error](pubToAny),
		E.Chain(marshalPKIXPublicKeyE),
		E.Map[error](func(data []byte) []byte {
			return pem.EncodeToMemory(
				&pem.Block{
					Type:  "PUBLIC KEY",
					Bytes: data,
				},
			)
		}),
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

// decryptPKCS1v15 creates a function that decrypts a piece of text using a private key
func decryptPKCS1v15(pub *rsa.PrivateKey) func([]byte) E.Either[error, []byte] {
	return func(ciphertext []byte) E.Either[error, []byte] {
		return E.TryCatchError(func() ([]byte, error) {
			return rsa.DecryptPKCS1v15(rand.Reader, pub, ciphertext)
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

// CryptoAsymmetricDecrypt decrypts a piece of text using a private key
var CryptoAsymmetricDecrypt = cryptoAsymmetricDecrypt(privToRsaKey)

// cbcEncrypt creates a new encrypter and then encrypts a plaintext into a ciphertext
func cbcEncrypt(b cipher.Block, iv []byte) func([]byte) []byte {
	return func(src []byte) []byte {
		ciphertext := make([]byte, len(src))
		cipher.NewCBCEncrypter(b, iv).CryptBlocks(ciphertext, src)
		return ciphertext
	}
}

// cbcDecrypt creates a new decryptor and then decrypts ciphertext into plaintext
func cbcDecrypt(b cipher.Block, iv []byte) func([]byte) []byte {
	return func(src []byte) []byte {
		plaintext := make([]byte, len(src))
		cipher.NewCBCDecrypter(b, iv).CryptBlocks(plaintext, src)
		return plaintext
	}
}

// CryptoSymmetricEncrypt encrypts a set of bytes using a password
func CryptoSymmetricEncrypt(srcPlainBytes []byte) func([]byte) E.Either[error, string] {
	// Pad plaintext to a multiple of BlockSize with random padding.
	bytesToPad := aes.BlockSize - (len(srcPlainBytes) % aes.BlockSize)
	// pad the byte array
	paddedPlainBytes := B.Monoid.Concat(srcPlainBytes, RA.Replicate(bytesToPad, byte(bytesToPad)))
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
		ciphertextE := E.Sequence2(func(b cipher.Block, iv []byte) E.Either[error, []byte] {
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

// OpenSSLDecryptBasic implements basic decryption using golang crypto libraries given the private key
func CryptoDecryptBasic(privKey []byte) func(string) E.Either[error, []byte] {
	return DecryptBasic(CryptoAsymmetricDecrypt(privKey), CryptoSymmetricDecrypt)
}

func shaToBytes(sha [32]byte) []byte {
	return sha[:]
}

func privToPub(privKey *rsa.PrivateKey) *rsa.PublicKey {
	return &privKey.PublicKey
}

func pubToAny(pubKey *rsa.PublicKey) any {
	return pubKey
}

func privKeyToPem(privKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)
}

// CryptoPrivateKey generates a private key
func CryptoPrivateKey() E.Either[error, []byte] {
	return F.Pipe1(
		E.TryCatchError(func() (*rsa.PrivateKey, error) {
			return rsa.GenerateKey(rand.Reader, 4096)
		}),
		E.Map[error](privKeyToPem),
	)
}

// implements the signing operation in a functional way
func signPKCS1v15(privateKey *rsa.PrivateKey) func([]byte) E.Either[error, []byte] {
	return func(digest []byte) E.Either[error, []byte] {
		return E.TryCatchError(func() ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
		})
	}
}

// CryptoSignDigest generates a signature across the sha256
func CryptoSignDigest(privKey []byte) func([]byte) E.Either[error, []byte] {
	// parse the private key and derive the signer from it
	signE := F.Pipe3(
		privKey,
		pemDecodeE,
		E.Chain(parsePKCS1PrivateKeyE),
		E.Map[error](signPKCS1v15),
	)
	return func(data []byte) E.Either[error, []byte] {
		// compute the digest
		digestE := F.Pipe2(
			data,
			sha256.Sum256,
			shaToBytes,
		)
		// apply the signer
		return F.Pipe1(
			signE,
			E.Chain(I.Ap[[]byte, E.Either[error, []byte]](digestE)),
		)
	}
}

// implements the validation operation in a functional way
func verifyPKCS1v15(pubKey *rsa.PublicKey) func([]byte) func([]byte) O.Option[error] {
	return func(data []byte) func([]byte) O.Option[error] {
		digest := sha256.Sum256(data)
		return func(signature []byte) O.Option[error] {
			return common.FromErrorO(rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, digest[:], signature))
		}
	}
}

// errorValidator returns a validator that returns the orignal error
var errorValidator = F.Flow3(
	O.Of[error],
	F.Constant1[[]byte, O.Option[error]],
	F.Constant1[[]byte, func([]byte) O.Option[error]],
)

func unpad(data []byte) []byte {
	size := len(data)
	count := int(data[size-1])
	return data[0 : size-count]
}

// CryptoSymmetricDecrypt encrypts a set of bytes using a password
func CryptoSymmetricDecrypt(srcText string) func([]byte) E.Either[error, []byte] {
	// some offsets
	offSalt := len(salted)
	offciphertext := offSalt + saltlen
	// decode the source (would start with `salted`)
	srcBytesE := common.Base64DecodeE(srcText)
	// get the salt
	saltE := F.Pipe1(
		srcBytesE,
		E.Map[error](B.Slice(offSalt, offciphertext)),
	)
	// get the ciphertext
	ciphertextE := F.Pipe1(
		srcBytesE,
		E.Map[error](B.SliceRight(offciphertext)),
	)

	return func(password []byte) E.Either[error, []byte] {
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
		// decrypt
		return E.Sequence3(func(b cipher.Block, iv []byte, ciphertext []byte) E.Either[error, []byte] {
			return F.Pipe3(
				cbcDecrypt(b, iv),
				I.Ap[[]byte, []byte](ciphertext),
				unpad,
				E.Of[error, []byte],
			)
		})(blockE, ivE, ciphertextE)
	}
}

// cryptoAsymmetricDecrypt creates a function that encrypts a piece of text using a private key
func cryptoAsymmetricDecrypt(decKey func([]byte) E.Either[error, *rsa.PrivateKey]) func(privKey []byte) func(string) E.Either[error, []byte] {
	// prepare the decryption callback
	dec := F.Flow2(
		decKey,
		E.Map[error](decryptPKCS1v15),
	)
	return func(privKey []byte) func(string) E.Either[error, []byte] {
		// decode the input to an RSA public key
		decE := F.Pipe1(
			privKey,
			dec,
		)
		// returns the encryption function
		return func(data string) E.Either[error, []byte] {
			return F.Pipe2(
				decE,
				E.Ap[error, []byte, E.Either[error, []byte]](common.Base64DecodeE(data)),
				E.Flatten[error, []byte],
			)
		}
	}
}
