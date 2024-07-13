package argocd

import (
	"context"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

type Connection struct {
	Address string
	Token   string
}

type Client struct {
	client apiclient.Client
}

type AppClient struct {
	appClient application.ApplicationServiceClient
}

type AppSyncRequest struct {
	Name         *string
	DryRun       *bool
	Prune        *bool
	SyncStrategy *AppSyncStrategyRequest
}

type AppSyncStrategyRequest struct {
	Force bool
}

func NewClient(c *Connection) (*Client, error) {
	// set a new argocd client
	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: c.Address,
		AuthToken:  c.Token,
		Insecure:   true,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) NewAppClient() (*AppClient, error) {
	// set a new argocd app client
	_, appClient, err := c.client.NewApplicationClient()
	if err != nil {
		return nil, err
	}

	return &AppClient{
		appClient: appClient,
	}, nil
}

func (ac *AppClient) ExistsArgoCDAppCheck(r *string) error {
	// check if argocd app is null
	if r == nil || *r == "" {
		return fmt.Errorf("failed because of argocd app name is set to an empty value")
	}

	// check if a specific argocd app exists
	_, err := ac.appClient.Get(context.Background(), &application.ApplicationQuery{
		Name: r,
	})
	if err != nil {
		return fmt.Errorf("failed to get argocd app: %s, because does not exist app", *r)
	}

	return nil
}

func (ac *AppClient) Sync(s *AppSyncRequest) error {
	// check if argocd app exists
	err := ac.ExistsArgoCDAppCheck(s.Name)
	if err != nil {
		return err
	}

	// sync a specific argocd app
	_, err = ac.appClient.Sync(context.Background(), &application.ApplicationSyncRequest{
		Name:   s.Name,
		Prune:  s.Prune,
		DryRun: s.DryRun,
		Strategy: &v1alpha1.SyncStrategy{
			Apply: &v1alpha1.SyncStrategyApply{
				Force: s.SyncStrategy.Force,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
