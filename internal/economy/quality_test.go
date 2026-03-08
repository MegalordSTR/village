package economy

import (
	"math/rand"
	"testing"
)

func TestQualityString(t *testing.T) {
	tests := []struct {
		q    QualityTier
		want string
	}{
		{QualityPoor, "Poor"},
		{QualityNormal, "Normal"},
		{QualityGood, "Good"},
		{QualityExcellent, "Excellent"},
		{QualityMasterwork, "Masterwork"},
	}
	for _, tt := range tests {
		got := tt.q.String()
		if got != tt.want {
			t.Errorf("QualityTier(%d).String() = %q, want %q", tt.q, got, tt.want)
		}
	}
}

func TestDurabilityMultiplier(t *testing.T) {
	tests := []struct {
		q    QualityTier
		want float64
	}{
		{QualityPoor, 0.5},
		{QualityNormal, 1.0},
		{QualityGood, 1.5},
		{QualityExcellent, 2.0},
		{QualityMasterwork, 3.0},
	}
	for _, tt := range tests {
		got := DurabilityMultiplier(tt.q)
		if got != tt.want {
			t.Errorf("DurabilityMultiplier(%v) = %v, want %v", tt.q, got, tt.want)
		}
	}
}

func TestProductionSpeedMultiplier(t *testing.T) {
	tests := []struct {
		q    QualityTier
		want float64
	}{
		{QualityPoor, 0.8},
		{QualityNormal, 1.0},
		{QualityGood, 1.1},
		{QualityExcellent, 1.2},
		{QualityMasterwork, 1.3},
	}
	for _, tt := range tests {
		got := ProductionSpeedMultiplier(tt.q)
		if got != tt.want {
			t.Errorf("ProductionSpeedMultiplier(%v) = %v, want %v", tt.q, got, tt.want)
		}
	}
}

func TestSpoilageResistanceMultiplier(t *testing.T) {
	tests := []struct {
		q    QualityTier
		want float64
	}{
		{QualityPoor, 1.0},
		{QualityNormal, 1.0},
		{QualityGood, 0.7},
		{QualityExcellent, 0.5},
		{QualityMasterwork, 0.3},
	}
	for _, tt := range tests {
		got := SpoilageResistanceMultiplier(tt.q)
		if got != tt.want {
			t.Errorf("SpoilageResistanceMultiplier(%v) = %v, want %v", tt.q, got, tt.want)
		}
	}
}

func TestMaterialQualityBonus(t *testing.T) {
	tests := []struct {
		rt   ResourceType
		want float64
	}{
		{ResourceIron, 0.2},
		{ResourceWood, 0.1},
		{ResourceGrain, 0.0},
		{ResourceStone, 0.0},
	}
	for _, tt := range tests {
		got := MaterialQualityBonus(tt.rt)
		if got != tt.want {
			t.Errorf("MaterialQualityBonus(%v) = %v, want %v", tt.rt, got, tt.want)
		}
	}
}

func TestSpecialMaterialInheritance(t *testing.T) {
	// Test that iron provides a quality bonus in production.
	// Use the forge_tools recipe with iron input.
	ctx := ProductionContext{
		RecipeID: "forge_tools",
		InputResources: map[ResourceType]Resource{
			ResourceIron: {
				Type:     ResourceIron,
				Quantity: 10,
				Quality:  QualityNormal, // 0.25
			},
		},
		WorkerSkill:       0.5,
		BuildingCondition: 1.0,
		RandSeed:          42,
	}
	result := CalculateProduction(ctx)
	if !result.Success {
		t.Fatal("Production should succeed")
	}
	if len(result.OutputResources) != 1 {
		t.Fatal("Expected one output")
	}
	outQuality := result.OutputResources[0].Quality
	// Without bonus, quality would be Normal? Let's compute roughly.
	// Input quality Normal (0.25) + bonus 0.2 = 0.45 effective.
	// Inheritance factor 0.9 => inherited = 0.45 * 0.9 = 0.405
	// skill contrib = 0.5 * (1-0.9) = 0.05
	// final = 0.455 -> QualityGood (0.5)
	expected := QualityGood
	if outQuality != expected {
		t.Errorf("Output quality with iron bonus = %v, want %v", outQuality, expected)
	}
}

func TestQualityFromSkill(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	tests := []struct {
		skill float64
		want  QualityTier
	}{
		{0.0, QualityPoor},
		{0.1, QualityPoor},
		{0.2, QualityNormal},
		{0.3, QualityNormal},
		{0.4, QualityGood},
		{0.5, QualityGood},
		{0.6, QualityExcellent},
		{0.7, QualityExcellent},
		{0.8, QualityMasterwork},
		{0.9, QualityMasterwork},
		{1.0, QualityMasterwork},
	}
	for _, tt := range tests {
		got := QualityFromSkill(tt.skill, rng)
		if got != tt.want {
			t.Errorf("QualityFromSkill(%v) = %v, want %v", tt.skill, got, tt.want)
		}
	}
}

func TestSkillDefinitions(t *testing.T) {
	// Ensure skill IDs exist
	skillIDs := []string{SkillFarming, SkillMining, SkillSmithing, SkillWeaving, SkillBaking, SkillMilling, SkillCarpentry, SkillTailoring}
	for _, id := range skillIDs {
		if !IsValidSkill(id) {
			t.Errorf("Skill ID %q should be valid", id)
		}
	}
	// Invalid skill ID
	if IsValidSkill("nonexistent") {
		t.Errorf("Invalid skill ID should return false")
	}
}

func TestSkillProgression(t *testing.T) {
	skill := Skill{ID: SkillSmithing, Level: 0.2, XP: 0}
	// Practice increases XP
	skill.Practice(50.0)
	if skill.XP != 50.0 {
		t.Errorf("Practice should increase XP, got %v want 50", skill.XP)
	}
	// Level unchanged before UpdateLevel
	if skill.Level != 0.2 {
		t.Errorf("Level should not change before UpdateLevel")
	}
	// UpdateLevel recomputes level
	skill.UpdateLevel()
	// XP=50 => level=0.5
	if skill.Level != 0.5 {
		t.Errorf("UpdateLevel should set Level based on XP, got %v want 0.5", skill.Level)
	}
	// Cap at 1.0
	skill.XP = 200.0
	skill.UpdateLevel()
	if skill.Level != 1.0 {
		t.Errorf("Level should cap at 1.0, got %v", skill.Level)
	}
	// Negative practice does nothing
	skill.XP = 100.0
	skill.Practice(-10.0)
	if skill.XP != 100.0 {
		t.Errorf("Negative practice should not change XP")
	}
}
