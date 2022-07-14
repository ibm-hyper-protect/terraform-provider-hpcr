package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-provider-hpcr/datasource"
)

func provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			// "hpcr_tgz":            dataSourceTgz(),
			// "hpcr_tgz_encrypted":  dataSourceTgzEncrypted(),
			"hpcr_text":           datasource.DataSourceText(),
			"hpcr_text_encrypted": datasource.DataSourceTextEncrypted(),
			"hpcr_json":           datasource.DataSourceJson(),
			"hpcr_json_encrypted": datasource.DataSourceJsonEncrypted(),
		},
	}
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider,
	})
}
