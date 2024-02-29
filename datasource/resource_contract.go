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
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/contract"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
)

var (
	contractBytes = F.Flow2(
		getContractE,
		common.MapStgToBytesE,
	)
)

func ResourceContractEncrypted() *schema.Resource {
	return &schema.Resource{
		Create: contractEncrypted.F1,
		Read:   contractEncrypted.F2,
		Delete: contractEncrypted.F3,
		Schema: map[string]*schema.Schema{
			// input parameters
			common.KeyContract:       &schemaContractIn,
			common.KeyCert:           &schemaCertIn,
			common.KeyPrivKey:        &schemaPrivKeyIn,
			common.KeyCertExpiryDays: &schemaCertExpiryIn,
			common.KeyCaPrivKey:      &schemaCAPrivateKeyIn,
			common.KeyCaCert:         &schemaCACertificateIn,
			// output parameters
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
			common.KeyChecksum: &schemaChecksumOut,
		},
		Description: "Generates an encrypted and signed user data field",
	}
}

func encryptAndSignContract(
	enc func([]byte) func([]byte) E.Either[error, string],
	signer func([]byte) func([]byte) E.Either[error, []byte],
	pubKey func([]byte) E.Either[error, []byte],
) func([]byte) func([]byte) func(contract.RawMap) E.Either[error, contract.RawMap] {
	return func(cert []byte) func([]byte) func(contract.RawMap) E.Either[error, contract.RawMap] {
		return contract.EncryptAndSignContract(enc(cert), signer, pubKey)
	}
}

// callback to update a resource using encryption base64 encoding
func updateContract(ctx *Context) func(d fp.ResourceData) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
	// contextualize the encrypter
	encryptAndSign := encryptAndSignContract(ctx.EncryptBasic, ctx.SignDigest, ctx.PubKey)

	return func(d fp.ResourceData) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
		return updateResource(d)(func(data []byte) E.Either[error, string] {

			// marshal key or create the private key
			privKeyE := F.Pipe2(
				getPrivKeyE(d),
				common.MapStgToBytesE,
				E.Alt(ctx.PrivKey),
			)

			// deserialize the contract into a map
			contractE := F.Pipe1(
				data,
				contract.ParseRawMapE,
			)

			// create the function that can execute the signature
			resE := F.Pipe8(
				d,
				getCertificateE,
				common.MapStgToBytesE,
				E.Map[error](encryptAndSign),
				E.Ap[func(contract.RawMap) E.Either[error, contract.RawMap]](privKeyE),
				E.Ap[E.Either[error, contract.RawMap]](contractE),
				E.Flatten[error, contract.RawMap],
				E.Chain(contract.StringifyRawMapE),
				common.MapBytesToStgE,
			)

			return resE
		})
	}
}

func resourceEncContract(ctx *Context) func(d fp.ResourceData) ResourceDataE {

	// contextualize
	hashWithCertAndPrivateKey := createHashWithCertAndPrivateKey(ctx)
	update := updateContract(ctx)

	return func(d fp.ResourceData) ResourceDataE {
		// marshal input text
		contractE := contractBytes(d)

		return F.Pipe2(
			contractE,
			E.Chain(hashWithCertAndPrivateKey(d)),
			E.Chain(F.Flow3(
				checksumMatchO(d),
				update(d)(contractE),
				getResourceData(d),
			),
			),
		)
	}
}

var (
	contractEncrypted = resourceLifeCycle(resourceEncContract)
)
