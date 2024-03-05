package git

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/younggwon1/gitops-golang/file"
)

func (gitCli *GitClient) Clone(org, repo string) error {
	exists := file.Exists("/tmp/" + repo)
	if exists {
		err := os.RemoveAll("/tmp/" + repo)
		if err != nil {
			fmt.Println(err)
		}
	}

	r, err := git.PlainClone("/tmp/"+repo, false, &git.CloneOptions{

		Auth:     gitCli.gitAuth,
		URL:      "https://github.com/" + org + "/" + repo,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}

	gitCli.repo = r

	return nil
}

func (gitCli *GitClient) Add(repo string) error {
	directory := "/tmp/" + repo
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

func (gitCli *GitClient) Commit(user, email, file string) error {
	w, err := gitCli.repo.Worktree()
	if err != nil {
		return err
	}

	commit, err := w.Commit("changed image tag in "+file, &git.CommitOptions{
		Author: &object.Signature{
			Name:  user,
			Email: email,
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

func (gitCli *GitClient) Push() error {
	err := gitCli.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       gitCli.gitAuth,
	})
	if err != nil {
		return err
	}

	return nil
}
