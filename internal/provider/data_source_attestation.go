package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ibm-hyper-protect/contract-go/attestation"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/internal/common"
)

var _ datasource.DataSource = &AttestationDatasource{}

func HpcrAttestationDataSource() datasource.DataSource {
	return &AttestationDatasource{}
}

type AttestationDatasource struct{}

type AttestationDatasourceModel struct {
	Id          types.String `tfsdk:"id"`
	Attestation types.String `tfsdk:"attestation"`
	PrivKey     types.String `tfsdk:"privkey"`
	Checksums   types.String `tfsdk:"checksums"`
}

func (d *AttestationDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + common.DataSourceAttestationName
}

// Function to define schema of hpcr_attestation
func (d *AttestationDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: common.DataSourceAttestationDescription,
		Attributes: map[string]schema.Attribute{
			common.DataSourceAttestationInputName: schema.StringAttribute{
				Description: common.DataSourceAttestationInputDescription,
				Required:    true,
			},
			common.AttributePrivKeyName: schema.StringAttribute{
				Description: common.AttributePrivKeyDescription,
				Required:    true,
			},
			common.AttributeIdName: schema.StringAttribute{
				Description: common.AttributeIdDescription,
				Computed:    true,
			},
			common.DataSourceAttestationChecksumsName: schema.StringAttribute{
				Description: common.DataSourceAttestationChecksumsDescription,
				Computed:    true,
			},
		},
	}
}

// Function to create datasource - terraform apply
func (d *AttestationDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AttestationDatasourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	encryptedAttestation := data.Attestation.ValueString()
	privKey := data.PrivKey.ValueString()

	id, err := common.GenerateUuid()
	if err != nil {
		resp.Diagnostics.AddError(
			common.UuidGenerateFailureShortDescription,
			common.UUidGenerateFailureLongDescription,
		)

		return
	}

	decryptedAttestationRecords, err := attestation.HpcrGetAttestationRecords(encryptedAttestation, privKey)
	if err != nil {
		resp.Diagnostics.AddError(
			common.DataSourceAttestationFailureShortDescription,
			err.Error(),
		)

		return
	}

	data.Id = types.StringValue(id)
	data.Checksums = types.StringValue(decryptedAttestationRecords)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
