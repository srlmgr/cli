package importsvc

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/importsvc/bookings"
	"github.com/srlmgr/cli/cmd/importsvc/preview"
	"github.com/srlmgr/cli/cmd/importsvc/resolve"
	"github.com/srlmgr/cli/cmd/importsvc/upload"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "ImportService commands",
		Long:  "Commands for backend.import.v1.ImportService",
	}

	cmd.AddCommand(upload.NewCmd())
	cmd.AddCommand(resolve.NewCmd())
	cmd.AddCommand(preview.NewCmd())
	cmd.AddCommand(bookings.NewCmd())

	return cmd
}
