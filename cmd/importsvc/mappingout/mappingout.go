// Package mappingout renders unresolved mappings including their detailed cause.
package mappingout

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
)

// Write prints unresolved mappings with common fields and concise error details.
func Write(out io.Writer, mappings []*commonv1.UnresolvedMapping) error {
	if len(mappings) == 0 {
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "SOURCE\tTYPE\tENTRY ID\tRESULT IDX\tERROR"); err != nil {
		return fmt.Errorf("write unresolved mapping header: %w", err)
	}

	for _, m := range mappings {
		entryID, resultIdx := "", ""
		if m.HasResultEntryId() {
			entryID = fmt.Sprintf("%d", m.GetResultEntryId())
		}
		if m.HasResultIndex() {
			resultIdx = fmt.Sprintf("%d", m.GetResultIndex())
		}

		if _, err := fmt.Fprintf(w,
			"%q\t%s\t%s\t%s\t%s\n",
			m.GetSourceValue(), m.GetMappingType(), entryID, resultIdx, Reason(m),
		); err != nil {
			return fmt.Errorf("write unresolved mapping output: %w", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("write unresolved mapping output: %w", err)
	}
	return nil
}

func source(e *commonv1.SourceEntity) string {
	return fmt.Sprintf("%q (id=%s)", e.GetName(), e.GetId())
}

// Reason returns a human readable description of the mapping failure cause.
//
//nolint:gocyclo // one branch per reason
func Reason(m *commonv1.UnresolvedMapping) string {
	switch {
	case m.HasSourceEntityNotFound():
		r := m.GetSourceEntityNotFound()
		return fmt.Sprintf("Entity not found: %s %s",
			r.GetEntityType(), source(r.GetSource()))
	case m.HasSeasonDriverNotActive():
		r := m.GetSeasonDriverNotActive()
		return fmt.Sprintf("Driver not active: #%d %s",
			r.GetDriverId(), source(r.GetSourceDriver()))
	case m.HasSeasonCarNotConfigured():
		r := m.GetSeasonCarNotConfigured()
		return fmt.Sprintf("Car not configured: variant #%d %s",
			r.GetCarModelVariantId(), source(r.GetSourceCar()))
	case m.HasTeamNotFound():
		r := m.GetTeamNotFound()
		drivers := make([]string, 0, len(r.GetDriversChecked()))
		for _, d := range r.GetDriversChecked() {
			drivers = append(drivers,
				fmt.Sprintf("#%d %s", d.GetDriverId(), source(d.GetSourceDriver())))
		}
		return fmt.Sprintf("Team not found: %q (checked: %s)",
			r.GetSourceTeamName(), strings.Join(drivers, ", "))
	case m.HasTeamDriversMissing():
		return fmt.Sprintf("Team drivers missing: %q",
			m.GetTeamDriversMissing().GetSourceTeamName())
	case m.HasCarClassNotConfigured():
		r := m.GetCarClassNotConfigured()
		return fmt.Sprintf("Car class not configured: variant #%d %s",
			r.GetCarModelVariantId(), source(r.GetSourceCar()))
	case m.HasCarNumberMismatch():
		r := m.GetCarNumberMismatch()
		owner := ""
		if r.HasDriverId() {
			owner = fmt.Sprintf("driver #%d, ", r.GetDriverId())
		} else if r.HasTeamId() {
			owner = fmt.Sprintf("team #%d, ", r.GetTeamId())
		}
		return fmt.Sprintf("Car number mismatch: %simported %s, configured %s",
			owner, r.GetImportedNumber(), r.GetConfiguredNumber())
	default:
		return "Unknown mapping error"
	}
}
