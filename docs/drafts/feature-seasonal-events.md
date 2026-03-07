# Seasonal & Event System

**Vision:** Create a living medieval world where seasons dramatically affect village life and unexpected events—both good and bad—create emergent stories that players will remember and share.

**Status:** Draft  
**Created:** March 7, 2026  
**Priority:** P0 (MVP Atmosphere)

---

## Executive Summary

The seasonal system imposes a relentless, predictable rhythm on village life: spring planting, summer growth, autumn harvest, winter survival. Layered on this are random events—plagues, marriages, discoveries, accidents—that create unique stories for each playthrough. Together, they transform mechanical resource management into an emotional narrative where players prepare for known challenges and adapt to unexpected crises.

**Success Definition:** Players feel the changing seasons in their bones (spring's hope, winter's dread) and recount specific events ("Remember when the blacksmith found that iron vein right before winter?") as personal village legends.

---

## User Stories

### US-001: Seasonal Cycle
**As:** A player experiencing medieval rural life  
**I want:** Distinct seasons that dramatically affect gameplay  
**So that:** I must plan differently for spring planting vs. winter survival

**Acceptance Criteria:**
- [ ] Four visually distinct seasons (spring, summer, autumn, winter)
- [ ] Seasonal effects: planting/growth/harvest cycles, temperature, daylight
- [ ] Winter: no agriculture, increased fuel needs, health risks from cold
- [ ] Summer: maximum growth, potential drought, festival opportunities
- [ ] Seasonal preparation feels necessary, not optional

### US-002: Weather System
**As:** A player dependent on medieval agriculture  
**I want:** Weather that affects crop yields and village operations  
**So that:** I must adapt to favorable and unfavorable conditions

**Acceptance Criteria:**
- [ ] Weekly weather generation: rain, sun, storms, snow
- [ ] Rain increases crop growth but may cause flooding
- [ ] Drought reduces yields, increases fire risk
- [ ] Storms damage buildings, delay outdoor work
- [ ] Snow blocks travel, increases heating needs
- [ ] Weather forecasts (traditional signs) with partial accuracy

### US-003: Predictable Seasonal Events
**As:** A player learning medieval rhythms  
**I want:** Annual events that repeat each year  
**So that:** I can anticipate and prepare for them

**Acceptance Criteria:**
- [ ] Spring: Planting festival, livestock breeding
- [ ] Summer: Midsummer festival, haymaking
- [ ] Autumn: Harvest festival, slaughter season
- [ ] Winter: Christmas/Yule, darkest day rituals
- [ ] Events affect resident happiness, productivity, relationships
- [ ] Player can choose to participate or ignore (with consequences)

### US-004: Random Events
**As:** A player seeking emergent stories  
**I want:** Unexpected events that create unique narratives  
**So that:** Each playthrough feels different and memorable

**Acceptance Criteria:**
- [ ] Positive events: resource discovery, skill breakthrough, beneficial stranger
- [ ] Negative events: disease outbreak, accident, crime, natural disaster
- [ ] Neutral events: traveling merchant, marriage proposal, unusual occurrence
- [ ] Events chain together (drought → famine → illness → death)
- [ ] Event frequency adjustable by difficulty setting

### US-005: Event Decision Making
**As:** A player faced with unexpected situations  
**I want:** Meaningful choices when events occur  
**So that:** My decisions have lasting consequences

**Acceptance Criteria:**
- [ ] Events present 2-4 options with different costs/benefits
- [ ] Decisions affect: resource costs, resident morale, relationships, future events
- [ ] Some decisions have delayed consequences (help stranger now, get help later)
- [ ] Moral dilemmas with no clear "best" choice
- [ ] Decision history tracked and affects village reputation

### US-006: Historical Events
**As:** A history enthusiast  
**I want:** Events based on actual medieval occurrences  
**So that:** I learn about history through gameplay

**Acceptance Criteria:**
- [ ] Plague outbreaks with realistic spread patterns
- [ ] Little Ice Age effects on agriculture
- [ ] Religious events (pilgrimages, heresy accusations)
- [ ] Feudal obligations (tax collector visits, military levy)
- [ ] Events can be enabled/disabled for historical accuracy vs. gameplay

