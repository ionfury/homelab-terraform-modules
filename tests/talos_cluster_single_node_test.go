package test

import (
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

type VLAN struct {
	VlanID          int      `mapstructure:"vlanId"`
	Addresses       []string `mapstructure:"addresses"`
	DHCPRouteMetric int      `mapstructure:"dhcp_routeMetric"`
}

type Interface struct {
	HardwareAddr    string   `mapstructure:"hardwareAddr"`
	Addresses       []string `mapstructure:"addresses"`
	DHCPRouteMetric int      `mapstructure:"dhcp_routeMetric"`
	VLANs           []VLAN   `mapstructure:"vlans"`
}

type IPMI struct {
	IP  string `mapstructure:"ip"`
	MAC string `mapstructure:"mac"`
}

type Disk struct {
	Install string `mapstructure:"install"`
}

type Cluster struct {
	Member string `mapstructure:"member"`
	Role   string `mapstructure:"role"`
}

type Host struct {
	Cluster    Cluster     `mapstructure:"cluster"`
	Disk       Disk        `mapstructure:"disk"`
	Interfaces []Interface `mapstructure:"interfaces"`
	IPMI       IPMI        `mapstructure:"ipmi"`
}

type Hosts map[string]Host

func TestTalosClusterSingleNode(t *testing.T) {
	terraformOptions := createTalosClusterSingleNodeOptions()
	defer func() {
		resetClusterToMaintenanceMode(t, terraformOptions)
		destroyTerraformState(t, terraformOptions)
	}()
	terraform.InitAndApply(t, terraformOptions)

	t.Run("group", func(t *testing.T) {
		t.Run("validateTalosHostDnsConfig", func(t *testing.T) {
			t.Parallel()
			validateTalosHostDnsConfig(t, terraformOptions)
		})

		t.Run("validateTalosHostnameConfig", func(t *testing.T) {
			t.Parallel()
			validateTalosHostnameConfig(t, terraformOptions)
		})

		t.Run("validateTalosInstallDiskConfig", func(t *testing.T) {
			t.Parallel()
			validateTalosInstallDiskConfig(t, terraformOptions)
		})

		t.Run("validateTalosNameserversConfig", func(t *testing.T) {
			t.Parallel()
			validateTalosNameserversConfig(t, terraformOptions)
		})

		t.Run("validateTalosNTPServersConfig", func(t *testing.T) {
			t.Parallel()
			validateTalosNTPServersConfig(t, terraformOptions)
		})

		t.Run("validateTalosControlPlaneSchedulingConfig", func(t *testing.T) {
			t.Parallel()
			validateTalosControlPlaneSchedulingConfig(t, terraformOptions)
		})

		t.Run("validateKubernetesVersionConfig", func(t *testing.T) {
			t.Parallel()
			validateKubernetesVersionConfig(t, terraformOptions)
		})
	})
}

func destroyTerraformState(t *testing.T, terraformOptions *terraform.Options) {
	destroyCmd := shell.Command{
		Command:    "terraform",
		Args:       []string{"destroy", "-auto-approve", "-refresh=false"},
		Env:        terraformOptions.EnvVars,
		WorkingDir: terraformOptions.TerraformDir,
	}

	err := shell.RunCommandE(t, destroyCmd)
	assert.NoError(t, err, "Failed to destroy resources")
}

func createTalosClusterSingleNodeOptions() *terraform.Options {
	clusterName := "talos-cluster-" + random.UniqueId()
	endpoint := "https://192.168.10.246:6443"
	kubernetes_version := "1.30.1"
	talos_version := "v1.8.4"
	talos_config_path := "~/.talos"
	kubernetes_config_path := "~/.kube"
	nameservers := []string{"192.168.10.1"}
	ntp_servers := []string{"0.pool.ntp.org", "1.pool.ntp.org"}
	cluster_vip := "192.168.10.5"
	//ingress_firewall_enabled := true
	//cluster_subnet := "192.168.10.0/24"
	//cni_vxlan_port := "8473"
	allow_scheduling_on_controlplane := true
	host_dns_enabled := true
	host_dns_resolveMemberNames := true
	host_dns_forwardKubeDNSToHost := true
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
			"interfaces": []map[string]interface{}{
				{
					"hardwareAddr":     "ac:1f:6b:2d:c0:22",
					"addresses":        []string{"192.168.10.246"},
					"dhcp_routeMetric": 50,
					"vlans": []map[string]interface{}{
						{
							"vlanId":           20,
							"addresses":        []string{"192.168.20.20"},
							"dhcp_routeMetric": 100,
						},
					},
				},
			},
			"ipmi": map[string]interface{}{
				"ip":  "192.168.10.231",
				"mac": "ac:1f:6b:68:2b:e1",
			},
		},
	}

	return &terraform.Options{
		TerraformDir: "../modules/talos-cluster",
		Vars: map[string]interface{}{
			"name":                          clusterName,
			"endpoint":                      endpoint,
			"kubernetes_version":            kubernetes_version,
			"talos_version":                 talos_version,
			"talos_config_path":             talos_config_path,
			"kubernetes_config_path":        kubernetes_config_path,
			"nameservers":                   nameservers,
			"ntp_servers":                   ntp_servers,
			"cluster_vip":                   cluster_vip,
			"host_dns_enabled":              host_dns_enabled,
			"host_dns_resolveMemberNames":   host_dns_resolveMemberNames,
			"host_dns_forwardKubeDNSToHost": host_dns_forwardKubeDNSToHost,
			//"ingress_firewall_enabled":         ingress_firewall_enabled,
			//"cluster_subnet":                   cluster_subnet,
			"allow_scheduling_on_controlplane": allow_scheduling_on_controlplane,
			"hosts":                            hosts,
		},
	}
}

