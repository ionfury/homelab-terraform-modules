resource "talos_machine_secrets" "this" {}


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
  endpoints            = [for node_key, node in var.nodes : node.ip if node.machine_type == "controlplane"]
}
