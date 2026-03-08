package simulation

import (
	"math/rand"
	"testing"
)

func TestProductionSystemImplementsSystem(t *testing.T) {
	var _ System = (*ProductionSystem)(nil)
}

func TestNewProductionSystem(t *testing.T) {
	prod := NewProductionSystem()
	if prod == nil {
		t.Fatal("NewProductionSystem returned nil")
	}
}

func TestProductionUpdateDeterministic(t *testing.T) {
	prod := NewProductionSystem()

	// Create two states with same seed
	state1 := NewGameState("test1", 123)
	state2 := NewGameState("test2", 123)

	// Process first week
	events1 := prod.Update(1, state1, state1.RNG.Rand())
	events2 := prod.Update(1, state2, state2.RNG.Rand())

	// Should generate same events
	if len(events1) != len(events2) {
		t.Errorf("event count mismatch: %d vs %d", len(events1), len(events2))
	}
	// TODO: add more deterministic checks once we have production outputs
}

func TestProductionUpdateChangesResources(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 456)

	events := prod.Update(state.Calendar.Week, state, state.RNG.Rand())

	// Should generate some events (maybe)
	_ = events
	// Should change resources count (if there are production buildings)
	// For now, just ensure no panic
}

func TestProductionAgriculture(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 789)

	// Add a farm building
	state.AddBuilding(Building{
		Type:       "farm",
		Location:   "field1",
		Level:      1,
		Workers:    []string{},
		Production: []Production{},
		Metadata:   nil,
	})

	initialResources := len(state.Resources)

	// Run several weeks to allow growth and harvest
	for week := 1; week <= 10; week++ {
		prod.Update(week, state, state.RNG.Rand())
	}

	// Expect at least some resources added
	if len(state.Resources) <= initialResources {
		t.Errorf("expected resources to increase, got %d (initial %d)",
			len(state.Resources), initialResources)
	}

	// Ensure at least one wheat resource exists
	wheatFound := false
	for _, res := range state.Resources {
		if res.Type == "wheat" {
			wheatFound = true
			break
		}
	}
	if !wheatFound {
		t.Error("expected wheat resource to be produced")
	}
}

func TestProductionMining(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 123)

	// Add a mine building
	state.AddBuilding(Building{
		Type:       "mine",
		Location:   "iron_mine",
		Level:      1,
		Workers:    []string{},
		Production: []Production{},
		Metadata:   nil,
	})

	initialResources := len(state.Resources)

	// Run several weeks
	for week := 1; week <= 10; week++ {
		prod.Update(week, state, state.RNG.Rand())
	}

	// Expect ore resources added
	oreFound := false
	for _, res := range state.Resources {
		if res.Type == "ore" {
			oreFound = true
			break
		}
	}
	if !oreFound {
		t.Error("expected ore resource to be produced by mining")
	}
	if len(state.Resources) <= initialResources {
		t.Errorf("expected resources to increase, got %d (initial %d)",
			len(state.Resources), initialResources)
	}
}

func TestProductionCrafting(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 456)

	// Add ore resources
	state.AddResource(Resource{Type: "ore", Quantity: 10, Quality: 1.0})
	// Add a workshop building
	state.AddBuilding(Building{
		Type:       "workshop",
		Location:   "workshop1",
		Level:      1,
		Workers:    []string{},
		Production: []Production{},
		Metadata:   nil,
	})

	// Run one week
	prod.Update(1, state, state.RNG.Rand())

	// Expect tool resources added
	toolFound := false
	for _, res := range state.Resources {
		if res.Type == "tool" {
			toolFound = true
			break
		}
	}
	if !toolFound {
		t.Error("expected tool resource to be produced by crafting")
	}
	// Ore should be consumed (10 - 2 = 8 remaining)
	oreCount := 0
	for _, res := range state.Resources {
		if res.Type == "ore" {
			oreCount += res.Quantity
		}
	}
	if oreCount != 8 {
		t.Errorf("expected 8 ore remaining after crafting, got %d", oreCount)
	}
}

func TestProductionConstruction(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 789)

	// Add required materials
	state.AddResource(Resource{Type: "wood", Quantity: 10, Quality: 1.0})
	state.AddResource(Resource{Type: "stone", Quantity: 10, Quality: 1.0})
	// Add a construction site with workers
	state.AddBuilding(Building{
		Type:       "construction_site",
		Location:   "site1",
		Level:      1,
		Workers:    []string{"worker1", "worker2"},
		Production: []Production{},
		Metadata:   nil,
	})

	initialProgress := 0.0
	// Run several weeks
	for week := 1; week <= 5; week++ {
		prod.Update(week, state, state.RNG.Rand())
		// Check progress in metadata
		for _, b := range state.Buildings {
			if b.Type == "construction_site" {
				if prog, ok := b.Metadata["progress"].(float64); ok {
					initialProgress = prog
				}
			}
		}
	}
	// Expect progress > 0
	if initialProgress <= 0 {
		t.Error("expected construction progress to increase")
	}
}

