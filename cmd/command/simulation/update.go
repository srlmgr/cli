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
func NewUpdateCmd() *cobra.Command {
	var simulationID uint32
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a simulation",
		Long:  "Update an existing simulation via backend.command.v1.CommandService.UpdateSimulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := &updateSimulationCommand{
				simulationID: simulationID,
				name:         name,
				out:          cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32Var(&simulationID, "simulation-id", 0, "ID of the simulation to update")
	if err := cmd.MarkFlagRequired("simulation-id"); err != nil {
		panic(err)
	}
	cmd.Flags().StringVar(&name, "name", "", "New name for the simulation")

	return cmd
}

type updateSimulationCommand struct {
	simulationID uint32
	name         string
	out          io.Writer
}

func (c *updateSimulationCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	cl := client.NewCommandServiceClient(config.APIAddr, config.APIToken, logger)

	resp, err := cl.UpdateSimulation(
		ctx,
		connect.NewRequest(&commandv1.UpdateSimulationRequest{
			SimulationId: c.simulationID,
			Name:         c.name,
		}),
	)
	if err != nil {
		return fmt.Errorf("update simulation: %w", err)
	}

	sim := resp.Msg.GetSimulation()
	if _, err = fmt.Fprintf(c.out, "Updated simulation: id=%d name=%s active=%t\n",
		sim.GetId(),
		sim.GetName(),
		sim.GetIsActive(),
	); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
