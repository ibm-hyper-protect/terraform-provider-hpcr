// Copyright 2026 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hashicorp/go-uuid"
	"gopkg.in/yaml.v3"
)

// GenerateID generates a random UUID to be used as a Terraform resource or data source ID.
func GenerateID() (string, error) {
	return uuid.GenerateUUID()
}

// GeneratePrivateKey generates a 4096-bit RSA private key using OpenSSL.
// It respects the OPENSSL_BIN environment variable for the OpenSSL binary path.
func GeneratePrivateKey() (string, error) {
	// Get OpenSSL binary path from environment variable or use default
	opensslBin := os.Getenv("OPENSSL_BIN")
	if opensslBin == "" {
		opensslBin = "openssl"
	}

	// Generate 4096-bit RSA private key
	cmd := exec.Command(opensslBin, "genrsa", "4096")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// ReadFileData reads the contents of a file and returns it as a string.
// Returns an error if the file does not exist or cannot be read.
func ReadFileData(filePath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	// Read file contents
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	return string(content), nil
}

func FilterChecksum(input string) map[string]string {
	// SHA-256: 64 hex chars
	re := regexp.MustCompile(`^([a-fA-F0-9]{64})\s+(.+)$`)

	scanner := bufio.NewScanner(strings.NewReader(input))
	checksumMap := make(map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		checksumMap[matches[1]] = matches[2]
	}

	return checksumMap
}

func RefineContract(yamlStr string) (string, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal input YAML: %v", err)
	}

	output := make(map[string]*yaml.Node)

	for key, val := range data {
		// Marshal the value of this top-level key back to YAML
		contentBytes, err := yaml.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("failed to marshal key %s: %v", key, err)
		}

		// Wrap it in a block literal style node (|)
		node := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Style: yaml.LiteralStyle,
			Value: string(contentBytes),
		}
		output[key] = node
	}

	finalNode := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{},
	}

	for key, val := range output {
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
		}
		finalNode.Content = append(finalNode.Content, keyNode, val)
	}

	// Step 4: Marshal final YAML
	resultBytes, err := yaml.Marshal(finalNode)
	if err != nil {
		return "", fmt.Errorf("failed to marshal final YAML: %v", err)
	}

	return string(resultBytes), nil
}
