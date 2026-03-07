# Save System & Modding Support

**Vision:** Create a robust foundation for player creativity and long-term engagement through reliable save games and extensive modding capabilities that transform a single game into a platform for community content.

**Status:** Draft  
**Created:** March 7, 2026  
**Priority:** P0 (MVP Longevity)

---

## Executive Summary

The save system ensures players can invest hundreds of hours into a village knowing their progress is safe, with version compatibility, corruption recovery, and cloud backup options. The modding system transforms the game from a fixed experience into a platform where the community can add new resources, events, buildings, and even game mechanics, extending the game's lifespan indefinitely through user-generated content.

**Success Definition:** Players share elaborate save game stories ("My 100-year village survived three plagues!") and modders create content that becomes essential to the community, with the most popular mods downloaded thousands of times.

---

## User Stories

### US-001: Reliable Save Games
**As:** A player investing dozens of hours into a village  
**I want:** My game progress saved automatically and reliably  
**So that:** I never lose significant progress to crashes or mistakes

**Acceptance Criteria:**
- [ ] Automatic saving after each game week
- [ ] Manual save slots (10+ slots per village)
- [ ] Save file integrity checking on load
- [ ] Corruption recovery (load backup if primary save corrupt)
- [ ] Save file versioning with migration for old versions
- [ ] Save file size optimization (<5MB for 100-year village)

### US-002: Save Game Management
**As:** A player with multiple villages  
**I want:** To organize, browse, and manage my save games  
**So that:** I can easily find and continue specific games

**Acceptance Criteria:**
- [ ] Save game browser with village preview (map, population, date)
- [ ] Sort/filter by: village name, date, population, play time
- [ ] Save metadata: play time, real-world date saved, version
- [ ] Delete, rename, duplicate save games
- [ ] Import/export save files for sharing
- [ ] Cloud save synchronization (optional)

### US-003: Mod Installation & Management
**As:** A player wanting to enhance my game  
**I want:** To easily install and manage mods without technical knowledge  
**So that:** I can customize my gameplay experience

**Acceptance Criteria:**
- [ ] In-game mod browser/downloader
- [ ] One-click mod installation/activation
- [ ] Mod dependency resolution (mod A requires mod B)
- [ ] Load order management for conflicting mods
- [ ] Mod conflict detection and warnings
- [ ] Mod disabling without uninstallation

### US-004: Mod Creation Tools
**As:** A community member wanting to create content  
**I want:** Accessible tools for creating mods without programming  
**So that:** I can share my ideas with other players

**Acceptance Criteria:**
- [ ] JSON-based configuration for: resources, buildings, events, recipes
- [ ] Mod template system for common mod types
- [ ] Visual editor for simple mods (event editor, resource editor)
- [ ] Mod validation tools to catch errors before sharing
- [ ] Documentation and examples for mod creators
- [ ] Testing mode for mods in isolated environment

### US-005: Save Game Sharing
**As:** A player wanting to share my village story  
**I want:** To export my village with metadata and story  
**So that:** Others can experience my village's history

**Acceptance Criteria:**
- [ ] Export save with optional story/notes
- [ ] Screenshot embedding in save file
- [ ] Event history preservation for shared games
- [ ] Mod compatibility checking for shared games
- [ ] Privacy options (hide resident names, remove personal notes)
- [ ] Share to community hub with ratings/comments

### US-006: Version Compatibility
**As:** A player updating the game  
**I want:** My old save games to work with new versions  
**So that:** I can continue playing without starting over

**Acceptance Criteria:**
- [ ] Save file version detection
- [ ] Automatic migration for old save formats
- [ ] Migration report showing changes applied
- [ ] Option to keep backup of original save
- [ ] Graceful handling of missing mods in updated saves
- [ ] Warning if migration may cause issues

---

## Technical Specifications

