package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Push(repo *git.Repository, UserName, AccessToken string) error {
	err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: UserName,
			Password: AccessToken,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
