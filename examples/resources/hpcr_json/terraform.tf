terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

resource "hpcr_json" "json_data1" {
  json = <<JSON
  {
    "workload" : "value1",
    "env" : "value2"
  }
  JSON
}

resource "hpcr_json" "json_data2" {
  json = <<JSON
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
}

output "json1_rendered" {
  value = hpcr_json.json_data1.rendered
}

output "json1_sha256_in" {
  value = hpcr_json.json_data1.sha256_in
}

output "json1_sha256_out" {
  value = hpcr_json.json_data1.sha256_out
}

output "json2_rendered" {
  value = hpcr_json.json_data2.rendered
}

output "json2_sha256_in" {
  value = hpcr_json.json_data2.sha256_in
}

output "json2_sha256_out" {
  value = hpcr_json.json_data2.sha256_out
}
