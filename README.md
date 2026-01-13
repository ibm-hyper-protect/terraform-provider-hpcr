# Terraform Provider HPCR

[![Terraform Registry](https://img.shields.io/badge/Terraform-Registry-623CE4?logo=terraform)](https://registry.terraform.io/providers/ibm-hyper-protect/hpcr/latest)
[![GitHub Actions](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/workflows/CI/badge.svg)](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ibm-hyper-protect/terraform-provider-hpcr)](https://goreportcard.com/report/github.com/ibm-hyper-protect/terraform-provider-hpcr)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Terraform provider for automating IBM Hyper Protect confidential computing workloads. It supports three platforms:
- **HPVS** (Hyper Protect Virtual Servers)
- **HPCR** (Hyper Protect Container Runtime for Red Hat Virtualization)
- **HPCC** (Hyper Protect Confidential Container for Red Hat OpenShift Peer Pods)

## Table of Contents

- [Features](#features)
- [About Contract-Go](#about-contract-go)
- [Installation](#installation)
  - [Requirements](#requirements)
- [Quick Start Examples](#quick-start-examples)
  - [Create TGZ Archive](#create-tgz-archive)
  - [Encrypt Contract Sections](#encrypt-contract-sections)
  - [Select HPCR Image](#select-hpcr-image)
- [Resources and Data Sources](#resources-and-data-sources)
  - [hpcr_tgz](#hpcr_tgz)
  - [hpcr_text_encrypted](#hpcr_text_encrypted)
  - [hpcr_image](#hpcr_image)
- [Documentation Resources](#documentation-resources)
- [Related Tools](#related-tools)
- [Support](#support)
- [How to Contribute](#how-to-contribute)
- [License](#license)
- [References](#references)

## Features

The Terraform Provider HPCR provides comprehensive support for deploying and managing IBM Hyper Protect confidential computing workloads through Infrastructure-as-Code. Key capabilities include:

**Archive Management**
- Create Base64-encoded tgz archives from docker-compose folders
- Automatic compression and encoding for contract workload sections
- Support for complex multi-container applications

**Encryption Operations**
- Encrypt contract sections (workload, environment) using HPVS encryption certificates
- Automatic retrieval of encryption certificates from IBM Cloud
- Support for latest or specific encryption certificate versions
- Secure handling of sensitive configuration data

**Image Selection & Validation**
- Retrieve the latest HPCR stock images from IBM Cloud VPC
- Semantic versioning support with version constraints (e.g., `>=1.1.0`, `~>1.0`)
- Automatic validation of image availability and checksums
- Support for public IBM Cloud images

**Contract Generation**
- Streamlined workflow for generating encrypted HPCR contracts
- Integration with [contract-go library](https://github.com/ibm-hyper-protect/contract-go) for contract operations
- Support for signed and encrypted contracts
- YAML-based contract composition

## About Contract-Go

This Terraform provider is built on top of the **[contract-go library](https://github.com/ibm-hyper-protect/contract-go)**, which provides the underlying functionality for IBM Hyper Protect contract operations.

The contract-go library is a Go library for automating IBM Hyper Protect confidential computing workloads across HPVS, HPCR, and HPCC platforms. It handles:

- **Contract Generation**: Create signed and encrypted contracts for secure enclaves
- **Certificate Operations**: Download and manage HPVS encryption certificates from IBM Cloud
- **Image Selection**: Retrieve and validate HPCR images with semantic versioning
- **Archive Management**: Generate Base64 tar archives from docker-compose and pods configurations
- **Attestation**: Decrypt attestation records from secure enclaves

### Why Use the Terraform Provider?

While the contract-go library can be used directly in Go applications, this Terraform provider offers:

- **Infrastructure-as-Code**: Manage HPCR deployments using Terraform workflows
- **State Management**: Track and manage contract resources with Terraform state
- **Integration**: Seamlessly integrate with other Terraform providers (IBM Cloud, AWS, Azure)
- **Declarative Syntax**: Define contracts using Terraform's HCL syntax
- **Automated Workflows**: Combine contract generation with infrastructure provisioning

For direct Go integration or command-line usage, consider:
- **[contract-go](https://github.com/ibm-hyper-protect/contract-go)**: Use directly in Go applications
- **[contract-cli](https://github.com/ibm-hyper-protect/contract-cli)**: Command-line tool for manual contract generation

### Contract-Go Documentation

- [User Documentation](https://ibm-hyper-protect.github.io/contract-go)
- [Go Package Docs](https://pkg.go.dev/github.com/ibm-hyper-protect/contract-go/v2)
- [Sample Configurations](https://github.com/ibm-hyper-protect/contract-go/tree/main/samples)

## Installation

This provider is available on the [Terraform Registry](https://registry.terraform.io/providers/ibm-hyper-protect/hpcr/latest).

Add the provider to your Terraform configuration:

```terraform
terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "~> 2.0"
    }
  }
}

provider "hpcr" {}
```

Then run:
```bash
terraform init
```

### Requirements

- [Terraform](https://www.terraform.io/downloads) 0.13 or later
- [OpenSSL](https://www.openssl.org/) binary (not LibreSSL)
  - **Linux**: `apt-get install openssl`
  - **macOS**: `brew install openssl`
  - **Windows**: Download from [Win32OpenSSL](https://slproweb.com/products/Win32OpenSSL.html)
- **Optional**: Set `OPENSSL_BIN` environment variable if OpenSSL isn't in system PATH

## Quick Start Examples

### Create TGZ Archive

Use the `hpcr_tgz` resource to create a tgz archive of your docker-compose folder:

```terraform
resource "hpcr_tgz" "compose" {
  folder = var.FOLDER
}
```

You can access the Base64-encoded content via the `rendered` property.

### Encrypt Contract Sections

Use the `hpcr_text_encrypted` resource to encrypt contract sections. By default, it uses the encryption key of the latest HPCR image:

```terraform
resource "hpcr_text_encrypted" "workload" {
  text = yamlencode({
    "compose" : {
      "archive" : resource.hpcr_tgz.compose.rendered
    }
  })
}
```

The typical use case is to encrypt the `workload` and `env` sections separately and pass the YAML-encoded contract as input.

### Select HPCR Image

Use the `hpcr_image` data source to find the matching HPCR stock image:

```terraform
data "ibm_is_images" "hyper_protect_images" {
  visibility = "public"
  status     = "available"
}

data "hpcr_image" "selected_image" {
  images = jsonencode(data.ibm_is_images.hyper_protect_images.images)
  spec   = ">=1.1.0"  # optional version constraint
}
```

## Resources and Data Sources

### hpcr_tgz

Creates a Base64-encoded tgz archive from a folder containing docker-compose configuration.

**Arguments:**
- `folder` (Required) - Path to the folder containing docker-compose files

**Attributes:**
- `rendered` - Base64-encoded tgz archive content

### hpcr_text_encrypted

Encrypts text content using HPVS encryption certificates.

**Arguments:**
- `text` (Required) - Plain text content to encrypt
- `cert` (Optional) - Specific encryption certificate to use (defaults to latest)

**Attributes:**
- `rendered` - Encrypted and Base64-encoded text

### hpcr_image

Data source to select the appropriate HPCR stock image from IBM Cloud VPC.

**Arguments:**
- `images` (Required) - JSON-encoded list of available VPC images
- `spec` (Optional) - [Semantic version constraint](https://github.com/Masterminds/semver#checking-version-constraints) (e.g., `>=1.1.0`, `~>1.0`)

**Attributes:**
- `image` - ID of the selected image
- `version` - Semantic version string of the selected image (e.g., `1.0.8`)

## Documentation Resources

- [Terraform Registry Documentation](https://registry.terraform.io/providers/ibm-hyper-protect/hpcr/latest/docs)
- [IBM Cloud Hyper Protect Virtual Server Documentation](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se)
- [Contract-Go Library Documentation](https://ibm-hyper-protect.github.io/contract-go)

## Related Tools

This provider is part of the IBM Hyper Protect ecosystem:

- **[contract-go](https://github.com/ibm-hyper-protect/contract-go)** - Go library for contract automation (underlying library for this provider)
- **[contract-cli](https://github.com/ibm-hyper-protect/contract-cli)** - Command-line contract generation tool
- **[k8s-operator-hpcr](https://github.com/ibm-hyper-protect/k8s-operator-hpcr)** - Kubernetes operator for HPCR workloads
- **[linuxone-vsi-automation-samples](https://github.com/ibm-hyper-protect/linuxone-vsi-automation-samples)** - Infrastructure-as-code examples and automation samples

## Support

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/issues)
- **Security**: Report vulnerabilities via [GitHub Security Advisories](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/security/advisories) (never public issues)
- **Discussions**: Community discussions available on the [repository discussions page](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/discussions)

## How to Contribute

This repository uses [semantic-release](https://github.com/semantic-release/semantic-release). Please author commit messages accordingly using conventional commits:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks

## License

[Apache 2.0](LICENSE)

## References

- [How to Publish a Terraform Provider](https://learn.hashicorp.com/tutorials/terraform/provider-release-publish?in=terraform/providers)
- [Terraform Provider Development](https://www.terraform.io/docs/extend/writing-custom-providers.html)
- [IBM Cloud Hyper Protect Services](https://www.ibm.com/cloud/hyper-protect-services)
