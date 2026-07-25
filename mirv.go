package cs2gsi

import (
	"encoding/json"
	"errors"
	"fmt"

	models "github.com/nescabir/go-cs2-gsi/models"
	rawModels "github.com/nescabir/go-cs2-gsi/raw"
)

const (
	MIRVEventPlayerDeath = "player_death"
	MIRVEventPlayerHurt  = "player_hurt"
)

var (
	ErrMIRVNoPriorState = errors.New("digest must be called before DigestMIRV")
	ErrMIRVUnknownEvent = errors.New("unknown MIRV event type")
)

type MIRVResult struct {
	Kill *models.KillEvent
	Hurt *models.HurtEvent
}

func (gsi *CS2GSI) DigestMIRV(raw []byte, eventType string) (*MIRVResult, error) {
	if gsi.last == nil {
		return nil, ErrMIRVNoPriorState
	}

	switch eventType {
	case MIRVEventPlayerDeath:
		var rawKill rawModels.RawKill
		if err := json.Unmarshal(raw, &rawKill); err != nil {
			return nil, fmt.Errorf("failed to decode MIRV kill event: %w", err)
		}
		kill, err := gsi.digestMIRVKill(&rawKill)
		if err != nil {
			return nil, err
		}
		if kill != nil {
			publishKill(kill)
		}
		return &MIRVResult{Kill: kill}, nil
	case MIRVEventPlayerHurt:
		var rawHurt rawModels.RawHurt
		if err := json.Unmarshal(raw, &rawHurt); err != nil {
			return nil, fmt.Errorf("failed to decode MIRV hurt event: %w", err)
		}
		hurt, err := gsi.digestMIRVHurt(&rawHurt)
		if err != nil {
			return nil, err
		}
		if hurt != nil {
			publishHurt(hurt)
		}
		return &MIRVResult{Hurt: hurt}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrMIRVUnknownEvent, eventType)
	}
}

func (gsi *CS2GSI) digestMIRVKill(raw *rawModels.RawKill) (*models.KillEvent, error) {
	data := raw.Keys
	killer := gsi.findPlayerBySteamID(data.Attacker.Xuid)
	victim := gsi.findPlayerBySteamID(data.Userid.Xuid)
	assister := gsi.findPlayerBySteamID(data.Assister.Xuid)

	if victim == nil {
		return nil, nil
	}

	if data.Assister.Xuid == "0" {
		assister = nil
	}

	if killer == nil && (data.Weapon == "trigger_hurt" || data.Weapon == "worldspawn") {
		killer = victim
	}

	return &models.KillEvent{
		Attacker:      killer,
		Victim:        victim,
		Weapon:        &models.Weapon{Name: data.Weapon},
		Assister:      assister,
		Flashed:       data.Assistedflash,
		Headshot:      data.Headshot,
		Wallbang:      data.Penetrated > 0,
		AttackerBlind: data.Attackerblind,
		ThruSmoke:     data.Thrusmoke,
		NoScope:       data.Noscope,
		AttackerInAir: data.Attackerinair,
	}, nil
}

func (gsi *CS2GSI) digestMIRVHurt(raw *rawModels.RawHurt) (*models.HurtEvent, error) {
	data := raw.Keys
	attacker := gsi.findPlayerBySteamID(data.Attacker.Xuid)
	victim := gsi.findPlayerBySteamID(data.Userid.Xuid)

	if attacker == nil || victim == nil {
		return nil, nil
	}

	return &models.HurtEvent{
		Attacker:  attacker,
		Victim:    victim,
		Weapon:    &models.Weapon{Name: data.Weapon},
		Health:    data.Health,
		Armor:     data.Armor,
		DmgHealth: data.DmgHealth,
		DmgArmor:  data.DmgArmor,
		HitGroup:  data.Hitgroup,
	}, nil
}

func (gsi *CS2GSI) findPlayerBySteamID(steamID string) *models.Player {
	if gsi.last == nil || steamID == "" {
		return nil
	}
	return gsi.last.AllPlayers[steamID]
}
