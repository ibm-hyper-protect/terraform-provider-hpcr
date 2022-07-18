variable "FOLDER" {
  type = string
}

resource "hpcr_tgz_encrypted" "sample" {
  folder = var.FOLDER
}

output "result" {
  value = resource.hpcr_tgz_encrypted.sample.rendered
}
