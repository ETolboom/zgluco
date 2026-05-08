package treatments

import (
	"fmt"
	"time"
)

type SiteChangeTreatment struct {
	enteredBy string
	createdAt time.Time
}

func (s SiteChangeTreatment) TreatmentKind() TreatmentType {
	return SiteChange
}

func (s SiteChangeTreatment) TreatmentTime() time.Time {
	return s.createdAt
}

func (s SiteChangeTreatment) TreatmentEnteredBy() string {
	return s.enteredBy
}

func NewSiteChangeTreatment(createdAt time.Time, enteredBy string) SiteChangeTreatment {
	return SiteChangeTreatment{createdAt: createdAt, enteredBy: enteredBy}
}

func (s SiteChangeTreatment) String() string {
	return fmt.Sprintf("%s Site Change", s.createdAt.Format(time.RFC3339))
}
