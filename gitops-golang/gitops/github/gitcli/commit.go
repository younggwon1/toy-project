package gitcli

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Commit(repo *git.Repository, gitUser string, gitUserEmail string, yamlFile string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	commit, err := w.Commit("changed image tag in "+yamlFile, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitUser,
			Email: gitUserEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		return err
	}
	fmt.Println(obj)

	return nil
}
