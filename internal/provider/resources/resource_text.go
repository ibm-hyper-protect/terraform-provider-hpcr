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

var _ resource.Resource = &TextResource{}

func NewTextResource() resource.Resource {
	return &TextResource{}
}

type TextResource struct{}

type TextResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Text      types.String `tfsdk:"text"`
	Rendered  types.String `tfsdk:"rendered"`
	Sha256In  types.String `tfsdk:"sha256_in"`
	Sha256Out types.String `tfsdk:"sha256_out"`
}

func (r *TextResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_text"
}

func (r *TextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a base64 encoded token from text input.",
		Description:         "Generates a base64 encoded token from text input.",

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

func (r *TextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TextResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the input text
	plainText := data.Text.ValueString()

	// Encode text using the contract-go library
	encoded, inputHash, outputHash, err := contract.HpcrText(plainText)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to encode text",
			fmt.Sprintf("Error encoding text: %s", err.Error()),
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
	data.Rendered = types.StringValue(encoded)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TextResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TextResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the input text
	plainText := data.Text.ValueString()

	// Encode text using the contract-go library
	encoded, inputHash, outputHash, err := contract.HpcrText(plainText)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to encode text",
			fmt.Sprintf("Error encoding text: %s", err.Error()),
		)
		return
	}

	// Set the computed fields (keep the existing ID)
	data.Rendered = types.StringValue(encoded)
	data.Sha256In = types.StringValue(inputHash)
	data.Sha256Out = types.StringValue(outputHash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op
}
