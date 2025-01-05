package test

import (
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func validateKubernetesVersionConfig(t *testing.T, terraformOptions *terraform.Options) {
	kubeConfigPath := terraform.Output(t, terraformOptions, "kubernetes_config_file_path")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		t.Fatalf("Failed to build Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	versionInfo, err := clientset.Discovery().ServerVersion()
	if err != nil {
		t.Fatalf("Failed to get Kubernetes version: %v", err)
	}

	expectedVersion := "v" + terraformOptions.Vars["kubernetes_version"].(string)
	assert.Equal(t, expectedVersion, versionInfo.GitVersion, "Kubernetes version does not match the provided version")
}
