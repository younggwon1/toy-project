package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDeployerConfig(t *testing.T) {
	input := `
spec:
  kubernetes:
    deploys:
    - helm:
        url: "helm-url-1"
        org: "helm-org-1"
        repo: "helm-repo-1"
        values:
        - file: "value-file-1"
        - file: "value-file-2"
      argocd:
        url: "argocd-url"
        apps:
        - name: "app-name-1"
        - name: "app-name-2"
  amplify:
    deploys:
    - app-id: app-id-1
      branch: branch-name-1
    - app-id: app-id-2
      branch: branch-name-2
  cdn:
    deploys:
    - s3-bucket: "bucket-name-1"
      cloudfront-distribution-id: "distribution-id-1"
    - s3-bucket: "bucket-name-2"
      cloudfront-distribution-id: "distribution-id-2"
`

	output := Deployer{}
	err := yaml.Unmarshal([]byte(input), &output)
	assert.NoError(t, err)
	assert.Equal(t, Deployer{
		Spec: DeploySpec{
			Kubernetes: &Kubernetes{
				KubernetesDeploys: []KubernetesDeploy{
					{
						Helm: &HelmDeploy{
							Url:          "helm-url-1",
							Organization: "helm-org-1",
							Repository:   "helm-repo-1",
							Values: []HelmValue{
								{
									File: "value-file-1",
								},
								{
									File: "value-file-2",
								},
							},
						},
						ArgoCD: &ArgoCD{
							Url: "argocd-url",
							Apps: []ArgoCDApp{
								{
									Name: "app-name-1",
								},
								{
									Name: "app-name-2",
								},
							},
						},
					},
				},
			},
			Amplify: &Amplify{
				AmplifyDeploys: []AmplifyDeploy{
					{
						AppId:  "app-id-1",
						Branch: "branch-name-1",
					},
					{
						AppId:  "app-id-2",
						Branch: "branch-name-2",
					},
				},
			},
			CDN: &CDN{
				CDNDeploys: []CDNDeploy{
					{
						S3Bucket:                 "bucket-name-1",
						CloudFrontDistributionId: "distribution-id-1",
					},
					{
						S3Bucket:                 "bucket-name-2",
						CloudFrontDistributionId: "distribution-id-2",
					},
				},
			},
		},
	}, output)
}

func TestDeployerKubernetesConfig(t *testing.T) {
	input := `
spec:
  kubernetes:
    deploys:
    - helm:
        url: "helm-url-1"
        org: "helm-org-1"
        repo: "helm-repo-1"
        values:
        - file: "value-file-1"
        - file: "value-file-2"
      argocd:
        url: "argocd-url"
        apps:
        - name: "app-name-1"
        - name: "app-name-2"
`
	output := Deployer{}
	err := yaml.Unmarshal([]byte(input), &output)
	assert.NoError(t, err)
	assert.Equal(t, Deployer{
		Spec: DeploySpec{
			Kubernetes: &Kubernetes{
				KubernetesDeploys: []KubernetesDeploy{
					{
						Helm: &HelmDeploy{
							Url:          "helm-url-1",
							Organization: "helm-org-1",
							Repository:   "helm-repo-1",
							Values: []HelmValue{
								{
									File: "value-file-1",
								},
								{
									File: "value-file-2",
								},
							},
						},
						ArgoCD: &ArgoCD{
							Url: "argocd-url",
							Apps: []ArgoCDApp{
								{
									Name: "app-name-1",
								},
								{
									Name: "app-name-2",
								},
							},
						},
					},
				},
			},
			Amplify: nil,
			CDN:     nil,
		},
	}, output)
}

func TestDeployerAmplifyConfig(t *testing.T) {
	input := `
spec:
  amplify:
    deploys:
    - app-id: app-id-1
      branch: branch-name-1
    - app-id: app-id-2
      branch: branch-name-2
`
	output := Deployer{}
	err := yaml.Unmarshal([]byte(input), &output)
	assert.NoError(t, err)
	assert.Equal(t, Deployer{
		Spec: DeploySpec{
			Kubernetes: nil,
			Amplify: &Amplify{
				AmplifyDeploys: []AmplifyDeploy{
					{
						AppId:  "app-id-1",
						Branch: "branch-name-1",
					},
					{
						AppId:  "app-id-2",
						Branch: "branch-name-2",
					},
				},
			},
			CDN: nil,
		},
	}, output)
}

func TestDeployerCDNConfig(t *testing.T) {
	input := `
spec:
  cdn:
    deploys:
    - s3-bucket: "bucket-name-1"
      cloudfront-distribution-id: "distribution-id-1"
    - s3-bucket: "bucket-name-2"
      cloudfront-distribution-id: "distribution-id-2"
`
	output := Deployer{}
	err := yaml.Unmarshal([]byte(input), &output)
	assert.NoError(t, err)
	assert.Equal(t, Deployer{
		Spec: DeploySpec{
			Kubernetes: nil,
			Amplify:    nil,
			CDN: &CDN{
				CDNDeploys: []CDNDeploy{
					{
						S3Bucket:                 "bucket-name-1",
						CloudFrontDistributionId: "distribution-id-1",
					},
					{
						S3Bucket:                 "bucket-name-2",
						CloudFrontDistributionId: "distribution-id-2",
					},
				},
			},
		},
	}, output)
}
