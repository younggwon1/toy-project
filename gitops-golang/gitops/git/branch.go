package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/uuid"
)

func (gitCli *GitClient) Create() (plumbing.ReferenceName, error) {
	headRef, err := gitCli.repo.Head()
	gitCli.logger.Info().Msgf("headRef.Hash() : %s", headRef.Hash().String())
	if err != nil {
		return "", err
	}

	gitCli.branch = plumbing.NewBranchReferenceName("socops-deploy-" + uuid.New().String()[:8])
	gitCli.branchRef = plumbing.NewHashReference(gitCli.branch, headRef.Hash())

	if err := gitCli.repo.Storer.SetReference(gitCli.branchRef); err != nil {
		return "", err
	}

	return gitCli.branch, nil
}

func (gitCli *GitClient) Checkout() error {
	worktree, err := gitCli.repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(gitCli.branch),
	})

	if err != nil {
		return err
	}
	return nil
}

func (gitCli *GitClient) Delete(userName, accessToken string) error {
	// Delete remote branch
	gitCli.logger.Info().Msg("Start remote branch delete")
	pushOpts := &git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(":" + gitCli.branch)},
		Auth: &http.BasicAuth{
			Username: userName,
			Password: accessToken,
		},
	}

	err := gitCli.repo.Push(pushOpts)
	if err != nil {
		return err
	}
	gitCli.logger.Info().Msg("End remote branch delete")

	return nil
}
