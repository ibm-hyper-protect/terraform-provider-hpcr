variable "FOLDER" {
  type = string
}

resource "hpcr_tgz" "sample" {
  folder = var.FOLDER
}

output "result" {
  value = resource.hpcr_tgz.sample.rendered
}
