package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog"
)

type GitClient struct {
	repo      *git.Repository
	branch    plumbing.ReferenceName
	branchRef *plumbing.Reference
	logger    zerolog.Logger
}

func NewGitClient(repo *git.Repository) *GitClient {
	return &GitClient{
		repo:      repo,
		branch:    plumbing.ReferenceName(""),
		branchRef: &plumbing.Reference{},
		logger:    zerolog.New(os.Stdout),
	}
}
