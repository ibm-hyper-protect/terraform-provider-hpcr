package encrypt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// SimpleExecCommand - function to run os commands
func SimpleExecCommand(name string, stdinInput string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	fmt.Println("CMD -", cmd)

	// Check for standard input
	if stdinInput != "" {
		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}
		defer stdinPipe.Close()

		go func() {
			defer stdinPipe.Close()
			stdinPipe.Write([]byte(stdinInput))
		}()
	}

	// Buffer to capture the output from the command.
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command.
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Return the output from the command and nil for the error.
	return out.String(), nil
}

// CreateTempFile - Function to create temp file
func CreateTempFile(data string) (string, error) {

	trimmedData := strings.TrimSpace(data)
	tmpFile, err := os.CreateTemp("", "hpvs-")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Write the data to the temp file.
	_, err = tmpFile.WriteString(trimmedData)
	if err != nil {
		return "", err
	}

	// Return the path to the temp file.
	return tmpFile.Name(), nil
}

// EncodeToBase64 - function to encode string as base64
func EncodeToBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// MapToYaml - function to convert string map to YAML
func MapToYaml(m map[string]interface{}) (string, error) {
	// Marshal the map into a YAML string.
	yamlBytes, err := yaml.Marshal(m)
	if err != nil {
		return "", err // Return the error to the caller.
	}
	return string(yamlBytes), nil
}

// OpensslCheck - function to check if openssl exists
func OpensslCheck() error {
	_, err := SimpleExecCommand("openssl", "", "version")

	if err != nil {
		return err
	}

	return nil
}

// RandomPasswordGenerator - function to generate random password
func RandomPasswordGenerator() (string, string, error) {
	randomPassword, err := SimpleExecCommand("openssl", "", "rand", fmt.Sprint(keylen))
	if err != nil {
		return "", "", err
	}

	return randomPassword, EncodeToBase64(randomPassword), nil
}

// EncryptPassword - function to encrypt password
func EncryptPassword(password, cert string) (string, error) {
	encryptCertPath, err := CreateTempFile(cert)
	if err != nil {
		return "", err
	}

	result, err := SimpleExecCommand("openssl", password, "rsautl", "-encrypt", "-inkey", encryptCertPath, "-certin")
	if err != nil {
		return "", err
	}

	err = os.Remove(encryptCertPath)
	if err != nil {
		return "", err
	}

	return EncodeToBase64(result), nil
}

// EncryptContract - function to encrypt contract
func EncryptContract(password string, section map[string]interface{}) (string, error) {
	contract, err := MapToYaml(section)
	if err != nil {
		return "", err
	}

	contractPath, err := CreateTempFile(contract)
	if err != nil {
		return "", err
	}

	result, err := SimpleExecCommand("openssl", password, "enc", "-aes-256-cbc", "-pbkdf2", "-pass", "stdin", "-in", contractPath)
	if err != nil {
		return "", err
	}

	err = os.Remove(contractPath)
	if err != nil {
		return "", err
	}

	return EncodeToBase64(result), nil
}

// EncryptFinalStr - function to get final encrypted section
func EncryptFinalStr(encryptedPassword, encryptedContract string) string {
	return fmt.Sprintf("hyper-protect-basic.%s.%s", encryptedPassword, encryptedContract)
}

// SignContract - function to sign encrypted contract
// func SignContract(encryptedWorkload, encryptedEnv, privateKey, csrData string) (string, error) {

// }
