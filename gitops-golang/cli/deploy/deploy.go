package deploy

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/younggwon1/gitops-golang/external/github"
	f "github.com/younggwon1/gitops-golang/file"
	"github.com/younggwon1/gitops-golang/util"
)

var (
	user         string
	email        string
	gitUrl       string
	organisation string
	repository   string
	file         string
	values       string
)

const (
	DefaultGitURL = "https://github.com"
)

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "run deployer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init repoPath and targetPath
		repoPath := path.Join("/tmp", repository)
		targetPath := path.Join(repoPath, file)

		// init logger
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

		// retrieve username, password from env vars
		username := util.GetEnv("GIT_USERNAME", "")
		if username == "" {
			return fmt.Errorf("failed to retrieve `GIT_USERNAME` env var")
		}
		password := util.GetEnv("GIT_PASSWORD", "")
		if password == "" {
			return fmt.Errorf("failed to retrieve `GIT_PASSWORD` env var")
		}

		// retrieve git url with org, repo names
		gitUrl := util.DefaultStr(gitUrl, DefaultGitURL)
		gitCloneUrl, err := url.JoinPath(gitUrl, organisation, repository)
		if err != nil {
			return err
		}

		// clone git repository
		repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: username,
				Password: password,
			},
			URL:      gitCloneUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("cloned %s successfully", gitCloneUrl)

		// create git branch
		headRef, err := repo.Head()
		if err != nil {
			return err
		}
		branchRefName := plumbing.NewBranchReferenceName("gitops-deploy-" + uuid.New().String()[:8])
		branchRef := plumbing.NewHashReference(branchRefName, headRef.Hash())
		err = repo.Storer.SetReference(branchRef)
		if err != nil {
			return err
		}
		logger.Info().Msgf("created branch %s successfully", branchRefName.Short())

		// checkout git branch
		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: branchRefName,
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("checkout branch %s successfully", branchRefName.Short())

		// modify values.yaml file
		err = f.ModifyYamlFile(targetPath, values)
		if err != nil {
			return err
		}
		logger.Info().Msg("modified values successfully")

		// check if anything has changed
		status, err := worktree.Status()
		if err != nil {
			return err
		}
		if status.IsClean() {
			logger.Info().Msg("found nothing to commit, working tree clean")
			return nil
		}

		// add changed files to git
		err = worktree.AddWithOptions(&git.AddOptions{
			All:  true,
			Path: repoPath,
		})
		if err != nil {
			return err
		}
		logger.Info().Msg("added changed files to git successfully")

		// commit added files to git
		message := "changed image tag in " + branchRef.Hash().String() + " " + branchRefName.Short()
		commit, err := worktree.Commit(message, &git.CommitOptions{
			Author: &object.Signature{
				Name:  user,
				Email: email,
				When:  time.Now(),
			},
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("committed %s successfully", commit.String())

		// push committed files to git
		err = repo.Push(&git.PushOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username: username,
				Password: password,
			},
		})
		if err != nil {
			return err
		}
		logger.Info().Msg("pushed successfully")

		// auto create and merge pull request
		err = github.AutoCreateAndMergePR(&logger, branchRefName, password, organisation, repository)
		if err != nil {
			return err
		}

		// delete git branch
		err = repo.Push(&git.PushOptions{
			RefSpecs: []config.RefSpec{config.RefSpec(":" + branchRefName)},
			Auth: &http.BasicAuth{
				Username: username,
				Password: password,
			},
		})
		if err != nil {
			logger.Debug().Msg("failed to delete branch but not an error")
		}
		logger.Info().Msgf("deleted branch %s successfully", branchRefName.Short())

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&user, "user", "u", "", "git user")
	Cmd.Flags().StringVarP(&email, "email", "e", "", "git user email")
	Cmd.Flags().StringVarP(&gitUrl, "gitUrl", "g", "", "git url")
	Cmd.Flags().StringVarP(&organisation, "organisation", "o", "", "git organisation")
	Cmd.Flags().StringVarP(&repository, "repository", "r", "", "git repository name")
	Cmd.Flags().StringVarP(&file, "file", "f", "", "values yaml file ")
	Cmd.Flags().StringVarP(&values, "values", "v", "", "image tag values")
}
