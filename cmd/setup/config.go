package setup

import (
	"fmt"
	"os"
	"time"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	"gopkg.in/yaml.v3"
)

// SetupConfig is the root configuration for the setup command.
type SetupConfig struct {
	Drivers          []DriverConfig          `yaml:"drivers"`
	PointSystems     []PointSystemConfig     `yaml:"pointSystems"`
	Simulations      []SimulationConfig      `yaml:"simulations"`
	CarManufacturers []CarManufacturerConfig `yaml:"carManufacturers"`
	CarClasses       []CarClassConfig        `yaml:"carClasses"`
	Tracks           []TrackConfig           `yaml:"tracks"`
}

// PointSystemConfig defines a point system to be created.
type PointSystemConfig struct {
	Name         string                    `yaml:"name"`
	Description  string                    `yaml:"description"`
	Eligible     PointEligibilityConfig    `yaml:"eligible"`
	RaceSettings []PointRaceSettingsConfig `yaml:"raceSettings"`
}

// PointEligibilityConfig defines point eligibility filters.
type PointEligibilityConfig struct {
	Guests      bool    `yaml:"guests"`
	RaceDistPct float64 `yaml:"raceDistPct"`
}

// PointRaceSettingsConfig defines points configuration for one race setup.
type PointRaceSettingsConfig struct {
	Name     string                      `yaml:"name"`
	Policies []PointPolicySettingsConfig `yaml:"policies"`
}

// PointPolicySettingsConfig defines one point policy and its configuration.
type PointPolicySettingsConfig struct {
	Name     string                       `yaml:"name"`
	Points   [][]int32                    `yaml:"points"`
	Settings []ThresholdPenaltyRuleConfig `yaml:"settings"`
}

// ThresholdPenaltyRuleConfig defines threshold-based penalty settings.
type ThresholdPenaltyRuleConfig struct {
	Threshold      uint32  `yaml:"threshold"`
	PenaltyPercent float64 `yaml:"penaltyPercent"`
}

// EntitySimulationConfig holds a simulation reference
// with optional aliases for a specific entity.
type EntitySimulationConfig struct {
	Name    string   `yaml:"name"`
	Aliases []string `yaml:"aliases"`
}

// DriverConfig defines a driver to be created.
type DriverConfig struct {
	Name        string                   `yaml:"name"`
	ExternalID  string                   `yaml:"externalId"`
	IsActive    bool                     `yaml:"isActive"`
	Simulations []EntitySimulationConfig `yaml:"simulations"`
}

// SimulationConfig defines a simulation and its child series.
type SimulationConfig struct {
	Name             string                            `yaml:"name"`
	IsActive         bool                              `yaml:"isActive"`
	SupportedFormats []SimulationSupportedFormatConfig `yaml:"supportedFormats"`
	Series           []SeriesConfig                    `yaml:"series"`
}

// SimulationSupportedFormatConfig defines one supported import format for a simulation.
type SimulationSupportedFormatConfig struct {
	Format               string `yaml:"format"`
	AllowMultipleUploads bool   `yaml:"allowMultipleUploads"`
}

// SeriesConfig defines a series and its child seasons.
type SeriesConfig struct {
	Name    string         `yaml:"name"`
	Seasons []SeasonConfig `yaml:"seasons"`
}

// SeasonConfig defines a season and its associated point system name.
type SeasonConfig struct {
	Name                   string                        `yaml:"name"`
	PointSystem            string                        `yaml:"pointSystem"`
	StartsAt               string                        `yaml:"startsAt"`
	EndsAt                 string                        `yaml:"endsAt"`
	HasTeams               bool                          `yaml:"hasTeams"`
	Multiclass             bool                          `yaml:"multiclass"`
	TeamBased              bool                          `yaml:"teamBased"`
	TeamPointsTopN         int32                         `yaml:"teamPointsTopN"`
	SkipEvents             int32                         `yaml:"skipEvents"`
	DefaultCarModelVariant string                        `yaml:"defaultCarModelVariant"`
	DefaultJoinedAt        string                        `yaml:"defaultJoinedAt"`
	Drivers                []SeasonDriverConfig          `yaml:"drivers"`
	Events                 []EventConfig                 `yaml:"events"`
	Teams                  []TeamConfig                  `yaml:"teams"`
	CarClasses             []SeasonCarClassConfig        `yaml:"carClasses"`
	CarModelVariants       []SeasonCarModelVariantConfig `yaml:"carModelVariants"`
}

