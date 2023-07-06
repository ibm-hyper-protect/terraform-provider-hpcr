package example

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//go:embed certs.tf
var ConfigCerts string

func TestDatasourceCerts(t *testing.T) {
	t.Skip()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigCerts,
				Check: resource.ComposeTestCheckFunc(
					TestCheckOutput("cert_version", parseVersion),
				),
			},
		},
	})
}
