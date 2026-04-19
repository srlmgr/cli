package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type TeamStandingsFormatter struct{}

func NewTeamStandingsFormatter() *TeamStandingsFormatter {
	return &TeamStandingsFormatter{}
}

//nolint:whitespace // editor/linter issue
func (f *TeamStandingsFormatter) FormatTeamStandings(
	w io.Writer,
	format string,
	standingsResp *queryv1.GetTeamStandingsResponse,
	teams map[uint32]*commonv1.Team,
) error {
	switch strings.ToLower(format) {
	case JSONOutputFormat:
		return f.formatTeamStandingsJSON(w, standingsResp)
	case TableOutputFormat:
		return f.formatTeamStandingsTable(w, standingsResp, teams)
	default:
		return fmt.Errorf("unsupported output format %q (supported: table, json)", format)
	}
}

//nolint:whitespace // editor/linter issue
func (f *TeamStandingsFormatter) formatTeamStandingsJSON(
	w io.Writer,
	resp *queryv1.GetTeamStandingsResponse,
) error {
	payload, err := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	if _, err = fmt.Fprintln(w, string(payload)); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

//nolint:whitespace,funlen // readability
func (f *TeamStandingsFormatter) formatTeamStandingsTable(
	w io.Writer,
	standingsResp *queryv1.GetTeamStandingsResponse,
	teams map[uint32]*commonv1.Team,
) error {
	standings := standingsResp.GetStandings()
	if len(standings) == 0 {
		if _, err := fmt.Fprintln(w, "No team standings found."); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(
		tw,
		"POS\tTEAM_ID\tTEAM_NAME\tPOINTS\tPREV_POS\tWINS\tPODIUMS\tTOP5\tTOP10",
	); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, standing := range standings {
		teamID := standing.GetTeamId()
		team := teams[teamID]
		teamName := fmt.Sprintf("Team-%d", teamID)
		if team != nil {
			teamName = team.GetName()
		}

		data := standing.GetData()
		if _, err := fmt.Fprintf(
			tw,
			"%d\t%d\t%s\t%d\t%d\t%d\t%d\t%d\t%d\n",
			data.GetPosition(),
			teamID,
			teamName,
			data.GetTotalPoints(),
			data.GetPrevPosition(),
			data.GetNumWins(),
			data.GetNumPodiums(),
			data.GetNumTop5(),
			data.GetNumTop10(),
		); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	if err := tw.Flush(); err != nil {
		return fmt.Errorf("flush output: %w", err)
	}
	return nil
}
