package cs2gsi

import (
	"fmt"
	"strconv"
	"strings"

	rawModels "github.com/nescabir/go-cs2-gsi/raw"
)

func (gsi *CS2GSI) isValidGameState(rawState *rawModels.State) bool {
	return rawState != nil &&
		rawState.Map != nil &&
		rawState.Player != nil &&
		rawState.AllPlayers != nil &&
		rawState.Phase_countdowns != nil &&
		rawState.Round != nil &&
		rawState.Bomb != nil &&
		rawState.Auth != nil &&
		rawState.Grenades != nil
}

// validateRawData validates the raw game state data
func (gsi *CS2GSI) validateRawData(rawState *rawModels.State) error {
	// Check for nil state
	if rawState == nil {
		gsi.logger.Error("raw state is nil")
		return fmt.Errorf("raw state is nil")
	}

	// Validate players
	if err := gsi.validatePlayers(rawState.AllPlayers); err != nil {
		return err
	}

	// Validate map data
	if rawState.Map == nil {
		gsi.logger.Debug("map data is nil (likely main menu)")
		return nil // Return early for main menu state
	}

	if err := gsi.validateMap(rawState.Map); err != nil {
		return err
	}

	// Validate phase countdowns
	if err := gsi.validatePhaseCountdowns(rawState.Phase_countdowns); err != nil {
		return err
	}

	// Validate provider (optional but should be consistent)
	if rawState.Provider != nil {
		if err := gsi.validateProvider(rawState.Provider); err != nil {
			gsi.logger.Warn("provider validation failed", "error", err)
			// Don't return error for provider validation as it's not critical
		}
	}

	// Validate round data (optional but should be consistent)
	if rawState.Round != nil {
		if err := gsi.validateRound(rawState.Round); err != nil {
			gsi.logger.Warn("round validation failed", "error", err)
			// Don't return error for round validation as it's not critical
		}
	}

	// Validate player data (observed player)
	if rawState.Player == nil {
		gsi.logger.Debug("observed player is nil (likely main menu)")
		return nil // Return early for main menu state
	}

	if err := gsi.validateObservedPlayer(rawState.Player); err != nil {
		return err
	}

	return nil
}

// validatePlayers validates all players in the game state
func (gsi *CS2GSI) validatePlayers(allPlayers map[string]*rawModels.Player) error {
	if len(allPlayers) == 0 {
		gsi.logger.Debug("no players found in game state (likely main menu)")
		return nil // Return early for main menu state
	}

	// Check for reasonable player count (CS2 max is 10 players)
	if len(allPlayers) > 10 {
		gsi.logger.Warn("unexpected number of players", "count", len(allPlayers))
	}

	// Validate each player
	for steamId, player := range allPlayers {
		if err := gsi.validatePlayer(player, steamId); err != nil {
			gsi.logger.Error("player validation failed", "steamId", steamId, "error", err)
			return fmt.Errorf("player %s validation failed: %w", steamId, err)
		}
	}

	return nil
}

// validatePlayer validates a single player
func (gsi *CS2GSI) validatePlayer(player *rawModels.Player, steamId string) error {
	if player == nil {
		return fmt.Errorf("player is nil")
	}

	// Validate Steam ID format (should be numeric and reasonable length)
	if steamId == "" {
		return fmt.Errorf("steam ID is empty")
	}

	// Validate player name
	if player.Name == "" {
		gsi.logger.Warn("player has empty name", "steamId", steamId)
	}

	if len(player.Name) > 100 {
		gsi.logger.Warn("player name is unusually long", "steamId", steamId, "name", player.Name)
	}

	// Validate player state
	if player.State != nil {
		if err := gsi.validatePlayerState(player.State); err != nil {
			return fmt.Errorf("player state validation failed: %w", err)
		}
	}

	// Validate team assignment
	if player.Team != rawModels.CTSide && player.Team != rawModels.TSide && player.Team != rawModels.NilSide {
		gsi.logger.Warn("invalid team assignment", "steamId", steamId, "team", player.Team)
	}

	// Validate observer slot
	if player.Observer_slot < -1 || player.Observer_slot > 10 {
		gsi.logger.Warn("invalid observer slot", "steamId", steamId, "slot", player.Observer_slot)
	}

	return nil
}

