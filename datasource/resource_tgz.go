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
	"bytes"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/archive"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

var (
	// marshal input folder
	tarFolder = F.Flow4(
		getFolderE,
		E.Map[error](archive.TarFolder[*bytes.Buffer]),
		E.Chain(func(tar func(*bytes.Buffer) E.Either[error, *bytes.Buffer]) E.Either[error, *bytes.Buffer] {
			return tar(new(bytes.Buffer))
		}),
		E.Map[error]((*bytes.Buffer).Bytes),
	)
)

func ResourceTgzEncrypted() *schema.Resource {
	return &schema.Resource{
		Create: tgzEncrypted.F1,
		Read:   tgzEncrypted.F2,
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
		Create: tgzUnencrypted.F1,
		Read:   tgzUnencrypted.F2,
		Delete: tgzUnencrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyFolder:   &schemaFolderIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates a base64 encoded string from the TGZed files in the folder.",
	}
}

func resourceEncTgz(d fp.ResourceData) ResourceDataE {

	// marshal input folder
	tarE := tarFolder(d)

	return F.Pipe2(
		tarE,
		E.Chain(createHashWithCert(d)),
		E.Chain(F.Flow3(
			checksumMatchO(d),
			updateEncryptedResource(d)(tarE),
			getResourceData(d),
		),
		),
	)
}

func resourceTgz(d fp.ResourceData) ResourceDataE {

	// marshal input folder
	tarE := tarFolder(d)

	return F.Pipe2(
		tarE,
		createHashE,
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