// SeasonCarClassConfig defines a car class under a season.
type SeasonCarClassConfig struct {
	Name string `yaml:"name"`
}
type SeasonCarModelVariantConfig struct {
	Name string `yaml:"name"`
}
type SeasonDriverConfig struct {
	Name            string `yaml:"name"`
	CarModelVariant string `yaml:"carModelVariant"`
	CarNumber       string `yaml:"carNumber"`
	JoinedAt        string `yaml:"joinedAt"`
	LeftAt          string `yaml:"leftAt"`
}

// TeamConfig defines a team under a season.
type TeamConfig struct {
	Name            string             `yaml:"name"`
	CarModelVariant string             `yaml:"carModelVariant"`
	CarNumber       string             `yaml:"carNumber"`
	Drivers         []TeamDriverConfig `yaml:"drivers"`
	JoinedAt        string             `yaml:"joinedAt"`
	LeftAt          string             `yaml:"leftAt"`
}

// TeamDriverConfig defines a driver under a team.
type TeamDriverConfig struct {
	Name     string `yaml:"name"`
	JoinedAt string `yaml:"joinedAt"`
	LeftAt   string `yaml:"leftAt"`
}

// EventConfig defines an event under a season.
type EventConfig struct {
	Name            string                     `yaml:"name"`
	SequenceNo      uint32                     `yaml:"sequenceNo"`
	TrackLayout     string                     `yaml:"trackLayout"`
	Date            string                     `yaml:"date"`
	Status          EventStatusConfig          `yaml:"status"`
	ProcessingState EventProcessingStateConfig `yaml:"processingState"`
	Races           []RaceConfig               `yaml:"races"`
}

// EventStatusConfig stores backend.common.v1.EventStatus enum literals from YAML.
type EventStatusConfig string

// EventProcessingStateConfig stores backend.common.v1.EventProcessingState
// enum literals from YAML.
type EventProcessingStateConfig string

func (s EventStatusConfig) String() string {
	if s == "" {
		return commonv1.EventStatus_EVENT_STATUS_UNSPECIFIED.String()
	}

	return string(s)
}

func (s EventProcessingStateConfig) String() string {
	if s == "" {
		return commonv1.EventProcessingState_EVENT_PROCESSING_STATE_UNSPECIFIED.String()
	}

	return string(s)
}

// RaceConfig defines a race under an event.
type RaceConfig struct {
	Name        string           `yaml:"name"`
	SessionType string           `yaml:"sessionType"`
	SequenceNo  int32            `yaml:"sequenceNo"`
	Grids       []RaceGridConfig `yaml:"grids"`
}

// RaceGridConfig defines a race grid under a race.
type RaceGridConfig struct {
	Name        string `yaml:"name"`
	SessionType string `yaml:"sessionType"`
	SequenceNo  int32  `yaml:"sequenceNo"`
}

// CarManufacturerConfig defines a car manufacturer and its child brands.
type CarManufacturerConfig struct {
	Name   string        `yaml:"name"`
	Brands []BrandConfig `yaml:"brands"`
	Models []ModelConfig `yaml:"models"`
}

// BrandConfig defines a car brand and its child models.

// Deprecated: Use ModelConfig instead.
type BrandConfig struct {
	Name   string           `yaml:"name"`
	Models []ModelConfigOld `yaml:"models"`
}

// ModelConfig defines a car model.
type ModelConfig struct {
	Name     string               `yaml:"name"`
	Variants []ModelVariantConfig `yaml:"variants"`
}

// ModelConfig defines a car model.
type ModelVariantConfig struct {
	Name        string                   `yaml:"name"`
	Simulations []EntitySimulationConfig `yaml:"simulations"`
}

// ModelConfigOld defines a car model.
//
// Deprecated: Use ModelConfig instead.
type ModelConfigOld struct {
	Name        string                   `yaml:"name"`
	Simulations []EntitySimulationConfig `yaml:"simulations"`
}

