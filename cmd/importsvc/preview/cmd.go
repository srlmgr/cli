package preview

import (
	"context"
	"fmt"
	"io"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	importclient "github.com/srlmgr/cli/cmd/importsvc/client"
	"github.com/srlmgr/cli/cmd/config"
	queryclient "github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/log"
)

func NewCmd() *cobra.Command {
	var raceID uint32

	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview preprocessed results for a race",
		Long:  "Fetch a preprocessing preview via backend.import.v1.ImportService.GetPreprocessPreview",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			runner := &previewCommand{
				raceID: raceID,
				out:    cmd.OutOrStdout(),
				qrySvc: queryclient.NewQueryServiceClient(
					config.APIAddr, logger,
				),
				importSvc: importclient.NewImportServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32Var(&raceID, "race-id", 0, "ID of the race to preview results for")
	if err := cmd.MarkFlagRequired("race-id"); err != nil {
		panic(fmt.Sprintf("failed to mark 'race-id' flag as required: %v", err))
	}

	return cmd
}

type queryClient interface {
	GetRace(context.Context, *connect.Request[queryv1.GetRaceRequest]) (*connect.Response[queryv1.GetRaceResponse], error)
}

type importClient interface {
	GetPreprocessPreview(context.Context, *connect.Request[importv1.GetPreprocessPreviewRequest]) (*connect.Response[importv1.GetPreprocessPreviewResponse], error)
}

type previewCommand struct {
	raceID    uint32
	out       io.Writer
	qrySvc    queryClient
	importSvc importClient
}

func (c *previewCommand) run(ctx context.Context) error {
	raceResp, err := c.qrySvc.GetRace(
		ctx,
		connect.NewRequest(&queryv1.GetRaceRequest{Id: c.raceID}),
	)
	if err != nil {
		return fmt.Errorf("get race: %w", err)
	}

	eventID := raceResp.Msg.GetRace().GetEventId()

	resp, err := c.importSvc.GetPreprocessPreview(
		ctx,
		connect.NewRequest(&importv1.GetPreprocessPreviewRequest{
			EventId: eventID,
			RaceId:  c.raceID,
		}),
	)
	if err != nil {
		return fmt.Errorf("get preprocess preview: %w", err)
	}

	rows := resp.Msg.GetRows()
	unresolvedMappings := resp.Msg.GetUnresolvedMappings()

	if _, err = fmt.Fprintf(c.out, "Preview: rows=%d unresolved_mappings=%d\n",
		len(rows), len(unresolvedMappings),
	); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	for _, row := range rows {
		if _, err = fmt.Fprintf(c.out,
			"  row: id=%d race_id=%d driver_id=%d car_model_id=%d position=%d laps=%d\n",
			row.GetId(),
			row.GetRaceId(),
			row.GetDriverId(),
			row.GetCarModelId(),
			row.GetFinishingPosition(),
			row.GetCompletedLaps(),
		); err != nil {
			return fmt.Errorf("write row output: %w", err)
		}
	}

	for _, m := range unresolvedMappings {
		if _, err = fmt.Fprintf(c.out,
			"  unresolved: source_value=%s mapping_type=%s\n",
			m.GetSourceValue(),
			m.GetMappingType(),
		); err != nil {
			return fmt.Errorf("write unresolved mapping output: %w", err)
		}
	}

	return nil
}
