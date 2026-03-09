package economy

import (
	"math"
	"math/rand"
	"testing"
)

func TestScarcity(t *testing.T) {
	// Run all scarcity subtests
	t.Run("MineDepletion", TestScarcityMineDepletion)
	t.Run("ForestRegrowth", TestScarcityForestRegrowth)
	t.Run("MineYieldModifier", TestScarcityMineYieldModifier)
	t.Run("ForestLogging", TestScarcityForestLogging)
	t.Run("ForestHealthBounds", TestScarcityForestHealthBounds)
	t.Run("Substitution", TestScarcitySubstitution)
	t.Run("Scavenge", TestScarcityScavenge)
	t.Run("StrategicTradeOff", TestScarcityStrategicTradeOff)
	t.Run("ResourceDiscovery", TestScarcityResourceDiscovery)
	t.Run("ApplyScarcityModifier", TestScarcityApplyScarcityModifier)
	t.Run("Integration", TestScarcityIntegration)
}

func TestScarcityMineDepletion(t *testing.T) {
	m := Mine{
		TotalOre:      1000,
		RemainingOre:  1000,
		DepletionRate: 0.1,
	}
	m.Deplete(100)
	if m.RemainingOre != 900 {
		t.Errorf("expected remaining ore 900, got %v", m.RemainingOre)
	}
}

func TestScarcityForestRegrowth(t *testing.T) {
	f := Forest{
		Health:       0.5,
		RegrowthRate: 0.1,
	}
	f.Regrow(1)
	if f.Health != 0.6 {
		t.Errorf("expected health 0.6, got %v", f.Health)
	}
}

func TestScarcityMineYieldModifier(t *testing.T) {
	m := Mine{
		TotalOre:     1000,
		RemainingOre: 1000,
	}
	// full mine -> rich multiplier 1.5
	got := m.YieldModifier()
	if math.Abs(got-1.5) > 1e-9 {
		t.Errorf("full mine: expected 1.5, got %v", got)
	}
	m.RemainingOre = 0
	got = m.YieldModifier()
	if math.Abs(got-0.3) > 1e-9 {
		t.Errorf("depleted mine: expected 0.3, got %v", got)
	}
	m.RemainingOre = 500
	got = m.YieldModifier()
	expected := 0.3 + 1.2*(500.0/1000.0)
	if math.Abs(got-expected) > 1e-9 {
		t.Errorf("half mine: expected %v, got %v", expected, got)
	}
}

func TestScarcityForestLogging(t *testing.T) {
	f := Forest{
		Health:       0.8,
		RegrowthRate: 0.1,
	}
	f.Log(0.3)
	if math.Abs(f.Health-0.5) > 1e-9 {
		t.Errorf("after logging: expected health 0.5, got %v", f.Health)
	}
	// logging cannot reduce health below 0
	f.Log(1.0)
	if f.Health != 0.0 {
		t.Errorf("health should clamp at 0, got %v", f.Health)
	}
}

func TestScarcityForestHealthBounds(t *testing.T) {
	f := Forest{
		Health:       0.9,
		RegrowthRate: 0.2,
	}
	f.Regrow(1) // should become 1.1 but clamped to 1.0
	if f.Health != 1.0 {
		t.Errorf("regrowth should clamp at 1.0, got %v", f.Health)
	}
	f.Log(2.0) // should go to -1.0 but clamped to 0.0
	if f.Health != 0.0 {
		t.Errorf("logging should clamp at 0.0, got %v", f.Health)
	}
}

func TestScarcitySubstitution(t *testing.T) {
	ResetSubstitutionRules()
	// Register substitution for wood -> stone
	RegisterSubstitution(ResourceWood, Alternative{
		Type:              ResourceStone,
		TimeMultiplier:    2.0,
		QualityMultiplier: 0.8,
	})
	alts := GetAlternatives(ResourceWood)
	if len(alts) != 1 {
		t.Fatalf("expected 1 alternative, got %d", len(alts))
	}
	alt := alts[0]
	if alt.Type != ResourceStone {
		t.Errorf("expected alternative type stone, got %v", alt.Type)
	}
	if math.Abs(alt.TimeMultiplier-2.0) > 1e-9 {
		t.Errorf("expected time multiplier 2.0, got %v", alt.TimeMultiplier)
	}
	if math.Abs(alt.QualityMultiplier-0.8) > 1e-9 {
		t.Errorf("expected quality multiplier 0.8, got %v", alt.QualityMultiplier)
	}
	// No alternatives for unknown resource
	alts2 := GetAlternatives(ResourceIron)
	if len(alts2) != 0 {
		t.Errorf("expected nil or empty slice for unknown resource, got %v", alts2)
	}
}

