terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

resource "hpcr_tgz" "compose_b64" {
  folder = "compose"
}

output "b64_rendered" {
  value = hpcr_tgz.compose_b64.rendered
}

output "b64_sha256_in" {
  value = hpcr_tgz.compose_b64.sha256_in
}

output "b64_sha256_out" {
  value = hpcr_tgz.compose_b64.sha256_out
}