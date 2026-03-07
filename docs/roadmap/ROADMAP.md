# Development Roadmap: Medieval Village Simulator

## Overview
This roadmap outlines the development journey from concept to launch and beyond. The project follows agile methodology with 2-week sprints, focusing on delivering a playable MVP within 6 months, then expanding based on player feedback.

## Timeline Summary
```
Q2 2026 (Apr-Jun): Discovery & Foundation
├── Sprint 1-2: Technical spike & prototype
├── Sprint 3-4: Core simulation engine
└── Sprint 5-6: Basic frontend integration

Q3 2026 (Jul-Sep): MVP Development
├── Sprint 7-8: Resident management system
├── Sprint 9-10: Resource chains & production
├── Sprint 11-12: Seasons & events
└── Sprint 13-14: Polish & alpha testing

Q4 2026 (Oct-Dec): Launch & Iteration
├── Sprint 15-16: Beta testing & feedback
├── Sprint 17-18: Bug fixes & performance
├── Launch: Early Access (December 2026)
└── Sprint 19-20: Post-launch support

2027+: Expansion
├── Q1 2027: Social systems & modding
├── Q2 2027: Advanced economics
└── H2 2027: Multi-village & historical events
```

## Phase 1: Discovery & Foundation (Q2 2026 - 8 weeks)

### Sprint 1-2: Technical Spike & Prototype (April 2026)
**Goal:** Validate technical feasibility and establish development environment

#### Deliverables:
- [ ] **Technical Proof of Concept:** Basic deterministic simulation in Go
- [ ] **Architecture Documentation:** Detailed system design
- [ ] **Development Environment:** Docker setup, CI/CD pipeline
- [ ] **Technology Decisions:** Finalize libraries, frameworks, tools
- [ ] **Risk Assessment:** Identify major technical challenges

#### Key Activities:
1. Create barebones simulation engine that processes basic resource production
2. Evaluate Go libraries for game development (ECS frameworks, etc.)
3. Set up PostgreSQL 18 with test data schema
4. Create Angular 21 skeleton with basic routing
5. Establish coding standards and review process

### Sprint 3-4: Core Simulation Engine (May 2026)
**Goal:** Build deterministic foundation for all game systems

#### Deliverables:
- [ ] **Simulation Core:** Time-based event processing system
- [ ] **Game State Model:** Resident, building, resource entities
- [ ] **Deterministic Random:** Seed-based RNG with serialization
- [ ] **Unit Test Suite:** 80%+ coverage of simulation logic
- [ ] **Performance Benchmarks:** Baseline for 50-resident village

#### Key Activities:
1. Implement core game loop with weekly turn processing
2. Create entity system for residents with skills and needs
3. Build deterministic random number generator
4. Design save/load system architecture
5. Create comprehensive test suite

### Sprint 5-6: Basic Frontend Integration (June 2026)
**Goal:** Connect frontend to backend with minimal viable interface

#### Deliverables:
- [ ] **REST API:** Complete endpoint specification
- [ ] **Frontend-Backend Integration:** Basic game load/save
- [ ] **Minimal UI:** Angular components for core views
- [ ] **Development Workflow:** Hot reload, debugging tools
- [ ] **End-to-End Tests:** Basic game flow automation

#### Key Activities:
1. Implement REST API for game state management
2. Create Angular services for API communication
3. Build basic village map visualization
4. Implement resident list with drag-drop assignment
5. Create turn control (plan/execute/review cycle)

## Phase 2: MVP Development (Q3 2026 - 14 weeks)

### Sprint 7-8: Resident Management System (July 2026)
**Goal:** Complete system for individual villager simulation

#### Deliverables:
- [ ] **Resident AI:** Needs calculation, skill development
- [ ] **Relationship System:** Family ties, friendships, conflicts
- [ ] **Assignment UI:** Intuitive drag-drop interface
- [ ] **Resident Detail View:** Skills, needs, relationships display
- [ ] **Life Cycle:** Birth, aging, death with inheritance

