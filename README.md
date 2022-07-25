# Terraform Provider HPCR

Implementation of a [terraform provider](https://www.terraform.io/language/providers) to support working with [IBM Cloud Hyper Protect Virtual Server for IBM Cloud VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se).

## Prerequisite

- Installation of [terraform](https://www.terraform.io/downloads) for your platform
- [OpenSSL](https://www.openssl.org/) binary (not LibreSSL), the path to the binary can be configured via the `OPENSSL_BIN` environment variable

## Usage

The [terraform provider](https://www.terraform.io/language/providers) exposes a set of [resources](https://www.terraform.io/language/resources) that help assemble the user_data section for a contract:

### hpcr_tgz

Use this resource to create a tgz archive of your `docker-compose` folder. You can access the `base64` encoded content via the `rendered` property.

```terraform
resource "hpcr_tgz" "compose" {
  folder = var.FOLDER
}
```

### hpcr_text_encrypted

Use this resource to encrypt a string, per default the implementation uses encryption key of the latest HPCR image.

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


## License

[Apache 2.0](LICENSE)

## References

- [How to Publish a Terraform Provider](https://learn.hashicorp.com/tutorials/terraform/provider-release-publish?in=terraform/providers)