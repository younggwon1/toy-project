package deployer

import (
	"github.com/rs/zerolog"

	c "github.com/younggwon1/gitops-golang/config"
)

func AmplifyProcess(logger *zerolog.Logger, cfg []c.AmplifyDeploy) error {
	// 1. Check if Amplify APP Exists and Amplify APP Name to ENV
	// 2. Deploy to Amplify with the AWS CLI

	for _, spec := range cfg {
		logger.Info().Msgf("Deploying Amplify app %s on branch %s", spec.AppId, spec.Branch)

	}

	return nil
}
