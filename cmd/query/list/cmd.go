package list

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/query/list/series"
	"github.com/srlmgr/cli/cmd/query/list/simulation"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List commands",
		Long:  "Commands for listing resources from backend.query.v1.QueryService",
	}

	cmd.AddCommand(simulation.NewCmd())
	cmd.AddCommand(series.NewCmd())
	return cmd
}
