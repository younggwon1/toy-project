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

func (ac *AppClient) ExistsCheck(r *AppSyncRequest) error {
	// check if argocd app is null
	if r.Name == nil {
		return fmt.Errorf("failed if argocd app name is null")
	}

	// check if a specific argocd app exists
	_, err := ac.appClient.Get(context.Background(), &application.ApplicationQuery{
		Name: r.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (ac *AppClient) Sync(s *AppSyncRequest) error {
	// check if argocd app exists
	err := ac.ExistsCheck(s)
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
