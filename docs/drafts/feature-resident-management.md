# Resident Management System

**Vision:** Create believable medieval villagers with individual skills, needs, relationships, and life cycles that make players care about their digital population while providing deep strategic control.

**Status:** Draft  
**Created:** March 7, 2026  
**Priority:** P0 (MVP Core)

---

## Executive Summary

The resident management system transforms abstract "population numbers" into relatable individuals with personalities, capabilities, and connections. Each resident has unique skills that develop through work, basic needs that must be met, relationships that affect productivity, and a life cycle from birth to death. Players manage these individuals directly or through policies, creating emotional investment and strategic depth.

**Success Definition:** Players know residents by name, make assignments based on individual strengths, mourn losses, celebrate successes, and feel responsible for their villagers' wellbeing.

---

## User Stories

### US-001: Individual Resident Profiles
**As:** A player examining my village  
**I want:** To see detailed information about each resident  
**So that:** I can make informed decisions about their assignments

**Acceptance Criteria:**
- [ ] Clicking a resident shows: name, age, portrait, skills, needs, relationships
- [ ] Skills displayed with numeric values (0-100) and progress toward next level
- [ ] Needs shown as bars (hunger, health, happiness) with color coding
- [ ] Relationships shown as network (family, friends, conflicts)
- [ ] Work history shows past assignments and performance

### US-002: Skill-Based Assignment
**As:** A player optimizing village productivity  
**I want:** To assign residents to jobs matching their skills  
**So that:** My village operates efficiently

**Acceptance Criteria:**
- [ ] Drag resident to workplace or use assignment interface
- [ ] System suggests optimal assignments based on skills
- [ ] Skill mismatch shows predicted productivity penalty
- [ ] Multiple residents can be assigned to same workplace
- [ ] Unassigned residents perform basic maintenance tasks

### US-003: Needs System
**As:** A resident in the simulation  
**I want:** My basic needs to be tracked and affect my productivity  
**So that:** Players must maintain village wellbeing

**Acceptance Criteria:**
- [ ] Hunger increases daily, decreases when eating
- [ ] Health affected by food, shelter, medical care
- [ ] Happiness affected by relationships, work satisfaction, leisure
- [ ] Low needs reduce work efficiency, cause protests, or trigger emigration
- [ ] Needs visible to player with warning indicators

### US-004: Skill Development
**As:** A resident working regularly  
**I want:** My skills to improve through practice  
**So that:** Players are rewarded for consistent assignments

**Acceptance Criteria:**
- [ ] Working at job increases relevant skill by 0.1-1.0 points per week
- [ ] Higher skill levels require more experience to advance
- [ ] Natural aptitude affects learning rate (some residents learn faster)
- [ ] Skill degradation when not used for extended periods
- [ ] Mastery bonuses at skill levels 25, 50, 75, 100

### US-005: Life Cycle Events
**As:** A player observing village over time  
**I want:** Residents to age, form relationships, and experience life events  
**So that:** The village feels like a living community

**Acceptance Criteria:**
- [ ] Residents age 1 year per 52 game weeks
- [ ] Marriage between compatible residents creates families
- [ ] Pregnancy lasts 40 weeks, then child birth
- [ ] Children grow into adults over 16 years
- [ ] Death from old age (60-80), illness, or accidents
- [ ] Inheritance passes skills and relationships to children

### US-006: Relationship Network
**As:** A resident interacting with others  
**I want:** Relationships that affect my behavior and happiness  
**So that:** Social dynamics emerge naturally

**Acceptance Criteria:**
- [ ] Family relationships (parent, child, sibling, spouse)
- [ ] Friendship develops from working/living together
- [ ] Conflicts arise from competition, injustice, or personality clashes
- [ ] Relationships affect: work efficiency when together, happiness, event triggers
- [ ] Network visualization shows connections between residents

---

## Technical Specifications

### Data Model
```go
// Resident represents a single villager
type Resident struct {
    ID          string          `json:"id"`
    Name        string          `json:"name"`
    Age         int             `json:"age"`           // Years
    Gender      Gender          `json:"gender"`
    
    // Physical state
    Health      Need            `json:"health"`        // 0-100
    Hunger      Need            `json:"hunger"`        // 0-100
    Energy      Need            `json:"energy"`        // 0-100
    
    // Psychological state
    Happiness   Need            `json:"happiness"`     // 0-100
    Stress      float64         `json:"stress"`        // 0-1
    
    // Skills (0-100 scale)
    Skills      map[SkillType]Skill `json:"skills"`
    
    // Work & status
    Assignment  *Assignment     `json:"assignment"`    // Current job
    Home        *Building       `json:"home"`         // Residence
    
    // Relationships
    Family      FamilyRelations `json:"family"`
    Friends     []Relationship  `json:"friends"`      // Positive
    Conflicts   []Relationship  `json:"conflicts"`    // Negative
    
    // Life events
    Pregnant    *Pregnancy      `json:"pregnant,omitempty"`
    Alive       bool            `json:"alive"`
    
    // Traits (personality modifiers)
    Traits      []Trait         `json:"traits"`
}

// Skill with experience tracking
type Skill struct {
    Level       int     `json:"level"`       // 0-100
    Experience  float64 `json:"experience"`  // Toward next level
    Aptitude    float64 `json:"aptitude"`    // 0.5-1.5 learning modifier
}

// Relationship between two residents
type Relationship struct {
    TargetID    string  `json:"targetId"`
    Type        RelType `json:"type"`       // Friend, Rival, Lover, etc.
    Strength    float64 `json:"strength"`   // 0-100
    History     []Interaction `json:"history"`
}
```

