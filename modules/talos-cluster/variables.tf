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

variable "cilium_version" {
  description = "The version of Cilium to use."
  type        = string
  default     = "1.16.5"
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

variable "cluster_vip" {
  description = "The VIP to use for the Talos cluster. Applied to the first interface of control plane hosts."
  type        = string
  default     = "192.168.10.5"
}

variable "cluster_subnet" {
  description = "The subnet to use for the Talos cluster nodes."
  type        = string
  default     = "192.168.10.0/24"
}

variable "pod_subnet" {
  description = "The pod subnet to use for the Talos cluster."
  type        = string
  default     = "172.16.0.0/16"
}

variable "service_subnet" {
  description = "The pod subnet to use for the Talos cluster."
  type        = string
  default     = "172.17.0.0/16"
}

variable "allow_scheduling_on_controlplane" {
  description = "Whether to allow scheduling on the controlplane."
  type        = bool
  default     = true
}

variable "host_dns_enabled" {
  description = "Whether to enable host DNS."
  type        = bool
  default     = true
}

variable "host_dns_resolveMemberNames" {
  description = "Whether to resolve member names."
  type        = bool
  default     = true
}

variable "host_dns_forwardKubeDNSToHost" {
  description = "Whether to forward kube DNS to the host."
  type        = bool
  default     = true
}

variable "hosts" {
  description = "A map of current hosts from which to build the Talos cluster."
  type = map(object({
    cluster = object({
      member = string
      role   = string
    })
    install = object({
      diskSelector = list(string) # https://www.talos.dev/v1.9/reference/configuration/v1alpha1/config/#Config.machine.install.diskSelector
    })
    interfaces = list(object({
      hardwareAddr     = string
      addresses        = list(string)
      dhcp_routeMetric = number
      vlans = list(object({
        vlanId           = number
        addresses        = list(string)
        dhcp_routeMetric = number
      }))
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
      install = {
        diskSelector = ["type: 'ssd'"]
      }
      interfaces = [{
        hardwareAddr     = "ac:1f:6b:2d:c0:22"
        addresses        = ["192.168.10.246"]
        dhcp_routeMetric = 100
        vlans = [{
          vlanId           = 10
          addresses        = ["192.168.20.10"]
          dhcp_routeMetric = 100
        }]
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
