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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTgzResource_Metadata(t *testing.T) {
	r := NewTgzResource()

	req := resource.MetadataRequest{
		ProviderTypeName: "hpcr",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.TODO(), req, resp)

	if resp.TypeName != "hpcr_tgz" {
		t.Errorf("Expected TypeName to be 'hpcr_tgz', got '%s'", resp.TypeName)
	}
}

func TestTgzResource_Schema(t *testing.T) {
	r := NewTgzResource()

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.TODO(), req, resp)

	// Verify schema has required attributes
	if resp.Schema.Attributes == nil {
		t.Fatal("Schema attributes should not be nil")
	}

	requiredAttrs := []string{"id", "folder", "rendered", "sha256_in", "sha256_out"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected schema to have attribute '%s'", attr)
		}
	}

	// Verify folder is required
	folderAttr := resp.Schema.Attributes["folder"]
	if folderAttr.IsRequired() == false {
		t.Error("Expected 'folder' attribute to be required")
	}

	// Verify computed attributes
	computedAttrs := []string{"id", "rendered", "sha256_in", "sha256_out"}
	for _, attr := range computedAttrs {
		if resp.Schema.Attributes[attr].IsComputed() == false {
			t.Errorf("Expected '%s' attribute to be computed", attr)
		}
	}
}

func TestNewTgzResource(t *testing.T) {
	r := NewTgzResource()
	if r == nil {
		t.Fatal("NewTgzResource should not return nil")
	}

	// Verify it implements the Resource interface
	var _ resource.Resource = r
}

func TestTgzResource_SchemaDescriptions(t *testing.T) {
	r := NewTgzResource()

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.TODO(), req, resp)

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

func TestTgzResource_Delete(t *testing.T) {
	r := &TgzResource{}

	req := resource.DeleteRequest{}
	resp := &resource.DeleteResponse{}

	// Delete should be a no-op and not produce any errors
	r.Delete(context.TODO(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Error("Delete should not produce errors")
	}
}
