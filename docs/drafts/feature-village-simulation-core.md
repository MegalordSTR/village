# Village Simulation Core

**Vision:** Create the foundational system that simulates a medieval village as a living, breathing entity where time passes in weekly turns and every element interacts.

**Status:** Draft  
**Created:** March 7, 2026  
**Priority:** P0 (MVP Foundation)

---

## Executive Summary

The village simulation core is the heart of the game - a deterministic engine that processes one week of game time, updating all residents, resources, buildings, and environmental factors based on player decisions. This system must be performant (<3s per turn for 50 residents), completely deterministic (identical inputs → identical outputs), and modular enough to support future expansions.

**Success Definition:** Player makes strategic decisions, clicks "Advance Week," and watches as their choices play out through interconnected systems, feeling both in control and surprised by emergent outcomes.

---

## User Stories

### US-001: Weekly Turn Processing
**As:** A player who has made decisions for the upcoming week  
**I want:** To advance time by one week and see the results  
**So that:** I can observe how my strategic choices affect the village

**Acceptance Criteria:**
- [ ] Clicking "Advance Week" processes exactly 7 days of game time
- [ ] All systems update in correct order: environment → production → needs → events
- [ ] Processing completes in <3 seconds for village with 50 residents
- [ ] Results are deterministic - same inputs always produce same outputs
- [ ] Visual feedback shows time progression (day counter, visual changes)

### US-002: Deterministic Game State
**As:** A player who wants reliable save/load functionality  
**I want:** My game state to be completely reproducible from save files  
**So that:** I can save and resume games without unexpected changes

**Acceptance Criteria:**
- [ ] Game state serialized to JSON includes all simulation data
- [ ] Random events use seed-based RNG that's saved and restored
- [ ] Loading a save file recreates identical game state
- [ ] Simulation produces identical results when re-run from same state
- [ ] Save files are human-readable for debugging

### US-003: Performance with Scale
**As:** A player managing a large village  
**I want:** Turn processing to remain fast as my village grows  
**So that:** I'm not waiting minutes between decisions

**Acceptance Criteria:**
- [ ] 10-resident village processes in <1 second
- [ ] 50-resident village processes in <3 seconds
- [ ] 100-resident village processes in <5 seconds
- [ ] Processing time scales linearly, not exponentially
- [ ] Background processing during player planning phase

### US-004: Modular System Architecture
**As:** A developer adding new game systems  
**I want:** Clear interfaces between simulation subsystems  
**So that:** I can add features without breaking existing functionality

**Acceptance Criteria:**
- [ ] Each system (economy, social, environmental) has defined API
- [ ] Systems can be enabled/disabled without breaking game
- [ ] New resource types can be added via configuration
- [ ] Event system allows custom events without code changes
- [ ] All systems participate in save/load cycle

---

## Technical Specifications

### Data Model
```go
// Core game state structure
type GameState struct {
    ID          string    `json:"id"`
    Version     string    `json:"version"`
    Seed        int64     `json:"seed"`       // For deterministic RNG
    Calendar    Calendar  `json:"calendar"`   // Current date, season
    Village     Village   `json:"village"`    // Physical layout
    Residents   []Resident `json:"residents"` // Population
    Resources   Resources `json:"resources"`  // Inventory
    Buildings   []Building `json:"buildings"` // Structures
    History     []Event   `json:"history"`    // Past events
    Policies    Policies  `json:"policies"`   // Player decisions
}

// Weekly turn processing order
type TurnProcessor struct {
    systems []System // Ordered list of systems to update
}

// System interface for modular design
type System interface {
    Update(week int, state *GameState, rng *rand.Rand) []Event
    LoadConfig(config json.RawMessage) error
}
```

### Processing Order
1. **Environmental Systems** (week 0-6)
   - Day 0: Advance calendar, update season/weather
   - Day 1: Process temperature effects on residents
   - Day 2: Update soil fertility, plant growth
   - Day 3: Natural resource regeneration (forests, mines)
   - Day 4: Wildlife population changes
   - Day 5: Weather effects on buildings
   - Day 6: Seasonal transition checks

2. **Production Systems** (parallel processing)
   - Agriculture: Crop growth, harvesting
   - Mining: Resource extraction
   - Crafting: Goods production
   - Construction: Building progress

3. **Social Systems**
   - Resident needs calculation (hunger, warmth, happiness)
   - Skill development through work
   - Relationship changes
   - Life events (birth, aging, death)

