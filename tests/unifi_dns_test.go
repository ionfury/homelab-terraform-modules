package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestUnifiDNS(t *testing.T) {
	terraformOptions := createUnifiDNSOptions()
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	t.Run("group", func(t *testing.T) {
		t.Run("validateUnifiDNS", func(t *testing.T) {
			t.Parallel()
			validateUnifiDNS(t, terraformOptions)
		})
	})
}

func createUnifiDNSOptions() *terraform.Options {
	return &terraform.Options{
		TerraformDir: "../modules/unifi-dns",
		Vars: map[string]interface{}{
			"unifi_dns_records": map[string]interface{}{
				"record": map[string]interface{}{
					"name":  random.UniqueId() + ".com",
					"value": "192.168.111.111",
				},
			},
		},
	}
}

func validateUnifiDNS(t *testing.T, terraformOptions *terraform.Options) {
	records := make(map[string]string)
	for _, value := range terraformOptions.Vars["unifi_dns_records"].(map[string]interface{}) {
		record := value.(map[string]interface{})
		records[record["name"].(string)] = record["value"].(string)
	}
	time.Sleep(10 * time.Second)
	for expectedValue, ip := range records {
		cmd := shell.Command{
			Command: "nslookup",
			Args:    []string{expectedValue},
		}
		output, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			t.Fatalf("Failed to run nslookup: %v", err)
		}

		assert.Contains(t, string(output), ip, "Expected IP %s for %s, but got %s", ip, expectedValue, string(output))
	}
}
