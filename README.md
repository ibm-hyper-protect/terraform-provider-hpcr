# Terraform Provider HPCR

[![Terraform Registry](https://img.shields.io/badge/Terraform-Registry-623CE4?logo=terraform)](https://registry.terraform.io/providers/ibm-hyper-protect/hpcr/latest)
[![GitHub Actions](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/actions/workflows/build.yml/badge.svg)](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ibm-hyper-protect/terraform-provider-hpcr)](https://goreportcard.com/report/github.com/ibm-hyper-protect/terraform-provider-hpcr)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Terraform provider to automate generating workloads for IBM Hyper Protect confidential computing platforms.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Supported Platforms](#supported-platforms)
- [Examples](#examples)
- [Related Projects](#related-projects)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Overview

This Terraform provider is built on top of the **[contract-go library](https://github.com/ibm-hyper-protect/contract-go)**, which provides the core functionality for IBM Hyper Protect contract operations.

The contract-go library automates IBM Hyper Protect confidential computing workloads across HPVS, HPCR, and HPCC platforms, providing capabilities for:

- **Contract Generation**: Create signed and encrypted contracts for secure enclaves
- **Certificate Operations**: Download and manage HPVS encryption certificates from IBM Cloud
- **Image Selection**: Retrieve and validate HPCR images with semantic versioning
- **Archive Management**: Generate Base64 tar archives from docker-compose and pods configurations
- **Attestation**: Decrypt attestation records from secure enclaves

### What are Hyper Protect Services?

IBM Hyper Protect services provide confidential computing capabilities that protect data in use by leveraging the Secure Execution feature of IBM Z and LinuxONE.

**Learn more:**
- [Confidential computing with LinuxONE](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se)
- [IBM Hyper Protect Virtual Servers](https://www.ibm.com/docs/en/hpvs/2.2.x)
- [IBM Hyper Protect Confidential Container for Red Hat OpenShift](https://www.ibm.com/docs/en/hpcc/1.1.x)

### Why Use the Terraform Provider?

While the contract-go library can be used directly in Go applications, this Terraform provider offers:

- **Infrastructure-as-Code**: Manage HPCR deployments using Terraform workflows
- **State Management**: Track and manage contract resources with Terraform state
- **Integration**: Seamlessly integrate with other Terraform providers (IBM Cloud, libvirt, openstack)
- **Declarative Syntax**: Define contracts using Terraform's HCL syntax
- **Automated Workflows**: Combine contract generation with infrastructure provisioning

**For direct Go integration or command-line usage, consider:**
- **[contract-go](https://github.com/ibm-hyper-protect/contract-go)**: Use directly in Go applications
- **[contract-cli](https://github.com/ibm-hyper-protect/contract-cli)**: Command-line tool for manual contract generation

## Features

The Terraform Provider HPCR provides comprehensive support for deploying and managing IBM Hyper Protect confidential computing workloads through Infrastructure-as-Code. Key capabilities include:

**Archive Management**
- Create Base64-encoded tgz archives from folders containing docker-compose and podman play configurations
- Automatic compression and encoding for contract workload sections
- Support for complex multi-container applications

**Encryption Operations**
- Encrypt contract sections (workload, env) using HPVS encryption certificates
- Automatic retrieval of encryption certificates from IBM Cloud
- Support for latest or specific encryption certificate versions
- Secure handling of sensitive configuration data

**Image Selection & Validation**
- Select appropriate HPCR stock images from IBM Cloud VPC
- Semantic versioning support with flexible version constraints (e.g., `>=1.1.0`, `~>1.0`)
- Automatic validation of image availability
- Support for public IBM Cloud Hyper Protect images

**Contract Generation**
- Streamlined workflow for generating encrypted HPCR contracts
- Integration with [contract-go library](https://github.com/ibm-hyper-protect/contract-go) for contract operations
- Support for signed and encrypted contracts
- YAML-based contract composition

## Installation

This provider is published on the [Terraform Registry](https://registry.terraform.io/providers/ibm-hyper-protect/hpcr/latest) and can be installed automatically by Terraform.

Add the provider configuration to your Terraform files:

```terraform
terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "~> 0.16.2"
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

#### Optional: Custom OpenSSL Path

If OpenSSL is not in your system PATH, set the `OPENSSL_BIN` environment variable:

```bash
# Linux/macOS
export OPENSSL_BIN=/usr/bin/openssl

# Windows (PowerShell)
$env:OPENSSL_BIN="C:\Program Files\OpenSSL-Win64\bin\openssl.exe"
```

## Quick Start

### Create TGZ Archive

Use the `hpcr_tgz` resource to create a tgz archive from your docker-compose or podman play folder:

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

Use the `hpcr_image` data source to select the appropriate HPCR stock image from IBM Cloud VPC with optional version constraints:

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

## Documentation

### Provider Documentation

- [Terraform Registry Documentation](https://registry.terraform.io/providers/ibm-hyper-protect/hpcr/latest/docs)
- [Confidential computing with LinuxONE](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se)

### Contract-Go Documentation

- [User Documentation](https://ibm-hyper-protect.github.io/contract-go)
- [Go Package Docs](https://pkg.go.dev/github.com/ibm-hyper-protect/contract-go/v2)
- [Sample Configurations](https://github.com/ibm-hyper-protect/contract-go/tree/main/samples)

### Additional Resources

- [IBM Hyper Protect Virtual Servers](https://www.ibm.com/docs/en/hpvs/2.2.x)
- [IBM Hyper Protect Confidential Container for Red Hat OpenShift](https://www.ibm.com/docs/en/hpcc/1.1.x)
- [Terraform Provider Development Guide](https://www.terraform.io/docs/extend/writing-custom-providers.html)

## Supported Platforms

| Platform | Description | Support Status |
|----------|-------------|----------------|
| **HPVS** | [Hyper Protect Virtual Servers](https://www.ibm.com/docs/en/hpvs/2.2.x) - Confidential computing VMs on IBM Cloud | Supported |
| **HPCR-RHVS** | Hyper Protect Container Runtime for Red Hat Virtualization - Docker containers in secure enclaves | Supported |
| **HPCC-PeerPod** | [Hyper Protect Confidential Container Peer Pods](https://www.ibm.com/docs/en/hpcc/1.1.x) - Kubernetes confidential computing for OpenShift | Supported |

## Examples

### Resource: hpcr_tgz

Creates a Base64-encoded tgz archive from a folder containing docker-compose or podman play configuration.

**Arguments:**
- `folder` (Required) - Path to the folder containing docker-compose or podman play files

**Attributes:**
- `rendered` - Base64-encoded tgz archive content

**Example:**
```terraform
resource "hpcr_tgz" "compose" {
  folder = "./docker-compose"
}

output "archive" {
  value = resource.hpcr_tgz.compose.rendered
}
```

### Resource: hpcr_text_encrypted

Encrypts text content using HPVS encryption certificates.

**Arguments:**
- `text` (Required) - Plain text content to encrypt
- `cert` (Optional) - Specific encryption certificate to use (defaults to latest)

**Attributes:**
- `rendered` - Encrypted and Base64-encoded text

**Example:**
```terraform
resource "hpcr_text_encrypted" "workload" {
  text = yamlencode({
    "compose" : {
      "archive" : resource.hpcr_tgz.compose.rendered
    }
  })
}

resource "hpcr_text_encrypted" "env" {
  text = yamlencode({
    "type" : "env",
    "env" : {
      "MY_SECRET" : "sensitive_value"
    }
  })
}
```

### Data Source: hpcr_image

Selects the appropriate HPCR stock image from IBM Cloud VPC.

**Arguments:**
- `images` (Required) - JSON-encoded list of available VPC images
- `spec` (Optional) - [Semantic version constraint](https://github.com/Masterminds/semver#checking-version-constraints) (e.g., `>=1.1.0`, `~>1.0`)

**Attributes:**
- `image` - ID of the selected image
- `version` - Semantic version string of the selected image (e.g., `1.0.8`)

**Example:**
```terraform
data "ibm_is_images" "hyper_protect_images" {
  visibility = "public"
  status     = "available"
}

data "hpcr_image" "selected_image" {
  images = jsonencode(data.ibm_is_images.hyper_protect_images.images)
  spec   = ">=1.1.0"
}

output "selected_image_id" {
  value = data.hpcr_image.selected_image.image
}

output "selected_image_version" {
  value = data.hpcr_image.selected_image.version
}
```

## Related Projects

This provider is part of the IBM Hyper Protect ecosystem:

- **[contract-go](https://github.com/ibm-hyper-protect/contract-go)** - Go library for contract automation (underlying library for this provider)
- **[contract-cli](https://github.com/ibm-hyper-protect/contract-cli)** - Command-line contract generation tool
- **[k8s-operator-hpcr](https://github.com/ibm-hyper-protect/k8s-operator-hpcr)** - Kubernetes operator for HPCR workloads
- **[linuxone-vsi-automation-samples](https://github.com/ibm-hyper-protect/linuxone-vsi-automation-samples)** - Infrastructure-as-code examples and automation samples

## Contributing

Contributions are welcome! This repository uses [semantic-release](https://github.com/semantic-release/semantic-release) for automated versioning and releases.

Please follow the [Conventional Commits](https://www.conventionalcommits.org/) specification when authoring commit messages:

- `feat:` - New features (triggers minor version bump)
- `fix:` - Bug fixes (triggers patch version bump)
- `docs:` - Documentation changes only
- `chore:` - Maintenance tasks, refactoring, or dependency updates
- `BREAKING CHANGE:` - Include in commit footer for breaking changes (triggers major version bump)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/issues)
- **Security**: Report vulnerabilities via [GitHub Security Advisories](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/security/advisories) (never public issues)
- **Discussions**: Community discussions available on the [repository discussions page](https://github.com/ibm-hyper-protect/terraform-provider-hpcr/discussions)
