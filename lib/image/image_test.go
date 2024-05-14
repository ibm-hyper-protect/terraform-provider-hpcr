package image

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"

	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

const (
	ibmCloudImageListPath = "../samples/image.json"
	sampleVersion         = "1.0.8"

	sampleArchitecture = "s390x"
	sampleId           = "r042-45544dce-eff3-42dc-b149-6a33c2764e2d"
	sampleName         = "ibm-hyper-protect-container-runtime-1-0-s390x-8"
	sampleOs           = "hyper-protect-1-0-s390x-hpcr"
	sampleStatus       = "available"
	sampleVisibility   = "public"
	sampleChecksum     = "8c14f9676e727f21b31e6b0131d561b85b694cec050a7461d57e8fe8d94a70b8"
)

// Testcase to check SelectImage() is able to fetch the latest hyper protect image
func TestSelectImage(t *testing.T) {
	imageJsonList, err := gen.ReadDataFromFile(ibmCloudImageListPath)
	if err != nil {
		t.Errorf("failed to read data from file - %v", err)
	}

	imageId, imageName, imageChecksum, ImageVersion, err := HpcrSelectImage(imageJsonList, sampleVersion)
	if err != nil {
		t.Errorf("failed to select HPCR image - %v", err)
	}

	assert.Equal(t, imageId, sampleId)
	assert.Equal(t, imageName, sampleName)
	assert.Equal(t, imageChecksum, sampleChecksum)
	assert.Equal(t, ImageVersion, sampleVersion)
}

// Testcase to check if TestIsCandidateImage() can correctly identify if given data is hyper protect image data
func TestIsCandidateImage(t *testing.T) {
	sampleImageData := Image{
		Architecture: sampleArchitecture,
		ID:           sampleId,
		Name:         sampleName,
		OS:           sampleOs,
		Status:       sampleStatus,
		Visibility:   sampleVisibility,
		Checksum:     sampleChecksum,
	}

	result := IsCandidateImage(sampleImageData)

	assert.Equal(t, result, true)
}

// Testcase to check if PickLatestImage() is able to pick the latest image
func TestPickLatestImage(t *testing.T) {
	version, err := semver.NewVersion(sampleVersion)
	if err != nil {
		t.Errorf("failed to generate semantic version - %v", err)
	}

	var image []ImageVersion

	image = append(image, ImageVersion{ID: sampleId, Name: sampleName, Checksum: sampleChecksum, Version: version})

	imageId, imageName, imageChecksum, imageVersion, err := PickLatestImage(image, sampleVersion)
	if err != nil {
		t.Errorf("failed to pick latest image - %v", err)
	}

	assert.Equal(t, imageId, sampleId)
	assert.Equal(t, imageName, sampleName)
	assert.Equal(t, imageChecksum, sampleChecksum)
	assert.Equal(t, imageVersion, sampleVersion)
}
