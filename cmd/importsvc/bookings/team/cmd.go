package team

import (
	"context"
	"fmt"
	"io"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	importclient "github.com/srlmgr/cli/cmd/importsvc/client"
	"github.com/srlmgr/cli/log"
)

//nolint:lll // readability
func NewCmd() *cobra.Command {
	var eventID uint32
	cmd := &cobra.Command{
		Use:   "team",
		Short: "create bookings for teams based on an import batch",
		Long:  "Create bookings for teams via backend.import.v1.ImportService.CreateTeamBookings",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			runner := &resolveCommand{
				eventID: eventID,
				out:     cmd.OutOrStdout(),
				importSvc: importclient.NewImportServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32Var(&eventID, "event-id", 0, "ID of the event to create bookings for")
	if err := cmd.MarkFlagRequired("event-id"); err != nil {
		panic(fmt.Sprintf("failed to mark 'event-id' flag as required: %v", err))
	}

	return cmd
}

type importClient interface {
	ComputeTeamBookingEntries(
		context.Context,
		*connect.Request[importv1.ComputeTeamBookingEntriesRequest],
	) (*connect.Response[importv1.ComputeTeamBookingEntriesResponse], error)
}

type resolveCommand struct {
	eventID   uint32
	out       io.Writer
	importSvc importClient
}

func (c *resolveCommand) run(ctx context.Context) error {
	resp, err := c.importSvc.ComputeTeamBookingEntries(
		ctx,
		connect.NewRequest(&importv1.ComputeTeamBookingEntriesRequest{
			EventId: c.eventID,
		}),
	)
	if err != nil {
		return fmt.Errorf("compute team booking entries: %w", err)
	}

	if _, err = fmt.Fprintf(
		c.out,
		"Computed team booking entries: created_entries=%d\n",
		resp.Msg.GetCreatedEntries(),
	); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
