package test

import (
	"sync"
	"testing"

	"github.com/go-viper/mapstructure/v2"
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

type Install struct {
	DiskSelector string `mapstructure:"diskSelector"`
}

type Cluster struct {
	Member string `mapstructure:"member"`
	Role   string `mapstructure:"role"`
}

type Host struct {
	Cluster    Cluster     `mapstructure:"cluster"`
	Install    Install     `mapstructure:"install"`
	Interfaces []Interface `mapstructure:"interfaces"`
	IPMI       IPMI        `mapstructure:"ipmi"`
}

type Hosts map[string]Host

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

/*
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
*/
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

	var wg sync.WaitGroup
	for hostName := range hosts {
		wg.Add(1)
		go func(hostName string) {
			defer wg.Done()
			talosctlCmd := shell.Command{
				Command: "talosctl",
				Args:    []string{"--talosconfig", talosConfigPath, "--nodes", hostName, "reset", "--reboot", "--graceful=false", "--wait=false"},
			}

			err := shell.RunCommandE(t, talosctlCmd)
			assert.NoError(t, err, "Failed to reset the cluster node %s into maintenance mode", hostName)
		}(hostName)
	}
	wg.Wait()
}

func validateTalosNodeIpConfig(t *testing.T, terraformOptions *terraform.Options) {}

func validateTalosLinkConfig(t *testing.T, terraformOptions *terraform.Options) {}

func validateTalosRouteConfig(t *testing.T, terraformOptions *terraform.Options) {}