#### Key Activities:
1. Implement resident need system (food, shelter, happiness)
2. Create skill progression through work experience
3. Build relationship network with mutual effects
4. Design UI for resident management
5. Implement life events (birth, marriage, death)

### Sprint 9-10: Resource Chains & Production (July-August 2026)
**Goal:** Implement 15-resource economy with production chains

#### Deliverables:
- [ ] **Resource System:** 15 resources with properties
- [ ] **Production Chains:** Raw → processed → advanced goods
- [ ] **Building System:** Workplace assignment, production rates
- [ ] **Storage Management:** Inventory with spoilage (food)
- [ ] **Construction:** Building placement requiring materials/labor

#### Key Activities:
1. Define all 15 resources and their relationships
2. Implement production recipes and requirements
3. Create building system with worker slots
4. Build storage system with capacity limits
5. Implement construction interface

### Sprint 11-12: Seasons & Events (August 2026)
**Goal:** Add environmental factors and random events

#### Deliverables:
- [ ] **Seasonal System:** Weather, temperature, growing seasons
- [ ] **Event System:** Disease, accidents, good/bad harvests
- [ ] **Agricultural Cycle:** Planting, growing, harvesting
- [ ] **Visual Feedback:** Season effects on map appearance
- [ ] **Event History:** Log of past occurrences

#### Key Activities:
1. Implement seasonal weather patterns
2. Create event triggering system
3. Build agricultural calendar
4. Add visual season indicators
5. Create event history interface

### Sprint 13-14: Polish & Alpha Testing (September 2026)
**Goal:** Refine gameplay and conduct first external tests

#### Deliverables:
- [ ] **Alpha Release:** Playable end-to-end experience
- [ ] **Tutorial System:** Guided introduction for new players
- [ ] **UI Polish:** Improved visuals, animations, feedback
- [ ] **Balance Pass:** Tuned resource values, event frequencies
- [ ] **Bug Fixes:** Critical issues from internal testing

#### Key Activities:
1. Create interactive tutorial
2. Polish UI with animations and transitions
3. Balance game economy through playtesting
4. Fix critical bugs and performance issues
5. Prepare alpha release for external testers

## Phase 3: Launch & Iteration (Q4 2026 - 12 weeks)

### Sprint 15-16: Beta Testing & Feedback (October 2026)
**Goal:** Gather player feedback and iterate on core systems

#### Deliverables:
- [ ] **Beta Release:** Public testing version
- [ ] **Analytics Integration:** Player behavior tracking
- [ ] **Feedback Collection:** In-game survey system
- [ ] **Community Channels:** Discord, forums, bug reporting
- [ ] **Prioritized Backlog:** Features and fixes from feedback

#### Key Activities:
1. Launch beta to strategy game communities
2. Implement analytics for player behavior
3. Create feedback collection system
4. Set up community communication channels
5. Analyze feedback and update roadmap

### Sprint 17-18: Bug Fixes & Performance (November 2026)
**Goal:** Stabilize game and optimize for launch

#### Deliverables:
- [ ] **Performance Optimization:** Turn processing < 3 seconds
- [ ] **Bug Squashing:** Zero critical bugs
- [ ] **Browser Compatibility:** All major browsers supported
- [ ] **Accessibility Improvements:** WCAG 2.1 AA compliance
- [ ] **Localization Foundation:** i18n system for future translations

#### Key Activities:
1. Profile and optimize simulation performance
2. Fix all critical and high-priority bugs
3. Test on all target browsers
4. Implement accessibility features
5. Prepare localization framework

### Sprint 19-20: Launch & Post-Launch Support (December 2026)
**Goal:** Release Early Access and establish support processes

#### Deliverables:
- [ ] **Early Access Launch:** Live on project website
- [ ] **Documentation:** Player guide, modding docs
- [ ] **Support System:** Bug tracking, player help
- [ ] **Launch Marketing:** Announcements, trailers, press kits
- [ ] **Monitoring:** Server performance, error tracking

