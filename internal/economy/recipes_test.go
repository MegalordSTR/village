package economy

import (
	"math"
	"testing"
)

func TestRecipeStruct(t *testing.T) {
	recipe := Recipe{
		ID:         "test",
		Name:       "Test Recipe",
		Building:   "test_building",
		Skill:      "test_skill",
		Inputs:     []ResourceRequirement{},
		Outputs:    []ResourceOutput{},
		Time:       1.0,
		BaseYield:  1.0,
		MaxWorkers: 1,
	}
	if recipe.ID != "test" {
		t.Errorf("expected ID test, got %s", recipe.ID)
	}
	if recipe.Name != "Test Recipe" {
		t.Errorf("expected Name Test Recipe, got %s", recipe.Name)
	}
	if recipe.Building != "test_building" {
		t.Errorf("expected Building test_building, got %s", recipe.Building)
	}
	if recipe.Skill != "test_skill" {
		t.Errorf("expected Skill test_skill, got %s", recipe.Skill)
	}
	if recipe.Time != 1.0 {
		t.Errorf("expected Time 1.0, got %f", recipe.Time)
	}
	if recipe.BaseYield != 1.0 {
		t.Errorf("expected BaseYield 1.0, got %f", recipe.BaseYield)
	}
	if recipe.MaxWorkers != 1 {
		t.Errorf("expected MaxWorkers 1, got %d", recipe.MaxWorkers)
	}
}

func TestResourceRequirement(t *testing.T) {
	req := ResourceRequirement{
		Type:       ResourceGrain,
		Quantity:   10,
		MinQuality: QualityNormal,
	}
	if req.Type != ResourceGrain {
		t.Errorf("expected Type grain, got %s", req.Type)
	}
	if req.Quantity != 10 {
		t.Errorf("expected Quantity 10, got %f", req.Quantity)
	}
	if req.MinQuality != QualityNormal {
		t.Errorf("expected MinQuality Normal, got %v", req.MinQuality)
	}
}

func TestResourceOutput(t *testing.T) {
	out := ResourceOutput{
		Type:               ResourceBread,
		Quantity:           5,
		QualityInheritance: 0.9,
	}
	if out.Type != ResourceBread {
		t.Errorf("expected Type bread, got %s", out.Type)
	}
	if out.Quantity != 5 {
		t.Errorf("expected Quantity 5, got %f", out.Quantity)
	}
	if out.QualityInheritance != 0.9 {
		t.Errorf("expected QualityInheritance 0.9, got %f", out.QualityInheritance)
	}
}

func TestAllRecipesCount(t *testing.T) {
	recipes := AllRecipes()
	// At least 10 recipes covering all 4 chains
	if len(recipes) < 10 {
		t.Errorf("expected at least 10 recipes, got %d", len(recipes))
	}
	// Ensure each chain represented
	chains := map[string]bool{"food": false, "wood": false, "metal": false, "textile": false}
	for _, r := range recipes {
		for _, out := range r.Outputs {
			switch out.Type {
			case ResourceFlour, ResourceBread:
				chains["food"] = true
			case ResourcePlanks, ResourceFurniture:
				chains["wood"] = true
			case ResourceIron, ResourceTools, ResourceWeapons:
				chains["metal"] = true
			case ResourceCloth, ResourceClothing:
				chains["textile"] = true
			}
		}
	}
	for chain, present := range chains {
		if !present {
			t.Errorf("chain %s not represented in recipes", chain)
		}
	}
}

func TestFindRecipeByID(t *testing.T) {
	recipe := FindRecipeByID("mill_flour")
	if recipe == nil {
		t.Fatal("expected recipe mill_flour not found")
	}
	if recipe.Name != "Mill Flour" {
		t.Errorf("expected Name Mill Flour, got %s", recipe.Name)
	}
	// Not found
	notFound := FindRecipeByID("nonexistent")
	if notFound != nil {
		t.Error("expected nil for nonexistent recipe")
	}
}

func TestProductionChain(t *testing.T) {
	chain := ProductionChain(ResourceBread)
	if len(chain) == 0 {
		t.Fatal("expected at least one step for bread chain")
	}
	// Should include mill_flour and bakery_bread
	var foundMill, foundBakery bool
	for _, step := range chain {
		if step.RecipeID == "mill_flour" {
			foundMill = true
		}
		if step.RecipeID == "bakery_bread" {
			foundBakery = true
		}
	}
	if !foundMill || !foundBakery {
		t.Errorf("missing expected recipes: mill %v, bakery %v", foundMill, foundBakery)
	}
	// Check depth ordering: raw materials first (higher depth?)
	// Actually depth is computed as distance from final product? We'll just ensure ordering.
	for i := 1; i < len(chain); i++ {
		if chain[i].Depth < chain[i-1].Depth {
			t.Errorf("chain not sorted by depth ascending at step %d", i)
		}
	}
}

func TestBottleneckResources(t *testing.T) {
	chain := ProductionChain(ResourceBread)
	bottlenecks := BottleneckResources(chain)
	// No bottleneck threshold set high, so likely empty.
	// Just ensure function doesn't panic.
	_ = bottlenecks
}

