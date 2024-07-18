package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
)

type AWSClient struct {
	ctx           context.Context
	cfg           config.Config
	amplifyClient *amplify.Client
}

type AmplifyClient struct {
	AmplifyClient *amplify.Client
}

func NewAWSClient(ctx context.Context) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &AWSClient{
		ctx:           ctx,
		cfg:           cfg,
		amplifyClient: amplify.NewFromConfig(cfg),
	}, nil
}

func (c *AWSClient) NewAmplifyClient() *AmplifyClient {
	return &AmplifyClient{
		AmplifyClient: c.amplifyClient,
	}
}
