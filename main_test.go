// Copyright (c) 2026 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/ibm-hyper-protect/terraform-provider-hpcr/internal/provider"
)

func TestVersion(t *testing.T) {
	// Test default version
	if version == "" {
		t.Error("Version should have a default value")
	}

	// Version should be "dev" by default (set in main.go)
	if version != "dev" {
		t.Logf("Version is set to: %s", version)
	}
}

func TestCommit(t *testing.T) {
	// Commit can be empty by default
	// This is set by goreleaser during release builds
	if commit != "" {
		t.Logf("Commit is set to: %s", commit)
	}
}

func TestProviderInitialization(t *testing.T) {
	// Test that we can create a provider using the New function
	testVersion := "test-version"
	providerFunc := provider.New(testVersion)

	if providerFunc == nil {
		t.Fatal("provider.New() should not return nil")
	}

	p := providerFunc()
	if p == nil {
		t.Fatal("Provider function should not return nil provider")
	}
}

func TestProviderFactory(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"dev version", "dev"},
		{"test version", "test"},
		{"release version", "1.0.0"},
		{"custom version", "2.5.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			providerFunc := provider.New(tt.version)
			if providerFunc == nil {
				t.Fatalf("provider.New(%s) should not return nil", tt.version)
			}

			p := providerFunc()
			if p == nil {
				t.Fatalf("Provider function for version %s should not return nil provider", tt.version)
			}
		})
	}
}
