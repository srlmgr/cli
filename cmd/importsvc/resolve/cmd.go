package resolve

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
	"github.com/srlmgr/cli/cmd/support"
	"github.com/srlmgr/cli/log"
)

//nolint:lll // readability
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve driver/car mappings for an import batch",
		Long:  "Resolve driver and car mappings via backend.import.v1.ImportService.ResolveMappings",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")
			r := support.NewResolve(cmd.Context())
			var resp *queryv1.ResolveSeasonTargetResponse
			var err error
			resp, err = r.AutoResolve(cmd.Flags())
			if err != nil {
				return err
			}
			runner := &resolveCommand{
				raceGridID: resp.RaceGridId,
				out:        cmd.OutOrStdout(),
				importSvc: importclient.NewImportServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32("season-id", 0, "ID of the season to resolve mappings for")
	cmd.Flags().Uint32("grid-id", 0, "ID of the race grid to resolve mappings for")
	cmd.Flags().Uint32("event-no", 0, "Number of the event to resolve mappings for")
	cmd.Flags().Uint32("race-no", 0, "Number of the race to resolve mappings for")
	cmd.Flags().Uint32("grid-no", 0, "Number of the race grid to resolve mappings for")

	return cmd
}

type importClient interface {
	ResolveMappings(
		context.Context,
		*connect.Request[importv1.ResolveMappingsRequest],
	) (*connect.Response[importv1.ResolveMappingsResponse], error)
}

type resolveCommand struct {
	raceGridID uint32
	out        io.Writer
	importSvc  importClient
}

func (c *resolveCommand) run(ctx context.Context) error {
	resp, err := c.importSvc.ResolveMappings(
		ctx,
		connect.NewRequest(&importv1.ResolveMappingsRequest{
			RaceGridId: c.raceGridID,
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
