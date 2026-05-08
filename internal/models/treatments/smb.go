package treatments

import (
	"fmt"
	"time"
)

type SmbTreatment struct {
	createdAt     time.Time
	insulin       float32
	enteredBy     string
	windowMinutes int
}

func (s SmbTreatment) TreatmentKind() TreatmentType {
	return SMB
}

func (s SmbTreatment) TreatmentTime() time.Time {
	return s.createdAt
}

func (s SmbTreatment) TreatmentEnteredBy() string {
	return s.enteredBy
}

func (s SmbTreatment) Insulin() float32 {
	return s.insulin
}

func NewSmbTreatment(createdAt time.Time, enteredBy string, insulin float32) SmbTreatment {
	return SmbTreatment{createdAt: createdAt, enteredBy: enteredBy, insulin: insulin}
}

func NewAggregatedSmbTreatment(createdAt time.Time, insulin float32, enteredBy string, windowMinutes int) SmbTreatment {
	return SmbTreatment{createdAt: createdAt, insulin: insulin, enteredBy: enteredBy, windowMinutes: windowMinutes}
}

func (s SmbTreatment) String() string {
	if s.windowMinutes > 0 {
		return fmt.Sprintf("%s SMB %.2fU (aggregated over %dmin)", s.createdAt.Format(time.RFC3339), s.insulin, s.windowMinutes)
	}
	return fmt.Sprintf("%s SMB %.2fU", s.createdAt.Format(time.RFC3339), s.insulin)
}
