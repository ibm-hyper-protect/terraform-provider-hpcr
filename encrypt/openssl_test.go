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
	"testing"

	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func TestVersion(t *testing.T) {

	res := openSSLVersion

	fmt.Println(res)
}

func TestRandomPassword(t *testing.T) {

	genPwd := RandomPassword(32)

	pwd := genPwd()

	fmt.Println(pwd)
}

func TestEncryptPassword(t *testing.T) {

	//	genPwd := RandomPassword(32)

}

func TestPrivateKey(t *testing.T) {
	privKey := PrivateKey()

	pubKey := F.Pipe2(
		privKey,
		E.Chain(PublicKey),
		E.Map[error](B.ToString),
	)

	fmt.Println(pubKey)
}
