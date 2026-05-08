package treatments

import (
	"fmt"
	"strings"
	"time"
)

type TemporaryOverrideEvent struct {
	enteredBy   string
	createdAt   time.Time
	reason      string
	scaleFactor float32
	duration    float32 // -1 is indefinite
}

func (t TemporaryOverrideEvent) TreatmentKind() TreatmentType {
	return TemporaryOverride
}

func (t TemporaryOverrideEvent) TreatmentTime() time.Time {
	return t.createdAt
}

func (t TemporaryOverrideEvent) TreatmentEnteredBy() string {
	return t.enteredBy
}

// ScaleFactor returns the factor with which the insulin needs are scaled with.
func (t TemporaryOverrideEvent) ScaleFactor() float32 {
	return t.scaleFactor
}

// Duration returns the override duration in minutes
func (t TemporaryOverrideEvent) Duration() float32 {
	return t.duration
}

func NewTemporaryOverrideEvent(createdAt time.Time, enteredBy string, reason string, scaleFactor float32, duration float32) TemporaryOverrideEvent {
	return TemporaryOverrideEvent{createdAt: createdAt, enteredBy: enteredBy, reason: reason, scaleFactor: scaleFactor, duration: duration}
}

func (t TemporaryOverrideEvent) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s Temporary Override %.0f%% for %.0fmin", t.createdAt.Format(time.RFC3339), t.scaleFactor*100, t.duration)

	if t.reason != "" {
		fmt.Fprintf(&sb, " (%s)", t.reason)
	}

	return sb.String()
}
