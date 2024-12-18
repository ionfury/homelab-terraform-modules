variable "name" {
  description = "A name to provide for the Talos cluster."
  type        = string
  default     = "cluster"
}

variable "endpoint" {
  description = "The endpoint for the Talos cluster."
  type        = string
  default     = "https://192.168.10.246:6443"
}

variable "kubernetes_version" {
  description = "The version of kubernetes to deploy."
  type        = string
  default     = "1.30.1"
}

variable "talos_version" {
  description = "The version of Talos to use."
  type        = string
  default     = "v1.8.3"
}

variable "talos_config_path" {
  description = "The path to the Talos configuration file."
  type        = string
  default     = "~/.talos"
}

variable "kubernetes_config_path" {
  description = "The path to the Kubernetes configuration file."
  type        = string
  default     = "~/.kube"
}

variable "nameservers" {
  description = "A list of nameservers to use for the Talos cluster."
  type        = list(string)
  default     = ["1.1.1.1", "1.0.0.1"]
}

variable "ntp_servers" {
  description = "A list of NTP servers to use for the Talos cluster."
  type        = list(string)
  default     = ["0.pool.ntp.org", "1.pool.ntp.org"]
}
/*
variable "ingress_firewall_enabled" {
  description = "Whether to enable the ingress firewall for the Talos cluster."
  type        = bool
  default     = true
}

variable "cluster_subnet" {
  description = "The subnet to use for the Talos cluster."
  type        = string
  default     = "192.168.10.0/24"
}

variable "cni_vxlan_port" {
  description = "The port to use for the CNI VXLAN."
  type        = string
  default     = "8472" # Cilium default
}
*/
variable "allow_scheduling_on_controlplane" {
  description = "Whether to allow scheduling on the controlplane."
  type        = bool
  default     = true
}

variable "host_dns" {
  description = "The DNS server to use for the Talos cluster."
  type = object({
    enabled              = bool
    resolveMemberNames   = bool
    forwardKubeDNSToHost = bool
  })
  default = {
    enabled              = true
    resolveMemberNames   = true
    forwardKubeDNSToHost = true
  }
}

variable "hosts" {
  description = "A map of current hosts from which to build the Talos cluster."
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
  default = {
    node46 = {
      cluster = {
        member = "cluster"
        role   = "controlplane"
      }
      disk = {
        install = "/dev/sda"
      }
      lan = [{
        ip  = "192.168.10.246"
        mac = "ac:1f:6b:2d:c0:22"
      }]
      ipmi = {
        ip  = "192.168.10.231"
        mac = "ac:1f:6b:68:2b:e1"
      }
    }
  }

  validation {
    condition     = alltrue([for host in var.hosts : host.cluster.role == "worker" || host.cluster.role == "controlplane" || host.cluster.role == "null"])
    error_message = "The cluster.role must be either 'worker', 'controlplane', or 'null'."
  }
}
