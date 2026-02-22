# AGENTS.md - Project Guide for AI Assistants

This document provides comprehensive information about the AAAAXY project to help AI coding assistants effectively understand and work with the codebase.

## Project Overview

**AAAAXY** is a nonlinear 2D puzzle platformer game that explores impossible spaces through non-Euclidean geometry. The game takes place across Klein Bottles, Möbius strips, and other mathematical surfaces, creating a unique gaming experience where topology is a core gameplay mechanic.

- **Genre**: Puzzle Platformer with Non-Euclidean Geometry
- **License**: Apache 2.0 (Open Source)
- **Primary Language**: Go (Golang)
- **Game Engine**: Ebitengine (v2.9.8)
- **Playtime**: 4-6 hours (new player), ~1 hour (experienced), ~15 minutes (speedrun)
- **Platforms**: Windows, macOS, Linux, Android, iOS, HTML5/WASM

## Technology Stack

### Core Technologies
- **Go 1.24.0+** - Primary programming language
- **Ebitengine** - 2D game engine for graphics and input handling
- **Oto v3** - Cross-platform audio playback
- **Kage** - Shader language for custom graphics effects
- **TMX/Tiled** - Level editor format for map design

### Key Libraries
- `github.com/hajimehoshi/ebiten/v2` - Game engine
- `github.com/lafriks/go-tiled` - TMX map parser
- `github.com/jeandeaual/go-locale` - Internationalization
- `github.com/leonelquinteros/gotext` - Translation framework
- `github.com/ncruces/zenity` - Native dialogs
- `github.com/lucasb-eyer/go-colorful` - Color manipulation

### Build System
- **Make** - Primary build automation
- **go generate** - Asset generation pipeline
- **Profile-Guided Optimization** - Performance optimization with default.pgo
- **Cross-compilation** - Single codebase for 7+ platforms

## Project Structure

```
aaaaxy-main/
├── main.go                      # Application entry point
├── go.mod / go.sum              # Go dependencies
├── Makefile                     # Build automation
├── default.pgo                  # Profile-guided optimization data
│
├── internal/                    # Internal Go packages (core codebase)
│   ├── aaaaxy/                 # Main game initialization and loop
│   ├── engine/                 # Core game engine
│   │   ├── world.go            # World state, entity management (1235 lines)
│   │   ├── renderer.go         # Rendering pipeline with visibility masking
│   │   ├── entity.go           # Entity system with incarnations
│   │   ├── trace*.go           # Collision detection and ray tracing
│   │   ├── polygon.go          # Geometric primitives
│   │   └── polyline.go         # Polyline handling
│   ├── game/                   # Game-specific entities and logic
│   │   ├── player/             # Player controller (583-line player.go)
│   │   ├── checkpoint/         # Save system and respawn points
│   │   ├── trigger/            # Trigger entities (switches, jump pads, etc.)
│   │   ├── target/             # Event targets (sound, music, sequences)
│   │   ├── riser/              # Moving platform system
│   │   └── mixins/             # Reusable entity behaviors
│   ├── level/                  # Level loading and management
│   │   ├── level.go            # TMX parsing, warp zones (1189 lines)
│   │   ├── tile.go             # Tile data structures
│   │   ├── entity.go           # Entity spawning
│   │   └── checkpoint_locations.go  # Checkpoint graph
│   ├── input/                  # Multi-platform input handling
│   ├── menu/                   # Menu system
│   ├── vfs/                    # Virtual file system (platform-aware)
│   ├── font/                   # Font rendering
│   ├── sound/                  # Sound system
│   ├── music/                  # Music playback
│   ├── m/                      # Fixed-point math utilities
│   ├── demo/                   # Demo recording/playback
│   ├── shader/                 # Shader management
│   ├── palette/                # Palette and dithering effects
│   └── [30+ other packages]
│
├── assets/                      # Game assets
│   ├── credits/                # Credits by category
│   ├── demos/                  # Demo recordings for regression testing
│   ├── locales/                # Translation files (.po/.pot)
│   ├── maps/                   # Tiled map files (level.tmx)
│   ├── shaders/                # Kage shader templates
│   ├── sounds/                 # OGG audio files (50+ sounds)
│   ├── sprites/                # PNG sprite images (360+ files)
│   ├── tiles/                  # Tile graphics
│   └── splash/                 # Loading screen assets
│
├── cmd/                         # Command-line utilities
│   ├── dumpcps/                # Checkpoint dumper
│   ├── dumpcplocs/             # Checkpoint location dumper
│   └── dumpluts/               # Lookup table generator
│
├── scripts/                     # Build and utility scripts (39 shell scripts)
├── third_party/                # Third-party assets and licenses
├── AndroidStudioProjects/      # Android build configuration
├── XcodeProjects/              # iOS build configuration
├── snap/                       # Snapcraft packaging for Linux
├── fastlane/                   # App store metadata
└── .github/workflows/          # CI/CD pipelines (13 workflows)
```

## Core Architecture

### Main Execution Flow

