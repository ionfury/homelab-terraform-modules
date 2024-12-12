variable "aws" {
  description = "AWS account information"
  type = object({
    region  = string
    profile = string
  })
}

variable "unifi" {
  description = "Unifi controller information"
  type = object({
    address        = string
    username_store = string
    password_store = string
    site           = string
  })
}

variable "unifi_users" {
  description = "List of users to add to the Unifi controller."
  type = map(object({
    ip  = string
    mac = string
  }))
}
