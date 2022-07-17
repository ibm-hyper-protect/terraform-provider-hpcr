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

// implements basic encryption using openSSL given the public key
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

// implements basic decryption using openSSL given the private key
func OpenSSLDecryptBasic(privKey []byte) func(string) E.Either[error, []byte] {
	return DecryptBasic(AsymmerticDecrypt(privKey), SymmetricDecrypt)
}