func TestProductionDeterministicAcrossAllSystems(t *testing.T) {
	seed := int64(999)
	prod1 := NewProductionSystem()
	prod2 := NewProductionSystem()
	state1 := NewGameState("test1", seed)
	state2 := NewGameState("test2", seed)

	// Add same buildings to both states
	buildings := []Building{
		{Type: "farm", Location: "farm1", Level: 1, Metadata: nil},
		{Type: "mine", Location: "mine1", Level: 1, Metadata: nil},
		{Type: "workshop", Location: "workshop1", Level: 1, Metadata: nil},
		{Type: "construction_site", Location: "site1", Level: 1, Workers: []string{"w1"}, Metadata: nil},
	}
	for _, b := range buildings {
		state1.AddBuilding(b)
		state2.AddBuilding(b)
	}
	// Add some resources
	state1.AddResource(Resource{Type: "ore", Quantity: 20, Quality: 1.0})
	state1.AddResource(Resource{Type: "wood", Quantity: 20, Quality: 1.0})
	state1.AddResource(Resource{Type: "stone", Quantity: 20, Quality: 1.0})
	state2.AddResource(Resource{Type: "ore", Quantity: 20, Quality: 1.0})
	state2.AddResource(Resource{Type: "wood", Quantity: 20, Quality: 1.0})
	state2.AddResource(Resource{Type: "stone", Quantity: 20, Quality: 1.0})

	// Run 4 weeks
	for week := 1; week <= 4; week++ {
		prod1.Update(week, state1, state1.RNG.Rand())
		prod2.Update(week, state2, state2.RNG.Rand())
	}

	// Compare resource counts
	if len(state1.Resources) != len(state2.Resources) {
		t.Errorf("resource count mismatch: %d vs %d",
			len(state1.Resources), len(state2.Resources))
	}
	// Compare total quantity per type
	quant1 := make(map[string]int)
	quant2 := make(map[string]int)
	for _, res := range state1.Resources {
		quant1[res.Type] += res.Quantity
	}
	for _, res := range state2.Resources {
		quant2[res.Type] += res.Quantity
	}
	for typ, q1 := range quant1 {
		q2 := quant2[typ]
		if q1 != q2 {
			t.Errorf("resource quantity mismatch for %s: %d vs %d", typ, q1, q2)
		}
	}
}

func TestCalculateGrowthChance(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	env := Environment{
		SoilFertility: 0.8,
		Rainfall:      15.0,
		Temperature:   20.0,
	}
	chance := calculateGrowthChance(env, rng)
	if chance <= 0 || chance > 1 {
		t.Errorf("growth chance out of range: %f", chance)
	}

	// Test low fertility
	env.SoilFertility = 0.1
	chance2 := calculateGrowthChance(env, rng)
	if chance2 <= 0 || chance2 > 1 {
		t.Errorf("growth chance out of range: %f", chance2)
	}
	// chance2 should be lower than chance
	if chance2 >= chance {
		t.Errorf("expected lower chance with low fertility, got %f >= %f", chance2, chance)
	}

	// Test extreme rainfall
	env.Rainfall = 40.0
	chance3 := calculateGrowthChance(env, rng)
	if chance3 <= 0 || chance3 > 1 {
		t.Errorf("growth chance out of range: %f", chance3)
	}
}

func TestCalculateYield(t *testing.T) {
	rng := rand.New(rand.NewSource(456))
	env := Environment{
		SoilFertility: 0.7,
		Rainfall:      12.0,
		Temperature:   18.0,
	}
	yield := calculateYield(env, 2, rng)
	if yield <= 0 {
		t.Errorf("yield should be positive, got %d", yield)
	}
	// Yield should be at least 1
	if yield < 1 {
		t.Errorf("yield less than 1: %d", yield)
	}
	// Test with low fertility
	env.SoilFertility = 0.1
	yield2 := calculateYield(env, 2, rng)
	if yield2 <= 0 {
		t.Errorf("yield with low fertility should be positive, got %d", yield2)
	}
	// yield2 should be lower than yield
	if yield2 >= yield {
		t.Errorf("expected lower yield with low fertility, got %d >= %d", yield2, yield)
	}
}

