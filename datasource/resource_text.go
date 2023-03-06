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
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
)

var (
	textBytes = F.Flow2(
		getTextE,
		common.MapStgToBytesE,
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

func resourceEncText(ctx *Context) func(d fp.ResourceData) ResourceDataE {

	// get the update method depending on the context
	update := updateEncryptedResource(ctx)
	hashWithCert := createHashWithCert(ctx)

	return func(d fp.ResourceData) ResourceDataE {
		// marshal input text
		textE := textBytes(d)

		return F.Pipe2(
			textE,
			E.Chain(hashWithCert(d)),
			E.Chain(F.Flow3(
				checksumMatchO(d),
				update(d)(textE),
				getResourceData(d),
			),
			),
		)
	}
}

func resourceText(ctx *Context) func(d fp.ResourceData) ResourceDataE {

	return func(d fp.ResourceData) ResourceDataE {
		// marshal input text
		textE := textBytes(d)

		return F.Pipe2(
			textE,
			createHashE,
			E.Chain(F.Flow3(
				checksumMatchO(d),
				updatePlainTextResource(d)(textE),
				getResourceData(d),
			),
			),
		)
	}
}

var (
	textUnencrypted = resourceLifeCycle(resourceText)
	textEncrypted   = resourceLifeCycle(resourceEncText)
)
