package simulation

import (
	"context"
	"fmt"
	"io"

	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	"github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/cmd/query/output"
	"github.com/srlmgr/cli/log"
)

func NewCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "simulation",
		Short: "List all simulations",
		Long:  "Fetch simulations from backend.query.v1.QueryService.ListSimulations",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := &listSimulationsCommand{
				apiBaseURL:   config.APIAddr,
				apiToken:     config.APIToken,
				outputFormat: outputFormat,
				out:          cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&outputFormat,
		"output",
		"o",
		"table",
		"Output format (table or json)")
	return cmd
}

type listSimulationsCommand struct {
	apiBaseURL   string
	apiToken     string
	outputFormat string
	out          io.Writer
}

func (c *listSimulationsCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	svc := client.NewQueryServiceClient(c.apiBaseURL, logger)

	resp, err := svc.ListSimulations(
		ctx,
		&queryv1.ListSimulationsRequest{},
	)
	if err != nil {
		return fmt.Errorf("list simulations: %w", err)
	}

	formatter := output.NewSimulationFormatter()
	return formatter.FormatSimulations(c.out, c.outputFormat, resp)
}
