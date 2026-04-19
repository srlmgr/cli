package team

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
	"github.com/srlmgr/cli/log"
)

func NewCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "team",
		Short: "Get team standings",
		Long:  "Fetch team standings from backend.query.v1.QueryService.GetTeamStandings",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := cmd.Flags().GetUint32("event-id")
			if err != nil {
				return fmt.Errorf("parse event-id flag: %w", err)
			}

			runner := &getTeamStandingsCommand{
				apiBaseURL:   config.APIAddr,
				apiToken:     config.APIToken,
				eventID:      eventID,
				outputFormat: outputFormat,
				out:          cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32("event-id", 0, "ID of the event to get team standings for")
	if err := cmd.MarkFlagRequired("event-id"); err != nil {
		panic(err)
	}
	cmd.Flags().StringVarP(&outputFormat,
		"output",
		"o",
		"table",
		"Output format (table or json)")
	return cmd
}

type getTeamStandingsCommand struct {
	apiBaseURL   string
	apiToken     string
	eventID      uint32
	outputFormat string
	out          io.Writer
}

func (c *getTeamStandingsCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	svc := client.NewQueryServiceClient(c.apiBaseURL, logger)

	// Get the team standings
	standingsResp, err := svc.GetTeamStandings(
		ctx,
		connect.NewRequest(&queryv1.GetTeamStandingsRequest{
			EventId: c.eventID,
		}),
	)
	if err != nil {
		return fmt.Errorf("get team standings: %w", err)
	}

	// Resolve team names - first we need to get the season ID from the event
	// For simplicity, we'll make individual GetTeam calls for each team
	// since ListTeams only supports filtering by season_id currently
	teams := make(map[uint32]*commonv1.Team)
	for _, standing := range standingsResp.Msg.GetStandings() {
		teamID := standing.GetTeamId()
		if _, exists := teams[teamID]; !exists {
			teamResp, err := svc.GetTeam(
				ctx,
				connect.NewRequest(&queryv1.GetTeamRequest{
					Id: teamID,
				}),
			)
			if err != nil {
				// If we can't get the team name, we'll just use the ID
				continue
			}
			teams[teamID] = teamResp.Msg.GetTeam()
		}
	}

	formatter := output.NewTeamStandingsFormatter()
	return formatter.FormatTeamStandings(c.out, c.outputFormat, standingsResp.Msg, teams)
}
