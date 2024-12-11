## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_ansible"></a> [ansible](#requirement\_ansible) | ~> 1.3.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~>5.80.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_ansible"></a> [ansible](#provider\_ansible) | ~> 1.3.0 |
| <a name="provider_aws"></a> [aws](#provider\_aws) | ~>5.80.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [ansible_playbook.setup_iptables](https://registry.terraform.io/providers/ansible/ansible/latest/docs/resources/playbook) | resource |
| [ansible_playbook.setup_ipxe](https://registry.terraform.io/providers/ansible/ansible/latest/docs/resources/playbook) | resource |
| [ansible_playbook.setup_tftp_server](https://registry.terraform.io/providers/ansible/ansible/latest/docs/resources/playbook) | resource |
| [aws_ssm_parameter.pxeboot_password](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ssm_parameter) | data source |
| [aws_ssm_parameter.pxeboot_user](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ssm_parameter) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws"></a> [aws](#input\_aws) | AWS account information. | <pre>object({<br>    region  = string<br>    profile = string<br>  })</pre> | n/a | yes |
| <a name="input_pxeboot_host"></a> [pxeboot\_host](#input\_pxeboot\_host) | Name of the raspberry pi to use as the host for pxebootings | `string` | n/a | yes |
| <a name="input_raspberry_pis"></a> [raspberry\_pis](#input\_raspberry\_pis) | Map of raspberry pis with their IP and MAC addresses and ssh credential stores | <pre>map(object({<br>    ip  = string<br>    mac = string<br>    ssh = object({<br>      user_store = string<br>      pass_store = string<br>    })<br>  }))</pre> | n/a | yes |
| <a name="input_schematics_dir"></a> [schematics\_dir](#input\_schematics\_dir) | Directory containing schematics to be copied to the raspberry pi | `string` | n/a | yes |
| <a name="input_scripts_dir"></a> [scripts\_dir](#input\_scripts\_dir) | Directory containing scripts to be copied to the raspberry pi | `string` | n/a | yes |

## Outputs

No outputs.
