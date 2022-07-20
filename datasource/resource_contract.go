//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
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
