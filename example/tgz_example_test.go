package example

import (
	_ "embed"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-hpcr/provider"
)

var (
	providerName      = "hpcr"
	providerFactories = map[string]func() (*schema.Provider, error){
		providerName: func() (*schema.Provider, error) { return provider.Provider(), nil },
	}
)

//go:embed example1.tf
var ConfigExample1 string

func TestAccTgz(t *testing.T) {

	folder, _ := filepath.Abs("../samples/nginx-golang")

	t.Setenv("TF_VAR_FOLDER", folder)

	outputRe := regexp.MustCompile(`^((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigExample1,
				Check:  resource.TestMatchOutput("result", outputRe),
			},
		},
	})
}
