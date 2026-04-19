package get

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/query/get/simulation"
	"github.com/srlmgr/cli/cmd/query/get/standings"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get commands",
		Long:  "Commands for getting resources from backend.query.v1.QueryService",
	}

	cmd.AddCommand(simulation.NewCmd())
	cmd.AddCommand(standings.NewCmd())
	return cmd
}
