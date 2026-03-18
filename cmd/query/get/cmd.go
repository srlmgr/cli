package get

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get commands",
		Long:  "Commands for getting resources from backend.query.v1.QueryService",
	}

	return cmd
}
