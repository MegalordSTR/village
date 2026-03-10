package economy

import (
	"math"
)

// IsValidType returns true if the resource type is one of the defined constants.
func IsValidType(rt ResourceType) bool {
	switch rt {
	case ResourceGrain, ResourceVegetables, ResourceWood, ResourceStone,
		ResourceIronOre, ResourceWool, ResourceFlour, ResourceBread,
		ResourcePlanks, ResourceIron, ResourceCloth, ResourceTools,
		ResourceFurniture, ResourceWeapons, ResourceClothing:
		return true
	default:
		return false
	}
}

// IsPerishable returns true if the resource spoils over time.
func IsPerishable(rt ResourceType) bool {
	return SpoilageRate(rt) > 0.0
}

// GetCategory returns the category of a resource type (alias for CategoryForResource).
func GetCategory(rt ResourceType) ResourceCategory {
	return CategoryForResource(rt)
}

// NewResource creates a new resource with default values.
// Default quality is Normal, location empty, produced date zero, value base value.
func NewResource(rt ResourceType, quantity float64) Resource {
	return Resource{
		Type:     rt,
		Quantity: quantity,
		Quality:  QualityNormal,
		Location: "",
		Produced: GameDate{},
		Value:    BaseValue(rt),
		Spoiled:  0.0,
	}
}

// Validate checks if the resource is valid (type valid, quantity non-negative finite number).
func (r Resource) Validate() bool {
	if !IsValidType(r.Type) {
		return false
	}
	if math.IsNaN(r.Quantity) || math.IsInf(r.Quantity, 0) {
		return false
	}
	if r.Quantity < 0 {
		return false
	}
	return true
}

// ComputeValue calculates the total value of the resource based on quantity, quality, and base value.
// Quality multiplier: Poor 0.8, Normal 1.0, Good 1.2, Excellent 1.5, Masterwork 2.0
func (r Resource) ComputeValue() float64 {
	multiplier := 1.0
	switch r.Quality {
	case QualityPoor:
		multiplier = 0.8
	case QualityNormal:
		multiplier = 1.0
	case QualityGood:
		multiplier = 1.2
	case QualityExcellent:
		multiplier = 1.5
	case QualityMasterwork:
		multiplier = 2.0
	}
	return r.Quantity * r.Value * multiplier
}

// QualityToFloat converts a QualityTier to a float64 representation (0.0 to 1.0).
func QualityToFloat(q QualityTier) float64 {
	return float64(q) / 4.0
}

// FloatToQuality converts a float64 to the nearest QualityTier.
// Clamps to valid range.
func FloatToQuality(f float64) QualityTier {
	scaled := f * 4.0
	rounded := int(math.Round(scaled))
	if rounded < 0 {
		rounded = 0
	}
	if rounded > 4 {
		rounded = 4
	}
	return QualityTier(rounded)
}
