// Copyright 2023 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package datasource

package datasource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

var (
	schemaAttestationIn = schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The encrypted or unencrypted attestation record",
	}
)

func DatasourceAttestation() *schema.Resource {
	return &schema.Resource{
		Read: handleAttestation,
		Schema: map[string]*schema.Schema{
			common.KeyAttestation: &schemaAttestationIn,
		},
		Description: "handles the analysis of an attestation record.",
	}
}

func handleAttestation(data *schema.ResourceData, ctx any) error {

	attestation, ok := data.GetOk(common.KeyAttestation)
	if !ok {
		return fmt.Errorf("input missing for [%s]", common.KeyAttestation)
	}

	fmt.Printf("attestation: [%v]", attestation)

	return nil
}
