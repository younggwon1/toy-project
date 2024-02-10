package github

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v58/github"
)

func getDefaultBranch(organisation, repoName, AccessToken string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(AccessToken)

	repo, _, err := client.Repositories.Get(ctx, organisation, repoName)
	if err != nil {
		return "", err
	}

	return *repo.DefaultBranch, nil
}

func AutoCreateAndMerge(branch plumbing.ReferenceName, organisation, helmRepo, yamlFile, AccessToken string) error {
	defaultBranch, err := getDefaultBranch(organisation, helmRepo, AccessToken)
	if err != nil {
		return err
	}

	repoOwner := organisation
	repoName := helmRepo
	baseBranch := defaultBranch
	headBranch := branch.String()
	title := "Updated value in " + yamlFile
	body := "Updated image tag value in " + yamlFile

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(AccessToken)

	pr := &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Head:  &headBranch,
		Base:  &baseBranch,
	}

	pullRequest, _, err := client.PullRequests.Create(ctx, repoOwner, repoName, pr)
	if err != nil {
		return err
	}

	fmt.Printf("Pull request created: %s\n", pullRequest.GetHTMLURL())
	prNumber := pullRequest.GetNumber()

	_, _, err = client.PullRequests.Merge(ctx, repoOwner, repoName, prNumber, "", &github.PullRequestOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Pull request merged: %s\n", pullRequest.GetHTMLURL())

	return nil
}
