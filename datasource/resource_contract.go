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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/contract"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	O "github.com/terraform-provider-hpcr/fp/option"
	S "github.com/terraform-provider-hpcr/fp/string"
	Y "github.com/terraform-provider-hpcr/fp/yaml"
)

var (
	contractBytes = F.Flow2(
		getContractE,
		E.Map[error](S.ToBytes),
	)
)

func ResourceContractEncrypted() *schema.Resource {
	return &schema.Resource{
		Create: contractEncrypted.F1,
		Read:   contractEncrypted.F2,
		Delete: contractEncrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyContract: &schemaContractIn,
			common.KeyCert:     &schemaCertIn,
			common.KeyPrivKey:  &schemaPrivKeyIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an encrypted and signed user data field",
	}
}

// callback to update a resource using encryption base64 encoding
func updateContract(d fp.ResourceData) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
	return updateResource(d)(func(data []byte) E.Either[error, string] {

		// marshal key or create the private key
		privKeyE := F.Pipe2(
			getPrivKeyE(d),
			E.Map[error](S.ToBytes),
			E.Alt(encrypt.PrivateKey),
		)

		// deserialize the contract into a map
		contractE := F.Pipe2(
			data,
			Y.Parse[contract.RawMap],
			E.Map[error](F.Deref[contract.RawMap]),
		)

		// create the function that can execute the signature
		resE := F.Pipe10(
			d,
			getCertificateE,
			E.Map[error](S.ToBytes),
			E.Map[error](encrypt.OpenSSLEncryptBasic),
			E.Map[error](contract.EncryptAndSignContract),
			E.Ap[error, []byte, func(contract.RawMap) E.Either[error, contract.RawMap]](privKeyE),
			E.Ap[error, contract.RawMap, E.Either[error, contract.RawMap]](contractE),
			E.Flatten[error, contract.RawMap],
			E.Map[error](F.Ref[contract.RawMap]),
			E.Chain(Y.Stringify[contract.RawMap]),
			E.Map[error](B.ToString),
		)

		return resE
	})
}

func resourceEncContract(d fp.ResourceData) ResourceDataE {

	// marshal input text
	contractE := contractBytes(d)

	return F.Pipe2(
		contractE,
		E.Chain(createHashWithCert(d)),
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateContract(d)(contractE),
			getResourceData(d),
		),
		),
	)
}

var (
	contractEncrypted = resourceLifeCycle(resourceEncContract)
)