// validatePlayerState validates player state data
func (gsi *CS2GSI) validatePlayerState(state *rawModels.PlayerState) error {
	if state == nil {
		return fmt.Errorf("player state is nil")
	}

	// Validate health (0-100)
	if state.Health < 0 || state.Health > 100 {
		gsi.logger.Warn("invalid health value", "health", state.Health)
	}

	// Validate armor (0-100)
	if state.Armor < 0 || state.Armor > 100 {
		gsi.logger.Warn("invalid armor value", "armor", state.Armor)
	}

	// Validate money (reasonable range for CS2)
	if state.Money < 0 || state.Money > 16000 {
		gsi.logger.Warn("unusual money value", "money", state.Money)
	}

	// Validate round stats
	if state.Round_kills < 0 || state.Round_kills > 5 {
		gsi.logger.Warn("unusual round kills", "kills", state.Round_kills)
	}

	if state.Round_killhs < 0 || state.Round_killhs > state.Round_kills {
		gsi.logger.Warn("invalid round headshots", "headshots", state.Round_killhs, "kills", state.Round_kills)
	}

	if state.Round_totaldmg < 0 || state.Round_totaldmg > 1000 {
		gsi.logger.Warn("unusual round damage", "damage", state.Round_totaldmg)
	}

	// Validate equipment value
	if state.Equip_value < 0 || state.Equip_value > 16000 {
		gsi.logger.Warn("unusual equipment value", "equip_value", state.Equip_value)
	}

	return nil
}

// validateMap validates map data
func (gsi *CS2GSI) validateMap(mapData *rawModels.Map) error {
	if mapData == nil {
		gsi.logger.Error("map data is missing")
		return fmt.Errorf("map data is missing")
	}

	// Validate map name
	if mapData.Name == "" {
		gsi.logger.Error("map name is empty")
		return fmt.Errorf("map name is empty")
	}

	// Check for valid map names (CS2 competitive maps)
	validMaps := map[string]bool{
		"de_ancient":  true,
		"de_anubis":   true,
		"de_inferno":  true,
		"de_mirage":   true,
		"de_nuke":     true,
		"de_overpass": true,
		"de_vertigo":  true,
		"de_cache":    true,
		"de_dust2":    true,
		"de_train":    true,
	}

	// Extract map name from path if needed
	mapName := mapData.Name
	if strings.Contains(mapName, "/") {
		mapName = mapName[strings.LastIndex(mapName, "/")+1:]
	}

	if !validMaps[mapName] {
		gsi.logger.Warn("unknown map", "map", mapName)
	}

	// Validate round number
	if mapData.Round < 0 || mapData.Round > 50 {
		gsi.logger.Warn("unusual round number", "round", mapData.Round)
	}

	// Validate team data
	if err := gsi.validateTeam(mapData.Team_ct, "CT"); err != nil {
		return fmt.Errorf("ct team validation failed: %w", err)
	}

	if err := gsi.validateTeam(mapData.Team_t, "T"); err != nil {
		return fmt.Errorf("t team validation failed: %w", err)
	}

	// Validate spectators count
	if mapData.Current_spectators < 0 || mapData.Current_spectators > 100 {
		gsi.logger.Warn("unusual spectator count", "spectators", mapData.Current_spectators)
	}

	// Validate series matches
	if mapData.Num_matches_to_win_series < 0 || mapData.Num_matches_to_win_series > 5 {
		gsi.logger.Warn("unusual series matches", "matches", mapData.Num_matches_to_win_series)
	}

	return nil
}

// validateTeam validates team data
func (gsi *CS2GSI) validateTeam(team *rawModels.Team, side string) error {
	if team == nil {
		return fmt.Errorf("%s team is nil", side)
	}

	// Validate team score
	if team.Score < 0 || team.Score > 30 {
		gsi.logger.Warn("unusual team score", "side", side, "score", team.Score)
	}

	// Validate consecutive losses
	if team.Consecutive_round_losses < 0 || team.Consecutive_round_losses > 15 {
		gsi.logger.Warn("unusual consecutive losses", "side", side, "losses", team.Consecutive_round_losses)
	}

	// Validate timeouts remaining
	if team.Timeouts_remaining < 0 || team.Timeouts_remaining > 4 {
		gsi.logger.Warn("unusual timeouts remaining", "side", side, "timeouts", team.Timeouts_remaining)
	}

	// Validate matches won
	if team.Matches_won_this_series < 0 || team.Matches_won_this_series > 3 {
		gsi.logger.Warn("unusual matches won", "side", side, "matches", team.Matches_won_this_series)
	}

	return nil
}

