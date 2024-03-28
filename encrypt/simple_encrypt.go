package encrypt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// SimpleExecCommand - function to run os commands
func SimpleExecCommand(name string, stdinInput string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

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
		return "", err
	}
	return string(yamlBytes), nil
}

// KeyValueInjector - function to inject key value pair in YAML
func KeyValueInjector(contract map[string]interface{}, key, value string) (string, error) {
	contract[key] = value

	modifiedYAMLBytes, err := yaml.Marshal(contract)
	if err != nil {
		return "", err
	}

	return string(modifiedYAMLBytes), nil
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
func RandomPasswordGenerator() (string, error) {
	randomPassword, err := SimpleExecCommand("openssl", "", "rand", fmt.Sprint(keylen))
	if err != nil {
		return "", err
	}

	return randomPassword, nil
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

// CreateSigningCert - function to generate Signing Certificate
func CreateSigningCert(privateKey, cacert, cakey, csrData, csrPemData string, expiryDays int) (string, error) {
	var csr string
	if csrPemData == "" {
		privateKeyPath, err := CreateTempFile(privateKey)
		if err != nil {
			return "", err
		}

		var csrDataMap map[string]interface{}
		err = json.Unmarshal([]byte(csrData), &csrDataMap)
		if err != nil {
			return "", err
		}

		csrParam := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%sC/emailAddress=%s", csrDataMap["country"], csrDataMap["state"], csrDataMap["location"], csrDataMap["org"], csrDataMap["unit"], csrDataMap["domain"], csrDataMap["mail"])

		csr, err = SimpleExecCommand("openssl", "", "req", "-new", "-key", privateKeyPath, "-subj", csrParam)
		if err != nil {
			return "", err
		}

		err = os.Remove(privateKeyPath)
		if err != nil {
			return "", err
		}

	} else {
		csr = csrPemData
	}

	csrPath, err := CreateTempFile(csr)
	if err != nil {
		return "", err
	}

	caCertPath, err := CreateTempFile(cacert)
	if err != nil {
		return "", err
	}
	caKeyPath, err := CreateTempFile(cakey)
	if err != nil {
		return "", err
	}

	signingCert, err := CreateCert(csrPath, caCertPath, caKeyPath, expiryDays)
	if err != nil {
		return "", err
	}

	for _, path := range []string{csrPath, caCertPath, caKeyPath} {
		err := os.Remove(path)
		if err != nil {
			return "", err
		}
	}

	return EncodeToBase64(signingCert), nil
}

func CreateCert(csrPath, caCertPath, caKeyPath string, expiryDays int) (string, error) {
	signingCert, err := SimpleExecCommand("openssl", "", "x509", "-req", "-in", csrPath, "-CA", caCertPath, "-CAkey", caKeyPath, "-CAcreateserial", "-days", fmt.Sprintf("%d", expiryDays))
	if err != nil {
		return "", err
	}

	return signingCert, nil
}

// SignContract - function to sign encrypted contract
func SignContract(encryptedWorkload, encryptedEnv, privateKey string) (string, error) {
	combinedContract := encryptedWorkload + encryptedEnv

	privateKeyPath, err := CreateTempFile(privateKey)
	if err != nil {
		return "", err
	}

	workloadEnvSignature, err := SimpleExecCommand("openssl", combinedContract, "dgst", "-sha256", "-sign", privateKeyPath)
	if err != nil {
		return "", err
	}

	err = os.Remove(privateKeyPath)
	if err != nil {
		return "", err
	}

	return EncodeToBase64(workloadEnvSignature), nil
}

// GenFinalSignedContract - function to generate the final contract
func GenFinalSignedContract(workload, env, workloadEnvSig string) (string, error) {
	contract := map[string]interface{}{
		"workload":             workload,
		"env":                  env,
		"envWorkloadSignature": workloadEnvSig,
	}

	finalContract, err := MapToYaml(contract)
	if err != nil {
		return "", err
	}

	return finalContract, nil
}
