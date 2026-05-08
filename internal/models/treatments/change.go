package treatments

import (
	"fmt"
	"time"
)

type ChangeType int

const (
	Insulin ChangeType = iota
	Site
)

type ChangeEvent struct {
	createdAt  time.Time
	enteredBy  string
	changeType ChangeType
}

func (c ChangeEvent) TreatmentKind() TreatmentType {
	return InsulinChange
}

func (c ChangeEvent) TreatmentTime() time.Time {
	return c.createdAt
}

func (c ChangeEvent) TreatmentEnteredBy() string {
	return c.enteredBy
}

func (c ChangeEvent) String() string {
	switch c.changeType {
	case Insulin:
		return fmt.Sprintf("%s Insulin Change", c.createdAt.Format(time.RFC3339))
	case Site:
		return fmt.Sprintf("%s Site Change", c.createdAt.Format(time.RFC3339))
	default:
		panic("invalid changeType")
	}
}

func NewChangeEvent(createdAt time.Time, enteredBy string, changeType ChangeType) ChangeEvent {
	return ChangeEvent{createdAt, enteredBy, changeType}
}
