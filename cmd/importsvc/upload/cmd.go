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
	"github.com/srlmgr/cli/conversion"
	"github.com/srlmgr/cli/log"
)

//nolint:lll,funlen // readability
func NewCmd() *cobra.Command {
	var gridID uint32
	var importFormat string
	var filename string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a results file",
		Long:  "Upload a results file via backend.import.v1.ImportService.UploadResultsFile",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			payload, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("read file %q: %w", filename, err)
			}

			runner := &uploadCommand{
				raceGridID:   gridID,
				importFormat: importFormat,
				payload:      payload,
				out:          cmd.OutOrStdout(),
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

	cmd.Flags().Uint32Var(&gridID, "grid-id", 0, "ID of the race grid to upload results for")
	if err := cmd.MarkFlagRequired("grid-id"); err != nil {
		panic(fmt.Sprintf("failed to mark 'grid-id' flag as required: %v", err))
	}

	cmd.Flags().
		StringVar(&importFormat, "import-format", "", "Import format of the file (e.g. json, csv)")
	if err := cmd.MarkFlagRequired("import-format"); err != nil {
		panic(fmt.Sprintf("failed to mark 'import-format' flag as required: %v", err))
	}

	cmd.Flags().StringVar(&filename, "filename", "", "Path to the results file to upload")
	if err := cmd.MarkFlagRequired("filename"); err != nil {
		panic(fmt.Sprintf("failed to mark 'filename' flag as required: %v", err))
	}

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
	raceGridID   uint32
	importFormat string
	payload      []byte
	out          io.Writer
	qrySvc       queryClient
	importSvc    importClient
}

func (c *uploadCommand) run(ctx context.Context) error {
	format, err := conversion.ParseImportFormat(c.importFormat)
	if err != nil {
		return fmt.Errorf("parse import format: %w", err)
	}

	resp, err := c.importSvc.UploadResultsFile(
		ctx,
		connect.NewRequest(&importv1.UploadResultsFileRequest{
			RaceGridId:   c.raceGridID,
			ImportFormat: format,
			Payload:      c.payload,
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
