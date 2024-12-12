package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
)

func fetchHCLFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseHCL(hclContent string) (map[string]interface{}, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "inventory.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	spec := &hcldec.ObjectSpec{
		"raspberry_pis": &hcldec.ObjectSpec{
			"*": &hcldec.ObjectSpec{
				"ssh": &hcldec.ObjectSpec{
					"user_store": &hcldec.AttrSpec{Type: hcldec.String},
					"pass_store": &hcldec.AttrSpec{Type: hcldec.String},
				},
				"lan": &hcldec.ObjectSpec{
					"ip": &hcldec.AttrSpec{Type: hcldec.String},
				},
			},
		},
	}

	result, diags := hcldec.Decode(file.Body, spec, nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode HCL: %s", diags.Error())
	}

	return result.(map[string]interface{}), nil
}

func TestPxePiModule(t *testing.T) {
	t.Parallel()

	// Fetch and parse the HCL file
	hclURL := "https://raw.githubusercontent.com/ionfury/homelab-infrastructure/refs/heads/main/terraform/inventory.hcl"
	hclContent, err := fetchHCLFile(hclURL)
	if err != nil {
		t.Fatalf("Failed to fetch HCL file: %s", err)
	}

	hclData, err := parseHCL(hclContent)
	if err != nil {
		t.Fatalf("Failed to parse HCL file: %s", err)
	}

	raspberryPis := hclData["raspberry_pis"].(map[string]interface{})

	// Define the Terraform options with variables
	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/pxe-pi",
		Vars: map[string]interface{}{
			"raspberry_pi":  "pi1",
			"raspberry_pis": raspberryPis,
		},
	}

	// Clean up resources with "terraform destroy" at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply"
	terraform.InitAndApply(t, terraformOptions)

	// Get the IP address of the provisioned instance
	instanceIP := terraform.Output(t, terraformOptions, "instance_ip")

	// Define the SSH host with password authentication
	sshHost := ssh.Host{
		Hostname:    instanceIP,
		SshUserName: "ubuntu",
		SshPassword: "<your-password>",
	}

	// Test setup_iptables.yaml
	testSetupIptables(t, sshHost)

	// Test setup_ipxe.yaml
	testSetupIpxe(t, sshHost)

	// Test setup_tftp_server.yaml
	testSetupTftpServer(t, sshHost)
}

func testSetupIptables(t *testing.T, sshHost ssh.Host) {
	// Check if iptables rules are set correctly
	iptablesRules := []string{
		"iptables -C INPUT -p udp --dport 69 -s 192.168.10.0/24 -j ACCEPT",
		"iptables -C INPUT -p udp --dport 69 -s 192.168.5.0/24 -j ACCEPT",
		"iptables -C INPUT -p udp --dport 69 -j DROP",
	}

	for _, rule := range iptablesRules {
		output, err := ssh.CheckSshCommandE(t, sshHost, rule)
		assert.NoError(t, err, "Failed to verify iptables rule: %s", rule)
		assert.Empty(t, output, "Unexpected output for iptables rule: %s", rule)
	}
}

func testSetupIpxe(t *testing.T, sshHost ssh.Host) {
	// Check if iPXE files are present
	ipxeFiles := []string{
		"/srv/tftp/ipxe-i386.kpxe",
		"/etc/talospxe/generate_ipxe_menu.sh",
	}

	for _, file := range ipxeFiles {
		output, err := ssh.CheckSshCommandE(t, sshHost, "test -f "+file)
		assert.NoError(t, err, "Failed to verify iPXE file: %s", file)
		assert.Empty(t, output, "Unexpected output for iPXE file: %s", file)
	}
}

func testSetupTftpServer(t *testing.T, sshHost ssh.Host) {
	// Check if TFTP server is running
	output, err := ssh.CheckSshCommandE(t, sshHost, "systemctl is-active tftpd-hpa")
	assert.NoError(t, err, "Failed to verify TFTP server status")
	assert.Equal(t, "active", output, "TFTP server is not active")
}
