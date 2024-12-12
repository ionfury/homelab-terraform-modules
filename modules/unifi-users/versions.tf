terraform {
  required_version = ">= 1.6.6"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>5.80.0"
    }
    unifi = {
      source  = "paultyng/unifi"
      version = "~>0.41.0"
    }
  }
}
