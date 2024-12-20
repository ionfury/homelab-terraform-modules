<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.6.6 |
| <a name="requirement_local"></a> [local](#requirement\_local) | 2.1.0 |
| <a name="requirement_talos"></a> [talos](#requirement\_talos) | 0.7.0-alpha.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_local"></a> [local](#provider\_local) | 2.1.0 |
| <a name="provider_talos"></a> [talos](#provider\_talos) | 0.7.0-alpha.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [local_file.kubeconfig](https://registry.terraform.io/providers/hashicorp/local/2.1.0/docs/resources/file) | resource |
| [local_file.talosconfig](https://registry.terraform.io/providers/hashicorp/local/2.1.0/docs/resources/file) | resource |
| [talos_cluster_kubeconfig.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0-alpha.0/docs/resources/cluster_kubeconfig) | resource |
| [talos_machine_bootstrap.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0-alpha.0/docs/resources/machine_bootstrap) | resource |
| [talos_machine_configuration_apply.hosts](https://registry.terraform.io/providers/siderolabs/talos/0.7.0-alpha.0/docs/resources/machine_configuration_apply) | resource |
| [talos_machine_secrets.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0-alpha.0/docs/resources/machine_secrets) | resource |
| [talos_client_configuration.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0-alpha.0/docs/data-sources/client_configuration) | data source |
| [talos_cluster_health.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0-alpha.0/docs/data-sources/cluster_health) | data source |
| [talos_machine_configuration.control_plane](https://registry.terraform.io/providers/siderolabs/talos/0.7.0-alpha.0/docs/data-sources/machine_configuration) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_allow_scheduling_on_controlplane"></a> [allow\_scheduling\_on\_controlplane](#input\_allow\_scheduling\_on\_controlplane) | Whether to allow scheduling on the controlplane. | `bool` | `true` | no |
| <a name="input_cluster_vip"></a> [cluster\_vip](#input\_cluster\_vip) | The VIP to use for the Talos cluster. Applied to the first interface of control plane hosts. | `string` | `"192.168.10.5"` | no |
| <a name="input_endpoint"></a> [endpoint](#input\_endpoint) | The endpoint for the Talos cluster. | `string` | `"https://192.168.10.246:6443"` | no |
| <a name="input_host_dns_enabled"></a> [host\_dns\_enabled](#input\_host\_dns\_enabled) | Whether to enable host DNS. | `bool` | `true` | no |
| <a name="input_host_dns_forwardKubeDNSToHost"></a> [host\_dns\_forwardKubeDNSToHost](#input\_host\_dns\_forwardKubeDNSToHost) | Whether to forward kube DNS to the host. | `bool` | `true` | no |
| <a name="input_host_dns_resolveMemberNames"></a> [host\_dns\_resolveMemberNames](#input\_host\_dns\_resolveMemberNames) | Whether to resolve member names. | `bool` | `true` | no |
| <a name="input_hosts"></a> [hosts](#input\_hosts) | A map of current hosts from which to build the Talos cluster. | <pre>map(object({<br/>    cluster = object({<br/>      member = string<br/>      role   = string<br/>    })<br/>    disk = object({<br/>      install = string<br/>    })<br/>    interfaces = list(object({<br/>      hardwareAddr     = string<br/>      addresses        = list(string)<br/>      dhcp_routeMetric = number<br/>      vlans = list(object({<br/>        vlanId           = number<br/>        addresses        = list(string)<br/>        dhcp_routeMetric = number<br/>      }))<br/>    }))<br/>    ipmi = object({<br/>      ip  = string<br/>      mac = string<br/>    })<br/>  }))</pre> | <pre>{<br/>  "node46": {<br/>    "cluster": {<br/>      "member": "cluster",<br/>      "role": "controlplane"<br/>    },<br/>    "disk": {<br/>      "install": "/dev/sda"<br/>    },<br/>    "interfaces": [<br/>      {<br/>        "addresses": [<br/>          "192.168.10.246"<br/>        ],<br/>        "dhcp_routeMetric": 100,<br/>        "hardwareAddr": "ac:1f:6b:2d:c0:22",<br/>        "vlans": [<br/>          {<br/>            "addresses": [<br/>              "192.168.20.10"<br/>            ],<br/>            "dhcp_routeMetric": 100,<br/>            "vlanId": 10<br/>          }<br/>        ]<br/>      }<br/>    ],<br/>    "ipmi": {<br/>      "ip": "192.168.10.231",<br/>      "mac": "ac:1f:6b:68:2b:e1"<br/>    }<br/>  }<br/>}</pre> | no |
| <a name="input_kubernetes_config_path"></a> [kubernetes\_config\_path](#input\_kubernetes\_config\_path) | The path to the Kubernetes configuration file. | `string` | `"~/.kube"` | no |
| <a name="input_kubernetes_version"></a> [kubernetes\_version](#input\_kubernetes\_version) | The version of kubernetes to deploy. | `string` | `"1.30.1"` | no |
| <a name="input_name"></a> [name](#input\_name) | A name to provide for the Talos cluster. | `string` | `"cluster"` | no |
| <a name="input_nameservers"></a> [nameservers](#input\_nameservers) | A list of nameservers to use for the Talos cluster. | `list(string)` | <pre>[<br/>  "1.1.1.1",<br/>  "1.0.0.1"<br/>]</pre> | no |
| <a name="input_ntp_servers"></a> [ntp\_servers](#input\_ntp\_servers) | A list of NTP servers to use for the Talos cluster. | `list(string)` | <pre>[<br/>  "0.pool.ntp.org",<br/>  "1.pool.ntp.org"<br/>]</pre> | no |
| <a name="input_talos_config_path"></a> [talos\_config\_path](#input\_talos\_config\_path) | The path to the Talos configuration file. | `string` | `"~/.talos"` | no |
| <a name="input_talos_version"></a> [talos\_version](#input\_talos\_version) | The version of Talos to use. | `string` | `"v1.8.3"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_kubernetes_config_file_path"></a> [kubernetes\_config\_file\_path](#output\_kubernetes\_config\_file\_path) | n/a |
| <a name="output_talos_config_file_path"></a> [talos\_config\_file\_path](#output\_talos\_config\_file\_path) | n/a |
<!-- END_TF_DOCS -->