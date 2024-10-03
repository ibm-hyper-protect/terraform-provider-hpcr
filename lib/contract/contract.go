package contract

import (
	"fmt"

	"gopkg.in/yaml.v3"

	enc "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/encrypt"
	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

const (
	emptyParameterErrStatement = "required parameter is empty"
)

// HpcrText - function to generate base64 data and checksum from string
func HpcrText(plainText string) (string, string, error) {
	if gen.CheckIfEmpty(plainText) {
		return "", "", fmt.Errorf(emptyParameterErrStatement)
	}

	return gen.EncodeToBase64(plainText), gen.GenerateSha256(plainText), nil
}

// HpcrJson - function to generate base64 data and checksum from JSON string
func HpcrJson(plainJson string) (string, string, error) {
	if !gen.IsJSON(plainJson) {
		return "", "", fmt.Errorf("not a JSON data")
	}
	return gen.EncodeToBase64(plainJson), gen.GenerateSha256(plainJson), nil
}

// HpcrTextEncrypted - function to generate encrypted Hyper protect data and SHA256 from plain text
func HpcrTextEncrypted(plainText, encryptionCertificate string) (string, string, error) {
	if gen.CheckIfEmpty(plainText) {
		return "", "", fmt.Errorf(emptyParameterErrStatement)
	}

	return Encrypter(plainText, encryptionCertificate)
}

// HpcrJsonEncrypted - function to generate encrypted hyper protect data and SHA256 from plain JSON data
func HpcrJsonEncrypted(plainJson, encryptionCertificate string) (string, string, error) {
	if !gen.IsJSON(plainJson) {
		return "", "", fmt.Errorf("contract is not a JSON data")
	}
	return Encrypter(plainJson, encryptionCertificate)
}

// HpcrTgz - function to generate base64 of tar.tgz which was prepared from docker compose/podman files
func HpcrTgz(folderPath string) (string, error) {
	if gen.CheckIfEmpty(folderPath) {
		return "", fmt.Errorf(emptyParameterErrStatement)
	}

	if !gen.CheckFileFolderExists(folderPath) {
		return "", fmt.Errorf("folder doesn't exists - %s", folderPath)
	}

	filesFoldersList, err := gen.ListFoldersAndFiles(folderPath)
	if err != nil {
		return "", fmt.Errorf("failed to get files and folder under path - %v", err)
	}

	tgzBase64, err := gen.GenerateTgzBase64(filesFoldersList)
	if err != nil {
		return "", fmt.Errorf("failed to get base64 tgz - %v", err)
	}

	return tgzBase64, nil
}

// HpcrTgzEncrypted - function to generate encrypted tgz
func HpcrTgzEncrypted(folderPath, encryptionCertificate string) (string, string, error) {
	if gen.CheckIfEmpty(folderPath) {
		return "", "", fmt.Errorf(emptyParameterErrStatement)
	}

	tgzBase64, err := HpcrTgz(folderPath)
	if err != nil {
		return "", "", err
	}

	return Encrypter(tgzBase64, encryptionCertificate)
}

// HpcrContractSignedEncrypted - function to generate Signed and Encrypted contract
func HpcrContractSignedEncrypted(contract, encryptionCertificate, privateKey string) (string, error) {
	if gen.CheckIfEmpty(contract, privateKey, encryptionCertificate) {
		return "", fmt.Errorf(emptyParameterErrStatement)
	}

	publicKey, err := enc.GeneratePublicKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate public key - %v", err)
	}

	signedEncryptContract, err := EncryptWrapper(contract, encryptionCertificate, privateKey, publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign and encrypt contract - %v", err)
	}

	return signedEncryptContract, nil
}

// HpcrContractSignedEncryptedContractExpiry - function to generate sign with contract expiry enabled and encrypt contract (with CSR parameters and CSR file)
func HpcrContractSignedEncryptedContractExpiry(contract, encryptionCertificate, privateKey, cacert, caKey, csrDataStr, csrPemData string, expiryDays int) (string, error) {
	if gen.CheckIfEmpty(contract, privateKey, cacert, caKey) {
		return "", fmt.Errorf(emptyParameterErrStatement)
	}

	if csrPemData == "" && csrDataStr == "" || len(csrPemData) > 0 && len(csrDataStr) > 0 {
		return "", fmt.Errorf("the CSR parameters and CSR PEM file are parsed together or both are nil")
	}

	signingCert, err := enc.CreateSigningCert(privateKey, cacert, caKey, csrDataStr, csrPemData, expiryDays)
	if err != nil {
		return "", fmt.Errorf("failed to generate signing certificate - %v", err)
	}

	finalContract, err := EncryptWrapper(contract, encryptionCertificate, privateKey, signingCert)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed and encrypted contract - %v", err)
	}

	return finalContract, nil
}

// EncryptWrapper - wrapper function to sign (with and without contract expiry) and encrypt contract
func EncryptWrapper(contract, encryptionCertificate, privateKey, publicKey string) (string, error) {
	if gen.CheckIfEmpty(contract, privateKey, publicKey, encryptionCertificate) {
		return "", fmt.Errorf(emptyParameterErrStatement)
	}

	var contractMap map[string]interface{}

	err := yaml.Unmarshal([]byte(contract), &contractMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML - %v", err)
	}

	workloadData, err := gen.MapToYaml(contractMap["workload"].(map[string]interface{}))
	if err != nil {
		return "", fmt.Errorf("failed to convert MAP to YAML - %v", err)
	}

	encryptedWorkload, _, err := Encrypter(workloadData, encryptionCertificate)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt workload - %v", err)
	}

	updatedEnv, err := gen.KeyValueInjector(contractMap["env"].(map[string]interface{}), "signingKey", gen.EncodeToBase64(publicKey))
	if err != nil {
		return "", fmt.Errorf("failed to inject signingKey to env - %v", err)
	}

	encryptedEnv, _, err := Encrypter(updatedEnv, encryptionCertificate)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt env - %v", err)
	}

	workloadEnvSignature, err := enc.SignContract(encryptedWorkload, encryptedEnv, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign contract - %v", err)
	}

	finalContract, err := enc.GenFinalSignedContract(encryptedWorkload, encryptedEnv, workloadEnvSignature)
	if err != nil {
		return "", fmt.Errorf("failed to generate final contract - %v", err)
	}

	return finalContract, nil
}

// Encrypter - function to generate encrypted hyper protect data from plain string
func Encrypter(stringText, encryptionCertificate string) (string, string, error) {
	if gen.CheckIfEmpty(stringText, encryptionCertificate) {
		return "", "", fmt.Errorf(emptyParameterErrStatement)
	}

	password, err := enc.RandomPasswordGenerator()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random password - %v", err)
	}

	encodedEncryptedPassword, err := enc.EncryptPassword(password, encryptionCertificate)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt password - %v", err)
	}

	encryptedString, err := enc.EncryptString(password, stringText)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt key - %v", err)
	}

	return enc.EncryptFinalStr(encodedEncryptedPassword, encryptedString), gen.GenerateSha256(stringText), nil
}
