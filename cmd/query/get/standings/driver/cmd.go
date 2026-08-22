package driver

import (
	"context"
	"fmt"
	"io"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	"github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/cmd/query/output"
	"github.com/srlmgr/cli/conversion"
	"github.com/srlmgr/cli/log"
)

func NewCmd() *cobra.Command {
	var outputFormat string
	var skipMode string

	cmd := &cobra.Command{
		Use:   "driver",
		Short: "Get driver standings",
		Long:  "Fetch driver standings from backend.query.v1.QueryService.GetDriverStandings",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := cmd.Flags().GetUint32("event-id")
			if err != nil {
				return fmt.Errorf("parse event-id flag: %w", err)
			}

			runner := &getDriverStandingsCommand{
				apiBaseURL:   config.APIAddr,
				apiToken:     config.APIToken,
				eventID:      eventID,
				outputFormat: outputFormat,
				skipMode:     skipMode,
				out:          cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32("event-id", 0, "ID of the event to get driver standings for")
	if err := cmd.MarkFlagRequired("event-id"); err != nil {
		panic(err)
	}
	cmd.Flags().StringVarP(&outputFormat,
		"output",
		"o",
		"table",
		"Output format (table or json)")
	cmd.Flags().StringVarP(&skipMode,
		"skip-mode",
		"s",
		"never",
		"Skip mode (never, always, when-applicable)")
	return cmd
}

type getDriverStandingsCommand struct {
	apiBaseURL   string
	apiToken     string
	eventID      uint32
	outputFormat string
	skipMode     string
	out          io.Writer
}

//nolint:funlen // lot of things to do
func (c *getDriverStandingsCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	svc := client.NewQueryServiceClient(c.apiBaseURL, logger)

	// Get the driver standings
	standingsResp, err := svc.GetDriverStandings(
		ctx,
		connect.NewRequest(&queryv1.GetDriverStandingsRequest{
			EventId: c.eventID,
			SkipMode: func() commonv1.SkipMode {
				skipMode, err := conversion.ParseSkipMode(c.skipMode)
				if err != nil {
					panic(err)
				}
				return skipMode
			}(),
		}),
	)
	if err != nil {
		return fmt.Errorf("get driver standings: %w", err)
	}

	// Extract driver IDs for name resolution
	driverIDs := make([]uint32, len(standingsResp.Msg.GetStandings()))
	for i, standing := range standingsResp.Msg.GetStandings() {
		driverIDs[i] = standing.GetDriverId()
	}

	// Resolve driver names
	var driversResp *queryv1.ListDriversResponse
	if len(driverIDs) > 0 {
		resp, err := svc.ListDrivers(
			ctx,
			connect.NewRequest(&queryv1.ListDriversRequest{
				Filter: &queryv1.ListDriversRequest_MultipleDrivers{
					MultipleDrivers: &queryv1.ListMultipleDrivers{
						DriverIds: driverIDs,
					},
				},
			}),
		)
		if err != nil {
			return fmt.Errorf("list drivers: %w", err)
		}
		driversResp = resp.Msg
	}

	formatter := output.NewDriverStandingsFormatter()
	return formatter.FormatDriverStandings(
		c.out, c.outputFormat, standingsResp.Msg, driversResp)
}
