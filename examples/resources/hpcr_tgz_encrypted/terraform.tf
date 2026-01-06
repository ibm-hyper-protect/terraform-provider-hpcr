terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

resource "hpcr_tgz_encrypted" "compose_b64_enc" {
  folder = "compose"
}

resource "hpcr_tgz_encrypted" "compose_b64_enc_platform" {
  folder = "compose"
  platform = "hpvs"
}

resource "hpcr_tgz_encrypted" "compose_b64_enc_cert" {
  folder = "compose"
  cert = file("./cert/encrypt.crt")
}

output "b64_enc_rendered" {
  value = hpcr_tgz_encrypted.compose_b64_enc.rendered
}

output "b64_enc_sha256" {
  value = hpcr_tgz_encrypted.compose_b64_enc.sha256_in
}

output "b64_enc_checksum" {
  value = hpcr_tgz_encrypted.compose_b64_enc.sha256_out
}

output "b64_enc_platform_rendered" {
  value = hpcr_tgz_encrypted.compose_b64_enc_platform.rendered
}

output "b64_enc_platform_sha256" {
  value = hpcr_tgz_encrypted.compose_b64_enc_platform.sha256_in
}

output "b64_enc_platform_checksum" {
  value = hpcr_tgz_encrypted.compose_b64_enc_platform.sha256_out
}

output "b64_enc_cert_rendered" {
  value = hpcr_tgz_encrypted.compose_b64_enc_cert.rendered
}

output "b64_enc_cert_sha256" {
  value = hpcr_tgz_encrypted.compose_b64_enc_cert.sha256_in
}

output "b64_enc_cert_checksum" {
  value = hpcr_tgz_encrypted.compose_b64_enc_cert.sha256_out
}
