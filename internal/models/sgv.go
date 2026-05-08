package models

import (
	"time"
)

type GlucoseDirection int

const (
	RisingRapidly GlucoseDirection = iota
	RisingFast
	Flat
	FallingFast
	FallingRapidly
)

func (d GlucoseDirection) String() string {
	switch d {
	case RisingRapidly:
		return "rising rapidly"
	case RisingFast:
		return "rising fast"
	case Flat:
		return "flat"
	case FallingFast:
		return "falling fast"
	case FallingRapidly:
		return "falling rapidly"
	default:
		return "unknown"
	}
}

func GlucoseDirectionFromString(s string) GlucoseDirection {
	switch s {
	case "SingleUp":
		return RisingRapidly
	case "FortyFiveUp":
		return RisingFast
	case "FortyFiveDown":
		return FallingFast
	case "SingleDown":
		return FallingRapidly
	case "Flat":
		fallthrough
	default:
		return Flat
	}
}

type SensorGlucoseValue struct {
	CreatedAt time.Time        `json:"created_at"`
	Glucose   float64          `json:"glucose"`
	Direction GlucoseDirection `json:"direction"`
}
