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

	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TgzEncryptedResource{}

func NewTgzEncryptedResource() resource.Resource {
	return &TgzEncryptedResource{}
}

// TgzEncryptedResource defines the resource implementation.
type TgzEncryptedResource struct{}

// TgzEncryptedResourceModel describes the resource data model.
type TgzEncryptedResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Folder    types.String `tfsdk:"folder"`
	Cert      types.String `tfsdk:"cert"`
	Platform  types.String `tfsdk:"platform"`
	Rendered  types.String `tfsdk:"rendered"`
	Sha256In  types.String `tfsdk:"sha256_in"`
	Sha256Out types.String `tfsdk:"sha256_out"`
}

func (r *TgzEncryptedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tgz_encrypted"
}

func (r *TgzEncryptedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates an encrypted token from the TGZed files in the folder.",
		Description:         "Generates an encrypted token from the TGZed files in the folder.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
				Description:         "Resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"folder": schema.StringAttribute{
				MarkdownDescription: "Path to the folder to encrypt",
				Description:         "Path to the folder to encrypt",
				Required:            true,
			},
			"cert": schema.StringAttribute{
				MarkdownDescription: "Certificate used to encrypt the JSON document, in PEM format. Defaults to the latest HPCR image certificate if not specified.",
				Description:         "Certificate used to encrypt the JSON document, in PEM format",
				Optional:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "Hyper Protect platform where this contract will be deployed. Defaults to hpvs",
				Description:         "Hyper Protect platform where this contract will be deployed",
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

func (r *TgzEncryptedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TgzEncryptedResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the folder path
	folderPath := data.Folder.ValueString()

	// Get optional parameters
	platform := data.Platform.ValueString()

	cert := ""
	if !data.Cert.IsNull() && !data.Cert.IsUnknown() {
		cert = data.Cert.ValueString()

		// check expiry of the encryption certificate
		expiryInfo, err := certificate.HpcrValidateEncryptionCertificate(cert)
		if err != nil {
			resp.Diagnostics.AddError(
				"Fail to encrypt text",
				fmt.Sprintf("Encryption certificate has expired: %s", err.Error()),
			)
			return
		}

		resp.Diagnostics.AddWarning(
			"Encryption certificate validity",
			expiryInfo,
		)
	}

	// Encrypt TGZ archive using the contract-go library
	encrypted, inputHash, outputHash, err := contract.HpcrTgzEncrypted(folderPath, platform, cert)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to encrypt TGZ archive",
			fmt.Sprintf("Error encrypting TGZ archive from folder '%s': %s", folderPath, err.Error()),
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
	data.Rendered = types.StringValue(encrypted)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TgzEncryptedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TgzEncryptedResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Integrate your business logic here if needed
	// For most use cases, no action is needed in Read

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TgzEncryptedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TgzEncryptedResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the folder path
	folderPath := data.Folder.ValueString()

	// Get optional parameters
	platform := data.Platform.ValueString()

	cert := ""
	if !data.Cert.IsNull() && !data.Cert.IsUnknown() {
		cert = data.Cert.ValueString()

		// check expiry of the encryption certificate
		expiryInfo, err := certificate.HpcrValidateEncryptionCertificate(cert)
		if err != nil {
			resp.Diagnostics.AddError(
				"Fail to encrypt text",
				fmt.Sprintf("Encryption certificate has expired: %s", err.Error()),
			)
			return
		}

		resp.Diagnostics.AddWarning(
			"Encryption certificate validity",
			expiryInfo,
		)
	}

	// Encrypt TGZ archive using the contract-go library
	encrypted, inputHash, outputHash, err := contract.HpcrTgzEncrypted(folderPath, platform, cert)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to encrypt TGZ archive",
			fmt.Sprintf("Error encrypting TGZ archive from folder '%s': %s", folderPath, err.Error()),
		)
		return
	}

	// Set the computed fields (keep the existing ID)
	data.Rendered = types.StringValue(encrypted)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TgzEncryptedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op for this resource type
}
