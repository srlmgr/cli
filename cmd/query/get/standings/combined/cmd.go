package combined

import (
	"context"
	"fmt"
	"io"

	"buf.build/gen/go/srlmgr/api/connectrpc/go/backend/query/v1/queryv1connect"
	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/samber/lo"
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
		Use:   "combined",
		Short: "Get combined standings (driver and team)",
		Long:  "Fetch combined standings from backend.query.v1.QueryService.GetStandings",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := cmd.Flags().GetUint32("event-id")
			if err != nil {
				return fmt.Errorf("parse event-id flag: %w", err)
			}

			runner := &getCombinedStandingsCommand{
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

	cmd.Flags().Uint32("event-id", 0, "ID of the event to get combined standings for")
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

type getCombinedStandingsCommand struct {
	apiBaseURL   string
	apiToken     string
	eventID      uint32
	outputFormat string
	skipMode     string
	out          io.Writer
	querySvc     queryv1connect.QueryServiceClient
	event        *commonv1.Event
	season       *commonv1.Season
}

//nolint:funlen // lot of things to do
func (c *getCombinedStandingsCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	c.querySvc = client.NewQueryServiceClient(c.apiBaseURL, logger)
	standingsSvc := client.NewStandingsServiceClient(c.apiBaseURL, logger)

	// Get the driver standings
	standingsResp, err := standingsSvc.GetStandings(
		ctx,
		connect.NewRequest(&queryv1.GetStandingsRequest{
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

	if err = c.loadEvent(ctx); err != nil {
		return fmt.Errorf("load event: %w", err)
	}
	if err = c.loadSeason(ctx); err != nil {
		return fmt.Errorf("load season: %w", err)
	}

	if len(standingsResp.Msg.GetPrimaryStandings()) == 0 {
		if _, err := fmt.Fprintln(c.out, "No primary standings found."); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}
	var primaryResolver output.IDResolver
	//nolint:exhaustive // we only handle driver and team standings for now
	switch standingsResp.Msg.GetPrimaryStandings()[0].StandingsType {
	case queryv1.StandingsType_STANDINGS_TYPE_DRIVER:
		resolver, err := c.driverResolver(ctx, standingsResp.Msg.GetPrimaryStandings())
		if err != nil {
			return fmt.Errorf("resolve driver names: %w", err)
		}
		primaryResolver = resolver

	case queryv1.StandingsType_STANDINGS_TYPE_TEAM:
		resolver, err := c.teamResolver(ctx, standingsResp.Msg.GetPrimaryStandings())
		if err != nil {
			return fmt.Errorf("resolve team names: %w", err)
		}
		primaryResolver = resolver
	}

	formatter := output.NewPrimaryStandingsFormatter()
	return formatter.FormatPrimaryStandings(
		c.out, c.outputFormat, standingsResp.Msg, primaryResolver)
}

//nolint:whitespace // readability
func (c *getCombinedStandingsCommand) driverResolver(
	ctx context.Context,
	standings []*queryv1.Standing,
) (output.IDResolver, error) {
	driverIDs := make([]uint32, len(standings))
	for i, standing := range standings {
		driverIDs[i] = standing.GetReferenceId()
	}

	// Resolve driver names
	var driversResp *queryv1.ListDriversResponse
	if len(driverIDs) > 0 {
		resp, err := c.querySvc.ListDrivers(
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
			return nil, fmt.Errorf("list drivers: %w", err)
		}
		driversResp = resp.Msg
	}
	lookup := lo.SliceToMap(driversResp.GetItems(),
		func(item *commonv1.Driver) (uint32, string) {
			return item.GetId(), item.GetName()
		})
	return func(id uint32) string {
		ret, _ := lo.Coalesce(lookup[id], "n.a.")
		return ret
	}, nil
}

//nolint:whitespace // readability
func (c *getCombinedStandingsCommand) teamResolver(
	ctx context.Context,
	standings []*queryv1.Standing,
) (output.IDResolver, error) {
	teamIDs := make([]uint32, len(standings))
	for i, standing := range standings {
		teamIDs[i] = standing.GetReferenceId()
	}

	// Resolve team names
	var teamsResp *queryv1.ListTeamsResponse
	if len(teamIDs) > 0 {
		resp, err := c.querySvc.ListTeams(
			ctx,
			connect.NewRequest(&queryv1.ListTeamsRequest{
				SeasonId: c.season.GetId(),
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("list teams: %w", err)
		}
		teamsResp = resp.Msg
	}
	lookup := lo.SliceToMap(teamsResp.GetItems(),
		func(item *commonv1.Team) (uint32, string) {
			return item.GetId(), item.GetName()
		})
	return func(id uint32) string {
		ret, _ := lo.Coalesce(lookup[id], "n.a.")
		return ret
	}, nil
}

func (c *getCombinedStandingsCommand) loadEvent(ctx context.Context) (err error) {
	var resp *connect.Response[queryv1.GetEventResponse]
	resp, err = c.querySvc.GetEvent(ctx, connect.NewRequest(&queryv1.GetEventRequest{
		Id: c.eventID,
	}))
	if err != nil {
		return fmt.Errorf("get event: %w", err)
	}
	c.event = resp.Msg.GetEvent()
	return nil
}

func (c *getCombinedStandingsCommand) loadSeason(ctx context.Context) (err error) {
	var resp *connect.Response[queryv1.GetSeasonResponse]
	resp, err = c.querySvc.GetSeason(ctx, connect.NewRequest(&queryv1.GetSeasonRequest{
		Id: c.event.GetSeasonId(),
	}))
	if err != nil {
		return fmt.Errorf("get season: %w", err)
	}
	c.season = resp.Msg.GetSeason()
	return nil
}
