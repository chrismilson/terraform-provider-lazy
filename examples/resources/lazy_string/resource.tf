terraform {
  required_providers {
    lazy = {
      source = "chrismilson/lazy"
    }
  }
}

variable "enable_public_ingress" {
  type    = bool
  default = null
}

resource "lazy_string" "enable_public_ingress" {
  initially  = false
  explicitly = var.enable_public_ingress
}

resource "kubernetes_ingress_v1" "example" {
  for_each = lazy_string.enable_public_ingress ? { enabled = true } : {}
  // ...
}
