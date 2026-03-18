package simulation

import (
	"context"
	"fmt"
	"io"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/command/client"
	"github.com/srlmgr/cli/cmd/config"
	"github.com/srlmgr/cli/log"
)

//nolint:lll // readability
func NewDeleteCmd() *cobra.Command {
	var simulationID uint32

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a simulation",
		Long:  "Delete a simulation via backend.command.v1.CommandService.DeleteSimulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := &deleteSimulationCommand{
				simulationID: simulationID,
				out:          cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32Var(&simulationID, "simulation-id", 0, "ID of the simulation to delete")
	if err := cmd.MarkFlagRequired("simulation-id"); err != nil {
		panic(err)
	}

	return cmd
}

type deleteSimulationCommand struct {
	simulationID uint32
	out          io.Writer
}

func (c *deleteSimulationCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	cl := client.NewCommandServiceClient(config.APIAddr, config.APIToken, logger)

	_, err := cl.DeleteSimulation(
		ctx,
		connect.NewRequest(&commandv1.DeleteSimulationRequest{
			SimulationId: c.simulationID,
		}),
	)
	if err != nil {
		return fmt.Errorf("delete simulation: %w", err)
	}

	if _, err = fmt.Fprintf(c.out,
		"Deleted simulation: id=%d\n",
		c.simulationID); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
