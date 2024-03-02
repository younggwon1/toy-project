package deploy

import (
	"github.com/younggwon1/gitops-golang/file"
	"github.com/younggwon1/gitops-golang/gitops/git"
	"github.com/younggwon1/gitops-golang/gitops/github"
	"github.com/younggwon1/gitops-golang/util"

	"github.com/spf13/cobra"
)

var (
	user         string
	userEmail    string
	organisation string
	helmRepo     string
	yamlFile     string
	values       string
)

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "run deployer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		userName := util.GetEnv("UserName", "")
		accessToken := util.GetEnv("AccessToken", "")

		// github repository clone
		r, err := git.Clone(
			userName,
			accessToken,
			organisation,
			helmRepo,
		)
		if err != nil {
			return err
		}

		// init GitClient struct
		gitCli := git.NewGitClient(r)

		// create git branch
		b, err := gitCli.Create()
		if err != nil {
			return err
		}

		// checkout git branch
		err = gitCli.Checkout()
		if err != nil {
			return err
		}

		// modify values.yaml file
		err = file.ModifyFromYamlFile(helmRepo, yamlFile, values)
		if err != nil {
			return err
		}

		// add modified files to git
		err = gitCli.Add(helmRepo)
		if err != nil {
			return err
		}

		// commit added files to git
		err = gitCli.Commit(user, userEmail, yamlFile)
		if err != nil {
			return err
		}

		// push committed files to git
		err = gitCli.Push(userName, accessToken)
		if err != nil {
			return err
		}

		// auto create and merge pull request
		err = github.AutoCreateAndMerge(b, organisation, helmRepo, yamlFile, accessToken)
		if err != nil {
			return err
		}

		// delete git branch
		err = gitCli.Delete(userName, accessToken)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&user, "user", "u", "", "github user")
	Cmd.Flags().StringVarP(&userEmail, "email", "e", "", "github user email")
	Cmd.Flags().StringVarP(&organisation, "organisation", "o", "", "github organisation")
	Cmd.Flags().StringVarP(&helmRepo, "repository", "r", "", "helm repository name")
	Cmd.Flags().StringVarP(&yamlFile, "file", "f", "", "values yaml file ")
	Cmd.Flags().StringVarP(&values, "values", "v", "", "image tag values")
}
