package example

import (
	_ "embed"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-provider-hpcr/validation"
)

//go:embed example1.tf
var ConfigExample1 string

//go:embed example2.tf
var ConfigExample2 string

func TestAccTgz(t *testing.T) {

	folder, _ := filepath.Abs("../samples/nginx-golang")

	t.Setenv("TF_VAR_FOLDER", folder)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigExample1,
				Check:  resource.TestMatchOutput("result", validation.Base64Re),
			},
		},
	})
}

func TestAccTgzEncrypted(t *testing.T) {

	folder, _ := filepath.Abs("../samples/nginx-golang")

	t.Setenv("TF_VAR_FOLDER", folder)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigExample2,
				Check:  resource.TestMatchOutput("result", validation.TokenRe),
			},
		},
	})
}
