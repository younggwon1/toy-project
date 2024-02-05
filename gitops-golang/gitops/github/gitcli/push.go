package gitcli

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Push(repo *git.Repository) error {
	err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: os.Getenv("Username"),
			Password: os.Getenv("AccessToken"),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
