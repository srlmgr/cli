package deletepenalty

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

type importClient interface {
	DeletePenalty(
		context.Context,
		*connect.Request[importv1.DeletePenaltyRequest],
	) (*connect.Response[importv1.DeletePenaltyResponse], error)
}

type deletePenaltyCommand struct {
	penaltyID uint32
	out       io.Writer
	importSvc importClient
}

func NewCmd() *cobra.Command {
	var penaltyID uint32

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a penalty",
		Long:  "Delete a penalty via backend.import.v1.ImportService.DeletePenalty",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.GetFromContext(cmd.Context()).Named("rpc")

			runner := &deletePenaltyCommand{
				penaltyID: penaltyID,
				out:       cmd.OutOrStdout(),
				importSvc: importclient.NewImportServiceClient(
					config.APIAddr, config.APIToken, logger,
				),
			}
			return runner.run(cmd.Context())
		},
	}

	cmd.Flags().Uint32Var(&penaltyID, "penalty-id", 0, "ID of the penalty to delete")
	if err := cmd.MarkFlagRequired("penalty-id"); err != nil {
		panic(fmt.Sprintf("failed to mark 'penalty-id' flag as required: %v", err))
	}

	return cmd
}

func (c *deletePenaltyCommand) run(ctx context.Context) error {
	req := &importv1.DeletePenaltyRequest{
		PenaltyId: c.penaltyID,
	}

	_, err := c.importSvc.DeletePenalty(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("delete penalty: %w", err)
	}

	fmt.Fprintf(c.out, "Successfully deleted penalty with ID: %d\n", c.penaltyID)
	return nil
}
