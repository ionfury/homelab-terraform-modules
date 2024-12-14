resource "talos_machine_secrets" "this" {
  talos_version = var.talos_version
}

data "talos_machine_configuration" "control_plane" {
  cluster_name       = var.name
  cluster_endpoint   = var.endpoint
  kubernetes_version = var.kubernetes_version
  talos_version      = var.talos_version
  machine_type       = "controlplane"
  machine_secrets    = talos_machine_secrets.this.machine_secrets
}

data "talos_client_configuration" "this" {
  cluster_name         = var.name
  client_configuration = talos_machine_secrets.this.client_configuration
  endpoints            = [for host_key, host in var.hosts : host.lan[0].ip if host.cluster.role == "controlplane" && host.cluster.member == var.name]
}

resource "local_file" "kubeconfig" {
  content  = talos_cluster_kubeconfig.this.kubeconfig_raw
  filename = pathexpand("${var.kubernetes_config_path}/${var.name}")

  file_permission = "0644"
}

resource "local_file" "talosconfig" {
  content  = data.talos_client_configuration.this.talos_config
  filename = pathexpand("${var.talos_config_path}/${var.name}")

  file_permission = "0644"
}
