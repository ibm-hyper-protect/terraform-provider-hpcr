package certificate

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

const (
	defaultEncCertUrlTemplate    = "https://cloud.ibm.com/media/docs/downloads/hyper-protect-container-runtime/ibm-hyper-protect-container-runtime-{{.Major}}-{{.Minor}}-s390x-{{.Patch}}-encrypt.crt"
	missingParameterErrStatement = "required parameter is missing"
)

type CertSpec struct {
	Major string
	Minor string
	Patch string
}

// HpcrGetEncryptionCertificateFromJson - function to get encryption certificate from encryption certificate JSON data
func HpcrGetEncryptionCertificateFromJson(encryptionCertificateJson, version string) (string, string, error) {
	if gen.CheckIfEmpty(encryptionCertificateJson, version) {
		return "", "", fmt.Errorf(missingParameterErrStatement)
	}

	return gen.GetDataFromLatestVersion(encryptionCertificateJson, version)
}

// HpcrDownloadEncryptionCertificates - function to download encryption certificates for specified versions
func HpcrDownloadEncryptionCertificates(versionList []string) (string, error) {
	if gen.CheckIfEmpty(versionList) {
		return "", fmt.Errorf(missingParameterErrStatement)
	}

	var verCertMap = make(map[string]string)

	for _, version := range versionList {
		verSpec := strings.Split(version, ".")

		urlTemplate := template.New("url")
		urlTemplate, err := urlTemplate.Parse(defaultEncCertUrlTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to create url template - %v", err)
		}

		builder := &strings.Builder{}
		err = urlTemplate.Execute(builder, CertSpec{verSpec[0], verSpec[1], verSpec[2]})
		if err != nil {
			return "", fmt.Errorf("failed to apply template - %v", err)
		}

		url := builder.String()
		status, err := gen.CheckUrlExists(url)
		if err != nil {
			return "", fmt.Errorf("failed to check if URL exists - %v", err)
		}
		if !status {
			return "", fmt.Errorf("encryption certificate doesn't exist in %s", url)
		}

		cert, err := gen.CertificateDownloader(url)
		if err != nil {
			return "", fmt.Errorf("failed to download encryption certificate - %v", err)
		}

		verCertMap[version] = cert
	}

	jsonBytes, err := json.Marshal(verCertMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON - %v", err)
	}

	return string(jsonBytes), nil
}