### Need Calculation (Weekly)
```go
func (r *Resident) CalculateNeeds(week int, village *Village) {
    // Hunger increases daily, decreases with food consumption
    dailyHunger := 2.0 // Base hunger increase per day
    if r.Assignment != nil && r.Assignment.LaborIntensity > 0.5 {
        dailyHunger *= 1.5 // Physical work increases hunger
    }
    
    // Food consumption reduces hunger
    foodAvailable := village.GetFoodPerResident()
    hungerReduction := math.Min(foodAvailable * 10, 30.0)
    
    r.Hunger.Value = clamp(r.Hunger.Value + dailyHunger*7 - hungerReduction)
    
    // Health affected by hunger and shelter
    if r.Hunger.Value > 80 {
        r.Health.Value -= 5 // Starvation damage
    }
    
    // Happiness affected by multiple factors
    happinessModifiers := 0.0
    happinessModifiers += r.calculateWorkSatisfaction()
    happinessModifiers += r.calculateSocialSatisfaction()
    happinessModifiers -= r.Stress * 20
    
    r.Happiness.Value = clamp(50 + happinessModifiers) // Base 50 +/- modifiers
}
```

### Skill Development Algorithm
```go
func (r *Resident) UpdateSkills(week int) {
    if r.Assignment == nil {
        return
    }
    
    // Get relevant skill for current assignment
    skillType := r.Assignment.RequiredSkill
    skill := r.Skills[skillType]
    
    // Calculate experience gain
    baseGain := 1.0 // Base XP per week of work
    if r.Happiness.Value > 70 {
        baseGain *= 1.2 // Happy workers learn faster
    }
    if r.Health.Value < 30 {
        baseGain *= 0.5 // Sick workers learn slower
    }
    
    // Apply aptitude modifier
    baseGain *= skill.Aptitude
    
    // Add experience
    skill.Experience += baseGain
    
    // Level up if enough experience
    xpForNextLevel := 10 + skill.Level * 2 // Increasing requirement
    if skill.Experience >= xpForNextLevel {
        skill.Level = min(skill.Level + 1, 100)
        skill.Experience -= xpForNextLevel
    }
    
    r.Skills[skillType] = skill
}
```

### Relationship Dynamics
```go
func (r *Resident) UpdateRelationships(week int, otherResidents []*Resident) {
    for _, other := range otherResidents {
        if other.ID == r.ID {
            continue
        }
        
        // Check for interactions this week
        interacted := r.checkInteraction(other, week)
        
        if interacted {
            // Modify relationship based on interaction type
            rel := r.getRelationship(other.ID)
            
            // Positive interactions increase friendship
            if interactionWasPositive {
                rel.Strength = min(rel.Strength + 5, 100)
                if rel.Strength > 50 && rel.Type != Friend {
                    rel.Type = Friend
                }
            }
            
            // Negative interactions create conflicts
            if interactionWasNegative {
                rel.Strength = max(rel.Strength - 10, -100)
                if rel.Strength < -30 && rel.Type != Conflict {
                    rel.Type = Conflict
                }
            }
        }
    }
}
```

### UI Interface Requirements

#### Resident List View
- Sortable columns: Name, Age, Assignment, Skills, Needs
- Filter by: Assignment status, skill level, needs threshold
- Group by: Family, workplace, skill specialty
- Bulk actions: Assign multiple residents to same workplace

#### Resident Detail Modal
- **Basic Info:** Portrait, name, age, gender, traits
- **Status:** Needs bars with numerical values
- **Skills:** Grid of skill icons with levels, progress bars
- **Assignment:** Current job, productivity, satisfaction
- **Relationships:** Network visualization, list of connections
- **History:** Life events, work history, notable achievements

#### Assignment Interface
- Drag residents from list to workplace on map
- Workplace shows assigned residents with skill indicators
- Suggested assignments based on skill matching
- Productivity prediction for different assignments
- Conflict warnings for assigning incompatible residents together

#### Needs Dashboard
- Overview of village needs status
- Warning indicators for residents in critical need
- Trend lines showing needs over time
- Correlation between needs and productivity

---

## Integration Points

### With Village Simulation Core
- Receives weekly update call to process all residents
- Returns resident status changes for event system
- Provides productivity multipliers based on needs/skills

