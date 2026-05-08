package sources

import (
	"zgluco/internal/models"
	"zgluco/internal/models/profile"
	"zgluco/internal/models/treatments"
	"zgluco/internal/types"
)

type Source interface {
	FetchProfile(timeRange types.TimeRange) (*profile.Profile, error)
	FetchTreatments(timeRange types.TimeRange) ([]treatments.Treatment, error)
	FetchSensorGlucoseValues(timeRange types.TimeRange) ([]*models.SensorGlucoseValue, error)
}
