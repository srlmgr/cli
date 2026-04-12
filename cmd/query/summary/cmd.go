package summary

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	"github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/conversion"
	"github.com/srlmgr/cli/log"
)

func NewCmd() *cobra.Command {
	var eventID uint32
	var raceID uint32
	var outputFormat string
	var summaryTargetType string
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "get summary for race or event",
		Long:  "Commands for getting summary from backend.query.v1.QueryService",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := &querySummaryCommand{
				apiBaseURL:        config.APIAddr,
				apiToken:          config.APIToken,
				outputFormat:      outputFormat,
				eventID:           eventID,
				raceID:            raceID,
				summaryTargetType: summaryTargetType,
				out:               cmd.OutOrStdout(),
			}
			return runner.run(cmd.Context())
		},
	}
	cmd.Flags().StringVarP(&outputFormat,
		"output",
		"o",
		"table",
		"Output format (table or csv)")
	cmd.Flags().Uint32Var(&raceID,
		"race-id",
		0,
		"ID of a race to get summary for")
	cmd.Flags().Uint32Var(&eventID,
		"event-id",
		0,
		"ID of an event to get summary for")
	cmd.Flags().StringVar(&summaryTargetType,
		"type",
		"driver",
		"Summary target type (driver or team)")
	cmd.MarkFlagsMutuallyExclusive("race-id", "event-id")
	cmd.MarkFlagsOneRequired("event-id", "race-id")
	return cmd
}

type querySummaryCommand struct {
	apiBaseURL        string
	apiToken          string
	outputFormat      string
	eventID, raceID   uint32
	summaryTargetType string
	out               io.Writer
}

//nolint:funlen // much work to do
func (c *querySummaryCommand) run(ctx context.Context) error {
	logger := log.GetFromContext(ctx).Named("rpc")
	svc := client.NewQueryServiceClient(c.apiBaseURL, logger)

	sel := commonv1.SummarySelector{}
	if c.raceID != 0 {
		sel.Scope = &commonv1.SummarySelector_RaceId{RaceId: c.raceID}
	} else if c.eventID != 0 {
		sel.Scope = &commonv1.SummarySelector_EventId{EventId: c.eventID}
	}
	targetType, err := conversion.ParseSummaryTarget(c.summaryTargetType)
	if err != nil {
		return fmt.Errorf("parse summary target type: %w", err)
	}
	sel.Type = targetType
	resp, err := svc.GetSummary(
		ctx,
		connect.NewRequest(&queryv1.GetSummaryRequest{
			Selector: &sel,
		}),
	)
	if err != nil {
		return fmt.Errorf("get summary: %w", err)
	}

	dIDs := lo.Map(resp.Msg.Summaries, func(s *commonv1.Summary, _ int) uint32 {
		return s.GetReferenceId()
	})

	dResp, err := svc.ListDrivers(ctx, connect.NewRequest(&queryv1.ListDriversRequest{
		Filter: &queryv1.ListDriversRequest_MultipleDrivers{
			MultipleDrivers: &queryv1.ListMultipleDrivers{
				DriverIds: dIDs,
			},
		},
	}))
	if err != nil {
		return fmt.Errorf("list drivers: %w", err)
	}

	driverByID := make(map[uint32]*commonv1.Driver, len(dResp.Msg.GetItems()))
	for _, d := range dResp.Msg.GetItems() {
		driverByID[d.GetId()] = d
	}

	switch c.outputFormat {
	case "csv":
		return c.outputCSV(resp.Msg.Summaries, driverByID)
	default:
		return c.outputTable(resp.Msg.Summaries, driverByID)
	}
}

//nolint:whitespace // editor/linter issue
func (c *querySummaryCommand) outputTable(
	summaries []*commonv1.Summary,
	driverByID map[uint32]*commonv1.Driver,
) error {
	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDRIVER\tPOINTS\tBONUS POINTS\tTOTAL POINTS")
	for _, s := range summaries {
		name := resolveDriverName(s.GetReferenceId(), driverByID)
		fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%d\n",
			s.GetReferenceId(),
			name,
			s.GetPoints(),
			s.GetBonusPoints(),
			s.GetTotalPoints(),
		)
	}
	return w.Flush()
}

//nolint:whitespace // editor/linter issue
func (c *querySummaryCommand) outputCSV(
	summaries []*commonv1.Summary,
	driverByID map[uint32]*commonv1.Driver,
) error {
	w := csv.NewWriter(c.out)
	if err := w.Write(
		[]string{"id", "driver", "points", "bonus_points", "total_points"},
	); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	for _, s := range summaries {
		name := resolveDriverName(s.GetReferenceId(), driverByID)
		if err := w.Write([]string{
			strconv.Itoa(int(s.GetReferenceId())),
			name,
			strconv.Itoa(int(s.GetPoints())),
			strconv.Itoa(int(s.GetBonusPoints())),
			strconv.Itoa(int(s.GetTotalPoints())),
		}); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}

func resolveDriverName(id uint32, driverByID map[uint32]*commonv1.Driver) string {
	if d, ok := driverByID[id]; ok {
		return d.GetName()
	}
	return strconv.Itoa(int(id))
}
