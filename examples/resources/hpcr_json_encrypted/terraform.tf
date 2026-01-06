terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

resource "hpcr_json_encrypted" "json_data" {
  json = <<JSON
  {
    "workload" : "value1",
    "env" : "value2"
  }
  JSON
}

resource "hpcr_json_encrypted" "json_data_platform" {
  json     = <<JSON
  {
    "workload": {
        "compose": {
            "archive": "testing"
        }
    },
    "env": {
        "logging": "testing"
    }
  }
  JSON

  platform = "hpvs"
}

resource "hpcr_json_encrypted" "json_data_cert" {
  json = <<JSON
  {
    "workload" : "value1",
    "env" : "value2"
  }
  JSON

  cert = file("./cert/encrypt.crt")
}

output "json_data_rendered" {
    value = hpcr_json_encrypted.json_data.rendered
}

output "json_data_sha256_in" {
    value = hpcr_json_encrypted.json_data.sha256_in
}

output "json_data_sha256_out" {
    value = hpcr_json_encrypted.json_data.sha256_out
}

output "json_data_platform_rendered" {
    value = hpcr_json_encrypted.json_data_platform.rendered
}

output "json_data_platform_sha256_in" {
    value = hpcr_json_encrypted.json_data_platform.sha256_in
}

output "json_data_platform_sha256_out" {
    value = hpcr_json_encrypted.json_data_platform.sha256_out
}

output "json_data_cert_rendered" {
    value = hpcr_json_encrypted.json_data_cert.rendered
}

output "json_data_cert_sha256_in" {
    value = hpcr_json_encrypted.json_data_cert.sha256_in
}

output "json_data_cert_sha256_out" {
    value = hpcr_json_encrypted.json_data_cert.sha256_out
}