locals {
  bootstrap_node     = [for host_key, host in var.hosts : host_key if host.cluster.role == "controlplane"][0]
  bootstrap_endpoint = [for host_key, host in var.hosts : host.lan[0].ip if host.cluster.role == "controlplane"][0]
}

resource "talos_machine_configuration_apply" "hosts" {
  for_each = var.hosts

  client_configuration        = talos_machine_secrets.this.client_configuration
  machine_configuration_input = data.talos_machine_configuration.control_plane.machine_configuration
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
  node                 = local.bootstrap_node
  endpoint             = local.bootstrap_endpoint
}

resource "talos_cluster_kubeconfig" "this" {
  depends_on = [talos_machine_bootstrap.this]

  client_configuration = talos_machine_secrets.this.client_configuration
  node                 = local.bootstrap_node
  endpoint             = local.bootstrap_endpoint
}