---

## Technical Specifications

### Calendar & Season System
```go
type Calendar struct {
    Year        int            `json:"year"`        // Starting ~1300 AD
    Season      Season         `json:"season"`      // Spring, Summer, Autumn, Winter
    Week        int            `json:"week"`        // 1-52
    Day         int            `json:"day"`         // 1-7 (Monday-Sunday)
    Temperature float64        `json:"temperature"` // Celsius
    Weather     WeatherType    `json:"weather"`     // Current conditions
}

type Season int

const (
    SeasonSpring Season = iota // Weeks 1-13
    SeasonSummer              // Weeks 14-26
    SeasonAutumn              // Weeks 27-39
    SeasonWinter              // Weeks 40-52
)

// Seasonal modifiers applied to various systems
type SeasonalModifiers struct {
    Season          Season
    GrowthModifier  float64 // Crop growth rate
    WorkModifier    float64 // Outdoor work efficiency
    HungerModifier  float64 // Food consumption rate
    HealthModifier  float64 // Sickness probability
    MoodModifier    float64 // Resident happiness
}
```

### Weather System
```go
type WeatherType int

const (
    WeatherClear WeatherType = iota
    WeatherCloudy
    WeatherRain
    WeatherStorm
    WeatherSnow
    WeatherFog
)

type WeatherForecast struct {
    Current     WeatherType    `json:"current"`
    Tomorrow    WeatherType    `json:"tomorrow"`
    Accuracy    float64        `json:"accuracy"` // 0-1 how reliable forecast is
    Signs       []WeatherSign  `json:"signs"`    // Medieval forecasting signs
}

// Generate weather for a week
func GenerateWeather(season Season, rng *rand.Rand) []WeatherType {
    weather := make([]WeatherType, 7)
    
    baseProbabilities := map[Season]map[WeatherType]float64{
        SeasonSpring: {
            WeatherClear: 0.3,
            WeatherCloudy: 0.4,
            WeatherRain: 0.25,
            WeatherStorm: 0.05,
        },
        SeasonSummer: {
            WeatherClear: 0.6,
            WeatherCloudy: 0.3,
            WeatherRain: 0.08,
            WeatherStorm: 0.02,
        },
        // ... similar for autumn and winter
    }
    
    for day := 0; day < 7; day++ {
        // Markov chain: today's weather affects tomorrow's
        if day > 0 {
            // Rain often continues, clear often stays clear
            weather[day] = transitionWeather(weather[day-1], season, rng)
        } else {
            weather[day] = randomWeather(season, rng)
        }
    }
    
    return weather
}
```

### Event System Architecture
```go
// Event definition
type Event struct {
    ID          string          `json:"id"`
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Type        EventType       `json:"type"`       // Seasonal, Random, Historical
    Triggers    []Trigger       `json:"triggers"`   // Conditions for firing
    
    // Player choices
    Choices     []EventChoice   `json:"choices"`
    
    // Automatic effects (if no player choice needed)
    Effects     []EventEffect   `json:"effects"`
    
    // Visual/audio
    Image       string          `json:"image"`      // Event illustration
    Sound       string          `json:"sound"`      // Audio cue
}

// Event trigger conditions
type Trigger struct {
    Type        TriggerType     `json:"type"`
    Condition   interface{}     `json:"condition"` // Type-specific condition
    Probability float64         `json:"probability"` // Chance per check
}

type TriggerType int

const (
    TriggerSeasonal TriggerType = iota // Specific season/week
    TriggerResourceLevel               // Resource above/below threshold
    TriggerResidentState               // Resident condition (sick, skilled, etc.)
    TriggerBuildingCount               // Number of specific buildings
    TriggerRandom                      // Pure random chance
    TriggerPreviousEvent               // Follow-up to previous event
)

// Player choice during event
type EventChoice struct {
    Text        string          `json:"text"`        // Choice description
    Requirements []Requirement  `json:"requirements"` // Resources, skills needed
    Effects     []EventEffect   `json:"effects"`     // Immediate consequences
    Consequences []Consequence  `json:"consequences"` // Delayed/long-term effects
}

// Event effect on game state
type EventEffect struct {
    Type        EffectType      `json:"type"`
    Target      string          `json:"target"`      // Resource, resident, building ID
    Amount      float64         `json:"amount"`      // Change amount
    Duration    int             `json:"duration"`    // Weeks effect lasts (0 = permanent)
}
```

