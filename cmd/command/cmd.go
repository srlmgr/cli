package command

import (
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/command/simulation"
)

func NewCmd() *cobra.Command {
	commandCmd := &cobra.Command{
		Use:   "command",
		Short: "CommandService commands",
		Long:  "Commands for backend.command.v1.CommandService",
	}

	commandCmd.AddCommand(simulation.NewSimulationCmd())

	return commandCmd
}
