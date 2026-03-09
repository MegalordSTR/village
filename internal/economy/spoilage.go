package economy

// WeeksSince returns the number of weeks elapsed between two GameDates.
// If the produced date is after current date, returns 0.
func WeeksSince(current, produced GameDate) int {
	if current.Year < produced.Year || (current.Year == produced.Year && current.Week < produced.Week) {
		return 0
	}
	return (current.Year-produced.Year)*52 + (current.Week - produced.Week)
}

// ApplySpoilage calculates spoilage for a resource given storage conditions and current date.
// Returns the amount spoiled this call (added to resource.Spoiled) and whether quality was degraded.
// Modifies the resource in place: reduces Quantity, increases Spoiled, may lower Quality.
func ApplySpoilage(resource *Resource, storage *StorageBuilding, currentDate GameDate) (spoiledQty float64, degraded bool) {
	if resource.Quantity <= 0 {
		return 0.0, false
	}
	baseRate := SpoilageRate(resource.Type)
	if baseRate <= 0 {
		return 0.0, false
	}
	ageWeeks := WeeksSince(currentDate, resource.Produced)
	if ageWeeks <= 0 {
		return 0.0, false
	}
	// Total spoilage rate = base * storage multiplier
	storageMult := storage.SpoilageMultiplier()
	totalRate := baseRate * storageMult
	// Approximate spoilage: linear decay over age (simplified)
	// Spoiled fraction = totalRate * ageWeeks (capped at 1.0)
	spoiledFraction := totalRate * float64(ageWeeks)
	if spoiledFraction > 1.0 {
		spoiledFraction = 1.0
	}
	originalQty := resource.Quantity + resource.Spoiled
	spoiledQty = originalQty * spoiledFraction
	// Subtract already spoiled amount (resource.Spoiled) from total spoiled to get new spoilage
	newSpoiled := spoiledQty - resource.Spoiled
	if newSpoiled <= 0 {
		return 0.0, false
	}
	// Update resource
	resource.Quantity -= newSpoiled
	if resource.Quantity < 0 {
		resource.Quantity = 0
	}
	resource.Spoiled = spoiledQty

	// Check quality degradation threshold (20% spoilage)
	threshold := 0.2
	if resource.Spoiled/originalQty >= threshold {
		// Degrade quality by one tier, if not already poorest
		if resource.Quality > QualityPoor {
			resource.Quality--
			degraded = true
		}
	}
	return newSpoiled, degraded
}

// ApplySpoilageToInventory applies spoilage to all resources in inventory.
// Returns total spoiled quantity per resource type (or location).
func ApplySpoilageToInventory(inv *Inventory, storageReg *StorageRegistry, currentDate GameDate) map[string]float64 {
	totals := make(map[string]float64)
	for location, resources := range inv.resources {
		storage := storageReg.GetBuilding(location)
		if storage == nil {
			continue // no storage building, no spoilage? maybe default outdoor pile?
			// TODO apply additional spoilage if resources stored without storage, it must be much more than in proper storage building
		}
		for i := range resources {
			spoiled, _ := ApplySpoilage(&resources[i], storage, currentDate)
			if spoiled > 0 {
				rt := string(resources[i].Type)
				totals[rt] += spoiled
			}
		}
		// Update the slice in inventory (since we modified elements in place)
		inv.resources[location] = resources
	}
	return totals
}
