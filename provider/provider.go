package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/datasource"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"hpcr_tgz":            datasource.DataSourceTgz(),
			"hpcr_tgz_encrypted":  datasource.DataSourceTgzEncrypted(),
			"hpcr_text":           datasource.DataSourceText(),
			"hpcr_text_encrypted": datasource.DataSourceTextEncrypted(),
			"hpcr_json":           datasource.DataSourceJson(),
			"hpcr_json_encrypted": datasource.DataSourceJsonEncrypted(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"hpcr_tgz": datasource.ResourceTgz(),
		},
	}
}
