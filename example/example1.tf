variable "FOLDER" {
  type        = string
}

data "hpcr_tgz" "sample" {
  folder = var.FOLDER
}

output "result" {
  value = data.hpcr_tgz.sample.rendered
}