// validatePhaseCountdowns validates phase countdown data
func (gsi *CS2GSI) validatePhaseCountdowns(phaseCountdowns *rawModels.PhaseCountdown) error {
	if phaseCountdowns == nil {
		gsi.logger.Error("phase countdown data is missing")
		return fmt.Errorf("phase countdown data is missing")
	}

	// Validate phase type
	validPhases := map[rawModels.PhaseType]bool{
		rawModels.PhaseTypeFreezetime: true,
		rawModels.PhaseTypeBomb:       true,
		rawModels.PhaseTypeWarmup:     true,
		rawModels.PhaseTypeLive:       true,
		rawModels.PhaseTypeOver:       true,
		rawModels.PhaseTypeDefuse:     true,
		rawModels.PhaseTypePaused:     true,
		rawModels.PhaseTypeTimeoutCT:  true,
		rawModels.PhaseTypeTimeoutT:   true,
		rawModels.PhaseTypeNil:        true,
	}

	if !validPhases[phaseCountdowns.Phase] {
		gsi.logger.Warn("unknown phase type", "phase", phaseCountdowns.Phase)
	}

	// Validate phase ends in (should be reasonable time)
	if phaseCountdowns.Phase_ends_in != "" {
		timeLeft, err := strconv.ParseFloat(phaseCountdowns.Phase_ends_in, 64)
		if err != nil {
			gsi.logger.Warn("invalid phase ends in value", "value", phaseCountdowns.Phase_ends_in)
		} else if timeLeft < -10 || timeLeft > 300 {
			gsi.logger.Warn("unusual phase time remaining", "time", timeLeft)
		}
	}

	return nil
}

// validateProvider validates provider data
func (gsi *CS2GSI) validateProvider(provider *rawModels.Provider) error {
	if provider == nil {
		return fmt.Errorf("provider is nil")
	}

	// Validate app ID (CS2 app ID is 730)
	if provider.AppId != 730 {
		gsi.logger.Warn("unexpected app ID", "appId", provider.AppId)
	}

	// Validate timestamp (should be recent)
	if provider.Timestamp > 0 {
		// Check if timestamp is within reasonable range (not too old, not in future)
		// This would require time.Now() comparison, but we'll keep it simple for now
		if provider.Timestamp < 1000000000 { // Before 2001
			gsi.logger.Warn("unusual timestamp", "timestamp", provider.Timestamp)
		}
	}

	return nil
}

// validateRound validates round data
func (gsi *CS2GSI) validateRound(round *rawModels.Round) error {
	if round == nil {
		return fmt.Errorf("round is nil")
	}

	// Validate round phase
	validRoundPhases := map[rawModels.RoundPhase]bool{
		rawModels.RoundPhaseFreezeTime: true,
		rawModels.RoundPhaseLive:       true,
		rawModels.RoundPhaseOver:       true,
		rawModels.RoundPhaseNil:        true,
	}

	if !validRoundPhases[round.Phase] {
		gsi.logger.Warn("unknown round phase", "phase", round.Phase)
	}

	// Validate win team
	if round.Win_team != "" {
		if round.Win_team != rawModels.CTSide && round.Win_team != rawModels.TSide {
			gsi.logger.Warn("invalid win team", "team", round.Win_team)
		}
	}

	// Validate bomb state
	if round.Bomb != "" {
		validBombStates := map[rawModels.BombRoundState]bool{
			rawModels.BombRoundStatePlanted:  true,
			rawModels.BombRoundStateExploded: true,
			rawModels.BombRoundStateDefused:  true,
			rawModels.BombRoundStateNil:      true,
		}

		if !validBombStates[round.Bomb] {
			gsi.logger.Warn("unknown bomb state", "bomb", round.Bomb)
		}
	}

	return nil
}

// validateObservedPlayer validates observed player data
func (gsi *CS2GSI) validateObservedPlayer(player *rawModels.PlayerObserved) error {
	if player == nil {
		gsi.logger.Error("observed player data is missing")
		return fmt.Errorf("observed player data is missing")
	}

	// Validate the underlying player data
	if err := gsi.validatePlayer(&player.Player, "observed"); err != nil {
		return fmt.Errorf("observed player validation failed: %w", err)
	}

	// Validate spectarget (optional field)
	if player.Spectarget != "" {
		if len(player.Spectarget) < 10 || len(player.Spectarget) > 20 {
			gsi.logger.Warn("unusual spectarget Steam ID length", "spectarget", player.Spectarget)
		}
	}

	return nil
}
