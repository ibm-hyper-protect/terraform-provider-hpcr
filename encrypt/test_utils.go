// Copyright 2022 IBM Corp.
// Licensed under the Apache License, Version 2.0 (the "License");
//
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package encrypt

import (
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/stretchr/testify/assert"
)

func SignatureTest(
	privateKey func() E.Either[error, []byte],
	pubKey func([]byte) E.Either[error, []byte],
	randomData func() E.Either[error, []byte],
	signer func([]byte) func([]byte) E.Either[error, []byte],
	validator func([]byte) func([]byte) func([]byte) O.Option[error],
) func(t *testing.T) {
	return func(t *testing.T) {
		// generate a random key
		privKeyE := privateKey()
		// generate some random data
		dataE := randomData()
		// construct the signature
		signE := F.Pipe1(
			privKeyE,
			E.Map[error](signer),
		)
		// signature
		resE := F.Pipe2(
			signE,
			E.Ap[E.Either[error, []byte]](dataE),
			E.Flatten[error, []byte],
		)
		// validate the signature
		validO := F.Pipe5(
			privKeyE,
			E.Chain(pubKey),
			E.Map[error](validator),
			E.Ap[func([]byte) O.Option[error]](dataE),
			E.Ap[O.Option[error]](resE),
			E.GetOrElse(O.Of[error]),
		)
		// handle the option
		assert.Equal(t, O.None[error](), validO)
	}
}
