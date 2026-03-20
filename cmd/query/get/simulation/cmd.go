package simulation

import (
	"context"
	"fmt"
	"io"

	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
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
		Short: "Get a simulation",
		Long:  "Fetch a simulation from backend.query.v1.QueryService.GetSimulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			simulationID, err := cmd.Flags().GetUint32("simulation-id")
			if err != nil {
				return fmt.Errorf("parse simulation-id flag: %w", err)
			}

			runner := &getSimulationCommand{
				apiBaseURL:   config.APIAddr,
				apiToken:     config.APIToken,
				simulationID: simulationID,
				outputFormat: outputFormat,
				out:          cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32("simulation-id", 0, "ID of the simulation to get")
	if err := cmd.MarkFlagRequired("simulation-id"); err != nil {
		panic(err)
	}
	cmd.Flags().StringVarP(&outputFormat,
		"output",
		"o",
		"table",
		"Output format (table or json)")
	return cmd
}

type getSimulationCommand struct {
	apiBaseURL   string
	apiToken     string
	simulationID uint32
	outputFormat string
	out          io.Writer
}

func (c *getSimulationCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	svc := client.NewQueryServiceClient(c.apiBaseURL, logger)

	resp, err := svc.GetSimulation(
		ctx,
		connect.NewRequest(&queryv1.GetSimulationRequest{
			Id: c.simulationID,
		}),
	)
	if err != nil {
		return fmt.Errorf("get simulation: %w", err)
	}

	formatter := output.NewSimulationFormatter()
	return formatter.FormatSimulation(c.out, c.outputFormat, resp.Msg)
}
