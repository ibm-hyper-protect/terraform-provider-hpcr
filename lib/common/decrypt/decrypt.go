package decrypt

import (
	"fmt"

	enc "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/encrypt"
	gen "github.com/ibm-hyper-protect/terraform-provider-hpcr/lib/common/general"
)

// DecryptPassword - function to decrypt encrypted string with private key
func DecryptPassword(base64EncryptedData, privateKey string) (string, error) {
	err := enc.OpensslCheck()
	if err != nil {
		return "", fmt.Errorf("openssl not found - %v", err)
	}

	decodedEncryptedData, err := gen.DecodeBase64String(base64EncryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode Base64 - %v", err)
	}

	encryptedDataPath, err := gen.CreateTempFile(decodedEncryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to generate temp file - %v", err)
	}

	privateKeyPath, err := gen.CreateTempFile(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file - %v", err)
	}

	result, err := gen.ExecCommand("openssl", "", "pkeyutl", "-decrypt", "-inkey", privateKeyPath, "-in", encryptedDataPath)
	if err != nil {
		return "", fmt.Errorf("failed to execute openssl command - %v", err)
	}

	for _, path := range []string{encryptedDataPath, privateKeyPath} {
		err := gen.RemoveTempFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to remove tmp file - %v", err)
		}
	}

	return result, nil
}

// DecryptWorkload - function to decrypt workload using password
func DecryptWorkload(password, encryptedWorkload string) (string, error) {
	err := enc.OpensslCheck()
	if err != nil {
		return "", fmt.Errorf("openssl not found - %v", err)
	}

	decodedEncryptedWorkload, err := gen.DecodeBase64String(encryptedWorkload)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data - %v", err)
	}

	encryptedDataPath, err := gen.CreateTempFile(decodedEncryptedWorkload)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file - %v", err)
	}

	result, err := gen.ExecCommand("openssl", password, "aes-256-cbc", "-d", "-pbkdf2", "-in", encryptedDataPath, "-pass", "stdin")
	if err != nil {
		return "", fmt.Errorf("failed to execute openssl command - %v", err)
	}

	err = gen.RemoveTempFile(encryptedDataPath)
	if err != nil {
		return "", fmt.Errorf("failed to remove temp file - %v", err)
	}

	return result, nil
}