// TrackConfig defines a track and its child layouts.
type TrackConfig struct {
	Name    string         `yaml:"name"`
	Layouts []LayoutConfig `yaml:"layouts"`
}

// LayoutConfig defines a track layout.
type LayoutConfig struct {
	Name        string                   `yaml:"name"`
	Simulations []EntitySimulationConfig `yaml:"simulations"`
}

// CarClassConfig defines a car class.
type CarClassConfig struct {
	Name          string                `yaml:"name"`
	ModelVariants []CarClassModelConfig `yaml:"modelVariants"`
}
type CarClassModelConfig struct {
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
	for i := range items {
		ps := items[i]
		if ps.Name == "" {
			return fmt.Errorf("pointSystems[%d]: name is required", i)
		}

		for j := range ps.RaceSettings {
			for k := range ps.RaceSettings[j].Policies {
				policy := ps.RaceSettings[j].Policies[k]
				if policy.Name == "" {
					return fmt.Errorf(
						"pointSystems[%d].raceSettings[%d].policies[%d]: name is required",
						i,
						j,
						k,
					)
				}

				if _, ok := commonv1.PointPolicy_value[policy.Name]; !ok {
					return fmt.Errorf(
						"pointSystems[%d].raceSettings[%d].policies[%d]: unknown policy %q",
						i,
						j,
						k,
						policy.Name,
					)
				}
			}
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

		if _, err := parseImportConfigsFromSetup(items[i].SupportedFormats); err != nil {
			return fmt.Errorf("simulations[%d].supportedFormats: %w", i, err)
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

//nolint:lll,funlen // readability
func validateSeasonList(simIdx, serIdx int, seasons []SeasonConfig) error {
	for k := range seasons {
		if seasons[k].Name == "" {
			return fmt.Errorf(
				"simulations[%d].series[%d].seasons[%d]: name is required",
				simIdx, serIdx, k,
			)
		}

		var startsAt time.Time
		if seasons[k].StartsAt != "" {
			var parseErr error
			startsAt, parseErr = parseSeasonTimestamp(seasons[k].StartsAt)
			if parseErr != nil {
				return fmt.Errorf(
					"simulations[%d].series[%d].seasons[%d]: %w",
					simIdx,
					serIdx,
					k,
					parseErr,
				)
			}
		}

		if seasons[k].EndsAt != "" {
			endsAt, parseErr := parseSeasonTimestamp(seasons[k].EndsAt)
			if parseErr != nil {
				return fmt.Errorf(
					"simulations[%d].series[%d].seasons[%d].ends_at: %w",
					simIdx,
					serIdx,
					k,
					parseErr,
				)
			}

			if !startsAt.IsZero() && endsAt.Before(startsAt) {
				return fmt.Errorf(
					"simulations[%d].series[%d].seasons[%d]: endsAt must be greater than or equal to startsAt",
					simIdx,
					serIdx,
					k,
				)
			}
		}

		if err := validateEvents(simIdx, serIdx, k, seasons[k].Events); err != nil {
			return err
		}
	}

	return nil
}

func parseSeasonTimestamp(value string) (time.Time, error) {
	if ts, err := time.Parse(time.RFC3339, value); err == nil {
		return ts, nil
	}

	ts, err := time.Parse(time.DateOnly, value)
	if err == nil {
		return ts, nil
	}

	return time.Time{}, fmt.Errorf(
		"invalid value %q (expected RFC3339 or YYYY-MM-DD)",
		value,
	)
}

//nolint:lll // readability
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

		if err := validateRaceGrids(
			simIdx,
			serIdx,
			snIdx,
			evIdx,
			i,
			races[i].Grids,
		); err != nil {
			return err
		}
	}

	return nil
}

//nolint:whitespace // editor/linter issue
func validateRaceGrids(
	simIdx, serIdx, snIdx, evIdx, raceIdx int,
	grids []RaceGridConfig,
) error {
	for i := range grids {
		if grids[i].Name == "" {
			return fmt.Errorf(
				"simulations[%d].series[%d].seasons[%d].events[%d].races[%d]."+
					"grids[%d]: name is required",
				simIdx,
				serIdx,
				snIdx,
				evIdx,
				raceIdx,
				i,
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

func validateModels(mfrIdx, brandIdx int, models []ModelConfigOld) error {
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
