locals {
  endpoints          = [for host_key, host in var.hosts : host.lan[0].ip if host.cluster.role == "controlplane"]
  nodes              = [for host_key, host in var.hosts : host_key]
  controlplane_nodes = [for host_key, host in var.hosts : host_key if host.cluster.role == "controlplane"]
  worker_nodes       = [for host_key, host in var.hosts : host_key if host.cluster.role == "worker"]
}

resource "talos_machine_secrets" "this" {
  talos_version = var.talos_version
}

data "talos_machine_configuration" "control_plane" {
  machine_type = "controlplane"

  cluster_name       = var.name
  cluster_endpoint   = var.endpoint
  kubernetes_version = var.kubernetes_version
  talos_version      = var.talos_version
  machine_secrets    = data.talos_machine_secrets.this.machine_secrets
}

data "talos_machine_configuration" "worker" {
  machine_type = "worker"

  cluster_name       = var.name
  cluster_endpoint   = var.endpoint
  kubernetes_version = var.kubernetes_version
  talos_version      = var.talos_version
  machine_secrets    = data.talos_machine_secrets.this.machine_secrets
}

data "talos_client_configuration" "this" {
  cluster_name         = var.name
  client_configuration = data.talos_machine_secrets.this.client_configuration
  endpoints            = local.endpoints
  nodes                = local.nodes
}

data "talos_cluster_health" "this" {
  client_configuration   = data.talos_client_configuration.this.client_configuration
  endpoints              = local.endpoints
  control_plane_nodes    = local.controlplane_nodes
  worker_nodes           = local.worker_nodes
  skip_kubernetes_checks = false

  timeouts = {
    read = "5m"
  }
}

resource "local_file" "kubeconfig" {
  depends_on = [data.talos_cluster_health.this]
  content    = talos_cluster_kubeconfig.this.kubeconfig_raw
  filename   = pathexpand("${var.kubernetes_config_path}/${var.name}")

  file_permission = "0644"
}


resource "local_file" "talosconfig" {
  depends_on = [data.talos_cluster_health.this]
  content    = data.talos_client_configuration.this.talos_config
  filename   = pathexpand("${var.talos_config_path}/${var.name}")

  file_permission = "0644"
}
