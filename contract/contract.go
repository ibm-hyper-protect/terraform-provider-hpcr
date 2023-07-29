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

	B "github.com/IBM/fp-go/bytes"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	O "github.com/IBM/fp-go/option"
	R "github.com/IBM/fp-go/record"
	S "github.com/IBM/fp-go/string"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	Y "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/yaml"
)

type RawMap = map[string]any

var (
	KeyEnv                  = "env"
	KeyWorkload             = "workload"
	KeyAttestationPublicKey = "attestationPublicKey"
	KeySigningKey           = "signingKey"
	KeyEnvWorkloadSignature = "envWorkloadSignature"

	getEnv        = R.Lookup[string, any](KeyEnv)
	getWorkload   = R.Lookup[string, any](KeyWorkload)
	getSigningKey = R.Lookup[string, any](KeySigningKey)

	ParseRawMapE     = Y.Parse[RawMap]
	StringifyRawMapE = Y.Stringify[RawMap]
	MapDerefRawMapE  = E.Map[error](F.Deref[RawMap])
	MapRefRawMapE    = E.Map[error](F.Ref[RawMap])

	// converts an arbitrary value to YAML
	anyToYAML = Y.Stringify[any]

	// converts a string value to bytes
	anyToString = F.Flow2(
		common.ToTypeE[string],
		E.Map[error](S.ToBytes),
	)

	// function that accepts a map, transforms the given key and returns a map with the key encrypted
	upsertYAMLEncrypted = upsertEncrypted(anyToYAML)

	// function that accepts a map, transforms the given key and returns a map with the key encrypted
	upsertStringEncrypted = upsertEncrypted(anyToString)
)

func toAny[A any](a A) any {
	return a
}

func anyToBytes(a any) []byte {
	return []byte(fmt.Sprintf("%s", a))
}

// computes the signature across workload and env
func createEnvWorkloadSignature(signer func([]byte) func([]byte) E.Either[error, []byte]) func([]byte) func(RawMap) E.Either[error, string] {
	return func(privKey []byte) func(RawMap) E.Either[error, string] {
		// callback to construct the digest
		sign := signer(privKey)

		// lookup workload and env
		getEnvO := F.Flow2(
			getEnv,
			O.Map(anyToBytes),
		)
		getWorkloadO := F.Flow2(
			getWorkload,
			O.Map(anyToBytes),
		)
		seqE := O.Sequence2(func(left, right []byte) O.Option[[]byte] {
			return O.Some(B.Monoid.Concat(left, right))
		})
		// combine into a digest
		return func(contract RawMap) E.Either[error, string] {
			// lookup the
			return F.Pipe3(
				seqE(getWorkloadO(contract), getEnvO(contract)),
				E.FromOption[error, []byte](func() error {
					return fmt.Errorf("the contract is missing [%s] or [%s] or both", KeyEnv, KeyWorkload)
				}),
				E.Chain(sign),
				E.Map[error](common.Base64Encode),
			)
		}
	}
}

// constructs a workload across workload and env and adds this to the map
func upsertEnvWorkloadSignature(signer func([]byte) func([]byte) E.Either[error, []byte]) func([]byte) func(RawMap) E.Either[error, RawMap] {
	// callback that can creates the signature
	envWorkloadSignature := createEnvWorkloadSignature(signer)

	return func(privKey []byte) func(RawMap) E.Either[error, RawMap] {
		// callback to create the signature
		create := envWorkloadSignature(privKey)
		setSignature := F.Bind1st(R.UpsertAt[string, any], KeyEnvWorkloadSignature)

		return func(contract RawMap) E.Either[error, RawMap] {
			return F.Pipe4(
				contract,
				create,
				E.Map[error](toAny[string]),
				E.Map[error](setSignature),
				E.Map[error](I.Ap[RawMap, RawMap](contract)),
			)
		}
	}
}

// returns a function that adds the public part of the key to the input mapping
func addSigningKey(pubKey func([]byte) E.Either[error, []byte]) func(key []byte) func(RawMap) E.Either[error, RawMap] {
	// callback to decode the key
	getPemE := F.Flow4(
		pubKey,
		common.MapBytesToStgE,
		E.Map[error](toAny[string]),
		E.Map[error](F.Bind1st(R.UpsertAt[string, any], KeySigningKey)),
	)

	return func(key []byte) func(RawMap) E.Either[error, RawMap] {
		// function to add the pkey into a map
		pemE := F.Pipe1(
			key,
			getPemE,
		)
		// actually work on a map
		return func(data RawMap) E.Either[error, RawMap] {
			// insert into the map
			return F.Pipe1(
				pemE,
				E.Map[error](I.Ap[RawMap, RawMap](data)),
			)
		}
	}
}

// upsertSigningKey returns a function that adds the public part of the signing key
func upsertSigningKey(pubKey func([]byte) E.Either[error, []byte]) func([]byte) func(RawMap) E.Either[error, RawMap] {
	// bind the function to the callback that can extract a public key
	addSigningKeyE := addSigningKey(pubKey)
	setEnv := F.Bind1st(R.UpsertAt[string, any], KeyEnv)

	return func(privKey []byte) func(RawMap) E.Either[error, RawMap] {
		// adds the signing key to the env map
		addKeyE := addSigningKeyE(privKey)

		return func(contract RawMap) E.Either[error, RawMap] {
			// get the env part, fall back to the empty map, then insert the signature
			return F.Pipe7(
				contract,
				getEnv,
				O.Chain(common.ToTypeO[RawMap]),
				O.GetOrElse(F.Constant(make(RawMap))),
				addKeyE,
				E.Map[error](toAny[RawMap]),
				E.Map[error](setEnv),
				E.Map[error](I.Ap[RawMap, RawMap](contract)),
			)
		}
	}
}

// function that accepts a map, transforms the given key and returns a map with the key encrypted
func upsertEncrypted(serializer func(any) E.Either[error, []byte]) func(enc func(data []byte) E.Either[error, string]) func(string) func(RawMap) E.Either[error, RawMap] {
	return func(enc func(data []byte) E.Either[error, string]) func(string) func(RawMap) E.Either[error, RawMap] {
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
					O.Map(F.Flow5(
						serializer,
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
}

// EncryptAndSignContract returns a function that signs the workload and env part of a contract and that adds the public key of the signature to the map
// - enc encrypts a piece of data
// - signer signs a piece of data
// - pubKey extracts the public key from the private key
func EncryptAndSignContract(
	enc func(data []byte) E.Either[error, string],
	signer func([]byte) func([]byte) E.Either[error, []byte],
	pubKey func([]byte) E.Either[error, []byte],
) func([]byte) func(RawMap) E.Either[error, RawMap] {
	// the upserter
	upsertKey := upsertSigningKey(pubKey)
	upsertSig := upsertEnvWorkloadSignature(signer)
	// the function that encrypts fields
	encrypterYAML := upsertYAMLEncrypted(enc)
	encEnv := encrypterYAML(KeyEnv)
	encWorkload := encrypterYAML(KeyWorkload)
	encrypterString := upsertStringEncrypted(enc)
	encAttPubKey := encrypterString(KeyAttestationPublicKey)
	// callback to handle signature
	return func(privKey []byte) func(RawMap) E.Either[error, RawMap] {
		// the signature callback
		addPubKey := upsertKey(privKey)
		addSignature := upsertSig(privKey)
		// execute one step after the other
		return F.Flow5(
			addPubKey,
			E.Chain(encEnv),
			E.Chain(encWorkload),
			E.Chain(encAttPubKey),
			E.Chain(addSignature),
		)
	}
}
