package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/aws/aws-sdk-go-v2/service/amplify/types"
)

type StartAmplifyJobInput struct {
	JobType       types.JobType
	CommitId      *string
	CommitMessage *string
}

func (cli *AWSClient) GetAmplifyApp(appId string) (*amplify.GetAppOutput, error) {
	output, err := cli.amplifyClient.GetApp(cli.ctx, &amplify.GetAppInput{
		AppId: &appId,
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (cli *AWSClient) GetAmplifyBranch(appId, branchName string) (*amplify.GetBranchOutput, error) {
	output, err := cli.amplifyClient.GetBranch(cli.ctx, &amplify.GetBranchInput{
		AppId:      &appId,
		BranchName: &branchName,
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (cli *AWSClient) CreateAmplifyBranch(appId, branchName string) (*amplify.CreateBranchOutput, error) {
	output, err := cli.amplifyClient.CreateBranch(cli.ctx, &amplify.CreateBranchInput{
		AppId:      &appId,
		BranchName: &branchName,
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (cli *AWSClient) StartAmplifyJob(appId, branchName string, s *StartAmplifyJobInput) (*amplify.StartJobOutput, error) {
	output, err := cli.amplifyClient.StartJob(cli.ctx, &amplify.StartJobInput{
		AppId:         &appId,
		BranchName:    &branchName,
		JobType:       s.JobType,
		CommitId:      s.CommitId,
		CommitMessage: s.CommitMessage,
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}