### Example Events

#### Seasonal: Spring Planting Festival
```json
{
  "id": "spring_planting_festival",
  "name": "Spring Planting Festival",
  "description": "The villagers gather to bless the fields and celebrate the start of the growing season. A successful festival boosts morale and ensures a good harvest.",
  "type": "seasonal",
  "triggers": [
    {
      "type": "seasonal",
      "condition": {"season": "spring", "week": 3},
      "probability": 1.0
    }
  ],
  "choices": [
    {
      "text": "Hold a lavish festival (cost: 50 food, 10 ale)",
      "requirements": [
        {"type": "resource", "resource": "food", "quantity": 50},
        {"type": "resource", "resource": "ale", "quantity": 10}
      ],
      "effects": [
        {"type": "happiness", "target": "all", "amount": 20, "duration": 4},
        {"type": "crop_yield", "target": "all", "amount": 1.15, "duration": 12}
      ]
    },
    {
      "text": "Simple blessings only (cost: 10 food)",
      "requirements": [
        {"type": "resource", "resource": "food", "quantity": 10}
      ],
      "effects": [
        {"type": "happiness", "target": "all", "amount": 5, "duration": 2},
        {"type": "crop_yield", "target": "all", "amount": 1.05, "duration": 12}
      ]
    },
    {
      "text": "Skip the festival entirely",
      "effects": [
        {"type": "happiness", "target": "all", "amount": -10, "duration": 4}
      ],
      "consequences": [
        {
          "type": "event",
          "event_id": "poor_harvest_complaints",
          "delay_weeks": 20,
          "probability": 0.7
        }
      ]
    }
  ]
}
```

#### Random: Iron Vein Discovery
```json
{
  "id": "iron_vein_discovery",
  "name": "Iron Vein Discovery",
  "description": "Miners have found a rich vein of iron ore while expanding the mine! This could solve our tool shortage for years.",
  "type": "random",
  "triggers": [
    {
      "type": "resource",
      "condition": {"resource": "tools", "comparison": "<", "threshold": 10},
      "probability": 0.05
    },
    {
      "type": "building",
      "condition": {"building": "mine", "count": ">=1"},
      "probability": 1.0
    }
  ],
  "choices": [
    {
      "text": "Redirect miners to exploit the vein immediately",
      "requirements": [
        {"type": "workers", "skill": "mining", "count": 3}
      ],
      "effects": [
        {"type": "resource", "resource": "iron_ore", "amount": 200},
        {"type": "happiness", "target": "miners", "amount": 15, "duration": 8}
      ],
      "consequences": [
        {
          "type": "event",
          "event_id": "mine_collapse_risk",
          "delay_weeks": 4,
          "probability": 0.3
        }
      ]
    },
    {
      "text": "Extract carefully to preserve the mine structure",
      "requirements": [
        {"type": "workers", "skill": "mining", "count": 2}
      ],
      "effects": [
        {"type": "resource", "resource": "iron_ore", "amount": 120},
        {"type": "happiness", "target": "miners", "amount": 5, "duration": 8}
      ]
    }
  ]
}
```

#### Historical: Plague Outbreak
```json
{
  "id": "plague_outbreak",
  "name": "The Coughing Sickness",
  "description": "A traveler brought sickness to the village. People are coughing, feverish, and several have already died. The disease seems to spread quickly.",
  "type": "historical",
  "triggers": [
    {
      "type": "random",
      "probability": 0.02
    },
    {
      "type": "resident",
      "condition": {"state": "traveler_recent", "count": ">=1"},
      "probability": 0.1
    }
  ],
  "choices": [
    {
      "text": "Quarantine the sick in the old barn",
      "requirements": [
        {"type": "building", "building": "barn", "count": 1}
      ],
      "effects": [
        {"type": "disease_spread", "modifier": 0.5},
        {"type": "happiness", "target": "all", "amount": -20, "duration": 8}
      ]
    },
    {
      "text": "Pray for divine intervention",
      "requirements": [
        {"type": "resource", "resource": "herbs", "quantity": 5}
      ],
      "effects": [
        {"type": "disease_spread", "modifier": 0.8},
        {"type": "faith", "target": "all", "amount": 10, "duration": 52}
      ]
    },
    {
      "text": "Continue normal life and hope it passes",
      "effects": [
        {"type": "disease_spread", "modifier": 1.5},
        {"type": "happiness", "target": "all", "amount": -5, "duration": 4}
      ]
    }
  ]
}
```

