<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.6.6 |
| <a name="requirement_helm"></a> [helm](#requirement\_helm) | 2.17.0 |
| <a name="requirement_local"></a> [local](#requirement\_local) | 2.1.0 |
| <a name="requirement_talos"></a> [talos](#requirement\_talos) | 0.7.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_helm"></a> [helm](#provider\_helm) | 2.17.0 |
| <a name="provider_local"></a> [local](#provider\_local) | 2.1.0 |
| <a name="provider_talos"></a> [talos](#provider\_talos) | 0.7.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [helm_release.cilium](https://registry.terraform.io/providers/hashicorp/helm/2.17.0/docs/resources/release) | resource |
| [helm_release.prometheus_crds](https://registry.terraform.io/providers/hashicorp/helm/2.17.0/docs/resources/release) | resource |
| [helm_release.spegel](https://registry.terraform.io/providers/hashicorp/helm/2.17.0/docs/resources/release) | resource |
| [local_file.kubeconfig](https://registry.terraform.io/providers/hashicorp/local/2.1.0/docs/resources/file) | resource |
| [local_file.talosconfig](https://registry.terraform.io/providers/hashicorp/local/2.1.0/docs/resources/file) | resource |
| [talos_cluster_kubeconfig.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/resources/cluster_kubeconfig) | resource |
| [talos_image_factory_schematic.host_schematic](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/resources/image_factory_schematic) | resource |
| [talos_machine_bootstrap.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/resources/machine_bootstrap) | resource |
| [talos_machine_configuration_apply.hosts](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/resources/machine_configuration_apply) | resource |
| [talos_machine_secrets.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/resources/machine_secrets) | resource |
| [talos_client_configuration.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/data-sources/client_configuration) | data source |
| [talos_cluster_health.available](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/data-sources/cluster_health) | data source |
| [talos_cluster_health.this](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/data-sources/cluster_health) | data source |
| [talos_image_factory_extensions_versions.host_version](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/data-sources/image_factory_extensions_versions) | data source |
| [talos_image_factory_urls.host_image_url](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/data-sources/image_factory_urls) | data source |
| [talos_machine_configuration.control_plane](https://registry.terraform.io/providers/siderolabs/talos/0.7.0/docs/data-sources/machine_configuration) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_allow_scheduling_on_controlplane"></a> [allow\_scheduling\_on\_controlplane](#input\_allow\_scheduling\_on\_controlplane) | Whether to allow scheduling on the controlplane. | `bool` | `true` | no |
| <a name="input_cilium_version"></a> [cilium\_version](#input\_cilium\_version) | The version of Cilium to use. | `string` | `"1.16.5"` | no |
| <a name="input_cluster_id"></a> [cluster\_id](#input\_cluster\_id) | An ID to provide for the Talos cluster. | `number` | `1` | no |
| <a name="input_cluster_vip"></a> [cluster\_vip](#input\_cluster\_vip) | The VIP to use for the Talos cluster. Applied to the first interface of control plane hosts. | `string` | `"192.168.10.5"` | no |
| <a name="input_endpoint"></a> [endpoint](#input\_endpoint) | The endpoint for the Talos cluster. | `string` | `"https://192.168.10.246:6443"` | no |
| <a name="input_host_dns_enabled"></a> [host\_dns\_enabled](#input\_host\_dns\_enabled) | Whether to enable host DNS. | `bool` | `true` | no |
| <a name="input_host_dns_forwardKubeDNSToHost"></a> [host\_dns\_forwardKubeDNSToHost](#input\_host\_dns\_forwardKubeDNSToHost) | Whether to forward kube DNS to the host. | `bool` | `true` | no |
| <a name="input_host_dns_resolveMemberNames"></a> [host\_dns\_resolveMemberNames](#input\_host\_dns\_resolveMemberNames) | Whether to resolve member names. | `bool` | `true` | no |
| <a name="input_hosts"></a> [hosts](#input\_hosts) | A map of current hosts from which to build the Talos cluster. | <pre>map(object({<br/>    role = string<br/>    install = object({<br/>      diskSelector    = list(string) # https://www.talos.dev/v1.9/reference/configuration/v1alpha1/config/#Config.machine.install.diskSelector<br/>      extraKernelArgs = optional(list(string), [])<br/>      extensions      = optional(list(string), [])<br/>      secureboot      = optional(bool, false)<br/>      wipe            = optional(bool, false)<br/>      architecture    = optional(string, "amd64")<br/>      platform        = optional(string, "metal")<br/><br/>    })<br/>    interfaces = list(object({<br/>      hardwareAddr     = string<br/>      addresses        = list(string)<br/>      dhcp_routeMetric = number<br/>      vlans = list(object({<br/>        vlanId           = number<br/>        addresses        = list(string)<br/>        dhcp_routeMetric = number<br/>      }))<br/>    }))<br/>  }))</pre> | <pre>{<br/>  "node46": {<br/>    "install": {<br/>      "diskSelector": [<br/>        "type: 'ssd'"<br/>      ],<br/>      "extensions": [<br/>        "iscsi-tools",<br/>        "util-linux-tools"<br/>      ],<br/>      "extraKernelArgs": [<br/>        "apparmor=0"<br/>      ],<br/>      "secureboot": false,<br/>      "wipe": false<br/>    },<br/>    "interfaces": [<br/>      {<br/>        "addresses": [<br/>          "192.168.10.246"<br/>        ],<br/>        "dhcp_routeMetric": 100,<br/>        "hardwareAddr": "ac:1f:6b:2d:c0:22",<br/>        "vlans": [<br/>          {<br/>            "addresses": [<br/>              "192.168.20.10"<br/>            ],<br/>            "dhcp_routeMetric": 100,<br/>            "vlanId": 10<br/>          }<br/>        ]<br/>      }<br/>    ],<br/>    "role": "controlplane"<br/>  }<br/>}</pre> | no |
| <a name="input_kubernetes_config_path"></a> [kubernetes\_config\_path](#input\_kubernetes\_config\_path) | The path to the Kubernetes configuration file. | `string` | `"~/.kube"` | no |
| <a name="input_kubernetes_version"></a> [kubernetes\_version](#input\_kubernetes\_version) | The version of kubernetes to deploy. | `string` | `"1.30.1"` | no |
| <a name="input_name"></a> [name](#input\_name) | A name to provide for the Talos cluster. | `string` | `"cluster"` | no |
| <a name="input_nameservers"></a> [nameservers](#input\_nameservers) | A list of nameservers to use for the Talos cluster. | `list(string)` | <pre>[<br/>  "1.1.1.1",<br/>  "1.0.0.1"<br/>]</pre> | no |
| <a name="input_node_subnet"></a> [node\_subnet](#input\_node\_subnet) | The subnet to use for the Talos cluster nodes. | `string` | `"192.168.10.0/24"` | no |
| <a name="input_ntp_servers"></a> [ntp\_servers](#input\_ntp\_servers) | A list of NTP servers to use for the Talos cluster. | `list(string)` | <pre>[<br/>  "0.pool.ntp.org",<br/>  "1.pool.ntp.org"<br/>]</pre> | no |
| <a name="input_pod_subnet"></a> [pod\_subnet](#input\_pod\_subnet) | The pod subnet to use for pods on the Talos cluster. | `string` | `"172.16.0.0/16"` | no |
| <a name="input_prometheus_crd_version"></a> [prometheus\_crd\_version](#input\_prometheus\_crd\_version) | The version of the Prometheus CRD to use. | `string` | `"17.0.2"` | no |
| <a name="input_service_subnet"></a> [service\_subnet](#input\_service\_subnet) | The pod subnet to use for services on the Talos cluster. | `string` | `"172.17.0.0/16"` | no |
| <a name="input_spegal_version"></a> [spegal\_version](#input\_spegal\_version) | The version of Spegal to use. | `string` | `"v0.0.28"` | no |
| <a name="input_talos_config_path"></a> [talos\_config\_path](#input\_talos\_config\_path) | The path to the Talos configuration file. | `string` | `"~/.talos"` | no |
| <a name="input_talos_version"></a> [talos\_version](#input\_talos\_version) | The version of Talos to use. | `string` | `"v1.8.3"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_kubernetes_config_file_path"></a> [kubernetes\_config\_file\_path](#output\_kubernetes\_config\_file\_path) | n/a |
| <a name="output_talos_config_file_path"></a> [talos\_config\_file\_path](#output\_talos\_config\_file\_path) | n/a |
<!-- END_TF_DOCS -->