package profile

import (
	"time"
	"zgluco/internal/models"
)

type TimeValueDiff struct {
	TimeAsSeconds int
	OldValue      float32
	NewValue      float32
	Removed       bool
}

type Change struct {
	Type        ChangeType
	EffectiveAt *time.Time
	Diffs       []TimeValueDiff
}

type TimeValue struct {
	StartsAtSeconds int
	Value           float32
}

type GlucoseTarget struct {
	StartsAtSeconds int
	LowerBound      float32
	UpperBound      float32
}

type Profile struct {
	PreferredUnits            models.BgUnitType
	BasalRates                []*TimeValue
	CarbRatios                []*TimeValue
	InsulinSensitivityFactors []*TimeValue
	GlucoseTargets            []*GlucoseTarget
	Changes                   []*Change
}

func (p Profile) GetTotalBasalDose() float32 {
	var totalBasalDose float32
	for _, rate := range p.BasalRates {
		totalBasalDose += rate.Value
	}

	return totalBasalDose
}
