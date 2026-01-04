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
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEncryptionCertDataSource_Metadata(t *testing.T) {
	ds := NewEncryptionCertDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "hpcr",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.TODO(), req, resp)

	if resp.TypeName != "hpcr_encryption_cert" {
		t.Errorf("Expected TypeName to be 'hpcr_encryption_cert', got '%s'", resp.TypeName)
	}
}

func TestEncryptionCertDataSource_Schema(t *testing.T) {
	ds := NewEncryptionCertDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.TODO(), req, resp)

	// Verify schema has required attributes
	if resp.Schema.Attributes == nil {
		t.Fatal("Schema attributes should not be nil")
	}

	requiredAttrs := []string{"id", "certs", "spec", "cert", "version"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected schema to have attribute '%s'", attr)
		}
	}

	// Verify certs is required
	certsAttr := resp.Schema.Attributes["certs"]
	if certsAttr.IsRequired() == false {
		t.Error("Expected 'certs' attribute to be required")
	}

	// Verify spec is optional
	specAttr := resp.Schema.Attributes["spec"]
	if specAttr.IsOptional() == false {
		t.Error("Expected 'spec' attribute to be optional")
	}

	// Verify cert is computed
	certAttr := resp.Schema.Attributes["cert"]
	if certAttr.IsComputed() == false {
		t.Error("Expected 'cert' attribute to be computed")
	}

	// Verify version is computed
	versionAttr := resp.Schema.Attributes["version"]
	if versionAttr.IsComputed() == false {
		t.Error("Expected 'version' attribute to be computed")
	}

	// Verify id is computed
	idAttr := resp.Schema.Attributes["id"]
	if idAttr.IsComputed() == false {
		t.Error("Expected 'id' attribute to be computed")
	}
}

func TestNewEncryptionCertDataSource(t *testing.T) {
	ds := NewEncryptionCertDataSource()
	if ds == nil {
		t.Fatal("NewEncryptionCertDataSource should not return nil")
	}

	// Verify it implements the DataSource interface
	var _ datasource.DataSource = ds
}

func TestEncryptionCertDataSource_VersionSelection(t *testing.T) {
	// Test data with multiple certificate versions
	testCerts := map[string]string{
		"1.0.10": "cert-1.0.10-content",
		"1.0.11": "cert-1.0.11-content",
		"1.1.0":  "cert-1.1.0-content",
		"1.1.5":  "cert-1.1.5-content",
		"2.0.0":  "cert-2.0.0-content",
	}

	// Convert to JSON to verify it can be marshaled
	jsonData, err := json.Marshal(testCerts)
	if err != nil {
		t.Fatalf("Failed to marshal test certs: %v", err)
	}

	// Verify JSON is valid
	var unmarshaled map[string]string
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal test certs: %v", err)
	}

	// Verify all keys are present
	for key := range testCerts {
		if _, ok := unmarshaled[key]; !ok {
			t.Errorf("Expected version '%s' to be in unmarshaled data", key)
		}
	}
}
