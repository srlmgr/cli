//nolint:dupl // methods are similar by design
package setup

import (
	"context"
	"fmt"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
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
	ctx context.Context, name string,
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
					IsActive: true,
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
