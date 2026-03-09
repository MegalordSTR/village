package simulation

import "math/rand"

// EconomicSystem implements the economic simulation system.
// It handles resource consumption, inventory management, trade, and wealth distribution.
type EconomicSystem struct {
	// wealth tracks each resident's wealth (optional)
	wealth map[string]float64
}

// NewEconomicSystem creates a new economic system.
func NewEconomicSystem() *EconomicSystem {
	return &EconomicSystem{
		wealth: make(map[string]float64),
	}
}

// Update processes one week of economic simulation.
func (e *EconomicSystem) Update(week int, state *GameState, rng *rand.Rand) []Event {
	events := make([]Event, 0, 20)

	// 1. Resource consumption based on resident needs
	events = append(events, e.consumeResourcesByResidents(week, state, rng)...)

	// 2. Building maintenance resource consumption
	events = append(events, e.consumeMaintenanceResources(week, state, rng)...)

	// 3. Inventory management with storage limits and spoilage
	events = append(events, e.applySpoilage(week, state, rng)...)

	// 4. Enforce storage limits (discard excess)
	events = append(events, e.enforceStorageLimits(week, state, rng)...)

	// 5. Trade between villages (placeholder)
	// No implementation yet, returns empty slice

	// 6. Wealth distribution affecting happiness and productivity
	e.updateWealth(state, rng)
	e.adjustHappinessForWealth(state)

	// 7. Resource price fluctuations based on scarcity (optional)
	// Not implemented yet

	return events
}

// consumeResourcesByResidents handles food consumption based on resident needs.
func (e *EconomicSystem) consumeResourcesByResidents(week int, state *GameState, rng *rand.Rand) []Event {
	events := make([]Event, 0, len(state.Residents)*2)
	if len(state.Residents) == 0 {
		return events
	}

	// Calculate total food needed
	foodNeeded := len(state.Residents) // each resident needs 1 food per week
	// Adjust based on hunger need? For now simple.

	// Find food resources
	foodAvailable := 0
	foodIndex := -1
	for i, r := range state.Resources {
		if r.Type == "food" {
			foodAvailable = r.Quantity
			foodIndex = i
			break
		}
	}

	// Determine how much to consume
	consumeAmount := foodNeeded
	if consumeAmount > foodAvailable {
		consumeAmount = foodAvailable
	}

	// Consume
	if foodIndex >= 0 && consumeAmount > 0 {
		state.Resources[foodIndex].Quantity -= consumeAmount
		// Remove resource if quantity zero
		if state.Resources[foodIndex].Quantity <= 0 {
			// Remove from slice
			state.Resources = append(state.Resources[:foodIndex], state.Resources[foodIndex+1:]...)
		}

		// Generate consumption event
		events = append(events, Event{
			ID:        generateEventID("consumption", week),
			Type:      "consumption",
			Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
			Data: map[string]interface{}{
				"food_consumed": consumeAmount,
				"residents":     len(state.Residents),
			},
		})
	}

	return events
}

// consumeMaintenanceResources handles building maintenance.
func (e *EconomicSystem) consumeMaintenanceResources(week int, state *GameState, rng *rand.Rand) []Event {
	events := make([]Event, 0, len(state.Buildings))
	if len(state.Buildings) == 0 {
		return events
	}

	// Maintenance requirements per building type per level
	type maintenanceReq struct {
		wood  int
		stone int
		tools int
	}
	requirements := map[string]maintenanceReq{
		"house":     {wood: 1, stone: 0, tools: 0},
		"farm":      {wood: 1, stone: 0, tools: 1},
		"mine":      {wood: 2, stone: 1, tools: 2},
		"workshop":  {wood: 2, stone: 1, tools: 1},
		"warehouse": {wood: 3, stone: 2, tools: 1},
		// construction_site consumes resources via construction system, not maintenance
	}

	// Total required resources
	requiredWood := 0
	requiredStone := 0
	requiredTools := 0
	for _, b := range state.Buildings {
		req, ok := requirements[b.Type]
		if !ok {
			continue
		}
		requiredWood += req.wood * b.Level
		requiredStone += req.stone * b.Level
		requiredTools += req.tools * b.Level
	}

	// Helper to consume a specific resource type
	consumeResource := func(resType string, amount int) int {
		consumed := 0
		for i := range state.Resources {
			if state.Resources[i].Type == resType && state.Resources[i].Quantity > 0 {
				available := state.Resources[i].Quantity
				take := amount - consumed
				if take > available {
					take = available
				}
				state.Resources[i].Quantity -= take
				consumed += take
				if state.Resources[i].Quantity <= 0 {
					// Remove zero quantity resource
					state.Resources = append(state.Resources[:i], state.Resources[i+1:]...)
					break // need to restart iteration but fine for now
				}
				if consumed >= amount {
					break
				}
			}
		}
		return consumed
	}

	woodConsumed := consumeResource("wood", requiredWood)
	stoneConsumed := consumeResource("stone", requiredStone)
	toolsConsumed := consumeResource("tool", requiredTools) // note: resource type "tool" not "tools"

	if woodConsumed > 0 || stoneConsumed > 0 || toolsConsumed > 0 {
		events = append(events, Event{
			ID:        generateEventID("maintenance", week),
			Type:      "maintenance",
			Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
			Data: map[string]interface{}{
				"wood_consumed":  woodConsumed,
				"stone_consumed": stoneConsumed,
				"tools_consumed": toolsConsumed,
				"buildings":      len(state.Buildings),
			},
		})
	}

	return events
}

