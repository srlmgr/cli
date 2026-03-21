package setup

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// SetupConfig is the root configuration for the setup command.
type SetupConfig struct {
	Drivers          []DriverConfig          `yaml:"drivers"`
	PointSystems     []PointSystemConfig     `yaml:"pointSystems"`
	Simulations      []SimulationConfig      `yaml:"simulations"`
	CarManufacturers []CarManufacturerConfig `yaml:"carManufacturers"`
	Tracks           []TrackConfig           `yaml:"tracks"`
}

// PointSystemConfig defines a point system to be created.
type PointSystemConfig struct {
	Name string `yaml:"name"`
}

// DriverConfig defines a driver to be created.
type DriverConfig struct {
	Name       string `yaml:"name"`
	ExternalID string `yaml:"externalId"`
	IsActive   bool   `yaml:"isActive"`
}

// SimulationConfig defines a simulation and its child series.
type SimulationConfig struct {
	Name     string         `yaml:"name"`
	IsActive bool           `yaml:"isActive"`
	Series   []SeriesConfig `yaml:"series"`
}

// SeriesConfig defines a series and its child seasons.
type SeriesConfig struct {
	Name    string         `yaml:"name"`
	Seasons []SeasonConfig `yaml:"seasons"`
}

// SeasonConfig defines a season and its associated point system name.
type SeasonConfig struct {
	Name        string        `yaml:"name"`
	PointSystem string        `yaml:"pointSystem"`
	Events      []EventConfig `yaml:"events"`
}

// EventConfig defines an event under a season.
type EventConfig struct {
	Name            string       `yaml:"name"`
	TrackLayout     string       `yaml:"trackLayout"`
	Date            string       `yaml:"date"`
	Status          string       `yaml:"status"`
	ProcessingState string       `yaml:"processingState"`
	Races           []RaceConfig `yaml:"races"`
}

// RaceConfig defines a race under an event.
type RaceConfig struct {
	Name        string `yaml:"name"`
	SessionType string `yaml:"sessionType"`
	SequenceNo  int32  `yaml:"sequenceNo"`
}

// CarManufacturerConfig defines a car manufacturer and its child brands.
type CarManufacturerConfig struct {
	Name   string        `yaml:"name"`
	Brands []BrandConfig `yaml:"brands"`
}

// BrandConfig defines a car brand and its child models.
type BrandConfig struct {
	Name   string        `yaml:"name"`
	Models []ModelConfig `yaml:"models"`
}

// ModelConfig defines a car model.
type ModelConfig struct {
	Name string `yaml:"name"`
}

// TrackConfig defines a track and its child layouts.
type TrackConfig struct {
	Name    string         `yaml:"name"`
	Layouts []LayoutConfig `yaml:"layouts"`
}

// LayoutConfig defines a track layout.
type LayoutConfig struct {
	Name string `yaml:"name"`
}

// loadConfig reads, parses, and validates the YAML setup file.
func loadConfig(filePath string) (*SetupConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg SetupConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *SetupConfig) validate() error {
	if err := validateDrivers(c.Drivers); err != nil {
		return err
	}

	if err := validatePointSystems(c.PointSystems); err != nil {
		return err
	}

	if err := validateSimulations(c.Simulations); err != nil {
		return err
	}

	if err := validateManufacturers(c.CarManufacturers); err != nil {
		return err
	}

	return validateTracks(c.Tracks)
}

func validatePointSystems(items []PointSystemConfig) error {
	for i, ps := range items {
		if ps.Name == "" {
			return fmt.Errorf("pointSystems[%d]: name is required", i)
		}
	}

	return nil
}

func validateDrivers(items []DriverConfig) error {
	for i := range items {
		d := items[i]
		if d.Name == "" {
			return fmt.Errorf("drivers[%d]: name is required", i)
		}
	}

	return nil
}

func validateSimulations(items []SimulationConfig) error {
	for i := range items {
		if items[i].Name == "" {
			return fmt.Errorf("simulations[%d]: name is required", i)
		}

		if err := validateSeriesList(i, items[i].Series); err != nil {
			return err
		}
	}

	return nil
}

func validateSeriesList(simIdx int, series []SeriesConfig) error {
	for j := range series {
		if series[j].Name == "" {
			return fmt.Errorf(
				"simulations[%d].series[%d]: name is required",
				simIdx, j,
			)
		}

		if err := validateSeasonList(simIdx, j, series[j].Seasons); err != nil {
			return err
		}
	}

	return nil
}

func validateSeasonList(simIdx, serIdx int, seasons []SeasonConfig) error {
	for k := range seasons {
		if seasons[k].Name == "" {
			return fmt.Errorf(
				"simulations[%d].series[%d].seasons[%d]: name is required",
				simIdx, serIdx, k,
			)
		}

		if err := validateEvents(simIdx, serIdx, k, seasons[k].Events); err != nil {
			return err
		}
	}

	return nil
}

func validateEvents(simIdx, serIdx, snIdx int, events []EventConfig) error {
	for i := range events {
		if events[i].Name == "" {
			return fmt.Errorf(
				"simulations[%d].series[%d].seasons[%d].events[%d]: name is required",
				simIdx, serIdx, snIdx, i,
			)
		}

		if events[i].Date != "" {
			if _, err := time.Parse(time.DateOnly, events[i].Date); err != nil {
				//nolint:lll // readability
				return fmt.Errorf(
					"simulations[%d].series[%d].seasons[%d].events[%d]: invalid date %q (expected YYYY-MM-DD)",
					simIdx,
					serIdx,
					snIdx,
					i,
					events[i].Date,
				)
			}
		}

		if err := validateRaces(simIdx, serIdx, snIdx, i, events[i].Races); err != nil {
			return err
		}
	}

	return nil
}

func validateRaces(simIdx, serIdx, snIdx, evIdx int, races []RaceConfig) error {
	for i := range races {
		if races[i].Name == "" {
			return fmt.Errorf(
				"simulations[%d].series[%d].seasons[%d].events[%d].races[%d]: name is required",
				simIdx, serIdx, snIdx, evIdx, i,
			)
		}
	}

	return nil
}

func validateManufacturers(items []CarManufacturerConfig) error {
	for i := range items {
		if items[i].Name == "" {
			return fmt.Errorf("carManufacturers[%d]: name is required", i)
		}

		if err := validateBrands(i, items[i].Brands); err != nil {
			return err
		}
	}

	return nil
}

func validateBrands(mfrIdx int, brands []BrandConfig) error {
	for j := range brands {
		if brands[j].Name == "" {
			return fmt.Errorf(
				"carManufacturers[%d].brands[%d]: name is required",
				mfrIdx, j,
			)
		}

		if err := validateModels(mfrIdx, j, brands[j].Models); err != nil {
			return err
		}
	}

	return nil
}

func validateModels(mfrIdx, brandIdx int, models []ModelConfig) error {
	for k := range models {
		if models[k].Name == "" {
			return fmt.Errorf(
				"carManufacturers[%d].brands[%d].models[%d]: name is required",
				mfrIdx, brandIdx, k,
			)
		}
	}

	return nil
}

func validateTracks(items []TrackConfig) error {
	for i := range items {
		if items[i].Name == "" {
			return fmt.Errorf("tracks[%d]: name is required", i)
		}

		for j := range items[i].Layouts {
			if items[i].Layouts[j].Name == "" {
				return fmt.Errorf(
					"tracks[%d].layouts[%d]: name is required", i, j,
				)
			}
		}
	}

	return nil
}
