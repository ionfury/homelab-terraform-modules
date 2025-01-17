variable "name" {
  description = "A name to provide for the Talos cluster."
  type        = string
  default     = "cluster"

  validation {
    condition     = length(var.name) <= 32 && can(regex("^([a-z0-9]+-)*[a-z0-9]+$", var.name))
    error_message = "The name must contain at most 32 characters, begin and end with a lower case alphanumeric character, and may contain lower case alphanumeric characters and dashes between."
  }
}

variable "cluster_id" {
  description = "An ID to provide for the Talos cluster."
  type        = number
  default     = 1

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

variable "prometheus_crd_version" {
  description = "The version of the Prometheus CRD to use."
  type        = string
  default     = "17.0.2"
}

variable "spegal_version" {
  description = "The version of Spegal to use."
  type        = string
  default     = "v0.0.28"
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

variable "node_subnet" {
  description = "The subnet to use for the Talos cluster nodes."
  type        = string
  default     = "192.168.10.0/24"
}

variable "pod_subnet" {
  description = "The pod subnet to use for pods on the Talos cluster."
  type        = string
  default     = "172.16.0.0/16"
}

variable "service_subnet" {
  description = "The pod subnet to use for services on the Talos cluster."
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
    role = string
    install = object({
      diskSelector    = list(string) # https://www.talos.dev/v1.9/reference/configuration/v1alpha1/config/#Config.machine.install.diskSelector
      extraKernelArgs = optional(list(string), [])
      extensions      = optional(list(string), [])
      secureboot      = optional(bool, false)
      wipe            = optional(bool, false)
      architecture    = optional(string, "amd64")
      platform        = optional(string, "metal")

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
  }))

  default = {
    node46 = {
      role = "controlplane"
      install = {
        diskSelector    = ["type: 'ssd'"]
        extraKernelArgs = ["apparmor=0"]
        extensions      = ["iscsi-tools", "util-linux-tools"]
        secureboot      = false
        wipe            = false
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
    }
  }

  validation {
    condition     = alltrue([for host in var.hosts : host.role == "worker" || host.role == "controlplane"])
    error_message = "The host role must be either 'worker', 'controlplane'."
  }
}
