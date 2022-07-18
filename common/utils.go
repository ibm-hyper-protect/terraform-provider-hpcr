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

	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

var (
	Base64Encode  = base64.StdEncoding.EncodeToString
	Base64DecodeE = E.Eitherize1(base64.StdEncoding.DecodeString)

	Base64EncodeE = F.Flow2(
		Base64Encode,
		E.Of[error, string],
	)

	CreateTempE = E.Eitherize2(os.CreateTemp)
)
