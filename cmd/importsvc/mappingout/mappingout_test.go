package mappingout

import (
	"bytes"
	"strings"
	"testing"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

func TestWriteResultColumns(t *testing.T) {
	t.Parallel()

	withIDs := &commonv1.UnresolvedMapping{}
	withIDs.SetSourceValue("with-ids")
	withIDs.SetMappingType("driver")
	withIDs.SetSimulationId(1)
	withIDs.SetEventId(2)
	withIDs.SetRaceGridId(3)
	withIDs.SetResultEntryId(17)
	withIDs.SetResultIndex(4)

	withoutIDs := &commonv1.UnresolvedMapping{}
	withoutIDs.SetSourceValue("without-ids")
	withoutIDs.SetMappingType("team")

	var out bytes.Buffer
	if err := Write(&out, []*commonv1.UnresolvedMapping{withIDs, withoutIDs}); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "ENTRY ID") || !strings.Contains(output, "RESULT IDX") {
		t.Fatalf("expected Entry ID and Result Idx columns, got:\n%s", output)
	}
	for _, ignored := range []string{"Simulation", "Event", "Race grid"} {
		if strings.Contains(output, ignored) {
			t.Errorf("output unexpectedly contains %q:\n%s", ignored, output)
		}
	}

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	header := lines[0]
	entryStart := strings.Index(header, "ENTRY ID")
	resultStart := strings.Index(header, "RESULT IDX")
	errorStart := strings.Index(header, "ERROR")
	if entryStart < 0 || resultStart < 0 || errorStart < 0 {
		t.Fatalf("could not locate result columns in header: %q", header)
	}

	row := func(source string) string {
		for _, line := range lines[1:] {
			if strings.Contains(line, `"`+source+`"`) {
				return line
			}
		}
		t.Fatalf("no output row for source %q:\n%s", source, output)
		return ""
	}
	value := func(line string, start, end int) string {
		return strings.TrimSpace(line[start:end])
	}

	if got := value(row("with-ids"), entryStart, resultStart); got != "17" {
		t.Errorf("Entry ID = %q, want 17", got)
	}
	if got := value(row("with-ids"), resultStart, errorStart); got != "4" {
		t.Errorf("Result Idx = %q, want 4", got)
	}
	if got := value(row("without-ids"), entryStart, resultStart); got != "" {
		t.Errorf("missing Entry ID column = %q, want blank", got)
	}
	if got := value(row("without-ids"), resultStart, errorStart); got != "" {
		t.Errorf("missing Result Idx column = %q, want blank", got)
	}
}
