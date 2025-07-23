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
