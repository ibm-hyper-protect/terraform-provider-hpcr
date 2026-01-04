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

package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ibm-hyper-protect/contract-go/v2/image"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var _ datasource.DataSource = &ImageDataSource{}

func NewImageDataSource() datasource.DataSource {
	return &ImageDataSource{}
}

type ImageDataSource struct{}

type ImageDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Images    types.String `tfsdk:"images"`
	Spec      types.String `tfsdk:"spec"`
	ImageID   types.String `tfsdk:"image"`
	ImageName types.String `tfsdk:"name"`
	Version   types.String `tfsdk:"version"`
	Sha256    types.String `tfsdk:"sha256"`
}

func (d *ImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (d *ImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Selects an HPCR image from a JSON formatted list of images based on semantic versioning.",
		Description:         "Selects an HPCR image from a JSON formatted list of images.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier",
			},
			"images": schema.StringAttribute{
				MarkdownDescription: "List of images in JSON format",
				Description:         "List of images in JSON format",
				Required:            true,
			},
			"spec": schema.StringAttribute{
				MarkdownDescription: "Semantic version range defining the HPCR image. Defaults to '*' (latest).",
				Description:         "Semantic version range defining the HPCR image",
				Optional:            true,
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "ID of the selected image",
				Description:         "ID of the selected image",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the selected image",
				Description:         "Name of the selected image",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version number of the selected image",
				Description:         "Version number of the selected image",
				Computed:            true,
			},
			"sha256": schema.StringAttribute{
				MarkdownDescription: "SHA256 checksum of the selected image",
				Description:         "SHA256 checksum of the selected image",
				Computed:            true,
			},
		},
	}
}

func (d *ImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ImageDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the images JSON data
	imageJsonData := data.Images.ValueString()
	if imageJsonData == "" {
		resp.Diagnostics.AddError(
			"Empty images data",
			"The images field must contain valid JSON data",
		)
		return
	}

	// Get the spec value (default to "*" for latest if not provided)
	spec := "*"
	if !data.Spec.IsNull() && !data.Spec.IsUnknown() {
		spec = data.Spec.ValueString()
	}

	// Select the best matching image using the contract-go library
	imageID, imageName, checksum, version, err := image.HpcrSelectImage(imageJsonData, spec)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to select image",
			fmt.Sprintf("Error selecting HPCR image with spec '%s': %s", spec, err.Error()),
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
	data.ImageID = types.StringValue(imageID)
	data.ImageName = types.StringValue(imageName)
	data.Version = types.StringValue(version)
	data.Sha256 = types.StringValue(checksum)
	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
