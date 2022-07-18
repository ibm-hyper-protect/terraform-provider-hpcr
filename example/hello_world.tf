variable "FOLDER" {
  type = string
}

variable "LOGDNA_INGESTION_KEY" {
  type = string
}

variable "LOGDNA_INGESTION_HOSTNAME" {
  type = string
}

resource "hpcr_tgz" "compose" {
  folder = var.FOLDER
}

resource "hpcr_text_encrypted" "env" {
  text = yamlencode({
    "logging" : {
      "logDNA" : {
        "ingestionKey" : var.LOGDNA_INGESTION_KEY,
        "hostname" : var.LOGDNA_INGESTION_HOSTNAME
      }
    }
  })
}

resource "hpcr_text_encrypted" "workload" {
  text = yamlencode({
    "compose" : {
      "archive" : resource.hpcr_tgz.compose.rendered
    }
  })
}

output "user_data" {
  value = yamlencode({
    "workload" : resource.hpcr_text_encrypted.workload.rendered,
    "enc" : resource.hpcr_text_encrypted.env.rendered,
  })
}
