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

func (r *setupRunner) setupSimulations(
	ctx context.Context,
	sims []SimulationConfig,
	psIDs map[string]uint32,
) error {
	for _, sim := range sims {
		simID, created, err := r.ensureSimulation(ctx, sim.Name)
		if err != nil {
			return fmt.Errorf("simulation %q: %w", sim.Name, err)
		}

		if err := r.printResult("simulation", sim.Name, simID, created); err != nil {
			return err
		}

		if err := r.setupSeriesList(ctx, simID, sim.Series, psIDs); err != nil {
			return fmt.Errorf("simulation %q series: %w", sim.Name, err)
		}
	}

	return nil
}

func (r *setupRunner) setupSeriesList(
	ctx context.Context,
	simID uint32,
	series []SeriesConfig,
	psIDs map[string]uint32,
) error {
	for _, sr := range series {
		srID, created, err := r.ensureSeries(ctx, simID, sr.Name)
		if err != nil {
			return fmt.Errorf("series %q: %w", sr.Name, err)
		}

		if err := r.printResult("series", sr.Name, srID, created); err != nil {
			return err
		}

		if err := r.setupSeasonList(ctx, srID, sr.Seasons, psIDs); err != nil {
			return fmt.Errorf("series %q seasons: %w", sr.Name, err)
		}
	}

	return nil
}

func (r *setupRunner) setupSeasonList(
	ctx context.Context,
	seriesID uint32,
	seasons []SeasonConfig,
	psIDs map[string]uint32,
) error {
	for _, sn := range seasons {
		psID := psIDs[sn.PointSystem]

		snID, created, err := r.ensureSeason(ctx, seriesID, psID, sn.Name)
		if err != nil {
			return fmt.Errorf("season %q: %w", sn.Name, err)
		}

		if err := r.printResult("season", sn.Name, snID, created); err != nil {
			return err
		}
	}

	return nil
}

func (r *setupRunner) setupCarManufacturers(
	ctx context.Context,
	mfrs []CarManufacturerConfig,
) error {
	for _, mfr := range mfrs {
		mfrID, created, err := r.ensureCarManufacturer(ctx, mfr.Name)
		if err != nil {
			return fmt.Errorf("car manufacturer %q: %w", mfr.Name, err)
		}

		if err := r.printResult("car-manufacturer", mfr.Name, mfrID, created); err != nil {
			return err
		}

		if err := r.setupBrandList(ctx, mfrID, mfr.Brands); err != nil {
			return fmt.Errorf("car manufacturer %q brands: %w", mfr.Name, err)
		}
	}

	return nil
}

func (r *setupRunner) setupBrandList(
	ctx context.Context,
	mfrID uint32,
	brands []BrandConfig,
) error {
	for _, b := range brands {
		brandID, created, err := r.ensureCarBrand(ctx, mfrID, b.Name)
		if err != nil {
			return fmt.Errorf("car brand %q: %w", b.Name, err)
		}

		if err := r.printResult("car-brand", b.Name, brandID, created); err != nil {
			return err
		}

		if err := r.setupModelList(ctx, mfrID, brandID, b.Models); err != nil {
			return fmt.Errorf("car brand %q models: %w", b.Name, err)
		}
	}

	return nil
}

func (r *setupRunner) setupModelList(
	ctx context.Context,
	mfrID, brandID uint32,
	models []ModelConfig,
) error {
	for _, m := range models {
		modelID, created, err := r.ensureCarModel(ctx, mfrID, brandID, m.Name)
		if err != nil {
			return fmt.Errorf("car model %q: %w", m.Name, err)
		}

		if err := r.printResult("car-model", m.Name, modelID, created); err != nil {
			return err
		}
	}

	return nil
}

func (r *setupRunner) setupTracks(
	ctx context.Context,
	tracks []TrackConfig,
) error {
	for _, t := range tracks {
		trackID, created, err := r.ensureTrack(ctx, t.Name)
		if err != nil {
			return fmt.Errorf("track %q: %w", t.Name, err)
		}

		if err := r.printResult("track", t.Name, trackID, created); err != nil {
			return err
		}

		if err := r.setupLayoutList(ctx, trackID, t.Layouts); err != nil {
			return fmt.Errorf("track %q layouts: %w", t.Name, err)
		}
	}

	return nil
}

func (r *setupRunner) setupLayoutList(
	ctx context.Context,
	trackID uint32,
	layouts []LayoutConfig,
) error {
	for _, l := range layouts {
		layID, created, err := r.ensureTrackLayout(ctx, trackID, l.Name)
		if err != nil {
			return fmt.Errorf("track layout %q: %w", l.Name, err)
		}

		if err := r.printResult("track-layout", l.Name, layID, created); err != nil {
			return err
		}
	}

	return nil
}

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
