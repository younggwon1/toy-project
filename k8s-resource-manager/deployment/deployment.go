package deployment

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// request "github.com/younggwon1/k8s-resource-manager/deployment/config/request"
	response "github.com/younggwon1/k8s-resource-manager/deployment/config/response"
)

func ErrorDeployments(clientSet *kubernetes.Clientset) {
	deployments, err := clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}

	var errorDeploymentsList []response.ResponseErrorDeployments
	for _, deployment := range deployments.Items {
		// Check the deployment error status
		if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
			errorDeployment := &response.ResponseErrorDeployments{
				Name:      deployment.Name,
				NameSpace: deployment.Namespace,
				Reason:    string(deployment.Status.Conditions[1].Reason),
				Message:   deployment.Status.Conditions[1].Message,
				Age:       deployment.ObjectMeta.GetCreationTimestamp().Format("2006-01-02 15:04:05"),
			}
			errorDeploymentsList = append(errorDeploymentsList, *errorDeployment)
		}
	}

	if len(errorDeploymentsList) != 0 {
		errorDeploymentsJson, _ := json.Marshal(errorDeploymentsList)
		fmt.Println(errorDeploymentsJson)
	} else {
		fmt.Println("No error deployments")
	}
}

func DeleteErrorDeployments(clientSet *kubernetes.Clientset) {
	fmt.Println("Delete Error Deployments")
}
