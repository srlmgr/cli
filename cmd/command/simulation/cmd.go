package simulation

import "github.com/spf13/cobra"

func NewSimulationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulation",
		Short: "Simulation commands",
		Long:  "Commands for managing simulations via backend.command.v1.CommandService",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())
	return cmd
}
