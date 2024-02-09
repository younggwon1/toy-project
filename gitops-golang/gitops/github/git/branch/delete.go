package branch

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Delete(repo *git.Repository, branch plumbing.ReferenceName, UserName, AccessToken string) error {
	// Delete remote branch
	fmt.Println("Start remote branch delete")
	pushOpts := &git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(":" + branch)},
		Auth: &http.BasicAuth{
			Username: UserName,
			Password: AccessToken,
		},
	}

	err := repo.Push(pushOpts)
	if err != nil {
		return err
	}
	fmt.Println("End remote branch delete")

	return nil
}