#### Key Activities:
1. Final launch preparation and testing
2. Create player documentation
3. Set up support channels
4. Execute marketing plan
5. Monitor launch and respond to issues

## Phase 4: Expansion (2027+)

### Q1 2027: Social Systems & Modding
**Focus:** Deepen simulation and enable community creation

#### Planned Features:
- Social classes and hierarchies
- Crime and justice system
- Religion and ceremonies
- Modding API and tools
- Steam Workshop integration

### Q2 2027: Advanced Economics
**Focus:** Expand economic simulation depth

#### Planned Features:
- Trade routes and caravans
- Currency and banking
- Taxation and feudal obligations
- Market price dynamics
- Crafting specialization

### H2 2027: Multi-Village & Historical Events
**Focus:** Scale simulation and add historical context

#### Planned Features:
- Multiple interacting villages
- Historical events (Black Death, Little Ice Age)
- Technology progression across centuries
- Asynchronous multiplayer
- Historical commentary mode

## Success Metrics & Milestones

### Alpha (September 2026)
- ✅ 50 internal testers complete 20+ game weeks
- ✅ 80%+ positive feedback on core gameplay loop
- ✅ Turn processing < 5 seconds for 50-resident village
- ✅ No game-breaking bugs in 10+ hours of play

### Beta (October 2026)
- ✅ 500+ external testers
- ✅ Average session length > 60 minutes
- ✅ 70% retention after first week
- ✅ NPS > 20 from strategy game enthusiasts

### Early Access Launch (December 2026)
- ✅ 1000+ active players
- ✅ 75% positive reviews on launch platforms
- ✅ < 1% crash rate
- ✅ Community mods appearing within first month

### Version 1.0 (Q2 2027)
- ✅ All MVP features complete and polished
- ✅ Modding ecosystem with 50+ community mods
- ✅ 5000+ active players
- ✅ Revenue covering ongoing development

## Resource Planning

### Development Team (MVP Phase)
- **Backend Developer (Go):** Simulation engine, database, API
- **Frontend Developer (Angular):** UI, visualization, UX
- **Game Designer:** Systems design, balance, content
- **Part-time:** Art/assets, testing, community management

### Infrastructure (Initial)
- **Web Server:** Single instance capable of 10,000+ games
- **Database:** PostgreSQL 18 with regular backups
- **Storage:** S3-compatible for save games and assets
- **CDN:** For frontend asset delivery
- **Monitoring:** Error tracking, performance metrics

### Budget Estimate (6 months to MVP)
- **Development:** $120,000 (2.5 FTEs)
- **Infrastructure:** $2,000/month
- **Services/Tools:** $1,000/month
- **Contingency (20%):** $25,000
- **Total MVP Budget:** ~$150,000

## Risk Management

### High Priority Risks
1. **Simulation Complexity:** Break into modular subsystems with clear interfaces
2. **Performance Issues:** Early profiling, algorithmic optimization focus
3. **Player Comprehension:** Extensive tutorial and progressive disclosure
4. **Scope Creep:** Strict MVP definition, post-MVP features in roadmap

### Mitigation Strategies
- Weekly playtesting from sprint 3 onward
- Modular architecture allowing feature removal if needed
- Community involvement from beta phase
- Regular roadmap reviews and adjustments

## Dependencies & Assumptions

### External Dependencies
- Angular 21 stability and long-term support
- PostgreSQL 18 feature compatibility
- Browser vendors maintaining Web Workers/WebAssembly support
- Strategy game community interest and engagement

### Internal Assumptions
- Team familiarity with Go and Angular
- Historical research resources available
- Player willingness to learn complex systems
- Modding community development with proper tools

---

*Roadmap Version 1.0 · Created: March 7, 2026 · Next Review: After Alpha Testing (September 2026)*  
*This is a living document updated quarterly based on development progress and player feedback.*