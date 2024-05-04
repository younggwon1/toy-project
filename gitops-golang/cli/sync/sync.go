package sync

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/younggwon1/gitops-golang/external/argocd"
)

var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "run syncer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init logger
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

		// retrieve argocd address, token from env vars
		argocdAddress := os.Getenv("ARGOCD_ADDRESS")
		if argocdAddress == "" {
			return fmt.Errorf("failed to retrieve `ARGOCD_ADDRESS` env var")
		}
		argocdToken := os.Getenv("ARGOCD_TOKEN")
		if argocdToken == "" {
			return fmt.Errorf("failed to retrieve `ARGOCD_TOKEN` env var")
		}

		// init argocd client
		cli, err := argocd.NewClient(&argocd.Connection{
			Address: argocdAddress,
			Token:   argocdToken,
		}, logger)
		if err != nil {
			return err
		}

		// init argocd app client
		appCli, err := cli.NewAppClient()
		if err != nil {
			return err
		}

		// sync argocd app
		err = appCli.Sync()
		if err != nil {
			return err
		}

		return nil
	},
}
