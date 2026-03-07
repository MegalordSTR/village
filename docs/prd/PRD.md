# Product Requirements Document: Medieval Village Simulator

## 1. Overview
**Project Codename:** Village  
**Type:** Browser-based historical economic simulator  
**Core Loop:** Weekly turn-based strategy with individual resident management  
**Target Platform:** Web (PWA capable), potential future desktop release  
**Tech Stack:** Go backend, PostgreSQL 18, Angular 21 frontend, S3-compatible storage  

## 2. Problem Statement
Strategy gamers seeking deep systemic complexity lack historically-accurate simulations that balance authentic medieval economics with engaging gameplay. Existing solutions either oversimplify (mobile games) or focus on fantasy elements rather than historical realism.

## 3. Goals & Success Criteria

### 3.1 Business Goals
- **MVP Launch:** Fully playable single village simulation within 6 months
- **Player Engagement:** Average session duration > 90 minutes, weekly retention > 60%
- **Community Growth:** 1000+ active players within 3 months of launch
- **Modding Ecosystem:** 50+ community-created mods within first year

### 3.2 Technical Goals
- **Performance:** Process 1 week of game time for 50-resident village in < 3 seconds
- **Determinism:** 100% reproducible game state from save files
- **Modularity:** Clear API boundaries for mod developers
- **Scalability:** Support 10,000+ concurrent games on modest hardware

### 3.3 User Experience Goals
- **Learnability:** New players understand basic systems within 30 minutes
- **Depth:** Players discover new interactions after 10+ hours
- **Clarity:** Complex relationships visualized without spreadsheet overload
- **Satisfaction:** "Aha!" moments when system connections become clear

## 4. User Personas

### 4.1 The History Enthusiast (Primary)
- **Demographics:** 25-45, university education, reads historical non-fiction
- **Needs:** Authentic medieval experience, educational value, attention to detail
- **Pain Points:** Historical inaccuracies in games, oversimplified economics
- **Success:** Feeling they've learned about medieval life while having fun

### 4.2 The Systems Thinker (Primary)
- **Demographics:** 20-35, STEM background, enjoys programming/logic puzzles
- **Needs:** Complex interconnected systems, emergent behavior, optimization challenges
- **Pain Points:** Games that are complex but not deep, lack of meaningful choices
- **Success:** Discovering unexpected system interactions and mastering them

### 4.3 The Storyteller (Secondary)
- **Demographics:** 18-30, creative background, active on social media/Twitch
- **Needs:** Emergent narratives, shareable moments, character development
- **Pain Points:** Games that are too abstract to generate stories
- **Success:** Creating compelling narratives from gameplay to share online

## 5. Core Features

### 5.1 MVP Features (P0)

#### 5.1.1 Village Simulation Core
- Single village with 10-50 residents
- Weekly turn-based progression
- 15 core resources with production chains:
  - **Basic:** Food (grain, vegetables, meat), Wood, Stone
  - **Processed:** Flour (from grain), Bread (from flour), Tools (from iron + wood)
  - **Advanced:** Iron (requires mining + smelting), Textiles (from wool + processing)
- Seasonal cycle (spring, summer, autumn, winter) affecting agriculture
- Basic needs system: food, shelter, warmth

#### 5.1.2 Resident Management
- Individual residents with:
  - Skills (farming, mining, crafting, etc.)
  - Needs (hunger, health, happiness)
  - Relationships (family, friendships, conflicts)
- Assignment system: drag residents to workplaces
- Skill development through practice
- Life cycle: birth, aging, death

#### 5.1.3 Building & Production
- 10+ building types: houses, farm fields, pastures, workshops, storage
- Production chains: raw materials → processed goods → advanced goods
- Resource storage with spoilage (food decays)
- Construction requiring materials and labor

#### 5.1.4 User Interface
- Angular 21 SPA with responsive design
- Main views:
  - Village overview (map with buildings)
  - Resident management (list/details)
  - Resource management (inventory, production)
  - Turn controls (plan/execute/review)
- Real-time visualization of turn results
- Save/load functionality

#### 5.1.5 Game Systems
- Weather system affecting agriculture
- Basic events: disease outbreaks, good/bad harvests, accidents
- Economy: internal resource valuation, basic trade (future expansion)
- Save system with export/import

### 5.2 Post-MVP Features (P1)

#### 5.2.1 Social Systems
- Social classes (peasants, craftsmen, clergy, nobility)
- Crime and punishment system
- Religious beliefs and ceremonies
- Education and knowledge transmission

#### 5.2.2 Advanced Economics
- Trade routes with neighboring villages
- Currency system (coins minted from precious metals)
- Taxation and feudal obligations
- Market price fluctuations based on supply/demand

#### 5.2.3 Environmental Depth
- Soil fertility depletion and recovery
- Animal husbandry with breeding
- Forestry management (replanting trees)
- Water management (wells, irrigation)

#### 5.2.4 Modding Support
- JSON-based configuration for all game data
- JavaScript API for custom game logic
- Asset pipeline for custom graphics
- Steam Workshop integration (future)

### 5.3 Future Features (P2+)
- Multiple villages with trade/conflict
- Historical events (Black Death, Little Ice Age)
- Technology tree spanning centuries
- Multiplayer (asynchronous turn sharing)
- Mobile app companion (view village status)

## 6. Technical Requirements

### 6.1 Backend (Go)
- **Simulation Engine:** Deterministic game state processing
- **REST API:** JSON endpoints for frontend communication
- **Database Layer:** PostgreSQL 18 with schema optimized for game state
- **Authentication:** Simple token-based auth (single player)
- **File Storage:** S3-compatible for save games and mods

