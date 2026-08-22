package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	PrimaryStandingsFormatter struct{}
	IDResolver                func(id uint32) string
)

func NewPrimaryStandingsFormatter() *PrimaryStandingsFormatter {
	return &PrimaryStandingsFormatter{}
}

//nolint:whitespace // editor/linter issue
func (f *PrimaryStandingsFormatter) FormatPrimaryStandings(
	w io.Writer,
	format string,
	standingsResp *queryv1.GetStandingsResponse,
	idResolver IDResolver,
) error {
	switch strings.ToLower(format) {
	case JSONOutputFormat:
		return f.formatPrimaryStandingsJSON(w, standingsResp)
	case TableOutputFormat:
		return f.formatPrimaryStandingsTable(w, standingsResp, idResolver)
	default:
		return fmt.Errorf("unsupported output format %q (supported: table, json)", format)
	}
}

//nolint:whitespace // editor/linter issue
func (f *PrimaryStandingsFormatter) formatPrimaryStandingsJSON(
	w io.Writer,
	resp *queryv1.GetStandingsResponse,
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
func (f *PrimaryStandingsFormatter) formatPrimaryStandingsTable(
	w io.Writer,
	standingsResp *queryv1.GetStandingsResponse,
	idResolver IDResolver,
) error {
	standings := standingsResp.GetPrimaryStandings()
	if len(standings) == 0 {
		if _, err := fmt.Fprintln(w, "No primary standings found."); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(
		tw,
		"POS\tPREV_POS\tREF_ID\tNAME\tPOINTS\tEVENTS\tRACES\tWINS\tPODIUMS\tTOP5\tTOP10",
	); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, standing := range standings {

		data := standing.GetData()
		if _, err := fmt.Fprintf(
			tw,
			"%d\t%d\t%d\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			data.GetPosition(),
			data.GetPrevPosition(),
			standing.GetReferenceId(),
			idResolver(standing.GetReferenceId()),
			data.GetTotalPoints(),
			data.GetNumEvents(),
			data.GetNumRaces(),
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
