package economy

// String returns the human-readable name of the quality tier.
func (q QualityTier) String() string {
	switch q {
	case QualityPoor:
		return "Poor"
	case QualityNormal:
		return "Normal"
	case QualityGood:
		return "Good"
	case QualityExcellent:
		return "Excellent"
	case QualityMasterwork:
		return "Masterwork"
	default:
		return "Unknown"
	}
}

// DurabilityMultiplier returns the multiplier for tool durability based on quality.
// Poor: 0.5x, Normal: 1.0x, Good: 1.5x, Excellent: 2.0x, Masterwork: 3.0x
func DurabilityMultiplier(q QualityTier) float64 {
	switch q {
	case QualityPoor:
		return 0.5
	case QualityNormal:
		return 1.0
	case QualityGood:
		return 1.5
	case QualityExcellent:
		return 2.0
	case QualityMasterwork:
		return 3.0
	default:
		return 1.0
	}
}

// ProductionSpeedMultiplier returns the multiplier for production speed based on quality.
// Higher quality tools and materials speed up production.
// Poor: 0.8x, Normal: 1.0x, Good: 1.1x, Excellent: 1.2x, Masterwork: 1.3x
func ProductionSpeedMultiplier(q QualityTier) float64 {
	switch q {
	case QualityPoor:
		return 0.8
	case QualityNormal:
		return 1.0
	case QualityGood:
		return 1.1
	case QualityExcellent:
		return 1.2
	case QualityMasterwork:
		return 1.3
	default:
		return 1.0
	}
}

// SpoilageResistanceMultiplier returns the multiplier for spoilage rate based on quality.
// Higher quality reduces spoilage rate (lower multiplier).
// Poor: 1.0x (no reduction), Normal: 1.0x, Good: 0.7x, Excellent: 0.5x, Masterwork: 0.3x
func SpoilageResistanceMultiplier(q QualityTier) float64 {
	switch q {
	case QualityPoor:
		return 1.0
	case QualityNormal:
		return 1.0
	case QualityGood:
		return 0.7
	case QualityExcellent:
		return 0.5
	case QualityMasterwork:
		return 0.3
	default:
		return 1.0
	}
}

// MaterialQualityBonus returns a bonus to quality for special materials.
// Fine iron provides +0.2, seasoned wood provides +0.1, other materials have no bonus.
func MaterialQualityBonus(rt ResourceType) float64 {
	switch rt {
	case ResourceIron:
		return 0.2
	case ResourceWood:
		return 0.1
	default:
		return 0.0
	}
}
