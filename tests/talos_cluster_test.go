package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTalosClusterSingleNode(t *testing.T) {
	//t.Parallel()

	clusterName := "talos-cluster-" + random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/talos-cluster",
		Vars: map[string]interface{}{
			"name": clusterName,
		},
	}

	defer func() {
		resetClusterToMaintenanceMode(t, terraformOptions)
		terraform.Destroy(t, terraformOptions)
	}()

	terraform.InitAndApply(t, terraformOptions)

	confirmClusterWithTalosctl(t, terraformOptions)

	confirmClusterWithKubectl(t, terraformOptions)
}

func confirmClusterWithTalosctl(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")

	talosctlCmd := shell.Command{
		Command: "talosctl",
		Args:    []string{"--talosconfig", talosConfigFilePath, "-n", "node46", "version"},
	}

	output, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
	assert.NoError(t, err, "Failed to run talosctl command")
	assert.Contains(t, output, "Client Version", "Talos cluster is not up and functional")
}

func confirmClusterWithKubectl(t *testing.T, terraformOptions *terraform.Options) {
	kubeConfigPath := terraform.Output(t, terraformOptions, "kubernetes_config_file_path")

	kubectlCmd := shell.Command{
		Command: "kubectl",
		Args:    []string{"--kubeconfig", kubeConfigPath, "get", "nodes"},
	}

	output, err := shell.RunCommandAndGetOutputE(t, kubectlCmd)
	assert.NoError(t, err, "Failed to run kubectl command")
	assert.Contains(t, output, "Ready", "Kubernetes cluster is not up and functional")
}

func resetClusterToMaintenanceMode(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")

	talosctlCmd := shell.Command{
		Command: "talosctl",
		Args:    []string{"--talosconfig", talosConfigFilePath, "-n", "node46", "reset", "--graceful=false", "--reboot"},
	}

	err := shell.RunCommandE(t, talosctlCmd)
	assert.NoError(t, err, "Failed to reset the cluster into maintenance mode")
}
