// Copyright 2022 IBM Corp.
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
package example

import (
	_ "embed"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//go:embed hello_world.tf
var ConfigHelloWorld string

func TestAccHelloWorld(t *testing.T) {

	folder, _ := filepath.Abs("../samples/hello-world")

	t.Setenv("TF_VAR_FOLDER", folder)
	t.Setenv("TF_VAR_LOGDNA_INGESTION_KEY", "00000000000000000000000")
	t.Setenv("TF_VAR_LOGDNA_INGESTION_HOSTNAME", "syslog-x.ibm.com")

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigHelloWorld,
				Check: resource.ComposeTestCheckFunc(
					TestCheckOutput("user_data", validateUserData),
				),
			},
		},
	})
}
