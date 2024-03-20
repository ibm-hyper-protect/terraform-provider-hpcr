package encrypt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
