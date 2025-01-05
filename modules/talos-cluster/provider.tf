provider "helm" {
  kubernetes {
    config_path = local_file.kubeconfig.filename
  }
}

provider "helm" {
  alias = "bootstrap"
}
