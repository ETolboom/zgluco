package treatments

import (
	"fmt"
	"time"
)

type ExternalInsulinEvent struct {
	createdAt time.Time
	enteredBy string
	insulin   float32
}

func (e ExternalInsulinEvent) TreatmentKind() TreatmentType {
	return ExternalInsulin
}

func (e ExternalInsulinEvent) TreatmentTime() time.Time {
	return e.createdAt
}

func (e ExternalInsulinEvent) TreatmentEnteredBy() string {
	return e.enteredBy
}

func (e ExternalInsulinEvent) String() string {
	return fmt.Sprintf("%s External Insulin %.2fU", e.createdAt.Format(time.RFC3339), e.insulin)
}

func NewExternalInsulinEvent(createdAt time.Time, enteredBy string, insulin float32) ExternalInsulinEvent {
	return ExternalInsulinEvent{createdAt: createdAt, enteredBy: enteredBy, insulin: insulin}
}
