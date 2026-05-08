package formatters

import (
	"fmt"
	"strings"

	"zgluco/internal/models/profile"
	"zgluco/internal/models/treatments"
	"zgluco/internal/types"
)

func FormatTreatments(sb *strings.Builder, p *profile.Profile, t []treatments.Treatment, tr types.TimeRange) {
	fmt.Fprintf(sb, "═══ Treatments (%d events) ═══\n\n", len(t))
	fmt.Fprint(sb, tr.String())

	for _, treatment := range t {
		if tt, ok := treatment.(treatments.TemporaryTargetEvent); ok {
			sb.WriteString(tt.StringWithUnit(p.PreferredUnits))
		} else {
			sb.WriteString(treatment.String())
		}
		sb.WriteRune('\n')
	}
}
