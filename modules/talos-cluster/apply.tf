locals {
  bootstrap_node     = [for host_key, host in var.hosts : host_key if host.cluster.role == "controlplane"][0]
  bootstrap_endpoint = [for host_key, host in var.hosts : host.interfaces[0].addresses[0] if host.cluster.role == "controlplane"][0]
}

resource "talos_machine_configuration_apply" "hosts" {
  for_each = var.hosts

  client_configuration        = talos_machine_secrets.this.client_configuration
  machine_configuration_input = data.talos_machine_configuration.control_plane.machine_configuration
  node                        = each.key
  endpoint                    = each.value.interfaces[0].addresses[0]

  config_patches = [
    each.value.cluster.role == "controlplane" ? templatefile("${path.module}/resources/templates/scheduling_cp.yaml.tmpl", {
      allow_scheduling_on_controlplane = var.allow_scheduling_on_controlplane
    }) : null,

    templatefile("${path.module}/resources/templates/hostname.yaml.tmpl", {
      hostname = each.key
    }),

    templatefile("${path.module}/resources/templates/install_disk.yaml.tmpl", {
      install_disk = each.value.disk.install
    }),

    templatefile("${path.module}/resources/templates/nameservers.yaml.tmpl", {
      nameservers = var.nameservers
    }),

    templatefile("${path.module}/resources/templates/ntp_servers.yaml.tmpl", {
      ntp_servers = var.ntp_servers
    }),

    templatefile("${path.module}/resources/templates/host_dns.yaml.tmpl", {
      host_dns_enabled         = var.host_dns_enabled
      resolve_member_names     = var.host_dns_resolveMemberNames
      forward_kube_dns_to_host = var.host_dns_forwardKubeDNSToHost
    }),

    templatefile("${path.module}/resources/templates/interfaces.yaml.tmpl", {
      cluster_vip = var.cluster_vip
      interfaces  = each.value.interfaces
    }),
    /*
    var.ingress_firewall_enabled == true && each.value.cluster.role == "controlplane" ? templatefile("${path.module}/resources/templates/firewall_cp.yaml.tmpl", {
      cni_vxlan_port    = var.cni_vxlan_port
      cluster_subnet    = var.cluster_subnet
      control_plane_ips = local.controlplane_ips
    }) : null,

    var.ingress_firewall_enabled == true && each.value.cluster.role == "worker" ? templatefile("${path.module}/resources/templates/firewall_worker.yaml.tmpl", {
      cni_vxlan_port = var.cni_vxlan_port
      cluster_subnet = var.cluster_subnet
    }) : null,
*/
    file("${path.module}/resources/files/longhorn.yaml"),
  ]
}

resource "talos_machine_bootstrap" "this" {
  depends_on = [talos_machine_configuration_apply.hosts]

  client_configuration = talos_machine_secrets.this.client_configuration
  node                 = local.bootstrap_node
  endpoint             = local.bootstrap_endpoint
}

resource "talos_cluster_kubeconfig" "this" {
  depends_on = [talos_machine_bootstrap.this]

  client_configuration = talos_machine_secrets.this.client_configuration
  node                 = local.bootstrap_node
  endpoint             = local.bootstrap_endpoint
}
