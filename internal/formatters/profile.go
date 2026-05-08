package formatters

import (
	"fmt"
	"strings"
	"time"
	"zgluco/internal/models/profile"
	"zgluco/internal/types"
)

func FormatProfile(sb *strings.Builder, p *profile.Profile, tr types.TimeRange, loc *time.Location) {
	unitStr := p.PreferredUnits.String()

	fmt.Fprintf(sb, "═══ Profile (%s) ═══\n\n", unitStr)
	fmt.Fprint(sb, tr.String())

	fmt.Fprintln(sb, "Basal Rates (U/hr)")
	for _, r := range p.BasalRates {
		fmt.Fprintf(sb, "  %s  %.3f U/hr\n", secondsToHHMM(r.StartsAtSeconds), r.Value)
	}
	fmt.Fprintf(sb, "(Total daily: %.2f U)\n", p.GetTotalBasalDose())

	fmt.Fprintln(sb, "\nCarb Ratios (g/U)")
	for _, r := range p.CarbRatios {
		fmt.Fprintf(sb, "  %s  %.1f g/U\n", secondsToHHMM(r.StartsAtSeconds), r.Value)
	}

	fmt.Fprintf(sb, "\nInsulin Sensitivity (%s/U)\n", unitStr)
	for _, r := range p.InsulinSensitivityFactors {
		fmt.Fprintf(sb, "  %s  %.1f %s/U\n", secondsToHHMM(r.StartsAtSeconds), r.Value, unitStr)
	}

	fmt.Fprintf(sb, "\nGlucose Targets (%s)\n", unitStr)
	for _, t := range p.GlucoseTargets {
		fmt.Fprintf(sb, "  %s  %.1f–%.1f %s\n", secondsToHHMM(t.StartsAtSeconds), t.LowerBound, t.UpperBound, unitStr)
	}

	if len(p.Changes) > 0 {
		fmt.Fprintln(sb, "\nChanges (Oldest to Newest)")
		for _, c := range p.Changes {
			t := c.EffectiveAt.In(loc)
			fmt.Fprintf(sb, "\t[%s] %s\n", t.Format("Jan 02 2006 15:04"), c.Type)
			for _, d := range c.Diffs {
				if d.Removed {
					fmt.Fprintf(sb, "\t\t%s\t%.3f -> (removed)\n", secondsToHHMM(d.TimeAsSeconds), d.OldValue)
				} else {
					fmt.Fprintf(sb, "\t\t%s\t%.3f -> %.3f\n", secondsToHHMM(d.TimeAsSeconds), d.OldValue, d.NewValue)
				}
			}
		}
	}
}

func secondsToHHMM(s int) string {
	return fmt.Sprintf("%02d:%02d", s/3600, (s%3600)/60)
}
