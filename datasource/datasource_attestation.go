// Copyright 2023 IBM Corp.
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
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/attestation"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/data"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/validation"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	O "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/option"
	S "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/string"
)

var (
	schemaAttestationIn = schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The encrypted or unencrypted attestation record",
	}
	schemaAttestationCertIn = schema.Schema{
		Type:             schema.TypeString,
		Description:      "Certificate used to validate the attestation signature, in PEM format",
		Optional:         true,
		Default:          data.DefaultCertificate,
		ValidateDiagFunc: validation.DiagCertificate,
	}
	schemaAttestationPrivKeyIn = schema.Schema{
		Type:             schema.TypeString,
		Description:      "Private key used to decrypt an encrypted attestation record. If missing the attestation record is assumed to be unencrypted.",
		Optional:         true,
		Sensitive:        true,
		ValidateDiagFunc: validation.DiagPrivateKey,
	}
	schemaChecksumsOut = schema.Schema{
		Type:        schema.TypeMap,
		Description: "Map from filename to checksum of the attestation record.",
		Computed:    true,
	}
)

func DatasourceAttestation() *schema.Resource {
	return &schema.Resource{
		Read: handleAttestation,
		Schema: map[string]*schema.Schema{
			// input parameters
			common.KeyAttestation: &schemaAttestationIn,
			common.KeyPrivKey:     &schemaAttestationPrivKeyIn,
			common.KeyCert:        &schemaAttestationCertIn,
			// output parameters
			common.KeyChecksums: &schemaChecksumsOut,
		},
		Description: "handles the analysis of an attestation record.",
	}
}

// parseAttestationRecord
var parseAttestationRecord = F.Flow2(
	attestation.ParseAttestationRecord,
	E.Of[error, attestation.ChecksumMap],
)

func handleAttestationWithContext(ctx *Context) func(data fp.ResourceData) ResourceDataE {
	// the decryptor
	decryptAttestation := attestation.DecryptAttestation(ctx.DecryptBasic)

	return func(data fp.ResourceData) ResourceDataE {
		// some applicatives
		apE := fp.ResourceDataAp[fp.ResourceData](data)
		// final result
		resE := E.MapTo[error, []fp.ResourceData](data)

		// the attestation record
		attestationE := F.Pipe1(
			data,
			getAttestationE,
		)

		// attestation parser function
		attestationParser := F.Pipe2(
			data,
			getPrivKeyO,
			O.Fold(F.Constant(parseAttestationRecord), F.Flow2(
				S.ToBytes,
				decryptAttestation,
			)),
		)

		// output records
		checksumMapE := F.Pipe3(
			attestationE,
			E.Chain(attestationParser),
			E.Map[error](setChecksums),
			apE,
		)

		// combine all outputs
		return F.Pipe2(
			[]ResourceDataE{checksumMapE},
			seqResourceData,
			resE,
		)
	}
}

// handleAttestation is the data source callback, it dispatches to a more convenient API
func handleAttestation(data *schema.ResourceData, ctx any) error {
	// lift f into the context
	return F.Pipe5(
		ctx,
		toContextE,
		E.Map[error](handleAttestationWithContext),
		E.Ap[error, fp.ResourceData, ResourceDataE](F.Pipe2(
			data,
			fp.CreateResourceDataProxy,
			setUniqueID,
		)),
		E.Flatten[error, fp.ResourceData],
		E.ToError[fp.ResourceData],
	)
}
