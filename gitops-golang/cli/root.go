package cli

import (
	"github.com/younggwon1/gitops-golang/cli/deploy"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gogit",
	Short: "gogit operations",
}

func init() {
	cobra.EnableCommandSorting = false
	RootCmd.AddCommand(
		deploy.Cmd,
	)
}
