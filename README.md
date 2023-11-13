# Terraform Provider HPCR

Implementation of a [terraform provider](https://www.terraform.io/language/providers) to support working with [IBM Cloud Hyper Protect Virtual Server for IBM Cloud VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se).

## Prerequisite

- Installation of [terraform](https://www.terraform.io/downloads) for your platform
- [OpenSSL](https://www.openssl.org/) binary (not LibreSSL), the path to the binary can be configured via the `OPENSSL_BIN` environment variable

## Usage

The [terraform provider](https://www.terraform.io/language/providers) exposes a set of [resources](https://www.terraform.io/language/resources) that help assemble the user_data section for a contract:

### hpcr_tgz

Use this [resource](https://developer.hashicorp.com/terraform/language/resources) to create a tgz archive of your `docker-compose` folder. You can access the `base64` encoded content via the `rendered` property.

```terraform
resource "hpcr_tgz" "compose" {
  folder = var.FOLDER
}
```

### hpcr_text_encrypted

Use this [resource](https://developer.hashicorp.com/terraform/language/resources) to encrypt a string, per default the implementation uses encryption key of the latest HPCR image.

```terraform
resource "hpcr_text_encrypted" "workload" {
  text = yamlencode({
    "compose" : {
      "archive" : resource.hpcr_tgz.compose.rendered
    }
  })
}
```

The typical usecase is to encrypt the `workload` and the `env` section separately and to pass in the yml encoded contract as an input.

### hpcr_image

Use this [datasource](https://developer.hashicorp.com/terraform/language/data-sources) to find the matching HPCR stock image. 

```terraform
data "ibm_is_images" "hyper_protect_images" {
  visibility = "public"
  status     = "available"

}

data "hpcr_image" "selected_image" {
  images= jsonencode(data.ibm_is_images.hyper_protect_images.images)
}
```

This data source accepts a list of available VPC image (e.g. from the VPC [is_images](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/is_images) datasource). The list needs to be serialized to JSON.

Optionally the datasource takes a `spec` parameter that can be used as a [version constraint](https://github.com/Masterminds/semver#checking-version-constraints).

The result of the lookup can be accessed via the following attributes:

- `image`: ID of the selected image
- `version`: semantic version string of the selected image (e.g. `1.0.8`)


## License

[Apache 2.0](LICENSE)

## How to Contribute

The repository uses [semantic-release](https://github.com/semantic-release/semantic-release). Please author the commit messages accordingly.

## References

- [How to Publish a Terraform Provider](https://learn.hashicorp.com/tutorials/terraform/provider-release-publish?in=terraform/providers)