terraform {
  required_version = ">= 1.6.6"
  required_providers {
    talos = {
      source  = "siderolabs/talos"
      version = "0.7.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.1.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.17.0"
    }
    #cilium = {
    #  source  = "littlejo/cilium"
    #  version = "0.2.14"
    #}

  }
}
