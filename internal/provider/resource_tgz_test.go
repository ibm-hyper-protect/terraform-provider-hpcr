// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	testresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/internal/common"
)

const (
	hpcrTgzTestName          = "hpcr_tgz.test"
	sampleValidComposePath   = "../../samples/tgz"
	sampleInvalidComposePath = "../test"
	sampleBase64Tgz          = "H4sIAAAAAAAA/+zSwU7DMAwG4Jx5irzANrtN06QnXiV2HVotVVAymPb2aIDQDiAulRBav8uvKP8hsTxmPkrZcV6ec5X9JSxJrQ0AwBpzTew7uM13DaLCDnqLpkVABWgNgNKw+ku+8VJPoSiAGup0DqfjT73f7j//8pX/RJXyOrPU4UHrSVLKu3MuabwetZ6X8CSD/tiR/ZwPaaYSyuVw03ysU2g6O3gDbNFHGg266KkhJNtK043OWYyeGclYaSM7AnTMURz1LqInx2ziXw9is9ls7sxbAAAA//9BV11AAAgAAA=="

	sampleTgzTerraformConfig = `
provider "hpcr" {}

resource "hpcr_tgz" "test" {
	folder = "%s"
}
`
)

// Testcase to validate schema of hpcr_tgz
func TestTgzResourceSchema(t *testing.T) {
	r := HpcrTgzResource()

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("schema method returned errors: %v", resp.Diagnostics.Errors())
	}

	if resp.Schema.Attributes == nil {
		t.Fatal("schema attributes should not be nil")
	}

	attrs := []string{common.AttributeTgzFolderName, common.AttributeIdName, common.AttributeRenderedName, common.AttributeSha256InName, common.AttributeSha256OutName}
	for _, attr := range attrs {
		_, exists := resp.Schema.Attributes[attr]

		if !exists {
			t.Errorf("expected attribute %s not in schema", attr)
		}
	}
}

// Testcase to check generateTgz() is able to generate base64 of TGZed files
func TestTgzResourceGenerateTgz(t *testing.T) {
	r := &TgzResource{}
	ctx := context.Background()

	testCases := []struct {
		name         string
		folder       string
		rendered     string
		expectResult string
		expectErr    bool
	}{
		{
			name:         "Positive testcase",
			folder:       sampleValidComposePath,
			expectResult: sampleBase64Tgz,
		}, {
			name:      "Negative testcase",
			folder:    sampleInvalidComposePath,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := &TgzResourceModel{
				FolderPath: types.StringValue(tc.folder),
			}

			diags := r.generateTgz(ctx, data)

			if tc.expectErr {
				if !diags.HasError() {
					t.Errorf("expected error but got none")
				}

				return
			}

			if diags.HasError() {
				t.Errorf("unexpected error: %v", diags.Errors())
				return
			}

			if data.Rendered.ValueString() != tc.expectResult {
				t.Errorf("expected output %s, got %s", tc.expectResult, data.Rendered.ValueString())
			}
		})
	}
}

// Testcase to check if create, update and read are working as expected
func TestTgzResourceSuccess(t *testing.T) {
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			// Create and Read testing
			{
				Config: testTgzResourceConfig(sampleValidComposePath),
				Check: testresource.ComposeAggregateTestCheckFunc(
					testresource.TestCheckResourceAttr(hpcrTgzTestName, common.AttributeTgzFolderName, sampleValidComposePath),
					testresource.TestCheckResourceAttr(hpcrTgzTestName, common.AttributeRenderedName, sampleBase64Tgz),
					testresource.TestCheckResourceAttrSet(hpcrTgzTestName, common.AttributeIdName),
				),
			},
			// Update and Read testing
			{
				Config: testTgzResourceConfig(sampleValidComposePath),
				Check: testresource.ComposeAggregateTestCheckFunc(
					testresource.TestCheckResourceAttr(hpcrTgzTestName, common.AttributeTgzFolderName, sampleValidComposePath),
					testresource.TestCheckResourceAttr(hpcrTgzTestName, common.AttributeRenderedName, sampleBase64Tgz),
					testresource.TestCheckResourceAttrSet(hpcrTgzTestName, common.AttributeIdName),
				),
			},
			// ImportState testing
			{
				ResourceName:      hpcrTgzTestName,
				ImportState:       true,
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile("Import Not Implemented"),
			},
		},
	})
}

// common function to provide sample tf syntax
func testTgzResourceConfig(folderPath string) string {
	return fmt.Sprintf(sampleTgzTerraformConfig, folderPath)
}
