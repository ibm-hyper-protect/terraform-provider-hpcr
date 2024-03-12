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
	"crypto/sha256"
	"encoding/hex"

	B "github.com/IBM/fp-go/bytes"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	O "github.com/IBM/fp-go/option"
	P "github.com/IBM/fp-go/predicate"
	S "github.com/IBM/fp-go/string"
	T "github.com/IBM/fp-go/tuple"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/data"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/validation"
)

// shortcuts
type ResourceDataE = E.Either[error, fp.ResourceData]
type ResourceLifeCycle = T.Tuple3[func(*schema.ResourceData, any) error, func(*schema.ResourceData, any) error, func(*schema.ResourceData, any) error]

// produces a new UUID
var (
	uuidE      = E.Eitherize0(uuid.GenerateUUID)
	toContextE = common.ToTypeE[*Context]
)

// assigns a new uuid to a resource
func setUniqueID(d fp.ResourceData) ResourceDataE {
	return F.Pipe1(
		uuidE(),
		E.Map[error](func(id string) fp.ResourceData {
			d.SetID(id)
			return d
		}),
	)
}

func createHash(data []byte) string {
	checksum := sha256.Sum256(data)
	return hex.EncodeToString(checksum[:])
}

var (
	seqResourceData   = E.SequenceArray[error, fp.ResourceData]
	setRendered       = fp.ResourceDataSet[string](common.KeyRendered)
	setSha256         = fp.ResourceDataSet[string](common.KeySha256)
	setChecksum       = fp.ResourceDataSet[string](common.KeyChecksum)
	setChecksums      = fp.ResourceDataSet[map[string]string](common.KeyChecksums)
	getJsonE          = fp.ResourceDataGetE[any](common.KeyJSON)
	getTextE          = fp.ResourceDataGetE[string](common.KeyText)
	getContractE      = fp.ResourceDataGetE[string](common.KeyContract)
	getPrivKeyE       = fp.ResourceDataGetE[string](common.KeyPrivKey)
	getCertExpiryDays = fp.ResourceDataGetE[string](common.KeyCertExpiryDays)
	getCaPrivKeyE     = fp.ResourceDataGetE[string](common.KeyCaPrivKey)
	getCaCertE        = fp.ResourceDataGetE[string](common.KeyCaCert)
	getFolderE        = fp.ResourceDataGetE[string](common.KeyFolder)
	getAttestationE   = fp.ResourceDataGetE[string](common.KeyAttestation)
	getCertificateE   = fp.ResourceDataGetE[string](common.KeyCert)

	getSha256O  = fp.ResourceDataGetO[string](common.KeySha256)
	getPrivKeyO = fp.ResourceDataGetO[string](common.KeyPrivKey)

	createHashE = E.Map[error](createHash)

	schemaJsonIn = schema.Schema{
		Type:        schema.TypeMap,
		Required:    true,
		ForceNew:    true,
		Sensitive:   true,
		Description: "JSON Document to archive",
	}

	schemaTextIn = schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Sensitive:   true,
		Description: "Text to archive",
	}

	schemaContractIn = schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Sensitive:        true,
		Description:      "YAML serialization of the contract",
		ValidateDiagFunc: validation.DiagContract,
	}

	schemaFolderIn = schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Path to the folder to encrypt",
		ValidateDiagFunc: validation.DiagFolder,
	}

	schemaCertIn = schema.Schema{
		Type:             schema.TypeString,
		Description:      "Certificate used to encrypt the JSON document, in PEM format",
		Optional:         true,
		ForceNew:         true,
		Default:          data.DefaultCertificate,
		ValidateDiagFunc: validation.DiagCertificate,
	}

	schemaPrivKeyIn = schema.Schema{
		Type:             schema.TypeString,
		Description:      "Private key used to sign the contract. If omitted, a temporally signing key is created.",
		Optional:         true,
		ForceNew:         true,
		Sensitive:        true,
		ValidateDiagFunc: validation.DiagPrivateKey,
	}

	schemaCertExpiryIn = schema.Schema{
		Type:        schema.TypeInt,
		Description: "Number of days for expiring the contract. If omitted, Signing key will be public key",
		Optional:    true,
		ForceNew:    true,
		Sensitive:   false,
	}

	schemaCaPrivKeyIn = schema.Schema{
		Type:             schema.TypeString,
		Description:      "Path of CA Private Key. If omitted, Signing key will be public key",
		Optional:         true,
		ForceNew:         true,
		Sensitive:        true,
		ValidateDiagFunc: validation.DiagPrivateKey,
	}

	schemaCaCertIn = schema.Schema{
		Type:             schema.TypeString,
		Description:      "Path of CA Certificate. If omitted, Signing key will be public key",
		Optional:         true,
		ForceNew:         true,
		Sensitive:        true,
		ValidateDiagFunc: validation.DiagCertificate,
	}

	schemaRenderedOut = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Rendered output of the resource",
	}

	schemaSha256Out = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "SHA256 of the input",
	}

	schemaChecksumOut = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "SHA256 of the output",
	}
)

func resourceDeleteNoOp(d *schema.ResourceData, m any) error {
	return nil
}

// returns a predicate that checks if the sha256 value in resource data matches the given value
func checksumMatch(d fp.ResourceData) func(string) bool {
	// get the sha
	sha256O := F.Pipe1(
		d,
		getSha256O,
	)
	// returns the comparator
	return func(checksum string) bool {
		return F.Pipe1(
			sha256O,
			O.Fold(F.ConstFalse, S.Equals(checksum)),
		)
	}
}