func TestCalculateProductionFullBatch(t *testing.T) {
	ctx := ProductionContext{
		RecipeID: "mill_flour",
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 100, Quality: QualityNormal},
		},
		WorkerSkill:       0.8,
		BuildingCondition: 1.0,
		RandSeed:          42,
	}
	result := CalculateProduction(ctx)
	if !result.Success {
		t.Fatal("expected successful production")
	}
	if len(result.OutputResources) != 1 {
		t.Fatalf("expected 1 output resource, got %d", len(result.OutputResources))
	}
	out := result.OutputResources[0]
	if out.Type != ResourceFlour {
		t.Errorf("expected output type flour, got %s", out.Type)
	}
	// With 100 grain, requirement 10 per batch, max batches = 10.
	// Base output quantity = 8 per batch * 10 = 80.
	// Yield multiplier includes skill bonus (0.8 skill => +16%) and random variation +/-5%.
	// Compute expected multiplier: skill factor = 1 + 0.8*0.2 = 1.16
	// Random variation with seed 42: 0.95 + rng.Float64()*0.1 = 0.987303 (computed earlier)
	// total multiplier = 1.16 * 0.987303 = 1.14527048
	expectedMultiplier := 1.14527048
	expectedQuantity := 8.0 * 10.0 * expectedMultiplier
	// Allow tiny floating point difference
	if math.Abs(out.Quantity-expectedQuantity) > 1e-4 {
		t.Errorf("expected output quantity %f, got %f", expectedQuantity, out.Quantity)
	}
	// Check consumed inputs
	consumed, ok := result.ConsumedInputs[ResourceGrain]
	if !ok {
		t.Fatal("expected grain consumed")
	}
	if consumed != 100.0 {
		t.Errorf("expected consumed grain 100, got %f", consumed)
	}
}

func TestCalculateProductionPartialBatch(t *testing.T) {
	ctx := ProductionContext{
		RecipeID: "mill_flour",
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 5, Quality: QualityNormal},
		},
		WorkerSkill:       0.8,
		BuildingCondition: 1.0,
		RandSeed:          42,
	}
	result := CalculateProduction(ctx)
	if !result.Success {
		t.Fatal("expected successful partial production")
	}
	if len(result.OutputResources) != 1 {
		t.Fatalf("expected 1 output resource, got %d", len(result.OutputResources))
	}
	out := result.OutputResources[0]
	// 5 grain vs required 10, scale = 0.5.
	// Yield multiplier same as full batch test: 1.14527048
	expectedMultiplier := 1.14527048
	expectedQuantity := 8.0 * 0.5 * expectedMultiplier
	if math.Abs(out.Quantity-expectedQuantity) > 1e-4 {
		t.Errorf("expected output quantity %f, got %f", expectedQuantity, out.Quantity)
	}
	// consumed should be 5
	consumed, ok := result.ConsumedInputs[ResourceGrain]
	if !ok {
		t.Fatal("expected grain consumed")
	}
	if consumed != 5.0 {
		t.Errorf("expected consumed grain 5, got %f", consumed)
	}
}

func TestCalculateProductionInsufficientInputs(t *testing.T) {
	ctx := ProductionContext{
		RecipeID: "mill_flour",
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 0, Quality: QualityNormal},
		},
		WorkerSkill:       0.8,
		BuildingCondition: 1.0,
		RandSeed:          42,
	}
	result := CalculateProduction(ctx)
	if result.Success {
		t.Error("expected production to fail with zero inputs")
	}
}

func TestQualityInheritance(t *testing.T) {
	ctx := ProductionContext{
		RecipeID: "mill_flour",
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 100, Quality: QualityExcellent},
		},
		WorkerSkill:       0.5, // low skill
		BuildingCondition: 1.0,
		RandSeed:          42,
	}
	result := CalculateProduction(ctx)
	if !result.Success {
		t.Fatal("expected successful production")
	}
	out := result.OutputResources[0]
	// Input quality Excellent = 3 (0 Poor,1 Normal,2 Good,3 Excellent,4 Masterwork)
	// QualityToFloat(Excellent) = 3/4 = 0.75
	// QualityInheritance factor 0.8 => inherited = 0.75 * 0.8 = 0.6
	// Skill contribution = 0.5 * (1-0.8)=0.1
	// final quality = 0.7 => FloatToQuality(0.7) = round(0.7*4) = round(2.8) = 3 => Excellent
	expectedQuality := QualityExcellent
	if out.Quality != expectedQuality {
		t.Errorf("expected output quality %v, got %v", expectedQuality, out.Quality)
	}
}

func TestDeterministicOutput(t *testing.T) {
	ctx1 := ProductionContext{
		RecipeID: "mill_flour",
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 100, Quality: QualityNormal},
		},
		WorkerSkill:       0.7,
		BuildingCondition: 0.9,
		RandSeed:          123,
	}
	ctx2 := ProductionContext{
		RecipeID: "mill_flour",
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 100, Quality: QualityNormal},
		},
		WorkerSkill:       0.7,
		BuildingCondition: 0.9,
		RandSeed:          123,
	}
	result1 := CalculateProduction(ctx1)
	result2 := CalculateProduction(ctx2)
	if result1.OutputResources[0].Quantity != result2.OutputResources[0].Quantity {
		t.Errorf("determinism broken: quantities differ %f vs %f",
			result1.OutputResources[0].Quantity, result2.OutputResources[0].Quantity)
	}
	if result1.OutputResources[0].Quality != result2.OutputResources[0].Quality {
		t.Errorf("determinism broken: qualities differ %v vs %v",
			result1.OutputResources[0].Quality, result2.OutputResources[0].Quality)
	}
}

func TestBuildingConditionAffectsSkill(t *testing.T) {
	ctx := ProductionContext{
		RecipeID: "mill_flour",
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 100, Quality: QualityNormal},
		},
		WorkerSkill:       1.0,
		BuildingCondition: 0.5, // halves effective skill
		RandSeed:          42,
	}
	result := CalculateProduction(ctx)
	// Effective skill should be 0.5
	if result.EffectiveSkill != 0.5 {
		t.Errorf("expected effective skill 0.5, got %f", result.EffectiveSkill)
	}
}
