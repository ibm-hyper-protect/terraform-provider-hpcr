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

package datasource

import (
	"fmt"
	"testing"

	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/contract"
	D "github.com/terraform-provider-hpcr/data"
	"github.com/terraform-provider-hpcr/encrypt"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
)

func TestHashWithCertAndKey(t *testing.T) {

	privKeyE := encrypt.PrivateKey()

	dataE := F.Pipe2(
		privKeyE,
		E.Map[error](func(key []byte) contract.RawMap {
			data := make(map[string]any)

			// prepare input data
			data[common.KeyCert] = D.DefaultCertificate
			data[common.KeyPrivKey] = string(key)

			return data
		}),
		E.Map[error](CreateResourceDataMock),
	)

	test := []byte("Carsten")

	hashE := F.Pipe2(
		dataE,
		E.Map[error](createHashWithCertAndPrivateKey),
		E.Chain(I.Ap[[]byte, E.Either[error, string]](test)),
	)

	fmt.Println(hashE)
}

func TestHashWithCertAndNoKey(t *testing.T) {

	privKeyE := encrypt.PrivateKey()

	dataE := F.Pipe2(
		privKeyE,
		E.Map[error](func(key []byte) contract.RawMap {
			data := make(map[string]any)

			// prepare input data
			data[common.KeyCert] = D.DefaultCertificate

			return data
		}),
		E.Map[error](CreateResourceDataMock),
	)

	test := []byte("Carsten")

	hashE := F.Pipe2(
		dataE,
		E.Map[error](createHashWithCertAndPrivateKey),
		E.Chain(I.Ap[[]byte, E.Either[error, string]](test)),
	)

	fmt.Println(hashE)
}

func TestHashWithCert(t *testing.T) {

	data := make(map[string]any)

	// prepare input data
	data[common.KeyCert] = D.DefaultCertificate

	test := []byte("Carsten")

	hashE := F.Pipe3(
		data,
		CreateResourceDataMock,
		createHashWithCert(&defaultContext),
		I.Ap[[]byte, E.Either[error, string]](test),
	)

	fmt.Println(hashE)
}
