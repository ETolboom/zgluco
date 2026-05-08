package treatments

import (
	"fmt"
	"strings"
	"time"
	"zgluco/internal/models"
)

type TemporaryTargetEvent struct {
	enteredBy  string
	createdAt  time.Time
	reason     string
	duration   string
	targetLow  float32
	targetHigh float32
}

func NewTemporaryTargetEvent(createdAt time.Time, enteredBy string, reason string, duration string, targetLow float32, targetHigh float32) TemporaryTargetEvent {
	return TemporaryTargetEvent{createdAt: createdAt, enteredBy: enteredBy, reason: reason, duration: duration, targetLow: targetLow, targetHigh: targetHigh}
}

func (t TemporaryTargetEvent) TreatmentKind() TreatmentType {
	return TemporaryTarget
}

func (t TemporaryTargetEvent) TreatmentTime() time.Time {
	return t.createdAt
}

func (t TemporaryTargetEvent) TreatmentEnteredBy() string {
	return t.enteredBy
}

func (t TemporaryTargetEvent) Reason() string {
	return t.reason
}

// TargetLow returns the lower bound of the target range in mg/dL
func (t TemporaryTargetEvent) TargetLow() float32 {
	return t.targetLow
}

// TargetHigh returns the upper bound of the target range in mg/dL
func (t TemporaryTargetEvent) TargetHigh() float32 {
	return t.targetHigh
}

func (t TemporaryTargetEvent) String() string {
	// By default, Trio provides targetTop and targetBottom in mg/dL.
	return t.StringWithUnit(models.Mgdl)
}

func (t TemporaryTargetEvent) StringWithUnit(unit models.BgUnitType) string {
	var sb strings.Builder

	sb.WriteString(t.createdAt.Format(time.RFC3339))
	sb.WriteString(" Temporary Target")

	low, high := float64(t.targetLow), float64(t.targetHigh)
	var valStr string
	if unit == models.Mmol {
		low, high = low/18, high/18
		if high == low {
			valStr = fmt.Sprintf("%.1f", high)
		} else {
			valStr = fmt.Sprintf("%.1f-%.1f", low, high)
		}
	} else {
		if high == low {
			valStr = fmt.Sprintf("%d", int(high))
		} else {
			valStr = fmt.Sprintf("%d-%d", int(low), int(high))
		}
	}

	fmt.Fprintf(&sb, " %s %s", valStr, unit)

	if t.reason != "" {
		fmt.Fprintf(&sb, " (%s)", t.reason)
	}

	fmt.Fprintf(&sb, " %s minutes", t.duration)

	return sb.String()
}
