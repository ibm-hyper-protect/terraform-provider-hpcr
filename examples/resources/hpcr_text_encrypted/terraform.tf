terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

resource "hpcr_text_encrypted" "text" {
  text = "hello world"
}

resource "hpcr_text_encrypted" "text_platform" {
  text = "hello world"
  platform = "hpvs"
}

resource "hpcr_text_encrypted" "text_cert" {
  text = "hello world"
  cert = file("./cert/encrypt.crt")
}

output "hpcr_text_rendered" {
  value = hpcr_text_encrypted.text.rendered
}

output "hpcr_text_sha256_in" {
  value = hpcr_text_encrypted.text.sha256_in
}

output "hpcr_text_sha256_out" {
  value = hpcr_text_encrypted.text.sha256_out
}

output "hpcr_text_platform_rendered" {
  value = hpcr_text_encrypted.text_platform.rendered
}

output "hpcr_text_platform_sha256_in" {
  value = hpcr_text_encrypted.text_platform.sha256_in
}

output "hpcr_text_platform_sha256_out" {
  value = hpcr_text_encrypted.text_platform.sha256_out
}

output "hpcr_text_cert_rendered" {
  value = hpcr_text_encrypted.text_cert.rendered
}

output "hpcr_text_cert_sha256_in" {
  value = hpcr_text_encrypted.text_cert.sha256_in
}

output "hpcr_text_cert_sha256_out" {
  value = hpcr_text_encrypted.text_cert.sha256_out
}
