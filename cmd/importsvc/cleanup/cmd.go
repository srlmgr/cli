package cleanup

import (
	"context"
	"fmt"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	importclient "github.com/srlmgr/cli/cmd/importsvc/client"
	"github.com/srlmgr/cli/log"
)

type importClient interface {
	CleanupProcessingData(
		context.Context,
		*connect.Request[importv1.CleanupProcessingDataRequest],
	) (*connect.Response[importv1.CleanupProcessingDataResponse], error)
}

type cleanupCommand struct {
	// cleanup target flags (oneof)
	eventID            uint32
	raceID             uint32
	gridID             uint32
	includeManualEdits bool
	recomputeResults   bool
	importSvc          importClient
}

//nolint:lll // readability
func NewCmd() *cobra.Command {
	var (
		eventID            uint32
		raceID             uint32
		gridID             uint32
		includeManualEdits bool
		recomputeResults   bool
	)

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup processing data",
		Long:  "Cleanup processing data via backend.import.v1.ImportService.CleanupProcessingData",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			runner := &cleanupCommand{
				eventID: eventID,
				raceID:  raceID,
				gridID:  gridID,
				importSvc: importclient.NewImportServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
				includeManualEdits: includeManualEdits,
				recomputeResults:   recomputeResults,
			}
			return runner.run(cmd.Context())
		},
	}

	// target flags (mutually exclusive)
	cmd.Flags().Uint32Var(&eventID, "event-id", 0, "ID of the event to apply penalty to")
	cmd.Flags().Uint32Var(&raceID, "race-id", 0, "ID of the race to apply penalty to")
	cmd.Flags().Uint32Var(&gridID, "grid-id", 0, "ID of the race grid to apply penalty to")
	cmd.Flags().BoolVar(&includeManualEdits, "include-manuals", false, "Include manual edits in the cleanup")
	cmd.Flags().BoolVar(&recomputeResults, "recompute", false, "Recompute results after cleanup")
	return cmd
}

//nolint:funlen,lll // many parameters
func (c *cleanupCommand) run(ctx context.Context) error {
	// Validate target flags (exactly one must be set)
	targetFlags := 0
	if c.eventID != 0 {
		targetFlags++
	}
	if c.raceID != 0 {
		targetFlags++
	}
	if c.gridID != 0 {
		targetFlags++
	}
	if targetFlags != 1 {
		return fmt.Errorf(
			"exactly one target flag must be specified (--event-id, --race-id, or --grid-id)",
		)
	}

	// Build cleanup target
	target := &importv1.CleanupTarget{}

	// Set scope
	switch {
	case c.eventID != 0:
		target.SetEventId(c.eventID)
	case c.raceID != 0:
		target.SetRaceId(c.raceID)
	case c.gridID != 0:
		target.SetRaceGridId(c.gridID)
	}

	req := &importv1.CleanupProcessingDataRequest{
		CleanupTarget:      target,
		IncludeManualEdits: c.includeManualEdits,
		PerformRecompute:   c.recomputeResults,
	}

	_, err := c.importSvc.CleanupProcessingData(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("cleanup processing data: %w", err)
	}

	return nil
}
