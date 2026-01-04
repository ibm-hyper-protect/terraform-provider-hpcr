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

func TestEncryptionCertsDataSource_Metadata(t *testing.T) {
	ds := NewEncryptionCertsDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "hpcr",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.TODO(), req, resp)

	if resp.TypeName != "hpcr_encryption_certs" {
		t.Errorf("Expected TypeName to be 'hpcr_encryption_certs', got '%s'", resp.TypeName)
	}
}

func TestEncryptionCertsDataSource_Schema(t *testing.T) {
	ds := NewEncryptionCertsDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.TODO(), req, resp)

	// Verify schema has required attributes
	if resp.Schema.Attributes == nil {
		t.Fatal("Schema attributes should not be nil")
	}

	requiredAttrs := []string{"id", "template", "versions", "certs"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected schema to have attribute '%s'", attr)
		}
	}

	// Verify versions is required
	versionsAttr := resp.Schema.Attributes["versions"]
	if versionsAttr.IsRequired() == false {
		t.Error("Expected 'versions' attribute to be required")
	}

	// Verify template is optional
	templateAttr := resp.Schema.Attributes["template"]
	if templateAttr.IsOptional() == false {
		t.Error("Expected 'template' attribute to be optional")
	}

	// Verify certs is computed
	certsAttr := resp.Schema.Attributes["certs"]
	if certsAttr.IsComputed() == false {
		t.Error("Expected 'certs' attribute to be computed")
	}

	// Verify id is computed
	idAttr := resp.Schema.Attributes["id"]
	if idAttr.IsComputed() == false {
		t.Error("Expected 'id' attribute to be computed")
	}
}

func TestNewEncryptionCertsDataSource(t *testing.T) {
	ds := NewEncryptionCertsDataSource()
	if ds == nil {
		t.Fatal("NewEncryptionCertsDataSource should not return nil")
	}

	// Verify it implements the DataSource interface
	var _ datasource.DataSource = ds
}
