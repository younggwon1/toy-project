package kubernetes

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) AllScaleDown(namespace string) ([]string, error) {
	// get the deployments in the namespace
	deployments, err := c.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	c.logger.Info().Msgf("succeed to get deployments %s namespace", namespace)

	// Check the deployment error status
	var scaleDownList []string
	for _, deployment := range deployments.Items {
		if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
			scaleDownList = append(scaleDownList, deployment.Name)
		}
	}
	c.logger.Info().Msg("succeed to check the deployment error status")

	// call scale down function for scale down deployments
	if len(scaleDownList) != 0 {
		err = c.ScaleDown(scaleDownList, namespace)
		if err != nil {
			return nil, err
		}
	} else {
		c.logger.Info().Msgf("failed to find scale down deployments in the %s", namespace)
	}

	return scaleDownList, nil
}

func (c *Client) ScaleDown(names []string, namespace string) error {
	// request confirmation from the user
	fmt.Printf("Are you sure you want to down scale the deployment '%s' in namespace '%s'? (y/n): ", names, namespace)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	response = strings.TrimSpace(response)

	// scale down deployments
	if response == "y" {
		if len(names) != 0 {
			wg := sync.WaitGroup{}
			for _, name := range names {
				wg.Add(1)
				go func(name string) {
					defer wg.Done()
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
						c.logger.Error().Err(err).Msgf("failed to scale down %s deployments", name)
					}
				}(name)
			}
			wg.Wait()
		} else {
			c.logger.Info().Msg("failed to find scale down deployments")
		}
	} else {
		return fmt.Errorf("aborted scale down of deployment")
	}

	return nil
}

func (c *Client) AllDelete(namespace string) ([]string, error) {
	// get the deployments in the namespace
	deployments, err := c.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	c.logger.Info().Msgf("succeed to get deployments %s namespace", namespace)

	// check if the deployment has zero replicas
	var deleteList []string
	for _, deployment := range deployments.Items {
		if *deployment.Spec.Replicas == 0 {
			deleteList = append(deleteList, deployment.Name)
		}
	}
	c.logger.Info().Msg("succeed to check the deployments zero replica")

	// call delete function for delete deployments
	if len(deleteList) != 0 {
		err = c.Delete(deleteList, namespace)
		if err != nil {
			return nil, err
		}
	} else {
		c.logger.Info().Msgf("failed to find delete deployments in the %s", namespace)
	}

	return deleteList, nil
}

func (c *Client) Delete(names []string, namespace string) error {
	// request confirmation from the user
	fmt.Printf("Are you sure you want to delete the deployment '%s' in namespace '%s'? (y/n): ", names, namespace)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	response = strings.TrimSpace(response)

	// delete deployments
	if response == "y" {
		if len(names) != 0 {
			wg := sync.WaitGroup{}
			for _, name := range names {
				wg.Add(1)
				go func(name string) {
					defer wg.Done()
					// Delete the deployment
					err := c.client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
					if err != nil {
						c.logger.Error().Err(err).Msgf("failed to delete %s deployments", name)
					}
				}(name)
			}
			wg.Wait()
		} else {
			c.logger.Info().Msg("failed to find delete deployments")
		}
	} else {
		return fmt.Errorf("aborted deletion of deployment")
	}

	return nil
}
