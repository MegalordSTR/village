package economy

// CropStage represents the growth stage of a crop.
type CropStage int

const (
	CropPlanted CropStage = iota
	CropGrowing
	CropReadyForHarvest
	CropHarvested
)

// Crop represents a planted crop in a field.
type Crop struct {
	Type            ResourceType // e.g., ResourceGrain, ResourceVegetables
	Stage           CropStage
	PlantedWeek     int     // week number when planted (1-52)
	GrowthProgress  float64 // 0.0 to 1.0
	YieldMultiplier float64 // multiplier based on soil, weather, etc.
}

// NewCrop creates a new crop planted at the given week.
func NewCrop(cropType ResourceType, plantedWeek int) *Crop {
	return &Crop{
		Type:            cropType,
		Stage:           CropPlanted,
		PlantedWeek:     plantedWeek,
		GrowthProgress:  0.0,
		YieldMultiplier: 1.0,
	}
}

// UpdateGrowth advances the crop growth based on elapsed weeks, season, and weather.
// Returns true if the crop reached next stage.
func (c *Crop) UpdateGrowth(currentWeek int, season Season, weatherMultiplier float64) bool {
	// Determine if crop can grow (not winter)
	if season == Winter {
		// No growth in winter
		return false
	}
	// Compute growth per week based on season and weather
	growthPerWeek := 0.1 // base growth per week
	growthPerWeek *= AgricultureYieldMultiplier(season)
	growthPerWeek *= weatherMultiplier

	weeksPassed := currentWeek - c.PlantedWeek
	if weeksPassed < 0 {
		weeksPassed += 52 // wrap around year
	}
	// Ensure we only grow up to 1.0
	newProgress := float64(weeksPassed) * growthPerWeek
	if newProgress > 1.0 {
		newProgress = 1.0
	}
	c.GrowthProgress = newProgress

	// Update stage based on progress
	oldStage := c.Stage
	if c.GrowthProgress >= 1.0 {
		c.Stage = CropReadyForHarvest
	} else if c.GrowthProgress >= 0.5 {
		c.Stage = CropGrowing
	} else {
		c.Stage = CropPlanted
	}
	return c.Stage != oldStage
}

// Harvest returns the amount of resource harvested and marks crop as harvested.
// Returns 0 if not ready.
func (c *Crop) Harvest() float64 {
	if c.Stage != CropReadyForHarvest {
		return 0.0
	}
	c.Stage = CropHarvested
	// Base yield per crop type (could be configurable)
	baseYield := 10.0 // placeholder
	return baseYield * c.YieldMultiplier
}

// CanPlant returns true if planting is allowed in the given week.
// Planting window: spring weeks 1-4.
func CanPlant(week int) bool {
	season := SeasonFromWeek(week)
	if season != Spring {
		return false
	}
	// Spring weeks 1-4
	weekInSeason := week // since spring starts at week 1
	if weekInSeason >= 1 && weekInSeason <= 4 {
		return true
	}
	return false
}

// CanHarvest returns true if harvest is allowed in the given week.
// Harvest window: autumn weeks 8-12.
func CanHarvest(week int) bool {
	season := SeasonFromWeek(week)
	if season != Autumn {
		return false
	}
	// Autumn weeks 8-12 (autumn starts at week 27)
	weekInSeason := week - 26 // week 27 -> 1, week 28 ->2, etc.
	if weekInSeason >= 8 && weekInSeason <= 12 {
		return true
	}
	return false
}
