terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

data "hpcr_attestation" "attestation_encrypted" {
  attestation = file("./cert/se-checksums.txt.enc")
  privkey     = file("./cert/private.pem")
}

data "hpcr_attestation" "attestation_unencrypted" {
  attestation = file("./cert/se-checksums.txt")
}

output "attestation_encrypted" {
  value = data.hpcr_attestation.attestation_encrypted.checksums
}

output "attestation_unencrypted" {
  value = data.hpcr_attestation.attestation_unencrypted.checksums
}
