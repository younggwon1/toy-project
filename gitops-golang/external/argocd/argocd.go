package argocd

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/rs/zerolog"
)

type Connection struct {
	Address string
	Token   string
}

type Client struct {
	client apiclient.Client
	logger zerolog.Logger
}

type AppClient struct {
	appClient application.ApplicationServiceClient
}

func NewClient(c *Connection, logger zerolog.Logger) (*Client, error) {
	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: c.Address,
		Insecure:   true,
		AuthToken:  c.Token,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

func (c *Client) NewAppClient() (*AppClient, error) {
	_, appClient, err := c.client.NewApplicationClient()
	if err != nil {
		return nil, err
	}

	return &AppClient{
		appClient: appClient,
	}, nil
}

func (ac *AppClient) Sync() error {
	_, err := ac.appClient.Sync(context.Background(), &application.ApplicationSyncRequest{})
	if err != nil {
		return err
	}

	return nil
}
