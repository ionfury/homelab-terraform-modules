variable "name" {
  description = "A name to provide for the Talos cluster"
  type        = string
  default     = "cluster"
}

variable "endpoint" {
  description = "The endpoint for the Talos cluster"
  type        = string
  default     = "https://cluster.local:6443"
}

variable "kubernetes_version" {
  description = "The version of kubernetes to deploy"
  type        = string
  default     = "1.30.1"
}

variable "talos_version" {
  description = "The version of Talos to use"
  type        = string
  default     = "1.8.0"
}

