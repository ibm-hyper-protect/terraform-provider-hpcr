package attestation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

const (
	encryptedChecksumPath      = "../samples/attestation/se-checksums.txt.enc"
	privateKeyPath             = "../samples/attestation/private.pem"
	sampleAttestationRecordKey = "baseimage"
)

// Testcase to check if HpcrGetAttestationRecords() retrieves attestation records from encrypted data
func TestHpcrGetAttestationRecords(t *testing.T) {
	encChecksum, err := gen.ReadDataFromFile(encryptedChecksumPath)
	if err != nil {
		t.Errorf("failed to get encrypted checksum - %v", err)
	}

	privateKeyData, err := gen.ReadDataFromFile(privateKeyPath)
	if err != nil {
		t.Errorf("failed to get private key - %v", err)
	}

	result, err := HpcrGetAttestationRecords(encChecksum, privateKeyData)
	if err != nil {
		t.Errorf("failed to decrypt attestation records - %v", err)
	}

	assert.Contains(t, result, sampleAttestationRecordKey)
}
