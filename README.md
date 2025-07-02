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
    ServerAddr:            ":3000",           // HTTP server address ex: 127.0.0.1:3000, localhost:4242
    REGULATION_MAX_ROUNDS: 13,               // Max rounds in regulation
    OVERTIME_MAX_ROUNDS:   3,                // Max rounds in overtime
    LogLevel:              slog.LevelInfo,   // Logging level
}

gsi := cs2gsi.New(config)
```

## 🎮 Available Events

The library provides type-safe events for all major game occurrences:

### Core Events

- `Data` - Raw game state updates
- `RoundEnd` - Round completion with winner information
- `MapEnd` - Map completion
- `MatchEnd` - Match completion

### Combat Events

- `Mvp` - MVP player selection

- MIRV needed (Not implemented)
  - `Kill` - Player kills with detailed weapon and damage info
  - `Hurt` - Player damage events

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
    REGULATION_MAX_ROUNDS: 30,  // Custom round limit
    OVERTIME_MAX_ROUNDS:   6,   // Custom overtime limit
    LogLevel:              slog.LevelWarn,
}

gsi := cs2gsi.New(config)
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