```
main.go
  → Parse command-line flags
  → Setup profiling (CPU, memory, mutex, block)
  → Initialize logging
  → aaaaxy.NewGame()
  → game.InitEbitengine() (early init)
  → game.InitFull() (full loading)
  → ebiten.RunGameWithOptions()
    ├─→ game.Update() (called at 60 FPS)
    │     ├─ Update input state
    │     ├─ Update world physics
    │     ├─ Update entities
    │     └─ Handle menu state
    └─→ game.Draw() (render frame)
          ├─ Render world with shaders
          ├─ Apply palette/dithering
          └─ Apply post-processing effects
```

### Entity System

The game uses an entity-component-like system:

- **Entities** (`internal/engine/entity.go`) - Base game objects
- **Incarnations** - Entity spawn configurations stored in levels
- **Mixins** (`internal/game/mixins/`) - Reusable behaviors:
  - `physics.go` - Physics simulation
  - `moving.go` - Moving platforms
  - `fadable.go` - Fade in/out effects
  - `track_player.go` - Player tracking behavior

### World and Warp Zones

The core innovation is the warp zone system (`internal/engine/world.go`):

- **Warp Zones** - Seamless portals connecting different parts of the level
- **Universal Cover** - Mathematical model for consistent non-Euclidean spaces
- **Visibility Masking** - Ensures objects appear only once on screen
- **Trace System** - Collision detection that works across warp zones

### Level Loading

Levels are created in Tiled Map Editor and loaded via TMX format:

1. Parse TMX file (`internal/level/level.go`)
2. Extract tiles and objects
3. Generate warp zone network
4. Validate checkpoint graph
5. Initialize entity spawners
6. Create world state

Level properties are extensively documented in `aaaaxy.tiled-project`.

### Fixed-Point Math

The game uses custom fixed-point arithmetic (`internal/m/fixed.go`) for deterministic physics:

- 12-bit fractional precision
- Ensures identical behavior across platforms
- Critical for demo replay and speedrunning
- Only tested code in the project (`fixed_test.go`)

## Key Subsystems

### Graphics Pipeline

1. **Rendering** (`internal/engine/renderer.go`)
   - Entity rendering with Z-ordering
   - Tile-based background rendering
   - Visibility masking shader

2. **Shaders** (`assets/shaders/*.kage`)
   - Visibility masking for warp zones
   - CRT effects
   - Scanline filters
   - Border stretching

3. **Palette System** (`internal/palette/`)
   - Multiple retro color modes (VGA, CGA, etc.)
   - Dither patterns (Bayer, checker, plastic, etc.)
   - Dynamic palette switching

### Input System

Multi-platform input support (`internal/input/`):
- Keyboard
- Mouse
- Gamepad (with SDL_GameControllerDB)
- Touch screen (with custom touch screen editor for mobile)
- Easter egg detection system

### Audio System

- **Sounds** (`internal/sound/`) - Sound effect playback
- **Music** (`internal/music/`) - Background music with transitions
- Uses Oto library for cross-platform audio

### Demo System

Unique testing approach using demo recordings (`internal/demo/`):
- Records entire game sessions frame-by-frame
- Replays with deterministic physics
- Used for regression testing instead of traditional unit tests
- Key demos:
  - `assets/demos/_anypercent.dem` - Any% speedrun
  - `assets/demos/benchmark.dem` - Performance benchmark

### Localization

Translation support for 10+ languages:
- GNU gettext format (.po files)
- Translation files in `assets/locales/`
- Transifex integration for community translations
- In-game validation of translations

## Build System

### Common Make Targets

```bash
make              # Build debug binary
make run          # Build and run
make vet          # Run linters and checks
make assets-update # Regenerate assets
make release      # Build release binary
make ziprelease   # Build with zipped assets
make embedrelease # Build with embedded assets
```

### Build Types

- **Debug** - Development build with debug features
- **Release** - Optimized release build
- **Zip Release** - Assets in separate .zip file
- **Embed Release** - Assets embedded in binary
- **FS Release** - Assets in filesystem

### Asset Generation

The build system generates several computed assets:
- Palette lookup tables (via `cmd/dumpluts/`)
- Checkpoint location graphs (via `cmd/dumpcplocs/`)
- Image load order optimization

## Testing Strategy

### Limited Traditional Testing
- Only one unit test file: `internal/m/fixed_test.go`
- Focus on testing the deterministic fixed-point math

### Demo-Based Regression Testing
- Entire game sessions recorded as demos
- Regression tests replay demos and verify success
- Scripts: `scripts/regression-test-demo.sh`, `scripts/rerecord-demo.sh`
- Ensures changes don't break core gameplay

### Validation at Load Time
- Extensive level validation when loading TMX files
- Checkpoint graph validation
- Warp zone consistency checks

### Static Analysis
- `make vet` runs multiple Go linters
- GitHub Actions run checks on every commit
- CodeQL security analysis

## Important Files to Understand

### Core Engine
- `internal/engine/world.go` (1235 lines) - The heart of the game, manages world state
- `internal/engine/entity.go` - Entity system foundation
- `internal/engine/renderer.go` - Graphics rendering
- `internal/engine/trace*.go` - Collision detection

