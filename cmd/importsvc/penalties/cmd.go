package penalties

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/importsvc/penalties/addpenalty"
	"github.com/srlmgr/cli/cmd/importsvc/penalties/deletepenalty"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "penalties",
		Short: "Penalties commands for ImportService",
		Long:  "Penalties commands for backend.import.v1.ImportService",
	}

	cmd.AddCommand(addpenalty.NewCmd())
	cmd.AddCommand(deletepenalty.NewCmd())

	return cmd
}
