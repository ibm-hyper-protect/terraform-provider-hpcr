package decrypt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

const (
	encryptedChecksumPath      = "../../samples/attestation/se-checksums.txt.enc"
	privateKeyPath             = "../../samples/attestation/private.pem"
	sampleAttestationRecordKey = "baseimage"
)

// Testcase to check if DecryptPassword() is able to decrypt password
func TestDecryptPassword(t *testing.T) {
	encChecksum, err := gen.ReadDataFromFile(encryptedChecksumPath)
	if err != nil {
		t.Errorf("failed to read encrypted checksum - %v", err)
	}

	encodedEncryptedData := strings.Split(encChecksum, ".")[1]

	privateKeyData, err := gen.ReadDataFromFile(privateKeyPath)
	if err != nil {
		t.Errorf("failed to read private key - %v", err)
	}

	_, err = DecryptPassword(encodedEncryptedData, privateKeyData)
	if err != nil {
		t.Errorf("failed to decrypt password - %v", err)
	}
}

// Testcase to check if DecryptWorkload() is able to decrypt workload
func TestDecryptWorkload(t *testing.T) {
	encChecksum, err := gen.ReadDataFromFile(encryptedChecksumPath)
	if err != nil {
		t.Errorf("failed to read encrypted checksum - %v", err)
	}

	encodedEncryptedPassword := strings.Split(encChecksum, ".")[1]
	encodedEncryptedData := strings.Split(encChecksum, ".")[2]

	privateKeyData, err := gen.ReadDataFromFile(privateKeyPath)
	if err != nil {
		t.Errorf("failed to read private key - %v", err)
	}

	password, err := DecryptPassword(encodedEncryptedPassword, privateKeyData)
	if err != nil {
		t.Errorf("failed to decrypt password - %v", err)
	}

	result, err := DecryptWorkload(password, encodedEncryptedData)
	if err != nil {
		t.Errorf("failed to decrypt workload - %v", err)
	}

	assert.Contains(t, result, sampleAttestationRecordKey)
}
