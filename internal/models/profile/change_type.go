package profile

type ChangeType int

const (
	Basal ChangeType = iota
	CarbRatio
	ISF
	TargetHigh
	TargetLow
	DIA
)

func (ct ChangeType) String() string {
	switch ct {
	case Basal:
		return "Basal"
	case CarbRatio:
		return "Carb Ratio"
	case ISF:
		return "ISF"
	case TargetHigh:
		return "Target High"
	case TargetLow:
		return "Target Low"
	case DIA:
		return "DIA"
	default:
		return "Unknown"
	}
}
