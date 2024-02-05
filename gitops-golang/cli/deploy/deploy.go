package deploy

import (
	"fmt"
	"log"
	"os"

	"github/younggwon1/gitops-golang/file"
	"github/younggwon1/gitops-golang/gitops/github/branchcli"
	"github/younggwon1/gitops-golang/gitops/github/gitcli"
	"github/younggwon1/gitops-golang/gitops/github/pr"
	"github/younggwon1/gitops-golang/gitops/github/repository"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	gitUser      string
	gitUserEmail string
	organisation string
	helmRepo     string
	yamlFile     string
	values       string
)

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "run deployer server",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		repo, err := repository.Clone(
			os.Getenv("Username"),
			os.Getenv("AccessToken"),
			organisation,
			helmRepo,
		)
		if err != nil {
			fmt.Println(err)
		}

		branch, err := branchcli.Create(repo)
		if err != nil {
			fmt.Println(err)
		}

		err = branchcli.Checkout(repo, branch)
		if err != nil {
			fmt.Println(err)
		}

		err = file.ModifyImageTagInYAMLFile(helmRepo, yamlFile, values)
		if err != nil {
			fmt.Println(err)
		}

		err = gitcli.Add(repo, helmRepo)
		if err != nil {
			fmt.Println(err)
		}

		err = gitcli.Commit(repo, gitUser, gitUserEmail, yamlFile)
		if err != nil {
			fmt.Println(err)
		}

		err = gitcli.Push(repo)
		if err != nil {
			fmt.Println(err)
		}

		err = pr.AutoCreateAndMerge(organisation, helmRepo, branch, yamlFile)
		if err != nil {
			fmt.Println(err)
		}

		err = branchcli.Delete(repo, branch)
		if err != nil {
			fmt.Println(err)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&gitUser, "git-user", "u", "", "github user")
	Cmd.Flags().StringVarP(&gitUserEmail, "git-email", "e", "", "github user email")
	Cmd.Flags().StringVarP(&organisation, "organisation", "o", "", "github organisation")
	Cmd.Flags().StringVarP(&helmRepo, "repository-name", "r", "", "helm repository")
	Cmd.Flags().StringVarP(&yamlFile, "file", "f", "", "values yaml file")
	Cmd.Flags().StringVarP(&values, "values", "v", "", "image tag values")
}
