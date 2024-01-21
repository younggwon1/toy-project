package main

import (
	k8s "github.com/younggwon1/k8s-resource-manager/config/kubernetes"
	deployment "github.com/younggwon1/k8s-resource-manager/deployment"
)

func main() {
	kubernetesConfig := k8s.KubernetesCredentials() // Init kubernetes config

	deployment.ErrorDeployments(kubernetesConfig)
	deployment.DeleteErrorDeployments(kubernetesConfig)

}
