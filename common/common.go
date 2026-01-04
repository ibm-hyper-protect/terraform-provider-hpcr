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
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-uuid"
)

// GenerateID generates a random UUID to be used as a Terraform resource or data source ID
func GenerateID() (string, error) {
	return uuid.GenerateUUID()
}

// GeneratePrivateKey generates a 4096-bit RSA private key using OpenSSL
// It respects the OPENSSL_BIN environment variable for the OpenSSL binary path
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
