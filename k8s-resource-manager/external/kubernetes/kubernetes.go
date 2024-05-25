package kubernetes

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/younggwon1/k8s-resource-manager/deployment/config/response"
)

type Client struct {
	client *kubernetes.Clientset
	logger *zerolog.Logger
}

func NewClient(logger *zerolog.Logger) (*Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

func (c *Client) ScaleDownErrorDeployment(namespace string) error {
	// get the deployments in the namespace
	deployments, err := c.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	c.logger.Info().Msg("succeed to get deployments")

	var errorList []response.ErrorStatus
	for _, deployment := range deployments.Items {
		// Check the deployment error status
		if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
			errorDeployment := &response.ErrorStatus{
				Name:      deployment.Name,
				NameSpace: deployment.Namespace,
				Reason:    string(deployment.Status.Conditions[1].Reason),
				Message:   deployment.Status.Conditions[1].Message,
				Age:       deployment.ObjectMeta.GetCreationTimestamp().Format("2006-01-02 15:04:05"),
			}
			errorList = append(errorList, *errorDeployment)
		}
	}
	c.logger.Info().Msg("succeed to check the deployment error status")

	if len(errorList) != 0 {
		wg := sync.WaitGroup{}
		for _, errorDeployment := range errorList {
			wg.Add(1)
			go func(errorDeployment response.ErrorStatus) {
				defer wg.Done()
				err := c.ScaleDown(errorDeployment.Name, namespace)
				if err != nil {
					c.logger.Error().Err(err).Msg("failed to scale down " + errorDeployment.Name + " deployments")
				}
				c.logger.Info().Msg("succeed to scale down " + errorDeployment.Name + " deployments")
			}(errorDeployment)
		}
		wg.Wait()
	} else {
		c.logger.Info().Msg("No error deployments in the " + namespace + " namespace")
	}

	return nil
}

func (c *Client) ScaleDown(name, namespace string) error {
	// Scale down the deployment
	scale := &v1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.ScaleSpec{
			Replicas: 0,
		},
	}
	_, err := c.client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), name, scale, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