4. **Economic Systems**
   - Resource consumption
   - Inventory updates
   - Trade calculations
   - Wealth distribution

5. **Event System**
   - Check for random events based on conditions
   - Apply event consequences
   - Add to history log

6. **Cleanup & Validation**
   - Remove dead residents
   - Destroy ruined buildings
   - Validate game state integrity
   - Prepare for next turn

### Performance Targets
| Village Size | Max Processing Time | Memory Usage |
|--------------|---------------------|--------------|
| 10 residents | 1 second            | < 50 MB      |
| 50 residents | 3 seconds           | < 200 MB     |
| 100 residents| 5 seconds           | < 400 MB     |
| 200 residents| 8 seconds           | < 800 MB     |

### Determinism Requirements
- All random numbers from seeded RNG (crypto/rand with seed)
- Floating-point calculations with fixed precision
- Map iteration in consistent order (sorted keys)
- Time-based events use game clock, not real time
- External inputs (API calls) mocked during simulation

---

## Integration Points

### With Resident Management System
- Receives resident assignments from UI
- Updates resident states based on work and needs
- Returns resident status changes for visualization

### With Resource Economy System
- Receives production recipes and rates
- Updates resource quantities based on production
- Returns resource changes for UI updates

### With Seasonal Event System
- Receives event definitions and probabilities
- Triggers events based on game state
- Returns event results for player notification

### With Save/Load System
- Provides complete serializable game state
- Receives loaded state for continuation
- Validates state integrity on load

---

## Testing Strategy

### Unit Tests
- **Determinism Tests:** Verify identical outputs from same inputs
- **Performance Tests:** Benchmark processing times at different scales
- **Edge Cases:** Empty village, maximum residents, extreme values
- **System Isolation:** Test each subsystem independently

### Integration Tests
- **Full Turn Cycle:** Complete week processing with all systems
- **Save/Load Cycle:** Save → load → verify identical state
- **Scale Tests:** Gradually increase village size
- **Regression Tests:** Verify fixes don't break existing behavior

### Playtesting
- **Alpha Testers:** Internal team - focus on system stability
- **Beta Testers:** External strategy gamers - focus on feel and balance
- **Performance Monitoring:** Real-world timing data collection

---

## Risks & Mitigations

### Technical Risks
1. **Non-deterministic behavior**  
   **Mitigation:** Extensive determinism testing, avoid goroutine race conditions

2. **Performance degradation with scale**  
   **Mitigation:** Algorithmic optimization, profiling early, caching results

3. **Save file corruption**  
   **Mitigation:** Checksum validation, backup system, graceful recovery

### Design Risks
1. **Too abstract for players to understand**  
   **Mitigation:** Clear visualization of cause-effect chains, tutorial system

2. **Weekly pace feels too slow**  
   **Mitigation:** Option to process multiple weeks, meaningful events each week

3. **Lack of emergent storytelling**  
   **Mitigation:** Rich event system, resident personality traits, history tracking

### Mitigation Strategies
- **Early Prototyping:** Core simulation before UI polish
- **Continuous Playtesting:** Feedback from day one of development
- **Modular Design:** Ability to replace subsystems if needed
- **Performance Budgets:** Strict limits on processing time

---

## Success Metrics

### Development Metrics
- [ ] 80%+ unit test coverage for simulation engine
- [ ] All determinism tests pass 1000 iterations
- [ ] Performance targets met at all scale levels
- [ ] Zero memory leaks in 24-hour stress test

### Player Experience Metrics
- [ ] 95% of players understand turn-based concept within first session
- [ ] Average turn processing time <2.5 seconds (target <3s)
- [ ] Save/load reliability >99.9% (no corrupted saves)
- [ ] Players report feeling of "living world" in feedback

---

## Dependencies

### Required First
- Basic game state data structures
- Deterministic random number system
- System interface definitions

### Dependent Features
- Resident Management System (requires simulation core)
- Resource Economy System (requires simulation core)
- Seasonal Event System (requires simulation core)
- Save/Load System (requires serializable game state)

---

## Open Questions

1. **Processing Parallelism:** How much can be safely parallelized while maintaining determinism?
2. **State Size:** How much game history should be kept in memory vs. archived?
3. **Validation:** What level of game state validation is needed after each turn?
4. **Modding:** How much of the simulation order should be configurable for mods?

---

*Feature Version 1.0 · Owner: Backend Team · Estimated Effort: 6-8 sprints*