### Save File Format
```go
// Save file structure
type SaveFile struct {
    Header      SaveHeader     `json:"header"`
    GameState   GameState      `json:"gameState"`
    Metadata    SaveMetadata   `json:"metadata"`
    ModData     []ModSaveData  `json:"modData,omitempty"`
    Checksum    string         `json:"checksum"` // For integrity verification
}

type SaveHeader struct {
    Version     string    `json:"version"`     // Game version "1.0.0"
    Format      int       `json:"format"`      // Save format version
    Created     time.Time `json:"created"`     // Real-world creation time
    GameTime    GameDate  `json:"gameTime"`    // In-game date
    VillageName string    `json:"villageName"` // Player-given name
}

type SaveMetadata struct {
    PlayTime    time.Duration `json:"playTime"`    // Total play time
    RealDays    int           `json:"realDays"`    // Real days since start
    Screenshots []string      `json:"screenshots"` // Base64 encoded thumbnails
    PlayerNotes string        `json:"playerNotes"` // Player's notes about village
    Statistics  GameStats     `json:"statistics"`  // Summary stats
}

// Game state (simplified example)
type GameState struct {
    // Deterministic seed for recreation
    Seed        int64         `json:"seed"`
    
    // Core simulation state
    Calendar    Calendar      `json:"calendar"`
    Village     VillageState  `json:"village"`
    Residents   []Resident    `json:"residents"`
    Resources   ResourceMap   `json:"resources"`
    Buildings   []Building    `json:"buildings"`
    
    // Event and decision history
    EventLog    []EventRecord `json:"eventLog"`
    Decisions   []Decision    `json:"decisions"`
    
    // Mod-specific data
    ModStates   map[string]interface{} `json:"modStates"`
}
```

### Save File Optimization
```go
// Delta compression for frequent saves
func CreateDeltaSave(previous *SaveFile, current *GameState) DeltaSave {
    delta := DeltaSave{
        BaseVersion: previous.Header.Version,
        Changes: make(map[string]interface{}),
    }
    
    // Only save changes from previous state
    if !reflect.DeepEqual(previous.GameState.Residents, current.Residents) {
        // Save only changed residents
        delta.Changes["residents"] = calculateResidentDeltas(
            previous.GameState.Residents, 
            current.Residents,
        )
    }
    
    // Similar for resources, buildings, etc.
    
    return delta
}

// Apply delta to reconstruct full state
func ApplyDelta(base *SaveFile, delta DeltaSave) (*GameState, error) {
    state := base.GameState.DeepCopy()
    
    for key, change := range delta.Changes {
        switch key {
        case "residents":
            applyResidentDeltas(&state.Residents, change)
        case "resources":
            applyResourceDeltas(&state.Resources, change)
        // ... other state sections
        }
    }
    
    return state, nil
}
```

### Mod System Architecture
```go
// Mod definition file (mod.json)
type ModDefinition struct {
    ID          string            `json:"id"`          // "author.modname"
    Name        string            `json:"name"`        // "Medieval Furniture Pack"
    Version     string            `json:"version"`     // "1.2.0"
    Author      string            `json:"author"`
    Description string            `json:"description"`
    
    // Dependencies
    Requires    []Dependency      `json:"requires,omitempty"`
    Conflicts   []string          `json:"conflicts,omitempty"` // Mod IDs
    
    // Content definitions
    Resources   []ResourceDef     `json:"resources,omitempty"`
    Buildings   []BuildingDef     `json:"buildings,omitempty"`
    Events      []EventDef        `json:"events,omitempty"`
    Recipes     []RecipeDef       `json:"recipes,omitempty"`
    
    // Gameplay changes
    Config      ModConfig         `json:"config,omitempty"`
    Scripts     []ScriptDef       `json:"scripts,omitempty"` // JavaScript extensions
    
    // Assets
    Assets      AssetManifest     `json:"assets,omitempty"`
}

// Mod loading and integration
type ModManager struct {
    activeMods  map[string]*LoadedMod
    loadOrder   []string
    hooks       map[HookType][]ModHook
}

func (m *ModManager) LoadMod(path string) error {
    // Read mod definition
    def, err := readModDefinition(path)
    if err != nil {
        return err
    }
    
    // Check dependencies
    for _, dep := range def.Requires {
        if !m.IsModLoaded(dep.ID) {
            return fmt.Errorf("missing dependency: %s", dep.ID)
        }
    }
    
    // Load mod assets
    assets, err := loadModAssets(path, def.Assets)
    if err != nil {
        return err
    }
    
    // Register mod hooks
    for _, hook := range def.Hooks {
        m.RegisterHook(hook)
    }
    
    // Add to active mods
    m.activeMods[def.ID] = &LoadedMod{
        Definition: def,
        Assets: assets,
        State: make(map[string]interface{}),
    }
    
    // Recalculate load order
    m.UpdateLoadOrder()
    
    return nil
}
```

