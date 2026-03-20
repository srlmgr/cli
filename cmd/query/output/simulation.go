package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/srlmgr/cli/conversion"
)

// SimulationFormatter handles formatting simulation output in various formats.
type SimulationFormatter struct{}

// NewSimulationFormatter creates a new SimulationFormatter.
func NewSimulationFormatter() *SimulationFormatter {
	return &SimulationFormatter{}
}

// FormatSimulation writes a single simulation to the output in the specified format.
//
//nolint:whitespace // editor/linter issue
func (f *SimulationFormatter) FormatSimulation(
	w io.Writer,
	format string,
	resp *queryv1.GetSimulationResponse,
) error {
	switch strings.ToLower(format) {
	case "json":
		return f.formatSimulationJSON(w, resp)
	case "table":
		return f.formatSimulationTable(w, resp)
	default:
		return fmt.Errorf("unsupported output format %q (supported: table, json)", format)
	}
}

// FormatSimulations writes a list of simulations to the output in the specified format.
//
//nolint:whitespace // editor/linter issue
func (f *SimulationFormatter) FormatSimulations(
	w io.Writer,
	format string,
	resp *queryv1.ListSimulationsResponse,
) error {
	switch strings.ToLower(format) {
	case "json":
		return f.formatSimulationsJSON(w, resp)
	case "table":
		return f.formatSimulationsTable(w, resp)
	default:
		return fmt.Errorf("unsupported output format %q (supported: table, json)", format)
	}
}

//nolint:whitespace // editor/linter issue
func (f *SimulationFormatter) formatSimulationJSON(
	w io.Writer,
	resp *queryv1.GetSimulationResponse,
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

//nolint:whitespace // editor/linter issue
func (f *SimulationFormatter) formatSimulationTable(
	w io.Writer,
	resp *queryv1.GetSimulationResponse,
) error {
	sim := resp.GetSimulation()
	if sim == nil {
		return fmt.Errorf("empty simulation in response")
	}

	if _, err := fmt.Fprintf(w, "ID: %d\n", sim.GetId()); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	if _, err := fmt.Fprintf(w, "Name: %s\n", sim.GetName()); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	if _, err := fmt.Fprintf(w, "Active: %t\n", sim.GetIsActive()); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (f *SimulationFormatter) formatSimulationsJSON(
	w io.Writer,
	resp *queryv1.ListSimulationsResponse,
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

//nolint:whitespace // editor/linter issue
func (f *SimulationFormatter) formatSimulationsTable(
	w io.Writer,
	resp *queryv1.ListSimulationsResponse,
) error {
	if len(resp.GetItems()) == 0 {
		if _, err := fmt.Fprintln(w, "No simulations found."); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "ID\tNAME\tACTIVE\tSUPPORTED_FORMATS"); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, sim := range resp.GetItems() {
		if _, err := fmt.Fprintf(
			tw,
			"%d\t%s\t%t\t%s\n",
			sim.GetId(),
			sim.GetName(),
			sim.GetIsActive(),
			conversion.JoinImportFormats(sim.GetSupportedFormats()),
		); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("flush output: %w", err)
	}

	return nil
}
