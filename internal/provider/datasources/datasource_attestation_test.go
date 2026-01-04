// Copyright 2025 IBM Corp.
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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestAttestationDataSource_Metadata(t *testing.T) {
	ds := NewAttestationDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "hpcr",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.TODO(), req, resp)

	if resp.TypeName != "hpcr_attestation" {
		t.Errorf("Expected TypeName to be 'hpcr_attestation', got '%s'", resp.TypeName)
	}
}

func TestAttestationDataSource_Schema(t *testing.T) {
	ds := NewAttestationDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.TODO(), req, resp)

	// Verify schema has required attributes
	if resp.Schema.Attributes == nil {
		t.Fatal("Schema attributes should not be nil")
	}

	requiredAttrs := []string{"id", "attestation", "privkey", "cert", "checksums"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected schema to have attribute '%s'", attr)
		}
	}

	// Verify attestation is required
	attestationAttr := resp.Schema.Attributes["attestation"]
	if attestationAttr.IsRequired() == false {
		t.Error("Expected 'attestation' attribute to be required")
	}

	// Verify privkey is optional
	privkeyAttr := resp.Schema.Attributes["privkey"]
	if privkeyAttr.IsOptional() == false {
		t.Error("Expected 'privkey' attribute to be optional")
	}

	// Verify cert is optional
	certAttr := resp.Schema.Attributes["cert"]
	if certAttr.IsOptional() == false {
		t.Error("Expected 'cert' attribute to be optional")
	}

	// Verify computed attributes
	computedAttrs := []string{"id", "checksums"}
	for _, attr := range computedAttrs {
		if resp.Schema.Attributes[attr].IsComputed() == false {
			t.Errorf("Expected '%s' attribute to be computed", attr)
		}
	}
}

func TestAttestationDataSource_Schema_Sensitive(t *testing.T) {
	ds := NewAttestationDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.TODO(), req, resp)

	// Verify privkey is marked as sensitive
	privkeyAttr := resp.Schema.Attributes["privkey"]
	if privkeyAttr.IsSensitive() == false {
		t.Error("Expected 'privkey' attribute to be marked as sensitive")
	}
}

func TestNewAttestationDataSource(t *testing.T) {
	ds := NewAttestationDataSource()
	if ds == nil {
		t.Fatal("NewAttestationDataSource should not return nil")
	}

	// Verify it implements the DataSource interface
	var _ datasource.DataSource = ds
}

func TestAttestationDataSource_SchemaDescriptions(t *testing.T) {
	ds := NewAttestationDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.TODO(), req, resp)

	// Verify schema has descriptions
	if resp.Schema.Description == "" {
		t.Error("Expected schema to have a description")
	}

	if resp.Schema.MarkdownDescription == "" {
		t.Error("Expected schema to have a markdown description")
	}

	// Verify attributes have descriptions (either Description or MarkdownDescription)
	for name, attr := range resp.Schema.Attributes {
		desc := attr.GetDescription()
		mdDesc := attr.GetMarkdownDescription()
		if desc == "" && mdDesc == "" {
			t.Errorf("Expected attribute '%s' to have a description or markdown description", name)
		}
	}
}
