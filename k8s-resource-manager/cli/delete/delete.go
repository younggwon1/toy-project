package delete

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/younggwon1/k8s-resource-manager/external/kubernetes"
	"github.com/younggwon1/k8s-resource-manager/external/slack"
)

var (
	namespace string
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "kuberenetes resource delete operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init logger
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Info().Msg("init logger")

		// retrieve slackwebhook url from env vars
		slackWebHookUrl := os.Getenv("SLACK_WEBHOOK_URL")
		if slackWebHookUrl == "" {
			return fmt.Errorf("failed to retrieve `SLACK_WEBHOOK_URL` env var")
		}
		logger.Info().Msg("retrieve slack webhook url")

		// init kubernetes config
		cli, err := kubernetes.NewClient(
			&logger,
		)
		if err != nil {
			return err
		}
		logger.Info().Msg("init kubernetes credentials")

		// validate kubernetes namespace
		if namespace == "" {
			return fmt.Errorf("failed because of `kubernetes namespace` was set to an empty value")
		}
		logger.Info().Msg("validate kubernetes namespace")

		// delete replica 0 deployment
		deleteNames, err := cli.AllDelete(namespace)
		if err != nil {
			return err
		}
		logger.Info().Msgf("succeed to delete replica 0 deployments in %s namespace", namespace)

		// send slack message
		tmpl := `{"status": "Delete", "name": "{{.Name}}", "namespace": "{{.Namespace}}","time": "{{.Time}}"}`
		message := slack.Template{
			Name:      deleteNames,
			Namespace: namespace,
			Time:      time.Now().Format("2006-01-02 15:04:05"),
		}
		err = slack.SendMessage(slackWebHookUrl, tmpl, message)
		if err != nil {
			return err
		}
		logger.Info().Msg("succeed to send slack message")

		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&namespace, "namespace", "", "(required) kubernetes namespace")
}
