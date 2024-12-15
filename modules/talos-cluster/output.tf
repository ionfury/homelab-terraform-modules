output "talos_config_file_path" {
  value = local_file.talosconfig.filename
}

output "kubernetes_config_file_path" {
  value = local_file.kubeconfig.filename
}
