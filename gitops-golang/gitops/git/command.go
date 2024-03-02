package git

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/younggwon1/gitops-golang/file"
)

func Clone(userName, accessToken, org, helmRepo string) (*git.Repository, error) {
	exists := file.Exists("/tmp/" + helmRepo)
	if exists {
		err := os.RemoveAll("/tmp/" + helmRepo)
		if err != nil {
			fmt.Println(err)
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

func (gitCli *GitClient) Add(helmRepo string) error {
	directory := "/tmp/" + helmRepo
	w, err := gitCli.repo.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		gitCli.logger.Info().Msg("Nothing to commit, Working tree clean")
		return nil
	}

	gitCli.logger.Info().Msg(status.String())

	err = w.AddWithOptions(&git.AddOptions{
		All:  true,
		Path: directory,
	})
	if err != nil {
		return err
	}

	return nil
}

func (gitCli *GitClient) Commit(gitUser, gitUserEmail, yamlFile string) error {
	w, err := gitCli.repo.Worktree()
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

	obj, err := gitCli.repo.CommitObject(commit)
	if err != nil {
		return err
	}
	gitCli.logger.Info().Msg(obj.String())

	return nil
}

func (gitCli *GitClient) Push(userName, accessToken string) error {
	err := gitCli.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: userName,
			Password: accessToken,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
