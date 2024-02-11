package branch

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Checkout(repo *git.Repository, branch plumbing.ReferenceName) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branch),
	})

	if err != nil {
		return err
	}
	return nil
}