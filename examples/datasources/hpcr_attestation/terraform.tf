terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = ">= 1.2.0"
    }
  }
}

# Decrypt an encrypted attestation record
data "hpcr_attestation" "attestation_encrypted" {
  attestation = file("./cert/se-checksums.txt.enc")
  privkey     = file("./cert/private.pem")
}

# Parse an unencrypted attestation record (no private key required)
data "hpcr_attestation" "attestation_unencrypted" {
  attestation = file("./cert/se-checksums.txt")
}

# Verify the attestation signature of a plaintext (already-decrypted) record.
# cert and signature must always be provided together or both omitted.
# - cert:      IBM attestation certificate (PEM) obtained from the runtime image
# - signature: filebase64() is required because se-signature.bin is binary
data "hpcr_attestation" "attestation_verified" {
  attestation = file("./cert/se-checksums.txt")
  cert        = file("./cert/attestation.crt")
  signature   = filebase64("./cert/se-signature.bin")
}

output "attestation_encrypted" {
  value = data.hpcr_attestation.attestation_encrypted.checksums
}

output "attestation_unencrypted" {
  value = data.hpcr_attestation.attestation_unencrypted.checksums
}

output "attestation_verified" {
  value = data.hpcr_attestation.attestation_verified.checksums
}
