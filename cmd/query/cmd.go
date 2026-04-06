package query

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/query/get"
	"github.com/srlmgr/cli/cmd/query/list"
	"github.com/srlmgr/cli/cmd/query/summary"
)

func NewCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "QueryService commands",
		Long:  "Commands for backend.query.v1.QueryService",
	}

	queryCmd.AddCommand(list.NewCmd())
	queryCmd.AddCommand(get.NewCmd())
	queryCmd.AddCommand(summary.NewCmd())
	return queryCmd
}
