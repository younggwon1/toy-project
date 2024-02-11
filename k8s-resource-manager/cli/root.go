package cli

import (
	"github.com/spf13/cobra"
	"github.com/younggwon1/k8s-resource-manager/cli/manager"
)

var RootCmd = &cobra.Command{
	Use:   "k8s",
	Short: "k8s operations",
}

func init() {
	cobra.EnableCommandSorting = false
	RootCmd.AddCommand(
		manager.Cmd,
	)
}
