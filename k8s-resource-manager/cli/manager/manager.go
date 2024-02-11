package manager

import (
	"github.com/spf13/cobra"
	"github.com/younggwon1/k8s-resource-manager/config/util"
	"github.com/younggwon1/k8s-resource-manager/deployment"
)

var (
	namespace string
)

var Cmd = &cobra.Command{
	Use:   "manager",
	Short: "kuberenetes resource manager operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		k8sCfg := util.KubernetesCredentials() // Init kubernetes config

		err := deployment.ErrorStatus(k8sCfg, namespace)
		if err != nil {
			panic(err)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "flag for namespace to delete")
}
