package economy

import (
	"errors"
	"github.com/vano44/village/internal/simulation"
)

// Inventory manages resources across storage locations.
type Inventory struct {
	// location -> list of resources stored there
	resources map[string][]Resource
	// optional storage registry for capacity and type enforcement
	storage *StorageRegistry
}

// Alert represents a stock alert when quantity falls below threshold.
type Alert struct {
	Location  string
	Resource  ResourceType
	Current   float64
	Threshold float64
}

// AlertThreshold maps resource type to minimum quantity that triggers an alert.
type AlertThreshold map[ResourceType]float64

// NewInventory creates an empty inventory.
func NewInventory() *Inventory {
	return &Inventory{
		resources: make(map[string][]Resource),
	}
}

// NewInventoryWithStorage creates an empty inventory with a storage registry.
func NewInventoryWithStorage(storage *StorageRegistry) *Inventory {
	return &Inventory{
		resources: make(map[string][]Resource),
		storage:   storage,
	}
}

// SetStorage assigns a storage registry to the inventory.
func (inv *Inventory) SetStorage(storage *StorageRegistry) {
	inv.storage = storage
}

// CheckAlerts returns alerts for resources whose total quantity across all locations falls below threshold.
func (inv *Inventory) CheckAlerts(thresholds AlertThreshold) []Alert {
	var alerts []Alert
	for rt, threshold := range thresholds {
		total := 0.0
		for _, list := range inv.resources {
			for _, r := range list {
				if r.Type == rt {
					total += r.Quantity
				}
			}
		}
		if total < threshold {
			alerts = append(alerts, Alert{
				Location:  "", // empty indicates global total
				Resource:  rt,
				Current:   total,
				Threshold: threshold,
			})
		}
	}
	return alerts
}

// AddResource adds a resource to a location.
// If a storage registry is attached, it validates capacity and storage type.
func (inv *Inventory) AddResource(location string, r Resource) error {
	if inv.storage != nil {
		// Check if resource type can be stored at this location
		if !inv.storage.CanStoreResourceAt(location, r.Type) {
			return errors.New("resource type cannot be stored at this location")
		}
		// Check capacity
		available := inv.storage.AvailableCapacity(location, inv)
		if r.Quantity > available {
			return errors.New("insufficient storage capacity")
		}
	}
	inv.resources[location] = append(inv.resources[location], r)
	return nil
}

// GetAvailable returns the total quantity of a resource type at a location.
func (inv *Inventory) GetAvailable(location string, rt ResourceType) float64 {
	list, ok := inv.resources[location]
	if !ok {
		return 0.0
	}
	total := 0.0
	for _, r := range list {
		if r.Type == rt {
			total += r.Quantity
		}
	}
	return total
}

// RemoveResource removes the requested quantity of a resource type from a location.
// Returns the amount removed (equal to quantity) or an error if insufficient.
func (inv *Inventory) RemoveResource(location string, rt ResourceType, quantity float64) (float64, error) {
	available := inv.GetAvailable(location, rt)
	if available < quantity {
		return 0.0, errors.New("insufficient quantity")
	}
	list, ok := inv.resources[location]
	if !ok {
		return 0.0, nil // should not happen because available >= quantity > 0
	}
	remaining := quantity
	var newList []Resource
	for _, r := range list {
		if r.Type == rt && remaining > 0 {
			if r.Quantity <= remaining {
				remaining -= r.Quantity
				// skip this resource (fully consumed)
				continue
			} else {
				// partially consume this resource
				r.Quantity -= remaining
				newList = append(newList, r)
				remaining = 0
			}
		} else {
			newList = append(newList, r)
		}
	}
	inv.resources[location] = newList
	return quantity - remaining, nil
}

// TransferResource moves quantity of a resource type from source to destination location.
func (inv *Inventory) TransferResource(src, dst string, rt ResourceType, quantity float64) (float64, error) {
	// Remove from source
	removed, err := inv.RemoveResource(src, rt, quantity)
	if err != nil {
		return 0.0, err
	}
	if removed == 0.0 {
		return 0.0, nil
	}
	// Add to destination as a new resource with default quality and produced date
	inv.AddResource(dst, NewResource(rt, removed))
	return removed, nil
}

// LoadInventoryFromGameState imports simulation resources into inventory, assigning them to defaultLocation.
func LoadInventoryFromGameState(inv *Inventory, resources []simulation.Resource, defaultLocation string) error {
	for _, sr := range resources {
		r := FromSimulationResource(sr)
		r.Location = defaultLocation
		if err := inv.AddResource(defaultLocation, r); err != nil {
			return err
		}
	}
	return nil
}

// ExportInventoryToGameState exports all inventory resources as simulation resources (location information lost).
func ExportInventoryToGameState(inv *Inventory) []simulation.Resource {
	var out []simulation.Resource
	for _, list := range inv.resources {
		for _, r := range list {
			out = append(out, r.ToSimulationResource())
		}
	}
	return out
}
