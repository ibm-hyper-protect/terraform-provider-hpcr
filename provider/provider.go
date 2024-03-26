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
package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/datasource"
)

func Provider(version, commit string) func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"hpcr_tgz":                   datasource.ResourceTgz(),
				"hpcr_tgz_encrypted":         datasource.ResourceTgzEncrypted(),
				"hpcr_text":                  datasource.ResourceText(),
				"hpcr_text_encrypted":        datasource.ResourceTextEncrypted(),
				"hpcr_json":                  datasource.ResourceJSON(),
				"hpcr_json_encrypted":        datasource.ResourceJSONEncrypted(),
				"hpcr_contract_encrypted":    datasource.ResourceContractEncrypted(),
				"hpcr_contract_encrypted_ce": datasource.ResourceContractEncryptedSigningCert(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"hpcr_image":            datasource.DatasourceImage(),
				"hpcr_attestation":      datasource.DatasourceAttestation(),
				"hpcr_encryption_certs": datasource.DatasourceEncryptionCertificates(),
				"hpcr_encryption_cert":  datasource.DatasourceEncryptionCertificate(),
			},
			ConfigureContextFunc: datasource.ConfigureContext(version),
		}
	}
}
