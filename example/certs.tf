data "hpcr_encryption_certs" "enc_certs" {
  versions = ["1.0.10", "1.0.11"]
}

data "hpcr_encryption_cert" "enc_cert" {
  certs = data.hpcr_encryption_certs.enc_certs.certs
  spec = "1.0.10"
}

output "cert_version" {
  value = data.hpcr_encryption_cert.enc_cert.version
}