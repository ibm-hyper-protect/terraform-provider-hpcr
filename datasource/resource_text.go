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
	textBytes = F.Flow2(
		getTextE,
		E.Map[error](S.ToBytes),
	)
)

func ResourceText() *schema.Resource {
	return &schema.Resource{
		Create: textUnencrypted.F1,
		Read:   textUnencrypted.F2,
		Delete: textUnencrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyText:     &schemaTextIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an base64 encoded token from text input",
	}
}

func ResourceTextEncrypted() *schema.Resource {
	return &schema.Resource{
		Create: textEncrypted.F1,
		Read:   textEncrypted.F2,
		Delete: textEncrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyText:     &schemaTextIn,
			common.KeyCert:     &schemaCertIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an encrypted token from text input",
	}
}

func resourceEncText(d fp.ResourceData) ResourceDataE {

	// marshal input text
	textE := textBytes(d)

	return F.Pipe2(
		textE,
		E.Chain(createHashWithCert(d)),
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateEncryptedResource(d)(textE),
			getResourceData(d),
		),
		),
	)
}

func resourceText(d fp.ResourceData) ResourceDataE {

	// marshal input text
	textE := textBytes(d)

	return F.Pipe2(
		textE,
		createHashE,
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateBase64Resource(d)(textE),
			getResourceData(d),
		),
		),
	)
}

var (
	textUnencrypted = resourceLifeCycle(resourceText)
	textEncrypted   = resourceLifeCycle(resourceEncText)
)
