package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type DriverStandingsFormatter struct{}

func NewDriverStandingsFormatter() *DriverStandingsFormatter {
	return &DriverStandingsFormatter{}
}

//nolint:whitespace // editor/linter issue
func (f *DriverStandingsFormatter) FormatDriverStandings(
	w io.Writer,
	format string,
	standingsResp *queryv1.GetDriverStandingsResponse,
	driversResp *queryv1.ListDriversResponse,
) error {
	switch strings.ToLower(format) {
	case JSONOutputFormat:
		return f.formatDriverStandingsJSON(w, standingsResp)
	case TableOutputFormat:
		return f.formatDriverStandingsTable(w, standingsResp, driversResp)
	default:
		return fmt.Errorf("unsupported output format %q (supported: table, json)", format)
	}
}

//nolint:whitespace // editor/linter issue
func (f *DriverStandingsFormatter) formatDriverStandingsJSON(
	w io.Writer,
	resp *queryv1.GetDriverStandingsResponse,
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

//nolint:whitespace,funlen,lll // readability
func (f *DriverStandingsFormatter) formatDriverStandingsTable(
	w io.Writer,
	standingsResp *queryv1.GetDriverStandingsResponse,
	driversResp *queryv1.ListDriversResponse,
) error {
	standings := standingsResp.GetStandings()
	if len(standings) == 0 {
		if _, err := fmt.Fprintln(w, "No driver standings found."); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}

	// Create a map of driver ID to driver name for quick lookup
	driverNames := make(map[uint32]string)
	if driversResp != nil {
		for _, driver := range driversResp.GetItems() {
			driverNames[driver.GetId()] = driver.GetName()
		}
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(
		tw,
		"POS\tPREV_POS\tDRIVER_ID\tDRIVER_NAME\tPOINTS\tEVENTS\tRACES\tWINS\tPODIUMS\tTOP5\tTOP10",
	); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, standing := range standings {
		driverID := standing.GetDriverId()
		driverName := driverNames[driverID]
		if driverName == "" {
			driverName = fmt.Sprintf("Driver-%d", driverID)
		}

		data := standing.GetData()
		if _, err := fmt.Fprintf(
			tw,
			"%d\t%d\t%d\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			data.GetPosition(),
			data.GetPrevPosition(),
			driverID,
			driverName,
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