func validateTalosHostDnsConfig(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")

	hostDnsEnabled, ok := terraformOptions.Vars["host_dns_enabled"].(bool)
	if !ok {
		t.Fatalf("host_dns_enabled variable is not set or is not a boolean")
	}

	if hostDnsEnabled {
		hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
		if !ok {
			t.Fatalf("hosts variable is not set or is not a map")
		}

		resolveMemberNames, ok := terraformOptions.Vars["host_dns_resolveMemberNames"].(bool)
		if !ok {
			t.Fatalf("host_dns_resolveMemberNames variable is not set or is not a boolean")
		}

		for hostName := range hosts {
			talosctlCmd := shell.Command{
				Command: "talosctl",
				Args:    []string{"--talosconfig", talosConfigFilePath, "-n", hostName, "get", "hostdnsconfig", "-o", "json"},
			}

			json, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
			if err != nil {
				t.Fatalf("Failed to run talosctl command: %v", err)
			}

			assert.Equal(t, hostDnsEnabled, gjson.Get(json, "spec.enabled").Bool(), "Host DNS for host %s is not enabled", hostName)
			assert.Equal(t, resolveMemberNames, gjson.Get(json, "spec.resolveMemberNames").Bool(), "ResolveMemberNames for host %s does not match", hostName)
		}
	}
}

func validateTalosHostnameConfig(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")

	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	for hostName := range hosts {
		talosctlCmd := shell.Command{
			Command: "talosctl",
			Args:    []string{"--talosconfig", talosConfigFilePath, "-n", hostName, "get", "hostname", "-o", "json"},
		}

		json, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
		if err != nil {
			t.Fatalf("Failed to run talosctl command: %v", err)
		}

		assert.Equal(t, hostName, gjson.Get(json, "spec.hostname").String(), "Hostname for host %s does not match", hostName)
	}
}

func validateTalosInstallDiskConfig(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")
	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	var decodedHosts Hosts
	err := mapstructure.Decode(hosts, &decodedHosts)
	if err != nil {
		t.Fatalf("Failed to decode hosts: %v", err)
	}

	for hostName, hostConfig := range decodedHosts {
		talosctlCmd := shell.Command{
			Command: "talosctl",
			Args:    []string{"--talosconfig", talosConfigFilePath, "-n", hostName, "get", "systemdisk", "-o", "json"},
		}

		json, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
		if err != nil {
			t.Fatalf("Failed to run talosctl command: %v", err)
		}

		assert.Equal(t, hostConfig.Disk.Install, gjson.Get(json, "spec.devPath").String(), "Install disk for host %s does not match", hostName)
	}
}

func validateTalosInterfaceHardwareAddrConfig(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")
	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	var decodedHosts Hosts
	err := mapstructure.Decode(hosts, &decodedHosts)
	if err != nil {
		t.Fatalf("Failed to decode hosts: %v", err)
	}

	for hostName, hostConfig := range decodedHosts {
		for _, iface := range hostConfig.Interfaces {
			talosctlCmd := shell.Command{
				Command: "talosctl",
				Args:    []string{"--talosconfig", talosConfigFilePath, "-n", hostName, "get", "link", "-o", "json"},
			}

			json, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
			if err != nil {
				t.Fatalf("Failed to run talosctl command: %v", err)
			}

			assert.Equal(t, iface.HardwareAddr, gjson.Get(json, "spec.hardwareAddr").String(), "Hardware address for host %s does not match", hostName)
		}
	}
}

