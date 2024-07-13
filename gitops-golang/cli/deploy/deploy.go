package deploy

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	c "github.com/younggwon1/gitops-golang/config"
	"github.com/younggwon1/gitops-golang/external/argocd"
	"github.com/younggwon1/gitops-golang/server/deployer"
)

var (
	// git config
	user  string
	email string

	// helm charts
	values string

	// spec for deploy
	spec string

	// argocd sync flags
	dryRun bool
	prune  bool
	force  bool
)

type DeployProcess interface {
	kubernetes(logger *zerolog.Logger) error
}

type DeployKubernetesFlags struct {
	user                    string
	email                   string
	values                  string
	argoCDFlags             *argocd.AppSyncRequest
	kubernetesDeploysConfig interface{}
}

func (dkf *DeployKubernetesFlags) kubernetes(logger *zerolog.Logger) error {
	err := deployer.KubernetesProcess(logger, dkf.user, dkf.email, dkf.values, dkf.argoCDFlags, dkf.kubernetesDeploysConfig.([]c.KubernetesDeploy))
	if err != nil {
		return err
	}

	return nil
}

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "run deployer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init logger
		logger := zerolog.New(os.Stdout).With().
			Timestamp().
			Logger()

		// read config from file
		cfg := c.Deployer{}
		err := cfg.ReadFromFile(spec)
		if err != nil {
			return err
		}
		logger.Info().Msgf("read config from %s", spec)

		if cfg.Spec.Kubernetes != nil {
			var dp DeployProcess
			// set kubernetes flags
			k8sDeployFlags := DeployKubernetesFlags{
				user:   user,
				email:  email,
				values: values,
				argoCDFlags: &argocd.AppSyncRequest{
					DryRun: &dryRun,
					Prune:  &prune,
					SyncStrategy: &argocd.AppSyncStrategyRequest{
						Force: force,
					},
				},
				kubernetesDeploysConfig: cfg.Spec.Kubernetes.KubernetesDeploys,
			}

			// deploy kubernetes
			dp = &k8sDeployFlags
			err := dp.kubernetes(&logger)
			if err != nil {
				return err
			}
		}

		if cfg.Spec.CDN != nil {
			logger.Info().Msg("deploy cdn")
		}

		if cfg.Spec.Amplify != nil {
			logger.Info().Msg("deploy amplify")
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&user, "user", "", "git user")
	Cmd.Flags().StringVar(&email, "email", "", "git user email")
	Cmd.Flags().StringVar(&values, "values", "", "image tag values for modify")
	Cmd.Flags().StringVar(&spec, "spec", "", "spec template for deploy")
	Cmd.Flags().BoolVar(&dryRun, "dryRun", false, "(optional) argocd application dry run option, default: false")
	Cmd.Flags().BoolVar(&prune, "prune", false, "(optional) argocd application prune option, default: false")
	Cmd.Flags().BoolVar(&force, "force", false, "(optional) argocd application force option, default: false")
}
