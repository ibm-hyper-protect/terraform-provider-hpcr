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
	"regexp"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	R "github.com/terraform-provider-hpcr/fp/record"
	S "github.com/terraform-provider-hpcr/fp/string"
	Y "github.com/terraform-provider-hpcr/fp/yaml"
)

//go:embed samples/contract1.yaml
var Contract1 string

var (
	// keypair for testing
	privKey = encrypt.OpenSSLPrivateKey()
	pubKey  = F.Pipe1(
		privKey,
		E.Chain(encrypt.OpenSSLPublicKey),
	)

	// the encryption function based on the keys
	openSSLEncryptBasicE = F.Pipe1(
		pubKey,
		E.Map[error](func(pubKey []byte) func([]byte) E.Either[error, string] {
			return encrypt.EncryptBasic(encrypt.OpenSSLRandomPassword(32), encrypt.AsymmetricEncryptPub(pubKey), encrypt.SymmetricEncrypt)
		}),
	)
)

func openSSLEncryptAndSignContract(pubKey []byte) func([]byte) func(RawMap) E.Either[error, RawMap] {
	return EncryptAndSignContract(encrypt.EncryptBasic(encrypt.OpenSSLRandomPassword(32), encrypt.AsymmetricEncryptPub(pubKey), encrypt.SymmetricEncrypt), encrypt.OpenSSLSignDigest, encrypt.OpenSSLPublicKey)
}

func TestAddSigningKey(t *testing.T) {
	privKeyE := encrypt.OpenSSLPrivateKey()
	// add to key
	addKey := F.Pipe1(
		privKeyE,
		E.Map[error](addSigningKey(encrypt.OpenSSLPublicKey)),
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
		E.Chain(encrypt.OpenSSLPublicKey),
		common.MapBytesToStgE,
	)

	assert.Equal(t, pubE, augE)
}

// regular expression used to split the token
var tokenRe = regexp.MustCompile(`^hyper-protect-basic\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)

// regular expression used to check for the existence of a public key
var keyRe = regexp.MustCompile(`-----BEGIN PUBLIC KEY-----`)

func TestUpsertEncrypted(t *testing.T) {
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

	r := F.Pipe3(
		resE,
		E.Chain(getKeyE),
		E.Chain(common.ToTypeE[string]),
		E.Chain(E.FromPredicate(tokenRe.MatchString, func(s string) error {
			return fmt.Errorf("string [%s] is not a valid typer protect token", s)
		})),
	)

	assert.True(t, E.IsRight(r))
}

func TestUpsertSigningKey(t *testing.T) {
	privKeyE := encrypt.OpenSSLPrivateKey()
	// add to key
	upsertKeyE := F.Pipe1(
		privKeyE,
		E.Map[error](upsertSigningKey(encrypt.OpenSSLPublicKey)),
	)
	// prepare some contract without a key
	contractE := F.Pipe3(
		Contract1,
		S.ToBytes,
		Y.Parse[RawMap],
		E.Map[error](F.Deref[RawMap]),
	)
	// actually upsert
	resE := F.Pipe5(
		upsertKeyE,
		E.Ap[error, RawMap, E.Either[error, RawMap]](contractE),
		E.Flatten[error, RawMap],
		E.Map[error](F.Ref[RawMap]),
		E.Chain(Y.Stringify[RawMap]),
		common.MapBytesToStgE,
	)
	// check that the serialized form contains the key
	checkE := F.Pipe1(
		resE,
		E.Map[error](keyRe.MatchString),
	)

	assert.Equal(t, E.Of[error](true), checkE)
}

func TestEncryptAndSignContract(t *testing.T) {
	// the private key
	privKeyE := encrypt.OpenSSLPrivateKey()
	// the encryption function
	signerE := F.Pipe2(
		pubKey,
		E.Map[error](openSSLEncryptAndSignContract),
		E.Ap[error, []byte, func(RawMap) E.Either[error, RawMap]](privKeyE),
	)
	// prepare some contract without a key
	contractE := F.Pipe3(
		Contract1,
		S.ToBytes,
		Y.Parse[RawMap],
		E.Map[error](F.Deref[RawMap]),
	)
	// add signature and encrypt the fields
	resE := F.Pipe5(
		signerE,
		E.Ap[error, RawMap, E.Either[error, RawMap]](contractE),
		E.Flatten[error, RawMap],
		E.Map[error](F.Ref[RawMap]),
		E.Chain(Y.Stringify[RawMap]),
		common.MapBytesToStgE,
	)
	assert.True(t, E.IsRight(resE))

	fmt.Println(resE)
}

func TestEnvWorkloadSignature(t *testing.T) {
	// the private key
	privKeyE := encrypt.OpenSSLPrivateKey()

	signer := F.Pipe1(
		privKeyE,
		E.Map[error](createEnvWorkloadSignature(encrypt.OpenSSLSignDigest)),
	)

	// some sample data
	data := RawMap{
		KeyEnv:      "some env",
		KeyWorkload: "some workload",
	}

	// compute the signature
	signatureE := F.Pipe1(
		signer,
		E.Chain(I.Ap[RawMap, E.Either[error, string]](data)),
	)

	assert.True(t, E.IsRight(signatureE))
}