func validateTalosNodeIpConfig(t *testing.T, terraformOptions *terraform.Options) {}

func validateTalosLinkConfig(t *testing.T, terraformOptions *terraform.Options) {}

func validateTalosRouteConfig(t *testing.T, terraformOptions *terraform.Options) {}

func validateTalosNameserversConfig(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")
	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	expectedNameservers, ok := terraformOptions.Vars["nameservers"].([]string)
	if !ok {
		t.Fatalf("nameservers variable is not set or is not a list of strings")
	}

	for hostName := range hosts {
		talosctlCmd := shell.Command{
			Command: "talosctl",
			Args:    []string{"--talosconfig", talosConfigFilePath, "-n", hostName, "get", "resolvers", "-o", "json"},
		}

		json, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
		if err != nil {
			t.Fatalf("Failed to run talosctl command: %v", err)
		}

		nameservers := gjson.Get(json, "spec.dnsServers").Array()
		var actualNameservers []string
		for _, nameserver := range nameservers {
			actualNameservers = append(actualNameservers, nameserver.String())
		}

		assert.Equal(t, expectedNameservers, actualNameservers, "Configured nameservers for host %s do not match", hostName)
	}
}

func validateTalosNTPServersConfig(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")
	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	expectedNTPServers, ok := terraformOptions.Vars["ntp_servers"].([]string)
	if !ok {
		t.Fatalf("ntp_servers variable is not set or is not a list of strings")
	}

	for hostName := range hosts {
		talosctlCmd := shell.Command{
			Command: "talosctl",
			Args:    []string{"--talosconfig", talosConfigFilePath, "-n", hostName, "get", "timeservers", "-o", "json"},
		}

		json, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
		if err != nil {
			t.Fatalf("Failed to run talosctl command: %v", err)
		}

		timeservers := gjson.Get(json, "spec.timeServers").Array()
		var actualNTPServers []string
		for _, timeserver := range timeservers {
			actualNTPServers = append(actualNTPServers, timeserver.String())
		}

		assert.Equal(t, expectedNTPServers, actualNTPServers, "Configured NTP servers for host %s do not match", hostName)
	}
}

func validateTalosControlPlaneSchedulingConfig(t *testing.T, terraformOptions *terraform.Options) {
	talosConfigFilePath := terraform.Output(t, terraformOptions, "talos_config_file_path")
	hosts, ok := terraformOptions.Vars["hosts"].(map[string]interface{})
	if !ok {
		t.Fatalf("hosts variable is not set or is not a map")
	}

	allowSchedulingOnControlplane, ok := terraformOptions.Vars["allow_scheduling_on_controlplane"].(bool)
	if !ok {
		t.Fatalf("allow_scheduling_on_controlplane variable is not set or is not a boolean")
	}

	for hostName := range hosts {
		talosctlCmd := shell.Command{
			Command: "talosctl",
			Args:    []string{"--talosconfig", talosConfigFilePath, "-n", hostName, "get", "schedulerconfig", "-o", "json"},
		}

		json, err := shell.RunCommandAndGetOutputE(t, talosctlCmd)
		if err != nil {
			t.Fatalf("Failed to run talosctl command: %v", err)
		}

		assert.Equal(t, allowSchedulingOnControlplane, gjson.Get(json, "spec.enabled").Bool(), "AllowScheduling for host %s does not match", hostName)
	}
}

func validateKubernetesVersionConfig(t *testing.T, terraformOptions *terraform.Options) {
	kubeConfigPath := terraform.Output(t, terraformOptions, "kubernetes_config_file_path")

	kubectlCmd := shell.Command{
		Command: "kubectl",
		Args:    []string{"--kubeconfig", kubeConfigPath, "version", "-o", "json"},
	}

	output, err := shell.RunCommandAndGetOutputE(t, kubectlCmd)
	if err != nil {
		t.Fatalf("Failed to run kubectl command: %v", err)
	}

	expectedVersion := "v" + terraformOptions.Vars["kubernetes_version"].(string)
	assert.Equal(t, expectedVersion, gjson.Get(output, "serverVersion.gitVersion").String(), "Kubernetes version does not match the provided version")
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
