package simulation

import (
	"math/rand"
	"fmt"
)

// EnvironmentSystem implements the environmental simulation system.
type EnvironmentSystem struct{}

// NewEnvironmentSystem creates a new environmental system.
func NewEnvironmentSystem() *EnvironmentSystem {
	return &EnvironmentSystem{}
}

// Update processes one week of environmental simulation.
// It updates weather, seasons, temperature, soil fertility, natural resources, and wildlife.
func (e *EnvironmentSystem) Update(week int, state *GameState, rng *rand.Rand) []Event {
	var events []Event

	// Update season based on week of year (1-52)
	season := calculateSeason(week)
	state.Environment.Season = season

	// Update temperature and rainfall based on season
	temp, rain := calculateWeather(season, rng)
	state.Environment.Temperature = temp
	state.Environment.Rainfall = rain

	// Update soil fertility (simplified: improves with moderate rain, degrades with extremes)
	state.Environment.SoilFertility = updateSoilFertility(state.Environment.SoilFertility, rain, rng)

	// Natural resource regeneration
	state.Environment.ForestHealth = regenerateResource(state.Environment.ForestHealth, rng, 0.01) // slow regeneration
	state.Environment.MineQuality = regenerateResource(state.Environment.MineQuality, rng, 0.001)  // very slow

	// Wildlife population fluctuations
	state.Environment.WildlifePopulation = updateWildlife(
		state.Environment.WildlifePopulation,
		state.Environment.ForestHealth,
		rng,
	)

	// Generate events for significant changes
	events = append(events, GenerateWeatherEvents(temp, rain, week, state.Calendar.Year)...)

	return events
}

// calculateSeason determines the season based on week number (1-52).
func calculateSeason(week int) string {
	// 52 weeks per year, 13 weeks per season
	// Week 1-13: Spring, 14-26: Summer, 27-39: Autumn, 40-52: Winter
	switch {
	case week <= 0:
		return "winter" // should not happen
	case week <= 13:
		return "spring"
	case week <= 26:
		return "summer"
	case week <= 39:
		return "autumn"
	default:
		return "winter"
	}
}

// calculateWeather returns temperature (°C) and rainfall (mm/week) for a given season.
func calculateWeather(season string, rng *rand.Rand) (float64, float64) {
	// Base values per season (temperate climate)
	var baseTemp, baseRain float64
	switch season {
	case "spring":
		baseTemp = 10.0
		baseRain = 15.0
	case "summer":
		baseTemp = 22.0
		baseRain = 8.0
	case "autumn":
		baseTemp = 12.0
		baseRain = 18.0
	case "winter":
		baseTemp = 2.0
		baseRain = 12.0
	default:
		baseTemp = 15.0
		baseRain = 10.0
	}

	// Add some randomness (±5°C, ±5mm)
	temp := baseTemp + (rng.Float64()*10 - 5)
	rain := baseRain + (rng.Float64()*10 - 5)
	if rain < 0 {
		rain = 0
	}

	return temp, rain
}

// updateSoilFertility updates soil fertility based on rainfall.
// Simplified model: moderate rain (10-20mm) improves fertility,
// extreme dryness or heavy rain reduces it.
func updateSoilFertility(current float64, rainfall float64, rng *rand.Rand) float64 {
	change := 0.0
	switch {
	case rainfall > 0 && rainfall < 5:
		// Very dry - fertility decreases
		change = -0.01
	case rainfall >= 5 && rainfall < 10:
		// Slightly dry - small decrease
		change = -0.005
	case rainfall >= 10 && rainfall <= 20:
		// Ideal - fertility improves
		change = 0.01
	case rainfall > 20 && rainfall <= 30:
		// Heavy rain - washes nutrients
		change = -0.005
	case rainfall > 30:
		// Very heavy - significant decrease
		change = -0.02
	}

	// Add small random variation
	change += (rng.Float64()*0.02 - 0.01)

	newFertility := current + change
	// Clamp between 0.0 and 1.0
	if newFertility < 0 {
		return 0.0
	}
	if newFertility > 1.0 {
		return 1.0
	}
	return newFertility
}

// regenerateResource slowly regenerates a natural resource.
func regenerateResource(current float64, rng *rand.Rand, rate float64) float64 {
	// Resources slowly regenerate on their own
	change := rate + (rng.Float64()*0.01 - 0.005)
	newValue := current + change
	if newValue < 0 {
		return 0.0
	}
	if newValue > 1.0 {
		return 1.0
	}
	return newValue
}

// updateWildlife updates wildlife population based on forest health.
func updateWildlife(current, forestHealth float64, rng *rand.Rand) float64 {
	// Wildlife depends on forest health
	baseChange := 0.005
	// If forest health is good (>0.7), population grows faster
	if forestHealth > 0.7 {
		baseChange = 0.01
	}
	// If forest health is poor (<0.3), population declines
	if forestHealth < 0.3 {
		baseChange = -0.01
	}

	// Add randomness
	change := baseChange + (rng.Float64()*0.02 - 0.01)
	newValue := current + change
	// Clamp
	if newValue < 0 {
		return 0.0
	}
	if newValue > 1.0 {
		return 1.0
	}
	return newValue
}

// generateEventID creates a deterministic event ID.
func generateEventID(base string, week int) string {
	// Simple deterministic ID
	return base + "-" + string(rune(week%26+97)) // append a letter a-z based on week
}

// formatWeekTimestamp creates a timestamp string for events.
func formatWeekTimestamp(year, week int) string {
	return fmt.Sprintf("%d-W%d", year, week)
}

// GenerateWeatherEvents returns weather events based on temperature and rainfall.
func GenerateWeatherEvents(temp, rain float64, week, year int) []Event {
	var events []Event
	if rain > 20.0 {
		events = append(events, Event{
			ID:        generateEventID("heavy-rain", week),
			Type:      "weather",
			Timestamp: formatWeekTimestamp(year, week),
			Data: map[string]interface{}{
				"rainfall": rain,
				"severity": "heavy",
			},
		})
	}
	if temp < 0.0 {
		events = append(events, Event{
			ID:        generateEventID("freezing-temp", week),
			Type:      "weather",
			Timestamp: formatWeekTimestamp(year, week),
			Data: map[string]interface{}{
				"temperature": temp,
				"severity":    "freezing",
			},
		})
	}
	return events
}
