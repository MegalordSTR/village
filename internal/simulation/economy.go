package simulation

import (
	"github.com/vano44/village/internal/economy"
	"log"
	"math/rand"
	"strings"
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

// logRemoveResourceError logs an error from RemoveResource at appropriate severity.
// DEBUG for expected "insufficient quantity" errors, WARN for other errors.
func logRemoveResourceError(err error, operation string, resource economy.ResourceType, quantity float64) {
	if strings.Contains(err.Error(), "insufficient quantity") {
		log.Printf("DEBUG: operation=%s resource=%s quantity=%f message=\"resource not present, skipping\"", operation, resource, quantity)
	} else {
		log.Printf("WARN: operation=%s resource=%s quantity=%f error=%v", operation, resource, quantity, err)
	}
}

// Update processes one week of economic simulation.
func (e *EconomicSystem) Update(week int, state *GameState, rng *rand.Rand) []Event {
	events := make([]Event, 0, 20)

	// Attach storage registry to inventory for capacity enforcement
	state.Inventory.SetStorage(e.storage)

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

	// Inventory is always present after NewGameState

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
			// Log expected "insufficient quantity" as DEBUG, others as WARN
			logRemoveResourceError(err, "CalculateFoodConsumption", rt, float64(foodNeeded))
			continue
		}
		if removed > 0 {
			consumed[rt] = int(removed)
			foodNeeded -= int(removed)
		}
	}

	// No fallback to legacy Resources slice; inventory is the sole source of truth

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
			// Log expected "insufficient quantity" as DEBUG, others as WARN
			logRemoveResourceError(err, "calculateBuildingMaintenance", rt, float64(amount))
			return 0
		}
		return int(removed)
	}

	woodConsumed := consumeFromInventory(economy.ResourceWood, requiredWood)
	stoneConsumed := consumeFromInventory(economy.ResourceStone, requiredStone)
	toolsConsumed := consumeFromInventory(economy.ResourceTools, requiredTools)

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

	// Calculate total storage capacity (legacy logic)
	baseCapacity := 100
	warehouseCapacity := 0
	for _, b := range state.Buildings {
		if b.Type == "warehouse" {
			warehouseCapacity += b.Level * 50
		}
	}
	totalCapacity := baseCapacity + warehouseCapacity

	// Get resources at global location
	resources, ok := state.Inventory.ResourcesMap()["global"]
	if !ok {
		return events
	}

	// Enforce limit per resource type (each resource type cannot exceed total capacity)
	for _, r := range resources {
		totalQty := r.Quantity
		if totalQty <= float64(totalCapacity) {
			continue
		}
		excess := totalQty - float64(totalCapacity)
		// Remove excess from inventory (from global location)
		removed, err := state.Inventory.RemoveResource("global", r.Type, excess)
		if err != nil {
			// Should not happen under normal operation; log WARN
			log.Printf("WARN: operation=enforceStorageLimits resource=%s excess=%f error=%v", r.Type, excess, err)
			continue
		}
		if removed > 0 {
			events = append(events, Event{
				ID:        generateEventID("storage-limit", week),
				Type:      "storage",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"resource": string(r.Type),
					"excess":   int(removed),
					"capacity": totalCapacity,
				},
			})
		}
	}
	return events
}
