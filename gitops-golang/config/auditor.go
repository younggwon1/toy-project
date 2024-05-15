package config

// version: v0.0.1
// metadata:
//   name: gitops-golang
//   label:
//     executor: younggwon
// spec:
//   src:
//     code:
//       repo: https://github.com/younggwon1/gitops-golang
//       rev: develop
//     helm:
//       repo: git@github.com:younggwon1/gitops-golang-helm-chart.git
//       chart: charts/gitops-golang/values-dev.yaml
//       rev: ghhfv24hv2hv8389273423jhfkjdshf893298
//     jira:
//       ticket:
//         cr: CR-XX
//   dst:
//     argocd:
//       url: https://<argocd-url>/applications/gitops-golang
//       synced: 2024-05-06 21:18:10

type Auditor struct {
	Version  string `yaml:"version"`
	Metadata `yaml:"metadata"`
	Spec     `yaml:"spec"`
}

type Metadata struct {
	Name  string `yaml:"name"`
	Label Label  `yaml:"label"`
}

type Label struct {
	Executor string `yaml:"executor"`
}

type Spec struct {
	Source      Source      `yaml:"src"`
	Destination Destination `yaml:"dst"`
}

type Source struct {
	Code Code `yaml:"code"`
	Helm Helm `yaml:"helm"`
	Jira Jira `yaml:"jira"`
}

type Destination struct {
	ArgoCD ArgoCD `yaml:"argocd"`
}

type Code struct {
	Repo string `yaml:"repo"`
	Rev  string `yaml:"rev"`
}

type Helm struct {
	Repo  string `yaml:"repo"`
	Chart string `yaml:"chart"`
	Rev   string `yaml:"rev"`
}

type Jira struct {
	Ticket Ticket `yaml:"ticket"`
}

type Ticket struct {
	CR string `yaml:"cr"`
}

type ArgoCD struct {
	URL    string `yaml:"url"`
	Synced string `yaml:"synced"`
}
