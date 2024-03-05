package deploy

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	f "github.com/younggwon1/gitops-golang/file"
	"github.com/younggwon1/gitops-golang/gitops/git"
	"github.com/younggwon1/gitops-golang/gitops/github"
)

var (
	user   string
	email  string
	org    string
	repo   string
	file   string
	values string
)

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "run deployer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

		// init GitClient struct
		gitCli := git.NewGitClient()
		logger.Info().Msg("initialized Git Client")

		// clone github repository
		err := gitCli.Clone(
			org,
			repo,
		)
		if err != nil {
			return err
		}
		logger.Info().Msg("Cloned " + repo + " successfully")

		// create git branch
		b, err := gitCli.Create()
		if err != nil {
			return err
		}
		logger.Info().Msg("Created Branch " + b.String() + " successfully.\n")

		// checkout git branch
		err = gitCli.Checkout()
		if err != nil {
			return err
		}
		logger.Info().Msg("Checkout Branch " + b.String() + " successfully.\n")

		// modify values.yaml file
		err = f.ModifyFromYamlFile(repo, file, values)
		if err != nil {
			return err
		}
		logger.Info().Msg("Modified Values successfully.\n")

		// add modified files to git
		err = gitCli.Add(repo)
		if err != nil {
			return err
		}
		logger.Info().Msg("Changed Files added successfully to" + b.String() + " branch.\n")

		// commit added files to git
		err = gitCli.Commit(user, email, file)
		if err != nil {
			return err
		}
		logger.Info().Msg("Done a commit successfully.\n")

		// push committed files to git
		err = gitCli.Push()
		if err != nil {
			return err
		}
		logger.Info().Msg("Done a push successfully.\n")

		// auto create and merge pull request
		err = github.AutoCreateAndMerge(gitCli.Branch, org, repo, file)
		if err != nil {
			return err
		}
		logger.Info().Msg("Pull request created and merged successfully.\n")

		// delete git branch
		err = gitCli.Delete()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&user, "user", "u", "", "github user")
	Cmd.Flags().StringVarP(&email, "email", "e", "", "github user email")
	Cmd.Flags().StringVarP(&org, "organisation", "o", "", "github organisation")
	Cmd.Flags().StringVarP(&repo, "repository", "r", "", "git repository name")
	Cmd.Flags().StringVarP(&file, "file", "f", "", "values yaml file ")
	Cmd.Flags().StringVarP(&values, "values", "v", "", "image tag values")
}
