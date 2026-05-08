package nightscout

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"
	"zgluco/internal/models"
	"zgluco/internal/models/treatments"
	"zgluco/internal/types"
)

func (n *Nightscout) fetchTreatments(timeRange types.TimeRange) ([]treatments.Treatment, error) {
	res, err := n.doApiCall("treatments.json", "created_at", "", timeRange)
	if err != nil {
		return nil, fmt.Errorf("could not fetch treatments: %w", err)
	}

	var events []*Treatment
	err = json.Unmarshal(res, &events)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errors.New("no treatments found")
	}

	var parsedEvents []treatments.Treatment
	for _, e := range events {
		parsed := fromEventType(e)
		if parsed == nil {
			fmt.Printf("Unknown event type: %q\n", e.EventType)
			continue
		}
		parsedEvents = append(parsedEvents, parsed)
	}

	slices.SortFunc(parsedEvents, func(a, b treatments.Treatment) int {
		return a.TreatmentTime().Compare(b.TreatmentTime())
	})

	// First aggregate all consecutive temp zero basal events
	parsedEvents = mergeZeroTempBasals(parsedEvents)

	// Fix durations of single temp basal entries
	parsedEvents = fixTempBasalDurations(parsedEvents)

	// Merge consecutive SMBs and return
	return mergeSMBs(parsedEvents), nil
}

func fromEventType(t *Treatment) treatments.Treatment {
	switch t.EventType {
	case "SMB", "Correction Bolus":
		return treatments.NewSmbTreatment(t.CreatedAt, t.EnteredBy, t.Insulin)
	case "Temp Basal":
		return treatments.NewTempBasalEvent(t.CreatedAt, t.EnteredBy, int(t.Duration), t.Rate)
	case "Carb Correction":
		return treatments.NewCarbCorrectionEvent(t.CreatedAt, t.EnteredBy, t.Insulin, t.InsulinType, models.NewMeal(t.Carbs, t.Fat, t.Protein, t.AbsorptionTime))
	case "Temporary Override":
		return treatments.NewTemporaryOverrideEvent(t.CreatedAt, t.EnteredBy, t.Notes, float32(t.InsulinNeedsScaleFactor), t.Duration)
	case "Temporary Target":
		return treatments.NewTemporaryTargetEvent(t.CreatedAt, t.EnteredBy, t.Notes, fmt.Sprintf("%.0f", t.Duration), t.TargetBottom, t.TargetTop)
	case "Bolus":
		return treatments.NewBolusEvent(t.CreatedAt, t.EnteredBy, t.Insulin)
	case "Site Change":
		return treatments.NewChangeEvent(t.CreatedAt, t.EnteredBy, treatments.Site)
	case "Insulin Change":
		return treatments.NewChangeEvent(t.CreatedAt, t.EnteredBy, treatments.Insulin)
	case "Exercise", "Announcement", "Note":
		return treatments.NewNoteEvent(t.CreatedAt, t.EnteredBy, t.Notes, t.Duration)
	case "External Insulin":
		return treatments.NewExternalInsulinEvent(t.CreatedAt, t.EnteredBy, t.Insulin)
	default:
		return nil
	}
}

