package bookings

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/importsvc/bookings/auto"
	"github.com/srlmgr/cli/cmd/importsvc/bookings/driver"
	"github.com/srlmgr/cli/cmd/importsvc/bookings/team"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bookings",
		Short: "Bookings commands for ImportService",
		Long:  "Bookings commands for backend.import.v1.ImportService",
	}

	cmd.AddCommand(auto.NewCmd())
	cmd.AddCommand(driver.NewCmd())
	cmd.AddCommand(team.NewCmd())

	return cmd
}
