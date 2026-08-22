package conversion

import (
	"fmt"
	"strings"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

var SupportedSkipModeLiterals = []string{"never", "always", "when-applicable"}

// ParseSkipMode converts string literals from CLI input into protobuf enum values.
func ParseSkipMode(value string) (commonv1.SkipMode, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "never":
		return commonv1.SkipMode_SKIP_MODE_NEVER, nil
	case "always":
		return commonv1.SkipMode_SKIP_MODE_ALWAYS, nil
	case "when-applicable":
		return commonv1.SkipMode_SKIP_MODE_WHEN_APPLICABLE, nil
	default:
		return commonv1.SkipMode_SKIP_MODE_UNSPECIFIED, fmt.Errorf(
			"unsupported skip mode %q (supported: %s)",
			value,
			strings.Join(SupportedSkipModeLiterals, ", "),
		)
	}
}
