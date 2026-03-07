package simulation

import (
	"math/rand"
)

var _ *rand.Rand = nil

// TurnProcessor orchestrates weekly turn processing for all systems.
type TurnProcessor struct {
	systems map[SystemType]System
	order   []SystemType
}

// NewTurnProcessor creates a new TurnProcessor with built-in systems.
func NewTurnProcessor() *TurnProcessor {
	tp := &TurnProcessor{
		systems: make(map[SystemType]System),
		order:   SystemOrder,
	}
	// Initialize built-in systems (to be implemented in future workstreams)
	// For now, they are nil; will be replaced with actual implementations.
	return tp
}

// RegisterSystem adds or replaces a system for the given type.
func (tp *TurnProcessor) RegisterSystem(typ SystemType, sys System) {
	tp.systems[typ] = sys
}

// ProcessWeek advances the simulation by one week, updating all systems in order.
// Returns the events generated during the week.
func (tp *TurnProcessor) ProcessWeek(state *GameState) []Event {
	var allEvents []Event
	week := state.Calendar.Week

	for _, typ := range tp.order {
		sys, ok := tp.systems[typ]
		if !ok || sys == nil {
			// System not yet implemented; skip silently
			continue
		}
		// Create a *rand.Rand from the state's RNG for this system's update.
		// Each system gets its own Rand derived from the RNG state,
		// ensuring deterministic progression across systems.
		rng := state.RNG.Rand()
		events := sys.Update(week, state, rng)
		allEvents = append(allEvents, events...)
	}

	// Add all events to game state history
	for _, ev := range allEvents {
		state.AddEvent(ev)
	}

	// Advance calendar by one week
	state.Calendar.Week++
	if state.Calendar.Week > 52 {
		state.Calendar.Week = 1
		state.Calendar.Month++
		if state.Calendar.Month > 12 {
			state.Calendar.Month = 1
			state.Calendar.Year++
		}
	}

	return allEvents
}
