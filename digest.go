package cs2gsi

import (
	"encoding/json"
	"errors"
	"fmt"

	models "github.com/nescabir/go-cs2-gsi/models"
	rawModels "github.com/nescabir/go-cs2-gsi/raw"
)

// Digest parses a raw GSI JSON payload and updates internal state.
func (gsi *CS2GSI) Digest(raw []byte) error {
	publishRaw(raw)

	stateRaw := &rawModels.State{}
	if err := json.Unmarshal(raw, stateRaw); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	if err := gsi.validateAuthToken(stateRaw); err != nil {
		return err
	}

	return gsi.digest(stateRaw)
}

// Snapshot returns the last successfully parsed game state.
func (gsi *CS2GSI) Snapshot() *models.State {
	return gsi.current
}

func (gsi *CS2GSI) validateAuthToken(rawState *rawModels.State) error {
	if gsi.config.ExpectedToken == "" {
		return nil
	}
	token := ""
	if rawState.Auth != nil {
		token = rawState.Auth.Token
	}
	if token != gsi.config.ExpectedToken {
		return ErrInvalidToken
	}
	return nil
}

// IsInvalidToken reports whether err is an auth token validation failure.
func IsInvalidToken(err error) bool {
	return errors.Is(err, ErrInvalidToken)
}
