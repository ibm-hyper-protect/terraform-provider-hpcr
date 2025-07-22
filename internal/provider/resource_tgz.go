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
	resp.TypeName = req.ProviderTypeName + "_tgz"
}

func (r *TgzResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a base64 encoded string from the TGZed files in the folder",
		Attributes: map[string]schema.Attribute{
			"folder": schema.StringAttribute{
				Description: "Path to folder",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"rendered": schema.StringAttribute{
				Description: "Generated encoded string",
				Computed:    true,
			},
			"sha256_in": schema.StringAttribute{
				Description: "SHA256 of input",
				Computed:    true,
			},
			"sha256_out": schema.SetAttribute{
				Description: "SHA256 of output",
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
			"Failed to generate ID",
			"Failed to generate UUID using Terraform inbuilt function",
		)

		return diags
	}

	encodedTgz, inputSha256, outputSha256, err := contract.HpcrTgz(folderPath)
	if err != nil {
		diags.AddError(
			"Failed to generate encoded TGZ",
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

func (r *TgzResource) handleGenerateTgz(ctx context.Context, data *TgzResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(r.generateTgz(ctx, data)...)
	return diags
}

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

func (r *TgzResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
