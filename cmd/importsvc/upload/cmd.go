package upload

import (
	"context"
	"fmt"
	"io"
	"os"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	importclient "github.com/srlmgr/cli/cmd/importsvc/client"
	queryclient "github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/cmd/support"
	"github.com/srlmgr/cli/conversion"
	"github.com/srlmgr/cli/log"
)

//nolint:lll,funlen // readability
func NewCmd() *cobra.Command {
	var importFormat string
	var externalResultID string
	var filename string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a results file",
		Long:  "Upload a results file via backend.import.v1.ImportService.UploadResultsFile",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			r := support.NewResolve(cmd.Context())
			var resp *queryv1.ResolveSeasonTargetResponse
			var err error
			resp, err = r.AutoResolve(cmd.Flags())
			if err != nil {
				return err
			}

			runner := &uploadCommand{
				raceGridID:       resp.RaceGridId,
				importFormat:     importFormat,
				externalResultID: externalResultID,
				filename:         filename,
				out:              cmd.OutOrStdout(),
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

	cmd.Flags().Uint32("season-id", 0, "ID of the season to upload results for")
	cmd.Flags().Uint32("grid-id", 0, "ID of the race grid to upload results for")
	cmd.Flags().Uint32("event-no", 0, "Number of the event to upload results for")
	cmd.Flags().Uint32("race-no", 0, "Number of the race to upload results for")
	cmd.Flags().Uint32("grid-no", 0, "Number of the race grid to upload results for")
	cmd.Flags().
		StringVar(&importFormat, "import-format", "", "Import format of the file (e.g. json, csv)")
	if err := cmd.MarkFlagRequired("import-format"); err != nil {
		panic(fmt.Sprintf("failed to mark 'import-format' flag as required: %v", err))
	}

	cmd.Flags().StringVar(&filename, "filename", "", "Path to the results file to upload")

	cmd.Flags().
		StringVar(&externalResultID, "external-result-id", "", "External result ID associated with the upload")

	return cmd
}

type queryClient interface {
	GetRace(
		context.Context,
		*connect.Request[queryv1.GetRaceRequest],
	) (*connect.Response[queryv1.GetRaceResponse], error)
}

type importClient interface {
	UploadResultsFile(
		context.Context,
		*connect.Request[importv1.UploadResultsFileRequest],
	) (*connect.Response[importv1.UploadResultsFileResponse], error)
}

type uploadCommand struct {
	raceGridID       uint32
	importFormat     string
	externalResultID string // will be put into payload for import format API
	filename         string
	out              io.Writer
	qrySvc           queryClient
	importSvc        importClient
}

func (c *uploadCommand) run(ctx context.Context) error {
	format, err := conversion.ParseImportFormat(c.importFormat)
	if err != nil {
		return fmt.Errorf("parse import format: %w", err)
	}
	var payload []byte
	if c.filename != "" {
		payload, err = os.ReadFile(c.filename)
		if err != nil {
			return fmt.Errorf("read file %q: %w", c.filename, err)
		}
	}
	if c.externalResultID != "" {
		payload = []byte(c.externalResultID)
	}
	resp, err := c.importSvc.UploadResultsFile(
		ctx,
		connect.NewRequest(&importv1.UploadResultsFileRequest{
			RaceGridId:   c.raceGridID,
			ImportFormat: format,
			Payload:      payload,
		}),
	)
	if err != nil {
		return fmt.Errorf("upload results file: %w", err)
	}

	if _, err = fmt.Fprintf(
		c.out,
		"Uploaded results file: race_grid_id=%d processing_state=%s\n",
		resp.Msg.GetRaceGridId(),
		resp.Msg.GetProcessingState(),
	); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
