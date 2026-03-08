package simulation

import (
	"testing"
)

// TestEmptyVillage ensures simulation works with zero residents.
func TestEmptyVillage(t *testing.T) {
	state := NewGameState("empty", 123)
	// No residents added
	tp := NewTurnProcessor()
	tp.RegisterSystem(SystemEnvironment, NewEnvironmentSystem())
	tp.RegisterSystem(SystemProduction, NewProductionSystem())
	tp.RegisterSystem(SystemSocial, NewSocialSystem())
	tp.RegisterSystem(SystemEconomic, NewEconomicSystem())
	tp.RegisterSystem(SystemEvents, NewEventSystem())

	// Process a week - should not panic
	events := tp.ProcessWeek(state)
	// No residents, no social/economic events expected, but environmental events may occur
	_ = events
	// Verify state still valid
	if len(state.Residents) != 0 {
		t.Errorf("expected 0 residents, got %d", len(state.Residents))
	}
}

// TestMaximumResidents processes a village with many residents.
// The benchmark uses 100 residents; we test that the simulation scales.
func TestMaximumResidents(t *testing.T) {
	const numResidents = 100
	state := NewGameState("max", 456)
	for i := 0; i < numResidents; i++ {
		state.AddResident(Resident{
			ID:   string(rune('a' + i%26)),
			Name: "Resident",
			Age:  20 + i%60,
		})
	}
	tp := NewTurnProcessor()
	tp.RegisterSystem(SystemEnvironment, NewEnvironmentSystem())
	tp.RegisterSystem(SystemProduction, NewProductionSystem())
	tp.RegisterSystem(SystemSocial, NewSocialSystem())
	tp.RegisterSystem(SystemEconomic, NewEconomicSystem())
	tp.RegisterSystem(SystemEvents, NewEventSystem())

	events := tp.ProcessWeek(state)
	// Should generate some events
	_ = events
	if len(state.Residents) != numResidents {
		t.Errorf("resident count changed unexpectedly: %d -> %d", numResidents, len(state.Residents))
	}
}

// TestExtremeEnvironmentalValues ensures simulation handles extreme temperature and rainfall.
func TestExtremeEnvironmentalValues(t *testing.T) {
	env := NewEnvironmentSystem()
	state := NewGameState("extreme", 789)

	// Set extreme values
	state.Environment.Temperature = -50.0 // very cold
	state.Environment.Rainfall = 0.0      // drought
	state.Environment.SoilFertility = 0.0 // barren
	state.Environment.ForestHealth = 0.0
	state.Environment.MineQuality = 0.0
	state.Environment.WildlifePopulation = 0.0

	// Process a week - should not panic
	events := env.Update(state.Calendar.Week, state, state.RNG.Rand())
	_ = events

	// Values should stay within reasonable bounds
	if state.Environment.Temperature < -100 || state.Environment.Temperature > 100 {
		t.Errorf("temperature out of bounds: %f", state.Environment.Temperature)
	}
	if state.Environment.Rainfall < 0 {
		t.Errorf("rainfall negative: %f", state.Environment.Rainfall)
	}
	if state.Environment.SoilFertility < 0 || state.Environment.SoilFertility > 1 {
		t.Errorf("soil fertility out of bounds: %f", state.Environment.SoilFertility)
	}
}

// TestExtremeResourceQuantities tests production with huge resource stockpiles.
func TestExtremeResourceQuantities(t *testing.T) {
	state := NewGameState("extreme-resources", 999)
	// Add a massive amount of a resource
	state.AddResource(Resource{
		Type:     "food",
		Quantity: 1_000_000,
		Quality:  1.0,
	})
	prod := NewProductionSystem()
	events := prod.Update(state.Calendar.Week, state, state.RNG.Rand())
	_ = events
	// Should not panic
}
