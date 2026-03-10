package simulation

import (
	"github.com/vano44/village/internal/economy"
	"math/rand"
)

// EconomicSystem implements the economic simulation system.
// It handles resource consumption, inventory management, trade, and wealth distribution.
type EconomicSystem struct {
	// wealth tracks each resident's wealth (optional)
	wealth map[string]float64
	// storage registry for capacity and spoilage
	storage *economy.StorageRegistry
}

// NewEconomicSystem creates a new economic system.
func NewEconomicSystem() *EconomicSystem {
	storage := economy.NewStorageRegistry()
	// Add a default outdoor pile for global storage
	storage.AddBuilding(&economy.StorageBuilding{
		ID:      "global",
		Type:    economy.StorageGranary,
		Level:   1,
		Quality: economy.StorageQualityGood,
	})
	return &EconomicSystem{
		wealth:  make(map[string]float64),
		storage: storage,
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

	// Ensure inventory exists and is synced
	state.SyncInventory()

	// Food priority: bread > vegetables > grain
	foodTypes := []economy.ResourceType{economy.ResourceBread, economy.ResourceVegetables, economy.ResourceGrain}
	foodNeeded := len(state.Residents)
	consumed := make(map[economy.ResourceType]int)

	// Try to consume from inventory in priority order
	for _, rt := range foodTypes {
		if foodNeeded <= 0 {
			break
		}
		// Remove as much as possible from inventory at "global" location
		removed, err := state.Inventory.RemoveResource("global", rt, float64(foodNeeded))
		if err != nil {
			// Ignore error (e.g., resource not present)
			continue
		}
		if removed > 0 {
			consumed[rt] = int(removed)
			foodNeeded -= int(removed)
		}
	}

	// If we still need food after checking inventory, fall back to legacy Resources slice
	if foodNeeded > 0 && state.Resources != nil {
		// This should rarely happen if inventory is synced, but keep for backward compatibility
		for i := range state.Resources {
			if state.Resources[i].Type == economy.ResourceGrain {
				available := state.Resources[i].Quantity
				take := foodNeeded
				if take > available {
					take = available
				}
				state.Resources[i].Quantity -= take
				consumed[economy.ResourceGrain] = consumed[economy.ResourceGrain] + take // treat generic "food" as grain
				foodNeeded -= take
				if state.Resources[i].Quantity <= 0 {
					// Remove zero quantity resource
					state.Resources = append(state.Resources[:i], state.Resources[i+1:]...)
				}
				break
			}
		}
	}

	// After inventory modifications, sync Resources slice for any remaining legacy code
	state.SyncResources()

	// Generate consumption event
	totalConsumed := 0
	for _, amount := range consumed {
		totalConsumed += amount
	}
	if totalConsumed > 0 {
		// Convert map to simple map[string]interface{} for event data
		consumedData := make(map[string]int)
		for rt, amount := range consumed {
			consumedData[string(rt)] = amount
		}
		events = append(events, Event{
			ID:        generateEventID("consumption", week),
			Type:      "consumption",
			Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
			Data: map[string]interface{}{
				"food_consumed": totalConsumed,
				"residents":     len(state.Residents),
				"breakdown":     consumedData,
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

	// Ensure inventory exists and is synced
	state.SyncInventory()

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

	// Consume resources from inventory
	consumeFromInventory := func(rt economy.ResourceType, amount int) int {
		if amount <= 0 {
			return 0
		}
		removed, err := state.Inventory.RemoveResource("global", rt, float64(amount))
		if err != nil {
			return 0
		}
		return int(removed)
	}

	woodConsumed := consumeFromInventory(economy.ResourceWood, requiredWood)
	stoneConsumed := consumeFromInventory(economy.ResourceStone, requiredStone)
	toolsConsumed := consumeFromInventory(economy.ResourceTools, requiredTools)

	// If inventory didn't have enough, fall back to legacy Resources slice (should be synced already)
	if woodConsumed < requiredWood || stoneConsumed < requiredStone || toolsConsumed < requiredTools {
		// This should rarely happen if inventory is synced, but keep for backward compatibility
		// We'll just accept whatever was consumed from inventory (no further fallback)
		// In practice, SyncInventory loads all resources, so inventory should have everything.
	}

	// After inventory modifications, sync Resources slice for any remaining legacy code
	state.SyncResources()

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
	events := make([]Event, 0, 10)
	// Ensure inventory exists and is synced
	state.SyncInventory()
	if state.Inventory == nil {
		// Should not happen after SyncInventory, but fallback
		return events
	}

	// Convert simulation week to economy game date
	currentDate := economy.GameDate{
		Year: state.Calendar.Year,
		Week: week,
	}

	// Apply spoilage to inventory using storage registry (may be nil)
	spoiledTotals := economy.ApplySpoilageToInventory(state.Inventory, e.storage, currentDate)

	// Generate events for each resource type that spoiled
	for rtStr, spoiledQty := range spoiledTotals {
		if spoiledQty <= 0 {
			continue
		}
		events = append(events, Event{
			ID:        generateEventID("spoilage", week),
			Type:      "spoilage",
			Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
			Data: map[string]interface{}{
				"resource": rtStr,
				"amount":   int(spoiledQty),
			},
		})
	}

	// Sync Resources slice to reflect inventory changes
	state.SyncResources()

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
	// Ensure inventory exists and is synced
	state.SyncInventory()
	if state.Inventory == nil {
		// Should not happen after SyncInventory, but fallback
		return events
	}

	// Calculate total storage capacity (legacy logic)
	baseCapacity := 100
	warehouseCapacity := 0
	for _, b := range state.Buildings {
		if b.Type == "warehouse" {
			warehouseCapacity += b.Level * 50
		}
	}
	totalCapacity := baseCapacity + warehouseCapacity

	// Food resource types to enforce limits on
	foodTypes := []economy.ResourceType{
		economy.ResourceGrain,
		economy.ResourceVegetables,
		economy.ResourceBread,
	}

	// Iterate over food types
	for _, rt := range foodTypes {
		totalQty := state.Inventory.GetAvailable("global", rt)
		if totalQty > float64(totalCapacity) {
			excess := totalQty - float64(totalCapacity)
			// Remove excess from inventory (from global location)
			removed, err := state.Inventory.RemoveResource("global", rt, excess)
			if err != nil {
				// Should not happen
				continue
			}
			if removed > 0 {
				events = append(events, Event{
					ID:        generateEventID("storage-limit", week),
					Type:      "storage",
					Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
					Data: map[string]interface{}{
						"resource": string(rt),
						"excess":   int(removed),
						"capacity": totalCapacity,
					},
				})
			}
		}
	}

	// Also check legacy "food" resource in Resources slice (should be synced)
	// This is redundant but ensures backward compatibility
	for i := range state.Resources {
		r := &state.Resources[i]
		if StringToResourceType(string(r.Type)) == economy.ResourceGrain {
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

	// Sync Resources slice to reflect inventory changes
	state.SyncResources()

	return events
}