### Event Processing Algorithm
```go
func ProcessEvents(week int, gameState *GameState, rng *rand.Rand) []EventResult {
    var results []EventResult
    
    // Check seasonal events first
    for _, event := range seasonalEvents {
        if event.ShouldTrigger(week, gameState) {
            results = append(results, TriggerEvent(event, gameState, rng))
        }
    }
    
    // Check random events (weighted by probability)
    for _, event := range randomEvents {
        if rng.Float64() < event.BaseProbability {
            if event.ConditionsMet(gameState) {
                results = append(results, TriggerEvent(event, gameState, rng))
                // Limit random events per week
                if len(results) >= 2 {
                    break
                }
            }
        }
    }
    
    // Check triggered consequences from previous events
    for _, pending := range gameState.PendingConsequences {
        if pending.TriggerWeek == week {
            if rng.Float64() < pending.Probability {
                event := GetEventByID(pending.EventID)
                results = append(results, TriggerEvent(event, gameState, rng))
            }
        }
    }
    
    return results
}
```

### UI Interface Requirements

#### Season & Weather Display
- **Season indicator:** Large icon/color coding current season
- **Calendar view:** Year/week display with seasonal markers
- **Weather visualization:** Animated effects (rain, snow, sun)
- **Temperature display:** Numerical + feel (freezing, cold, mild, warm, hot)
- **Seasonal forecast:** Upcoming seasonal changes highlighted

#### Event Notification System
- **Pop-up events:** Modal dialog for important events with choices
- **Notification feed:** Sidebar showing recent events
- **Event history:** Complete log of past events with filters
- **Choice consequences:** Visual feedback on past decisions' outcomes

#### Seasonal Preparation Interface
- **Winter readiness checklist:** Food stores, fuel, warm clothing
- **Planting calendar:** Optimal times for different crops
- **Festival planner:** Upcoming seasonal events with preparation requirements
- **Weather preparedness:** Buildings needing repair before storms

#### Event Browser
- **Event encyclopedia:** All possible events with triggering conditions
- **Choice analysis:** Statistics on player choice outcomes
- **Historical context:** Real-world basis for historical events
- **Mod event integration:** Community-created events display

---

## Integration Points

### With Resource Economy System
- Seasonal effects on agriculture production
- Weather impacts on resource gathering
- Event consequences affecting resource quantities
- Festival costs consuming resources

### With Resident Management System
- Seasonal mood effects on residents
- Disease events affecting resident health
- Marriage/birth/death events
- Skill development through event participation

### With Building System
- Weather damage to buildings
- Seasonal building usage (granaries in autumn, hearths in winter)
- Event requirements for specific buildings
- Construction limitations in winter

### With Save/Load System
- Event history saved for continuity
- Pending consequences serialized
- Random seed preservation for deterministic events

---

## Balancing & Tuning

### Event Frequency
- **Small events:** 1-2 per week (weather, minor incidents)
- **Medium events:** 1-2 per season (festivals, discoveries)
- **Large events:** 1-2 per year (plague, famine, war)
- **Difficulty settings:** Adjust frequency and severity

### Seasonal Effects
| Season | Growth | Work Speed | Hunger | Health Risk | Mood |
|--------|--------|------------|--------|-------------|------|
| Spring | 110% | 90% | 100% | Medium | +5 |
| Summer | 150% | 80% (heat) | 105% | High (disease) | +10 |
| Autumn | 100% | 100% | 100% | Low | +0 |
| Winter | 0% | 70% (cold) | 115% | High (cold) | -10 |

