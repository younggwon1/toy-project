package branchcli

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Create(repo *git.Repository) (plumbing.ReferenceName, error) {
	fmt.Printf("Start creating branch.\n")

	headRef, err := repo.Head()
	fmt.Println(headRef.Hash())
	if err != nil {
		fmt.Println(err)
	}

	// newBranchRefName := plumbing.NewBranchReferenceName("gogit-deploy-" + headRef.Hash().String()[:7])
	newBranchRefName := plumbing.NewBranchReferenceName("gogit-deploy-" + uuid.New().String()[:8])
	newBranchRef := plumbing.NewHashReference(newBranchRefName, headRef.Hash())

	if err := repo.Storer.SetReference(newBranchRef); err != nil {
		return "", err
	}

	fmt.Printf("Branch %s created successfully.\n", newBranchRefName)
	return newBranchRefName, nil
}
