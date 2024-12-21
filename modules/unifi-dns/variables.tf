variable "aws" {
  description = "AWS account information."
  type = object({
    region  = string
    profile = string
  })
  default = {
    region  = "us-east-2"
    profile = "terragrunt"
  }
}

variable "unifi" {
  description = "Unifi controller information"
  type = object({
    address        = string
    username_store = string
    password_store = string
    site           = string
  })
  default = {
    address        = "https://192.168.1.1"
    username_store = "/homelab/unifi/terraform/username"
    password_store = "/homelab/unifi/terraform/password"
    site           = "default"
  }
}

variable "unifi_dns_records" {
  description = "List of DNS records to add to the Unifi controller."
  type = map(object({
    name        = optional(string, null)
    value       = string
    enabled     = optional(bool, true)
    port        = optional(number, 0)
    priority    = optional(number, 0)
    record_type = optional(string, "A")
    ttl         = optional(number, 0)
  }))
  validation {
    condition     = alltrue([for record in var.unifi_dns_records : can(regex("^((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d)\\.?\\b){4}$", record.value))])
    error_message = "Each DNS record value must be a valid IP address."
  }
}