func updateResource(d fp.ResourceData) func(func([]byte) E.Either[error, string]) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
	// compute the applicatives
	apE := fp.ResourceDataAp[fp.ResourceData](d)
	apI := I.Ap[ResourceDataE](d)
	// final result
	resE := E.MapTo[error, []fp.ResourceData](d)

	return func(serialize func([]byte) E.Either[error, string]) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
		// construct the serialization callback
		serE := E.Chain(serialize)

		return func(dataE E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {

			// serialize the content
			return O.Map(func(checksum string) ResourceDataE {

				// rendering result
				resultE := F.Pipe1(
					dataE,
					serE,
				)

				// render the content using the serializer
				renderedE := F.Pipe2(
					resultE,
					E.Map[error](setRendered),
					apE,
				)

				// encode as sha256
				sha256E := F.Pipe2(
					checksum,
					setSha256,
					apI,
				)

				// update the checksum
				checksumE := F.Pipe1(
					resultE,
					E.Chain(F.Flow4(
						S.ToBytes,
						createHash,
						setChecksum,
						apI,
					)),
				)

				return F.Pipe1(
					seqResourceData([]ResourceDataE{renderedE, sha256E, checksumE}),
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
		E.Of[error, fp.ResourceData],
		F.Constant[ResourceDataE],
		O.GetOrElse[ResourceDataE],
	)

	// callback to update a resource using simple base64 encoding
	updateBase64Resource = F.Flow2(
		updateResource,
		I.Ap[func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE]](common.Base64EncodeE),
	)

	// callback to update a resource using plain text encoding
	updatePlainTextResource = F.Flow2(
		updateResource,
		I.Ap[func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE]](common.PlainTextEncodeE),
	)
)

// callback to update a resource using encryption base64 encoding
func updateEncryptedResource(ctx *Context) func(d fp.ResourceData) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
	return func(d fp.ResourceData) func(E.Either[error, []byte]) func(O.Option[string]) O.Option[ResourceDataE] {
		return updateResource(d)(func(data []byte) E.Either[error, string] {
			return F.Pipe4(
				d,
				getCertificateE,
				common.MapStgToBytesE,
				E.Map[error](ctx.EncryptBasic),
				E.Chain(I.Ap[E.Either[error, string]](data)),
			)
		})
	}
}

func resourceLifeCycle(f func(ctx *Context) func(fp.ResourceData) ResourceDataE) ResourceLifeCycle {

	// lift f into the context
	withCtx := F.Flow2(
		toContextE,
		E.Map[error](f),
	)

	create := func(d *schema.ResourceData, m any) error {

		return F.Pipe4(
			m,
			withCtx,
			E.Ap[E.Either[error, fp.ResourceData]](F.Pipe2(
				d,
				fp.CreateResourceDataProxy,
				setUniqueID,
			)),
			E.Flatten[error, fp.ResourceData],
			E.ToError[fp.ResourceData],
		)

	}

	read := func(d *schema.ResourceData, m any) error {

		return F.Pipe3(
			m,
			withCtx,
			E.Chain(I.Ap[E.Either[error, fp.ResourceData]](F.Pipe1(
				d,
				fp.CreateResourceDataProxy,
			))),
			E.ToError[fp.ResourceData],
		)

	}

	delete := resourceDeleteNoOp

	return T.MakeTuple3(create, read, delete)
}

// computes a hash for the given bytes and includes the fingerprint of the certificate as part of the hash
func createHashWithCert(ctx *Context) func(d fp.ResourceData) func([]byte) E.Either[error, string] {
	return func(d fp.ResourceData) func([]byte) E.Either[error, string] {
		// get the fingerprint
		fpE := F.Pipe3(
			d,
			getCertificateE,
			common.MapStgToBytesE,
			E.Chain(ctx.CertFingerprint),
		)
		// combine the fingerprint with the actual data
		return func(data []byte) E.Either[error, string] {
			return F.Pipe1(
				fpE,
				E.Map[error](F.Flow2(
					F.Bind2nd(B.Monoid.Concat, data),
					createHash,
				)),
			)
		}
	}
}

// computes a hash for the given bytes and includes the fingerprint of the certificate as part of the hash
func createHashWithCertAndPrivateKey(ctx *Context) func(d fp.ResourceData) func([]byte) E.Either[error, string] {

	return func(d fp.ResourceData) func([]byte) E.Either[error, string] {
		// get the fingerprint for the certificate
		certE := F.Pipe3(
			d,
			getCertificateE,
			common.MapStgToBytesE,
			E.Chain(ctx.CertFingerprint),
		)

		// get the fingerprint for the private key
		privKeyE := F.Pipe4(
			d,
			getPrivKeyE,
			common.MapStgToBytesE,
			E.Chain(ctx.PrivKeyFingerprint),
			E.Alt(F.Constant(E.Of[error](B.Monoid.Empty()))),
		)
		// combine into one
		fp := E.Sequence2(func(left, right []byte) E.Either[error, []byte] {
			return E.Of[error](B.Monoid.Concat(left, right))
		})

		// combine the fingerprint with the actual data
		return func(data []byte) E.Either[error, string] {
			return F.Pipe1(
				fp(certE, privKeyE),
				E.Map[error](F.Flow2(
					F.Bind2nd(B.Monoid.Concat, data),
					createHash),
				),
			)
		}
	}
}
