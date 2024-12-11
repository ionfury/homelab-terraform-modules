variable "name" {
  description = "A name to provide for the Talos cluster"
  type        = string
  default     = "cluster"
}

variable "endpoint" {
  description = "The endpoint for the Talos cluster"
  type        = string
  default     = "10.0.0.1"
}

variable "kubernetes_version" {
  description = "The version of kubernetes to deploy"
  type        = string
  default     = "1.30.1"
}

variable "talos_version" {
  description = "The version of Talos to use"
  type        = string
  default     = "1.30.1"
}

variable "nodes" {
  description = "A map of node data describing the cluster."
  type = map(object({
    machine_type = string
    ip           = string
    install_disk = string
    hostname     = optional(string)
  }))

  default = {
    node1 = {
      machine_type = "controlplane"
      ip           = "10.0.0.1"
      install_disk = "/dev/sda"
      hostname     = "control-plane-1"
    }
  }

  validation {
    condition     = alltrue([for node in var.nodes : node.machine_type == "worker" || node.machine_type == "controlplane"])
    error_message = "The machine_type must be either 'worker' or 'controlplane'."
  }
}
