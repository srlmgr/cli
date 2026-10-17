package preview

import (
	"context"
	"fmt"
	"io"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	importclient "github.com/srlmgr/cli/cmd/importsvc/client"
	queryclient "github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/cmd/support"
	"github.com/srlmgr/cli/log"
)

//nolint:lll // readability
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview preprocessed results for a race",
		Long:  "Fetch a preprocessing preview via backend.import.v1.ImportService.GetPreprocessPreview",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")
			r := support.NewResolve(cmd.Context())
			var resp *queryv1.ResolveSeasonTargetResponse
			var err error
			resp, err = r.AutoResolve(cmd.Flags())
			if err != nil {
				return err
			}
			runner := &previewCommand{
				raceGridID: resp.RaceGridId,
				out:        cmd.OutOrStdout(),
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
	cmd.Flags().Uint32("season-id", 0, "ID of the season to view preview for")
	cmd.Flags().Uint32("grid-id", 0, "ID of the race grid to view preview for")
	cmd.Flags().Uint32("event-no", 0, "Number of the event to view preview for")
	cmd.Flags().Uint32("race-no", 0, "Number of the race to view preview for")
	cmd.Flags().Uint32("grid-no", 0, "Number of the race grid to view preview for")

	return cmd
}

type queryClient interface {
	GetRace(
		context.Context,
		*connect.Request[queryv1.GetRaceRequest],
	) (*connect.Response[queryv1.GetRaceResponse], error)
}

type importClient interface {
	GetPreprocessPreview(
		context.Context,
		*connect.Request[importv1.GetPreprocessPreviewRequest],
	) (*connect.Response[importv1.GetPreprocessPreviewResponse], error)
}

type previewCommand struct {
	raceGridID uint32
	out        io.Writer
	qrySvc     queryClient
	importSvc  importClient
}

func (c *previewCommand) run(ctx context.Context) error {
	resp, err := c.importSvc.GetPreprocessPreview(
		ctx,
		connect.NewRequest(&importv1.GetPreprocessPreviewRequest{
			RaceGridId: c.raceGridID,
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
		//nolint:lll // readability
		if _, err = fmt.Fprintf(
			c.out,
			"  row: id=%d race_grid_id=%d driver_id=%d car_model_variant_id=%d position=%d laps=%d\n",
			row.GetId(),
			row.GetRaceGridId(),
			row.GetDriverId(),
			row.GetCarModelVariantId(),
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