// mergeZeroTempBasals collapses consecutive zero-rate temp basal events that
// overlap in time into a single event. AID systems re-enact zero basals before
// they expire, so without merging the same suspended period would appear as
// many short overlapping entries instead of one correctly-timed span.
func mergeZeroTempBasals(events []treatments.Treatment) []treatments.Treatment {
	var merged []treatments.Treatment
	var zeroBasals = make([]treatments.TempBasalEvent, 0, 30)
	var pendingOthers = make([]treatments.Treatment, 0, 10)

	flush := func() {
		first := zeroBasals[0]
		last := zeroBasals[len(zeroBasals)-1]
		lastExpiry := last.TreatmentTime().Add(time.Duration(last.Duration()) * time.Minute)

		// Since we check for natural expiration we are aware of all zero temp basals that happen in a consecutive window
		duration := int(lastExpiry.Sub(first.TreatmentTime()).Minutes())

		merged = append(merged, treatments.NewTempBasalEvent(first.TreatmentTime(), first.TreatmentEnteredBy(), duration, 0))
		merged = append(merged, pendingOthers...)
		zeroBasals = zeroBasals[:0]
		pendingOthers = pendingOthers[:0]
	}

	for _, e := range events {
		tb, ok := e.(treatments.TempBasalEvent)

		if ok && tb.Rate() == 0 {
			// Zero-rate temp basal
			if len(zeroBasals) > 0 {
				last := zeroBasals[len(zeroBasals)-1]

				// Check if it naturally expired (actually took 30/60/90/120minutes)
				// The current temp basal event takes place after the last temp basal event has fully taken place
				lastExpiry := last.TreatmentTime().Add(time.Duration(last.Duration()) * time.Minute)
				if !tb.TreatmentTime().Before(lastExpiry) {
					flush()
				}
			}

			zeroBasals = append(zeroBasals, tb)
			continue
		}

		if ok {
			// Non-zero temp basal: true basal reset
			if len(zeroBasals) > 0 {
				flush()
			}

			// Note: the real duration of the non-zero temp basal is not properly recorded here.
			// it will be corrected afterward.
			merged = append(merged, e)
			continue
		}

		// Flush criterion: non-temp basal event
		if len(zeroBasals) > 0 {
			last := zeroBasals[len(zeroBasals)-1]

			// Again: check for natural expiration
			lastExpiry := last.TreatmentTime().Add(time.Duration(last.Duration()) * time.Minute)
			if e.TreatmentTime().Before(lastExpiry) {
				// Store the current event such that we can add it after merging
				// Such that we get the full zero basal followed by all the (non-zero temp basal) events after.
				pendingOthers = append(pendingOthers, e)
				continue
			}
			flush()
		}

		merged = append(merged, e)
	}

	// Flush remaining group
	if len(zeroBasals) > 0 {
		flush()
	}

	return merged
}

// mergeSMBs collapses consecutive SMB events into a single aggregated entry.
// A new group starts whenever the gap between two successive SMBs exceeds 10
// minutes (roughly two AID loop cycles), indicating the loop genuinely paused.
func mergeSMBs(events []treatments.Treatment) []treatments.Treatment {
	// We assume that the time between last two SMBs is at most within 10 minutes (about 2 cycles) of each other.
	const cycleGap = 2 * 5 * time.Minute

	var merged []treatments.Treatment
	var pending []treatments.SmbTreatment

	flush := func() {
		if len(pending) == 0 {
			return
		}

		first := pending[0]

		if len(pending) == 1 {
			merged = append(merged, first)
		} else {
			var total float32

			// Sum up SMB amount
			for _, s := range pending {
				total += s.Insulin()
			}

			span := int(pending[len(pending)-1].TreatmentTime().Sub(first.TreatmentTime()).Minutes())
			merged = append(merged, treatments.NewAggregatedSmbTreatment(first.TreatmentTime(), total, first.TreatmentEnteredBy(), span))
		}

		pending = pending[:0]
	}

	for _, e := range events {
		smb, ok := e.(treatments.SmbTreatment)
		if !ok {
			flush()
			merged = append(merged, e)
			continue
		}

		// If the gap between current and last pending is larger than the cycle gap, flush.
		if len(pending) > 0 && smb.TreatmentTime().Sub(pending[len(pending)-1].TreatmentTime()) >= cycleGap {
			flush()
		}

		// Queue up SMB
		pending = append(pending, smb)
	}

	// Flush remaining group
	flush()

	return merged
}

// fixTempBasalDurations corrects each temp basal's duration to the time until
// the next temp basal begins, if that is shorter than the nominal duration.
// AID systems encode a fixed duration (e.g. 30 min) on every entry even when
// the basal is replaced within minutes.
func fixTempBasalDurations(events []treatments.Treatment) []treatments.Treatment {
	for i, e := range events {
		tb, ok := e.(treatments.TempBasalEvent)
		if !ok {
			continue
		}

		expiry := tb.TreatmentTime().Add(time.Duration(tb.Duration()) * time.Minute)

		for j := i + 1; j < len(events); j++ {
			next, ok := events[j].(treatments.TempBasalEvent)
			if !ok {
				continue
			}
			if next.TreatmentTime().Before(expiry) {
				newDuration := int(next.TreatmentTime().Sub(tb.TreatmentTime()).Minutes())
				events[i] = treatments.NewTempBasalEvent(tb.TreatmentTime(), tb.TreatmentEnteredBy(), newDuration, tb.Rate())
			}
			break
		}
	}

	return events
}
