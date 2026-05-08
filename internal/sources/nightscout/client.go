package nightscout

import (
	"fmt"
	"net/url"
	"zgluco/internal/models"
	"zgluco/internal/models/profile"
	"zgluco/internal/models/treatments"
	"zgluco/internal/types"
)

type Nightscout struct {
	nightscoutUrl string
	apiKey        string
}

func New(nightscoutUrl string, apiKey string) (*Nightscout, error) {
	u, err := url.Parse(nightscoutUrl)
	if err != nil {
		return nil, fmt.Errorf("could not parse nightscout url: %w", err)
	}

	u = u.JoinPath("/api/v1/")

	return &Nightscout{nightscoutUrl: u.String(), apiKey: apiKey}, nil
}

func (n *Nightscout) FetchProfile(timeRange types.TimeRange) (*profile.Profile, error) {
	return n.fetchProfile(timeRange)
}

func (n *Nightscout) FetchTreatments(timeRange types.TimeRange) ([]treatments.Treatment, error) {
	return n.fetchTreatments(timeRange)
}

func (n *Nightscout) FetchSensorGlucoseValues(timeRange types.TimeRange) (values []*models.SensorGlucoseValue, err error) {
	return n.fetchSensorGlucoseValues(timeRange)
}
