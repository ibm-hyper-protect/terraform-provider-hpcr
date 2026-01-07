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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var _ resource.Resource = &ContractEncryptedResource{}

func NewContractEncryptedResource() resource.Resource {
	return &ContractEncryptedResource{}
}

type ContractEncryptedResource struct{}

type ContractEncryptedResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Contract  types.String `tfsdk:"contract"`
	Platform  types.String `tfsdk:"platform"`
	Cert      types.String `tfsdk:"cert"`
	PrivKey   types.String `tfsdk:"privkey"`
	Rendered  types.String `tfsdk:"rendered"`
	Sha256In  types.String `tfsdk:"sha256_in"`
	Sha256Out types.String `tfsdk:"sha256_out"`
}

func (r *ContractEncryptedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contract_encrypted"
}

func (r *ContractEncryptedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates an encrypted and signed user data field from an HPCR contract.",
		Description:         "Generates an encrypted and signed user data field from an HPCR contract.",

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
			"platform": schema.StringAttribute{
				MarkdownDescription: "Hyper Protect platform where this contract will be deployed. Defaults to hpvs",
				Description:         "Hyper Protect platform where this contract will be deployed",
				Optional:            true,
			},
			"cert": schema.StringAttribute{
				MarkdownDescription: "Certificate used to encrypt the contract, in PEM format. Defaults to the latest HPCR image certificate if not specified.",
				Description:         "Certificate used to encrypt the contract, in PEM format",
				Optional:            true,
			},
			"privkey": schema.StringAttribute{
				MarkdownDescription: "Private key used to sign the contract. If omitted, a temporary signing key is created.",
				Description:         "Private key used to sign the contract",
				Optional:            true,
				Sensitive:           true,
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

func (r *ContractEncryptedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ContractEncryptedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get required and optional parameters
	contractYAML := data.Contract.ValueString()
	cert := data.Cert.ValueString()
	platform := data.Platform.ValueString()
	privKey := data.PrivKey.ValueString()

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

	refinedContract, err := common.RefineContract(contractYAML)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to refine contract",
			fmt.Sprintf("Error refining contract: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Contract YAML:- \n%s", refinedContract))

	// Generate signed and encrypted contract using the contract-go library
	signedContract, inputHash, outputHash, err := contract.HpcrContractSignedEncrypted(refinedContract, platform, cert, privKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create signed encrypted contract",
			fmt.Sprintf("Error creating signed encrypted contract: %s", err.Error()),
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

func (r *ContractEncryptedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ContractEncryptedResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Integrate your business logic here if needed
	// For most use cases, no action is needed in Read

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractEncryptedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ContractEncryptedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get required and optional parameters
	contractYAML := data.Contract.ValueString()
	cert := data.Cert.ValueString()
	platform := data.Platform.ValueString()
	privKey := data.PrivKey.ValueString()

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

	// Generate signed and encrypted contract using the contract-go library
	signedContract, inputHash, outputHash, err := contract.HpcrContractSignedEncrypted(contractYAML, platform, cert, privKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create signed encrypted contract",
			fmt.Sprintf("Error creating signed encrypted contract: %s", err.Error()),
		)
		return
	}

	// Set the computed fields (keep the existing ID)
	data.Rendered = types.StringValue(signedContract)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractEncryptedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op
}