func TestProductionAgricultureNoFarm(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 789)
	// No farm buildings
	initialResources := len(state.Resources)
	prod.Update(1, state, state.RNG.Rand())
	if len(state.Resources) != initialResources {
		t.Error("resources should not change without farms")
	}
}

func TestProductionMiningDepletion(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 999)
	state.AddBuilding(Building{
		Type:     "mine",
		Location: "mine1",
		Level:    1,
		Metadata: map[string]interface{}{"depletion": 0.9},
	})
	// Run one week
	prod.Update(1, state, state.RNG.Rand())
	// Check depletion updated
	for _, b := range state.Buildings {
		if b.Type == "mine" {
			dep, ok := b.Metadata["depletion"].(float64)
			if !ok {
				t.Error("depletion metadata missing")
			}
			if dep <= 0 || dep > 1 {
				t.Errorf("depletion out of range: %f", dep)
			}
		}
	}
}

func TestProductionCraftingNoOre(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 111)
	state.AddBuilding(Building{
		Type:     "workshop",
		Location: "ws1",
		Level:    1,
	})
	// No ore resources
	initialResources := len(state.Resources)
	prod.Update(1, state, state.RNG.Rand())
	if len(state.Resources) != initialResources {
		t.Error("resources should not change without ore")
	}
}

func TestProductionConstructionNoMaterials(t *testing.T) {
	prod := NewProductionSystem()
	state := NewGameState("test", 222)
	state.AddBuilding(Building{
		Type:     "construction_site",
		Location: "site1",
		Workers:  []string{"w1"},
	})
	// No wood/stone
	initialResources := len(state.Resources)
	prod.Update(1, state, state.RNG.Rand())
	if len(state.Resources) != initialResources {
		t.Error("resources should not change without materials")
	}
	// Progress should still increase due to workers
	for _, b := range state.Buildings {
		if b.Type == "construction_site" {
			if prog, ok := b.Metadata["progress"].(float64); ok && prog <= 0 {
				t.Error("progress should increase with workers even without materials")
			}
		}
	}
}

func TestProductionIntegrationWithEnvironment(t *testing.T) {
	seed := int64(777)
	state := NewGameState("integration", seed)
	// Add buildings
	state.AddBuilding(Building{Type: "farm", Location: "farm1", Level: 1})
	state.AddBuilding(Building{Type: "mine", Location: "mine1", Level: 1})
	state.AddBuilding(Building{Type: "workshop", Location: "ws1", Level: 1})
	state.AddBuilding(Building{Type: "construction_site", Location: "site1", Level: 1, Workers: []string{"w1"}})
	// Add some resources
	state.AddResource(Resource{Type: "ore", Quantity: 20, Quality: 1.0})
	state.AddResource(Resource{Type: "wood", Quantity: 20, Quality: 1.0})
	state.AddResource(Resource{Type: "stone", Quantity: 20, Quality: 1.0})

	envSys := NewEnvironmentSystem()
	prodSys := NewProductionSystem()

	// Run environment update for week 1
	envEvents := envSys.Update(state.Calendar.Week, state, state.RNG.Rand())
	_ = envEvents
	// Run production update for week 1
	prodEvents := prodSys.Update(state.Calendar.Week, state, state.RNG.Rand())
	_ = prodEvents

	// Ensure environment fields changed
	if state.Environment.Season == "" {
		t.Error("season should be set by environment update")
	}
	// Ensure some production happened (maybe)
	// At least resources changed
	if len(state.Resources) == 0 {
		t.Error("resources should exist")
	}
}

