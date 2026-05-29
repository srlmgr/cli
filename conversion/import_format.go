package conversion

import (
	"fmt"
	"strconv"
	"strings"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

//nolint:goconst // by design
var supportedImportFormatLiterals = []string{"json", "csv", "xml"}

// ParseImportConfigs converts CLI values into protobuf ImportConfig entries.
// Supported input format is `import-format` or `import-format:true|false`.
func ParseImportConfigs(values []string) ([]*commonv1.ImportConfig, error) {
	if len(values) == 0 {
		return nil, nil
	}

	configs := make([]*commonv1.ImportConfig, 0, len(values))
	for _, value := range values {
		parts := strings.SplitN(strings.TrimSpace(value), ":", 2)
		format, err := ParseImportFormat(parts[0])
		if err != nil {
			return nil, err
		}

		allowMultipleUploads := false
		if len(parts) == 2 {
			allowMultipleUploads, err = strconv.ParseBool(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf(
					"invalid import config %q: expected format or format:true|false",
					value,
				)
			}
		}

		configs = append(configs, &commonv1.ImportConfig{
			Format:               format,
			AllowMultipleUploads: allowMultipleUploads,
		})
	}

	return configs, nil
}

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

// JoinImportConfigs renders protobuf import config entries in user-facing format.
func JoinImportConfigs(values []*commonv1.ImportConfig) string {
	if len(values) == 0 {
		return ""
	}

	items := make([]string, 0, len(values))
	for _, value := range values {
		items = append(items, fmt.Sprintf(
			"%s(multi=%t)",
			importFormatLiteral(value.GetFormat()),
			value.GetAllowMultipleUploads(),
		))
	}

	return strings.Join(items, ", ")
}

func importFormatLiteral(value commonv1.ImportFormat) string {
	switch value {
	case commonv1.ImportFormat_IMPORT_FORMAT_JSON:
		return "json"
	case commonv1.ImportFormat_IMPORT_FORMAT_CSV:
		return "csv"
	case commonv1.ImportFormat_IMPORT_FORMAT_XML:
		return "xml"
	case commonv1.ImportFormat_IMPORT_FORMAT_UNSPECIFIED:
		return "unspecified"
	default:
		return strings.ToLower(strings.TrimPrefix(value.String(), "IMPORT_FORMAT_"))
	}
}
