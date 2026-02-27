terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = ">= 1.2.0"
    }
  }
}

data "hpcr_encryption_certs" "encryption_cert" {
  versions = ["1.0.13", "1.0.14", "1.0.15"]
}

data "hpcr_encryption_cert" "cert" {
  certs = data.hpcr_encryption_certs.encryption_cert.certs
  spec  = "1.0.15"
}

output "cert" {
  value = data.hpcr_encryption_cert.cert.cert
}

output "expiry_days" {
  value = data.hpcr_encryption_cert.cert.expiry
}

output "expiry_status" {
  value = data.hpcr_encryption_cert.cert.status
}
