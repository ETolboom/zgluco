package types

import (
	"math"
	"time"
)

// FPUConfig holds the configurable parameters for Fat-Protein Unit processing,
// matching the Warsaw formula implementation used in Trio.
type FPUConfig struct {
	// Delay is the time before the first FPU carb equivalent entry (default: 60 min).
	Delay time.Duration
	// Interval is the time between each carb equivalent entry (default: 30 min).
	Interval time.Duration
	// TimeCap is the maximum duration for FPU absorption in hours (default: 8).
	TimeCap int
	// AdjustmentFactor scales the carb equivalents (default: 0.5 = half effect).
	AdjustmentFactor float32
}

// DefaultFPUConfig returns the default FPU configuration matching Trio defaults.
func DefaultFPUConfig() FPUConfig {
	return FPUConfig{
		Delay:            60 * time.Minute,
		Interval:         30 * time.Minute,
		TimeCap:          8,
		AdjustmentFactor: 0.5,
	}
}

// FPUResult holds the output of the Warsaw formula FPU calculation.
type FPUResult struct {
	// CarbEquivalents is the total grams of carb equivalents from fat and protein.
	CarbEquivalents float32

	// FPUs is the number of Fat-Protein Units (carbEquivalents / 10).
	FPUs float32

	// Duration is the computed absorption duration in hours.
	Duration int

	// EntrySize is the grams of carb equivalent per entry.
	EntrySize float32

	// Entries are the scheduled carb equivalent entries over time.
	Entries []FPUEntry
}

// FPUEntry represents a single scheduled carb equivalent dose.
type FPUEntry struct {
	Time  time.Time
	Carbs float32
}

// FatProteinUnits calculates Fat-Protein Units using the Warsaw formula (as in Trio)
// and distributes carb equivalents over time.
//
// Warsaw formula:
//  1. kcal = protein * 4 + fat * 9
//  2. carbEquivalents = (kcal / 10) * adjustmentFactor
//  3. fpus = carbEquivalents / 10
func FatProteinUnits(protein, fat float32, mealTime time.Time) FPUResult {
	cfg := DefaultFPUConfig()

	kcal := protein*4 + fat*9
	carbEquivalents := (kcal / 10) * cfg.AdjustmentFactor
	fpus := carbEquivalents / 10

	duration := ComputeFPUDuration(fpus, cfg.TimeCap)

	intervalHours := float32(cfg.Interval.Minutes()) / 60
	entrySize := carbEquivalents / float32(duration) * intervalHours

	if entrySize < 1.0 {
		entrySize = 1.0
	}

	// Round to 1 decimal place
	entrySize = float32(math.Round(float64(entrySize)*10)) / 10

	numberOfEntries := int(carbEquivalents / entrySize)

	entries := make([]FPUEntry, 0, numberOfEntries)
	t := mealTime.Add(cfg.Delay)
	for i := 0; i < numberOfEntries; i++ {
		entries = append(entries, FPUEntry{
			Time:  t,
			Carbs: entrySize,
		})
		t = t.Add(cfg.Interval)
	}

	return FPUResult{
		CarbEquivalents: carbEquivalents,
		FPUs:            fpus,
		Duration:        duration,
		EntrySize:       entrySize,
		Entries:         entries,
	}
}

// ComputeFPUDuration maps FPU count to absorption duration in hours, capped by timeCap.
func ComputeFPUDuration(fpus float32, timeCap int) int {
	var hours int
	switch {
	case fpus >= 4:
		hours = 8
	case fpus >= 3:
		hours = 7
	case fpus >= 2:
		hours = 6
	case fpus >= 1:
		hours = 4
	default:
		hours = 3
	}
	if hours > timeCap {
		hours = timeCap
	}
	return hours
}
