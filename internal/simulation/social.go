package simulation

import (
	"fmt"
	"github.com/vano44/village/internal/economy"
	"math/rand"
)

// SocialSystem implements the social simulation system.
// It handles resident needs, skill development, relationships, and life events.
type SocialSystem struct{}

// NewSocialSystem creates a new social system.
func NewSocialSystem() *SocialSystem {
	return &SocialSystem{}
}

// Update processes one week of social simulation.
// It updates resident needs, skills, relationships, and handles life events.
func (s *SocialSystem) Update(week int, state *GameState, rng *rand.Rand) []Event {
	var events []Event

	// Update each resident
	for i := range state.Residents {
		s.updateResident(week, &state.Residents[i], state, rng, &events)
	}

	// Handle relationship changes between residents
	s.updateRelationships(state, rng, &events)

	// Process life events (births, deaths)
	s.processLifeEvents(week, state, rng, &events)

	return events
}

// updateResident updates a single resident's needs and skills.
func (s *SocialSystem) updateResident(week int, resident *Resident, state *GameState, rng *rand.Rand, events *[]Event) {
	// Ensure needs exist
	if len(resident.Needs) == 0 {
		s.initializeNeeds(resident)
	}

	// Update needs based on work, shelter, food availability
	s.updateNeeds(resident, state, rng)

	// Improve skills through work experience
	s.improveSkills(resident, state, rng)

	// Age resident (1 week older)
	resident.Age++
}

// initializeNeeds sets up default needs for a resident.
func (s *SocialSystem) initializeNeeds(resident *Resident) {
	resident.Needs = []Need{
		{ID: "hunger", Name: "Hunger", Level: 0.3},
		{ID: "warmth", Name: "Warmth", Level: 0.5},
		{ID: "happiness", Name: "Happiness", Level: 0.7},
		{ID: "social", Name: "Social", Level: 0.4},
		{ID: "shelter", Name: "Shelter", Level: 0.6},
	}
}

// updateNeeds adjusts need levels based on current conditions.
func (s *SocialSystem) updateNeeds(resident *Resident, state *GameState, rng *rand.Rand) {
	for i := range resident.Needs {
		need := &resident.Needs[i]

		// Basic need decay
		decay := 0.01 + rng.Float64()*0.02
		need.Level += decay

		// Apply modifiers based on state
		switch need.ID {
		case "hunger":
			// Food availability reduces hunger
			foodAvailable := s.getFoodAvailability(state)
			need.Level -= foodAvailable * 0.05
		case "warmth":
			// Temperature affects warmth
			temp := state.Environment.Temperature
			if temp < 10 {
				need.Level += 0.03 // colder increases need
			} else if temp > 25 {
				need.Level -= 0.02 // warmer reduces need
			}
		case "happiness":
			// Influenced by other needs, relationships, work, and shelter
			avgNeed := s.averageNeedLevel(resident)
			relStrength := s.relationshipStrength(resident, state)
			need.Level += (0.5 - avgNeed) * 0.1
			need.Level += relStrength * 0.05
			if s.hasWork(resident, state) {
				need.Level += 0.05 // work increases happiness
			}
			if s.hasShelter(resident, state) {
				need.Level += 0.03 // shelter increases happiness
			}
		case "social":
			// More relationships reduce social need
			relCount := len(resident.Relationships)
			need.Level -= float64(relCount) * 0.02
		case "shelter":
			// Shelter need decreases if resident has a house
			if s.hasShelter(resident, state) {
				need.Level -= 0.1
			}
		}

		// Clamp to [0, 1]
		if need.Level < 0 {
			need.Level = 0
		}
		if need.Level > 1 {
			need.Level = 1
		}
	}
}

// getFoodAvailability estimates food availability from resources.
func (s *SocialSystem) getFoodAvailability(state *GameState) float64 {
	foodCount := int(state.Inventory.GetAvailable("global", economy.ResourceGrain))
	// Convert to per-resident availability
	residentCount := len(state.Residents)
	if residentCount == 0 {
		return 1.0
	}
	availability := float64(foodCount) / float64(residentCount*10)
	if availability > 1 {
		return 1
	}
	return availability
}

// averageNeedLevel calculates average need level for a resident.
func (s *SocialSystem) averageNeedLevel(resident *Resident) float64 {
	if len(resident.Needs) == 0 {
		return 0
	}
	sum := 0.0
	for _, n := range resident.Needs {
		sum += n.Level
	}
	return sum / float64(len(resident.Needs))
}

