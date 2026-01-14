terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = ">= 1.2.0"
    }
  }
}

resource "hpcr_tgz" "contract" {
  folder = "pods"
}

locals {
  # contract in clear text
  contract = yamlencode({
    "env" : {
      "type" : "env",
      "logging" : {
        "logRouter" : {
          "hostname" : "5c2d6b69-c7f0-41bd-b69b-240695369d6e.ingress.us-south.logs.cloud.ibm.com",
          "iamApiKey" : "ab00e3c09p1d4ff7fff9f04c12183413"
        }
      }
    },
    "workload" : {
      "type" : "workload",
      "play" : {
        "archive" : hpcr_tgz.contract.rendered
      }
    },
  })
}

resource "hpcr_contract_encrypted" "contract" {
  contract = local.contract
}

resource "hpcr_contract_encrypted" "contract_cert" {
  contract = local.contract
  cert     = file("./cert/encrypt.crt")
}

resource "hpcr_contract_encrypted" "contract_privkey" {
  contract = local.contract
  privkey  = file("./cert/private.pem")
}

resource "hpcr_contract_encrypted" "contract_platform" {
  contract = local.contract
  platform = "hpvs"
}

output "contract_rendered" {
  value = hpcr_contract_encrypted.contract.rendered
}

output "contract_sha256_in" {
  value = hpcr_contract_encrypted.contract.sha256_in
}

output "contract_sha256_out" {
  value = hpcr_contract_encrypted.contract.sha256_out
}

output "contract_cert_rendered" {
  value = hpcr_contract_encrypted.contract_cert.rendered
}

output "contract_cert_sha256_in" {
  value = hpcr_contract_encrypted.contract_cert.sha256_in
}

output "contract_cert_sha256_out" {
  value = hpcr_contract_encrypted.contract_cert.sha256_out
}

output "contract_privkey_rendered" {
  value = hpcr_contract_encrypted.contract_privkey.rendered
}

output "contract_privkey_sha256_in" {
  value = hpcr_contract_encrypted.contract_privkey.sha256_in
}

output "contract_privkey_sha256_out" {
  value = hpcr_contract_encrypted.contract_privkey.sha256_out
}

output "contract_platform_rendered" {
  value = hpcr_contract_encrypted.contract_platform.rendered
}

output "contract_platform_sha256_in" {
  value = hpcr_contract_encrypted.contract_platform.sha256_in
}

output "contract_platform_sha256_out" {
  value = hpcr_contract_encrypted.contract_platform.sha256_out
}
