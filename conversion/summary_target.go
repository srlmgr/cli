package conversion

import (
	"fmt"
	"strings"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

var supportedSummaryTargetLiterals = []string{"driver", "team"}

// ParseSummaryTarget converts string literals from CLI input into protobuf enum values.
func ParseSummaryTarget(value string) (commonv1.SummaryTargetType, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "driver":
		return commonv1.SummaryTargetType_SUMMARY_TARGET_TYPE_DRIVER, nil
	case "team":
		return commonv1.SummaryTargetType_SUMMARY_TARGET_TYPE_TEAM, nil
	default:
		return commonv1.SummaryTargetType_SUMMARY_TARGET_TYPE_UNSPECIFIED, fmt.Errorf(
			"unsupported summary target type %q (supported: %s)",
			value,
			strings.Join(supportedSummaryTargetLiterals, ", "),
		)
	}
}
