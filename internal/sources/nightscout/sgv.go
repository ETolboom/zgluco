package nightscout

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"zgluco/internal/models"
	"zgluco/internal/types"
)

func (n *Nightscout) fetchSensorGlucoseValues(timeRange types.TimeRange) ([]*models.SensorGlucoseValue, error) {
	res, err := n.doApiCall("entries.json", "dateString", "", timeRange)
	if err != nil {
		return nil, fmt.Errorf("could not fetch sensor glucose values: %w", err)
	}

	var sgvs []*Sgv
	err = json.Unmarshal(res, &sgvs)
	if err != nil {
		return nil, err
	}

	if len(sgvs) == 0 {
		return nil, errors.New("no sensor glucose values found")
	}

	var parsedSgvs []*models.SensorGlucoseValue
	for _, sgv := range sgvs {
		parsedSgvs = append(parsedSgvs, &models.SensorGlucoseValue{
			CreatedAt: sgv.DateString,
			Glucose:   sgv.Sgv,
			Direction: models.GlucoseDirectionFromString(sgv.Direction),
		})
	}

	slices.SortFunc(parsedSgvs, func(a, b *models.SensorGlucoseValue) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	return parsedSgvs, nil
}
