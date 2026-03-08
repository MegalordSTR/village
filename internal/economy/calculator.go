package economy

import (
	"math/rand"
)

// ProductionContext holds the inputs and configuration for a production calculation.
type ProductionContext struct {
	RecipeID          string
	InputResources    map[ResourceType]Resource // available inputs, keyed by type
	WorkerSkill       float64                   // 0..1, where 1 is expert
	BuildingCondition float64                   // 0..1, where 1 is perfect condition
	RandSeed          int64                     // seed for deterministic RNG
}

// ProductionResult holds the outputs and metadata of a production run.
type ProductionResult struct {
	OutputResources []Resource
	ConsumedInputs  map[ResourceType]float64 // quantity consumed per type
	ActualYield     float64                  // multiplier applied
	EffectiveSkill  float64                  // skill after building modifier
	QualityBonus    float64                  // bonus to output quality
	Success         bool                     // whether any output was produced
}

// computeEffectiveSkill returns worker skill adjusted by building condition.
func computeEffectiveSkill(workerSkill, buildingCondition float64) float64 {
	effective := workerSkill * buildingCondition
	if effective > 1.0 {
		effective = 1.0
	}
	return effective
}

// computeYieldMultiplier returns the multiplier for output quantities based on skill and random variation.
func computeYieldMultiplier(effectiveSkill, baseYield float64, rng *rand.Rand) float64 {
	// Skill bonus: up to +20%
	multiplier := baseYield * (1.0 + effectiveSkill*0.2)
	// Random variation +/-5%
	variation := 0.95 + rng.Float64()*0.1
	multiplier *= variation
	return multiplier
}

// computeOutputQuality calculates the quality tier for an output given inputs, inheritance factor, and effective skill.
func computeOutputQuality(inputs []ResourceRequirement, consumed map[ResourceType]float64, available map[ResourceType]Resource, inheritance float64, effectiveSkill float64) QualityTier {
	var totalQual float64
	var totalQuant float64
	for _, req := range inputs {
		avail := available[req.Type]
		qual := QualityToFloat(avail.Quality)
		// weight by consumed quantity (which is proportional to requirement)
		totalQual += qual * consumed[req.Type]
		totalQuant += consumed[req.Type]
	}
	if totalQuant == 0 {
		return QualityNormal
	}
	avgInputQuality := totalQual / totalQuant
	inherited := avgInputQuality * inheritance
	skillContrib := effectiveSkill * (1.0 - inheritance)
	finalQuality := inherited + skillContrib
	if finalQuality > 1.0 {
		finalQuality = 1.0
	}
	return FloatToQuality(finalQuality)
}

// CalculateProduction computes the outputs of a recipe given the context.
// It returns a ProductionResult describing what was produced and what inputs were consumed.
func CalculateProduction(ctx ProductionContext) ProductionResult {
	recipe := FindRecipeByID(ctx.RecipeID)
	if recipe == nil {
		return ProductionResult{Success: false}
	}

	// Deterministic RNG
	rng := rand.New(rand.NewSource(ctx.RandSeed))

	effectiveSkill := computeEffectiveSkill(ctx.WorkerSkill, ctx.BuildingCondition)

	// Determine how many batches can be produced based on inputs.
	// For each input requirement, compute max batches = floor(available / required).
	maxBatches := -1
	for _, req := range recipe.Inputs {
		avail, ok := ctx.InputResources[req.Type]
		if !ok {
			maxBatches = 0
			break
		}
		if avail.Quantity < req.Quantity {
			maxBatches = 0
			break
		}
		batches := int(avail.Quantity / req.Quantity)
		if maxBatches == -1 || batches < maxBatches {
			maxBatches = batches
		}
	}

	var scale float64
	var numBatches int
	if maxBatches <= 0 {
		// Partial production: compute fractional batch size.
		var minRatio float64 = 1.0
		for _, req := range recipe.Inputs {
			avail, ok := ctx.InputResources[req.Type]
			if !ok {
				minRatio = 0.0
				break
			}
			ratio := avail.Quantity / req.Quantity
			if ratio < minRatio {
				minRatio = ratio
			}
		}
		if minRatio <= 0.0 {
			return ProductionResult{Success: false}
		}
		scale = minRatio
		numBatches = 0 // indicates fractional batch
	} else {
		scale = 1.0
		numBatches = maxBatches
	}

	// Compute consumed inputs.
	consumed := make(map[ResourceType]float64)
	for _, req := range recipe.Inputs {
		if numBatches == 0 {
			consumed[req.Type] = req.Quantity * scale
		} else {
			consumed[req.Type] = req.Quantity * float64(numBatches)
		}
	}

	// Compute yield multiplier (same for all outputs).
	yieldMultiplier := computeYieldMultiplier(effectiveSkill, recipe.BaseYield, rng)

	// Produce outputs.
	var outputs []Resource
	for _, out := range recipe.Outputs {
		quality := computeOutputQuality(recipe.Inputs, consumed, ctx.InputResources, out.QualityInheritance, effectiveSkill)
		var quantity float64
		if numBatches == 0 {
			quantity = out.Quantity * scale * yieldMultiplier
		} else {
			quantity = out.Quantity * float64(numBatches) * yieldMultiplier
		}
		outputs = append(outputs, Resource{
			Type:     out.Type,
			Quantity: quantity,
			Quality:  quality,
			Location: "",
			Produced: GameDate{},
			Value:    BaseValue(out.Type),
		})
	}

	return ProductionResult{
		OutputResources: outputs,
		ConsumedInputs:  consumed,
		ActualYield:     yieldMultiplier,
		EffectiveSkill:  effectiveSkill,
		QualityBonus:    effectiveSkill,
		Success:         true,
	}
}
