package example

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//go:embed cert.tf
var ConfigCert string

//go:embed certificates.json
var Certificates string

func TestDatasourceCert(t *testing.T) {

	t.Setenv("TF_VAR_CERTIFICATES", Certificates)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigCert,
				Check: resource.ComposeTestCheckFunc(
					TestCheckOutput("cert_version", parseVersion),
				),
			},
		},
	})
}
