variable "CERTIFICATES" {
  type = string
}

data "hpcr_encryption_cert" "enc_cert" {
  certs = jsondecode(var.CERTIFICATES)
  spec = "1.0.10"
}

output "cert_version" {
  value = data.hpcr_encryption_cert.enc_cert.version
}