package cs2gsi

import "log/slog"

type Config struct {
	RegulationMaxRounds int
	OvertimeMaxRounds   int
	ServerAddr          string
	LogLevel            slog.Level
}

// NewConfig creates a new Config with sensible defaults
func NewConfig() Config {
	return Config{
		RegulationMaxRounds: 13,
		OvertimeMaxRounds:   3,
		ServerAddr:          ":3000",
		LogLevel:            slog.LevelInfo,
	}
}

// SetDefaults sets default values for any unset fields
func (c *Config) SetDefaults() {
	if c.RegulationMaxRounds <= 0 {
		c.RegulationMaxRounds = 13
	}
	if c.OvertimeMaxRounds <= 0 {
		c.OvertimeMaxRounds = 3
	}
	if c.ServerAddr == "" {
		c.ServerAddr = ":3000"
	}
	if c.LogLevel == 0 {
		c.LogLevel = slog.LevelInfo
	}
}
