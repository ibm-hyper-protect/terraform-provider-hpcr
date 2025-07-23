// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ibm-hyper-protect/contract-go/contract"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/internal/common"
)

var _ resource.Resource = &TgzResource{}

func HpcrTgzResource() resource.Resource {
	return &TgzResource{}
}

type TgzResource struct{}

type TgzResourceModel struct {
	Id         types.String `tfsdk:"id"`
	FolderPath types.String `tfsdk:"folder"`
	Rendered   types.String `tfsdk:"rendered"`
	Sha256In   types.String `tfsdk:"sha256_in"`
	Sha256Out  types.String `tfsdk:"sha256_out"`
}

func (r *TgzResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + common.ResourceTgzName
}

// Function to define schema of hpcr_tgz
func (r *TgzResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: common.ResourceTgzDescription,
		Attributes: map[string]schema.Attribute{
			common.AttributeTgzFolderName: schema.StringAttribute{
				Description: common.AttributeTgzFolderDescription,
				Required:    true,
			},
			common.AttributeIdName: schema.StringAttribute{
				Description: common.AttributeIdDescription,
				Computed:    true,
			},
			common.AttributeRenderedName: schema.StringAttribute{
				Description: common.AttributeTgzRenderedDescription,
				Computed:    true,
			},
			common.AttributeSha256InName: schema.StringAttribute{
				Description: common.AttributeSha256InDescription,
				Computed:    true,
			},
			common.AttributeSha256OutName: schema.StringAttribute{
				Description: common.AttributeSha256OutDescription,
				Computed:    true,
			},
		},
	}
}

// Function to generate TGZ
func (r *TgzResource) generateTgz(ctx context.Context, data *TgzResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	folderPath := data.FolderPath.ValueString()

	id, err := common.GenerateUuid()
	if err != nil {
		diags.AddError(
			common.UuidGenerateFailureShortDescription,
			common.UUidGenerateFailureLongDescription,
		)

		return diags
	}

	encodedTgz, inputSha256, outputSha256, err := contract.HpcrTgz(folderPath)
	if err != nil {
		diags.AddError(
			common.ResourceTgzFailureShortDescription,
			err.Error(),
		)

		return diags
	}

	data.Rendered = types.StringValue(encodedTgz)
	data.Id = types.StringValue(id)
	data.Sha256In = types.StringValue(inputSha256)
	data.Sha256Out = types.StringValue(outputSha256)

	return diags
}

// Handler function to create, update, read
func (r *TgzResource) handleGenerateTgz(ctx context.Context, data *TgzResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(r.generateTgz(ctx, data)...)
	return diags
}

// Function to create resource - terraform apply
func (r *TgzResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TgzResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.handleGenerateTgz(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Function to read resource - terraform plan/apply/destroy
func (r *TgzResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TgzResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.handleGenerateTgz(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Function to update resource - terraform apply
func (r *TgzResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TgzResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.handleGenerateTgz(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Function to delete resource - terraform destroy
func (r *TgzResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
