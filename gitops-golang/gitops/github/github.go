package github

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/younggwon1/gitops-golang/util"
)

type GithubClient struct {
	logger    zerolog.Logger
	authToken string
}

func NewGithubClient() *GithubClient {
	return &GithubClient{
		logger:    zerolog.New(os.Stdout).With().Timestamp().Logger(),
		authToken: util.GetEnv("PASSWORD", ""),
	}
}
