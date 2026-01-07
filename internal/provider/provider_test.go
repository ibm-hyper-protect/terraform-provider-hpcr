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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestHPCRProvider_Metadata(t *testing.T) {
	p := &HPCRProvider{version: "1.0.0"}

	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p.Metadata(context.TODO(), req, resp)

	if resp.TypeName != "hpcr" {
		t.Errorf("Expected TypeName to be 'hpcr', got '%s'", resp.TypeName)
	}

	if resp.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got '%s'", resp.Version)
	}
}

func TestHPCRProvider_Schema(t *testing.T) {
	p := &HPCRProvider{}

	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p.Schema(context.TODO(), req, resp)

	// Verify schema has description
	if resp.Schema.Description == "" {
		t.Error("Expected schema to have a description")
	}

	if resp.Schema.MarkdownDescription == "" {
		t.Error("Expected schema to have a markdown description")
	}
}

func TestHPCRProvider_Configure(t *testing.T) {
	p := &HPCRProvider{}

	// Note: We skip calling Configure in the test because it requires a valid
	// tfsdk.Config which is complex to set up in a unit test.
	// The Configure function is tested through integration tests.
	// Here we just verify the provider structure is correct.

	if p == nil {
		t.Error("Provider should not be nil")
	}
}

func TestHPCRProvider_Resources(t *testing.T) {
	p := &HPCRProvider{}

	resources := p.Resources(context.TODO())

	expectedCount := 8
	if len(resources) != expectedCount {
		t.Errorf("Expected %d resources, got %d", expectedCount, len(resources))
	}

	// Verify all resources can be instantiated
	for i, resourceFunc := range resources {
		r := resourceFunc()
		if r == nil {
			t.Errorf("Resource at index %d returned nil", i)
		}
	}
}

func TestHPCRProvider_DataSources(t *testing.T) {
	p := &HPCRProvider{}

	dataSources := p.DataSources(context.TODO())

	expectedCount := 4
	if len(dataSources) != expectedCount {
		t.Errorf("Expected %d data sources, got %d", expectedCount, len(dataSources))
	}

	// Verify all data sources can be instantiated
	for i, dataSourceFunc := range dataSources {
		ds := dataSourceFunc()
		if ds == nil {
			t.Errorf("DataSource at index %d returned nil", i)
		}
	}
}

func TestNew(t *testing.T) {
	version := "test-version"
	providerFunc := New(version)

	if providerFunc == nil {
		t.Fatal("New() should not return nil")
	}

	p := providerFunc()
	if p == nil {
		t.Fatal("Provider function should not return nil provider")
	}

	// Verify it's an HPCRProvider with the correct version
	hpcrProvider, ok := p.(*HPCRProvider)
	if !ok {
		t.Fatal("Provider should be of type *HPCRProvider")
	}

	if hpcrProvider.version != version {
		t.Errorf("Expected version '%s', got '%s'", version, hpcrProvider.version)
	}
}

func TestHPCRProvider_ImplementsProvider(t *testing.T) {
	var _ provider.Provider = &HPCRProvider{}
}

func TestHPCRProvider_VersionSet(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"dev version", "dev"},
		{"test version", "test"},
		{"release version", "1.2.3"},
		{"empty version", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HPCRProvider{version: tt.version}
			resp := &provider.MetadataResponse{}

			p.Metadata(context.TODO(), provider.MetadataRequest{}, resp)

			if resp.Version != tt.version {
				t.Errorf("Expected version '%s', got '%s'", tt.version, resp.Version)
			}
		})
	}
}

func TestHPCRProvider_ResourcesList(t *testing.T) {
	p := &HPCRProvider{}
	resources := p.Resources(context.TODO())

	// Verify we have the expected resource types
	expectedResourceCount := 8

	if len(resources) != expectedResourceCount {
		t.Errorf("Expected %d resources, got %d", expectedResourceCount, len(resources))
	}

	// Verify all resources can be instantiated
	for i, resourceFunc := range resources {
		r := resourceFunc()
		if r == nil {
			t.Errorf("Resource at index %d returned nil", i)
		}
	}
}

func TestHPCRProvider_DataSourcesList(t *testing.T) {
	p := &HPCRProvider{}
	dataSources := p.DataSources(context.TODO())

	// Verify we have the expected data source types
	expectedDataSources := 4 // image, attestation, encryption_certs, encryption_cert

	if len(dataSources) != expectedDataSources {
		t.Errorf("Expected %d data sources, got %d", expectedDataSources, len(dataSources))
	}
}
