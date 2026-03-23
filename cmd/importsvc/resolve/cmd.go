package resolve

import (
	"context"
	"fmt"
	"io"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/srlmgr/cli/cmd/config"
	importclient "github.com/srlmgr/cli/cmd/importsvc/client"
	"github.com/srlmgr/cli/log"
)

func NewCmd() *cobra.Command {
	var raceID uint32
	//nolint:lll // readability
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve driver/car mappings for an import batch",
		Long:  "Resolve driver and car mappings via backend.import.v1.ImportService.ResolveMappings",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			runner := &resolveCommand{
				raceID: raceID,
				out:    cmd.OutOrStdout(),
				importSvc: importclient.NewImportServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32Var(&raceID, "race-id", 0, "ID of the race to resolve mappings for")
	if err := cmd.MarkFlagRequired("race-id"); err != nil {
		panic(fmt.Sprintf("failed to mark 'race-id' flag as required: %v", err))
	}

	return cmd
}

type importClient interface {
	ResolveMappings(
		context.Context,
		*connect.Request[importv1.ResolveMappingsRequest],
	) (*connect.Response[importv1.ResolveMappingsResponse], error)
}

type resolveCommand struct {
	raceID    uint32
	out       io.Writer
	importSvc importClient
}

func (c *resolveCommand) run(ctx context.Context) error {
	resp, err := c.importSvc.ResolveMappings(
		ctx,
		connect.NewRequest(&importv1.ResolveMappingsRequest{
			RaceId: c.raceID,
		}),
	)
	if err != nil {
		return fmt.Errorf("resolve mappings: %w", err)
	}

	if _, err = fmt.Fprintf(
		c.out,
		"Resolved mappings: unresolved_mappings=%d\n",
		resp.Msg.GetUnresolvedMappings(),
	); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