### 6.2 Frontend (Angular 21 + TypeScript)
- **Framework:** Angular 21 with Material Design components
- **State Management:** NgRx for complex application state
- **Visualization:** SVG/D3.js for interactive village map
- **Performance:** Virtual scrolling for resident lists, Web Workers for background calculations
- **PWA:** Installable, offline-capable for save games

### 6.3 Data Model
```
Game
├── Village
│   ├── Residents[] (with skills, needs, relationships)
│   ├── Buildings[] (with type, condition, workers)
│   ├── Resources[] (with quantities, locations)
│   └── Policies (tax rates, work assignments)
├── Calendar (current date, season, weather)
└── History (past events, statistics)
```

### 6.4 Performance Targets
- **Turn Processing:** < 3 seconds for 50-resident village
- **UI Responsiveness:** < 100ms for all interactions
- **Initial Load:** < 2 seconds on broadband
- **Memory:** < 500MB for frontend with large villages

## 7. User Experience Specifications

### 7.1 Core Workflows

#### 7.1.1 Starting a New Game
1. Player selects starting conditions (village size, difficulty, historical period)
2. System generates initial village with residents and resources
3. Tutorial introduces basic interface and first decisions
4. Player makes initial assignments and starts first turn

#### 7.1.2 Weekly Turn Cycle
1. **Review Phase:** Player sees results of previous week
2. **Planning Phase:** Player assigns residents, starts constructions, sets policies
3. **Execution:** Player clicks "Advance Week," simulation runs
4. **Results:** Animation shows key events, summary screen appears
5. **Adjustment:** Player reviews detailed reports, prepares for next week

#### 7.1.3 Resident Management
1. Open resident list or click on villager in map
2. See skills, needs, current assignment
3. Drag to new workplace or change assignment via dropdown
4. View relationship network (family, friends, conflicts)

### 7.2 Interface Principles
- **Progressive Disclosure:** Basic info first, details on demand
- **Causal Visualization:** Show how decisions affect outcomes
- **Contextual Help:** Explanations appear when hovering complex elements
- **Multiple Views:** Same data available in map, list, and chart forms

### 7.3 Accessibility
- Color-blind friendly palettes
- Keyboard navigation for all functions
- Screen reader support for text content
- Adjustable text sizes

## 8. Non-Functional Requirements

### 8.1 Reliability
- Game state autosaved after each turn
- Save file validation to prevent corruption
- Graceful error recovery (corrupted saves can be repaired)
- Version compatibility for save files across updates

### 8.2 Security
- Client-side only for single player (no sensitive data)
- Save file encryption optional for privacy
- Input validation to prevent malicious save files
- No personal data collection beyond analytics opt-in

### 8.3 Maintainability
- Comprehensive unit tests for simulation engine
- Integration tests for full game loops
- API documentation for mod developers
- Database migration scripts for schema changes

### 8.4 Scalability
- Single server supports 10,000+ concurrent games
- Database sharding by game ID if needed
- Stateless API servers behind load balancer
- Asset caching via CDN for frontend

## 9. Constraints & Assumptions

### 9.1 Technical Constraints
- Must run in modern browsers (Chrome, Firefox, Safari, Edge)
- Backend must compile to single binary for easy deployment
- Database must be PostgreSQL 18+ for JSONB features
- Frontend must be Angular 21+ for long-term support

### 9.2 Business Constraints
- Development team: 2-3 developers for 6 months to MVP
- Budget: Limited, focus on core gameplay over polish
- Timeline: Playable alpha in 3 months, beta in 5, launch in 6
- Monetization: Premium one-time purchase, no microtransactions

### 9.3 User Assumptions
- Players are willing to learn complex systems
- Historical accuracy valued over game balance
- Long-term engagement expected (weeks/months of play)
- Community will create mods if tools provided

## 10. Success Metrics & Validation

### 10.1 MVP Validation Criteria
- **Playtest Completion:** 80% of testers complete 10+ game weeks
- **Understanding:** 70% can explain basic resource chains after 1 hour
- **Enjoyment:** Net Promoter Score > 30 from strategy game enthusiasts
- **Performance:** 95% of turns process in < 5 seconds on test hardware

### 10.2 Analytics Tracking
- Session length and frequency
- Most-used interface features
- Common failure points (where players give up)
- Resource chain comprehension progression
- Mod usage and creation rates

### 10.3 Quality Gates
- **Code Coverage:** > 80% for simulation engine
- **Load Testing:** 100 concurrent games stable for 24 hours
- **Browser Compatibility:** All major browsers pass functionality tests
- **Accessibility:** WCAG 2.1 AA compliance

## 11. Open Questions & Risks

### 11.1 Technical Risks
- Deterministic simulation proving difficult with complex systems
- Frontend performance with 50+ animated residents
- Save file size growing too large with game history
- Mod API being either too restrictive or too complex

### 11.2 Design Risks
- Simulation too complex for players to understand
- Historical accuracy making game unfun (too harsh)
- Weekly turn pace feeling too slow for some players
- Insufficient emergent storytelling opportunities

### 11.3 Business Risks
- Niche audience too small for sustainable development
- Competition from established titles (Banished, Rimworld)
- Development time exceeding estimates due to complexity
- Modding community not developing as hoped

### 11.4 Mitigation Strategies
- Early prototyping of simulation core
- Frequent playtesting with target audience
- Modular design allowing features to be cut
- Community engagement from earliest alpha

## 12. Glossary
- **Turn:** One week of game time
- **Resident:** Individual villager with skills and needs
- **Resource Chain:** Raw material → processed good → advanced good
- **Deterministic:** Same inputs always produce same outputs
- **JSONB:** PostgreSQL binary JSON storage for flexible schemas

---

*PRD Version 1.0 · Created: March 7, 2026 · Owner: Product Team*  
*Next Review: After MVP Playtesting (June 2026)*