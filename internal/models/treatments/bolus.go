package treatments

import (
	"fmt"
	"time"
)

type BolusEvent struct {
	enteredBy string
	createdAt time.Time
	insulin   float32
}

func (b BolusEvent) TreatmentKind() TreatmentType {
	return Bolus
}

func (b BolusEvent) TreatmentTime() time.Time {
	return b.createdAt
}

func (b BolusEvent) TreatmentEnteredBy() string {
	return b.enteredBy
}

func (b BolusEvent) String() string {
	return fmt.Sprintf("%s Manual Bolus %0.2fU", b.createdAt.Format(time.RFC3339), b.insulin)
}

func NewBolusEvent(createdAt time.Time, enteredBy string, insulin float32) BolusEvent {
	return BolusEvent{createdAt: createdAt, enteredBy: enteredBy, insulin: insulin}
}
