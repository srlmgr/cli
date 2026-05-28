package setup

import (
	"fmt"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

//nolint:whitespace // editor/linter issue
func buildCreatePointSystemRequest(
	cfg PointSystemConfig,
) (*commandv1.CreatePointSystemRequest, error) {
	eligibility, raceSettings, err := buildPointSystemMessages(cfg)
	if err != nil {
		return nil, err
	}

	return &commandv1.CreatePointSystemRequest{
		Name:         cfg.Name,
		Description:  cfg.Description,
		Eligibility:  eligibility,
		RaceSettings: raceSettings,
	}, nil
}

//nolint:whitespace // editor/linter issue
func buildPointSystemMessages(
	cfg PointSystemConfig,
) (*commonv1.PointEligibility, []*commonv1.PointRaceSettings, error) {
	eligibility := &commonv1.PointEligibility{
		Guests:                 cfg.Eligible.Guests,
		MinRaceDistancePercent: cfg.Eligible.RaceDistPct,
	}

	raceSettings := make([]*commonv1.PointRaceSettings, 0, len(cfg.RaceSettings))
	for i := range cfg.RaceSettings {
		raceCfg := cfg.RaceSettings[i]
		policies := make([]*commonv1.PointPolicySettings, 0, len(raceCfg.Policies))

		for j := range raceCfg.Policies {
			policy, err := buildPointPolicySettings(raceCfg.Policies[j])
			if err != nil {
				return nil, nil, fmt.Errorf(
					"raceSettings[%d].policies[%d]: %w",
					i,
					j,
					err,
				)
			}

			policies = append(policies, policy)
		}

		raceSettings = append(raceSettings, &commonv1.PointRaceSettings{
			Name:     raceCfg.Name,
			Policies: policies,
		})
	}

	return eligibility, raceSettings, nil
}

//nolint:whitespace // editor/linter issue
func buildPointPolicySettings(
	policyCfg PointPolicySettingsConfig,
) (*commonv1.PointPolicySettings, error) {
	policyNum, ok := commonv1.PointPolicy_value[policyCfg.Name]
	if !ok {
		return nil, fmt.Errorf("unknown policy %q", policyCfg.Name)
	}

	policy := &commonv1.PointPolicySettings{
		Name: commonv1.PointPolicy(policyNum),
	}
	//nolint:exhaustive // by design
	switch policy.GetName() {
	case commonv1.PointPolicy_POINT_POLICY_FINISH_POS:
		policy.Config = &commonv1.PointPolicySettings_FinishPos{
			FinishPos: buildPositionPointsConfig(policyCfg.Points),
		}
	case commonv1.PointPolicy_POINT_POLICY_QUALIFICATION_POS:
		policy.Config = &commonv1.PointPolicySettings_QualificationPos{
			QualificationPos: buildPositionPointsConfig(policyCfg.Points),
		}
	case commonv1.PointPolicy_POINT_POLICY_LEAST_INCIDENTS:
		policy.Config = &commonv1.PointPolicySettings_LeastIncidents{
			LeastIncidents: buildPositionPointsConfig(policyCfg.Points),
		}
	case commonv1.PointPolicy_POINT_POLICY_FASTEST_LAP:
		policy.Config = &commonv1.PointPolicySettings_FastestLap{
			FastestLap: buildPositionPointsConfig(policyCfg.Points),
		}
	case commonv1.PointPolicy_POINT_POLICY_TOP_N_FINISHER:
		policy.Config = &commonv1.PointPolicySettings_TopNFinisher{
			TopNFinisher: buildPositionPointsConfig(policyCfg.Points),
		}
	case commonv1.PointPolicy_POINT_POLICY_INCIDENTS_EXCEEDED:
		policy.Config = &commonv1.PointPolicySettings_IncidentsExceeded{
			IncidentsExceeded: buildThresholdPenaltyConfig(policyCfg.Settings),
		}
	}

	return policy, nil
}

func buildPositionPointsConfig(points [][]int32) *commonv1.PositionPointsConfig {
	tables := make([]*commonv1.PointTable, 0, len(points))
	for i := range points {
		tableValues := append([]int32(nil), points[i]...)
		tables = append(tables, &commonv1.PointTable{Values: tableValues})
	}

	return &commonv1.PositionPointsConfig{Tables: tables}
}

//nolint:whitespace // editor/linter issue
func buildThresholdPenaltyConfig(
	settings []ThresholdPenaltyRuleConfig,
) *commonv1.ThresholdPenaltyConfig {
	rules := make([]*commonv1.ThresholdPenaltyRule, 0, len(settings))
	for i := range settings {
		rules = append(rules, &commonv1.ThresholdPenaltyRule{
			Threshold:      settings[i].Threshold,
			PenaltyPercent: settings[i].PenaltyPercent,
		})
	}

	return &commonv1.ThresholdPenaltyConfig{Rules: rules}
}
