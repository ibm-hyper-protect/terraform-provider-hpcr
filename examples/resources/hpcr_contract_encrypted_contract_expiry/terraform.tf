terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = "0.15.0"
    }
  }
}

resource "hpcr_tgz" "contract" {
  folder = "pods"
}

locals {
  # contract in clear text
  contract = yamlencode({
    "env" : {
      "type" : "env",
      "logging" : {
        "logRouter" : {
          "hostname" : "5c2d6b69-c7f0-41bd-b69b-240695369d6e.ingress.us-south.logs.cloud.ibm.com",
          "iamApiKey" : "ab00e3c09p1d4ff7fff9f04c12183413"
        }
      }
    },
    "workload" : {
      "type" : "workload",
      "play" : {
        "archive" : hpcr_tgz.contract.rendered
      }
    },
  })

  csrParams = {
    "country" : "IN",
    "state" : "Karnataka",
    "location" : "Bangalore",
    "org" : "IBM",
    "unit" : "ISDL",
    "domain" : "Hyper Protect",
    "mail" : "Sashwat.K@ibm.com"
  }
}

resource "hpcr_contract_encrypted_contract_expiry" "contract" {
  contract = local.contract
  expiry = 30
  cakey = file("./cert/personal_ca.pem")
  cacert = file("./cert/personal_ca.crt")
  csrparams = local.csrParams
}

resource "hpcr_contract_encrypted_contract_expiry" "contract_csr" {
  contract = local.contract
  expiry = 30
  cakey = file("./cert/personal_ca.pem")
  cacert = file("./cert/personal_ca.crt")
  csr = file("./cert/csr.pem")
}

resource "hpcr_contract_encrypted_contract_expiry" "contract_cert" {
 contract = local.contract
 cert = file("./cert/encrypt.crt")
 expiry = 30
    cakey = file("./cert/personal_ca.pem")
  cacert = file("./cert/personal_ca.crt")
  csrparams = local.csrParams
}

resource "hpcr_contract_encrypted_contract_expiry" "contract_privkey" {
  contract = local.contract
  privkey = file("./cert/private.pem")
  expiry = 30
  cakey = file("./cert/personal_ca.pem")
  cacert = file("./cert/personal_ca.crt")
  csr = file("./cert/csr.pem")
}

output "contract_rendered" {
  value = hpcr_contract_encrypted_contract_expiry.contract.rendered
}

output "contract_sha256_in" {
  value = hpcr_contract_encrypted_contract_expiry.contract.sha256_in
}

output "contract_sha256_out" {
  value = hpcr_contract_encrypted_contract_expiry.contract.sha256_out
}

output "contract_csr_rendered" {
  value = hpcr_contract_encrypted_contract_expiry.contract_csr.rendered
}

output "contract_csr_sha256_in" {
  value = hpcr_contract_encrypted_contract_expiry.contract_csr.sha256_in
}

output "contract_csr_sha256_out" {
  value = hpcr_contract_encrypted_contract_expiry.contract_csr.sha256_out
}

output "contract_cert_rendered" {
  value = hpcr_contract_encrypted_contract_expiry.contract_cert.rendered
}

output "contract_cert_sha256_in" {
  value = hpcr_contract_encrypted_contract_expiry.contract_cert.sha256_in
}

output "contract_cert_sha256_out" {
  value = hpcr_contract_encrypted_contract_expiry.contract_cert.sha256_out
}

output "contract_privkey_rendered" {
  value = hpcr_contract_encrypted_contract_expiry.contract_privkey.rendered
}

output "contract_privkey_sha256_in" {
  value = hpcr_contract_encrypted_contract_expiry.contract_privkey.sha256_in
}

output "contract_privkey_sha256_out" {
  value = hpcr_contract_encrypted_contract_expiry.contract_privkey.sha256_out
}