// Copyright 2025 IBM Corp.
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
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var _ datasource.DataSource = &EncryptionCertsDataSource{}

func NewEncryptionCertsDataSource() datasource.DataSource {
	return &EncryptionCertsDataSource{}
}

type EncryptionCertsDataSource struct{}

type EncryptionCertsDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Template types.String `tfsdk:"template"`
	Versions types.List   `tfsdk:"versions"`
	Certs    types.Map    `tfsdk:"certs"`
}

func (d *EncryptionCertsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encryption_certs"
}

func (d *EncryptionCertsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Downloads encryption certificates from the official IBM Cloud Object Storage location for specified HPCR versions.",
		Description:         "Downloads the encryption certificates for the given version numbers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier",
			},
			"template": schema.StringAttribute{
				MarkdownDescription: "Template used to download the encryption certificate. " +
					"May contain placeholders: {{.Major}}, {{.Minor}}, {{.Patch}}. " +
					"Default: https://hpvsvpcubuntu.s3.us.cloud-object-storage.appdomain.cloud/s390x-{{.Patch}}/ibm-hyper-protect-container-runtime-{{.Major}}-{{.Minor}}-s390x-{{.Patch}}-encrypt.crt",
				Description: "Template used to download the encryption certificate",
				Optional:    true,
			},
			"versions": schema.ListAttribute{
				MarkdownDescription: "List of strings, each denoting the version number of the certificate to download (e.g., ['1.0.10', '1.0.11'])",
				Description:         "List of strings, each denoting the version number of the certificate to download",
				ElementType:         types.StringType,
				Required:            true,
			},
			"certs": schema.MapAttribute{
				MarkdownDescription: "Map of certificates from version to certificate content",
				Description:         "Map of certificates from version to certificate",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *EncryptionCertsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EncryptionCertsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the template value (empty string will use default from library)
	template := ""
	if !data.Template.IsNull() && !data.Template.IsUnknown() {
		template = data.Template.ValueString()
	}

	// Extract versions list from Terraform config
	var versionList []string
	resp.Diagnostics.Append(data.Versions.ElementsAs(ctx, &versionList, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Download certificates using the contract-go library
	certsJSON, err := certificate.HpcrDownloadEncryptionCertificates(versionList, "json", template)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to download encryption certificates",
			fmt.Sprintf("Error downloading certificates: %s", err.Error()),
		)
		return
	}

	// Parse JSON response into a Go map
	var certsMap map[string]string
	if err := json.Unmarshal([]byte(certsJSON), &certsMap); err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse certificates JSON",
			fmt.Sprintf("Error parsing JSON response: %s", err.Error()),
		)
		return
	}

	// Convert Go map to types.Map
	certsTypeMap, diags := types.MapValueFrom(ctx, types.StringType, certsMap)
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

	// Set the computed fields
	data.Certs = certsTypeMap
	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
