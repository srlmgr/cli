package simulation

import (
	"context"
	"fmt"
	"io"
	"strings"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/command/client"
	"github.com/srlmgr/cli/cmd/config"
	"github.com/srlmgr/cli/log"
)

//nolint:lll // readability
func NewCreateCmd() *cobra.Command {
	var name string
	var supportedFormats []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a simulation",
		Long:  "Create a new simulation via backend.command.v1.CommandService.CreateSimulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := &createSimulationCommand{
				name:             name,
				supportedFormats: supportedFormats,
				out:              cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the simulation to create")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	cmd.Flags().StringSliceVar(&supportedFormats, "supported-formats",
		nil, "Supported formats for the simulation")

	return cmd
}

type createSimulationCommand struct {
	name             string
	supportedFormats []string
	out              io.Writer
}

func (c *createSimulationCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	cl := client.NewCommandServiceClient(config.APIAddr, config.APIToken, logger)

	resp, err := cl.CreateSimulation(
		ctx,
		connect.NewRequest(&commandv1.CreateSimulationRequest{
			Name:             c.name,
			IsActive:         true,
			SupportedFormats: c.supportedFormats,
		}),
	)
	if err != nil {
		return fmt.Errorf("create simulation: %w", err)
	}

	sim := resp.Msg.GetSimulation()
	if _, err = fmt.Fprintf(c.out,
		"Created simulation: id=%d name=%s active=%t supported_formats=%s\n",
		sim.GetId(),
		sim.GetName(),
		sim.GetIsActive(),
		strings.Join(sim.GetSupportedFormats(), ", "),
	); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
