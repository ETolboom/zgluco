package treatments

import (
	"fmt"
	"strings"
	"time"
	"zgluco/internal/models"
)

type CarbCorrectionEvent struct {
	enteredBy   string
	createdAt   time.Time
	insulin     float32
	insulinType string
	meal        models.Meal
}

func (c CarbCorrectionEvent) TreatmentKind() TreatmentType {
	return CarbCorrection
}

func (c CarbCorrectionEvent) TreatmentTime() time.Time {
	return c.createdAt
}

func (c CarbCorrectionEvent) TreatmentEnteredBy() string {
	return c.enteredBy
}

func (c CarbCorrectionEvent) Insulin() float32 {
	return c.insulin
}

func (c CarbCorrectionEvent) InsulinType() string {
	return c.insulinType
}

func (c CarbCorrectionEvent) Meal() models.Meal {
	return c.meal
}

func NewCarbCorrectionEvent(createdAt time.Time, enteredBy string, insulin float32, insulinType string, meal models.Meal) CarbCorrectionEvent {
	return CarbCorrectionEvent{createdAt: createdAt, enteredBy: enteredBy, insulin: insulin, insulinType: insulinType, meal: meal}
}

func (c CarbCorrectionEvent) String() string {
	var sb strings.Builder
	sb.WriteString(c.createdAt.Format(time.RFC3339))
	sb.WriteString(" Meal")

	if c.insulin > 0 {
		fmt.Fprintf(&sb, " %.2fU insulin", c.insulin)
	}
	if c.meal.Carbs() > 0 {
		fmt.Fprintf(&sb, " %.0fg carbs", c.meal.Carbs())
	}
	if c.meal.Fat() > 0 {
		fmt.Fprintf(&sb, " %.0fg fat", c.meal.Fat())
	}
	if c.meal.Protein() > 0 {
		fmt.Fprintf(&sb, " %.0fg protein", c.meal.Protein())
	}
	if absorption := c.meal.AbsorptionTime(); absorption > 0 {
		fmt.Fprintf(&sb, " absorbed over %dmin", absorption)
	}

	return sb.String()
}
