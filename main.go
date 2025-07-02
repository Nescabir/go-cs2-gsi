package cs2gsi

import (
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
	models "github.com/nescabir/go-cs2-gsi/models"
)

type teams struct {
	ct *models.Team
	t  *models.Team
}

type CS2GSI struct {
	config                Config
	REGULATION_MAX_ROUNDS int
	OVERTIME_MAX_ROUNDS   int
	logger                *slog.Logger
	damage                []models.RoundDamage
	players               []models.Player
	teams                 *teams
	current               *models.State
	last                  *models.State
}

func New(config Config) *CS2GSI {
	// Set defaults for any unset fields
	config.SetDefaults()
	logHandler := slogpretty.New(os.Stdout, &slogpretty.Options{
		Level:      config.LogLevel,
		AddSource:  true,
		Colorful:   true,
		Multiline:  true,
		TimeFormat: slogpretty.DefaultTimeFormat,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	return &CS2GSI{
		config:                config,
		logger:                logger,
		REGULATION_MAX_ROUNDS: config.REGULATION_MAX_ROUNDS,
		OVERTIME_MAX_ROUNDS:   config.OVERTIME_MAX_ROUNDS,
		damage:                make([]models.RoundDamage, 0, 60),
		players:               make([]models.Player, 0, 16),
		teams: &teams{
			ct: nil,
			t:  nil,
		},
		current: nil,
		last:    nil,
	}
}
