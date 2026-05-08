package models

import "zgluco/internal/types"

type Meal struct {
	carbs          float32
	fat            float32
	protein        float32
	absorptionTime float32 // If provided, such as with Loop
}

func NewMeal(carbs, fat, protein float32, absorptionTime float32) Meal {
	return Meal{carbs: carbs, fat: fat, protein: protein, absorptionTime: absorptionTime}
}

func (m Meal) Carbs() float32 {
	return m.carbs
}

func (m Meal) Fat() float32 {
	return m.fat
}

func (m Meal) Protein() float32 {
	return m.protein
}

// AbsorptionTime returns the estimated absorption duration of a meal in minutes.
// If an explicit absorption time was provided (e.g. from Loop), it is returned directly.
// Otherwise, the duration is calculated from fat and protein using the Warsaw formula
// (as used in Trio), converting FPU-based hours to minutes.
func (m Meal) AbsorptionTime() int {
	if m.absorptionTime != 0 {
		// Return absorption time as provided by AID (e.g., Loop)
		return int(m.absorptionTime)
	}

	cfg := types.DefaultFPUConfig()

	kcal := m.protein*4 + m.fat*9
	carbEquivalents := (kcal / 10) * cfg.AdjustmentFactor
	fpus := carbEquivalents / 10

	duration := types.ComputeFPUDuration(fpus, cfg.TimeCap)

	return duration * 60
}
