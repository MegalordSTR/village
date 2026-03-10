package simulation

import (
	"github.com/vano44/village/internal/economy"
	"log"
	"math/rand"
)

// ProductionSystem implements the production simulation system.
// It handles agriculture, mining, crafting, and construction.
type ProductionSystem struct{}

// NewProductionSystem creates a new production system.
func NewProductionSystem() *ProductionSystem {
	return &ProductionSystem{}
}

// Update processes one week of production simulation.
// It updates crop growth, mining extraction, crafting, and construction progress.
func (p *ProductionSystem) Update(week int, state *GameState, rng *rand.Rand) []Event {
	events := make([]Event, 0, len(state.Buildings))

	events = append(events, processAgriculture(week, state, rng)...)
	events = append(events, processMining(week, state, rng)...)
	events = append(events, processCrafting(week, state, rng)...)
	events = append(events, processConstruction(week, state, rng)...)

	return events
}

// processAgriculture handles crop growth and harvesting for farm buildings.
func processAgriculture(week int, state *GameState, rng *rand.Rand) []Event {
	var events []Event

	for i := range state.Buildings {
		b := &state.Buildings[i]
		if b.Type != "farm" {
			continue
		}

		// Ensure metadata map exists
		if b.Metadata == nil {
			b.Metadata = make(map[string]interface{})
		}

		// Initialize crop if not present
		if _, ok := b.Metadata["crop_type"]; !ok {
			// Plant a crop (simplified: always plant wheat)
			b.Metadata["crop_type"] = "wheat"
			b.Metadata["growth_stage"] = 0
			b.Metadata["planted_week"] = week
			continue
		}

		growthStage, _ := b.Metadata["growth_stage"].(int)
		_, _ = b.Metadata["planted_week"].(int) // unused for now

		// Calculate growth chance based on environment
		growthChance := calculateGrowthChance(state.Environment, rng)
		if rng.Float64() < growthChance {
			growthStage++
			b.Metadata["growth_stage"] = growthStage
		}

		// If crop is harvestable (stage >= 3), produce yield and reset
		if growthStage >= 3 {
			yield := calculateYield(state.Environment, b.Level, rng)
			// Add resource
			cropTypeStr := b.Metadata["crop_type"].(string)
			rt, err := StringToResourceType(cropTypeStr)
			if err != nil {
				log.Printf("ERROR: unknown crop type %q, skipping harvest", cropTypeStr)
				// Reset crop (re-plant) and continue to next building
				b.Metadata["growth_stage"] = 0
				b.Metadata["planted_week"] = week
				continue
			}
			_ = AddProducedResource(state, b, rt, yield, 1.0) // TODO: vary quality based on conditions
			// Record event
			events = append(events, Event{
				ID:        generateEventID("harvest", week),
				Type:      "agriculture",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"crop":  b.Metadata["crop_type"],
					"yield": yield,
					"farm":  b.Location,
				},
			})
			// Reset crop (re-plant)
			b.Metadata["growth_stage"] = 0
			b.Metadata["planted_week"] = week
		}
	}

	return events
}

// calculateGrowthChance returns the probability (0-1) that a crop advances a growth stage.
func calculateGrowthChance(env Environment, rng *rand.Rand) float64 {
	// Base chance per week
	baseChance := 0.3

	// Soil fertility effect: higher fertility increases chance
	fertilityEffect := env.SoilFertility * 0.2

	// Rainfall effect: moderate rainfall (10-20mm) is best
	rainEffect := 0.0
	if env.Rainfall >= 10 && env.Rainfall <= 20 {
		rainEffect = 0.1
	} else if env.Rainfall > 20 {
		// Too much rain reduces growth
		rainEffect = -0.05
	} else {
		// Too little rain reduces growth
		rainEffect = -0.05
	}

	// Temperature effect: optimal around 20°C
	tempEffect := 0.0
	if env.Temperature >= 15 && env.Temperature <= 25 {
		tempEffect = 0.1
	} else if env.Temperature < 5 || env.Temperature > 35 {
		tempEffect = -0.1
	}

	// Combine
	chance := baseChance + fertilityEffect + rainEffect + tempEffect
	// Clamp between 0.05 and 0.8
	if chance < 0.05 {
		return 0.05
	}
	if chance > 0.8 {
		return 0.8
	}
	return chance
}

// calculateYield returns the amount of crop harvested.
func calculateYield(env Environment, farmLevel int, rng *rand.Rand) int {
	// Base yield per farm level
	baseYield := 10 * farmLevel

	// Soil fertility multiplier
	fertilityMultiplier := env.SoilFertility

	// Rainfall multiplier: optimal 10-20mm gives 1.0, otherwise less
	var rainMultiplier float64
	if env.Rainfall >= 10 && env.Rainfall <= 20 {
		rainMultiplier = 1.2
	} else if env.Rainfall > 20 {
		rainMultiplier = 0.8
	} else {
		rainMultiplier = 0.7
	}

	// Temperature multiplier: optimal 15-25°C
	tempMultiplier := 1.0
	if env.Temperature >= 15 && env.Temperature <= 25 {
		tempMultiplier = 1.1
	} else if env.Temperature < 5 || env.Temperature > 35 {
		tempMultiplier = 0.6
	}

	// Random variation ±20%
	variation := 0.8 + rng.Float64()*0.4

	yield := float64(baseYield) * fertilityMultiplier * rainMultiplier * tempMultiplier * variation
	if yield < 1 {
		return 1
	}
	return int(yield)
}

