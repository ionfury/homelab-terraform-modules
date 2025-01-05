package test

import (
	"context"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func validateHostnameUniqueness(t *testing.T, terraformOptions *terraform.Options) {
	kubeConfigPath := terraform.Output(t, terraformOptions, "kubernetes_config_file_path")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	assert.NoError(t, err, "Failed to build Kubernetes config")

	clientset, err := kubernetes.NewForConfig(config)
	assert.NoError(t, err, "Failed to create Kubernetes client")

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	assert.NoError(t, err, "kubectl get nodes failed - check KUBECONFIG setup")

	assert.NotEmpty(t, nodes.Items, "kubectl get nodes returned empty list - check KUBECONFIG setup")

	deduplicateHostnames := make(map[string]bool)
	numNodes := 0

	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == "Hostname" {
				numNodes++
				if _, exists := deduplicateHostnames[address.Address]; !exists {
					deduplicateHostnames[address.Address] = true
				}
			}
		}
	}

	assert.Equal(t, len(deduplicateHostnames), numNodes, "Nodes do not have unique hostnames")
}
