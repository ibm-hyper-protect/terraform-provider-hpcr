variable "IMAGES" {
  type = string
}

data "hpcr_image" "selected_image" {
  images= var.IMAGES
}

output "image_version" {
  value = data.hpcr_image.selected_image.version
}