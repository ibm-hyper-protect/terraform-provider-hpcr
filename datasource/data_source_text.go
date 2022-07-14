package datasource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	S "github.com/terraform-provider-hpcr/fp/string"
)

var (
	getTextBytes = F.Flow3(
		setUniqueID,
		getText,
		E.Map[error](S.ToBytes),
	)
)

func DataSourceText() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTextRead,
		Schema: map[string]*schema.Schema{
			common.KeyText:     &schemaTextIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an base64 encoded token from text input",
	}
}

func DataSourceTextEncrypted() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTextEncryptedRead,
		Schema: map[string]*schema.Schema{
			common.KeyText:     &schemaTextIn,
			common.KeyCert:     &schemaCertIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an encrypted token from text input",
	}
}

func dataSourceTextRead(d *schema.ResourceData, m any) error {

	// marshal data to bytes
	dataE := F.Pipe1(
		d,
		getTextBytes,
	)

	// encode and set as rendered
	renderedE := F.Pipe3(
		dataE,
		E.Map[error](common.Base64Encode),
		E.Map[error](setRendered),
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	// encode as sha256
	sha256E := F.Pipe2(
		dataE,
		computeSha256,
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	return F.Pipe1(
		seqResourceData([]E.Either[error, *schema.ResourceData]{renderedE, sha256E}),
		E.ToError[[]*schema.ResourceData],
	)
}

func dataSourceTextEncryptedRead(d *schema.ResourceData, m any) error {
	// marshal data to bytes
	dataE := F.Pipe1(
		d,
		getTextBytes,
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

	// encode as sha256
	sha256E := F.Pipe2(
		dataE,
		computeSha256,
		fp.ResourceDataAp[*schema.ResourceData](d),
	)

	return F.Pipe1(
		seqResourceData([]E.Either[error, *schema.ResourceData]{renderedE, sha256E}),
		E.ToError[[]*schema.ResourceData],
	)
}
