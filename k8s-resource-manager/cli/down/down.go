package down

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/younggwon1/k8s-resource-manager/external/kubernetes"
)

var (
	namespace string
	name      string
)

var Cmd = &cobra.Command{
	Use:   "down",
	Short: "kuberenetes resource down operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init logger
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Info().Msg("init logger")

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

		// scale down error deployment
		if name != "" {
			err = cli.ScaleDown(name, namespace)
			if err != nil {
				return err
			}
			logger.Info().Msgf("succeed to scale down %s deployment in %s namespace", name, namespace)
		} else {
			err = cli.ScaleDownErrorDeployment(namespace)
			if err != nil {
				return err
			}
			logger.Info().Msgf("succeed to scale down deployments in %s namespace", namespace)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&namespace, "namespace", "", "(required) kubernetes namespace")
	Cmd.Flags().StringVar(&name, "name", "", "(optional) kubernetes deployment name")
}