// relationshipStrength calculates total relationship strength for a resident.
func (s *SocialSystem) relationshipStrength(resident *Resident, state *GameState) float64 {
	total := 0.0
	for _, rel := range resident.Relationships {
		total += rel.Strength
	}
	return total
}

// hasShelter checks if resident has access to a house.
func (s *SocialSystem) hasShelter(resident *Resident, state *GameState) bool {
	for _, b := range state.Buildings {
		if b.Type == "house" {
			// Check if resident lives here (simplified: any house provides shelter for all)
			// In future, could track occupancy in metadata
			return true
		}
	}
	return false
}

// hasWork checks if resident is employed in any building.
func (s *SocialSystem) hasWork(resident *Resident, state *GameState) bool {
	for _, b := range state.Buildings {
		for _, workerID := range b.Workers {
			if workerID == resident.ID {
				return true
			}
		}
	}
	return false
}

// improveSkills increases skill levels based on work experience.
func (s *SocialSystem) improveSkills(resident *Resident, state *GameState, rng *rand.Rand) {
	hasWork := s.hasWork(resident, state)
	baseChance := 0.05 // 5% base chance
	if hasWork {
		baseChance = 0.2 // 20% chance if employed
	}

	for i := range resident.Skills {
		skill := &resident.Skills[i]

		// Chance to improve each week
		if rng.Float64() < baseChance {
			skill.Level++
		}
	}
}

// updateRelationships modifies relationships between residents.
func (s *SocialSystem) updateRelationships(state *GameState, rng *rand.Rand, events *[]Event) {
	// Simple random relationship formation
	for i := range state.Residents {
		for j := range state.Residents {
			if i >= j {
				continue // avoid duplicate pairs
			}

			r1 := &state.Residents[i]
			r2 := &state.Residents[j]

			// Check if relationship already exists
			exists := false
			for _, rel := range r1.Relationships {
				if rel.To == r2.ID {
					exists = true
					break
				}
			}

			if !exists && rng.Float64() < 0.05 { // 5% chance to form relationship
				rel := Relationship{
					From:     r1.ID,
					To:       r2.ID,
					Type:     "friend",
					Strength: 0.1 + rng.Float64()*0.2,
				}
				r1.Relationships = append(r1.Relationships, rel)
				// Add reciprocal relationship
				r2.Relationships = append(r2.Relationships, Relationship{
					From:     r2.ID,
					To:       r1.ID,
					Type:     "friend",
					Strength: rel.Strength,
				})
			}
		}
	}
}

// processLifeEvents handles births, aging, and deaths.
func (s *SocialSystem) processLifeEvents(week int, state *GameState, rng *rand.Rand, events *[]Event) {
	// Aging is handled in updateResident

	// Death from old age or critical needs
	for i := range state.Residents {
		resident := &state.Residents[i]

		// Death from old age
		deathProbability := 0.0
		if resident.Age > 80 {
			deathProbability = float64(resident.Age-80) * 0.001 // 0.1% per year over 80
		}

		// Death from critical hunger
		for _, need := range resident.Needs {
			if need.ID == "hunger" && need.Level > 0.9 {
				deathProbability += 0.01
			}
		}

		if rng.Float64() < deathProbability {
			// Record death
			*events = append(*events, Event{
				Type: "death",
				Data: map[string]interface{}{
					"resident_id": resident.ID,
					"cause":       "natural",
				},
			})
			// Remove resident (simplified: mark for removal)
			// In real implementation, would need to handle removal from slice
		}
	}

	// Births based on relationships
	// Simplified: couples produce children
	for i := range state.Residents {
		for j := range state.Residents {
			if i >= j {
				continue
			}
			r1 := &state.Residents[i]
			r2 := &state.Residents[j]

			// Check if they have strong relationship
			strongRel := false
			for _, rel := range r1.Relationships {
				if rel.To == r2.ID && rel.Type == "spouse" && rel.Strength > 0.5 {
					strongRel = true
					break
				}
			}

			if strongRel && rng.Float64() < 0.001 { // 0.1% chance per week
				// Create new resident
				baby := Resident{
					ID:   fmt.Sprintf("baby-%d-%d", week, rng.Int()),
					Name: "Child",
					Age:  0,
				}
				state.AddResident(baby)

				*events = append(*events, Event{
					Type: "birth",
					Data: map[string]interface{}{
						"parents": []string{r1.ID, r2.ID},
						"baby_id": baby.ID,
					},
				})
			}
		}
	}
}