### With Resource Economy System
- Consumes resources based on resident needs (food, clothing)
- Provides labor input for production calculations
- Skill levels affect production quality and quantity

### With Building System
- Residents assigned to workplaces in buildings
- Home buildings affect need satisfaction (shelter, warmth)
- Building capacity limits resident assignments

### With Event System
- Resident traits and relationships affect event probabilities
- Events modify resident states (illness, injury, mood)
- Life events trigger notifications (birth, death, marriage)

---

## Performance Considerations

### Optimization Strategies
- **Batch Processing:** Update residents in batches, not individually
- **Lazy Evaluation:** Only calculate needs for residents with changes
- **Caching:** Skill calculations cached until relevant changes occur
- **Delta Updates:** Only send changed resident data to frontend

### Scaling Targets
| Resident Count | Update Time | Memory per Resident |
|----------------|-------------|---------------------|
| 10 residents   | < 50ms      | ~2 KB               |
| 50 residents   | < 200ms     | ~2 KB               |
| 100 residents  | < 400ms     | ~2 KB               |
| 200 residents  | < 800ms     | ~2 KB               |

### Memory Management
- Dead residents moved to history archive after 52 weeks
- Relationship history trimmed to last 100 interactions
- Skill experience history aggregated weekly
- Old need values aggregated into daily averages

---

## Balancing & Tuning

### Skill Development Rates
- **Novice (0-25):** 1-2 weeks per level
- **Competent (25-50):** 3-4 weeks per level  
- **Expert (50-75):** 8-12 weeks per level
- **Master (75-100):** 16-24 weeks per level

### Need Decay Rates
- **Hunger:** +10 per day without food, -30 per meal
- **Health:** -1 per day if hunger > 80, +2 per day if all needs met
- **Happiness:** Base 50, +/- 20 based on conditions

### Relationship Building
- **Positive interaction:** +5 strength
- **Negative interaction:** -10 strength  
- **Daily decay:** -0.1 strength if no interaction
- **Thresholds:** Friend (>50), Conflict (<-30), Neutral (-30 to 50)

---

## Testing Strategy

### Unit Tests
- Skill development calculations under various conditions
- Need decay with different resource availability
- Relationship changes from interactions
- Life event probability calculations

### Integration Tests
- Full resident lifecycle from birth to death
- Skill-based assignment productivity
- Needs affecting work performance
- Relationship networks affecting village dynamics

### Playtesting Focus
- Emotional connection to residents
- Understanding of skill/need systems
- Satisfaction with assignment interface
- Perception of resident individuality

---

## Risks & Mitigations

### Technical Risks
1. **Performance with large populations**  
   **Mitigation:** Batch processing, efficient data structures, profiling

2. **Complex relationship networks slowing updates**  
   **Mitigation:** Limit active relationships, sparse updates, caching

3. **Save file size explosion**  
   **Mitigation:** Compression, archiving old data, delta encoding

### Design Risks
1. **Residents feel like spreadsheets, not people**  
   **Mitigation:** Personality traits, portraits, names, life stories

2. **Micro-management overwhelming players**  
   **Mitigation:** Policy-based management, group assignments, AI suggestions

3. **Predictable skill development**  
   **Mitigation:** Aptitude variations, learning events, skill synergy

### Mitigation Strategies
- **Early UI Testing:** Assignment interface tested before backend complete
- **Progressive Disclosure:** Basic needs first, advanced relationships later
- **Difficulty Settings:** Adjustable complexity for different player types
- **Mod Support:** Allow mods to adjust all tuning parameters

---

## Success Metrics

### Technical Metrics
- [ ] Resident update time <200ms for 50 residents
- [ ] Save file size <100KB per 50 residents
- [ ] Memory usage <200MB for 100 residents
- [ ] Zero memory leaks in resident lifecycle

### Player Experience Metrics
- [ ] 80% of players can name at least 5 residents
- [ ] 70% make assignments based on individual skills
- [ ] Average resident lifespan >40 years (realistic medieval)
- [ ] Players report emotional response to resident events

### Gameplay Metrics
- [ ] Skill distribution follows realistic curve (few masters, many novices)
- [ ] Needs system creates meaningful trade-offs (food vs. productivity)
- [ ] Relationships affect gameplay visibly (friends work better together)
- [ ] Resident individuality perceived within first 2 hours

---

## Dependencies

### Required First
- Basic resident data structure
- Need calculation system
- Skill definition framework

### Dependent Features
- Assignment UI (requires resident data)
- Needs visualization (requires need system)
- Relationship network (requires relationship data)
- Life events (requires age and relationship tracking)

---

## Open Questions

1. **Detail Level:** How detailed should resident psychology be? Depression, anxiety, personality disorders?
2. **Control Granularity:** Should players control individuals directly or through policies?
3. **Historical Accuracy:** How closely should life expectancy match medieval reality (30-40 years)?
4. **Accessibility:** How to make complex resident data understandable without overwhelming?

---

*Feature Version 1.0 · Owner: Game Design Team · Estimated Effort: 4-6 sprints*