package cs2gsi

import "log/slog"

type Config struct {
	REGULATION_MAX_ROUNDS int
	OVERTIME_MAX_ROUNDS   int
	ServerAddr            string
	LogLevel              slog.Level
}

// NewConfig creates a new Config with sensible defaults
func NewConfig() Config {
	return Config{
		REGULATION_MAX_ROUNDS: 13,
		OVERTIME_MAX_ROUNDS:   3,
		ServerAddr:            ":3000",
		LogLevel:              slog.LevelInfo,
	}
}

// SetDefaults sets default values for any unset fields
func (c *Config) SetDefaults() {
	if c.REGULATION_MAX_ROUNDS <= 0 {
		c.REGULATION_MAX_ROUNDS = 13
	}
	if c.OVERTIME_MAX_ROUNDS <= 0 {
		c.OVERTIME_MAX_ROUNDS = 3
	}
	if c.ServerAddr == "" {
		c.ServerAddr = ":3000"
	}
	if c.LogLevel == 0 {
		c.LogLevel = slog.LevelInfo
	}
}
