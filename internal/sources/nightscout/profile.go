package nightscout

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"
	"zgluco/internal/models"
	"zgluco/internal/models/profile"
	"zgluco/internal/types"
)

var knownUploaders = []string{"Trio", "Loop"}

func (n *Nightscout) fetchProfile(timeRange types.TimeRange) (*profile.Profile, error) {
	res, err := n.doApiCall("profiles", "startDate", "startDate", timeRange)
	if err != nil {
		return nil, fmt.Errorf("could not fetch profiles: %w", err)
	}

	var profiles []*Profile
	err = json.Unmarshal(res, &profiles)
	if err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, errors.New("no profiles found")
	}

	var changes []*profile.Change

	var unknownUploaders = make(map[string]struct{})
	if !slices.Contains(knownUploaders, profiles[0].EnteredBy) {
		unknownUploaders[profiles[0].EnteredBy] = struct{}{}
	}

	if len(profiles) > 1 {
		// We build the changelog based on the oldest entry
		for i := 1; i < len(profiles); i++ {
			// TODO: Future: Detect and handle sudden unit changes (normalize based on latest profile?)

			previousProfile := profiles[i-1].GetDefaultProfile()
			currentProfile := profiles[i].GetDefaultProfile()

			if !slices.Contains(knownUploaders, profiles[i].EnteredBy) {
				unknownUploaders[profiles[i].EnteredBy] = struct{}{}
			}

			currentStartDate := profiles[i].StartDate

			diaDiff := currentProfile.DIA - previousProfile.DIA
			if diaDiff != 0 {
				changes = append(changes, &profile.Change{
					Type:        profile.DIA,
					EffectiveAt: currentStartDate,
					Diffs: []profile.TimeValueDiff{
						{
							OldValue: previousProfile.DIA,
							NewValue: currentProfile.DIA,
						},
					},
				})
			}

			// CR, ISF, Basal
			changes = append(changes, compareTimeValueEntries(profile.CarbRatio, currentStartDate, currentProfile.CarbRatio, previousProfile.CarbRatio))
			changes = append(changes, compareTimeValueEntries(profile.ISF, currentStartDate, currentProfile.Sens, previousProfile.Sens))
			changes = append(changes, compareTimeValueEntries(profile.Basal, currentStartDate, currentProfile.Basal, previousProfile.Basal))
			changes = append(changes, compareTimeValueEntries(profile.TargetLow, currentStartDate, currentProfile.TargetLow, previousProfile.TargetLow))
			changes = append(changes, compareTimeValueEntries(profile.TargetHigh, currentStartDate, currentProfile.TargetHigh, previousProfile.TargetHigh))

			// Remove any nil entries
			changes = slices.DeleteFunc(changes, func(c *profile.Change) bool {
				return c == nil
			})
		}
	}

	newestProfile := profiles[len(profiles)-1].GetDefaultProfile()

	convertToRate := func(values []TimeValueEntry) []*profile.TimeValue {
		rates := make([]*profile.TimeValue, len(values))
		for i := range values {
			rates[i] = &profile.TimeValue{
				StartsAtSeconds: values[i].TimeAsSeconds,
				Value:           values[i].Value,
			}
		}

		return rates
	}

	convertToTarget := func(lowValues []TimeValueEntry, highValues []TimeValueEntry) []*profile.GlucoseTarget {
		targets := make([]*profile.GlucoseTarget, len(lowValues)/2)

		lowValuesByTime := make(map[int]float32, len(lowValues))
		for _, e := range lowValues {
			lowValuesByTime[e.TimeAsSeconds] = e.Value
		}

		highValuesByTime := make(map[int]float32, len(highValues))
		for _, e := range highValues {
			highValuesByTime[e.TimeAsSeconds] = e.Value
		}

		for t, _ := range lowValuesByTime {
			targets = append(targets, &profile.GlucoseTarget{
				StartsAtSeconds: t,
				LowerBound:      lowValuesByTime[t],
				UpperBound:      highValuesByTime[t],
			})
		}

		return targets
	}

	if len(unknownUploaders) > 0 {
		uploaders := slices.Sorted(maps.Keys(unknownUploaders))
		fmt.Printf("WARN: Found the following unknown uploaders %v\n", uploaders)
		fmt.Printf("Cannot guarantee proper parsing. Please open an issue with an example to ensure proper support.\n")
	}

	return &profile.Profile{
		PreferredUnits:            models.ParseBgUnitType(newestProfile.Units),
		BasalRates:                convertToRate(newestProfile.Basal),
		CarbRatios:                convertToRate(newestProfile.CarbRatio),
		InsulinSensitivityFactors: convertToRate(newestProfile.Sens),
		GlucoseTargets:            convertToTarget(newestProfile.TargetLow, newestProfile.TargetHigh),
		Changes:                   changes,
	}, nil
}

func compareTimeValueEntries(
	changeType profile.ChangeType,
	effectiveAt *time.Time,
	newEntries []TimeValueEntry,
	oldEntries []TimeValueEntry,
) *profile.Change {
	// Index old entries by TimeAsSeconds for O(1) lookup
	oldByTime := make(map[int]float32, len(oldEntries))
	for _, e := range oldEntries {
		oldByTime[e.TimeAsSeconds] = e.Value
	}

	// Index new entries the same way to detect removals
	newByTime := make(map[int]float32, len(newEntries))
	for _, e := range newEntries {
		newByTime[e.TimeAsSeconds] = e.Value
	}

	var diffs []profile.TimeValueDiff

	// Changed or added slots
	for _, e := range newEntries {
		oldVal, existed := oldByTime[e.TimeAsSeconds]
		if !existed || oldVal != e.Value {
			diffs = append(diffs, profile.TimeValueDiff{
				TimeAsSeconds: e.TimeAsSeconds,
				OldValue:      oldVal,
				NewValue:      e.Value,
				Removed:       existed && oldVal > 0 && e.Value == 0,
			})
		}
	}

	// Slots that disappeared entirely from the new profile
	for _, e := range oldEntries {
		if _, exists := newByTime[e.TimeAsSeconds]; !exists {
			diffs = append(diffs, profile.TimeValueDiff{
				TimeAsSeconds: e.TimeAsSeconds,
				OldValue:      e.Value,
				NewValue:      0,
				Removed:       true,
			})
		}
	}

	if len(diffs) == 0 {
		return nil
	}

	// Sort by time slot for deterministic output
	slices.SortFunc(diffs, func(a, b profile.TimeValueDiff) int {
		return cmp.Compare(a.TimeAsSeconds, b.TimeAsSeconds)
	})

	return &profile.Change{
		Type:        changeType,
		EffectiveAt: effectiveAt,
		Diffs:       diffs,
	}
}
