package deployment

import (
	"context"
	"fmt"
	"sync"

	"github.com/younggwon1/k8s-resource-manager/deployment/config/response"
	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ErrorStatus(clientSet *kubernetes.Clientset, nameSpace string) error {
	deployments, err := clientSet.AppsV1().Deployments(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errorDeploymentsList []response.ResponseErrorStatus
	for _, deployment := range deployments.Items {
		// Check the deployment error status
		if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
			errorDeployment := &response.ResponseErrorStatus{
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
		wg := sync.WaitGroup{}
		for _, errorDeployment := range errorDeploymentsList {
			wg.Add(1)
			go func(errorDeployment response.ResponseErrorStatus) {
				defer wg.Done()
				fmt.Println(errorDeployment.Name + " deployments has error status")
				_, err := ScaleDownErrorStatus(clientSet, nameSpace, errorDeployment.Name)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Scale down " + errorDeployment.Name + " deployments")
			}(errorDeployment)
		}
		wg.Wait()
	} else {
		fmt.Println("No error deployments in the " + nameSpace + " namespace")
	}

	return nil
}

func ScaleDownErrorStatus(clientSet *kubernetes.Clientset, nameSpace, name string) (*v1.Scale, error) {
	scale := &v1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: nameSpace,
		},
		Spec: v1.ScaleSpec{
			Replicas: 0,
		},
	}
	scale, err := clientSet.AppsV1().Deployments(nameSpace).UpdateScale(context.TODO(), name, scale, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return scale, nil
}
