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

package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var _ resource.Resource = &ContractEncryptedContractExpiryResource{}

func NewContractEncryptedContractExpiryResource() resource.Resource {
	return &ContractEncryptedContractExpiryResource{}
}

type ContractEncryptedContractExpiryResource struct{}

type ContractEncryptedContractExpiryResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Contract   types.String `tfsdk:"contract"`
	Platform   types.String `tfsdk:"platform"`
	Cert       types.String `tfsdk:"cert"`
	PrivKey    types.String `tfsdk:"privkey"`
	ExpiryDays types.Int64  `tfsdk:"expiry"`
	CaCert     types.String `tfsdk:"cacert"`
	CaKey      types.String `tfsdk:"cakey"`
	CsrParams  types.Map    `tfsdk:"csrparams"`
	CsrFile    types.String `tfsdk:"csrfile"`
	Rendered   types.String `tfsdk:"rendered"`
	Sha256In   types.String `tfsdk:"sha256_in"`
	Sha256Out  types.String `tfsdk:"sha256_out"`
}

func (r *ContractEncryptedContractExpiryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contract_encrypted_contract_expiry"
}

func (r *ContractEncryptedContractExpiryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates an encrypted and signed user data field with contract expiry enabled using a signing certificate.",
		Description:         "Generates an encrypted and signed user data field with contract expiry enabled.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contract": schema.StringAttribute{
				MarkdownDescription: "YAML serialization of the contract",
				Description:         "YAML serialization of the contract",
				Required:            true,
				Sensitive:           true,
			},
			"cert": schema.StringAttribute{
				MarkdownDescription: "Certificate used to encrypt the contract, in PEM format",
				Description:         "Certificate used to encrypt the contract, in PEM format",
				Required:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "Hyper Protect platform where this contract will be deployed. Defaults to hpvs",
				Description:         "Hyper Protect platform where this contract will be deployed",
				Optional:            true,
			},
			"privkey": schema.StringAttribute{
				MarkdownDescription: "Private key used to sign the contract. If omitted, a temporary signing key is created.",
				Description:         "Private key used to sign the contract",
				Optional:            true,
				Sensitive:           true,
			},
			"expiry": schema.Int64Attribute{
				MarkdownDescription: "Number of days for contract to expire",
				Description:         "Number of days for contract to expire",
				Required:            true,
			},
			"cacert": schema.StringAttribute{
				MarkdownDescription: "CA Certificate used to generate signing certificate",
				Description:         "CA Certificate used to generate signing certificate",
				Required:            true,
			},
			"cakey": schema.StringAttribute{
				MarkdownDescription: "CA Key used to generate signing certificate",
				Description:         "CA Key used to generate signing certificate",
				Required:            true,
			},
			"csrparams": schema.MapAttribute{
				MarkdownDescription: "CSR Parameters to generate signing certificate",
				Description:         "CSR Parameters to generate signing certificate",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"csrfile": schema.StringAttribute{
				MarkdownDescription: "CSR File to generate signing certificate",
				Description:         "CSR File to generate signing certificate",
				Optional:            true,
			},
			"rendered": schema.StringAttribute{
				MarkdownDescription: "Rendered output of the resource",
				Description:         "Rendered output of the resource",
				Computed:            true,
			},
			"sha256_in": schema.StringAttribute{
				MarkdownDescription: "SHA256 of the input",
				Description:         "SHA256 of the input",
				Computed:            true,
			},
			"sha256_out": schema.StringAttribute{
				MarkdownDescription: "SHA256 of the output",
				Description:         "SHA256 of the output",
				Computed:            true,
			},
		},
	}
}

func (r *ContractEncryptedContractExpiryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ContractEncryptedContractExpiryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get required and optional parameters
	contractYAML := data.Contract.ValueString()
	platform := data.Platform.ValueString()
	cert := data.Cert.ValueString()
	privKey := data.PrivKey.ValueString()
	expiryDays := int(data.ExpiryDays.ValueInt64())
	caCert := data.CaCert.ValueString()
	caKey := data.CaKey.ValueString()
	csrFilePath := data.CsrFile.ValueString()

	// Generate private key if not provided
	if privKey == "" {
		generatedKey, err := common.GeneratePrivateKey()
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to generate private key",
				fmt.Sprintf("Error generating private key: %s", err.Error()),
			)
			return
		}
		privKey = generatedKey
	}

	// Convert CSR params map to JSON string if provided
	var csrDataStr string
	if !data.CsrParams.IsNull() && !data.CsrParams.IsUnknown() {
		csrParamsMap := make(map[string]string)
		diag := data.CsrParams.ElementsAs(ctx, &csrParamsMap, false)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}

		csrDataBytes, err := json.Marshal(csrParamsMap)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to marshal CSR parameters",
				fmt.Sprintf("Error marshaling CSR parameters to JSON: %s", err.Error()),
			)
			return
		}
		csrDataStr = string(csrDataBytes)
	}

	// Read CSR file contents if provided
	var csrPemData string
	if csrFilePath != "" {
		contents, err := common.ReadFileData(csrFilePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read CSR file",
				fmt.Sprintf("Error reading CSR file: %s", err.Error()),
			)
			return
		}
		csrPemData = contents
	}

	// Validate that only one of csrparams or csrfile is provided
	if csrDataStr != "" && csrPemData != "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Only one of csrparams or csrfile can be provided, not both",
		)
		return
	}

	// Generate signed and encrypted contract with expiry using the contract-go library
	signedContract, inputHash, outputHash, err := contract.HpcrContractSignedEncryptedContractExpiry(contractYAML, platform, cert, privKey, caCert, caKey, csrDataStr, csrPemData, expiryDays)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create signed encrypted contract with expiry",
			fmt.Sprintf("Error creating signed encrypted contract with expiry: %s", err.Error()),
		)
		return
	}

	// Generate UUID for the resource ID
	id, err := common.GenerateID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate ID",
			fmt.Sprintf("Error generating ID for resource: %s", err.Error()),
		)
		return
	}

	// Set the computed fields
	data.ID = types.StringValue(id)
	data.Rendered = types.StringValue(signedContract)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractEncryptedContractExpiryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ContractEncryptedContractExpiryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Integrate your business logic here if needed
	// For most use cases, no action is needed in Read

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractEncryptedContractExpiryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ContractEncryptedContractExpiryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get required and optional parameters
	contractYAML := data.Contract.ValueString()
	platform := data.Platform.ValueString()
	cert := data.Cert.ValueString()
	privKey := data.PrivKey.ValueString()
	expiryDays := int(data.ExpiryDays.ValueInt64())
	caCert := data.CaCert.ValueString()
	caKey := data.CaKey.ValueString()
	csrFilePath := data.CsrFile.ValueString()

	// Generate private key if not provided
	if privKey == "" {
		generatedKey, err := common.GeneratePrivateKey()
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to generate private key",
				fmt.Sprintf("Error generating private key: %s", err.Error()),
			)
			return
		}
		privKey = generatedKey
	}

	// Convert CSR params map to JSON string if provided
	var csrDataStr string
	if !data.CsrParams.IsNull() && !data.CsrParams.IsUnknown() {
		csrParamsMap := make(map[string]string)
		diag := data.CsrParams.ElementsAs(ctx, &csrParamsMap, false)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}

		csrDataBytes, err := json.Marshal(csrParamsMap)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to marshal CSR parameters",
				fmt.Sprintf("Error marshaling CSR parameters to JSON: %s", err.Error()),
			)
			return
		}
		csrDataStr = string(csrDataBytes)
	}

	// Read CSR file contents if provided
	var csrPemData string
	if csrFilePath != "" {
		contents, err := common.ReadFileData(csrFilePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read CSR file",
				fmt.Sprintf("Error reading CSR file: %s", err.Error()),
			)
			return
		}
		csrPemData = contents
	}

	// Validate that only one of csrparams or csrfile is provided
	if csrDataStr != "" && csrPemData != "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Only one of csrparams or csrfile can be provided, not both",
		)
		return
	}

	// Generate signed and encrypted contract with expiry using the contract-go library
	signedContract, inputHash, outputHash, err := contract.HpcrContractSignedEncryptedContractExpiry(contractYAML, platform, cert, privKey, caCert, caKey, csrDataStr, csrPemData, expiryDays)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create signed encrypted contract with expiry",
			fmt.Sprintf("Error creating signed encrypted contract with expiry: %s", err.Error()),
		)
		return
	}

	// Set the computed fields (keep the existing ID)
	data.Rendered = types.StringValue(signedContract)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractEncryptedContractExpiryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op
}
