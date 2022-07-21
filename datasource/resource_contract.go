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
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	S "github.com/terraform-provider-hpcr/fp/string"
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

func resourceEncContract(d fp.ResourceData) ResourceDataE {

	// marshal input text
	contractE := contractBytes(d)

	return F.Pipe2(
		contractE,
		E.Chain(createHashWithCert(d)),
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateEncryptedResource(d)(contractE),
			getResourceData(d),
		),
		),
	)
}

var (
	contractEncrypted = resourceLifeCycle(resourceEncContract)
)
