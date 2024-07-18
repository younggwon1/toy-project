package deploy

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/amplify/types"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	c "github.com/younggwon1/gitops-golang/config"
	"github.com/younggwon1/gitops-golang/external/argocd"
	"github.com/younggwon1/gitops-golang/external/aws"
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

	// amplify flags
	jobType       string
	commitId      *string
	commitMessage *string
)

type Kubernetes interface {
	DeployKubernetes(logger *zerolog.Logger) error
}

type Amplify interface {
	DeployAmplify(logger *zerolog.Logger) error
}

type DeployKubernetesFlags struct {
	user                    string
	email                   string
	values                  string
	argoCDFlags             *argocd.AppSyncRequest
	kubernetesDeploysConfig interface{}
}

type DeployAmplifyFlags struct {
	amplifyFlags         *aws.StartAmplifyJobInput
	amplifyDeploysConfig interface{}
}

func (dkf *DeployKubernetesFlags) DeployKubernetes(logger *zerolog.Logger) error {
	err := deployer.KubernetesProcess(logger, dkf.user, dkf.email, dkf.values, dkf.argoCDFlags, dkf.kubernetesDeploysConfig.([]c.KubernetesDeploy))
	if err != nil {
		return err
	}

	return nil
}

func (daf *DeployAmplifyFlags) DeployAmplify(logger *zerolog.Logger) error {
	err := deployer.AmplifyProcess(logger, daf.amplifyFlags, daf.amplifyDeploysConfig.([]c.AmplifyDeploy))
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
			var k Kubernetes
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
			k = &k8sDeployFlags
			err := k.DeployKubernetes(&logger)
			if err != nil {
				return err
			}
		}

		if cfg.Spec.Amplify != nil {
			var a Amplify
			// set amplify flags
			amplifyDeployFlags := DeployAmplifyFlags{
				amplifyDeploysConfig: cfg.Spec.Amplify.AmplifyDeploys,
				amplifyFlags: &aws.StartAmplifyJobInput{
					JobType:       types.JobType(jobType),
					CommitId:      commitId,
					CommitMessage: commitMessage,
				},
			}

			//deploy amplify
			a = &amplifyDeployFlags
			err := a.DeployAmplify(&logger)
			if err != nil {
				return err
			}
		}

		if cfg.Spec.CDN != nil {
			logger.Info().Msg("deploy cdn")
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
	Cmd.Flags().StringVar(&jobType, "jobType", "RELEASE", "amplify job type, default: RELEASE")
	Cmd.Flags().StringVar(commitId, "commitId", "", "amplify commit id")
	Cmd.Flags().StringVar(commitMessage, "commitMessage", "", "amplify commit message")
}
