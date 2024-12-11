data "aws_ssm_parameter" "pxeboot_user" {
  name = lookup(var.raspberry_pis, var.pxeboot_host).ssh.user_store
}

data "aws_ssm_parameter" "pxeboot_password" {
  name = lookup(var.raspberry_pis, var.pxeboot_host).ssh.pass_store
}

resource "ansible_playbook" "setup_iptables" {
  playbook                = "${path.module}/playbooks/setup_iptables.yaml"
  name                    = lookup(var.raspberry_pis, var.pxeboot_host).ip
  replayable              = true
  ignore_playbook_failure = false
  extra_vars = {
    ansible_user     = data.aws_ssm_parameter.pxeboot_user.value
    ansible_password = data.aws_ssm_parameter.pxeboot_password.value
  }
}

resource "ansible_playbook" "setup_tftp_server" {
  playbook                = "${path.module}/playbooks/setup_tftp_server.yaml"
  name                    = lookup(var.raspberry_pis, var.pxeboot_host).ip
  replayable              = true
  ignore_playbook_failure = false
  extra_vars = {
    ansible_user     = data.aws_ssm_parameter.pxeboot_user.value
    ansible_password = data.aws_ssm_parameter.pxeboot_password.value
  }
}

resource "ansible_playbook" "setup_ipxe" {
  depends_on              = [ansible_playbook.setup_tftp_server]
  playbook                = "${path.module}/playbooks/setup_ipxe.yaml"
  name                    = lookup(var.raspberry_pis, var.pxeboot_host).ip
  replayable              = true
  ignore_playbook_failure = false
  extra_vars = {
    ansible_user     = data.aws_ssm_parameter.pxeboot_user.value
    ansible_password = data.aws_ssm_parameter.pxeboot_password.value
    schematics_dir   = "${path.module}/resources/schematics"
    scripts_dir      = "${path.module}/resources/scripts"
  }
}
