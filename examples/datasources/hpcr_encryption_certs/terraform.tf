terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

data "hpcr_encryption_certs" "encryption_cert" {
  versions = ["1.0.13", "1.0.14", "1.0.15"]
}

output "certs" {
  value = data.hpcr_encryption_certs.encryption_cert.certs
}

output "cert_15" {
  value = data.hpcr_encryption_certs.encryption_cert.certs["1.0.15"]["cert"]
}

output "cert_15_status" {
  value = data.hpcr_encryption_certs.encryption_cert.certs["1.0.15"]["status"]
}
