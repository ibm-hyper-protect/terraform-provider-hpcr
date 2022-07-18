//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package common

import (
	"encoding/base64"
	"os"

	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	O "github.com/terraform-provider-hpcr/fp/option"
)

var (
	Base64Encode  = base64.StdEncoding.EncodeToString
	Base64DecodeE = E.Eitherize1(base64.StdEncoding.DecodeString)

	Base64EncodeE = F.Flow2(
		Base64Encode,
		E.Of[error, string],
	)

	PlainTextEncodeE = F.Flow2(
		B.ToString,
		E.Of[error, string],
	)

	CreateTempE = E.Eitherize2(os.CreateTemp)
)

func ToTypeO[A any](data any) O.Option[A] {
	value, ok := data.(A)
	if ok {
		return O.Some(value)
	}
	return O.None[A]()
}
