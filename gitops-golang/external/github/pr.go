package github

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v58/github"
	"github.com/rs/zerolog"
)

func getTargetBranch(ctx context.Context, client *github.Client, org, repoName string) (*github.Branch, error) {
	repo, _, err := client.Repositories.Get(ctx, org, repoName)
	if err != nil {
		return nil, err
	}
	br, _, err := client.Repositories.GetBranch(ctx, org, repoName, *repo.DefaultBranch, 1)
	if err != nil {
		return nil, err
	}

	return br, nil
}

func AutoCreateAndMergePR(logger *zerolog.Logger, branch plumbing.ReferenceName, authtoken, org, repo string) error {
	// init github client
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(authtoken)

	// get current target branch
	br, err := getTargetBranch(ctx, client, org, repo)
	if err != nil {
		return err
	}

	repoOwner := org
	repoName := repo
	baseBranch := *br.Name
	headBranch := branch.String()
	title := "Updated value in " + branch.String()
	body := "Updated image tag value in " + branch.String() + " branch"
	pr := &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Head:  &headBranch,
		Base:  &baseBranch,
	}

	// create pull request
	pullRequest, _, err := client.PullRequests.Create(ctx, repoOwner, repoName, pr)
	if err != nil {
		return err
	}
	logger.Info().Msg("created pull request: " + pullRequest.GetHTMLURL())

	// get pull request number
	prNumber := pullRequest.GetNumber()

	// init ticker
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// check pull request status per each tick
	loop := true
	retryCount := 0
	for loop {
		select {
		case <-ctx.Done():
			loop = false
		case <-ticker.C:
			// merged pull request
			_, response, err := client.PullRequests.Merge(ctx, repoOwner, repoName, prNumber, "", &github.PullRequestOptions{})
			// issue : 405 Base branch was modified. Review and try the merge again.
			// A wait of more than 1 second is required between GitHub's API requests.
			// Therefore, when 405 occurs, implement logic to wait for more than 1s and then try again to ensure success.
			// https://docs.github.com/ko/rest/using-the-rest-api/best-practices-for-using-the-rest-api?apiVersion=2022-11-28#pause-between-mutative-requests
			if response.StatusCode == 405 {
				logger.Info().Msgf("occured 405 status code when merging a %s", pullRequest.GetHTMLURL())
				retryCount++
				if retryCount > 3 {
					loop = false
				} else {
					time.Sleep(2 * time.Second)
				}
			} else if err != nil {
				return err
			} else {
				loop = false
			}
		}
	}
	logger.Info().Msg("merged pull request: " + pullRequest.GetHTMLURL())

	return nil
}
