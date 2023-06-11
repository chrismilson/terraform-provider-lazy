terraform {
  required_providers {
    lazy = {
      source = "chrismilson/lazy"
    }
  }
}

variable "image_digest" {
  type    = string
  default = null
}

resource "lazy_string" "image_digest" {
  explicitly = var.image_digest
}

resource "kubernetes_pod" "example" {
  spec {
    container {
      image = "alpine@${lazy_string.image_digest}"
    }
  }
}
