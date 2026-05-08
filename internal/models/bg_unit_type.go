package models

import (
	"fmt"
	"strings"
)

type BgUnitType int

const (
	Mmol BgUnitType = 0
	Mgdl BgUnitType = 1
)

func ParseBgUnitType(s string) BgUnitType {
	switch strings.ToLower(strings.ReplaceAll(s, "/l", "")) {
	case "mmol":
		return Mmol
	case "mgdl", "mg/dl":
		return Mgdl
	default:
		panic(fmt.Errorf("unknown bg unit type: %q", s))
	}
}

func (b BgUnitType) String() string {
	switch b {
	case Mmol:
		return "mmol/L"
	case Mgdl:
		return "mg/dL"
	default:
		panic(fmt.Errorf("unknown bg unit type: %d", b))
	}
}
