package deploy

import (
	"github.com/younggwon1/gitops-golang/file"
	"github.com/younggwon1/gitops-golang/gitops/github"
	"github.com/younggwon1/gitops-golang/gitops/github/git"
	"github.com/younggwon1/gitops-golang/gitops/github/git/branch"
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
	Short: "run deployer server",
	RunE: func(cmd *cobra.Command, args []string) error {
		UserName := util.GetEnv("UserName", "")
		AccessToken := util.GetEnv("AccessToken", "")

		r, err := git.Clone(
			UserName,
			AccessToken,
			organisation,
			helmRepo,
		)
		if err != nil {
			return err
		}

		b, err := branch.Create(r)
		if err != nil {
			return err
		}

		err = branch.Checkout(r, b)
		if err != nil {
			return err
		}

		err = file.ModifyFromYamlFile(helmRepo, yamlFile, values)
		if err != nil {
			return err
		}

		err = git.Add(r, helmRepo)
		if err != nil {
			return err
		}

		err = git.Commit(r, user, userEmail, yamlFile)
		if err != nil {
			return err
		}

		err = git.Push(r, UserName, AccessToken)
		if err != nil {
			return err
		}

		err = github.AutoCreateAndMerge(b, organisation, helmRepo, yamlFile, AccessToken)
		if err != nil {
			return err
		}

		err = branch.Delete(r, b, UserName, AccessToken)
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
