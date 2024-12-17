terraform {
  required_version = ">= 1.6.6"
  required_providers {
    talos = {
      source  = "siderolabs/talos"
      version = "0.7.0-alpha.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.1.0"
    }
  }
}
