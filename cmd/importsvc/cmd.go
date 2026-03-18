package importsvc

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "ImportService commands",
		Long:  "Commands for backend.import.v1.ImportService",
	}
}
