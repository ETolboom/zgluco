package treatments

import "time"

type TreatmentType int

const (
	SMB = iota // Also called "Correction Bolus" in Loop
	TempBasal
	CarbCorrection
	TemporaryOverride // Loop e.g. "Temporary Override 🏋️‍♂️ Lifting 130% 7 - 8"
	TemporaryTarget   // Trio "Custom Target"
	Bolus
	InsulinChange
	Note
	ExternalInsulin
)

type Treatment interface {
	TreatmentKind() TreatmentType
	TreatmentTime() time.Time
	TreatmentEnteredBy() string
	String() string
}
