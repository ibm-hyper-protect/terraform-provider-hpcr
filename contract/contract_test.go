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
// limitations under the License.package datasource

package contract

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	R "github.com/terraform-provider-hpcr/fp/record"
)

var (
	// keypair for testing
	privKey = encrypt.PrivateKey()
	pubKey  = F.Pipe1(
		privKey,
		E.Chain(encrypt.PublicKey),
	)

	// the encryption function based on the keys
	openSSLEncryptBasicE = F.Pipe1(
		pubKey,
		E.Map[error](func(pubKey []byte) func([]byte) E.Either[error, string] {
			return encrypt.EncryptBasic(encrypt.RandomPassword(32), encrypt.AsymmetricEncryptPub(pubKey), encrypt.SymmetricEncrypt)
		}),
	)
)

func TestAddSigningKey(t *testing.T) {
	privKeyE := encrypt.PrivateKey()
	// add to key
	addKey := F.Pipe1(
		privKeyE,
		E.Map[error](addSigningKey),
	)
	// the target map
	var env RawMap

	augE := F.Pipe3(
		addKey,
		E.Chain(I.Ap[RawMap, E.Either[error, RawMap]](env)),
		E.ChainOptionK[error, RawMap, any](func() error {
			return fmt.Errorf("No key [%s]", KeySigningKey)
		})(getSigningKey),
		E.Chain(common.ToTypeE[string]),
	)

	pubE := F.Pipe2(
		privKeyE,
		E.Chain(encrypt.PublicKey),
		E.Map[error](B.ToString),
	)

	assert.Equal(t, pubE, augE)
}

// regular expression used to split the token
// var tokenRe = regexp.MustCompile(`^hyper-protect-basic\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)

func TestUpsetEncrypted(t *testing.T) {
	// the encryption function
	upsertE := F.Pipe1(
		openSSLEncryptBasicE,
		E.Map[error](upsertEncrypted),
	)
	// encrypt env
	encEnv := F.Pipe1(
		upsertE,
		E.Map[error](I.Ap[string, func(RawMap) E.Either[error, RawMap]](KeyEnv)),
	)
	// prepare some data
	data := RawMap{
		KeyEnv: RawMap{
			"type": "env",
		},
	}
	// encrypt the data
	resE := F.Pipe1(
		encEnv,
		E.Chain(I.Ap[RawMap, E.Either[error, RawMap]](data)),
	)
	// validate that the key exists and that it is a token
	getKeyE := F.Flow2(
		R.Lookup[string, any](KeyEnv),
		E.FromOption[error, any](func() error {
			return fmt.Errorf("Key not found")
		}),
	)

	r := F.Pipe2(
		resE,
		E.Chain(getKeyE),
		E.Chain(common.ToTypeE[string]),
	)

	fmt.Println(r)
}
