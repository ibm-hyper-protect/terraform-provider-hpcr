// Copyright 2022 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package encrypt

import (
	"fmt"
	"regexp"

	"github.com/terraform-provider-hpcr/common"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	T "github.com/terraform-provider-hpcr/fp/tuple"
)

type SplitToken = T.Tuple2[string, string]

func EncryptBasic(
	genPwd func() E.Either[error, []byte],
	asymmEncrypt func([]byte) E.Either[error, string],
	symmEncrypt func([]byte) func([]byte) E.Either[error, string],
) func([]byte) E.Either[error, string] {

	return func(data []byte) E.Either[error, string] {
		// generate the password
		pwdE := genPwd()

		// encode the password
		encPwd := F.Pipe1(
			pwdE,
			E.Chain(asymmEncrypt),
		)

		// encode the data
		encData := F.Pipe1(
			pwdE,
			E.Chain(symmEncrypt(data)),
		)

		// combine to a hyper protect token
		return E.Sequence2(func(pwd string, token string) E.Either[error, string] {
			return E.Of[error](fmt.Sprintf("%s.%s.%s", common.PrefixBasicEncoding, pwd, token))
		})(encPwd, encData)
	}
}

// OpenSSLEncryptBasic implements basic encryption using openSSL given the public key
func OpenSSLEncryptBasic(pubKey []byte) func([]byte) E.Either[error, string] {
	return EncryptBasic(RandomPassword(32), AsymmetricEncryptCert(pubKey), SymmetricEncrypt)
}

// regular expression used to split the token
var tokenRe = regexp.MustCompile(`^hyper-protect-basic\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)

var errNoMatch = E.Left[error, SplitToken](fmt.Errorf("token does not match the specification"))

func splitToken(token string) E.Either[error, SplitToken] {
	all := tokenRe.FindAllStringSubmatch(token, -1)
	if all == nil {
		return errNoMatch
	}
	match := all[0]
	return E.Of[error](T.MakeTuple2(match[1], match[2]))
}

var (
	getPwd   = T.FirstOf2[string, string]
	getToken = T.SecondOf2[string, string]
)

func DecryptBasic(
	asymmDecrypt func(string) E.Either[error, []byte],
	symmDecrypt func(string) func([]byte) E.Either[error, []byte],
) func(string) E.Either[error, []byte] {

	return func(data string) E.Either[error, []byte] {
		// split the string
		splitE := F.Pipe1(
			data,
			splitToken,
		)
		// get password
		pwdE := F.Pipe2(
			splitE,
			E.Map[error](getPwd),
			E.Chain(asymmDecrypt),
		)

		// get the token
		return F.Pipe4(
			splitE,
			E.Map[error](getToken),
			E.Map[error](symmDecrypt),
			E.Ap[error, []byte, E.Either[error, []byte]](pwdE),
			E.Flatten[error, []byte],
		)
	}
}

// OpenSSLDecryptBasic implements basic decryption using openSSL given the private key
func OpenSSLDecryptBasic(privKey []byte) func(string) E.Either[error, []byte] {
	return DecryptBasic(AsymmerticDecrypt(privKey), SymmetricDecrypt)
}
