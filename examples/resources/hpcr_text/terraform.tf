terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

resource "hpcr_text" "text" {
  text = "hello world"
}

output "hpcr_text_rendered" {
  value = hpcr_text.text.rendered
}

output "hpcr_text_sha256_in" {
  value = hpcr_text.text.sha256_in
}

output "hpcr_text_sha256_out" {
  value = hpcr_text.text.sha256_out
}
