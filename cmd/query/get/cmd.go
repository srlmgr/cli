package get

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/query/get/simulation"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get commands",
		Long:  "Commands for getting resources from backend.query.v1.QueryService",
	}

	cmd.AddCommand(simulation.NewCmd())
	return cmd
}
