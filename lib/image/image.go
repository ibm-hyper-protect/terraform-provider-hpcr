package image

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/Masterminds/semver/v3"

	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
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
		Name     string
		Version  *semver.Version
	}
)

var (
	// reHyperProtectOS tests if this is a hyper protect image
	reHyperProtectOS = regexp.MustCompile(`^hyper-protect-[\w-]+-s390x-hpcr$`)

	// reHyperProtectVersion tests if the name references a valid hyper protect version
	reHyperProtectName = regexp.MustCompile(`^ibm-hyper-protect-container-runtime-(\d+)-(\d+)-s390x-(\d+)$`)
)

const (
	emptyParameterErrStatement = "required parameter is empty"
)

// HpcrSelectImage - function to return the latest HPVS image
func HpcrSelectImage(imageJsonData, versionSpec string) (string, string, string, string, error) {
	if gen.CheckIfEmpty(imageJsonData, versionSpec) {
		return "", "", "", "", fmt.Errorf(emptyParameterErrStatement)
	}

	var images []Image
	var hyperProtectImages []ImageVersion

	err := json.Unmarshal([]byte(imageJsonData), &images)
	if err != nil {
		fmt.Println("Error:", err)
		return "", "", "", "", fmt.Errorf("failed to unmarshal JSON - %v", err)
	}

	for _, image := range images {
		if IsCandidateImage(image) {
			versionRegex := reHyperProtectName.FindStringSubmatch(image.Name)
			hyperProtectImages = append(hyperProtectImages, ImageVersion{
				ID:       image.ID,
				Name:     image.Name,
				Checksum: image.Checksum,
				Version:  semver.MustParse(fmt.Sprintf("%s.%s.%s", versionRegex[1], versionRegex[2], versionRegex[3])),
			})
		}
	}

	return PickLatestImage(hyperProtectImages, versionSpec)
}

// IsCandidateImage - function to check if image JSON data belong to Hyper Protect Image
func IsCandidateImage(img Image) bool {
	return img.Architecture == "s390x" && img.Status == "available" && img.Visibility == "public" && reHyperProtectOS.MatchString(img.OS) && reHyperProtectName.MatchString(img.Name)
}

// PickLatestImage - function to pick the latest Hyper Protect Image
func PickLatestImage(hyperProtectImages []ImageVersion, version string) (string, string, string, string, error) {
	if gen.CheckIfEmpty(hyperProtectImages, version) {
		return "", "", "", "", fmt.Errorf(emptyParameterErrStatement)
	}

	targetConstraint, err := semver.NewConstraint(version)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error parsing target version constraint - %v", err)
	}

	var matchingVersions []*semver.Version

	for _, image := range hyperProtectImages {
		if targetConstraint.Check(image.Version) {
			matchingVersions = append(matchingVersions, image.Version)
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(matchingVersions)))

	if len(matchingVersions) > 0 {
		for _, image := range hyperProtectImages {
			if image.Version.Equal(matchingVersions[0]) {
				return image.ID, image.Name, image.Checksum, image.Version.String(), nil
			}
		}
	}

	return "", "", "", "", fmt.Errorf("no Hyper Protect image matching version found for the given constraint")
}
