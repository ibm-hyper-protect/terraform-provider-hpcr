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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	S "github.com/terraform-provider-hpcr/fp/string"
)

func ResourceTgzEncrypted() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTgzEncryptedRead,
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
		Read:   resourceTgzRead,
		Create: resourceTgzCreate,
		Delete: resourceDeleteNoOp,
		Schema: map[string]*schema.Schema{
			common.KeyFolder:   &schemaFolderIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates a base64 encoded string from the TGZed files in the folder.",
	}
}

func resourceTgzEncryptedRead(d *schema.ResourceData, m any) error {
	// marshal input folder
	tarE := F.Pipe1(
		d,
		tarFolder,
	)

	// get the encryption function
	encryptedE := F.Pipe5(
		d,
		getPubKey,
		E.Map[error](S.ToBytes),
		E.Map[error](encrypt.OpenSSLEncryptBasic),
		E.Ap[error, []byte, E.Either[error, string]](tarE),
		E.Flatten[error, string],
	)

	// render the content
	renderedE := F.Pipe2(
		encryptedE,
		E.Map[error](setRendered),
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	// encode as sha256
	sha256E := F.Pipe2(
		tarE,
		computeSha256,
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	return F.Pipe1(
		seqResourceData([]E.Either[error, *schema.ResourceData]{renderedE, sha256E}),
		E.ToError[[]*schema.ResourceData],
	)
}

func resourceTgz(d *schema.ResourceData) E.Either[error, *schema.ResourceData] {

	// marshal input folder
	tarE := F.Pipe1(
		d,
		tarFolder,
	)

	// compute the checksum
	hashE := F.Pipe1(
		tarE,
		E.Map[error](createHash),
	)

	return F.Pipe1(
		hashE,
		E.Chain(func(checksum string) E.Either[error, *schema.ResourceData] {
			// get the sha256
			current := d.Get(common.KeySha256)
			if current != checksum {
				// requires update
				fmt.Println("updating resource")
				// render the content to base64
				renderedE := F.Pipe3(
					tarE,
					E.Map[error](common.Base64Encode),
					E.Map[error](setRendered),
					fp.ResourceDataAp[*schema.ResourceData](d),
				)

				// encode as sha256
				sha256E := F.Pipe2(
					hashE,
					E.Map[error](setSha256),
					fp.ResourceDataAp[*schema.ResourceData](d),
				)

				return F.Pipe1(
					seqResourceData([]E.Either[error, *schema.ResourceData]{renderedE, sha256E}),
					E.MapTo[error, []*schema.ResourceData](d),
				)
			}
			// nothing to do
			return E.Of[error](d)
		}),
	)
}

func resourceTgzCreate(d *schema.ResourceData, m any) error {
	return F.Pipe3(
		d,
		setUniqueID,
		E.Chain(resourceTgz),
		E.ToError[*schema.ResourceData],
	)
}

func resourceTgzRead(d *schema.ResourceData, m any) error {
	return F.Pipe2(
		d,
		resourceTgz,
		E.ToError[*schema.ResourceData],
	)
}
