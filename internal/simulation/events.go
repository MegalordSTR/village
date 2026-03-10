package simulation

import (
	"github.com/vano44/village/internal/economy"
	"math/rand"
)

// EventSystem implements the random event simulation system.
type EventSystem struct{}

// eventDef represents a possible random event.
type eventDef struct {
	eventType string
	trigger   func(state *GameState, week int, rng *rand.Rand) bool
	weight    func(state *GameState, week int) float64
	apply     func(state *GameState, week int, rng *rand.Rand) Event
}

// eventDefinitions holds all possible random events.
var eventDefinitions = []eventDef{
	{
		eventType: "festival",
		trigger: func(state *GameState, week int, rng *rand.Rand) bool {
			// Festival happens every 10 weeks
			return week%10 == 0
		},
		weight: func(state *GameState, week int) float64 {
			// Base weight, can be modified by happiness etc.
			return 1.0
		},
		apply: func(state *GameState, week int, rng *rand.Rand) Event {
			return Event{
				ID:        generateEventID("festival", week),
				Type:      "festival",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"week": week,
					"mood": "celebratory",
				},
			}
		},
	},
	{
		eventType: "disease_outbreak",
		trigger: func(state *GameState, week int, rng *rand.Rand) bool {
			// Disease outbreaks more likely in winter and poor forest health
			return state.Environment.Season == "winter" && state.Environment.ForestHealth < 0.3
		},
		weight: func(state *GameState, week int) float64 {
			base := 0.5
			// Increase weight with lower forest health
			return base * (1.0 - state.Environment.ForestHealth)
		},
		apply: func(state *GameState, week int, rng *rand.Rand) Event {
			// Reduce food resource by 10%
			adjustResource(state, economy.ResourceGrain, -0.1)
			return Event{
				ID:        generateEventID("disease", week),
				Type:      "disease_outbreak",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"severity":  "moderate",
					"affected":  "residents",
					"food_loss": 0.1,
				},
			}
		},
	},
	{
		eventType: "good_harvest",
		trigger: func(state *GameState, week int, rng *rand.Rand) bool {
			// Good harvest in autumn with high soil fertility
			return state.Environment.Season == "autumn" && state.Environment.SoilFertility > 0.7
		},
		weight: func(state *GameState, week int) float64 {
			// Weight proportional to soil fertility
			return state.Environment.SoilFertility * 0.8
		},
		apply: func(state *GameState, week int, rng *rand.Rand) Event {
			// Increase food resource by 20%
			adjustResource(state, economy.ResourceGrain, 0.2)
			return Event{
				ID:        generateEventID("good-harvest", week),
				Type:      "good_harvest",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"yield_increase": 0.2,
					"soil_fertility": state.Environment.SoilFertility,
				},
			}
		},
	},
	{
		eventType: "bad_harvest",
		trigger: func(state *GameState, week int, rng *rand.Rand) bool {
			// Bad harvest in summer with low rainfall
			return state.Environment.Season == "summer" && state.Environment.Rainfall < 5.0
		},
		weight: func(state *GameState, week int) float64 {
			// Weight inverse to rainfall
			return (10.0 - state.Environment.Rainfall) * 0.1
		},
		apply: func(state *GameState, week int, rng *rand.Rand) Event {
			// Decrease food resource by 15%
			adjustResource(state, economy.ResourceGrain, -0.15)
			return Event{
				ID:        generateEventID("bad-harvest", week),
				Type:      "bad_harvest",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"yield_loss": 0.15,
					"rainfall":   state.Environment.Rainfall,
				},
			}
		},
	},
	{
		eventType: "accident",
		trigger: func(state *GameState, week int, rng *rand.Rand) bool {
			// Accidents can happen any time, but more likely with many buildings
			return len(state.Buildings) > 0
		},
		weight: func(state *GameState, week int) float64 {
			// Weight increases with number of buildings
			return float64(len(state.Buildings)) * 0.05
		},
		apply: func(state *GameState, week int, rng *rand.Rand) Event {
			// Accident damages a random building (reduce level by 1)
			if len(state.Buildings) > 0 {
				idx := rng.Intn(len(state.Buildings))
				if state.Buildings[idx].Level > 0 {
					state.Buildings[idx].Level--
				}
			}
			return Event{
				ID:        generateEventID("accident", week),
				Type:      "accident",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"building_damaged": len(state.Buildings) > 0,
				},
			}
		},
	},
}

// adjustResource modifies a resource quantity by a percentage (positive or negative).
// Uses Inventory (always present). If the resource doesn't exist, it does nothing.
func adjustResource(state *GameState, resourceType economy.ResourceType, percent float64) {
	current := state.Inventory.GetAvailable("global", resourceType)
	if current <= 0 {
		return
	}
	change := current * percent
	if change > 0 {
		// Add resource
		r := economy.Resource{
			Type:     resourceType,
			Quantity: change,
			Quality:  economy.QualityNormal,
			Location: "global",
			Produced: economy.GameDate{Year: state.Calendar.Year, Week: state.Calendar.Week},
			Value:    economy.BaseValue(resourceType),
		}
		_ = state.Inventory.AddResource("global", r) // ignore error
	} else if change < 0 {
		// Remove resource (positive amount)
		amount := -change
		_, _ = state.Inventory.RemoveResource("global", resourceType, amount) // ignore error
	}
}

// NewEventSystem creates a new event system.
func NewEventSystem() *EventSystem {
	return &EventSystem{}
}

// Update processes one week of event simulation.
// It generates random events based on game state conditions.
func (e *EventSystem) Update(week int, state *GameState, rng *rand.Rand) []Event {
	var events []Event

	// Determine total weight of all triggered events
	var totalWeight float64
	var triggeredDefs []eventDef
	for _, def := range eventDefinitions {
		if def.trigger(state, week, rng) {
			weight := def.weight(state, week)
			if weight > 0 {
				totalWeight += weight
				triggeredDefs = append(triggeredDefs, def)
			}
		}
	}

	// No eligible events
	if totalWeight <= 0 {
		return events
	}

	// Select at most one event per week (simplification)
	// Use deterministic random selection based on RNG
	roll := rng.Float64() * totalWeight
	var accumulated float64
	for _, def := range triggeredDefs {
		weight := def.weight(state, week)
		accumulated += weight
		if roll < accumulated {
			events = append(events, def.apply(state, week, rng))
			break
		}
	}

	return events
}