// applySpoilage reduces food quantity over time.
func (e *EconomicSystem) applySpoilage(week int, state *GameState, rng *rand.Rand) []Event {
	events := make([]Event, 0, len(state.Resources))
	spoilageRate := 0.05 // 5% spoilage per week

	for i := range state.Resources {
		r := &state.Resources[i]
		// Only spoil food types
		if r.Type == "food" || r.Type == "grain" || r.Type == "meat" {
			spoiled := int(float64(r.Quantity) * spoilageRate)
			if spoiled < 1 && r.Quantity > 0 {
				spoiled = 1 // at least 1 unit spoils if any quantity
			}
			if spoiled > 0 {
				r.Quantity -= spoiled
				// If quantity becomes zero, mark for removal
				if r.Quantity <= 0 {
					state.Resources[i].Quantity = 0
				}
				events = append(events, Event{
					ID:        generateEventID("spoilage", week),
					Type:      "spoilage",
					Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
					Data: map[string]interface{}{
						"resource": r.Type,
						"amount":   spoiled,
					},
				})
			}
		}
	}

	// Remove zero quantity resources
	newResources := make([]Resource, 0, len(state.Resources))
	for _, r := range state.Resources {
		if r.Quantity > 0 {
			newResources = append(newResources, r)
		}
	}
	state.Resources = newResources

	return events
}

// updateWealth updates resident wealth based on work and skills.
func (e *EconomicSystem) updateWealth(state *GameState, rng *rand.Rand) {
	// Initialize wealth for residents missing from map
	for _, r := range state.Residents {
		if _, ok := e.wealth[r.ID]; !ok {
			e.wealth[r.ID] = 0.0
		}
	}

	// Increase wealth based on employment and skill level
	for _, r := range state.Residents {
		// Check if resident has work
		hasWork := false
		for _, b := range state.Buildings {
			for _, workerID := range b.Workers {
				if workerID == r.ID {
					hasWork = true
					break
				}
			}
			if hasWork {
				break
			}
		}

		if hasWork {
			// Base income
			income := 1.0
			// Add skill bonus
			skillBonus := 0.0
			for _, skill := range r.Skills {
				skillBonus += float64(skill.Level) * 0.1
			}
			e.wealth[r.ID] += income + skillBonus
		}
	}
}

// adjustHappinessForWealth modifies resident happiness need based on wealth.
func (e *EconomicSystem) adjustHappinessForWealth(state *GameState) {
	for i := range state.Residents {
		r := &state.Residents[i]
		wealth, ok := e.wealth[r.ID]
		if !ok {
			continue
		}

		// Find happiness need
		for j := range r.Needs {
			if r.Needs[j].ID == "happiness" {
				// Wealthier residents are happier (max effect +0.2)
				wealthEffect := wealth * 0.01
				if wealthEffect > 0.2 {
					wealthEffect = 0.2
				}
				r.Needs[j].Level -= wealthEffect // lower level = better
				if r.Needs[j].Level < 0 {
					r.Needs[j].Level = 0
				}
				break
			}
		}
	}
}

// enforceStorageLimits ensures resources do not exceed storage capacity.
func (e *EconomicSystem) enforceStorageLimits(week int, state *GameState, rng *rand.Rand) []Event {
	events := make([]Event, 0, len(state.Buildings))
	// Calculate total storage capacity
	baseCapacity := 100
	warehouseCapacity := 0
	for _, b := range state.Buildings {
		if b.Type == "warehouse" {
			warehouseCapacity += b.Level * 50
		}
	}
	totalCapacity := baseCapacity + warehouseCapacity

	// Apply to food resources
	for i := range state.Resources {
		r := &state.Resources[i]
		if r.Type == "food" || r.Type == "grain" || r.Type == "meat" {
			if r.Quantity > totalCapacity {
				excess := r.Quantity - totalCapacity
				r.Quantity = totalCapacity
				events = append(events, Event{
					ID:        generateEventID("storage-limit", week),
					Type:      "storage",
					Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
					Data: map[string]interface{}{
						"resource": r.Type,
						"excess":   excess,
						"capacity": totalCapacity,
					},
				})
			}
		}
	}
	return events
}
