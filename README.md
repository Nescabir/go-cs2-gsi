# Go CS2 GSI

[![Go Reference](https://pkg.go.dev/badge/github.com/nescabir/go-cs2-gsi.svg)](https://pkg.go.dev/github.com/nescabir/go-cs2-gsi)

A high-performance Go library for handling Counter-Strike 2 Game State Integration (GSI) with type-safe event handling and real-time game data processing.

## 🎯 Features

- **Type-safe Event System**: Subscribe to specific game events with compile-time type safety
- **Real-time Game Data**: Process live game state updates from CS2
- **HTTP Server**: Built-in HTTP server using Go 1.22+ ServeMux for handling GSI requests
- **Comprehensive Game Models**: Complete data structures for all CS2 game state information
- **Configurable**: Customizable server settings, round limits, and logging levels
- **Production Ready**: Proper error handling, validation, and security measures

## 🚀 Quick Start

### Installation

```bash
go get github.com/nescabir/go-cs2-gsi
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "log/slog"

    cs2gsi "github.com/nescabir/go-cs2-gsi"
    "github.com/nescabir/go-cs2-gsi/models"
)

func main() {
    // Create GSI instance with default configuration
    gsi := cs2gsi.New(cs2gsi.NewConfig())

    // Subscribe to events
    cs2gsi.Subscribe(cs2gsi.Mvp, func(event cs2gsi.Event[*models.Player]) {
      fmt.Printf("MVP: %s with %d kills (%d HS)\n",
        event.Data.Name, event.Data.State.Round_kills, event.Data.State.Round_killhs)
    })

    cs2gsi.Subscribe(cs2gsi.RoundEnd, func(event cs2gsi.Event[*models.Score]) {
        fmt.Printf("Round ended! Winner: %s\n", event.Data.Winner.Name)
    })

    // Start the server
    if err := gsi.Listen(); err != nil {
        log.Fatal(err)
    }
}
```

## 📋 Configuration

The library provides flexible configuration options:

```go
config := cs2gsi.Config{
    ServerAddr:          ":3000",         // HTTP server address ex: 127.0.0.1:3000, localhost:4242
    RegulationMaxRounds: 12,              // Rounds per half (CS2 MR12); use 15 for legacy MR15
    OvertimeMaxRounds:   3,               // Max rounds in overtime
    LogLevel:            slog.LevelInfo,  // Logging level
    ExpectedToken:       "",              // Optional: reject POSTs when auth token mismatch
    PlayerExtensions:    nil,             // Optional cloud player metadata merge
    TeamExtensions: cs2gsi.TeamExtensionsConfig{
        Left:  nil,
        Right: nil,
    },
}

gsi := cs2gsi.New(config)
```

## 🎮 Available Events

The library provides type-safe events for all major game occurrences:

### Core Events

- `Data` - Raw game state updates
- `RoundEnd` - Round completion with winner information
- `MatchEnd` - Match completion

### Combat Events

- `Mvp` - MVP player selection
- `Kill` - Player kills (via `DigestMIRV` + HLAE game events)
- `Hurt` - Player damage events (via `DigestMIRV` + HLAE game events)

### Other Events

- `Raw` - Raw JSON payload before parsing (mirrors csgogsi `raw` event)

### Round Events

- `FreezetimeStart/End` - Freeze time beginning/ending
- `IntermissionStart/End` - Intermission periods
- `TimeoutStart/End` - Team timeouts

### Bomb Events

- `BombPlantStart/Stop` - Bomb planting initiation/cancellation
- `BombPlanted` - Bomb successfully planted
- `BombDefused` - Bomb defused
- `BombExploded` - Bomb explosion
- `DefuseStart/End` - Defuse initiation/completion

## 🔧 CS2 Setup

1. **Copy the configuration template**:

   Copy the configuration template to your game's cfg folder (`steamapps/common/Counter-Strike Global Offensive/game/core/cfg/`)

2. **Configure the GSI file**:

   - Update the `uri` to match your server address
   - Modify the `token` for authentication (No token validation yet)
   - Enable/disable specific data feeds as needed

3. **Start your Go application** and launch CS2

## 📊 Data Models

The library provides comprehensive data structures for all CS2 game information:

### Player Information

- Steam ID, name, clan, team
- Health, armor, money, equipment
- Position, weapons, match statistics
- Activity status

### Game State

- Map information and phase
- Round details and outcomes
- Team scores and statistics
- Bomb state and position
- Grenade positions and effects

### Match Data

- Round history and outcomes
- Player damage tracking
- Weapon information and states
- Observer data

## 🛠️ Advanced Usage

### Programmatic digest (without HTTP)

```go
gsi := cs2gsi.New(cs2gsi.NewConfig())

if err := gsi.Digest(rawJSON); err != nil {
    log.Fatal(err)
}

state := gsi.Snapshot()
fmt.Println(state.Map.Name)
```

### MIRV / HLAE kill and hurt events

Kill and hurt events come from HLAE game event payloads, not standard GSI POST bodies. Call `DigestMIRV` on a separate feed after at least one successful `Digest`:

```go
cs2gsi.Subscribe(cs2gsi.Kill, func(event cs2gsi.Event[*models.KillEvent]) {
    fmt.Printf("%s killed %s with %s\n",
        event.Data.Attacker.Name,
        event.Data.Victim.Name,
        event.Data.Weapon.Name)
})

result, err := gsi.DigestMIRV(mirvJSON, cs2gsi.MIRVEventPlayerDeath)
```

### Raw payload hook

Subscribe to the raw JSON body before parsing:

```go
cs2gsi.Subscribe(cs2gsi.Raw, func(event cs2gsi.Event[[]byte]) {
    // log, forward, or custom-parse the payload
})
```

### Player and team extensions

Merge cloud metadata (avatars, custom names, team logos) like csgogsi:

```go
gsi := cs2gsi.New(cs2gsi.Config{
    PlayerExtensions: []models.PlayerExtension{{
        SteamId: "76561198000000001",
        Name:    "Broadcast Name",
        Avatar:  "https://example.com/avatar.png",
    }},
    TeamExtensions: cs2gsi.TeamExtensionsConfig{
        Left: &models.TeamExtension{Logo: "https://example.com/logo.png"},
    },
})
```

Parsed state also includes:
- `state.Previously` / `state.Added` — shallow GSI delta blocks
- `state.Damage` — per-round damage history used for ADR
- `phase_countdowns.timeout_team` — team that called a timeout
- `player.DefaultName` — raw GSI name before extension override

### Custom Event Handling

```go
// Subscribe to multiple events
cs2gsi.Subscribe(cs2gsi.BombPlanted, func(event cs2gsi.Event[*models.Player]) {
    fmt.Printf("Bomb planted by %s at site %s\n",
        event.Data.Name,
        event.Data.Team.Side)
})

cs2gsi.Subscribe(cs2gsi.Mvp, func(event cs2gsi.Event[*models.Player]) {
    fmt.Printf("MVP: %s with %d kills (%d headshots)\n",
        event.Data.Name,
        event.Data.State.Round_kills,
        event.Data.State.Round_killhs)
})
```

### Error Handling

```go
gsi := cs2gsi.New(cs2gsi.Config{
    ServerAddr: ":3000",
    LogLevel:   slog.LevelDebug,
})

if err := gsi.Listen(); err != nil {
    log.Fatalf("Failed to start server: %v", err)
}
```

### Custom Configuration

```go
config := cs2gsi.Config{
    ServerAddr:            ":8080",
    RegulationMaxRounds: 30,  // Custom round limit
    OvertimeMaxRounds:   6,   // Custom overtime limit
    LogLevel:              slog.LevelWarn,
}

gsi := cs2gsi.New(config)
```

## 🧪 Development

This project uses [Task](https://taskfile.dev):

```bash
task check          # fmt + vet + test
task release-check  # pre-release validation
```

## 📦 Releasing

Follow the [Go module publishing guide](https://go.dev/doc/modules/publishing):

1. Commit all changes on `master`
2. Run `task release-check`
3. Run `task release VERSION=v0.3.0`
4. Never retag an existing version — publish a new one instead

Consumers install with:

```bash
go get github.com/nescabir/go-cs2-gsi@v0.3.0
```

## 📈 Performance

- Built with Go 1.24+ for optimal performance
- Efficient memory management with pre-allocated slices
- Type-safe event system with minimal overhead
- Concurrent event handling with proper synchronization

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [osztenkurden's NodeJS implementation](https://github.com/osztenkurden/csgogsi) for the data processing which I heavly relied on
- Valve Corporation for the CS2 Game State Integration API
- The Go community for excellent tooling and libraries

---

**Note**: This library requires Counter-Strike 2 to be running with Game State Integration enabled. Make sure to properly configure the GSI file in your CS2 installation.
