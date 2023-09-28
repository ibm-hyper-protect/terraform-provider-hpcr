// Copyright 2023 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package datasource

package datasource

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/go-cty/cty"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	J "github.com/IBM/fp-go/json"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
)

type (
	Image struct {
		Architecture string `json:"architecture"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		OS           string `json:"os"`
		Status       string `json:"status"`
		Visibility   string `json:"visibility"`
		Checksum     string `json:"checksum"`
	}

	ImageVersion struct {
		ID       string
		Checksum string
		Version  *semver.Version
	}
)

var (
	schemaImagesIn = schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "List of images in JSON format",
	}
	schemaSpecIn = schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "*",
		Description:      "Semantic version range defining the HPCR image",
		ValidateDiagFunc: validateSpecFunc,
	}

	schemaIDOut = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID of the selected image",
	}

	schemaVersionOut = schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Version number of the selected image",
	}

	// parses a list of images from a string
	parseImages = F.Flow2(
		S.ToBytes,
		J.Unmarshal[[]Image],
	)

	// reHyperProtectOS tests if this is a hyper protect image
	reHyperProtectOS = regexp.MustCompile(`^hyper-protect-[\w-]+-s390x$`)

	// reHyperProtectVersion tests if the name references a valid hyper protect version
	reHyperProtectName = regexp.MustCompile(`^ibm-hyper-protect-container-runtime-(\d+)-(\d+)-s390x-(\d+)$`)
)

func validateSpecFunc(data any, path cty.Path) diag.Diagnostics {
	// test type
	specStr, ok := data.(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("unable to convert to string"))
	}
	// try to parse the spec
	_, err := semver.NewConstraint(specStr)
	if err != nil {
		return diag.FromErr(err)
	}
	// no problems
	return nil
}

// sortByVersion sorts the structs by version in inverse order, i.e. the first item will be the latest version
func sortByVersion(imgs []ImageVersion) []ImageVersion {
	sort.SliceStable(imgs, func(left, right int) bool {
		return imgs[left].Version.Compare(imgs[right].Version) > 0
	})
	return imgs
}

func imageVersionFomImage(img Image) ImageVersion {
	// assemble some infos
	parsed := reHyperProtectName.FindStringSubmatch(img.Name)
	version := semver.MustParse(fmt.Sprintf("%s.%s.%s", parsed[1], parsed[2], parsed[3]))
	// produce the result record
	return ImageVersion{ID: img.ID, Version: version, Checksum: img.Checksum}
}

// isCandidateImage tests if an image is a potential match for a hyper protect image
func isCandidateImage(img Image) bool {
	return img.Architecture == "s390x" &&
		img.Status == "available" &&
		img.Visibility == "public" &&
		reHyperProtectOS.MatchString(img.OS) &&
		reHyperProtectName.MatchString(img.Name)
}

func checkContraintPredicate(cstr *semver.Constraints) func(ImageVersion) bool {
	return func(img ImageVersion) bool {
		return cstr.Check(img.Version)
	}
}

// selectBySpec selects the latest version that matches the specification
func selectBySpec(spec string) func(img []ImageVersion) O.Option[ImageVersion] {
	cstr, err := semver.NewConstraint(spec)
	if err != nil {
		return F.Constant1[[]ImageVersion](O.None[ImageVersion]())
	}
	return F.Flow3(
		A.Filter(checkContraintPredicate(cstr)),
		I.Map(sortByVersion),
		A.Head[ImageVersion],
	)
}

func noMatchingVersionFound() error {
	return fmt.Errorf("unable to locate a matching version of the HPCR image")
}

func cannotConvertToString() error {
	return fmt.Errorf("unable to convert to string")
}

func selectImage(data *schema.ResourceData, ctx any) error {

	images, ok := data.GetOk(common.KeyImages)
	if !ok {
		return fmt.Errorf("input missing for [%s]", common.KeyImages)
	}
	spec, ok := data.GetOk(common.KeySpec)
	if !ok {
		return fmt.Errorf("input missing for [%s]", common.KeySpec)
	}

	return F.Pipe7(
		images,
		O.ToType[string],
		E.FromOption[string](cannotConvertToString),
		E.Chain(parseImages),
		E.Map[error](F.Flow2(
			A.Filter(isCandidateImage),
			A.Map(imageVersionFomImage),
		)),
		E.ChainOptionK[[]ImageVersion, ImageVersion](noMatchingVersionFound)(selectBySpec(spec.(string))),
		E.Map[error](func(version ImageVersion) ImageVersion {
			// update the data source
			data.SetId(time.Now().UTC().String())
			data.Set(common.KeyImageID, version.ID)
			data.Set(common.KeyVersion, version.Version.String())
			// some logging
			log.Printf("Selected image ID [%s] and version [%s]", version.ID, version.Version.String())

			return version
		}),
		E.ToError[ImageVersion],
	)
}

func DatasourceImage() *schema.Resource {
	return &schema.Resource{
		Read: selectImage,
		Schema: map[string]*schema.Schema{
			// input properties
			common.KeyImages: &schemaImagesIn,
			common.KeySpec:   &schemaSpecIn,
			// output properties
			common.KeyImageID: &schemaIDOut,
			common.KeyVersion: &schemaVersionOut,
			common.KeySha256:  &schemaSha256Out,
		},
		Description: "Selects an HPCR image from a JSON formatted list of images.",
	}
}
