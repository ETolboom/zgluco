package types

import (
	"fmt"
	"math"
	"time"
)

type TimeRange struct {
	From time.Time
	To   time.Time
}

func NewTimeRange(lookbackDays int) TimeRange {
	now := time.Now()
	return TimeRange{
		From: now.AddDate(0, 0, -lookbackDays),
		To:   now,
	}
}

func (tr TimeRange) String() string {
	loc := time.Now().Location() // Use machine TZ
	from := tr.From.In(loc)
	to := tr.To.In(loc)

	days := int(math.Ceil(tr.To.Sub(tr.From).Hours() / 24))
	return fmt.Sprintf("%s - %s (%d days)\n\n", from.Format("Jan 02 2006"), to.Format("Jan 02 2006"), days)
}
