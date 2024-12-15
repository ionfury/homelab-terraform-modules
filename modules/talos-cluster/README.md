<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.6.6 |
| <a name="requirement_local"></a> [local](#requirement\_local) | 2.1.0 |
| <a name="requirement_talos"></a> [talos](#requirement\_talos) | 0.6.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_local"></a> [local](#provider\_local) | 2.1.0 |
| <a name="provider_talos"></a> [talos](#provider\_talos) | 0.6.1 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [local_file.kubeconfig](https://registry.terraform.io/providers/hashicorp/local/2.1.0/docs/resources/file) | resource |
| [local_file.talosconfig](https://registry.terraform.io/providers/hashicorp/local/2.1.0/docs/resources/file) | resource |
| [talos_cluster_kubeconfig.this](https://registry.terraform.io/providers/siderolabs/talos/0.6.1/docs/resources/cluster_kubeconfig) | resource |
| [talos_machine_bootstrap.this](https://registry.terraform.io/providers/siderolabs/talos/0.6.1/docs/resources/machine_bootstrap) | resource |
| [talos_machine_configuration_apply.hosts](https://registry.terraform.io/providers/siderolabs/talos/0.6.1/docs/resources/machine_configuration_apply) | resource |
| [talos_machine_secrets.this](https://registry.terraform.io/providers/siderolabs/talos/0.6.1/docs/resources/machine_secrets) | resource |
| [talos_client_configuration.this](https://registry.terraform.io/providers/siderolabs/talos/0.6.1/docs/data-sources/client_configuration) | data source |
| [talos_machine_configuration.control_plane](https://registry.terraform.io/providers/siderolabs/talos/0.6.1/docs/data-sources/machine_configuration) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_endpoint"></a> [endpoint](#input\_endpoint) | The endpoint for the Talos cluster | `string` | `"https://cluster.local:6443"` | no |
| <a name="input_hosts"></a> [hosts](#input\_hosts) | A map of current hosts.  Hosts to join the cluster are determined by their cluster.member label matching var.name. | <pre>map(object({<br/>    cluster = object({<br/>      member = string<br/>      role   = string<br/>    })<br/>    disk = object({<br/>      install = string<br/>    })<br/>    lan = list(object({<br/>      ip  = string<br/>      mac = string<br/>    }))<br/>    ipmi = object({<br/>      ip  = string<br/>      mac = string<br/>    })<br/>  }))</pre> | n/a | yes |
| <a name="input_kubernetes_config_path"></a> [kubernetes\_config\_path](#input\_kubernetes\_config\_path) | The path to the Kubernetes configuration file | `string` | `"~/.kube"` | no |
| <a name="input_kubernetes_version"></a> [kubernetes\_version](#input\_kubernetes\_version) | The version of kubernetes to deploy | `string` | `"1.30.1"` | no |
| <a name="input_name"></a> [name](#input\_name) | A name to provide for the Talos cluster | `string` | `"cluster"` | no |
| <a name="input_talos_config_path"></a> [talos\_config\_path](#input\_talos\_config\_path) | The path to the Talos configuration file | `string` | `"~/.talos"` | no |
| <a name="input_talos_version"></a> [talos\_version](#input\_talos\_version) | The version of Talos to use | `string` | `"1.8.3"` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->