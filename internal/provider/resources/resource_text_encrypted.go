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

var _ resource.Resource = &TextEncryptedResource{}

func NewTextEncryptedResource() resource.Resource {
	return &TextEncryptedResource{}
}

type TextEncryptedResource struct{}

type TextEncryptedResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Text      types.String `tfsdk:"text"`
	Cert      types.String `tfsdk:"cert"`
	Platform  types.String `tfsdk:"platform"`
	Rendered  types.String `tfsdk:"rendered"`
	Sha256In  types.String `tfsdk:"sha256_in"`
	Sha256Out types.String `tfsdk:"sha256_out"`
}

func (r *TextEncryptedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_text_encrypted"
}

func (r *TextEncryptedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates an encrypted token from text input.",
		Description:         "Generates an encrypted token from text input.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"text": schema.StringAttribute{
				MarkdownDescription: "Text to archive",
				Description:         "Text to archive",
				Required:            true,
				Sensitive:           true,
			},
			"cert": schema.StringAttribute{
				MarkdownDescription: "Certificate used to encrypt the text, in PEM format. Defaults to the latest HPCR image certificate if not specified.",
				Description:         "Certificate used to encrypt the text, in PEM format",
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

func (r *TextEncryptedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TextEncryptedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the input text
	plainText := data.Text.ValueString()

	// Get the certificate (empty string will use default)
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

	platform := ""
	if !data.Platform.IsNull() && !data.Platform.IsUnknown() {
		platform = data.Platform.ValueString()
	}

	// Encrypt text using the contract-go library
	// Use empty string for hyperProtectOs to use default ("hpvs")
	encrypted, inputHash, outputHash, err := contract.HpcrTextEncrypted(plainText, platform, cert)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to encrypt text",
			fmt.Sprintf("Error encrypting text: %s", err.Error()),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TextEncryptedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TextEncryptedResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TextEncryptedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TextEncryptedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the input text
	plainText := data.Text.ValueString()

	// Get the certificate (empty string will use default)
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

	platform := ""
	if !data.Platform.IsNull() && !data.Platform.IsUnknown() {
		platform = data.Platform.ValueString()
	}

	// Encrypt text using the contract-go library
	// Use empty string for hyperProtectOs to use default ("hpvs")
	encrypted, inputHash, outputHash, err := contract.HpcrTextEncrypted(plainText, platform, cert)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to encrypt text",
			fmt.Sprintf("Error encrypting text: %s", err.Error()),
		)
		return
	}

	// Set the computed fields (keep the existing ID)
	data.Rendered = types.StringValue(encrypted)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TextEncryptedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op
}
