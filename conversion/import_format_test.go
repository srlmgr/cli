package conversion

import (
	"testing"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

func TestParseImportFormats(t *testing.T) {
	t.Parallel()

	formats, err := ParseImportFormats([]string{"json", "IMPORT_FORMAT_CSV"})
	if err != nil {
		t.Fatalf("ParseImportFormats returned error: %v", err)
	}

	if len(formats) != 2 {
		t.Fatalf("unexpected number of parsed formats: got %d want 2", len(formats))
	}

	if formats[0] != commonv1.ImportFormat_IMPORT_FORMAT_JSON {
		t.Fatalf("unexpected first format: got %v", formats[0])
	}

	if formats[1] != commonv1.ImportFormat_IMPORT_FORMAT_CSV {
		t.Fatalf("unexpected second format: got %v", formats[1])
	}
}

func TestParseImportFormat_Unknown(t *testing.T) {
	t.Parallel()

	_, err := ParseImportFormat("xml")
	if err == nil {
		t.Fatal("ParseImportFormat should return an error for unknown formats")
	}
}

func TestJoinImportFormats(t *testing.T) {
	t.Parallel()

	joined := JoinImportFormats([]commonv1.ImportFormat{
		commonv1.ImportFormat_IMPORT_FORMAT_JSON,
		commonv1.ImportFormat_IMPORT_FORMAT_CSV,
	})

	if joined != "json, csv" {
		t.Fatalf("unexpected joined import formats: got %q", joined)
	}
}