func TestProductionDeterministicWithEnvironment(t *testing.T) {
	seed := int64(888)
	state1 := NewGameState("det1", seed)
	state2 := NewGameState("det2", seed)
	// Add same buildings
	buildings := []Building{
		{Type: "farm", Location: "f1", Level: 1},
		{Type: "mine", Location: "m1", Level: 1},
		{Type: "workshop", Location: "w1", Level: 1},
		{Type: "construction_site", Location: "c1", Level: 1, Workers: []string{"worker"}},
	}
	for _, b := range buildings {
		state1.AddBuilding(b)
		state2.AddBuilding(b)
	}
	// Add same resources
	resources := []Resource{
		{Type: "ore", Quantity: 30, Quality: 1.0},
		{Type: "wood", Quantity: 30, Quality: 1.0},
		{Type: "stone", Quantity: 30, Quality: 1.0},
	}
	for _, r := range resources {
		state1.AddResource(r)
		state2.AddResource(r)
	}

	env1 := NewEnvironmentSystem()
	env2 := NewEnvironmentSystem()
	prod1 := NewProductionSystem()
	prod2 := NewProductionSystem()

	// Run 3 weeks
	for week := 1; week <= 3; week++ {
		env1.Update(week, state1, state1.RNG.Rand())
		env2.Update(week, state2, state2.RNG.Rand())
		prod1.Update(week, state1, state1.RNG.Rand())
		prod2.Update(week, state2, state2.RNG.Rand())
	}

	// Compare environment fields
	if state1.Environment.Season != state2.Environment.Season {
		t.Errorf("season mismatch: %s vs %s", state1.Environment.Season, state2.Environment.Season)
	}
	if state1.Environment.Temperature != state2.Environment.Temperature {
		t.Errorf("temperature mismatch: %f vs %f", state1.Environment.Temperature, state2.Environment.Temperature)
	}
	// Compare resource counts
	if len(state1.Resources) != len(state2.Resources) {
		t.Errorf("resource count mismatch: %d vs %d", len(state1.Resources), len(state2.Resources))
	}
	// Compare each resource type total quantity
	quant1 := make(map[string]int)
	quant2 := make(map[string]int)
	for _, res := range state1.Resources {
		quant1[res.Type] += res.Quantity
	}
	for _, res := range state2.Resources {
		quant2[res.Type] += res.Quantity
	}
	for typ, q1 := range quant1 {
		q2 := quant2[typ]
		if q1 != q2 {
			t.Errorf("resource quantity mismatch for %s: %d vs %d", typ, q1, q2)
		}
	}
}

func TestProductionEnvironmentHelpers(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	// Test calculateSeason
	if s := calculateSeason(1); s != "spring" {
		t.Errorf("week 1 should be spring, got %s", s)
	}
	if s := calculateSeason(14); s != "summer" {
		t.Errorf("week 14 should be summer, got %s", s)
	}
	if s := calculateSeason(15); s != "summer" {
		t.Errorf("week 15 should be summer, got %s", s)
	}
	if s := calculateSeason(40); s != "winter" {
		t.Errorf("week 40 should be winter, got %s", s)
	}
	if s := calculateSeason(52); s != "winter" {
		t.Errorf("week 52 should be winter, got %s", s)
	}

	// Test calculateWeather
	temp, rain := calculateWeather("spring", rng)
	if temp < -10 || temp > 40 {
		t.Errorf("spring temperature out of plausible range: %f", temp)
	}
	if rain < 0 {
		t.Errorf("rainfall negative: %f", rain)
	}

	// Test updateSoilFertility
	fertility := updateSoilFertility(0.5, 15.0, rng)
	if fertility < 0 || fertility > 1 {
		t.Errorf("soil fertility out of range: %f", fertility)
	}

	// Test regenerateResource
	resource := regenerateResource(0.3, rng, 0.01)
	if resource < 0 || resource > 1 {
		t.Errorf("regenerated resource out of range: %f", resource)
	}

	// Test updateWildlife
	wildlife := updateWildlife(0.4, 0.8, rng)
	if wildlife < 0 || wildlife > 1 {
		t.Errorf("wildlife population out of range: %f", wildlife)
	}

	// Test GenerateWeatherEvents
	events := GenerateWeatherEvents(25.0, 5.0, 1, 1)
	if len(events) != 0 {
		t.Errorf("expected no weather events, got %d", len(events))
	}
	events2 := GenerateWeatherEvents(-1.0, 25.0, 2, 1)
	if len(events2) < 2 { // freezing and heavy rain
		t.Errorf("expected at least 2 weather events, got %d", len(events2))
	}
}

func TestProductionRNGUsage(t *testing.T) {
	state := NewGameState("rngtest", 12345)
	// Use RNG methods
	n := state.RNG.Intn(100)
	if n < 0 || n >= 100 {
		t.Errorf("Intn out of range: %d", n)
	}
	f := state.RNG.Float64()
	if f < 0 || f >= 1 {
		t.Errorf("Float64 out of range: %f", f)
	}
	// Shuffle a slice
	slice := []int{1, 2, 3, 4, 5}
	state.RNG.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
	// Ensure slice still contains same elements (order may differ)
	sum := 0
	for _, v := range slice {
		sum += v
	}
	if sum != 15 {
		t.Errorf("shuffle corrupted slice")
	}
}

