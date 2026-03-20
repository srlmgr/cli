package setup

import (
	"fmt"

	"github.com/spf13/cobra"

	commandclient "github.com/srlmgr/cli/cmd/command/client"
	"github.com/srlmgr/cli/cmd/config"
	queryclient "github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/log"
)

// NewCmd returns the cobra command for the setup subcommand.
func NewCmd() *cobra.Command {
	var filePath string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup base data from a YAML file",
		Long: "Provision base data (simulations, tracks, car classes, etc.) " +
			"from a YAML configuration file. Uses list-then-create to avoid " +
			"creating duplicates when run multiple times.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("setup")

			runner := &setupRunner{
				filePath: filePath,
				dryRun:   dryRun,
				out:      cmd.OutOrStdout(),
				cmdSvc: commandclient.NewCommandServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
				qrySvc: queryclient.NewQueryServiceClient(
					config.APIAddr, logger,
				),
			}

			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(
		&filePath, "file", "f", "", "Path to the YAML setup file",
	)

	if err := cmd.MarkFlagRequired("file"); err != nil {
		panic(fmt.Sprintf("failed to mark 'file' flag as required: %v", err))
	}

	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false, "Preview actions without making changes",
	)

	return cmd
}
