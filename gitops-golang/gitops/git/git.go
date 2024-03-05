package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog"

	"github.com/younggwon1/gitops-golang/util"
)

type GitClient struct {
	logger    zerolog.Logger
	repo      *git.Repository
	Branch    plumbing.ReferenceName
	branchRef *plumbing.Reference
	gitAuth   *http.BasicAuth
}

func NewGitClient() *GitClient {
	return &GitClient{
		logger:    zerolog.New(os.Stdout).With().Timestamp().Logger(),
		repo:      &git.Repository{},
		Branch:    plumbing.ReferenceName(""),
		branchRef: &plumbing.Reference{},
		gitAuth: &http.BasicAuth{
			Username: util.GetEnv("USERNAME", ""),
			Password: util.GetEnv("PASSWORD", ""),
		},
	}
}
