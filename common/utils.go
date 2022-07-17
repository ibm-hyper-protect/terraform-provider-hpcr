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
)

func Base64Encode(buffer []byte) string {
	return base64.StdEncoding.EncodeToString(buffer)
}

var (
	Base64DecodeE = E.Eitherize1(base64.StdEncoding.DecodeString)

	CreateTempE = E.Eitherize2(os.CreateTemp)
)