// processMining handles resource extraction for mine buildings.
func processMining(week int, state *GameState, rng *rand.Rand) []Event {
	var events []Event

	for i := range state.Buildings {
		b := &state.Buildings[i]
		if b.Type != "mine" {
			continue
		}

		if b.Metadata == nil {
			b.Metadata = make(map[string]interface{})
		}

		// Get depletion level (0-1), where 1 means fully depleted
		depletion, _ := b.Metadata["depletion"].(float64)
		if depletion < 0 {
			depletion = 0
		}
		if depletion > 1 {
			depletion = 1
		}

		// Regenerate depletion slightly each week (natural recovery)
		regeneration := 0.01 + rng.Float64()*0.02
		depletion -= regeneration
		if depletion < 0 {
			depletion = 0
		}

		// Calculate extraction amount based on mine quality, level, and depletion
		baseExtraction := 5 * b.Level
		qualityMultiplier := state.Environment.MineQuality
		depletionMultiplier := 1.0 - depletion // less extraction as depletion increases
		variation := 0.8 + rng.Float64()*0.4

		extraction := float64(baseExtraction) * qualityMultiplier * depletionMultiplier * variation
		if extraction < 0 {
			extraction = 0
		}

		// Increase depletion based on extraction
		depletionIncrease := extraction / 100.0 // each unit extracted increases depletion by 0.01
		depletion += depletionIncrease
		if depletion > 1 {
			depletion = 1
		}

		// Store updated depletion
		b.Metadata["depletion"] = depletion

		if extraction > 0 {
			// Add resource (simplified: always produce "ore")
			_ = AddProducedResource(state, b, economy.ResourceIronOre, int(extraction), qualityMultiplier)

			events = append(events, Event{
				ID:        generateEventID("mining", week),
				Type:      "mining",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"ore":       int(extraction),
					"mine":      b.Location,
					"depletion": depletion,
				},
			})
		}
	}

	return events
}

// processCrafting transforms raw materials into goods at workshop buildings.
func processCrafting(week int, state *GameState, rng *rand.Rand) []Event {
	var events []Event

	for i := range state.Buildings {
		b := &state.Buildings[i]
		if b.Type != "workshop" {
			continue
		}

		if b.Metadata == nil {
			b.Metadata = make(map[string]interface{})
		}

		// Determine recipe (simplified: consume 2 ore, produce 1 tool)
		requiredOre := 2
		// Check if we have enough ore
		oreAvailable := GetAvailableResourceFromState(state, economy.ResourceIronOre)
		if oreAvailable < requiredOre {
			// Not enough resources
			continue
		}
		// Consume ore
		consumed, _ := ConsumeResourceFromState(state, economy.ResourceIronOre, requiredOre)
		if consumed < requiredOre {
			// Should not happen since we checked availability
			continue
		}

		// Produce tool
		toolsProduced := 1 * b.Level
		_ = AddProducedResource(state, b, economy.ResourceTools, toolsProduced, 1.0)

		events = append(events, Event{
			ID:        generateEventID("crafting", week),
			Type:      "crafting",
			Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
			Data: map[string]interface{}{
				"tools_produced": toolsProduced,
				"workshop":       b.Location,
			},
		})
	}

	return events
}

// processConstruction handles building progress at construction sites.
func processConstruction(week int, state *GameState, rng *rand.Rand) []Event { //nolint:gocognit
	var events []Event

	for i := range state.Buildings {
		b := &state.Buildings[i]
		if b.Type != "construction_site" {
			continue
		}

		if b.Metadata == nil {
			b.Metadata = make(map[string]interface{})
		}

		// Get current progress (0-100)
		progress, _ := b.Metadata["progress"].(float64)
		if progress < 0 {
			progress = 0
		}
		if progress >= 100 {
			// Already completed
			continue
		}

		// Determine workers assigned
		workerCount := len(b.Workers)
		// Each worker contributes 1 progress per week (base)
		workerContribution := float64(workerCount)

		// Check for required materials: wood and stone
		woodNeeded := 5
		stoneNeeded := 3
		woodAvailable := GetAvailableResourceFromState(state, economy.ResourceWood)
		stoneAvailable := GetAvailableResourceFromState(state, economy.ResourceStone)

		// If materials insufficient, reduce contribution
		materialMultiplier := 1.0
		if woodAvailable < woodNeeded {
			materialMultiplier *= 0.5
		}
		if stoneAvailable < stoneNeeded {
			materialMultiplier *= 0.5
		}

		// Consume materials if available
		if woodAvailable >= woodNeeded && stoneAvailable >= stoneNeeded {
			// Consume wood
			consumedWood, _ := ConsumeResourceFromState(state, economy.ResourceWood, woodNeeded)
			// Consume stone
			consumedStone, _ := ConsumeResourceFromState(state, economy.ResourceStone, stoneNeeded)

			_ = consumedWood
			_ = consumedStone // silence unused variable warnings
		}

		// Calculate progress increase
		progressIncrease := workerContribution * materialMultiplier
		progress += progressIncrease
		if progress > 100 {
			progress = 100
		}
		b.Metadata["progress"] = progress

		// If construction completed, generate event
		if progress >= 100 {
			events = append(events, Event{
				ID:        generateEventID("construction-complete", week),
				Type:      "construction",
				Timestamp: formatWeekTimestamp(state.Calendar.Year, week),
				Data: map[string]interface{}{
					"building": b.Location,
				},
			})
			// Optionally change building type? For now, keep as construction_site
		}
	}

	return events
}
