package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func Add(repo *git.Repository, helmRepo string) error {
	directory := "/tmp/" + helmRepo
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		fmt.Println("Nothing to commit, Working tree clean")
		return nil
	}

	fmt.Println(status)

	err = w.AddWithOptions(&git.AddOptions{
		All:  true,
		Path: directory,
	})
	if err != nil {
		return err
	}

	return nil
}
