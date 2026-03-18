package series

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "series",
		Short: "List series commands (placeholder)",
		Long:  "Commands for listing series from backend.query.v1.QueryService",
	}

	return cmd
}
