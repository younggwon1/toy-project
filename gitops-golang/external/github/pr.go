package github

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v58/github"
	"github.com/rs/zerolog"
)

func getTargetBranch(org, repoName, AccessToken string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(AccessToken)

	repo, _, err := client.Repositories.Get(ctx, org, repoName)
	if err != nil {
		return "", err
	}

	return *repo.DefaultBranch, nil
}

func AutoCreateAndMergePR(logger *zerolog.Logger, branch plumbing.ReferenceName, password, org, repo string) error {
	targetBranch, err := getTargetBranch(org, repo, password)
	if err != nil {
		return err
	}

	repoOwner := org
	repoName := repo
	baseBranch := targetBranch
	headBranch := branch.String()
	title := "Updated value in " + branch.String()
	body := "Updated image tag value in " + branch.String() + " branch"

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(password)

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

	logger.Info().Msg("created pull request: " + pullRequest.GetHTMLURL())
	prNumber := pullRequest.GetNumber()

	_, _, err = client.PullRequests.Merge(ctx, repoOwner, repoName, prNumber, "", &github.PullRequestOptions{})
	if err != nil {
		return err
	}

	logger.Info().Msg("merged pull request: " + pullRequest.GetHTMLURL())

	return nil
}
