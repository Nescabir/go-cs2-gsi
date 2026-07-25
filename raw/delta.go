package raw

import "encoding/json"

// DeltaState is a shallow GSI delta block (previously / added) without nested recursion.
type DeltaState struct {
	Player           *PlayerObserved     `json:"player"`
	AllPlayers       map[string]*Player  `json:"allplayers"`
	Bomb             *Bomb               `json:"bomb"`
	Round            *Round              `json:"round"`
	Grenades         map[string]*Grenade `json:"grenades"`
	Phase_countdowns *PhaseCountdown     `json:"phase_countdowns"`
}

func unmarshalDeltaField(data interface{}, target interface{}) error {
	if data == nil {
		return nil
	}
	if boolVal, ok := data.(bool); ok && !boolVal {
		return nil
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}

func unmarshalDeltaGrenades(values map[string]interface{}) map[string]*Grenade {
	if values == nil {
		return nil
	}
	grenades := make(map[string]*Grenade)
	for key, value := range values {
		if boolVal, ok := value.(bool); ok && !boolVal {
			continue
		}
		raw, err := json.Marshal(value)
		if err != nil {
			continue
		}
		var grenade Grenade
		if err := json.Unmarshal(raw, &grenade); err == nil {
			grenades[key] = &grenade
		}
	}
	return grenades
}

func unmarshalDeltaAllPlayers(values map[string]interface{}) map[string]*Player {
	if values == nil {
		return nil
	}
	players := make(map[string]*Player)
	for key, value := range values {
		if boolVal, ok := value.(bool); ok && !boolVal {
			continue
		}
		raw, err := json.Marshal(value)
		if err != nil {
			continue
		}
		var player Player
		if err := json.Unmarshal(raw, &player); err == nil {
			players[key] = &player
		}
	}
	return players
}

// UnmarshalJSON parses a shallow delta object from GSI previously/added blocks.
func (d *DeltaState) UnmarshalJSON(data []byte) error {
	type deltaAlias DeltaState
	aux := &struct {
		*deltaAlias
		Grenades         map[string]interface{} `json:"grenades"`
		AllPlayers       map[string]interface{} `json:"allplayers"`
		Player           interface{}            `json:"player"`
		Bomb             interface{}            `json:"bomb"`
		Round            interface{}            `json:"round"`
		Phase_countdowns interface{}            `json:"phase_countdowns"`
	}{
		deltaAlias: (*deltaAlias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Grenades != nil {
		d.Grenades = unmarshalDeltaGrenades(aux.Grenades)
	}
	if aux.AllPlayers != nil {
		d.AllPlayers = unmarshalDeltaAllPlayers(aux.AllPlayers)
	}
	if aux.Player != nil {
		var player PlayerObserved
		if err := unmarshalDeltaField(aux.Player, &player); err == nil {
			d.Player = &player
		}
	}
	if aux.Bomb != nil {
		var bomb Bomb
		if err := unmarshalDeltaField(aux.Bomb, &bomb); err == nil {
			d.Bomb = &bomb
		}
	}
	if aux.Round != nil {
		var round Round
		if err := unmarshalDeltaField(aux.Round, &round); err == nil {
			d.Round = &round
		}
	}
	if aux.Phase_countdowns != nil {
		var phase PhaseCountdown
		if err := unmarshalDeltaField(aux.Phase_countdowns, &phase); err == nil {
			d.Phase_countdowns = &phase
		}
	}

	return nil
}
