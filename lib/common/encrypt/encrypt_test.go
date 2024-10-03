package encrypt

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

const (
	certificateUrl       = "https://cloud.ibm.com/media/docs/downloads/hyper-protect-container-runtime/ibm-hyper-protect-container-runtime-1-0-s390x-15-encrypt.crt"
	simpleContractPath   = "../../samples/simple_contract.yaml"
	samplePrivateKeyPath = "../../samples/contract-expiry/private.pem"
	sampleCaCertPath     = "../../samples/contract-expiry/personal_ca.crt"
	sampleCaKeyPath      = "../../samples/contract-expiry/personal_ca.pem"
	sampleCsrFilePath    = "../../samples/contract-expiry/csr.pem"

	sampleCsrCountry  = "IN"
	sampleCsrState    = "Karnataka"
	sampleCsrLocation = "Bangalore"
	sampleCsrOrg      = "IBM"
	sampleCsrUnit     = "ISDL"
	sampleCsrDomain   = "HPVS"
	sampleCsrMailId   = "sashwat.k@ibm.com"

	sampleExpiryDays = 365

	simplePrivateKeyPath = "../../samples/encrypt/private.pem"
	simplePublicKeyPath  = "../../samples/encrypt/public.pem"
)

// Testcase to check if OpensslCheck() is able to check if openssl is present in the system or not
func TestOpensslCheck(t *testing.T) {
	err := OpensslCheck()
	if err != nil {
		t.Errorf("openssl check failed - %v", err)
	}
}

func TestGeneratePublicKey(t *testing.T) {
	privateKey, err := gen.ReadDataFromFile(simplePrivateKeyPath)
	if err != nil {
		t.Errorf("failed to read private key - %v", err)
	}

	publicKey, err := gen.ReadDataFromFile(simplePublicKeyPath)
	if err != nil {
		t.Errorf("failed to read public key - %v", err)
	}

	result, err := GeneratePublicKey(privateKey)
	if err != nil {
		t.Errorf("failed to generate public key - %v", err)
	}

	assert.Equal(t, result, publicKey)
}

// Testcase to check if RandomPasswordGenerator() is able to generate random password
func TestRandomPasswordGenerator(t *testing.T) {
	result, err := RandomPasswordGenerator()
	if err != nil {
		t.Errorf("failed to generate random password - %v", err)
	}

	assert.NotEmpty(t, result, "Random password did not get generated")
}

// Testcase to check if EncryptPassword() is able to encrypt password
func TestEncryptPassword(t *testing.T) {
	password, err := RandomPasswordGenerator()
	if err != nil {
		t.Errorf("failed to generate random password - %v", err)
	}

	encryptCertificate, err := gen.CertificateDownloader(certificateUrl)
	if err != nil {
		t.Errorf("failed to get encryption certificate - %v", err)
	}

	result, err := EncryptPassword(password, encryptCertificate)
	if err != nil {
		t.Errorf("failed to encrypt password - %v", err)
	}

	assert.NotEmpty(t, result, "Encrypted password did not get generated")
}

// Testcase to check if EncryptContract() is able to encrypt contract
func TestEncryptContract(t *testing.T) {
	var contractMap map[string]interface{}

	contract, err := gen.ReadDataFromFile(simpleContractPath)
	if err != nil {
		t.Errorf("failed to read contract - %v", err)
	}

	err = yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		t.Errorf("failed to unmarshal YAML - %v", err)
	}

	password, err := RandomPasswordGenerator()
	if err != nil {
		t.Errorf("failed to generate random password - %v", err)
	}

	result, err := EncryptContract(password, contractMap["workload"].(map[string]interface{}))
	if err != nil {
		t.Errorf("failed to encrypt contract - %v", err)
	}

	assert.NotEmpty(t, result, "Encrypted workload did not get generated")
}

// Testcase to check if EncryptString() is able to encrypt string
func TestEncryptString(t *testing.T) {
	password, err := RandomPasswordGenerator()
	if err != nil {
		t.Errorf("failed to generate random password - %v", err)
	}

	contract := `
	workload: |
		type: workload
	`

	result, err := EncryptString(password, contract)
	if err != nil {
		t.Errorf("failed to encrypt string - %v", err)
	}

	assert.NotEmpty(t, result, "Encrypted workload did not get generated")
}

// Testcase to check if EncryptFinalStr() is able to generate hyper-protect-basic.<password>.<workload>
func TestEncryptFinalStr(t *testing.T) {
	var contractMap map[string]interface{}

	contract, err := gen.ReadDataFromFile(simpleContractPath)
	if err != nil {
		t.Errorf("failed to get contract - %v", err)
	}

	err = yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		t.Errorf("failed to unmarshal YAML - %v", err)
	}

	password, err := RandomPasswordGenerator()
	if err != nil {
		t.Errorf("failed to generate random password - %v", err)
	}

	encryptCertificate, err := gen.CertificateDownloader(certificateUrl)
	if err != nil {
		t.Errorf("failed to get encryption certificate - %v", err)
	}

	encryptedRandomPassword, err := EncryptPassword(password, encryptCertificate)
	if err != nil {
		t.Errorf("failed to encrypt password - %v", err)
	}

	encryptedWorkload, err := EncryptContract(password, contractMap["workload"].(map[string]interface{}))
	if err != nil {
		t.Errorf("failed to encrypt workload - %v", err)
	}

	finalWorkload := EncryptFinalStr(encryptedRandomPassword, encryptedWorkload)

	assert.NotEmpty(t, finalWorkload, "Final workload did not get generated")
	assert.Contains(t, finalWorkload, "hyper-protect-basic.")
}

