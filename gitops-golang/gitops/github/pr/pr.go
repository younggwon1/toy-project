package pr

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v58/github"
)

func AutoCreateAndMerge(organisation string, helmRepo string, branch plumbing.ReferenceName, yamlFile string) error {
	repoOwner := organisation
	repoName := helmRepo
	baseBranch := "main"
	headBranch := branch.String()
	title := "Updated value in " + yamlFile
	body := "Updated image tag value in " + yamlFile

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(os.Getenv("AccessToken"))

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
