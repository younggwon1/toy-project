package github

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v58/github"
)

func getDefaultBranch(org, repoName, AccessToken string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(AccessToken)

	repo, _, err := client.Repositories.Get(ctx, org, repoName)
	if err != nil {
		return "", err
	}

	return *repo.DefaultBranch, nil
}

func AutoCreateAndMerge(branch plumbing.ReferenceName, org, repo, file string) error {
	githubCli := NewGithubClient()

	defaultBranch, err := getDefaultBranch(org, repo, githubCli.authToken)
	if err != nil {
		return err
	}

	repoOwner := org
	repoName := repo
	baseBranch := defaultBranch
	headBranch := branch.String()
	title := "Updated value in " + file
	body := "Updated image tag value in " + file

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(githubCli.authToken)

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

	githubCli.logger.Info().Msg("Pull request created: " + pullRequest.GetHTMLURL())
	prNumber := pullRequest.GetNumber()

	_, _, err = client.PullRequests.Merge(ctx, repoOwner, repoName, prNumber, "", &github.PullRequestOptions{})
	if err != nil {
		return err
	}

	githubCli.logger.Info().Msg("Pull request merged: " + pullRequest.GetHTMLURL())

	return nil
}
