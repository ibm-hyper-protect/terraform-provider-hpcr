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
	"os"
	"strings"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id, err := GenerateID()
	if err != nil {
		t.Fatalf("GenerateID() failed: %v", err)
	}

	if id == "" {
		t.Error("GenerateID() returned empty string")
	}

	// UUID should contain hyphens
	if !strings.Contains(id, "-") {
		t.Error("GenerateID() did not return a valid UUID format")
	}
}

func TestGeneratePrivateKey(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey() failed: %v", err)
	}

	if key == "" {
		t.Error("GeneratePrivateKey() returned empty string")
	}

	// Check for PEM format markers (could be RSA PRIVATE KEY or just PRIVATE KEY)
	if !strings.Contains(key, "BEGIN") || !strings.Contains(key, "PRIVATE KEY") {
		t.Error("GeneratePrivateKey() did not return a valid PEM formatted private key")
	}

	if !strings.Contains(key, "END") || !strings.Contains(key, "PRIVATE KEY") {
		t.Error("GeneratePrivateKey() did not return a complete PEM formatted private key")
	}
}

func TestGenerateID_Uniqueness(t *testing.T) {
	id1, err := GenerateID()
	if err != nil {
		t.Fatalf("GenerateID() failed: %v", err)
	}

	id2, err := GenerateID()
	if err != nil {
		t.Fatalf("GenerateID() failed: %v", err)
	}

	if id1 == id2 {
		t.Error("GenerateID() generated duplicate IDs")
	}
}

func TestReadFileData(t *testing.T) {
	// Create a temporary file with test content
	tmpFile := t.TempDir() + "/test.txt"
	testContent := "test content\nline 2"

	err := os.WriteFile(tmpFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Test reading existing file
	content, err := ReadFileData(tmpFile)
	if err != nil {
		t.Fatalf("ReadFileData() failed: %v", err)
	}

	if content != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, content)
	}
}

func TestReadFileData_NonExistentFile(t *testing.T) {
	_, err := ReadFileData("/path/to/nonexistent/file.txt")
	if err == nil {
		t.Error("ReadFileData() should return an error for non-existent file")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected error message to contain 'does not exist', got: %s", err.Error())
	}
}