func TestScarcityScavenge(t *testing.T) {
	// Test abandoned mine yields iron ore and/or stone
	resources := Scavenge("abandoned_mine", 1.0)
	if len(resources) == 0 {
		t.Error("expected at least one resource from abandoned mine, got none")
	}
	// Check that all resources are valid types
	for _, r := range resources {
		if !IsValidType(r.Type) {
			t.Errorf("invalid resource type %v", r.Type)
		}
	}
}

func TestScarcityStrategicTradeOff(t *testing.T) {
	// Test recommendation logic
	// currentModifier low (abundant), futureNeed low -> use now
	if !RecommendUse(0.3, 1.0) {
		t.Error("expected use now for low scarcity")
	}
	// currentModifier high (scarce), futureNeed high -> save
	if RecommendUse(1.5, 2.0) {
		t.Error("expected save for high scarcity and high future need")
	}
	// threshold case: 0.5*1.0 = 0.5 -> should return false (save)
	if RecommendUse(0.5, 1.0) {
		t.Error("expected save at threshold")
	}
}

func TestScarcityResourceDiscovery(t *testing.T) {
	rng := rand.New(rand.NewSource(42)) // #nosec G404
	// Discover mine vein
	mine := DiscoverResource("mine_vein", rng)
	if mine == nil {
		t.Fatal("expected mine, got nil")
	}
	m, ok := mine.(Mine)
	if !ok {
		t.Fatal("expected Mine type")
	}
	if m.TotalOre < 100 || m.TotalOre > 1000 {
		t.Errorf("total ore out of expected range: %v", m.TotalOre)
	}
	if m.RemainingOre != m.TotalOre {
		t.Errorf("remaining ore should equal total ore initially")
	}
	if m.DepletionRate < 0.05 || m.DepletionRate > 0.2 {
		t.Errorf("depletion rate out of expected range: %v", m.DepletionRate)
	}
	// Discover forest patch
	forest := DiscoverResource("forest_patch", rng)
	if forest == nil {
		t.Fatal("expected forest, got nil")
	}
	f, ok := forest.(Forest)
	if !ok {
		t.Fatal("expected Forest type")
	}
	if f.Health < 0.2 || f.Health > 1.0 {
		t.Errorf("health out of expected range: %v", f.Health)
	}
	if f.RegrowthRate < 0.05 || f.RegrowthRate > 0.15 {
		t.Errorf("regrowth rate out of expected range: %v", f.RegrowthRate)
	}
	// Unknown location returns nil
	if DiscoverResource("unknown", rng) != nil {
		t.Error("expected nil for unknown location")
	}
}

func TestScarcityApplyScarcityModifier(t *testing.T) {
	// Test that modifier multiplies yield correctly
	if got := ApplyScarcityModifier(10.0, 1.5); got != 15.0 {
		t.Errorf("expected 15.0, got %v", got)
	}
	if got := ApplyScarcityModifier(10.0, 0.3); got != 3.0 {
		t.Errorf("expected 3.0, got %v", got)
	}
	if got := ApplyScarcityModifier(0.0, 5.0); got != 0.0 {
		t.Errorf("zero yield should stay zero")
	}
}

func TestScarcityIntegration(t *testing.T) {
	// High-level integration test covering multiple scarcity aspects
	// Mine depletion affects yield modifier
	m := Mine{TotalOre: 1000, RemainingOre: 1000}
	if m.YieldModifier() != 1.5 {
		t.Error("full mine should have yield modifier 1.5")
	}
	m.Deplete(500)
	if m.RemainingOre != 500 {
		t.Error("depletion not working")
	}
	// Forest logging and regrowth
	f := Forest{Health: 0.5, RegrowthRate: 0.1}
	f.Regrow(1)
	if f.Health != 0.6 {
		t.Error("regrowth not working")
	}
	// Substitution
	ResetSubstitutionRules()
	RegisterSubstitution(ResourceWood, Alternative{ResourceStone, 2.0, 0.8})
	alts := GetAlternatives(ResourceWood)
	if len(alts) != 1 || alts[0].Type != ResourceStone {
		t.Error("substitution not working")
	}
	// Scavenge
	res := Scavenge("abandoned_mine", 1.0)
	if len(res) == 0 {
		t.Error("scavenge not working")
	}
	// Resource discovery
	rng := rand.New(rand.NewSource(42)) // #nosec G404
	_ = DiscoverResource("mine_vein", rng)
	// No assertion, just ensure no panic
	// Strategic trade-off
	if !RecommendUse(0.3, 1.0) {
		t.Error("recommend use logic error")
	}
	// Scarcity modifier application
	if ApplyScarcityModifier(100.0, 0.5) != 50.0 {
		t.Error("apply scarcity modifier error")
	}
}
