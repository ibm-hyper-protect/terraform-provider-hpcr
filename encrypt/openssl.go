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
	"io"
	"os"
	"strings"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	RA "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/array"
	B "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/bytes"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	FL "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/file"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	I "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/identity"
	O "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/option"
	P "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/predicate"
	S "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/string"
	T "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/tuple"
)

// OpenSSLVersion represents the openSSL version, including the path to the binary
type OpenSSLVersion = T.Tuple2[string, string]

var (
	// name of the environment variable carrying the openSSL binary
	KeyEnvOpenSSL = "OPENSSL_BIN"

	// default name of the openSSL binary
	defaultOpenSSL = "openssl"

	// the empty byte array
	emptyBytes = RA.Empty[byte]()

	// operator to extract stdout
	mapStdout = E.Map[error](common.GetStdOut)

	getPath    = T.FirstOf2[string, string]
	getVersion = T.SecondOf2[string, string]

	// operator to convert stdout to base64
	base64StdOut = F.Flow2(
		mapStdout,
		E.Map[error](common.Base64Encode),
	)

	// OpenSSLSignDigest signs the sha256 digest using a private key
	OpenSSLSignDigest = handle(signDigest)

	AsymmetricEncryptPub = handle(asymmetricEncryptPub)

	AsymmetricEncryptCert = handle(asymmetricEncryptCert)

	AsymmerticDecrypt = handle(asymmetricDecrypt)

	SymmetricEncrypt = handle(symmetricEncrypt)

	// OpenSSLPublicKey gets the public key from a private key
	OpenSSLPublicKey = F.Flow2(
		OpenSSL("rsa", "-pubout"),
		mapStdout,
	)

	// CertSerial gets the serial number from a certificate
	CertSerial = F.Flow2(
		OpenSSL("x509", "-serial", "-noout"),
		mapStdout,
	)

	// OpenSSLCertFingerprint gets the fingerprint of a certificate
	OpenSSLCertFingerprint = F.Flow4(
		OpenSSL("x509", "--outform", "DER"),
		mapStdout,
		E.Chain(OpenSSL("sha256", "--binary")),
		mapStdout,
	)

	// gets the fingerprint of the private key
	OpenSSLPrivKeyFingerprint = F.Flow4(
		OpenSSL("rsa", "-pubout", "-outform", "DER"),
		mapStdout,
		E.Chain(OpenSSL("sha256", "--binary")),
		mapStdout,
	)

	// tests if a string contains "OpenSSL"
	includesOpenSSL = S.Includes("OpenSSL")
)

// version string of the openSSL binary together with the binary
func openSSLVersion() E.Either[error, OpenSSLVersion] {
	// binary
	bin := openSSLBinary()
	// check the version
	return F.Pipe5(
		emptyBytes,
		common.ExecCommand(bin, "version"),
		mapStdout,
		common.MapBytesToStgE,
		E.Map[error](strings.TrimSpace),
		E.Map[error](F.Bind1st(T.MakeTuple2[string, string], bin)),
	)
}

// name of the open SSL binary either from the environment or a fallback
func openSSLBinary() string {
	return F.Pipe2(
		KeyEnvOpenSSL,
		O.FromValidation(os.LookupEnv),
		O.GetOrElse(F.Constant(defaultOpenSSL)),
	)
}

// command name of the valid openSSL binary
func validOpenSSL() E.Either[error, string] {
	return F.Pipe1(
		openSSLVersion(),
		E.Chain(func(version OpenSSLVersion) E.Either[error, string] {
			return F.Pipe3(
				version,
				O.FromPredicate(P.ContraMap(getVersion)(includesOpenSSL)),
				O.Map(getPath),
				E.FromOption[error, string](func() error {
					return fmt.Errorf("openSSL Version [%s] is unsupported", version)
				}),
			)
		}),
	)
}

// helper to safely write data into a file
func writeData[W io.Writer](data []byte) func(w W) E.Either[error, int] {
	return func(w W) E.Either[error, int] {
		return E.TryCatchError(func() (int, error) {
			return w.Write(data)
		})
	}
}

func OpenSSL(args ...string) func([]byte) E.Either[error, common.CommandOutput] {
	// validate the version of openssl and make sure to use the right one
	cmdE := F.Pipe1(
		validOpenSSL(),
		E.Map[error](func(cmd string) func([]byte) E.Either[error, common.CommandOutput] {
			return common.ExecCommand(cmd, args...)
		}),
	)
	// convert stdin to openssl output
	return func(dataIn []byte) E.Either[error, common.CommandOutput] {
		return F.Pipe1(
			cmdE,
			E.Chain(I.Ap[[]byte, E.Either[error, common.CommandOutput]](dataIn)),
		)
	}
}

