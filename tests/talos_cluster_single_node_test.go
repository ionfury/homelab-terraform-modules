package test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

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
		/*
			t.Run("validateTalosInstallDiskConfig", func(t *testing.T) {
				t.Parallel()
				validateTalosInstallDiskConfig(t, terraformOptions)
			})
		*/
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

		t.Run("validateHostnameUniqueness", func(t *testing.T) {
			t.Parallel()
			validateHostnameUniqueness(t, terraformOptions)
		})

		/*t.Run("validateMountPropagation", func(t *testing.T) {
			t.Parallel()
			waitLonghornEnvironmentCheckReady(t, terraformOptions)
			validateMountPropagation(t, terraformOptions)
		})*/
	})
}

func createTalosClusterSingleNodeOptions() *terraform.Options {
	clusterName := strings.ToLower("talos-cluster-single-node-" + random.UniqueId())
	endpoint := "https://192.168.10.246:6443"
	kubernetes_version := "1.30.1"
	talos_version := "v1.8.4"
	talos_config_path := "~/.talos"
	kubernetes_config_path := "~/.kube"
	nameservers := []string{"8.8.8.8", "1.1.1.1"}
	ntp_servers := []string{"0.pool.ntp.org", "1.pool.ntp.org"}
	cluster_vip := "192.168.10.5"
	allow_scheduling_on_controlplane := true
	host_dns_enabled := true
	host_dns_resolveMemberNames := true
	host_dns_forwardKubeDNSToHost := true
	// Sourced from: https://github.com/ionfury/homelab-infrastructure/blob/7847cc352ab553b5bb980c828264bbeba52c5e3a/terraform/inventory.hcl#L15
	// TODO: Reference this directly.
	hosts := map[string]interface{}{
		"node46": map[string]interface{}{
			"role": "controlplane",
			"install": map[string]interface{}{
				"diskSelector":    []string{"type: 'ssd'"},
				"extraKernelArgs": []string{"apparmor=0"},
				"extensions":      []string{"iscsi-tools", "util-linux-tools"},
				"secureboot":      false,
				"wipe":            false,
				"architecture":    "amd64",
				"platform":        "metal",
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
		},
	}

	return &terraform.Options{
		TerraformDir: "../modules/talos-cluster",
		Vars: map[string]interface{}{
			"name":                             clusterName,
			"endpoint":                         endpoint,
			"kubernetes_version":               kubernetes_version,
			"talos_version":                    talos_version,
			"talos_config_path":                talos_config_path,
			"kubernetes_config_path":           kubernetes_config_path,
			"nameservers":                      nameservers,
			"ntp_servers":                      ntp_servers,
			"cluster_vip":                      cluster_vip,
			"host_dns_enabled":                 host_dns_enabled,
			"host_dns_resolveMemberNames":      host_dns_resolveMemberNames,
			"host_dns_forwardKubeDNSToHost":    host_dns_forwardKubeDNSToHost,
			"allow_scheduling_on_controlplane": allow_scheduling_on_controlplane,
			"hosts":                            hosts,
		},
	}
}
