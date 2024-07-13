package config

type Deployer struct {
	Spec DeploySpec `yaml:"spec"`
}

type DeploySpec struct {
	Kubernetes *Kubernetes `yaml:"kubernetes"`
	Amplify    *Amplify    `yaml:"amplify"`
	CDN        *CDN        `yaml:"cdn"`
}

type Kubernetes struct {
	KubernetesDeploys []KubernetesDeploy `yaml:"deploys"`
}

type KubernetesDeploy struct {
	Helm   *HelmDeploy `yaml:"helm"`
	ArgoCD *ArgoCD     `yaml:"argocd"`
}

type HelmDeploy struct {
	Url          string      `yaml:"url"`
	Organization string      `yaml:"org"`
	Repository   string      `yaml:"repo"`
	Values       []HelmValue `yaml:"values"`
}

type HelmValue struct {
	File string `yaml:"file"`
}

type ArgoCD struct {
	Url  string      `yaml:"url"`
	Apps []ArgoCDApp `yaml:"apps"`
}

type ArgoCDApp struct {
	Name string `yaml:"name"`
}

type Amplify struct {
	AmplifyDeploys []AmplifyDeploy `yaml:"deploys"`
}

type AmplifyDeploy struct {
	AppId  string `yaml:"app-id"`
	Branch string `yaml:"branch"`
}

type CDN struct {
	CDNDeploys []CDNDeploy `yaml:"deploys"`
}

type CDNDeploy struct {
	S3Bucket                 string `yaml:"s3-bucket"`
	CloudFrontDistributionId string `yaml:"cloudfront-distribution-id"`
}

func (cfg *Deployer) ReadFromFile(filename string) error {
	return readFromFile(filename, cfg)
}
