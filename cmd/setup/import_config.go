package setup

import (
	"fmt"
	"strings"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"

	"github.com/srlmgr/cli/conversion"
)

func parseImportConfigsFromSetup(
	values []SimulationSupportedFormatConfig,
) ([]*commonv1.ImportConfig, error) {
	if len(values) == 0 {
		return nil, nil
	}

	configs := make([]*commonv1.ImportConfig, 0, len(values))
	for i := range values {
		format, err := conversion.ParseImportFormat(values[i].Format)
		if err != nil {
			return nil, fmt.Errorf("item[%d].importFormat: %w", i, err)
		}

		configs = append(configs, &commonv1.ImportConfig{
			Format:               format,
			AllowMultipleUploads: values[i].AllowMultipleUploads,
		})
	}

	if err := validateUniqueImportConfigFormats(configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func validateUniqueImportConfigFormats(values []*commonv1.ImportConfig) error {
	seen := make(map[commonv1.ImportFormat]struct{}, len(values))
	for i := range values {
		format := values[i].GetFormat()
		if _, exists := seen[format]; exists {
			return fmt.Errorf(
				"duplicate format %q",
				strings.ToLower(strings.TrimPrefix(format.String(), "IMPORT_FORMAT_")),
			)
		}
		seen[format] = struct{}{}
	}

	return nil
}
