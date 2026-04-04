package bookings

import (
	"github.com/spf13/cobra"
	"github.com/srlmgr/cli/cmd/importsvc/bookings/driver"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bookings",
		Short: "Bookings commands for ImportService",
		Long:  "Bookings commands for backend.import.v1.ImportService",
	}

	cmd.AddCommand(driver.NewCmd())

	return cmd
}
