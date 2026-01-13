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

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/internal/provider/datasources"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/internal/provider/resources"
)

// Ensure HPCRProvider satisfies various provider interfaces.
var _ provider.Provider = &HPCRProvider{}

// HPCRProvider defines the provider implementation.
type HPCRProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// HPCRProviderModel describes the provider data model.
type HPCRProviderModel struct {
	// Add provider-level configuration fields here if needed in the future
}

func (p *HPCRProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hpcr"
	resp.Version = p.version
}

func (p *HPCRProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for IBM Cloud Hyper Protect Virtual Server for VPC (HPCR). " +
			"This provider helps create encrypted contracts and user data for secure virtual servers.",
		MarkdownDescription: "Terraform provider for IBM Cloud Hyper Protect Virtual Server for VPC (HPCR). " +
			"This provider helps create encrypted contracts and user data for secure virtual servers.",
	}
}

func (p *HPCRProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config HPCRProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Initialize any shared client or configuration data here
	// For now, we'll pass nil as the provider doesn't have configuration
	// The datasource package context will be initialized per resource/datasource
	resp.DataSourceData = nil
	resp.ResourceData = nil
}

func (p *HPCRProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewTgzResource,
		resources.NewTgzEncryptedResource,
		resources.NewTextResource,
		resources.NewTextEncryptedResource,
		resources.NewJSONResource,
		resources.NewJSONEncryptedResource,
		resources.NewContractEncryptedResource,
		resources.NewContractEncryptedContractExpiryResource,
	}
}

func (p *HPCRProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewImageDataSource,
		datasources.NewAttestationDataSource,
		datasources.NewEncryptionCertsDataSource,
		datasources.NewEncryptionCertDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HPCRProvider{
			version: version,
		}
	}
}
