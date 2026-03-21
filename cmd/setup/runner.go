//nolint:dupl // methods are similar by design
package setup

import (
	"context"
	"fmt"
	"io"
)

type setupRunner struct {
	filePath string
	dryRun   bool
	out      io.Writer
	cmdSvc   commandClient
	qrySvc   queryClient
}

// run loads the config and provisions all entities in the required order.
func (r *setupRunner) run(ctx context.Context) error {
	cfg, err := loadConfig(r.filePath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	psIDs, err := r.setupPointSystems(ctx, cfg.PointSystems)
	if err != nil {
		return fmt.Errorf("setup point systems: %w", err)
	}

	if err := r.setupSimulations(ctx, cfg.Simulations, psIDs); err != nil {
		return fmt.Errorf("setup simulations: %w", err)
	}

	if err := r.setupCarManufacturers(ctx, cfg.CarManufacturers); err != nil {
		return fmt.Errorf("setup car manufacturers: %w", err)
	}

	if err := r.setupTracks(ctx, cfg.Tracks); err != nil {
		return fmt.Errorf("setup tracks: %w", err)
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupPointSystems(
	ctx context.Context,
	items []PointSystemConfig,
) (map[string]uint32, error) {
	ids := make(map[string]uint32, len(items))

	for _, ps := range items {
		id, created, err := r.ensurePointSystem(ctx, ps.Name)
		if err != nil {
			return nil, fmt.Errorf("point system %q: %w", ps.Name, err)
		}

		if err := r.printResult("point-system", ps.Name, id, created); err != nil {
			return nil, err
		}

		ids[ps.Name] = id
	}

	return ids, nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupSimulations(
	ctx context.Context,
	sims []SimulationConfig,
	psIDs map[string]uint32,
) error {
	for i := range sims {
		simID, created, err := r.ensureSimulation(ctx, sims[i].Name)
		if err != nil {
			return fmt.Errorf("simulation %q: %w", sims[i].Name, err)
		}

		if err := r.printResult("simulation", sims[i].Name, simID, created); err != nil {
			return err
		}

		if err := r.setupSeriesList(ctx, simID, sims[i].Series, psIDs); err != nil {
			return fmt.Errorf("simulation %q series: %w", sims[i].Name, err)
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupSeriesList(
	ctx context.Context,
	simID uint32,
	series []SeriesConfig,
	psIDs map[string]uint32,
) error {
	for i := range series {
		srID, created, err := r.ensureSeries(ctx, simID, series[i].Name)
		if err != nil {
			return fmt.Errorf("series %q: %w", series[i].Name, err)
		}

		if err := r.printResult("series", series[i].Name, srID, created); err != nil {
			return err
		}

		if err := r.setupSeasonList(ctx, srID, series[i].Seasons, psIDs); err != nil {
			return fmt.Errorf("series %q seasons: %w", series[i].Name, err)
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupSeasonList(
	ctx context.Context,
	seriesID uint32,
	seasons []SeasonConfig,
	psIDs map[string]uint32,
) error {
	for i := range seasons {
		psID := psIDs[seasons[i].PointSystem]

		snID, created, err := r.ensureSeason(ctx, seriesID, psID, seasons[i].Name)
		if err != nil {
			return fmt.Errorf("season %q: %w", seasons[i].Name, err)
		}

		if err := r.printResult("season", seasons[i].Name, snID, created); err != nil {
			return err
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupCarManufacturers(
	ctx context.Context,
	mfrs []CarManufacturerConfig,
) error {
	for i := range mfrs {
		mfrID, created, err := r.ensureCarManufacturer(ctx, mfrs[i].Name)
		if err != nil {
			return fmt.Errorf("car manufacturer %q: %w", mfrs[i].Name, err)
		}
		//nolint:lll // readability
		if err := r.printResult("car-manufacturer", mfrs[i].Name, mfrID, created); err != nil {
			return err
		}

		if err := r.setupBrandList(ctx, mfrID, mfrs[i].Brands); err != nil {
			return fmt.Errorf("car manufacturer %q brands: %w", mfrs[i].Name, err)
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupBrandList(
	ctx context.Context,
	mfrID uint32,
	brands []BrandConfig,
) error {
	for i := range brands {
		brandID, created, err := r.ensureCarBrand(ctx, mfrID, brands[i].Name)
		if err != nil {
			return fmt.Errorf("car brand %q: %w", brands[i].Name, err)
		}

		if err := r.printResult("car-brand", brands[i].Name, brandID, created); err != nil {
			return err
		}

		if err := r.setupModelList(ctx, mfrID, brandID, brands[i].Models); err != nil {
			return fmt.Errorf("car brand %q models: %w", brands[i].Name, err)
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupModelList(
	ctx context.Context,
	mfrID, brandID uint32,
	models []ModelConfig,
) error {
	for i := range models {
		modelID, created, err := r.ensureCarModel(ctx, mfrID, brandID, models[i].Name)
		if err != nil {
			return fmt.Errorf("car model %q: %w", models[i].Name, err)
		}

		if err := r.printResult("car-model", models[i].Name, modelID, created); err != nil {
			return err
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupTracks(
	ctx context.Context,
	tracks []TrackConfig,
) error {
	for i := range tracks {
		trackID, created, err := r.ensureTrack(ctx, tracks[i].Name)
		if err != nil {
			return fmt.Errorf("track %q: %w", tracks[i].Name, err)
		}

		if err := r.printResult("track", tracks[i].Name, trackID, created); err != nil {
			return err
		}

		if err := r.setupLayoutList(ctx, trackID, tracks[i].Layouts); err != nil {
			return fmt.Errorf("track %q layouts: %w", tracks[i].Name, err)
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupLayoutList(
	ctx context.Context,
	trackID uint32,
	layouts []LayoutConfig,
) error {
	for i := range layouts {
		layID, created, err := r.ensureTrackLayout(ctx, trackID, layouts[i].Name)
		if err != nil {
			return fmt.Errorf("track layout %q: %w", layouts[i].Name, err)
		}

		if err := r.printResult("track-layout", layouts[i].Name, layID, created); err != nil {
			return err
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) printResult(
	entityType, name string, id uint32, created bool,
) error {
	if r.dryRun && created {
		_, err := fmt.Fprintf(
			r.out, "dry-run: would create %s %q\n", entityType, name,
		)

		return err
	}

	status := "existing"
	if created {
		status = "created"
	}

	_, err := fmt.Fprintf(
		r.out, "%s %s %q id=%d\n", status, entityType, name, id,
	)

	return err
}
