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
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/data"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	O "github.com/terraform-provider-hpcr/fp/option"
	P "github.com/terraform-provider-hpcr/fp/predicate"
	S "github.com/terraform-provider-hpcr/fp/string"
	T "github.com/terraform-provider-hpcr/fp/tuple"
	"github.com/terraform-provider-hpcr/validation"
)

// shortcuts
type ResourceDataE = E.Either[error, *schema.ResourceData]
type ResourceLifeCycle = T.Tuple3[func(d *schema.ResourceData, m any) error, func(d *schema.ResourceData, m any) error, func(d *schema.ResourceData, m any) error]

var uuidE = E.Eitherize0(uuid.GenerateUUID)

// assigns a new uuid to a resource
func setUniqueID(d *schema.ResourceData) ResourceDataE {
	return F.Pipe1(
		uuidE(),
		E.Map[error](func(id string) *schema.ResourceData {
			d.SetId(id)
			return d
		}),
	)
}

func createHash(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

var (
	seqResourceData = E.SequenceArray[error, *schema.ResourceData]()
	setRendered     = fp.ResourceDataSet[string](common.KeyRendered)
	setText         = fp.ResourceDataSet[string](common.KeyText)
	setSha256       = fp.ResourceDataSet[string](common.KeySha256)
	getJsonE        = fp.ResourceDataGetE[any](common.KeyJson)
	getTextE        = fp.ResourceDataGetE[string](common.KeyText)
	getFolderE      = fp.ResourceDataGetE[string](common.KeyFolder)
	getCertificateE = fp.ResourceDataGetE[string](common.KeyCert)

	getSha256O = fp.ResourceDataGetO[string](common.KeySha256)

	// encode as sha256
	computeSha256 = F.Flow3(
		E.Map[error](sha256.Sum256),
		E.Map[error](func(hash [sha256.Size]byte) string { return fmt.Sprintf("%x", hash) }),
		E.Map[error](setSha256),
	)

	// encode as text
	computeText = F.Flow2(
		E.Map[error](B.ToString),
		E.Map[error](setText),
	)

	schemaCertIn = schema.Schema{
		Type:             schema.TypeString,
		Description:      "Certificate used to encrypt the JSON document in PEM format",
		Optional:         true,
		Default:          data.DefaultCertificate,
		ValidateDiagFunc: validation.DiagCertificate,
	}

	schemaTextOut = schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	schemaRenderedOut = schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	schemaSha256Out = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "SHA256 of the input",
	}

	schemaJsonIn = schema.Schema{
		Type:        schema.TypeMap,
		Required:    true,
		Description: "JSON Document to archive",
	}

	schemaTextIn = schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Text to archive",
	}

	schemaFolderIn = schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Path to the folder to encrypt",
		ValidateDiagFunc: validation.DiagFolder,
	}
)

func resourceDeleteNoOp(d *schema.ResourceData, m any) error {
	return nil
}

// returns a predicate that checks if the sha256 value in resource data matches the given value
func checksumMatch(d *schema.ResourceData) func(string) bool {
	// get the sha
	sha256O := F.Pipe1(
		d,
		getSha256O,
	)
	// returns the comparator
	return func(checksum string) bool {
		return F.Pipe1(
			sha256O,
			O.Fold(F.ConstFalse, F.Bind1st(S.Equals, checksum)),
		)
	}
}

func updateResource(d *schema.ResourceData) func(func([]byte) E.Either[error, string]) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
	// compute the applicatives
	apE := fp.ResourceDataAp[*schema.ResourceData](d)
	apI := I.Ap[*schema.ResourceData, ResourceDataE](d)
	// final result
	resE := E.MapTo[error, []*schema.ResourceData](d)

	return func(serialize func([]byte) E.Either[error, string]) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
		// construct the serialization callback
		serE := E.Chain(serialize)

		return func(dataE E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {

			// serialize the content
			return O.Map(func(checksum string) ResourceDataE {

				// render the content using the serializer
				renderedE := F.Pipe3(
					dataE,
					serE,
					E.Map[error](setRendered),
					apE,
				)

				// encode as sha256
				sha256E := F.Pipe2(
					checksum,
					setSha256,
					apI,
				)

				return F.Pipe1(
					seqResourceData([]ResourceDataE{renderedE, sha256E}),
					resE,
				)
			})
		}
	}
}

var (
	// checksum match as an optional
	checksumMatchO = F.Flow3(
		checksumMatch,
		P.Not[string],
		O.FromPredicate[string],
	)

	// common fallback
	getResourceData = F.Flow3(
		E.Of[error, *schema.ResourceData],
		F.Constant[ResourceDataE],
		O.GetOrElse[ResourceDataE],
	)

	// callback to update a resource using simple base64 encoding
	updateBase64Resource = F.Flow2(
		updateResource,
		I.Ap[func([]byte) E.Either[error, string], func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE]](common.Base64EncodeE),
	)
)

// callback to update a resource using encryption base64 encoding
func updateEncryptedResource(d *schema.ResourceData) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
	return updateResource(d)(func(data []byte) E.Either[error, string] {
		return F.Pipe4(
			d,
			getCertificateE,
			E.Map[error](S.ToBytes),
			E.Map[error](encrypt.OpenSSLEncryptBasic),
			E.Chain(I.Ap[[]byte, E.Either[error, string]](data)),
		)
	})
}

func resourceLifeCycle(f func(*schema.ResourceData) ResourceDataE) ResourceLifeCycle {

	create := func(d *schema.ResourceData, m any) error {
		return F.Pipe3(
			d,
			setUniqueID,
			E.Chain(f),
			E.ToError[*schema.ResourceData],
		)
	}

	read := func(d *schema.ResourceData, m any) error {
		return F.Pipe2(
			d,
			f,
			E.ToError[*schema.ResourceData],
		)
	}
	delete := resourceDeleteNoOp

	return T.MakeTuple3(create, read, delete)
}
