package formatters

import (
	"fmt"
	"math"
	"strings"
	"time"

	"zgluco/internal/models"
	"zgluco/internal/models/profile"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type glucoseBucket struct {
	Time             time.Time
	Avg              float64
	Min              float64
	Max              float64
	AverageDirection models.GlucoseDirection
	Count            int
}

func FormatSGVs(sb *strings.Builder, p *profile.Profile, sgvs []*models.SensorGlucoseValue, loc *time.Location) error {
	fmt.Fprintf(sb, "═══ Glucose Data (%s) ═══\n\n", p.PreferredUnits)

	buckets := bucketGlucose(sgvs)
	if len(buckets) == 0 {
		fmt.Fprintln(sb, "No glucose data available.")
		return nil
	}

	table := tablewriter.NewTable(sb,
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{
					PerColumn: []tw.Align{
						tw.AlignLeft,
						tw.AlignLeft,
						tw.AlignRight,
						tw.AlignRight,
						tw.AlignRight,
						tw.AlignLeft,
					},
				},
			},
		}),
	)
	defer table.Close()

	table.Header("Date", "Time", "Avg", "Min", "Max", "Trend")

	var lastDay string
	for i, b := range buckets {
		localTime := b.Time.In(loc)
		day := localTime.Format("2006-01-02")
		timeStr := localTime.Format("15:04")

		dayDisplay := day
		if day == lastDay {
			dayDisplay = ""
		}
		lastDay = day

		avg, _min, _max := b.Avg, b.Min, b.Max
		var avgStr, minStr, maxStr string
		if p.PreferredUnits == models.Mmol {
			avg, _min, _max = avg/18, _min/18, _max/18
			avgStr = fmt.Sprintf("%.1f", avg)
			minStr = fmt.Sprintf("%.1f", _min)
			maxStr = fmt.Sprintf("%.1f", _max)
		} else {
			avgStr = fmt.Sprintf("%.0f", avg)
			minStr = fmt.Sprintf("%.0f", _min)
			maxStr = fmt.Sprintf("%.0f", _max)
		}

		table.Append([]string{dayDisplay, timeStr, avgStr, minStr, maxStr, buckets[i].AverageDirection.String()})
	}

	return table.Render()
}

func bucketGlucose(sgvs []*models.SensorGlucoseValue) []glucoseBucket {
	const bucketDuration = 15 * time.Minute
	if len(sgvs) == 0 {
		return nil
	}

	first := sgvs[0].CreatedAt.Truncate(bucketDuration)
	last := sgvs[len(sgvs)-1].CreatedAt

	var buckets []glucoseBucket
	for t := first; !t.After(last); t = t.Add(bucketDuration) {
		bucketEnd := t.Add(bucketDuration)
		var (
			sum            float64
			lowest         = math.MaxFloat64
			highest        float64
			count          int
			directions     = make(map[models.GlucoseDirection]int)
			directionMode  models.GlucoseDirection
			directionCount int
		)

		for _, sgv := range sgvs {
			if sgv.CreatedAt.Before(t) || !sgv.CreatedAt.Before(bucketEnd) {
				continue
			}
			val := sgv.Glucose
			sum += val
			count++
			if val < lowest {
				lowest = val
			}
			if val > highest {
				highest = val
			}
			directions[sgv.Direction]++
			if directions[sgv.Direction] > directionCount {
				directionMode = sgv.Direction
				directionCount = directions[sgv.Direction]
			}
		}

		if count == 0 {
			continue
		}

		buckets = append(buckets, glucoseBucket{
			Time:             t,
			Avg:              sum / float64(count),
			Min:              lowest,
			Max:              highest,
			Count:            count,
			AverageDirection: directionMode,
		})
	}

	return buckets
}
