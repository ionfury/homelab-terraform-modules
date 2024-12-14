locals {
  cluster_hosts = {
    for host, details in var.hosts : host => details
    if details.cluster.member == var.name
  }
}


resource "talos_machine_configuration_apply" "hosts" {
  client_configuration        = talos_machine_secrets.this.client_configuration
  machine_configuration_input = data.talos_machine_configuration.control_plane.machine_configuration
  for_each                    = local.cluster_hosts
  node                        = each.key
  endpoint                    = each.value.lan[0].ip

  config_patches = [
    templatefile("${path.module}/resources/templates/install-disk-and-hostname.yaml.tmpl", {
      hostname     = each.key
      install_disk = each.value.disk.install
    }),
    each.value.cluster.role == "controlplane" ? file("${path.module}/resources/files/cp-scheduling.yaml") : null
  ]
}

resource "talos_machine_bootstrap" "this" {
  depends_on = [talos_machine_configuration_apply.hosts]

  client_configuration = talos_machine_secrets.this.client_configuration
  node                 = [for host_key, host in local.cluster_hosts : host_key if host.cluster.role == "controlplane"][0]
  endpoint             = [for host_key, host in local.cluster_hosts : host.lan[0].ip if host.cluster.role == "controlplane"][0]
}

resource "talos_cluster_kubeconfig" "this" {
  depends_on           = [talos_machine_bootstrap.this]
  client_configuration = talos_machine_secrets.this.client_configuration
  node                 = [for host_key, host in local.cluster_hosts : host_key if host.cluster.role == "controlplane"][0]
  endpoint             = [for host_key, host in local.cluster_hosts : host.lan[0].ip if host.cluster.role == "controlplane"][0]
}
