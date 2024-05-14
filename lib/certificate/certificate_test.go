package certificate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sampleJsonData = `{
		"1.0.0": "data1",
		"1.2.5": "data2",
		"2.0.5": "data3",
		"3.5.10": "data4",
		"4.0.0": "data5"
	}`
	sampleEncryptionCertVersions = []string{"1.0.13", "1.0.14", "1.0.15"}
)

// Testcase to check if GetEncryptionCertificateFromJson() gets encryption certificate as per version constraint
func TestGetEncryptionCertificateFromJson(t *testing.T) {
	version := "> 1.0.0"

	key, value, err := HpcrGetEncryptionCertificateFromJson(sampleJsonData, version)
	if err != nil {
		t.Errorf("failed to get encryption certificate from JSON - %v", err)
	}

	assert.Equal(t, key, "4.0.0")
	assert.Equal(t, value, "data5")
}

// Testcase to check if DownloadEncryptionCertificates() is able to download encryption certificates as per constraint
func TestDownloadEncryptionCertificates(t *testing.T) {
	certs, err := HpcrDownloadEncryptionCertificates(sampleEncryptionCertVersions)
	if err != nil {
		t.Errorf("failed to download HPCR encryption certificates - %v", err)
	}

	assert.Contains(t, certs, "1.0.13")
}

// Testcase to check both DownloadEncryptionCertificates() and GetEncryptionCertificateFromJson() together
func TestCombined(t *testing.T) {
	certs, err := HpcrDownloadEncryptionCertificates(sampleEncryptionCertVersions)
	if err != nil {
		t.Errorf("failed to download HPCR encryption certificates - %v", err)
	}

	version := "> 1.0.14"

	key, _, err := HpcrGetEncryptionCertificateFromJson(certs, version)
	if err != nil {
		t.Errorf("failed to get encryption certificate from JSON - %v", err)
	}

	assert.Equal(t, key, "1.0.15")
}
