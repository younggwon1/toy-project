package argocd

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

	"github.com/younggwon1/gitops-golang/config"
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

func (ac *AppClient) Sync(s *AppSyncRequest, executor, gitUrl, tag, ticket, argoCDAppUrl string) (*config.Auditor, error) {
	// check if argocd app exists
	err := ac.ExistsArgoCDAppCheck(s.Name)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// init ticker
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// check argocd app status per each tick after sync
	loop := true
	var audit *config.Auditor
	for loop {
		select {
		case <-context.Background().Done():
			loop = false
		case <-ticker.C:
			// get argocd app status
			response, err := ac.appClient.Get(context.Background(), &application.ApplicationQuery{
				Name: s.Name,
			})
			if err != nil {
				return nil, err
			}
			// set audit data
			if response.Status.OperationState != nil {
				audit = &config.Auditor{
					Version: config.Version{
						Version: "v0.0.1",
					},
					Metadata: config.Metadata{
						Name: *s.Name,
						Label: config.Label{
							Executor: executor,
						},
					},
					Spec: config.Spec{
						Source: config.Source{
							Code: config.Code{
								Repo: gitUrl,
								Rev:  tag,
							},
							Helm: config.Helm{
								Repo:  response.Status.OperationState.SyncResult.Source.RepoURL,
								Chart: path.Join(response.Status.OperationState.SyncResult.Source.Path, response.Status.OperationState.SyncResult.Source.Helm.ValueFiles[0]),
								Rev:   response.Status.OperationState.SyncResult.Revision,
							},
							Jira: config.Jira{
								Ticket: config.Ticket{
									CR: ticket,
								},
							},
						},
						Destination: config.Destination{
							ArgoCD: config.ArgoCD{
								URL:    argoCDAppUrl,
								Synced: response.Status.OperationState.FinishedAt.Time.Format("2006-01-02 15:04:05"),
							},
						},
					},
				}
				loop = false
			}
		}
	}

	return audit, nil
}
