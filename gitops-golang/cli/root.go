package cli

import (
	"github.com/spf13/cobra"

	"github.com/younggwon1/gitops-golang/cli/deploy"
	"github.com/younggwon1/gitops-golang/cli/sync"
)

var RootCmd = &cobra.Command{
	Use:   "gogit",
	Short: "gogit operations",
}

func init() {
	cobra.EnableCommandSorting = false
	RootCmd.AddCommand(
		deploy.Cmd,
		sync.Cmd,
	)
}
