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
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var _ datasource.DataSource = &EncryptionCertDataSource{}

func NewEncryptionCertDataSource() datasource.DataSource {
	return &EncryptionCertDataSource{}
}

type EncryptionCertDataSource struct{}

type EncryptionCertDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Certs   types.Map    `tfsdk:"certs"`
	Spec    types.String `tfsdk:"spec"`
	Cert    types.String `tfsdk:"cert"`
	Version types.String `tfsdk:"version"`
}

func (d *EncryptionCertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encryption_cert"
}

func (d *EncryptionCertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Selects the best matching certificate from a map based on semantic versioning.",
		Description:         "Selects the best matching certificate based on the semantic version.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier",
			},
			"certs": schema.MapAttribute{
				MarkdownDescription: "Map of certificates from version to certificate content",
				Description:         "Map of certificates from version to certificate",
				ElementType:         types.StringType,
				Required:            true,
			},
			"spec": schema.StringAttribute{
				MarkdownDescription: "Semantic version range defining the HPCR certificate. Defaults to '*' (latest).",
				Description:         "Semantic version range defining the HPCR certificate",
				Optional:            true,
			},
			"cert": schema.StringAttribute{
				MarkdownDescription: "Selected certificate content",
				Description:         "Selected certificate",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version number of the selected certificate",
				Description:         "Version number of the selected certificate",
				Computed:            true,
			},
		},
	}
}

func (d *EncryptionCertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EncryptionCertDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract certificates map from Terraform config
	var certsMap map[string]string
	resp.Diagnostics.Append(data.Certs.ElementsAs(ctx, &certsMap, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(certsMap) == 0 {
		resp.Diagnostics.AddError(
			"Empty certificates map",
			"The certs map must contain at least one certificate",
		)
		return
	}

	// Get the spec value (default to "*" for latest if not provided)
	spec := "*"
	if !data.Spec.IsNull() && !data.Spec.IsUnknown() {
		spec = data.Spec.ValueString()
	}

	// Parse the semantic version constraint
	constraint, err := semver.NewConstraint(spec)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid version constraint",
			fmt.Sprintf("Failed to parse spec '%s': %s", spec, err.Error()),
		)
		return
	}

	// Parse all versions and filter by constraint
	var matchingVersions []*semver.Version
	for versionStr := range certsMap {
		version, err := semver.NewVersion(versionStr)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Invalid version format",
				fmt.Sprintf("Skipping version '%s': %s", versionStr, err.Error()),
			)
			continue
		}

		if constraint.Check(version) {
			matchingVersions = append(matchingVersions, version)
		}
	}

	if len(matchingVersions) == 0 {
		resp.Diagnostics.AddError(
			"No matching versions",
			fmt.Sprintf("No versions found matching constraint '%s'", spec),
		)
		return
	}

	// Sort versions in descending order and select the latest
	sort.Sort(sort.Reverse(semver.Collection(matchingVersions)))
	selectedVersion := matchingVersions[0].String()

	// Convert certs map to JSON for the library function
	certsJSON, err := json.Marshal(certsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to marshal certificates",
			fmt.Sprintf("Error converting certificates to JSON: %s", err.Error()),
		)
		return
	}

	// Get the certificate using the contract-go library
	version, cert, err := certificate.HpcrGetEncryptionCertificateFromJson(string(certsJSON), selectedVersion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get encryption certificate",
			fmt.Sprintf("Error retrieving certificate for version '%s': %s", selectedVersion, err.Error()),
		)
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
	data.Version = types.StringValue(version)
	data.Cert = types.StringValue(cert)
	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
