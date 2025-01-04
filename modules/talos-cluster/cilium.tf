data "helm_template" "cilium" {
  repository   = "https://helm.cilium.io/"
  chart        = "cilium"
  name         = "cilium"
  version      = var.cilium_version
  kube_version = var.kubernetes_version
  namespace    = "kube-system"
  values = [
    templatefile("${path.module}/resources/templates/cilium.yaml.tmpl", {
      cluster_name   = var.name
      cluster_id     = var.cluster_id
      pod_subnet     = var.pod_subnet
      service_subnet = var.service_subnet
    })
  ]
}
