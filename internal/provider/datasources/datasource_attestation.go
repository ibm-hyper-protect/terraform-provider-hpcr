// Copyright 2026 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/ibm-hyper-protect/contract-go/v2/attestation"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var _ datasource.DataSource = &AttestationDataSource{}

func NewAttestationDataSource() datasource.DataSource {
	return &AttestationDataSource{}
}

type AttestationDataSource struct{}

type AttestationDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Attestation types.String `tfsdk:"attestation"`
	PrivKey     types.String `tfsdk:"privkey"`
	Cert        types.String `tfsdk:"cert"`
	Checksums   types.Map    `tfsdk:"checksums"`
}

func (d *AttestationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attestation"
}

func (d *AttestationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Handles the analysis of an HPCR attestation record (encrypted or unencrypted).",
		Description:         "Handles the analysis of an attestation record.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier",
			},
			"attestation": schema.StringAttribute{
				MarkdownDescription: "The encrypted or unencrypted attestation record",
				Description:         "The encrypted or unencrypted attestation record",
				Required:            true,
			},
			"privkey": schema.StringAttribute{
				MarkdownDescription: "Private key used to decrypt an encrypted attestation record. If missing the attestation record is assumed to be unencrypted.",
				Description:         "Private key used to decrypt an encrypted attestation record",
				Optional:            true,
				Sensitive:           true,
			},
			"cert": schema.StringAttribute{
				MarkdownDescription: "Certificate used to validate the attestation signature, in PEM format. Defaults to the default HPCR certificate if not specified.",
				Description:         "Certificate used to validate the attestation signature, in PEM format",
				Optional:            true,
			},
			"checksums": schema.MapAttribute{
				MarkdownDescription: "Map from filename to checksum of the attestation record",
				Description:         "Map from filename to checksum of the attestation record",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *AttestationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AttestationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	attestationData := data.Attestation.ValueString()
	privateKey := data.PrivKey.ValueString()

	var attestationRecords string
	var err error

	// Decrypt attestation if private key is provided
	if privateKey != "" {
		attestationRecords, err = attestation.HpcrGetAttestationRecords(attestationData, privateKey)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to decrypt attestation record",
				fmt.Sprintf("Error decrypting attestation: %s", err.Error()),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Decrypted attestation records - %s", attestationRecords))
	} else {
		// If no private key, assume attestation is already decrypted
		attestationRecords = attestationData
	}

	filteredAttestation := common.FilterChecksum(attestationRecords)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse attestation records",
			fmt.Sprintf("Error parsing attestation records as JSON: %s", err.Error()),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Filtered attestation records - %s", filteredAttestation))

	// Convert map[string]interface{} to map[string]string for Terraform
	checksumStrMap := make(map[string]string)
	for key, value := range filteredAttestation {
		checksumStrMap[key] = fmt.Sprintf("%v", value)
	}

	// Convert to Terraform types.Map
	checksums, diags := types.MapValueFrom(ctx, types.StringType, checksumStrMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate UUID for the data source ID
	id, err := common.GenerateID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate ID",
			fmt.Sprintf("Error generating ID for data source: %s", err.Error()),
		)
		return
	}

	// Set the computed values
	data.Checksums = checksums
	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
