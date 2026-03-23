//nolint:dupl // may happen
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
//
//nolint:lll // readability
func (r *setupRunner) run(ctx context.Context) error {
	cfg, err := loadConfig(r.filePath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	driverIDs, err := r.setupDrivers(ctx, cfg.Drivers)
	if err != nil {
		return fmt.Errorf("setup drivers: %w", err)
	}

	psIDs, err := r.setupPointSystems(ctx, cfg.PointSystems)
	if err != nil {
		return fmt.Errorf("setup point systems: %w", err)
	}

	layoutIDs, err := r.setupTracks(ctx, cfg.Tracks)
	if err != nil {
		return fmt.Errorf("setup tracks: %w", err)
	}

	simIDs, err := r.setupSimulations(ctx, cfg.Simulations, psIDs, layoutIDs)
	if err != nil {
		return fmt.Errorf("setup simulations: %w", err)
	}

	modelIDs, err := r.setupCarManufacturers(ctx, cfg.CarManufacturers)
	if err != nil {
		return fmt.Errorf("setup car manufacturers: %w", err)
	}

	if err := r.setupAliases(ctx, cfg, simIDs, driverIDs, modelIDs, layoutIDs); err != nil {
		return fmt.Errorf("setup aliases: %w", err)
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupDrivers(
	ctx context.Context,
	drivers []DriverConfig,
) (map[string]uint32, error) {
	ids := make(map[string]uint32, len(drivers))

	for i := range drivers {
		d := drivers[i]
		id, created, err := r.ensureDriver(ctx, d)
		if err != nil {
			return nil, fmt.Errorf("driver %q: %w", d.Name, err)
		}

		if err := r.printResult("driver", d.Name, id, created); err != nil {
			return nil, err
		}

		ids[d.Name] = id
	}

	return ids, nil
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
	layoutIDs map[string]uint32,
) (map[string]uint32, error) {
	simIDs := make(map[string]uint32, len(sims))

	for i := range sims {
		simID, created, err := r.ensureSimulation(ctx, sims[i].Name, sims[i].IsActive)
		if err != nil {
			return nil, fmt.Errorf("simulation %q: %w", sims[i].Name, err)
		}

		if err := r.printResult("simulation", sims[i].Name, simID, created); err != nil {
			return nil, err
		}

		simIDs[sims[i].Name] = simID

		if err := r.setupSeriesList(
			ctx,
			simID,
			sims[i].Series,
			psIDs,
			layoutIDs); err != nil {
			return nil, fmt.Errorf("simulation %q series: %w", sims[i].Name, err)
		}
	}

	return simIDs, nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupSeriesList(
	ctx context.Context,
	simID uint32,
	series []SeriesConfig,
	psIDs map[string]uint32,
	layoutIDs map[string]uint32,
) error {
	for i := range series {
		srID, created, err := r.ensureSeries(ctx, simID, series[i].Name)
		if err != nil {
			return fmt.Errorf("series %q: %w", series[i].Name, err)
		}

		if err := r.printResult("series", series[i].Name, srID, created); err != nil {
			return err
		}

		if err := r.setupSeasonList(
			ctx,
			srID,
			series[i].Seasons,
			psIDs,
			layoutIDs,
		); err != nil {
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
	layoutIDs map[string]uint32,
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

		if err := r.setupEventList(ctx, snID, seasons[i].Events, layoutIDs); err != nil {
			return fmt.Errorf("season %q events: %w", seasons[i].Name, err)
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupEventList(
	ctx context.Context,
	seasonID uint32,
	events []EventConfig,
	layoutIDs map[string]uint32,
) error {
	for i := range events {
		if events[i].TrackLayout != "" {
			if _, ok := layoutIDs[events[i].TrackLayout]; !ok {
				return fmt.Errorf(
					"event %q: track layout %q not found; ensure it is defined under tracks",
					events[i].Name, events[i].TrackLayout,
				)
			}
		}

		layoutID := layoutIDs[events[i].TrackLayout]

		evID, created, err := r.ensureEvent(ctx, seasonID, layoutID, &events[i])
		if err != nil {
			return fmt.Errorf("event %q: %w", events[i].Name, err)
		}

		if err := r.printResult("event", events[i].Name, evID, created); err != nil {
			return err
		}

		if err := r.setupRaceList(ctx, evID, events[i].Races); err != nil {
			return fmt.Errorf("event %q races: %w", events[i].Name, err)
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupRaceList(
	ctx context.Context,
	eventID uint32,
	races []RaceConfig,
) error {
	for i := range races {
		id, created, err := r.ensureRace(ctx, eventID, races[i])
		if err != nil {
			return fmt.Errorf("race %q: %w", races[i].Name, err)
		}

		if err := r.printResult("race", races[i].Name, id, created); err != nil {
			return err
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupCarManufacturers(
	ctx context.Context,
	mfrs []CarManufacturerConfig,
) (map[string]uint32, error) {
	modelIDs := make(map[string]uint32)

	for i := range mfrs {
		mfrID, created, err := r.ensureCarManufacturer(ctx, mfrs[i].Name)
		if err != nil {
			return nil, fmt.Errorf("car manufacturer %q: %w", mfrs[i].Name, err)
		}
		//nolint:lll // readability
		if err := r.printResult("car-manufacturer", mfrs[i].Name, mfrID, created); err != nil {
			return nil, err
		}

		if err := r.setupBrandList(ctx, mfrID, mfrs[i].Brands, modelIDs); err != nil {
			return nil, fmt.Errorf("car manufacturer %q brands: %w", mfrs[i].Name, err)
		}
	}

	return modelIDs, nil
}

//nolint:whitespace,lll // editor/linter issue
func (r *setupRunner) setupBrandList(
	ctx context.Context,
	mfrID uint32,
	brands []BrandConfig,
	modelIDs map[string]uint32,
) error {
	for i := range brands {
		brandID, created, err := r.ensureCarBrand(ctx, mfrID, brands[i].Name)
		if err != nil {
			return fmt.Errorf("car brand %q: %w", brands[i].Name, err)
		}

		if err := r.printResult("car-brand", brands[i].Name, brandID, created); err != nil {
			return err
		}

		if err := r.setupModelList(ctx, mfrID, brandID, brands[i].Models, modelIDs); err != nil {
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
	modelIDs map[string]uint32,
) error {
	for i := range models {
		modelID, created, err := r.ensureCarModel(ctx, mfrID, brandID, models[i].Name)
		if err != nil {
			return fmt.Errorf("car model %q: %w", models[i].Name, err)
		}

		if err := r.printResult("car-model", models[i].Name, modelID, created); err != nil {
			return err
		}

		modelIDs[models[i].Name] = modelID
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupTracks(
	ctx context.Context,
	tracks []TrackConfig,
) (map[string]uint32, error) {
	layoutIDs := make(map[string]uint32)

	for i := range tracks {
		trackID, created, err := r.ensureTrack(ctx, tracks[i].Name)
		if err != nil {
			return nil, fmt.Errorf("track %q: %w", tracks[i].Name, err)
		}

		if err := r.printResult("track", tracks[i].Name, trackID, created); err != nil {
			return nil, err
		}

		if err := r.setupLayoutList(ctx, trackID, tracks[i].Layouts, layoutIDs); err != nil {
			return nil, fmt.Errorf("track %q layouts: %w", tracks[i].Name, err)
		}
	}

	return layoutIDs, nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupLayoutList(
	ctx context.Context,
	trackID uint32,
	layouts []LayoutConfig,
	layoutIDs map[string]uint32,
) error {
	for i := range layouts {
		layID, created, err := r.ensureTrackLayout(ctx, trackID, layouts[i].Name)
		if err != nil {
			return fmt.Errorf("track layout %q: %w", layouts[i].Name, err)
		}

		if err := r.printResult("track-layout", layouts[i].Name, layID, created); err != nil {
			return err
		}

		layoutIDs[layouts[i].Name] = layID
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupAliases(
	ctx context.Context,
	cfg *SetupConfig,
	simIDs, driverIDs, modelIDs, layoutIDs map[string]uint32,
) error {
	if err := r.setupDriverAliases(ctx, cfg.Drivers, simIDs, driverIDs); err != nil {
		return err
	}

	if err := r.setupTrackLayoutAliases(ctx, cfg.Tracks, simIDs, layoutIDs); err != nil {
		return err
	}

	return r.setupCarModelAliases(ctx, cfg.CarManufacturers, simIDs, modelIDs)
}

//nolint:whitespace // editor/linter issue
func (r *setupRunner) setupDriverAliases(
	ctx context.Context,
	drivers []DriverConfig,
	simIDs, driverIDs map[string]uint32,
) error {
	for i := range drivers {
		d := &drivers[i]
		driverID := driverIDs[d.Name]

		for j := range d.Simulations {
			simCfg := &d.Simulations[j]
			if len(simCfg.Aliases) == 0 {
				continue
			}

			simID, ok := simIDs[simCfg.Name]
			if !ok {
				return fmt.Errorf(
					"driver %q: simulation %q not found; ensure it is defined under simulations",
					d.Name, simCfg.Name,
				)
			}

			if err := r.setSimulationDriverAliases(
				ctx,
				driverID,
				simID,
				simCfg.Aliases,
			); err != nil {
				return fmt.Errorf("driver %q simulation %q aliases: %w", d.Name, simCfg.Name, err)
			}

			if err := r.printAliasResult("driver", d.Name, simCfg.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

//nolint:whitespace,lll // editor/linter issue
func (r *setupRunner) setupTrackLayoutAliases(
	ctx context.Context,
	tracks []TrackConfig,
	simIDs, layoutIDs map[string]uint32,
) error {
	for i := range tracks {
		for j := range tracks[i].Layouts {
			layout := &tracks[i].Layouts[j]
			layoutID := layoutIDs[layout.Name]

			for k := range layout.Simulations {
				simCfg := &layout.Simulations[k]
				if len(simCfg.Aliases) == 0 {
					continue
				}

				simID, ok := simIDs[simCfg.Name]
				if !ok {
					return fmt.Errorf(
						"track layout %q: simulation %q not found; ensure it is defined under simulations",
						layout.Name,
						simCfg.Name,
					)
				}

				if err := r.setSimulationTrackLayoutAliases(
					ctx,
					layoutID,
					simID,
					simCfg.Aliases,
				); err != nil {
					return fmt.Errorf(
						"track layout %q simulation %q aliases: %w",
						layout.Name, simCfg.Name, err,
					)
				}

				if err := r.printAliasResult("track-layout", layout.Name, simCfg.Name); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//nolint:whitespace,funlen // editor/linter issue
func (r *setupRunner) setupCarModelAliases(
	ctx context.Context,
	mfrs []CarManufacturerConfig,
	simIDs, modelIDs map[string]uint32,
) error {
	for i := range mfrs {
		for j := range mfrs[i].Brands {
			for k := range mfrs[i].Brands[j].Models {
				model := &mfrs[i].Brands[j].Models[k]
				modelID := modelIDs[model.Name]

				for l := range model.Simulations {
					simCfg := &model.Simulations[l]
					if len(simCfg.Aliases) == 0 {
						continue
					}

					simID, ok := simIDs[simCfg.Name]
					if !ok {
						return fmt.Errorf(
							"car model %q: simulation %q not found; ensure it is defined under simulations",
							model.Name,
							simCfg.Name,
						)
					}

					if err := r.setSimulationCarAliases(
						ctx,
						modelID,
						simID,
						simCfg.Aliases,
					); err != nil {
						return fmt.Errorf(
							"car model %q simulation %q aliases: %w",
							model.Name, simCfg.Name, err,
						)
					}

					if err := r.printAliasResult("car-model", model.Name, simCfg.Name); err != nil {
						return err
					}
				}
			}
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

func (r *setupRunner) printAliasResult(entityType, name, simName string) error {
	if r.dryRun {
		_, err := fmt.Fprintf(
			r.out, "dry-run: would set %s %q aliases for simulation %q\n",
			entityType, name, simName,
		)

		return err
	}

	_, err := fmt.Fprintf(
		r.out, "set %s %q aliases for simulation %q\n",
		entityType, name, simName,
	)

	return err
}
