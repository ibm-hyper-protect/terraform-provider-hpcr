// Copyright (c) 2025 IBM Corp.
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

package common

const (
	TerraformProviderName = "hpcr"
	TerraformProviderLink = "registry.terraform.io/ibm-hyper-protect/hpcr"

	// common attributes
	AttributeIdName        = "id"
	AttributeIdDescription = "ID generated while executing resource"

	AttributeRenderedName = "rendered"

	AttributeSha256InName        = "sha256_in"
	AttributeSha256InDescription = "SHA256 of input"

	AttributeSha256OutName        = "sha256_out"
	AttributeSha256OutDescription = "SHA256 of output"

	// common error messages
	UuidGenerateFailureShortDescription = "Failed to generate ID"
	UUidGenerateFailureLongDescription  = "Failed to generate UUID using Terraform inbuilt function"

	// tgz resource
	ResourceTgzName                    = "_tgz"
	ResourceTgzDescription             = "Generates a base64 encoded string from the TGZed files in the folder"
	AttributeTgzFolderName             = "folder"
	AttributeTgzFolderDescription      = "Path to folder"
	AttributeTgzRenderedDescription    = "Encoded string of TGZed files"
	ResourceTgzFailureShortDescription = "Failed to generate encoded TGZ"
)
