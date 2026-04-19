package addpenalty

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

type importClient interface {
	AddPenalty(
		context.Context,
		*connect.Request[importv1.AddPenaltyRequest],
	) (*connect.Response[importv1.AddPenaltyResponse], error)
}

type addPenaltyCommand struct {
	// Scope flags (oneof)
	eventID uint32
	raceID  uint32
	gridID  uint32

	// Target flags (oneof)
	driverID uint32
	teamID   uint32

	// Penalty details
	penaltyPoints int32
	reason        string

	out       io.Writer
	importSvc importClient
}

//nolint:funlen // many parameters
func NewCmd() *cobra.Command {
	var (
		eventID       uint32
		raceID        uint32
		gridID        uint32
		driverID      uint32
		teamID        uint32
		penaltyPoints int32
		reason        string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a penalty",
		Long:  "Add a penalty via backend.import.v1.ImportService.AddPenalty",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			runner := &addPenaltyCommand{
				eventID:       eventID,
				raceID:        raceID,
				gridID:        gridID,
				driverID:      driverID,
				teamID:        teamID,
				penaltyPoints: penaltyPoints,
				reason:        reason,
				out:           cmd.OutOrStdout(),
				importSvc: importclient.NewImportServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
			}
			return runner.run(cmd.Context())
		},
	}

	// Scope flags (mutually exclusive)
	cmd.Flags().Uint32Var(&eventID, "event-id", 0, "ID of the event to apply penalty to")
	cmd.Flags().Uint32Var(&raceID, "race-id", 0, "ID of the race to apply penalty to")
	cmd.Flags().Uint32Var(&gridID, "grid-id", 0, "ID of the race grid to apply penalty to")

	// Target flags (mutually exclusive)
	cmd.Flags().Uint32Var(&driverID, "driver-id", 0, "ID of the driver to penalize")
	cmd.Flags().Uint32Var(&teamID, "team-id", 0, "ID of the team to penalize")

	// Penalty details
	cmd.Flags().Int32Var(&penaltyPoints, "penalty-points", 0, "Number of penalty points")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for the penalty")

	if err := cmd.MarkFlagRequired("penalty-points"); err != nil {
		panic(fmt.Sprintf("failed to mark 'penalty-points' flag as required: %v", err))
	}
	if err := cmd.MarkFlagRequired("reason"); err != nil {
		panic(fmt.Sprintf("failed to mark 'reason' flag as required: %v", err))
	}

	return cmd
}

//nolint:funlen,lll // many parameters
func (c *addPenaltyCommand) run(ctx context.Context) error {
	// Validate scope flags (exactly one must be set)
	scopeFlags := 0
	if c.eventID != 0 {
		scopeFlags++
	}
	if c.raceID != 0 {
		scopeFlags++
	}
	if c.gridID != 0 {
		scopeFlags++
	}
	if scopeFlags != 1 {
		return fmt.Errorf(
			"exactly one scope flag must be specified (--event-id, --race-id, or --grid-id)",
		)
	}

	// Validate target flags (exactly one must be set)
	targetFlags := 0
	if c.driverID != 0 {
		targetFlags++
	}
	if c.teamID != 0 {
		targetFlags++
	}
	if targetFlags != 1 {
		return fmt.Errorf("exactly one target flag must be specified (--driver-id or --team-id)")
	}

	// Build penalty target
	target := &importv1.PenaltyTarget{}

	// Set scope
	switch {
	case c.eventID != 0:
		target.SetEventId(c.eventID)
	case c.raceID != 0:
		target.SetRaceId(c.raceID)
	case c.gridID != 0:
		target.SetRaceGridId(c.gridID)
	}

	// Set target
	switch {
	case c.driverID != 0:
		target.SetDriverId(c.driverID)
	case c.teamID != 0:
		target.SetTeamId(c.teamID)
	}

	req := &importv1.AddPenaltyRequest{
		Target:        target,
		PenaltyPoints: c.penaltyPoints,
		Reason:        c.reason,
	}

	_, err := c.importSvc.AddPenalty(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("add penalty: %w", err)
	}

	fmt.Fprintf(c.out, "Successfully added penalty: %d points for %s\n",
		c.penaltyPoints, c.reason)
	return nil
}
