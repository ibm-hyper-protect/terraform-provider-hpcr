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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-provider-hpcr/common"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
)

func TestOpenSSLBinaryFromEnv(t *testing.T) {
	somepath := "/somepath/openssl.exe"
	t.Setenv(KeyEnvOpenSSL, somepath)

	assert.Equal(t, somepath, openSSLBinary())
}

func TestOpenSSLBinary(t *testing.T) {
	assert.NotEmpty(t, openSSLBinary())
}

func TestVersion(t *testing.T) {

	res := openSSLVersion()

	assert.NotEmpty(t, E.IsRight(res))
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

func TestSignDigest(t *testing.T) {
	// some key
	privKeyE := PrivateKey()
	// some input data
	data := []byte("Carsten")

	signE := F.Pipe1(
		privKeyE,
		E.Map[error](SignDigest),
	)

	resE := F.Pipe2(
		signE,
		E.Chain(I.Ap[[]byte, E.Either[error, []byte]](data)),
		E.Map[error](common.Base64Encode),
	)

	assert.True(t, E.IsRight(resE))
}
