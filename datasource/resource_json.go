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
	J "github.com/terraform-provider-hpcr/fp/json"
)

var (
	jsonBytes = F.Flow3(
		getJsonE,
		E.Map[error](F.Ref[any]),
		E.Chain(J.Stringify[any]),
	)
)

func ResourceJson() *schema.Resource {
	return &schema.Resource{
		Create: jsonUnencrypted.F1,
		Read:   jsonUnencrypted.F2,
		Delete: jsonUnencrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyJson:     &schemaJsonIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an base64 encoded token from the JSON serialization of the input",
	}
}

func ResourceJsonEncrypted() *schema.Resource {
	return &schema.Resource{
		Create: jsonEncrypted.F1,
		Read:   jsonEncrypted.F2,
		Delete: jsonEncrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyJson:     &schemaJsonIn,
			common.KeyCert:     &schemaCertIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an base64 encoded token from the JSON serialization of the input",
	}
}

func resourceEncJson(d fp.ResourceData) ResourceDataE {

	// marshal input text
	jsonE := jsonBytes(d)

	return F.Pipe2(
		jsonE,
		E.Chain(createHashWithCert(d)),
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateEncryptedResource(d)(jsonE),
			getResourceData(d),
		),
		),
	)
}

func resourceJson(d fp.ResourceData) ResourceDataE {

	// marshal input text
	jsonE := jsonBytes(d)

	return F.Pipe2(
		jsonE,
		createHashE,
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateBase64Resource(d)(jsonE),
			getResourceData(d),
		),
		),
	)
}

var (
	jsonUnencrypted = resourceLifeCycle(resourceJson)
	jsonEncrypted   = resourceLifeCycle(resourceEncJson)
)
