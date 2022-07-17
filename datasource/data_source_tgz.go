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
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/archive"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	S "github.com/terraform-provider-hpcr/fp/string"
)

var (
	// marshal input folder
	tarFolder = F.Flow4(
		getFolder,
		E.Map[error](archive.TarFolder[*bytes.Buffer]),
		E.Chain(I.Ap[*bytes.Buffer, E.Either[error, *bytes.Buffer]](new(bytes.Buffer))),
		E.Map[error]((*bytes.Buffer).Bytes),
	)
)

func DataSourceTgzEncrypted() *schema.Resource {
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

func DataSourceTgz() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTgzRead,
		Schema: map[string]*schema.Schema{
			common.KeyFolder:   &schemaFolderIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates a base64 encoded string from the TGZed files in the folder.",
	}
}

func dataSourceTgzEncryptedRead(d *schema.ResourceData, m any) error {
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

func dataSourceTgzRead(d *schema.ResourceData, m any) error {

	// marshal input folder
	tarE := F.Pipe1(
		d,
		tarFolder,
	)

	F.Pipe3(
		tarE,
		E.Map[error](sha256.Sum256),
		E.Map[error](func(hash [sha256.Size]byte) string { return fmt.Sprintf("%x", hash) }),
		E.Map[error](func(id string) string {
			d.SetId(id)
			return id
		}),
	)

	// render the content to base64
	renderedE := F.Pipe3(
		tarE,
		E.Map[error](common.Base64Encode),
		E.Map[error](setRendered),
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	// encode as sha256
	sha256E := F.Pipe2(
		tarE,
		computeSha256,
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	fmt.Println(d.HasChange(common.KeySha256))

	return F.Pipe1(
		seqResourceData([]E.Either[error, *schema.ResourceData]{renderedE, sha256E}),
		E.ToError[[]*schema.ResourceData],
	)
}