### Mod Hook System
```go
// Hook points where mods can inject behavior
type HookType string

const (
    HookGameStart      HookType = "game_start"
    HookWeekStart      HookType = "week_start"
    HookWeekEnd        HookType = "week_end"
    HookEventTrigger   HookType = "event_trigger"
    HookProductionCalc HookType = "production_calculate"
    HookNeedUpdate     HookType = "need_update"
    HookUIRender       HookType = "ui_render"
)

// Mod hook definition
type ModHook struct {
    ModID      string                 `json:"modId"`
    HookType   HookType               `json:"hookType"`
    Priority   int                    `json:"priority"` // Execution order
    Script     string                 `json:"script"`   // JavaScript code
    Condition  string                 `json:"condition,omitempty"`
}

// Hook execution
func ExecuteHooks(hookType HookType, context HookContext) (HookContext, error) {
    hooks := modManager.GetHooks(hookType)
    
    // Sort by priority
    sort.Slice(hooks, func(i, j int) bool {
        return hooks[i].Priority < hooks[j].Priority
    })
    
    // Execute each hook
    result := context
    for _, hook := range hooks {
        if hook.Condition != "" {
            if !EvaluateCondition(hook.Condition, result) {
                continue
            }
        }
        
        newResult, err := ExecuteScript(hook.Script, result)
        if err != nil {
            return result, fmt.Errorf("mod %s hook failed: %v", hook.ModID, err)
        }
        
        result = newResult
    }
    
    return result, nil
}
```

### Save File Migration
```go
// Migration system for version updates
type Migration struct {
    FromVersion string
    ToVersion   string
    Steps       []MigrationStep
}

type MigrationStep struct {
    Description string
    Apply       func(*SaveFile) error
    Rollback    func(*SaveFile) error // For backup restoration
}

func MigrateSave(save *SaveFile, targetVersion string) (*SaveFile, error) {
    current := save.Header.Version
    
    // Find migration path
    path, err := findMigrationPath(current, targetVersion)
    if err != nil {
        return nil, err
    }
    
    // Create backup
    backup := save.DeepCopy()
    backupPath := fmt.Sprintf("%s.backup.%s", savePath, current)
    
    // Apply each migration step
    for _, migration := range path {
        log.Printf("Migrating from %s to %s", migration.FromVersion, migration.ToVersion)
        
        for _, step := range migration.Steps {
            if err := step.Apply(save); err != nil {
                // Rollback on failure
                log.Printf("Migration failed: %v, rolling back", err)
                save = backup
                return nil, err
            }
        }
        
        // Update version in save file
        save.Header.Version = migration.ToVersion
    }
    
    // Save backup
    if err := saveBackup(backup, backupPath); err != nil {
        log.Printf("Warning: failed to save backup: %v", err)
    }
    
    return save, nil
}
```

### UI Interface Requirements

#### Save Game Browser
- **Grid/List view:** Save files with village preview image
- **Sorting:** Date, village name, population, play time
- **Filtering:** By mods used, by version, by tags
- **Quick actions:** Load, delete, duplicate, export
- **Search:** Village name, player notes, resident names

#### Mod Manager Interface
- **Available mods:** Browse local and online mods
- **Mod details:** Description, version, dependencies, screenshots
- **Installation status:** Installed, enabled, needs update
- **Conflict resolution:** Visual conflict detection and resolution
- **Load order editor:** Drag-drop ordering with dependency validation

#### Save Game Details
- **Village snapshot:** Map overview, key statistics
- **Timeline:** Major events plotted on timeline
- **Family trees:** Resident relationships visualization
- **Resource history:** Graphs of resource levels over time
- **Player journal:** Notes attached to specific dates

#### Mod Creation Tools
- **Template wizard:** Step-by-step mod creation for common types
- **Resource editor:** Visual editing of resource properties
- **Event editor:** Flowchart-style event creation
- **Recipe editor:** Drag-drop production chain creation
- **Test environment:** Isolated testing of mods without affecting saves

---

## Integration Points

### With Game State System
- Save system serializes complete game state
- Load system reconstructs deterministic simulation
- Migration system updates old game states
- Delta compression reduces save file size

### With Mod System
- Mod definitions extend game data structures
- Mod hooks inject behavior at key points
- Mod assets loaded on demand
- Mod state saved within save files

### With UI System
- Save browser integrates with main menu
- Mod manager accessible from game settings
- In-game mod configuration panels
- Mod conflict warnings during gameplay

### With Community Features
- Save file sharing to community hub
- Mod distribution through in-game browser
- Mod ratings and comments
- Automatic update checking for mods

---

## Performance Considerations

### Save File Size Optimization
| Village Size | Full Save | Delta Save | Compression |
|--------------|-----------|------------|-------------|
| 10 residents | ~500 KB   | ~50 KB     | ~100 KB (gzip) |
| 50 residents | ~2 MB     | ~200 KB    | ~400 KB (gzip) |
| 100 residents| ~5 MB     | ~500 KB    | ~1 MB (gzip)   |

