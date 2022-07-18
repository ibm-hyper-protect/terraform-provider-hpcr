package example

import (
	_ "embed"
	"path/filepath"
	"regexp"
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
				Check:  resource.TestMatchOutput("user_data", regexp.MustCompile(`^.*$`)), // TestCheckOutput("user_data", validateUserData),
			},
		},
	})
}
