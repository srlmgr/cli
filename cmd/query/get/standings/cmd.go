package standings

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/query/get/standings/driver"
	"github.com/srlmgr/cli/cmd/query/get/standings/team"
)

//nolint:lll // readability
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "standings",
		Short: "Get standings",
		Long:  "Commands for getting driver and team standings from backend.query.v1.QueryService",
	}

	cmd.AddCommand(driver.NewCmd())
	cmd.AddCommand(team.NewCmd())
	return cmd
}
