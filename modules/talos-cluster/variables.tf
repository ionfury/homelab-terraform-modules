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
  default     = "1.8.3"
}

variable "talos_config_path" {
  description = "The path to the Talos configuration file"
  type        = string
  default     = "~/.talos"
}

variable "kubernetes_config_path" {
  description = "The path to the Kubernetes configuration file"
  type        = string
  default     = "~/.kube"
}

variable "hosts" {
  description = "A map of current hosts.  Hosts to join the cluster are determined by their cluster.member label matching var.name."
  type = map(object({
    cluster = object({
      member = string
      role   = string
    })
    disk = object({
      install = string
    })
    lan = list(object({
      ip  = string
      mac = string
    }))
    ipmi = object({
      ip  = string
      mac = string
    })
  }))

  validation {
    condition     = alltrue([for host in var.hosts : host.cluster.role == "worker" || host.cluster.role == "controlplane" || host.cluster.role == "none"])
    error_message = "The cluster.role must be either 'worker', 'controlplane', or 'none'."
  }
}
