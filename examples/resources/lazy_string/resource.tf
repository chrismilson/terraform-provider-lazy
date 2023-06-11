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

output "public_ingress_status" {
  value = lazy_string.enable_public_ingress.result ? "enabled" : "disabled"
}
