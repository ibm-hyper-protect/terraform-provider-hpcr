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

	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TgzResource{}

func NewTgzResource() resource.Resource {
	return &TgzResource{}
}

// TgzResource defines the resource implementation.
type TgzResource struct{}

// TgzResourceModel describes the resource data model.
type TgzResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Folder    types.String `tfsdk:"folder"`
	Rendered  types.String `tfsdk:"rendered"`
	Sha256In  types.String `tfsdk:"sha256_in"`
	Sha256Out types.String `tfsdk:"sha256_out"`
}

func (r *TgzResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tgz"
}

func (r *TgzResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a base64 encoded string from the TGZed files in the folder.",
		Description:         "Generates a base64 encoded string from the TGZed files in the folder.",

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
				MarkdownDescription: "Path to the folder to archive",
				Description:         "Path to the folder to archive",
				Required:            true,
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

func (r *TgzResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TgzResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the folder path
	folderPath := data.Folder.ValueString()

	// Create TGZ archive using the contract-go library
	tgzBase64, inputHash, outputHash, err := contract.HpcrTgz(folderPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create TGZ archive",
			fmt.Sprintf("Error creating TGZ archive from folder '%s': %s", folderPath, err.Error()),
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
	data.Rendered = types.StringValue(tgzBase64)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TgzResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TgzResourceModel

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

func (r *TgzResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TgzResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the folder path
	folderPath := data.Folder.ValueString()

	// Create TGZ archive using the contract-go library
	tgzBase64, inputHash, outputHash, err := contract.HpcrTgz(folderPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create TGZ archive",
			fmt.Sprintf("Error creating TGZ archive from folder '%s': %s", folderPath, err.Error()),
		)
		return
	}

	// Set the computed fields (keep the existing ID)
	data.Rendered = types.StringValue(tgzBase64)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TgzResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op for this resource type
}