### Save/Load Performance Targets
| Operation | Target Time | Notes |
|-----------|-------------|-------|
| Full save | < 500ms | Includes compression |
| Delta save | < 100ms | Only changed data |
| Full load | < 1s | With decompression |
| Delta load | < 200ms | Apply to cached state |
| Migration | < 2s | Per version step |

### Mod System Performance
- **Mod loading:** < 50ms per mod (after first load)
- **Hook execution:** < 5ms per hook type per week
- **Asset loading:** Lazy loading with caching
- **Conflict detection:** < 100ms for full mod set

---

## Security & Safety

### Save File Validation
- Checksum verification on load
- Schema validation against expected structure
- Size limits to prevent memory exhaustion
- Malformed data isolation and recovery

### Mod Security Sandbox
- JavaScript hooks run in isolated VM
- File system access restrictions
- Network access disabled for mods
- Resource limits on mod execution
- Mod signing for verified authors

### User Data Protection
- Local save encryption option
- Cloud save encryption mandatory
- Personal data not included in shared saves
- Mod privacy settings (what data mods can access)

---

## Testing Strategy

### Unit Tests
- Save file serialization/deserialization
- Delta compression correctness
- Migration algorithm accuracy
- Mod dependency resolution
- Hook execution order

### Integration Tests
- Full save/load cycle with all game systems
- Mod installation/enabling/disabling
- Version migration of complex save files
- Mod conflict detection and resolution
- Cloud save synchronization

### Playtesting Focus
- Save reliability over long play sessions
- Mod discoverability and installation ease
- Performance impact of mods on gameplay
- User understanding of mod conflicts
- Satisfaction with sharing features

---

## Risks & Mitigations

### Technical Risks
1. **Save file corruption losing player progress**  
   **Mitigation:** Multiple backup system, checksum validation, corruption recovery

2. **Mod conflicts causing game instability**  
   **Mitigation:** Conflict detection, load order management, mod sandboxing

3. **Save file size explosion with long games**  
   **Mitigation:** Delta compression, history truncation, optional detail levels

4. **Version migration breaking old saves**  
   **Mitigation:** Extensive migration testing, backup preservation, rollback capability

### Design Risks
1. **Modding tools too complex for casual creators**  
   **Mitigation:** Template system, visual editors, example mods, tutorial videos

2. **Save system feeling intrusive (too many autosaves)**  
   **Mitigation:** Configurable autosave frequency, quiet saves, manual control

3. **Mod distribution leading to quality control issues**  
   **Mitigation:** Rating system, moderation, verified author program, quality badges

4. **Community fragmentation from incompatible mods**  
   **Mitigation:** Mod compatibility tagging, dependency system, conflict warnings

### Mitigation Strategies
- **Early save system implementation:** Test save/load from week 1 of development
- **Modding support from start:** Design game systems with modding in mind
- **Community involvement:** Early access to mod tools for dedicated players
- **Progressive complexity:** Basic modding first, advanced features later

---

## Success Metrics

### Technical Metrics
- [ ] Save file corruption rate < 0.1%
- [ ] Save/load success rate > 99.9%
- [ ] Migration success rate > 99% for all version pairs
- [ ] Mod loading error rate < 1%

### User Experience Metrics
- [ ] 95% of players understand save system within first session
- [ ] Average mod installation time < 2 minutes
- [ ] Mod conflict resolution success rate > 80%
- [ ] Player satisfaction with mod system > 4/5 stars

### Community Metrics
- [ ] 50+ community mods within 3 months of launch
- [ ] Top mods downloaded > 10,000 times
- [ ] Average mod rating > 4/5 stars
- [ ] Active mod creators > 100 within first year

---

## Dependencies

### Required First
- Complete game state serialization
- Basic file I/O system
- Mod definition format
- UI framework for mod manager

### Dependent Features
- Cloud save synchronization (requires save system)
- Mod distribution platform (requires mod system)
- Advanced modding tools (requires basic mod system)
- Save file sharing community (requires save export)

---

## Open Questions

1. **Mod scripting language:** JavaScript vs. Lua vs. custom DSL?
2. **Monetization:** Allow mod creators to charge for mods?
3. **Cross-platform saves:** Compatibility between web and potential future desktop version?
4. **Save file portability:** Should saves work between different mod configurations?

---

*Feature Version 1.0 · Owner: Platform Team · Estimated Effort: 6-8 sprints*