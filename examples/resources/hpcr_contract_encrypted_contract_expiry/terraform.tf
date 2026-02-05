terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = ">= 1.2.0"
    }
  }
}

resource "hpcr_tgz" "contract" {
  folder = "pods"
}

resource "hpcr_text" "attestation_public_key" {
  text = file("./cert/public.pem")
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
    "attestationPublicKey" : hpcr_text.attestation_public_key.rendered
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
  contract  = local.contract
  expiry    = 30
  cakey     = file("./cert/personal_ca.pem")
  cacert    = file("./cert/personal_ca.crt")
  csrparams = local.csrParams
}

resource "hpcr_contract_encrypted_contract_expiry" "contract_csr" {
  contract = local.contract
  expiry   = 30
  cakey    = file("./cert/personal_ca.pem")
  cacert   = file("./cert/personal_ca.crt")
  csr      = file("./cert/csr.pem")
}

resource "hpcr_contract_encrypted_contract_expiry" "contract_cert" {
  contract  = local.contract
  cert      = file("./cert/encrypt.crt")
  expiry    = 30
  cakey     = file("./cert/personal_ca.pem")
  cacert    = file("./cert/personal_ca.crt")
  csrparams = local.csrParams
}

resource "hpcr_contract_encrypted_contract_expiry" "contract_privkey" {
  contract = local.contract
  privkey  = file("./cert/private.pem")
  expiry   = 30
  cakey    = file("./cert/personal_ca.pem")
  cacert   = file("./cert/personal_ca.crt")
  csr      = file("./cert/csr.pem")
}

resource "local_file" "contract_rendered" {
  filename = "${path.module}/build/contract.yaml"
  content = hpcr_contract_encrypted_contract_expiry.contract.rendered
}

resource "local_file" "contract_sha256_in" {
  filename = "${path.module}/build/contract.in.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract.sha256_in
}

resource "local_file" "contract_sha256_out" {
  filename = "${path.module}/build/contract.out.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract.sha256_out
}

resource "local_file" "contract_csr_rendered" {
  filename = "${path.module}/build/contract_csr.yaml"
  content = hpcr_contract_encrypted_contract_expiry.contract_csr.rendered
}

resource "local_file" "contract_csr_sha256_in" {
  filename = "${path.module}/build/contract_csr.in.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract_csr.sha256_in
}

resource "local_file" "contract_csr_sha256_out" {
  filename = "${path.module}/build/contract_csr.out.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract_csr.sha256_out
}

resource "local_file" "contract_cert_rendered" {
  filename = "${path.module}/build/contract_cert.yaml"
  content = hpcr_contract_encrypted_contract_expiry.contract_cert.rendered
}

resource "local_file" "contract_cert_sha256_in" {
  filename = "${path.module}/build/contract_cert.in.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract_cert.sha256_in
}

resource "local_file" "contract_cert_sha256_out" {
  filename = "${path.module}/build/contract_cert.out.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract_cert.sha256_out
}

resource "local_file" "contract_privkey_rendered" {
  filename = "${path.module}/build/contract_privkey.yaml"
  content = hpcr_contract_encrypted_contract_expiry.contract_privkey.rendered
}

resource "local_file" "contract_privkey_sha256_in" {
  filename = "${path.module}/build/contract_privkey.in.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract_privkey.sha256_in
}

resource "local_file" "contract_privkey_sha256_out" {
  filename = "${path.module}/build/contract_privkey.out.sha256"
  content = hpcr_contract_encrypted_contract_expiry.contract_privkey.sha256_out
}