// Testcase to check if CreateSigningCert() is able to create signing certificate with CSR parameters
func TestCreateSigningCert(t *testing.T) {
	privateKey, err := gen.ReadDataFromFile(samplePrivateKeyPath)
	if err != nil {
		t.Errorf("failed to get private key - %v", err)
	}

	cacert, err := gen.ReadDataFromFile(sampleCaCertPath)
	if err != nil {
		t.Errorf("failed to get CA certificate - %v", err)
	}

	caKey, err := gen.ReadDataFromFile(sampleCaKeyPath)
	if err != nil {
		t.Errorf("failed to get CA Key - %v", err)
	}

	csrDataMap := map[string]interface{}{
		"country":  sampleCsrCountry,
		"state":    sampleCsrState,
		"location": sampleCsrLocation,
		"org":      sampleCsrOrg,
		"unit":     sampleCsrUnit,
		"domain":   sampleCsrDomain,
		"mail":     sampleCsrMailId,
	}
	csrDataStr, err := json.Marshal(csrDataMap)
	if err != nil {
		t.Errorf("failed to unmarshal JSON - %v", err)
	}

	signingCert, err := CreateSigningCert(privateKey, cacert, caKey, string(csrDataStr), "", sampleExpiryDays)
	if err != nil {
		t.Errorf("failed to create Signing certificate - %v", err)
	}

	assert.NotEmpty(t, signingCert, "Signing certificate did not get generated")
}

// Testcase to check if CreateSigningCert() is able to create signing certificate using CSR file
func TestCreateSigningCertCsrFile(t *testing.T) {
	privateKey, err := gen.ReadDataFromFile(samplePrivateKeyPath)
	if err != nil {
		t.Errorf("failed to get private key - %v", err)
	}

	cacert, err := gen.ReadDataFromFile(sampleCaCertPath)
	if err != nil {
		t.Errorf("failed to get CA certificate - %v", err)
	}

	caKey, err := gen.ReadDataFromFile(sampleCaKeyPath)
	if err != nil {
		t.Errorf("failed to get CA key - %v", err)
	}

	csr, err := gen.ReadDataFromFile(sampleCsrFilePath)
	if err != nil {
		t.Errorf("failed to get CSR file - %v", err)
	}

	signingCert, err := CreateSigningCert(privateKey, cacert, caKey, "", csr, sampleExpiryDays)
	if err != nil {
		t.Errorf("failed to create signing certificate - %v", err)
	}

	assert.NotEmpty(t, signingCert, "Signing certificate did not get generated")
}

// Testcase to check if SignContract() is able to sign the contract
func TestSignContract(t *testing.T) {
	var contractMap map[string]interface{}

	contract, err := gen.ReadDataFromFile(simpleContractPath)
	if err != nil {
		t.Errorf("failed to get contract - %v", err)
	}

	privateKey, err := gen.ReadDataFromFile(samplePrivateKeyPath)
	if err != nil {
		t.Errorf("failed to get private key - %v", err)
	}

	err = yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		t.Errorf("failed to unmarshal YAML - %v", err)
	}

	password, err := RandomPasswordGenerator()
	if err != nil {
		t.Errorf("failed to generate random password - %v", err)
	}

	encryptCertificate, err := gen.CertificateDownloader(certificateUrl)
	if err != nil {
		t.Errorf("failed to get encryption certificate - %v", err)
	}

	encryptedPassword, err := EncryptPassword(password, encryptCertificate)
	if err != nil {
		t.Errorf("failed to encrypt password - %v", err)
	}

	encryptedWorkload, err := EncryptContract(password, contractMap["workload"].(map[string]interface{}))
	if err != nil {
		t.Errorf("failed to encrypt workload - %v", err)
	}
	finalWorkload := EncryptFinalStr(encryptedPassword, encryptedWorkload)

	encryptedEnv, err := EncryptContract(password, contractMap["env"].(map[string]interface{}))
	if err != nil {
		t.Errorf("failed to encrypt env - %v", err)
	}

	finalEnv := EncryptFinalStr(encryptedPassword, encryptedEnv)

	workloadEnvSignature, err := SignContract(finalWorkload, finalEnv, privateKey)
	if err != nil {
		t.Errorf("failed to generate workload env signature - %v", err)
	}

	assert.NotEmpty(t, workloadEnvSignature, "workloadEnvSignature did not get generated")
}

// Testcase to check if GenFinalSignedContract() is able to generate signed contract
func TestGenFinalSignedContract(t *testing.T) {
	_, err := GenFinalSignedContract("test1", "test2", "test3")
	if err != nil {
		t.Errorf("failed to generate final signed and encrypted contract - %v", err)
	}
}