// OpenSSLRandomPassword creates a random password of given length using characters from the base64 alphabet only
func OpenSSLRandomPassword(count int) func() E.Either[error, []byte] {
	cmdE := OpenSSL("rand", fmt.Sprintf("%d", count))
	slice := B.Slice(0, count)
	return func() E.Either[error, []byte] {
		return F.Pipe4(
			emptyBytes,
			cmdE,
			base64StdOut,
			common.MapStgToBytesE,
			E.Map[error](slice),
		)
	}
}

// persists the data record for a minimal timespan in a temporary file and the invokes a callback
func handle[A, R any](cb func(string) func(A) E.Either[error, R]) func(data []byte) func(A) E.Either[error, R] {
	tmpFile := FL.WithTempFile[R]()
	// handle temp file
	return func(data []byte) func(A) E.Either[error, R] {
		writeDataE := writeData[*os.File](data)
		return func(key A) E.Either[error, R] {
			mapToA := E.MapTo[error, int](key)
			return tmpFile(func(f *os.File) E.Either[error, R] {
				enc := cb(f.Name())
				return F.Pipe3(
					f,
					writeDataE,
					mapToA,
					E.Chain(enc),
				)
			})
		}
	}
}

func signDigest(keyFile string) func([]byte) E.Either[error, []byte] {
	return F.Flow2(
		OpenSSL("dgst", "-sha256", "-sign", keyFile),
		mapStdout,
	)
}

func asymmetricDecrypt(keyFile string) func(string) E.Either[error, []byte] {
	return F.Flow3(
		common.Base64DecodeE,
		E.Chain(OpenSSL("rsautl", "-decrypt", "-inkey", keyFile)),
		mapStdout,
	)
}

func asymmetricEncryptPub(keyFile string) func([]byte) E.Either[error, string] {
	return F.Flow2(
		OpenSSL("rsautl", "-encrypt", "-pubin", "-inkey", keyFile),
		base64StdOut,
	)
}

func asymmetricEncryptCert(certFile string) func([]byte) E.Either[error, string] {
	return F.Flow2(
		OpenSSL("rsautl", "-encrypt", "-certin", "-inkey", certFile),
		base64StdOut,
	)
}

func symmetricEncrypt(dataFile string) func([]byte) E.Either[error, string] {
	return F.Flow2(
		OpenSSL("enc", "-aes-256-cbc", "-pbkdf2", "-in", dataFile, "-pass", "stdin"),
		base64StdOut,
	)
}

func symmetricDecrypt(dataFile string) func([]byte) E.Either[error, []byte] {
	return F.Flow2(
		OpenSSL("aes-256-cbc", "-d", "-pbkdf2", "-in", dataFile, "-pass", "stdin"),
		mapStdout,
	)
}

func SymmetricDecrypt(token string) func([]byte) E.Either[error, []byte] {
	// decode the token and produce the decryption function
	dec := F.Pipe2(
		token,
		common.Base64DecodeE,
		E.Map[error](handle(symmetricDecrypt)),
	)
	// decrypt using the provided password
	return func(pwd []byte) E.Either[error, []byte] {
		return F.Pipe1(
			dec,
			E.Chain(I.Ap[[]byte, E.Either[error, []byte]](pwd)),
		)
	}
}

// OpenSSLPrivateKey generates a private key
func OpenSSLPrivateKey() E.Either[error, []byte] {
	return F.Pipe2(
		emptyBytes,
		OpenSSL("genrsa", "4096"),
		mapStdout,
	)
}

// OpenSSLVerifyDigest verifies the signature of the input data against a signature
func OpenSSLVerifyDigest(pubKey []byte) func([]byte) func([]byte) O.Option[error] {
	// shortcut for the fold operation
	foldE := E.Fold(O.Of[error], F.Constant1[common.CommandOutput](O.None[error]()))
	// callback functions
	return func(data []byte) func([]byte) O.Option[error] {
		return func(signature []byte) O.Option[error] {
			return F.Pipe2(
				data,
				handle(func(pubKeyFile string) func([]byte) E.Either[error, common.CommandOutput] {
					return handle(func(signatureFile string) func([]byte) E.Either[error, common.CommandOutput] {
						return OpenSSL("dgst", "-verify", pubKeyFile, "-sha256", "-signature", signatureFile)
					})(signature)
				})(pubKey),
				foldE,
			)
		}
	}
}
