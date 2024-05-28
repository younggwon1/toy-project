package sync

import (
	"fmt"
	"net/url"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/younggwon1/gitops-golang/external/argocd"
	"github.com/younggwon1/gitops-golang/external/jira"
	"github.com/younggwon1/gitops-golang/external/slack"
	"github.com/younggwon1/gitops-golang/util"
)

var (
	// argocd sync flags
	server string
	token  string
	name   string
	dryRun bool
	prune  bool
	force  bool
	// set git flags
	executor   string
	email      string
	repository string
	tag        string
	ticket     string
)

var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "run syncer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init logger
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

		// retrieve slackwebhook url from env vars
		slackWebHookUrl := os.Getenv("SLACK_WEBHOOK_URL")
		if slackWebHookUrl == "" {
			return fmt.Errorf("failed to retrieve `SLACK_WEBHOOK_URL` env var")
		}
		// retrieve jira url, username, token from env vars
		jiraUrl := os.Getenv("JIRA_URL")
		if jiraUrl == "" {
			return fmt.Errorf("failed to retrieve `JIRA_URL` env var")
		}
		jiraUsername := os.Getenv("JIRA_USERNAME")
		if jiraUsername == "" {
			return fmt.Errorf("failed to retrieve `JIRA_USERNAME` env var")
		}
		jiraToken := os.Getenv("JIRA_TOKEN")
		if jiraToken == "" {
			return fmt.Errorf("failed to retrieve `JIRA_TOKEN` env var")
		}

		// validate argocd server address and token
		if server == "" {
			return fmt.Errorf("failed because of `argocd address` was set to an empty value")
		}
		if token == "" {
			return fmt.Errorf("failed because of `argocd token` was set to an empty value")
		}
		// init argocd client
		cli, err := argocd.NewClient(&argocd.Connection{
			Address: server,
			Token:   token,
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("created argocd client with address: %s", server)

		// init argocd app client
		appCli, err := cli.NewAppClient()
		if err != nil {
			return err
		}
		logger.Info().Msg("succeed argocd app client")

		// validate jira ticket
		result := util.ValidateTicket(ticket)
		if !result {
			return fmt.Errorf("failed to validate the entered jira ticket: %s", ticket)
		}
		logger.Info().Msgf("succeed to validate the entered jira ticket: %s", ticket)

		// validate jira ticket status for deploying service
		err = jira.TicketStatusCheck(&jira.JiraConnection{
			JiraUrl:  jiraUrl,
			Username: jiraUsername,
			Token:    jiraToken,
		}, email, ticket)
		if err != nil {
			return err
		}
		logger.Info().Msg("succeed to check the entered jira ticket status for deploying")

		// sync argocd app
		argoCDAppUrl, err := url.JoinPath("https://", server, "applications", name)
		if err != nil {
			return err
		}
		audit, err := appCli.Sync(&argocd.AppSyncRequest{
			Name:   &name,
			DryRun: &dryRun,
			Prune:  &prune,
			SyncStrategy: &argocd.AppSyncStrategyRequest{
				Force: force,
			},
		}, executor, repository, tag, ticket, argoCDAppUrl)
		if err != nil {
			return err
		}
		logger.Info().Msgf("synced argocd app: %s", name)

		// send audit slack message
		tmpl := `version: {{ .Version.Version }}
metadata:
  name: {{ .Metadata.Name }}
  label:
    executor: {{ .Metadata.Label.Executor }}
  spec:
    src:
	  code:
		repo: {{ .Spec.Source.Code.Repo }}
		rev: {{ .Spec.Source.Code.Rev }}
	  helm:
		repo: {{ .Spec.Source.Helm.Repo }}
		chart: {{ .Spec.Source.Helm.Chart }}
		rev: {{ .Spec.Source.Helm.Rev }}
	  jira:
		ticket:
		  cr: {{ .Spec.Source.Jira.Ticket.CR }}
	dst:
	  argocd:
		url: {{ .Spec.Destination.ArgoCD.URL }}
		synced: {{ .Spec.Destination.ArgoCD.Synced }}
`
		err = slack.SendMessage(slackWebHookUrl, tmpl, audit)
		if err != nil {
			return err
		}
		logger.Info().Msg("succeed to send audit slack message")

		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&server, "server", "", "(required) argocd server address")
	Cmd.Flags().StringVar(&token, "token", "", "(required) argocd server token")
	Cmd.Flags().StringVar(&name, "name", "", "(required) argocd application name")
	Cmd.Flags().BoolVar(&dryRun, "dryRun", false, "(optional) argocd application dry run option, default: false")
	Cmd.Flags().BoolVar(&prune, "prune", false, "(optional) argocd application prune option, default: false")
	Cmd.Flags().BoolVar(&force, "force", false, "(optional) argocd application force option, default: false")
	Cmd.Flags().StringVar(&executor, "executor", "", "(required) executor name who deployed the service")
	Cmd.Flags().StringVar(&email, "email", "", "(required) executor email who deployed the service")
	Cmd.Flags().StringVar(&repository, "repository", "", "(required) git repository url")
	Cmd.Flags().StringVar(&tag, "tag", "", "(required) tag name to deploy")
	Cmd.Flags().StringVar(&ticket, "ticket", "", "(required) ticket name")
}
