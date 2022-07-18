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
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func ResourceTgzEncrypted() *schema.Resource {
	return &schema.Resource{
		Read:   tgzEncrypted.F1,
		Create: tgzEncrypted.F2,
		Delete: tgzEncrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyFolder:   &schemaFolderIn,
			common.KeyCert:     &schemaCertIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates a encrypted token from the TGZed files in the folder.",
	}
}

func ResourceTgz() *schema.Resource {
	return &schema.Resource{
		Read:   tgzUnencrypted.F1,
		Create: tgzUnencrypted.F2,
		Delete: tgzUnencrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyFolder:   &schemaFolderIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates a base64 encoded string from the TGZed files in the folder.",
	}
}

func resourceEncTgz(d *schema.ResourceData) ResourceDataE {

	// marshal input folder
	tarE := tarFolder(d)

	return F.Pipe2(
		tarE,
		E.Map[error](createHash),
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateEncryptedResource(d)(tarE),
			getResourceData(d),
		),
		),
	)
}

func resourceTgz(d *schema.ResourceData) ResourceDataE {

	// marshal input folder
	tarE := tarFolder(d)

	return F.Pipe2(
		tarE,
		E.Map[error](createHash),
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateBase64Resource(d)(tarE),
			getResourceData(d),
		),
		),
	)
}

var (
	tgzUnencrypted = resourceLifeCycle(resourceTgz)
	tgzEncrypted   = resourceLifeCycle(resourceEncTgz)
)
