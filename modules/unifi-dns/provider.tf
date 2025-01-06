provider "aws" {
  region  = var.aws.region
  profile = var.aws.profile
}

data "aws_ssm_parameter" "unifi_username" {
  name = var.unifi.username_store
}

data "aws_ssm_parameter" "unifi_password" {
  name = var.unifi.password_store
}

provider "unifi" {
  api_url        = var.unifi.address
  username       = data.aws_ssm_parameter.unifi_username.value
  password       = data.aws_ssm_parameter.unifi_password.value
  allow_insecure = true
  site           = var.unifi.site
}
