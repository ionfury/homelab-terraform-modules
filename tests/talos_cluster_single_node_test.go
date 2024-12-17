package test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTalosClusterSingleNode(t *testing.T) {
	clusterName := "talos-cluster-" + random.UniqueId()
	endpoint := "https://192.168.10.246:6443"
	kubernetes_version := "1.30.1"
	talos_version := "v1.8.4"
	talos_config_path := "~/.talos"
	kubernetes_config_path := "~/.kube"

	// Sourced from: https://github.com/ionfury/homelab-infrastructure/blob/7847cc352ab553b5bb980c828264bbeba52c5e3a/terraform/inventory.hcl#L15
	// TODO: Reference this directly.
	hosts := map[string]interface{}{
		"node46": map[string]interface{}{
			"cluster": map[string]interface{}{
				"member": clusterName,
				"role":   "controlplane",
			},
			"disk": map[string]interface{}{
				"install": "/dev/sda",
			},
			"lan": []map[string]interface{}{
				{
					"ip":  "192.168.10.246",
					"mac": "ac:1f:6b:2d:c0:22",
				},
			},
			"ipmi": map[string]interface{}{
				"ip":  "192.168.10.231",
				"mac": "ac:1f:6b:68:2b:e1",
			},
		},
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/talos-cluster",
		Vars: map[string]interface{}{
			"name":                   clusterName,
			"endpoint":               endpoint,
			"kubernetes_version":     kubernetes_version,
			"talos_version":          talos_version,
			"talos_config_path":      talos_config_path,
			"kubernetes_config_path": kubernetes_config_path,
			"hosts":                  hosts,
		},
	}

	defer func() {
		resetClusterToMaintenanceMode(t, terraformOptions)
		//terraform.Destroy(t, terraformOptions)
		// Force destroy without refreshing the state.  Needed to destroy data.talos_cluster_health.this
		destroyCmd := shell.Command{
			Command:    "terraform",
			Args:       []string{"destroy", "-auto-approve", "-refresh=false"},
			Env:        terraformOptions.EnvVars,
			WorkingDir: terraformOptions.TerraformDir,
		}

		err := shell.RunCommandE(t, destroyCmd)
		assert.NoError(t, err, "Failed to destroy resources")
	}()

	terraform.InitAndApply(t, terraformOptions)

	confirmClusterWithTalosctl(t, terraformOptions)
	validateKubectlServerVersion(t, terraformOptions)
}

func validateNodes(t *testing.T, terraformOptions *terraform.Options) {
	kubeConfigPath := terraform.Output(t, terraformOptions, "kubernetes_config_file_path")

	kubectlCmd := shell.Command{
		Command: "kubectl",
		Args:    []string{"--kubeconfig", kubeConfigPath, "get", "nodes", "-o", "wide"},
	}

	output, err := shell.RunCommandAndGetOutputE(t, kubectlCmd)
	assert.NoError(t, err, "Failed to run kubectl command")

	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "NAME") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			t.Fatalf("Unexpected format in kubectl output: %s", line)
		}

		nodeName := fields[0]
		nodeStatus := fields[1]
		nodeVersion := fields[4]
		nodeRole := fields[2]
		nodeIP := fields[5]

		hostConfig, ok := hosts[nodeName]
		if !ok {
			t.Fatalf("Node %s not found in the hosts variable", nodeName)
		}

		hostConfigMap, ok := hostConfig.(map[string]interface{})
		if !ok {
			t.Fatalf("host configuration for %s is not a map", nodeName)
		}

		lanConfig, ok := hostConfigMap["lan"].([]interface{})
		if !ok || len(lanConfig) == 0 {
			t.Fatalf("lan configuration for %s is not set or is not a list", nodeName)
		}

		lanMap, ok := lanConfig[0].(map[string]interface{})
		if !ok {
			t.Fatalf("lan configuration for %s is not a map", nodeName)
		}

		ip, ok := lanMap["ip"].(string)
		if !ok || ip == "" {
			t.Fatalf("IP address for %s is not set or is not a string", nodeName)
		}

		// Validate the node details
		assert.Equal(t, ip, nodeIP, "IP address for node %s does not match", nodeName)
		assert.Equal(t, terraformOptions.Vars["kubernetes_version"].(string), nodeVersion, "Node version for %s does not match", nodeName)
		assert.Equal(t, "Ready", nodeStatus, "Node status for %s is not Ready", nodeName)
		assert.Equal(t, hostConfigMap["cluster"].(map[string]interface{})["role"].(string), nodeRole, "Node role for %s does not match", nodeName)
	}
}

func confirmClusterWithTalosctl(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")
	re := regexp.MustCompile(`Server:\s+NODE:\s+\w+\s+Tag:\s+(\S+)`)

	talosctlCmd := shell.Command{
		Command: "talosctl",
		Args:    []string{"--talosconfig", talosConfigFilePath, "version"},
	}

	output, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		t.Fatalf("Server.Tag not found in the talosctl output")
	}

	// serverTag := matches[1]
	// https://github.com/siderolabs/terraform-provider-talos/issues/196#issuecomment-2329652298
	//assert.Equal(t, terraformOptions.Vars["talos_version"].(string), serverTag, "Talos")
	assert.NoError(t, err, "Failed to run talosctl command")
	assert.Contains(t, output, "Server:", "Talos cluster is not up and functional")
}

func validateKubectlServerVersion(t *testing.T, terraformOptions *terraform.Options) {
	kubeConfigPath := terraform.Output(t, terraformOptions, "kubernetes_config_file_path")
	re := regexp.MustCompile(`Server Version:\s+v(\d+\.\d+\.\d+)`)

	kubectlCmd := shell.Command{
		Command: "kubectl",
		Args:    []string{"--kubeconfig", kubeConfigPath, "version"},
	}

	output, err := shell.RunCommandAndGetOutputE(t, kubectlCmd)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		t.Fatalf("Server Version not found in the kubectl output")
	}

	serverVersion := matches[1]

	assert.Equal(t, terraformOptions.Vars["kubernetes_version"].(string), serverVersion, "Cluster Server Version does not match the provided kubernetes version.")
	assert.NoError(t, err, "Failed to run kubectl command")
	assert.Contains(t, output, "Server Version", "Kubernetes cluster is not up and functional")
}

func resetClusterToMaintenanceMode(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigPath := terraform.Output(t, terraformOptions, "talos_config_file_path")

	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	for hostName := range hosts {
		talosctlCmd := shell.Command{
			Command: "talosctl",
			Args:    []string{"--talosconfig", talosConfigPath, "--nodes", hostName, "reset", "--reboot", "--graceful=false"},
		}

		err := shell.RunCommandE(t, talosctlCmd)
		assert.NoError(t, err, "Failed to reset the cluster node %s into maintenance mode", hostName)
		if err != nil {
			t.Fatalf("Failed to reset the cluster node %s into maintenance mode", hostName)
		}
	}
}
