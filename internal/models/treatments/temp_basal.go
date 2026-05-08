package treatments

import (
	"fmt"
	"time"
)

type TempBasalEvent struct {
	createdAt time.Time
	enteredBy string
	duration  int     // In minutes
	rate      float32 // Rate of increase/decrease
}

func (t TempBasalEvent) TreatmentKind() TreatmentType {
	return TempBasal
}

func (t TempBasalEvent) TreatmentTime() time.Time {
	return t.createdAt
}

func (t TempBasalEvent) TreatmentEnteredBy() string {
	return t.enteredBy
}

// Duration returns the duration of the temp basal event in minutes.
func (t TempBasalEvent) Duration() int {
	return t.duration
}

// Rate returns the rate which the basal rate is scaled by.
func (t TempBasalEvent) Rate() float32 {
	return t.rate
}

func NewTempBasalEvent(createdAt time.Time, enteredBy string, duration int, rate float32) TempBasalEvent {
	return TempBasalEvent{createdAt: createdAt, enteredBy: enteredBy, duration: duration, rate: rate}
}

func (t TempBasalEvent) String() string {
	return fmt.Sprintf("%s Temp Basal %.2fU/hr for %dmin", t.createdAt.Format(time.RFC3339), t.rate, t.duration)
}
