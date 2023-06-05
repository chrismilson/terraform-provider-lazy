terraform {
  required_providers {
    lazy = {
      source = "registry.terraform.io/chrismilson/lazy"
    }
  }
}

provider "lazy" {}

variable "value" {
  type = string
  default = null
}

resource "lazy_string" "this" {
  initially  = "hello"
  explicitly = var.value
}

output "value" {
  value = lazy_string.this.result
}

