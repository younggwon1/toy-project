package kubernetes

import (
	"github.com/rs/zerolog"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
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
