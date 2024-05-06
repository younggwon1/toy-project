package sync

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/younggwon1/gitops-golang/external/argocd"
)

var (
	appName string
	dryRun  bool
	prune   bool
	force   bool
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
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("created argocd client with address: %s", argocdAddress)

		// init argocd app client
		appCli, err := cli.NewAppClient()
		if err != nil {
			return err
		}
		logger.Info().Msg("created argocd app client")

		// sync argocd app
		err = appCli.Sync(&argocd.AppSyncRequest{
			Name:   &appName,
			DryRun: &dryRun,
			Prune:  &prune,
			SyncStrategy: &argocd.AppSyncStrategyRequest{
				Force: force,
			},
		})
		if err != nil {
			return err
		}
		logger.Info().Msgf("synced argocd app: %s", appName)

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&appName, "appName", "a", "", "(required) argocd application name")
	Cmd.Flags().BoolVarP(&dryRun, "dryRun", "d", false, "(optional) argocd application dry run option, default: false")
	Cmd.Flags().BoolVarP(&prune, "prune", "p", false, "(optional) argocd application prune option, default: false")
	Cmd.Flags().BoolVarP(&force, "force", "f", false, "(optional) argocd application force option, default: false")
}
