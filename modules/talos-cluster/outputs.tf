output "talosconfig" {
  description = "The talosconfig for the cluster."
  value       = data.talos_client_configuration.this.talos_config
  sensitive   = true
}

output "kubeconfig" {
  description = "The kubeconfig for the cluster."
  value       = talos_cluster_kubeconfig.this.kubeconfig_raw
  sensitive   = true
}
