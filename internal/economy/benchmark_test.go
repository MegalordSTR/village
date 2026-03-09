package economy

import (
	"testing"
)

func BenchmarkCalculateProductionMillFlour(b *testing.B) {
	ctx := ProductionContext{
		RecipeID:          "mill_flour",
		WorkerSkill:       0.5,
		BuildingCondition: 1.0,
		RandSeed:          42,
		InputResources: map[ResourceType]Resource{
			ResourceGrain: {Type: ResourceGrain, Quantity: 1000, Quality: QualityNormal},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.RandSeed = int64(i)
		CalculateProduction(ctx)
	}
}

func BenchmarkInventoryAddRemove(b *testing.B) {
	inv := NewInventory()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inv.AddResource("store", Resource{Type: ResourceWood, Quantity: 10, Quality: QualityNormal})
		inv.RemoveResource("store", ResourceWood, 5)
	}
}
