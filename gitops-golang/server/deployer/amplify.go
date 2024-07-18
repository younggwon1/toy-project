package deployer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	c "github.com/younggwon1/gitops-golang/config"
	"github.com/younggwon1/gitops-golang/external/aws"
)

func AmplifyProcess(logger *zerolog.Logger, amplifyFlags *aws.StartAmplifyJobInput, cfg []c.AmplifyDeploy) error {
	// init ctx
	ctx := context.Background()

	// init aws client
	awsCli, err := aws.NewAWSClient(ctx)
	if err != nil {
		return err
	}
	logger.Info().Msg("initialized aws client")

	for _, spec := range cfg {
		// check required fields
		if spec.AppId == "" {
			return fmt.Errorf("required setting AppId")
		}
		if spec.Branch == "" {
			return fmt.Errorf("required setting Branch")
		}

		// check if Amplify APP exists
		appOutput, err := awsCli.GetAmplifyApp(spec.AppId)
		if err != nil {
			return err
		}
		logger.Info().Msgf("exists amplify app %s", *appOutput.App.Name)

		// check if Amplify Branch exists
		branchOutput, err := awsCli.GetAmplifyBranch(spec.AppId, spec.Branch)
		if branchOutput == nil {
			// create Amplify Branch
			_, err = awsCli.CreateAmplifyBranch(spec.AppId, spec.Branch)
			if err != nil {
				return err
			}
			logger.Info().Msgf("created amplify branch %s", spec.Branch)
		}
		if err != nil {
			return err
		}

		// started Amplify Job
		_, err = awsCli.StartAmplifyJob(spec.AppId, spec.Branch, amplifyFlags)
		if err != nil {
			return err
		}
	}

	return nil
}
