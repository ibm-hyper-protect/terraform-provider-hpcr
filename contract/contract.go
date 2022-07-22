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
	"crypto/sha256"
	"fmt"

	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	O "github.com/terraform-provider-hpcr/fp/option"
	R "github.com/terraform-provider-hpcr/fp/record"
	Y "github.com/terraform-provider-hpcr/fp/yaml"
)

var (
	KeyEnv        = "env"
	KeyWorkload   = "workload"
	KeySigningKey = "signingKey"

	getEnv        = R.Lookup[string, any](KeyEnv)
	getWorkload   = R.Lookup[string, any](KeyWorkload)
	getSigningKey = R.Lookup[string, any](KeySigningKey)
)

type RawMap = map[string]any

func toAny[A any](a A) any {
	return a
}

// computes a sha256 as hex
func createHash(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// computes the signature across workload and env
func createEnvWorkloadSignature(privKey []byte) func(RawMap) E.Either[error, RawMap] {

	func(RawMap) E.Either[error, RawMap] {

	}

}

// returns a function that adds the public part of the key to the input mapping
func addSigningKey(key []byte) func(RawMap) E.Either[error, RawMap] {
	// function to add the pkey into a map
	pemE := F.Pipe4(
		key,
		encrypt.PublicKey,
		E.Map[error](B.ToString),
		E.Map[error](toAny[string]),
		E.Map[error](F.Bind1st(R.UpsertAt[string, any], KeySigningKey)),
	)

	return func(data RawMap) E.Either[error, RawMap] {
		// insert into the map
		return F.Pipe1(
			pemE,
			E.Map[error](I.Ap[RawMap, RawMap](data)),
		)
	}
}

func upsertSigningKey(key []byte) func(RawMap) E.Either[error, RawMap] {
	// adds the signing key to the env map
	addKey := addSigningKey(key)
	setEnv := F.Bind1st(R.UpsertAt[string, any], KeyEnv)

	return func(contract RawMap) E.Either[error, RawMap] {
		// get the env part, fall back to the empty map, then insert the signature
		return F.Pipe7(
			contract,
			getEnv,
			O.Chain(common.ToTypeO[RawMap]),
			O.GetOrElse(F.Constant(make(RawMap))),
			addKey,
			E.Map[error](toAny[RawMap]),
			E.Map[error](setEnv),
			E.Map[error](I.Ap[RawMap, RawMap](contract)),
		)
	}
}

// function that accepts a map, transforms the given key and returns a map with the key encrypted
func upsertEncrypted(enc func(data []byte) E.Either[error, string]) func(string) func(RawMap) E.Either[error, RawMap] {
	// callback that accepts the key
	return func(key string) func(RawMap) E.Either[error, RawMap] {
		// callback to insert the key into the target
		setKey := F.Bind1st(R.UpsertAt[string, any], key)
		getKey := R.Lookup[string, any](key)
		// returns the actual upserter
		return func(dst RawMap) E.Either[error, RawMap] {
			// lookup the original key
			return F.Pipe3(
				dst,
				getKey,
				O.Map(F.Flow6(
					F.Ref[any],
					Y.Stringify[any],
					E.Chain(enc),
					E.Map[error](toAny[string]),
					E.Map[error](setKey),
					E.Map[error](I.Ap[RawMap, RawMap](dst)),
				)),
				O.GetOrElse(F.Constant(E.Of[error](dst))),
			)
		}
	}
}

func EncryptAndSignContract(enc func(data []byte) E.Either[error, string]) func([]byte) func(RawMap) E.Either[error, RawMap] {
	// the function that encrypts fields
	encrypter := upsertEncrypted(enc)
	encEnv := encrypter(KeyEnv)
	encWorkload := encrypter(KeyWorkload)
	// callback to handle signature
	return func(privKey []byte) func(RawMap) E.Either[error, RawMap] {
		// the signature callback
		addKey := upsertSigningKey(privKey)
		// execute one step after the other
		return F.Flow3(
			addKey,
			E.Chain(encEnv),
			E.Chain(encWorkload),
		)
	}

}
