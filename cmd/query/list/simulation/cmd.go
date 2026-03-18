package simulation

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/srlmgr/cli/cmd/config"
	"github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/log"
)

func NewCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "simulations",
		Short: "List all simulations",
		Long:  "Fetch simulations from backend.query.v1.QueryService.ListSimulations",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := &listSimulationsCommand{
				apiBaseURL: config.APIAddr,
				apiToken:   config.APIToken,
				output:     output,
				out:        cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&output,
		"output",
		"o",
		"table",
		"Output format (table or json)")
	return cmd
}

type listSimulationsCommand struct {
	apiBaseURL string
	apiToken   string
	output     string
	out        io.Writer
}

func (c *listSimulationsCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	svc := client.NewQueryServiceClient(c.apiBaseURL, logger)

	resp, err := svc.ListSimulations(
		ctx,
		connect.NewRequest(&queryv1.ListSimulationsRequest{}),
	)
	if err != nil {
		return fmt.Errorf("list simulations: %w", err)
	}

	return c.writeOutput(resp.Msg)
}

//nolint:whitespace // editor/linter issue
func (c *listSimulationsCommand) writeOutput(
	resp *queryv1.ListSimulationsResponse,
) error {
	switch strings.ToLower(c.output) {
	case "json":
		return c.writeJSON(resp)
	case "table":
		return c.writeTable(resp)
	default:
		return fmt.Errorf("unsupported output format %q (supported: table, json)",
			c.output)
	}
}

//nolint:whitespace // editor/linter issue
func (c *listSimulationsCommand) writeJSON(
	resp *queryv1.ListSimulationsResponse,
) error {
	payload, err := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}

	if _, err = fmt.Fprintln(c.out, string(payload)); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

//nolint:whitespace // editor/linter issue
func (c *listSimulationsCommand) writeTable(
	resp *queryv1.ListSimulationsResponse,
) error {
	if len(resp.GetItems()) == 0 {
		if _, err := fmt.Fprintln(c.out, "No simulations found."); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "ID\tNAME\tACTIVE\tSUPPORTED_FORMATS"); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, sim := range resp.GetItems() {
		if _, err := fmt.Fprintf(
			w,
			"%d\t%s\t%t\t%s\n",
			sim.GetId(),
			sim.GetName(),
			sim.GetIsActive(),
			strings.Join(sim.GetSupportedFormats(), ", "),
		); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush output: %w", err)
	}

	return nil
}
