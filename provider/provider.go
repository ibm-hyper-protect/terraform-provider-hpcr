package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/datasource"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"hpcr_tgz":            datasource.ResourceTgz(),
			"hpcr_tgz_encrypted":  datasource.ResourceTgzEncrypted(),
			"hpcr_text":           datasource.ResourceText(),
			"hpcr_text_encrypted": datasource.ResourceTextEncrypted(),
			"hpcr_json":           datasource.ResourceJson(),
			"hpcr_json_encrypted": datasource.ResourceJsonEncrypted(),
		},
	}
}
