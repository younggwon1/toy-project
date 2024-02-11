package git

import (
	"os"

	"github.com/younggwon1/gitops-golang/file"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Clone(userName string, accessToken string, org string, helmRepo string) (*git.Repository, error) {
	exists := file.Exists("/tmp/" + helmRepo)
	if exists {
		err := os.RemoveAll("/tmp/" + helmRepo)
		if err != nil {
			return nil, err
		}
	}
	repo, err := git.PlainClone("/tmp/"+helmRepo, false, &git.CloneOptions{

		Auth: &http.BasicAuth{
			Username: userName,
			Password: accessToken,
		},

		URL:      "https://github.com/" + org + "/" + helmRepo,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}