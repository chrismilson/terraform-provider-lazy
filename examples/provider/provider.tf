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
  initially  = "sha256:200d1839f0645e5005181b6e6ebd46a040826dec2d6af9320b0f597ec9d27242"
  explicitly = var.image_digest
}

output "image" {
  value = "alpine@${lazy_string.image_digest}"
}
