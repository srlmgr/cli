package support

import (
	"context"
	"fmt"

	"buf.build/gen/go/srlmgr/api/connectrpc/go/backend/query/v1/queryv1connect"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"github.com/spf13/pflag"

	"github.com/srlmgr/cli/cmd/config"
	"github.com/srlmgr/cli/cmd/query/client"
	"github.com/srlmgr/cli/log"
)

func NewResolve(ctx context.Context) *Resolve {
	return &Resolve{
		ctx: ctx,
		l:   log.GetFromContext(ctx).Named("rpc"),
		querySvc: client.NewQueryServiceClient(
			config.APIAddr,
			log.GetFromContext(ctx).Named("rpc"),
		),
	}
}

type Resolve struct {
	ctx      context.Context
	l        *log.Logger
	querySvc queryv1connect.QueryServiceClient
}

// AutoResolve automatically resolves the season target based on the provided flags.
// It prioritizes grid, race and finally event
//
//nolint:whitespace // editor/linter issue
func (r *Resolve) AutoResolve(
	flags *pflag.FlagSet,
) (*queryv1.ResolveSeasonTargetResponse, error) {
	if flags.Changed("grid-id") || flags.Changed("grid-no") {
		return r.ResolveRaceGridID(flags)
	}
	if flags.Changed("race-id") || flags.Changed("race-no") {
		return r.ResolveRaceID(flags)
	}
	return r.ResolveEventID(flags)
}

//nolint:whitespace // editor/linter issue
func (r *Resolve) ResolveEventID(
	flags *pflag.FlagSet,
) (*queryv1.ResolveSeasonTargetResponse, error) {
	eventIDSet := flags.Changed("event-id")
	if eventIDSet {
		if eventID, err := flags.GetUint32("event-id"); err != nil {
			return nil, fmt.Errorf("get event-id flag: %w", err)
		} else {
			return &queryv1.ResolveSeasonTargetResponse{
				EventId: eventID,
			}, nil
		}
	}
	seasonIDSet := flags.Changed("season-id")
	eventNoSet := flags.Changed("event-no")
	//nolint:lll // readability
	if !seasonIDSet || !eventNoSet {
		return nil, fmt.Errorf("season-id and event-no flags are required if event-id is not set")
	}
	seasonID, err := flags.GetUint32("season-id")
	if err != nil {
		return nil, fmt.Errorf("get season-id flag: %w", err)
	}
	eventNo, err := flags.GetUint32("event-no")
	if err != nil {
		return nil, fmt.Errorf("get event-no flag: %w", err)
	}
	resp, err := r.querySvc.ResolveSeasonTarget(
		r.ctx,
		connect.NewRequest(&queryv1.ResolveSeasonTargetRequest{
			SeasonId: seasonID,
			EventNo:  eventNo,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("resolve season target: %w", err)
	}
	return resp.Msg, nil
}

//nolint:whitespace // editor/linter issue
func (r *Resolve) ResolveRaceID(
	flags *pflag.FlagSet,
) (*queryv1.ResolveSeasonTargetResponse, error) {
	raceIDSet := flags.Changed("race-id")
	if raceIDSet {
		if raceID, err := flags.GetUint32("race-id"); err != nil {
			return nil, fmt.Errorf("get race-id flag: %w", err)
		} else {
			return &queryv1.ResolveSeasonTargetResponse{
				RaceId: raceID,
			}, nil
		}
	}
	seasonIDSet := flags.Changed("season-id")
	eventNoSet := flags.Changed("event-no")
	raceNoSet := flags.Changed("race-no")

	if !seasonIDSet || !eventNoSet || !raceNoSet {
		return nil, fmt.Errorf(
			"season-id, event-no and race-no flags are required if race-id is not set",
		)
	}
	seasonID, err := flags.GetUint32("season-id")
	if err != nil {
		return nil, fmt.Errorf("get season-id flag: %w", err)
	}
	eventNo, err := flags.GetUint32("event-no")
	if err != nil {
		return nil, fmt.Errorf("get event-no flag: %w", err)
	}
	raceNo, err := flags.GetUint32("race-no")
	if err != nil {
		return nil, fmt.Errorf("get race-no flag: %w", err)
	}
	resp, err := r.querySvc.ResolveSeasonTarget(
		r.ctx,
		connect.NewRequest(&queryv1.ResolveSeasonTargetRequest{
			SeasonId: seasonID,
			EventNo:  eventNo,
			RaceNo:   raceNo,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("resolve season target: %w", err)
	}
	return resp.Msg, nil
}

//nolint:whitespace,funlen // editor/linter issue
func (r *Resolve) ResolveRaceGridID(
	flags *pflag.FlagSet,
) (*queryv1.ResolveSeasonTargetResponse, error) {
	raceIDSet := flags.Changed("grid-id")
	if raceIDSet {
		if gridID, err := flags.GetUint32("grid-id"); err != nil {
			return nil, fmt.Errorf("get grid-id flag: %w", err)
		} else {
			return &queryv1.ResolveSeasonTargetResponse{
				RaceGridId: gridID,
			}, nil
		}
	}
	seasonIDSet := flags.Changed("season-id")
	eventNoSet := flags.Changed("event-no")
	raceNoSet := flags.Changed("race-no")
	gridNoSet := flags.Changed("grid-no")

	if !seasonIDSet || !eventNoSet || !raceNoSet || !gridNoSet {
		return nil, fmt.Errorf(
			"season-id, event-no, race-no and grid-no flags are required if grid-id is not set",
		)
	}
	seasonID, err := flags.GetUint32("season-id")
	if err != nil {
		return nil, fmt.Errorf("get season-id flag: %w", err)
	}
	eventNo, err := flags.GetUint32("event-no")
	if err != nil {
		return nil, fmt.Errorf("get event-no flag: %w", err)
	}
	raceNo, err := flags.GetUint32("race-no")
	if err != nil {
		return nil, fmt.Errorf("get race-no flag: %w", err)
	}
	gridNo, err := flags.GetUint32("grid-no")
	if err != nil {
		return nil, fmt.Errorf("get grid-no flag: %w", err)
	}
	resp, err := r.querySvc.ResolveSeasonTarget(
		r.ctx,
		connect.NewRequest(&queryv1.ResolveSeasonTargetRequest{
			SeasonId:   seasonID,
			EventNo:    eventNo,
			RaceNo:     raceNo,
			RaceGridNo: gridNo,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("resolve season target: %w", err)
	}
	return resp.Msg, nil
}