### Weather Probabilities
| Season | Clear | Cloudy | Rain | Storm | Snow |
|--------|-------|--------|------|-------|------|
| Spring | 30% | 40% | 25% | 5% | 0% |
| Summer | 60% | 30% | 8% | 2% | 0% |
| Autumn | 40% | 40% | 18% | 2% | 0% |
| Winter | 20% | 30% | 10% | 0% | 40% |

---

## Performance Considerations

### Optimization Strategies
- **Event pre-filtering:** Only check events with possible triggers
- **Probability caching:** Cache calculated probabilities for repeated checks
- **Batch processing:** Process all events in single game state pass
- **Lazy loading:** Event images/sounds loaded on demand

### Scaling Targets
| Event Count | Processing Time | Memory Usage |
|-------------|-----------------|--------------|
| 50 events   | < 5ms           | < 2 MB       |
| 200 events  | < 20ms          | < 8 MB       |
| 1000 events | < 100ms         | < 40 MB      |

### Memory Management
- Event history trimmed after 100 events
- Old weather data aggregated daily
- Seasonal effects pre-calculated per season
- Event images unloaded after display

---

## Testing Strategy

### Unit Tests
- Season transition logic
- Weather generation probabilities
- Event trigger conditions
- Choice consequence application

### Integration Tests
- Full year seasonal cycle
- Event chains (drought → famine → plague)
- Save/load with event history
- Difficulty setting effects on event frequency

### Playtesting Focus
- Seasonal rhythm feeling natural, not arbitrary
- Event frequency creating stories without overwhelming
- Choices feeling meaningful with visible consequences
- Historical events feeling educational, not preachy

---

## Risks & Mitigations

### Technical Risks
1. **Event system becoming spaghetti code**  
   **Mitigation:** Strict interface between events and game systems, event DSL for definition

2. **Save file bloat from event history**  
   **Mitigation:** Optional detailed history, compression, archiving old events

3. **Probability calculations affecting performance**  
   **Mitigation:** Pre-computed probability tables, sampling optimization

### Design Risks
1. **Events feeling random/disconnected**  
   **Mitigation:** Event chains, cause-effect relationships, resident memory of events

2. **Seasonal cycle becoming predictable/boring**  
   **Mitigation:** Variable weather, unexpected seasonal events, climate change (long-term)

3. **Historical accuracy limiting fun**  
   **Mitigation:** Optional realism settings, "what-if" events, player agency in outcomes

### Mitigation Strategies
- **Event playtesting:** Test events individually before integration
- **Difficulty gradient:** Very easy mode reduces event frequency/severity
- **Player feedback:** Event rating system to identify popular/unpopular events
- **Mod support:** Allow players to create/share events

---

## Success Metrics

### Technical Metrics
- [ ] Event processing <10ms per week
- [ ] Event history save/load <50ms
- [ ] Memory usage <10MB for 200 events
- [ ] Zero event trigger bugs in 1000 test years

### Player Experience Metrics
- [ ] 80% of players can name a memorable event from their playthrough
- [ ] Seasonal changes feel significant, not cosmetic
- [ ] Event choices create anxiety/pleasure (emotional engagement)
- [ ] Players discuss events as stories ("my village survived the plague by...")

### Gameplay Metrics
- [ ] Event frequency: 1-2 memorable events per 4 hours of play
- [ ] Choice diversity: No single "best" choice for >70% of events
- [ ] Consequence visibility: Players can trace effects of decisions
- [ ] Replayability: Different events prominent in different playthroughs

---

## Dependencies

### Required First
- Basic calendar system
- Game state accessible to events
- UI event notification framework

### Dependent Features
- Historical event pack (requires event system)
- Dynamic event generation (requires resident relationships)
- Event modding tools (requires event definition format)
- Multi-event chains (requires consequence tracking)

---

## Open Questions

1. **Determinism vs. surprise:** Should events be completely deterministic from seed, or have some true randomness?
2. **Historical education:** How much historical explanation should accompany events?
3. **Player control:** Should players be able to disable certain event types?
4. **Event scaling:** Should event frequency scale with village size/age?

---

*Feature Version 1.0 · Owner: Narrative Design Team · Estimated Effort: 4-6 sprints*