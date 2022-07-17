//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package encrypt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/terraform-provider-hpcr/common"
	RA "github.com/terraform-provider-hpcr/fp/array"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	FL "github.com/terraform-provider-hpcr/fp/file"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	S "github.com/terraform-provider-hpcr/fp/string"
)

var (
	// the empty byte array
	emptyBytes = RA.Empty[byte]()

	// operator to extract stdout
	mapStdout = E.Map[error](common.GetStdOut)

	// version string of the openSSL binary
	openSSLVersion = F.Pipe4(
		emptyBytes,
		common.ExecCommand("openssl", "version"),
		mapStdout,
		E.Map[error](B.ToString),
		E.Map[error](strings.TrimSpace),
	)

	// command name of the valid openSSL binary
	validOpenSSL = F.Pipe1(
		openSSLVersion,
		E.Chain(func(version string) E.Either[error, string] {
			if strings.Contains(version, "OpenSSL") {
				return E.Of[error]("openssl")
			}
			return E.Left[error, string](fmt.Errorf("openSSL Version [%s] is unsupported", version))

		}),
	)

	// operator to convert stdout to base64
	base64StdOut = F.Flow2(
		mapStdout,
		E.Map[error](common.Base64Encode),
	)
)

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
		validOpenSSL,
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

func RandomPassword(count int) func() E.Either[error, []byte] {
	cmdE := OpenSSL("rand", fmt.Sprintf("%d", count))
	slice := B.Slice(0, count)
	return func() E.Either[error, []byte] {
		return F.Pipe4(
			emptyBytes,
			cmdE,
			base64StdOut,
			E.Map[error](S.ToBytes),
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

var AsymmetricEncryptPub = handle(asymmetricEncryptPub)

func asymmetricEncryptCert(certFile string) func([]byte) E.Either[error, string] {
	return F.Flow2(
		OpenSSL("rsautl", "-encrypt", "-certin", "-inkey", certFile),
		base64StdOut,
	)
}

var AsymmetricEncryptCert = handle(asymmetricEncryptCert)

var AsymmerticDecrypt = handle(asymmetricDecrypt)

func symmetricEncrypt(dataFile string) func([]byte) E.Either[error, string] {
	return F.Flow2(
		OpenSSL("enc", "-aes-256-cbc", "-pbkdf2", "-in", dataFile, "-pass", "stdin"),
		base64StdOut,
	)
}

var SymmetricEncrypt = handle(symmetricEncrypt)

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

// generates a private key
func PrivateKey() E.Either[error, []byte] {
	return F.Pipe2(
		emptyBytes,
		OpenSSL("genrsa", "4096"),
		mapStdout,
	)
}

// gets the public key from a private key
var PublicKey = F.Flow2(
	OpenSSL("rsa", "-pubout"),
	mapStdout,
)

// gets the serial number from a certificate
var CertSerial = F.Flow2(
	OpenSSL("x509", "-serial", "-noout"),
	mapStdout,
)
