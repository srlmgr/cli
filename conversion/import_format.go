package conversion

import (
	"fmt"
	"strings"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

var supportedImportFormatLiterals = []string{"json", "csv"}

// ParseImportFormats converts string literals from CLI input into protobuf enum values.
func ParseImportFormats(values []string) ([]commonv1.ImportFormat, error) {
	if len(values) == 0 {
		return nil, nil
	}

	formats := make([]commonv1.ImportFormat, 0, len(values))
	for _, value := range values {
		format, err := ParseImportFormat(value)
		if err != nil {
			return nil, err
		}
		formats = append(formats, format)
	}

	return formats, nil
}

// ParseImportFormat converts one string literal into a protobuf enum value.
func ParseImportFormat(value string) (commonv1.ImportFormat, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if normalized == "" {
		return commonv1.ImportFormat_IMPORT_FORMAT_UNSPECIFIED,
			fmt.Errorf(
				"import format cannot be empty (supported: %s)",
				strings.Join(supportedImportFormatLiterals, ", "),
			)
	}

	if enumValue, ok := commonv1.ImportFormat_value[normalized]; ok {
		return commonv1.ImportFormat(enumValue), nil
	}

	if !strings.HasPrefix(normalized, "IMPORT_FORMAT_") {
		if enumValue, ok := commonv1.ImportFormat_value["IMPORT_FORMAT_"+normalized]; ok {
			return commonv1.ImportFormat(enumValue), nil
		}
	}

	return commonv1.ImportFormat_IMPORT_FORMAT_UNSPECIFIED, fmt.Errorf(
		"unsupported import format %q (supported: %s)",
		value,
		strings.Join(supportedImportFormatLiterals, ", "),
	)
}

// JoinImportFormats renders protobuf enum values in user-facing format.
func JoinImportFormats(values []commonv1.ImportFormat) string {
	if len(values) == 0 {
		return ""
	}

	literals := make([]string, 0, len(values))
	for _, value := range values {
		literals = append(literals, importFormatLiteral(value))
	}

	return strings.Join(literals, ", ")
}

func importFormatLiteral(value commonv1.ImportFormat) string {
	switch value {
	case commonv1.ImportFormat_IMPORT_FORMAT_JSON:
		return "json"
	case commonv1.ImportFormat_IMPORT_FORMAT_CSV:
		return "csv"
	case commonv1.ImportFormat_IMPORT_FORMAT_UNSPECIFIED:
		return "unspecified"
	default:
		return strings.ToLower(strings.TrimPrefix(value.String(), "IMPORT_FORMAT_"))
	}
}
