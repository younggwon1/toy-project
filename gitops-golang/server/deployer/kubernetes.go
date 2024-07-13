package deployer

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	c "github.com/younggwon1/gitops-golang/config"
	"github.com/younggwon1/gitops-golang/external/argocd"
	"github.com/younggwon1/gitops-golang/external/github"
	f "github.com/younggwon1/gitops-golang/file"
	"github.com/younggwon1/gitops-golang/util"
)

const (
	DefaultGitURL string = "https://github.com"
)

func KubernetesProcess(logger *zerolog.Logger, user, email, values string, argoCDFlags *argocd.AppSyncRequest, cfg []c.KubernetesDeploy) error {
	// retrieve username, password from env vars
	username := util.GetEnv("GIT_USERNAME", "")
	if username == "" {
		return fmt.Errorf("failed to retrieve `GIT_USERNAME` env var")
	}
	password := util.GetEnv("GIT_PASSWORD", "")
	if password == "" {
		return fmt.Errorf("failed to retrieve `GIT_PASSWORD` env var")
	}
	// retrieve argocd server, token from env vars
	argoCDServer := os.Getenv("ARGOCD_SERVER")
	if argoCDServer == "" {
		return fmt.Errorf("failed to retrieve `ARGOCD_SERVER` env var")
	}
	argoCDToken := os.Getenv("ARGOCD_TOKEN")
	if argoCDToken == "" {
		return fmt.Errorf("failed to retrieve `ARGOCD_TOKEN` env var")
	}

	for _, spec := range cfg {
		// retrieve git url with org, repo names
		gitUrl := util.DefaultStr(spec.Helm.Url, DefaultGitURL)

		// init repoPath and targetPath
		repoPath := path.Join("/tmp", spec.Helm.Repository, uuid.New().String()[:8])

		// clone git repository
		gitCloneUrl, err := url.JoinPath(gitUrl, spec.Helm.Organization, spec.Helm.Repository)
		if err != nil {
			return err
		}
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
			Keep:   true, // Keep unstaged changes
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("checkout branch %s successfully", branchRefName.Short())

		// init paths
		var paths []string
		for _, value := range spec.Helm.Values {
			path := path.Join(repoPath, value.File)
			paths = append(paths, path)
		}

		// init wg
		wg := sync.WaitGroup{}

		// init goroutines for modifying values
		wg.Add(len(paths))

		// init goroutines
		for _, path := range paths {
			go func(path string) {
				defer wg.Done()
				err := f.ModifyYamlFile(path, values)
				if err != nil {
					logger.Err(err).Msg("")
				}
				logger.Info().Msgf("%s, modified values successfully", path)
			}(path)
		}

		// wait until all goroutines are done
		wg.Wait()

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
			RefSpecs:   []config.RefSpec{config.RefSpec("+" + branchRefName + ":" + branchRefName)},
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
		err = github.AutoCreateAndMergePR(logger, branchRefName, password, spec.Helm.Organization, spec.Helm.Repository)
		if err != nil {
			return err
		}
		logger.Info().Msg("created and merged pull request successfully")

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

		// init argocd client
		cli, err := argocd.NewClient(&argocd.Connection{
			Address: argoCDServer,
			Token:   argoCDToken,
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("created argocd client with address: %s", argoCDServer)

		// init argocd app client
		appCli, err := cli.NewAppClient()
		if err != nil {
			return err
		}
		logger.Info().Msg("succeed argocd app client")

		// init goroutines for argocd apps
		wg.Add(len(spec.ArgoCD.Apps))

		for _, app := range spec.ArgoCD.Apps {
			go func(app c.ArgoCDApp) {
				defer wg.Done()
				// sync argocd app
				err := appCli.Sync(&argocd.AppSyncRequest{
					Name:   &app.Name,
					DryRun: argoCDFlags.DryRun,
					Prune:  argoCDFlags.Prune,
					SyncStrategy: &argocd.AppSyncStrategyRequest{
						Force: argoCDFlags.SyncStrategy.Force,
					},
				})
				if err != nil {
					logger.Err(err).Msg("")
				}
				logger.Info().Msgf("synced argocd app: %s", app.Name)
			}(app)
		}

		// wait until all goroutines are done
		wg.Wait()
	}

	return nil
}
