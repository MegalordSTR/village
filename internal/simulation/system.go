package simulation

import "math/rand"

// System represents a simulation system that processes weekly updates.
type System interface {
	// Update processes one week of simulation for this system.
	// Returns events generated during the update.
	Update(week int, state *GameState, rng *rand.Rand) []Event
}

// SystemType identifies a built-in system.
type SystemType string

const (
	SystemEnvironment SystemType = "environment"
	SystemProduction  SystemType = "production"
	SystemSocial      SystemType = "social"
	SystemEconomic    SystemType = "economic"
	SystemEvents      SystemType = "events"
)

// SystemOrder defines the deterministic execution order of systems.
var SystemOrder = []SystemType{
	SystemEnvironment,
	SystemProduction,
	SystemSocial,
	SystemEconomic,
	SystemEvents,
}
