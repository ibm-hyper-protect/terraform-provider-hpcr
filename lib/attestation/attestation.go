package attestation

import (
	"fmt"

	attest "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/decrypt"
	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

const (
	missingParameterErrStatement = "required parameter is missing"
)

// HpcrGetAttestationRecords - function to get attestation records from encrypted data
func HpcrGetAttestationRecords(data, privateKey string) (string, error) {
	if gen.CheckIfEmpty(data, privateKey) {
		return "", fmt.Errorf(missingParameterErrStatement)
	}
	encodedEncryptedPassword, encodedEncryptedData := gen.GetEncryptPassWorkload(data)

	password, err := attest.DecryptPassword(encodedEncryptedPassword, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password - %v", err)
	}

	attestationRecords, err := attest.DecryptWorkload(password, encodedEncryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt attestation records - %v", err)
	}

	return attestationRecords, nil
}
