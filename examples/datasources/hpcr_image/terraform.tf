terraform {
  required_providers {
    hpcr = {
      source  = "ibm-hyper-protect/hpcr"
      version = ">= 1.2.0"
    }

    ibm = {
      source  = "IBM-Cloud/ibm"
      version = ">= 1.37.1"
    }
  }
}

provider "ibm" {
  region = "us-south"
  zone   = "us-south-3"
}

data "ibm_is_images" "ibm_images" {
  visibility = "public"
  status     = "available"
}

data "hpcr_image" "hyper_protect_image" {
  images = jsonencode(data.ibm_is_images.ibm_images.images)
}

output "hpcr_image_id" {
  value = data.hpcr_image.hyper_protect_image.id
}

output "hpcr_image_image" {
  value = data.hpcr_image.hyper_protect_image.image
}

output "hpcr_image_sha256" {
  value = data.hpcr_image.hyper_protect_image.sha256
}

output "hpcr_image_version" {
  value = data.hpcr_image.hyper_protect_image.version
}
