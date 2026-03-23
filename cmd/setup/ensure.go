//nolint:dupl // methods are similar by design
package setup

import (
	"context"
	"fmt"
	"sort"
	"time"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensurePointSystem(
	ctx context.Context, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.PointSystem, error) {
			resp, err := r.qrySvc.ListPointSystems(ctx,
				connect.NewRequest(&queryv1.ListPointSystemsRequest{}),
			)
			if err != nil {
				return nil, fmt.Errorf("list point systems: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(p *commonv1.PointSystem) bool { return p.GetName() == name },
		func(p *commonv1.PointSystem) uint32 { return p.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreatePointSystem(ctx,
				connect.NewRequest(&commandv1.CreatePointSystemRequest{
					Name: name,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create point system: %w", err)
			}

			return resp.Msg.GetPointSystem().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureSimulation(
	ctx context.Context, name string, isActive bool,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.Simulation, error) {
			resp, err := r.qrySvc.ListSimulations(ctx,
				connect.NewRequest(&queryv1.ListSimulationsRequest{}),
			)
			if err != nil {
				return nil, fmt.Errorf("list simulations: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(s *commonv1.Simulation) bool { return s.GetName() == name },
		func(s *commonv1.Simulation) uint32 { return s.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateSimulation(ctx,
				connect.NewRequest(&commandv1.CreateSimulationRequest{
					Name:     name,
					IsActive: isActive,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create simulation: %w", err)
			}

			return resp.Msg.GetSimulation().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureSeries(
	ctx context.Context, simID uint32, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.Series, error) {
			resp, err := r.qrySvc.ListSeries(ctx,
				connect.NewRequest(&queryv1.ListSeriesRequest{
					SimulationId: simID,
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("list series: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(s *commonv1.Series) bool { return s.GetName() == name },
		func(s *commonv1.Series) uint32 { return s.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateSeries(ctx,
				connect.NewRequest(&commandv1.CreateSeriesRequest{
					SimulationId: simID,
					Name:         name,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create series: %w", err)
			}

			return resp.Msg.GetSeries().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureSeason(
	ctx context.Context, seriesID, psID uint32, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.Season, error) {
			resp, err := r.qrySvc.ListSeasons(ctx,
				connect.NewRequest(&queryv1.ListSeasonsRequest{
					SeriesId: seriesID,
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("list seasons: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(s *commonv1.Season) bool { return s.GetName() == name },
		func(s *commonv1.Season) uint32 { return s.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateSeason(ctx,
				connect.NewRequest(&commandv1.CreateSeasonRequest{
					SeriesId:      seriesID,
					Name:          name,
					PointSystemId: psID,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create season: %w", err)
			}

			return resp.Msg.GetSeason().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureCarManufacturer(
	ctx context.Context, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.CarManufacturer, error) {
			resp, err := r.qrySvc.ListCarManufacturers(ctx,
				connect.NewRequest(&queryv1.ListCarManufacturersRequest{}),
			)
			if err != nil {
				return nil, fmt.Errorf("list car manufacturers: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(m *commonv1.CarManufacturer) bool { return m.GetName() == name },
		func(m *commonv1.CarManufacturer) uint32 { return m.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateCarManufacturer(ctx,
				connect.NewRequest(&commandv1.CreateCarManufacturerRequest{
					Name: name,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create car manufacturer: %w", err)
			}

			return resp.Msg.GetCarManufacturer().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureCarBrand(
	ctx context.Context, mfrID uint32, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.CarBrand, error) {
			resp, err := r.qrySvc.ListCarBrands(ctx,
				connect.NewRequest(&queryv1.ListCarBrandsRequest{
					ManufacturerId: mfrID,
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("list car brands: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(b *commonv1.CarBrand) bool { return b.GetName() == name },
		func(b *commonv1.CarBrand) uint32 { return b.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateCarBrand(ctx,
				connect.NewRequest(&commandv1.CreateCarBrandRequest{
					ManufacturerId: mfrID,
					Name:           name,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create car brand: %w", err)
			}

			return resp.Msg.GetCarBrand().GetId(), nil
		},
		r.dryRun,
	)
}

// ensureCarModel finds or creates a car model, scoped to both manufacturer and brand.
// ListCarModels filters by manufacturer; brand-level scoping is applied in code.
//
//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureCarModel(
	ctx context.Context, mfrID, brandID uint32, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.CarModel, error) {
			resp, err := r.qrySvc.ListCarModels(ctx,
				connect.NewRequest(&queryv1.ListCarModelsRequest{
					ManufacturerId: mfrID,
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("list car models: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(m *commonv1.CarModel) bool {
			return m.GetBrandId() == brandID && m.GetName() == name
		},
		func(m *commonv1.CarModel) uint32 { return m.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateCarModel(ctx,
				connect.NewRequest(&commandv1.CreateCarModelRequest{
					BrandId: brandID,
					Name:    name,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create car model: %w", err)
			}

			return resp.Msg.GetCarModel().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureTrack(
	ctx context.Context, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.Track, error) {
			resp, err := r.qrySvc.ListTracks(ctx,
				connect.NewRequest(&queryv1.ListTracksRequest{}),
			)
			if err != nil {
				return nil, fmt.Errorf("list tracks: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(t *commonv1.Track) bool { return t.GetName() == name },
		func(t *commonv1.Track) uint32 { return t.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateTrack(ctx,
				connect.NewRequest(&commandv1.CreateTrackRequest{
					Name: name,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create track: %w", err)
			}

			return resp.Msg.GetTrack().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureTrackLayout(
	ctx context.Context, trackID uint32, name string,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.TrackLayout, error) {
			resp, err := r.qrySvc.ListTrackLayouts(ctx,
				connect.NewRequest(&queryv1.ListTrackLayoutsRequest{
					TrackId: trackID,
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("list track layouts: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(l *commonv1.TrackLayout) bool { return l.GetName() == name },
		func(l *commonv1.TrackLayout) uint32 { return l.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateTrackLayout(ctx,
				connect.NewRequest(&commandv1.CreateTrackLayoutRequest{
					TrackId: trackID,
					Name:    name,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create track layout: %w", err)
			}

			return resp.Msg.GetTrackLayout().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureDriver(
	ctx context.Context, cfg DriverConfig,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.Driver, error) {
			resp, err := r.qrySvc.ListDrivers(ctx,
				connect.NewRequest(&queryv1.ListDriversRequest{}),
			)
			if err != nil {
				return nil, fmt.Errorf("list drivers: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(d *commonv1.Driver) bool { return d.GetName() == cfg.Name },
		func(d *commonv1.Driver) uint32 { return d.GetId() },
		func(ctx context.Context) (uint32, error) {
			resp, err := r.cmdSvc.CreateDriver(ctx,
				connect.NewRequest(&commandv1.CreateDriverRequest{
					Name:       cfg.Name,
					ExternalId: cfg.ExternalID,
					IsActive:   cfg.IsActive,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create driver: %w", err)
			}

			return resp.Msg.GetDriver().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace,funlen // editor/linter issue,much to do
func (r *setupRunner) ensureEvent(
	ctx context.Context, seasonID, trackLayoutID uint32, cfg *EventConfig,
) (uint32, bool, error) {
	if cfg == nil {
		return 0, false, fmt.Errorf("event config is nil")
	}

	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.Event, error) {
			resp, err := r.qrySvc.ListEvents(ctx,
				connect.NewRequest(&queryv1.ListEventsRequest{
					SeasonId: seasonID,
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("list events: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(e *commonv1.Event) bool { return e.GetName() == cfg.Name },
		func(e *commonv1.Event) uint32 { return e.GetId() },
		func(ctx context.Context) (uint32, error) {
			status, err := parseEventStatus(cfg.Status.String())
			if err != nil {
				return 0, err
			}

			processingState, err := parseEventProcessingState(cfg.ProcessingState.String())
			if err != nil {
				return 0, err
			}

			req := &commandv1.CreateEventRequest{
				SeasonId:        seasonID,
				TrackLayoutId:   trackLayoutID,
				Name:            cfg.Name,
				Status:          status,
				ProcessingState: processingState,
			}

			if cfg.Date != "" {
				eventDate, parseErr := time.Parse(time.DateOnly, cfg.Date)
				if parseErr != nil {
					return 0, fmt.Errorf("parse event date %q: %w", cfg.Date, parseErr)
				}

				req.SetEventDate(timestamppb.New(eventDate))
			}

			resp, err := r.cmdSvc.CreateEvent(ctx, connect.NewRequest(req))
			if err != nil {
				return 0, fmt.Errorf("create event: %w", err)
			}

			return resp.Msg.GetEvent().GetId(), nil
		},
		r.dryRun,
	)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) ensureRace(
	ctx context.Context, eventID uint32, cfg RaceConfig,
) (uint32, bool, error) {
	return findOrCreate(ctx,
		func(ctx context.Context) ([]*commonv1.Race, error) {
			resp, err := r.qrySvc.ListRaces(ctx,
				connect.NewRequest(&queryv1.ListRacesRequest{
					EventId: eventID,
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("list races: %w", err)
			}

			return resp.Msg.GetItems(), nil
		},
		func(rc *commonv1.Race) bool { return rc.GetName() == cfg.Name },
		func(rc *commonv1.Race) uint32 { return rc.GetId() },
		func(ctx context.Context) (uint32, error) {
			enumVal, ok := commonv1.RaceSessionType_value[cfg.SessionType]
			if !ok {
				return 0, fmt.Errorf(
					"unknown sessionType %q; valid values: %v",
					cfg.SessionType, validRaceSessionTypes(),
				)
			}

			resp, err := r.cmdSvc.CreateRace(ctx,
				connect.NewRequest(&commandv1.CreateRaceRequest{
					EventId:     eventID,
					Name:        cfg.Name,
					SessionType: commonv1.RaceSessionType(enumVal),
					SequenceNo:  cfg.SequenceNo,
				}),
			)
			if err != nil {
				return 0, fmt.Errorf("create race: %w", err)
			}

			return resp.Msg.GetRace().GetId(), nil
		},
		r.dryRun,
	)
}

func parseEventStatus(raw string) (commonv1.EventStatus, error) {
	enumVal, ok := commonv1.EventStatus_value[raw]
	if !ok {
		return commonv1.EventStatus_EVENT_STATUS_UNSPECIFIED, fmt.Errorf(
			"unknown status %q; valid values: %v",
			raw,
			validEventStatuses(),
		)
	}

	return commonv1.EventStatus(enumVal), nil
}

func parseEventProcessingState(raw string) (commonv1.EventProcessingState, error) {
	enumVal, ok := commonv1.EventProcessingState_value[raw]
	if !ok {
		return commonv1.EventProcessingState_EVENT_PROCESSING_STATE_UNSPECIFIED, fmt.Errorf(
			"unknown processingState %q; valid values: %v",
			raw,
			validEventProcessingStates(),
		)
	}

	return commonv1.EventProcessingState(enumVal), nil
}

func validEventStatuses() []string {
	names := make([]string, 0, len(commonv1.EventStatus_value))
	for name := range commonv1.EventStatus_value {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func validEventProcessingStates() []string {
	names := make([]string, 0, len(commonv1.EventProcessingState_value))
	for name := range commonv1.EventProcessingState_value {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// validRaceSessionTypes returns the list of valid session type enum names.
func validRaceSessionTypes() []string {
	names := make([]string, 0, len(commonv1.RaceSessionType_value))
	for name := range commonv1.RaceSessionType_value {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}
