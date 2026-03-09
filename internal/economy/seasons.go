package economy

// Season represents a season of the year.
type Season string

const (
	Spring Season = "spring"
	Summer Season = "summer"
	Autumn Season = "autumn"
	Winter Season = "winter"
)

// SeasonFromWeek returns the season for a given week number (1-52).
// Week 1-13: Spring, 14-26: Summer, 27-39: Autumn, 40-52: Winter.
func SeasonFromWeek(week int) Season {
	if week <= 0 {
		return Winter // should not happen
	}
	if week <= 13 {
		return Spring
	}
	if week <= 26 {
		return Summer
	}
	if week <= 39 {
		return Autumn
	}
	return Winter
}

// AgricultureYieldMultiplier returns the multiplier for agricultural yield in a given season.
// Based on AC: 0.0 in winter, 1.5 in summer, 1.0 in spring and autumn.
func AgricultureYieldMultiplier(season Season) float64 {
	switch season {
	case Summer:
		return 1.5
	case Spring, Autumn:
		return 1.0
	case Winter:
		return 0.0
	default:
		return 1.0
	}
}

// FuelConsumptionMultiplier returns the multiplier for fuel consumption in a given season.
// Based on AC: 1.0 summer, 2.5 winter, 1.0 other seasons.
func FuelConsumptionMultiplier(season Season) float64 {
	switch season {
	case Winter:
		return 2.5
	case Summer:
		return 1.0
	case Spring, Autumn:
		return 1.0
	default:
		return 1.0
	}
}

// PriceMultiplier returns the multiplier for market prices of food resources based on season and week.
// Food cheap after harvest (autumn weeks 8-12): 0.8, expensive in spring: 1.2, otherwise normal: 1.0.
func PriceMultiplier(resourceType ResourceType, season Season, week int) float64 {
	// Only food resources affected
	isFood := false
	switch resourceType {
	case ResourceGrain, ResourceVegetables, ResourceFlour, ResourceBread:
		isFood = true
	}
	if !isFood {
		return 1.0
	}
	// After harvest (autumn weeks 8-12)
	if season == Autumn {
		weekInAutumn := week - 26 // autumn starts at week 27
		if weekInAutumn >= 8 && weekInAutumn <= 12 {
			return 0.8
		}
	}
	// Expensive in spring
	if season == Spring {
		return 1.2
	}
	return 1.0
}

// RecommendedStockpiles returns recommended extra stockpile quantities for winter preparation.
// For winter season, recommends extra food and fuel.
func RecommendedStockpiles(season Season) map[ResourceType]float64 {
	if season != Winter {
		return nil
	}
	return map[ResourceType]float64{
		ResourceGrain:      100.0,
		ResourceVegetables: 50.0,
		ResourceWood:       200.0,
	}
}