func TestProductionStateHelpers(t *testing.T) {
	state := NewGameState("helpertest", 999)
	// Add resident
	state.AddResident(Resident{ID: "r1", Name: "Test"})
	if len(state.Residents) != 1 {
		t.Error("AddResident failed")
	}
	// Add event
	state.AddEvent(Event{ID: "e1", Type: "test"})
	if len(state.History) != 1 {
		t.Error("AddEvent failed")
	}
	// Add policy
	state.AddPolicy(Policy{ID: "p1", Name: "Test Policy"})
	if len(state.Policies) != 1 {
		t.Error("AddPolicy failed")
	}
}

func TestProductionEnvironmentEdgeCases(t *testing.T) {
	rng := rand.New(rand.NewSource(111))
	// calculateWeather for all seasons
	seasons := []string{"spring", "summer", "autumn", "winter", "unknown"}
	for _, season := range seasons {
		temp, rain := calculateWeather(season, rng)
		if rain < 0 {
			t.Errorf("negative rain for season %s: %f", season, rain)
		}
		// temperature could be negative in winter, but we just check plausible range
		if temp < -30 || temp > 50 {
			t.Errorf("implausible temperature for season %s: %f", season, temp)
		}
	}

	// updateSoilFertility with various rainfall values
	testCases := []float64{0, 2, 5, 10, 15, 25, 35, 50}
	for _, rain := range testCases {
		fertility := updateSoilFertility(0.5, rain, rng)
		if fertility < 0 || fertility > 1 {
			t.Errorf("fertility out of range for rain %f: %f", rain, fertility)
		}
	}

	// regenerateResource with extreme current values
	vals := []float64{0.0, 0.5, 1.0}
	for _, cur := range vals {
		res := regenerateResource(cur, rng, 0.01)
		if res < 0 || res > 1 {
			t.Errorf("regenerated resource out of range for cur %f: %f", cur, res)
		}
	}

	// updateWildlife with forest health extremes
	for _, forest := range []float64{0.0, 0.2, 0.5, 0.8, 1.0} {
		wildlife := updateWildlife(0.5, forest, rng)
		if wildlife < 0 || wildlife > 1 {
			t.Errorf("wildlife out of range for forest %f: %f", forest, wildlife)
		}
	}
}

func TestProductionCoreStub(t *testing.T) {
	// Ensure core stub works (no panic)
	core := NewCore()
	if core == nil {
		t.Error("NewCore returned nil")
	}
	err := core.ProcessWeek()
	if err != nil {
		t.Errorf("ProcessWeek returned error: %v", err)
	}
}

func TestProductionSerialization(t *testing.T) {
	// Test SaveRNG and LoadRNG (basic)
	state := NewGameState("serial", 123)
	// SaveRNG returns JSON bytes
	jsonBytes, err := SaveRNG(state.RNG)
	if err != nil {
		t.Errorf("SaveRNG error: %v", err)
	}
	if len(jsonBytes) == 0 {
		t.Error("SaveRNG returned empty bytes")
	}
	// LoadRNG should create a new RNG with same seed
	rng2, err := LoadRNG(jsonBytes)
	if err != nil {
		t.Errorf("LoadRNG error: %v", err)
	}
	if rng2 == nil {
		t.Error("LoadRNG returned nil")
	}
	// Both RNGs should produce same sequence
	n1 := state.RNG.Uint64()
	n2 := rng2.Uint64()
	if n1 != n2 {
		t.Error("RNG mismatch after load")
	}
}

func TestProductionTurnProcessor(t *testing.T) {
	tp := NewTurnProcessor()
	if tp == nil {
		t.Fatal("NewTurnProcessor returned nil")
	}
	// Register production system
	prod := NewProductionSystem()
	tp.RegisterSystem(SystemProduction, prod)
	// Create game state
	state := NewGameState("turn-test", 12345)
	// Add some buildings and resources
	state.AddBuilding(Building{Type: "farm", Location: "f1", Level: 1})
	state.AddResource(Resource{Type: "ore", Quantity: 10, Quality: 1.0})
	state.AddResource(Resource{Type: "wood", Quantity: 10, Quality: 1.0})
	state.AddResource(Resource{Type: "stone", Quantity: 10, Quality: 1.0})

	// Process a week
	events := tp.ProcessWeek(state)
	_ = events
	// Calendar should advance
	if state.Calendar.Week != 2 {
		t.Errorf("expected week 2, got %d", state.Calendar.Week)
	}
	// Resources may have changed
	if len(state.Resources) == 0 {
		t.Error("resources should exist")
	}
}
