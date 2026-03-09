package economy

import "errors"

// StorageType represents a type of storage building.
type StorageType string

const (
	StorageGranary     StorageType = "granary"
	StorageWarehouse   StorageType = "warehouse"
	StorageOutdoorPile StorageType = "outdoor_pile"
)

// StorageQuality represents the quality of storage conditions.
type StorageQuality string

const (
	StorageQualityGood   StorageQuality = "good"
	StorageQualityNormal StorageQuality = "normal"
	StorageQualityPoor   StorageQuality = "poor"
)

// StorageBuilding represents a physical storage building with capacity.
type StorageBuilding struct {
	ID      string         // unique identifier (matches location)
	Type    StorageType    // granary, warehouse, outdoor_pile
	Level   int            // building level (affects capacity)
	Quality StorageQuality // storage condition affecting spoilage
}

// Capacity returns the maximum amount this storage can hold.
// Units: kg for granary (food), generic units for warehouse, units for outdoor pile.
func (sb *StorageBuilding) Capacity() float64 {
	base := 0.0
	switch sb.Type {
	case StorageGranary:
		base = 1000.0 // kg
	case StorageWarehouse:
		base = 500.0 // units
	case StorageOutdoorPile:
		base = 200.0 // units
	default:
		return 0.0
	}
	// Level multiplies capacity (level 1 = base, level 2 = 2*base, etc)
	return base * float64(sb.Level)
}

// SpoilageMultiplier returns the multiplier for spoilage rate based on storage quality.
func (sb *StorageBuilding) SpoilageMultiplier() float64 {
	switch sb.Quality {
	case StorageQualityGood:
		return 0.5
	case StorageQualityNormal:
		return 1.0
	case StorageQualityPoor:
		return 2.0
	default:
		return 1.0
	}
}

// CanStoreResource checks if a resource can be stored in this storage type.
// Granary stores food resources (raw & processed food), warehouse stores goods, outdoor pile stores anything.
func (sb *StorageBuilding) CanStoreResource(rt ResourceType) bool {
	switch sb.Type {
	case StorageGranary:
		// Food resources: grain, vegetables, flour, bread
		return rt == ResourceGrain || rt == ResourceVegetables || rt == ResourceFlour || rt == ResourceBread
	case StorageWarehouse:
		// Goods: processed and advanced resources, but not raw materials (wood, stone, etc)
		category := CategoryForResource(rt)
		return category == CategoryProcessed || category == CategoryAdvanced
	case StorageOutdoorPile:
		// Anything can be stored outdoors (but spoilage high)
		return true
	default:
		return false
	}
}

// StorageRegistry maintains a collection of storage buildings.
type StorageRegistry struct {
	buildings map[string]*StorageBuilding // location ID -> building
}

// NewStorageRegistry creates an empty registry.
func NewStorageRegistry() *StorageRegistry {
	return &StorageRegistry{
		buildings: make(map[string]*StorageBuilding),
	}
}

// AddBuilding adds a storage building to the registry.
func (sr *StorageRegistry) AddBuilding(b *StorageBuilding) error {
	if b.ID == "" {
		return errors.New("storage building must have an ID")
	}
	sr.buildings[b.ID] = b
	return nil
}

// GetBuilding returns the storage building for a location, or nil if not found.
func (sr *StorageRegistry) GetBuilding(location string) *StorageBuilding {
	return sr.buildings[location]
}

// Capacity returns total capacity for a location (0 if building not found).
func (sr *StorageRegistry) Capacity(location string) float64 {
	b := sr.GetBuilding(location)
	if b == nil {
		return 1e9 // unlimited capacity for locations without storage building
	}
	return b.Capacity()
}

// CanStoreResourceAt checks if the resource can be stored at the given location.
func (sr *StorageRegistry) CanStoreResourceAt(location string, rt ResourceType) bool {
	b := sr.GetBuilding(location)
	if b == nil {
		return true // allow storage anywhere without building
	}
	return b.CanStoreResource(rt)
}

// CurrentOccupancy calculates the total quantity of resources stored at a location.
// Requires inventory to compute occupancy.
func (sr *StorageRegistry) CurrentOccupancy(location string, inv *Inventory) float64 {
	total := 0.0
	list, ok := inv.resources[location]
	if !ok {
		return 0.0
	}
	for _, r := range list {
		total += r.Quantity
	}
	return total
}

// AvailableCapacity returns remaining capacity at location.
func (sr *StorageRegistry) AvailableCapacity(location string, inv *Inventory) float64 {
	return sr.Capacity(location) - sr.CurrentOccupancy(location, inv)
}
