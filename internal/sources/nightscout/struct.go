package nightscout

import (
	"time"
)

// Profile represents a Nightscout profile document
type Profile struct {
	DefaultProfile string                  `json:"defaultProfile,omitempty"`
	Store          map[string]*ProfileData `json:"store"`
	EnteredBy      string                  `json:"enteredBy,omitempty"` // Trio / Loop
	StartDate      *time.Time              `json:"startDate,omitempty"`
}

// GetDefaultProfile returns the profile currently in use (default) by the user
func (p *Profile) GetDefaultProfile() *ProfileData {
	return p.Store[p.DefaultProfile]
}

// ProfileData represents the settings within a profile
type ProfileData struct {
	Units      string           `json:"units,omitempty"`
	DIA        float32          `json:"dia,omitempty"`
	TargetLow  []TimeValueEntry `json:"target_low,omitempty"`
	TargetHigh []TimeValueEntry `json:"target_high,omitempty"`
	CarbRatio  []TimeValueEntry `json:"carbratio,omitempty"`
	Sens       []TimeValueEntry `json:"sens,omitempty"`
	Basal      []TimeValueEntry `json:"basal,omitempty"`
}

// TimeValueEntry represents a time-indexed setting (e.g., basal rate, sensitivity)
type TimeValueEntry struct {
	Time          string  `json:"time"`
	TimeAsSeconds int     `json:"timeAsSeconds"`
	Value         float32 `json:"value"`
}

type Treatment struct {
	EnteredBy               string    `json:"enteredBy"`
	CreatedAt               time.Time `json:"created_at"`
	Duration                float32   `json:"duration,omitempty"`
	Absolute                float32   `json:"absolute,omitempty"`
	Rate                    float32   `json:"rate,omitempty"`
	EventType               string    `json:"eventType"`
	UtcOffset               float32   `json:"utcOffset"`
	Carbs                   float32   `json:"carbs"`
	Insulin                 float32   `json:"insulin"`
	Protein                 float32   `json:"protein,omitempty"`
	FoodType                string    `json:"foodType,omitempty"`
	Fat                     float32   `json:"fat,omitempty"`
	Notes                   string    `json:"notes,omitempty"`
	SyncIdentifier          string    `json:"syncIdentifier,omitempty"`
	Automatic               bool      `json:"automatic,omitempty"`
	Timestamp               time.Time `json:"timestamp,omitempty"`
	Temp                    string    `json:"temp,omitempty"`
	InsulinType             string    `json:"insulinType,omitempty"`
	Unabsorbed              float32   `json:"unabsorbed,omitempty"`
	Programmed              float64   `json:"programmed,omitempty"`
	Type                    string    `json:"type,omitempty"`
	UserEnteredAt           time.Time `json:"userEnteredAt,omitempty"`
	AbsorptionTime          float32   `json:"absorptionTime,omitempty"`
	Amount                  float32   `json:"amount,omitempty"`
	DurationType            string    `json:"durationType,omitempty"`
	Reason                  string    `json:"reason,omitempty"`
	InsulinNeedsScaleFactor float64   `json:"insulinNeedsScaleFactor,omitempty"`
	UserLastModifiedAt      time.Time `json:"userLastModifiedAt,omitempty"`
	TargetBottom            float32   `json:"targetBottom,omitempty"`
	TargetTop               float32   `json:"targetTop,omitempty"`
}

type Sgv struct {
	ID         string    `json:"_id"`
	Sgv        float64   `json:"sgv"`
	Direction  string    `json:"direction"`
	Type       string    `json:"type"`
	Filtered   float32   `json:"filtered"`
	Unfiltered float32   `json:"unfiltered"`
	Glucose    float32   `json:"glucose"`
	Date       float64   `json:"date"`
	DateString time.Time `json:"dateString"`
	UtcOffset  float32   `json:"utcOffset"`
	SysTime    time.Time `json:"sysTime"`
}
