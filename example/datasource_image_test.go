package example

import (
	_ "embed"
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/Masterminds/semver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//go:embed images.json
var Images string

//go:embed images.tf
var ConfigImages string

var parseVersion = F.Flow2(
	E.Eitherize1(semver.NewVersion),
	E.Fold(O.Some[error], F.Constant1[*semver.Version](O.None[error]())),
)

func TestDatasourceImage(t *testing.T) {

	t.Setenv("TF_VAR_IMAGES", Images)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigImages,
				Check: resource.ComposeTestCheckFunc(
					TestCheckOutput("image_version", parseVersion),
				),
			},
		},
	})
}
