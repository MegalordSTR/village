package economy

// WeatherPattern represents a weather pattern affecting agriculture.
type WeatherPattern string

const (
	WeatherDrought WeatherPattern = "drought"
	WeatherNormal  WeatherPattern = "normal"
	WeatherRain    WeatherPattern = "rain"
)

// YieldMultiplier returns the multiplier for agricultural yield for a given weather pattern.
// Based on AC: drought (0.5x yield), normal (1.0x), rain (1.3x yield but possible flooding)
func YieldMultiplier(wp WeatherPattern) float64 {
	switch wp {
	case WeatherDrought:
		return 0.5
	case WeatherNormal:
		return 1.0
	case WeatherRain:
		return 1.3
	default:
		return 1.0
	}
}

// FloodingProbability returns the probability of flooding causing crop loss for rainy weather.
// Based on AC: rain (1.3x yield but possible flooding). Assume 10% chance.
func FloodingProbability(wp WeatherPattern) float64 {
	if wp == WeatherRain {
		return 0.1
	}
	return 0.0
}

// ApplyFlooding simulates flooding effect on a crop, reducing yield.
// Returns the remaining yield multiplier after flooding.
func ApplyFlooding(yieldMultiplier float64) float64 {
	// Flooding reduces yield by 50%
	return yieldMultiplier * 0.5
}