### Game Logic
- `internal/game/player/player.go` (583 lines) - Player physics and abilities
- `internal/level/level.go` (1189 lines) - Level loading and warp zone setup
- `internal/aaaaxy/game.go` - Main game loop

### Entry Point
- `main.go` - Application startup and initialization

### Configuration
- `aaaaxy.tiled-project` - Complete entity and property definitions for level editor
- `go.mod` - Dependencies and Go version requirements
- `Makefile` - Build system and common commands

### Documentation
- `README.md` - Project overview and developer info
- `CONTRIBUTING.md` - Contribution guidelines
- `RELEASING.md` - Release process

## Unique Aspects

### Mathematical Innovation
- Implements non-Euclidean geometry in a 2D platformer
- Warp zones create seamless portals between spaces
- Objects only visible once on screen (no duplicate rendering)
- Universal cover mathematical model for consistency

### Technical Excellence
- Single codebase deploys to 7+ platforms
- Deterministic physics via fixed-point math
- Profile-guided optimization support
- Demo-based regression testing
- Extensive internationalization (10+ languages)

### Asset Management
- 360+ sprite files
- 50+ sound effects
- Multiple music tracks
- Procedurally generated lookup tables
- Optimized asset loading order

## Development Guidelines

### Code Organization
- All application code in `internal/` package
- Clear separation between engine and game logic
- Mixins for reusable entity behaviors
- Platform-specific code handled by build tags

### Performance Considerations
- Fixed-point math for determinism
- Profile-guided optimization enabled
- Asset loading optimization
- Efficient collision detection

### Multi-Platform Support
- Use `internal/vfs/` for file system operations
- Platform-specific code via build tags
- Test on target platforms
- Consider mobile input methods

### Adding New Entities
1. Create entity type in appropriate `internal/game/` subdirectory
2. Implement `Entity` interface methods
3. Add entity properties to `aaaaxy.tiled-project`
4. Register entity spawner in `internal/level/entity.go`
5. Test in Tiled Map Editor

### Modifying Levels
- Use Tiled Map Editor with `aaaaxy.tiled-project`
- Level file: `assets/maps/level.tmx`
- Follow existing entity property conventions
- Validate checkpoint graph after changes

### Working with Shaders
- Shader templates in `assets/shaders/*.kage`
- Kage is a subset of Go syntax
- Test across different palette modes
- Consider performance on mobile devices

## Common Debugging Tips

### Demo Playback Issues
- Demos rely on deterministic physics
- Any change to physics or timing can break demos
- Re-record demos after physics changes: `scripts/rerecord-demo.sh`

### Warp Zone Problems
- Validate checkpoint graph: `make dumpcplocs`
- Check for overlapping warp zones
- Ensure consistent tile connections
- Test player movement across boundaries

### Performance Issues
- Use built-in profiling flags: `-cpuprofile`, `-memprofile`
- Run timedemo: `scripts/timedemo.sh`
- Check asset loading order
- Profile shader performance

### Build Problems
- Ensure Go 1.24.0+ installed
- Run `go mod download` to fetch dependencies
- Check platform-specific requirements
- Use `make vet` to catch common issues

## CI/CD

The project has extensive CI/CD via GitHub Actions (13 workflows):
- Automated builds for all platforms
- Vet checks on every commit
- Translation validation
- Security analysis (CodeQL)
- Release automation

## Platform-Specific Notes

### Windows
- Uses zenity for native dialogs
- VFS stores state in AppData

### macOS
- App bundle creation via scripts
- VFS stores state in Application Support

### Linux
- Multiple distribution methods: Snap, Flatpak, AppImage
- Desktop file: `aaaaxy.desktop`
- AppStream metadata: `io.github.divverent.aaaaxy.metainfo.xml`

### Android
- Project in `AndroidStudioProjects/`
- Touch screen input with custom editor
- Special state storage location

### iOS
- Xcode project in `XcodeProjects/`
- Touch screen input support
- App Store metadata in `fastlane/`

### WebAssembly
- Special build process for browser deployment
- WASM-specific optimizations
- HTML wrapper: `aaaaxy.html`

## Resources

- **Repository**: https://github.com/divverent/aaaaxy
- **Website**: https://divverent.github.io/aaaaxy/
- **Translation**: https://www.transifex.com/aaaaxy/aaaaxy/
- **Ebitengine Docs**: https://ebitengine.org/
- **Tiled Editor**: https://www.mapeditor.org/

## Getting Started for AI Assistants

When helping with this project:

1. **Understand the non-Euclidean geometry** - Warp zones are the core innovation
2. **Respect deterministic physics** - Fixed-point math ensures replay consistency
3. **Test with demos** - Changes that break demos need careful consideration
4. **Consider all platforms** - Changes should work on desktop, mobile, and web
5. **Check the Tiled project file** - Entity properties are documented there
6. **Use the VFS** - Never access files directly, always use `internal/vfs/`
7. **Follow Go conventions** - Code should pass `make vet`
8. **Maintain localization** - Update translation templates when changing UI text

This project is a sophisticated example of Go game development with excellent engineering practices, innovative mathematical concepts, and comprehensive multi-platform support.
