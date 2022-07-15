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
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	J "github.com/terraform-provider-hpcr/fp/json"
	S "github.com/terraform-provider-hpcr/fp/string"
)

var (
	getJsonBytes = F.Flow4(
		setUniqueID,
		getJson,
		E.Map[error](F.Ref[any]),
		E.Chain(J.Stringify[any]),
	)
)

func DataSourceJson() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJsonRead,
		Schema: map[string]*schema.Schema{
			common.KeyJson:     &schemaJsonIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an base64 encoded token from the JSON serialization of the input",
	}
}

func DataSourceJsonEncrypted() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJsonEncryptedRead,
		Schema: map[string]*schema.Schema{
			common.KeyJson:     &schemaJsonIn,
			common.KeyCert:     &schemaCertIn,
			common.KeyText:     &schemaTextOut,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an base64 encoded token from the JSON serialization of the input",
	}
}

func dataSourceJsonRead(d *schema.ResourceData, m any) error {

	// marshal data to bytes
	dataE := F.Pipe1(
		d,
		getJsonBytes,
	)

	// encode and set as rendered
	renderedE := F.Pipe3(
		dataE,
		E.Map[error](common.Base64Encode),
		E.Map[error](setRendered),
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	// encode and set as text
	textE := F.Pipe2(
		dataE,
		computeText,
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	// encode as sha256
	sha256E := F.Pipe2(
		dataE,
		computeSha256,
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	return F.Pipe1(
		seqResourceData([]E.Either[error, *schema.ResourceData]{renderedE, textE, sha256E}),
		E.ToError[[]*schema.ResourceData],
	)
}

func dataSourceJsonEncryptedRead(d *schema.ResourceData, m any) error {

	// marshal data to bytes
	dataE := F.Pipe1(
		d,
		getJsonBytes,
	)

	// get the encryption function
	encryptE := F.Pipe3(
		d,
		getPubKey,
		E.Map[error](S.ToBytes),
		E.Map[error](encrypt.OpenSSLEncryptBasic),
	)

	// encrypt the data and set as rendered
	renderedE := F.Pipe4(
		encryptE,
		E.Ap[error, []byte, E.Either[error, string]](dataE),
		E.Flatten[error, string],
		E.Map[error](setRendered),
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	// encode and set as text
	textE := F.Pipe3(
		dataE,
		E.Map[error](B.ToString),
		E.Map[error](setText),
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	// encode as sha256
	sha256E := F.Pipe2(
		dataE,
		computeSha256,
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	return F.Pipe1(
		seqResourceData([]E.Either[error, *schema.ResourceData]{renderedE, textE, sha256E}),
		E.ToError[[]*schema.ResourceData],
	)
}
