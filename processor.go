package cs2gsi

import (
	"fmt"
	"math"
	"strings"

	models "github.com/nescabir/go-cs2-gsi/models"
	rawModels "github.com/nescabir/go-cs2-gsi/raw"
)

// digest processes the raw game state data
func (gsi *CS2GSI) digest(rawState *rawModels.State) error {
	if !gsi.isValidGameState(rawState) {
		gsi.logger.Debug("invalid game state (likely main menu)")
		return nil
	}

	gsi.logger.Debug("processing game state",
		"map", rawState.Map.Name,
		"round", rawState.Map.Round)

	// Validate input data
	if err := gsi.validateRawData(rawState); err != nil {
		return err
	}

	// Initialize and parse basic state
	state := gsi.initState()
	if err := gsi.parseBasicState(rawState, state); err != nil {
		return err
	}

	// Parse players and teams
	if err := gsi.parsePlayersAndTeams(rawState, state); err != nil {
		return err
	}

	// Process rounds and damage
	if err := gsi.processRoundsAndDamage(rawState, state); err != nil {
		return err
	}

	// Calculate player statistics
	if err := gsi.calculatePlayerStats(rawState, state); err != nil {
		return err
	}

	// Update state and detect events
	if err := gsi.updateStateAndDetectEvents(state); err != nil {
		return err
	}

	gsi.logger.Debug("game state processed successfully",
		"players", len(rawState.AllPlayers),
		"round", rawState.Map.Round)

	return nil
}

// parseBasicState parses the basic state information
func (gsi *CS2GSI) parseBasicState(rawState *rawModels.State, state *models.State) error {
	ctOrientation := getCTOrientation(rawState)

	state.Provider = gsi.parseProvider(rawState.Provider)
	state.Map = gsi.parseMap(rawState.Map, &ctOrientation)
	state.Round = gsi.parseRound(rawState.Round)
	state.Phase_countdowns = gsi.parsePhaseCountdown(rawState.Phase_countdowns)
	state.Auth = gsi.parseAuth(rawState.Auth)
	state.Grenades = gsi.parseGrenades(rawState.Grenades)

	// Set up teams
	gsi.teams = &teams{
		ct: state.Map.Team_ct,
		t:  state.Map.Team_t,
	}

	// Parse bomb and observer
	state.Bomb = gsi.parseBomb(rawState.Bomb, rawState.AllPlayers, gsi.teams, &rawState.Map.Name)
	state.Observer = &models.Observer{
		Activity:   models.PlayerActivity(rawState.Player.Activity),
		Spectarget: rawState.Player.Spectarget,
		Position:   parseVector(rawState.Player.Position),
		Forward:    parseVector(rawState.Player.Forward),
	}

	return nil
}

// parsePlayersAndTeams parses player and team data
func (gsi *CS2GSI) parsePlayersAndTeams(rawState *rawModels.State, state *models.State) error {
	if len(rawState.AllPlayers) == 0 {
		return nil
	}

	// Clear players slice for new state
	gsi.players = gsi.players[:0]

	for steamId, rawPlayer := range rawState.AllPlayers {
		parsedPlayer := gsi.parsePlayer(rawPlayer, gsi.teams)
		if parsedPlayer == nil {
			gsi.logger.Warn("failed to parse player", "steam_id", steamId)
			continue // Skip invalid players
		}

		parsedPlayer.SteamId = steamId

		// Set as current player if it's the observed player
		if rawState.Player != nil && steamId == rawState.Player.Steamid {
			state.Player = parsedPlayer
		}

		state.AllPlayers[steamId] = parsedPlayer
		gsi.players = append(gsi.players, *parsedPlayer)
	}

	return nil
}

// processRoundsAndDamage processes rounds and damage data
func (gsi *CS2GSI) processRoundsAndDamage(rawState *rawModels.State, state *models.State) error {
	// Process rounds
	if err := gsi.processRounds(rawState, state); err != nil {
		return err
	}

	// Reset damage when map changes
	if gsi.last != nil && gsi.last.Map.Name != rawState.Map.Name {
		gsi.damage = make([]models.RoundDamage, 0, 60)
	}

	// Process damage
	if err := gsi.processDamage(rawState, state); err != nil {
		return err
	}

	return nil
}

// processRounds processes round information
func (gsi *CS2GSI) processRounds(rawState *rawModels.State, state *models.State) error {
	if rawState.Round == nil || rawState.Map == nil || rawState.Map.Round_wins == nil {
		return nil
	}

	var currentRound int = rawState.Map.Round + 1
	if rawState.Round != nil && (rawState.Round.Phase == rawModels.RoundPhaseOver || rawState.Map.Phase == rawModels.MapPhaseGameOver) {
		currentRound = rawState.Map.Round
	}

	var rounds []models.RoundInfo
	if currentRound > 0 {
		for i := 1; i <= currentRound; i++ {
			result := getRoundWin(currentRound, gsi.teams, rawState.Map.Round_wins, i, gsi.regulationMaxRounds, gsi.overtimeMaxRounds)
			if result == nil {
				continue
			}
			rounds = append(rounds, *result)
		}
	}
	state.Map.Rounds = rounds

	return nil
}

// processDamage processes damage data for the current round
func (gsi *CS2GSI) processDamage(rawState *rawModels.State, state *models.State) error {
	currentRoundForDamage := gsi.getCurrentRoundForDamage(rawState)

	// Reset damage on warmup/freezetime
	if rawState.Map.Round == 0 && (rawState.Phase_countdowns.Phase == rawModels.PhaseTypeFreezetime || rawState.Phase_countdowns.Phase == rawModels.PhaseTypeWarmup) {
		gsi.damage = make([]models.RoundDamage, 0, 60)
		return nil
	}

	// Find or create current round damage
	currentRoundDamage := gsi.findOrCreateRoundDamage(currentRoundForDamage)
	if currentRoundDamage == nil {
		return nil
	}

	// Add player damage for current round
	if len(gsi.players) > 0 {
		for _, player := range gsi.players {
			currentRoundDamage.Players = append(currentRoundDamage.Players, models.RoundPlayerDamage{
				SteamId: player.SteamId,
				Damage:  player.State.Round_totaldmg,
			})
		}
	}

	return nil
}

// getCurrentRoundForDamage determines the current round for damage calculation
func (gsi *CS2GSI) getCurrentRoundForDamage(rawState *rawModels.State) int {
	currentRoundForDamage := rawState.Map.Round + 1
	if rawState.Round != nil && (rawState.Round.Phase == rawModels.RoundPhaseOver || rawState.Map.Phase == rawModels.MapPhaseGameOver) {
		currentRoundForDamage = rawState.Map.Round
	}
	return currentRoundForDamage
}

// findOrCreateRoundDamage finds or creates a round damage entry
func (gsi *CS2GSI) findOrCreateRoundDamage(currentRound int) *models.RoundDamage {
	// Find existing damage for this round
	for _, damage := range gsi.damage {
		if damage.Round == currentRound {
			return &damage
		}
	}

	// Create new damage entry
	newDamage := &models.RoundDamage{
		Round:   currentRound,
		Players: make([]models.RoundPlayerDamage, 0, 16),
	}
	gsi.damage = append(gsi.damage, *newDamage)
	return newDamage
}

// calculatePlayerStats calculates player statistics including ADR
func (gsi *CS2GSI) calculatePlayerStats(rawState *rawModels.State, state *models.State) error {
	if len(gsi.players) == 0 || gsi.current == nil {
		return nil
	}

	currentRoundForDamage := gsi.getCurrentRoundForDamage(rawState)

	for _, player := range gsi.players {
		adr := gsi.calculatePlayerADR(player, currentRoundForDamage)
		player.State.Adr = int(math.Floor(adr))
	}

	return nil
}

// calculatePlayerADR calculates the Average Damage per Round for a player
func (gsi *CS2GSI) calculatePlayerADR(player models.Player, currentRound int) float64 {
	// Get damage for previous rounds
	damageForRound := make([]models.RoundDamage, 0, 60)
	for _, damage := range gsi.damage {
		if damage.Round < currentRound {
			damageForRound = append(damageForRound, damage)
		}
	}

	if len(damageForRound) == 0 {
		return 0
	}

	// Calculate total damage for player
	var totalDamage int
	for _, damage := range damageForRound {
		for _, playerDamage := range damage.Players {
			if playerDamage.SteamId == player.SteamId {
				totalDamage += playerDamage.Damage
				break
			}
		}
	}

	// Calculate ADR
	if currentRound == 0 {
		return float64(totalDamage)
	}
	return float64(totalDamage) / float64(currentRound)
}

// updateStateAndDetectEvents updates the current state and detects events
func (gsi *CS2GSI) updateStateAndDetectEvents(state *models.State) error {
	// Update current state
	gsi.current = state

	// Handle first state
	if gsi.last == nil {
		gsi.last = state
		publishData(state)
		return nil
	}

	// Detect and publish events
	if err := gsi.detectRoundEvents(state); err != nil {
		return err
	}

	if err := gsi.detectBombEvents(state); err != nil {
		return err
	}

	if err := gsi.detectPhaseEvents(state); err != nil {
		return err
	}

	if err := gsi.detectTimeoutEvents(state); err != nil {
		return err
	}

	if err := gsi.detectMVPEvents(state); err != nil {
		return err
	}

	// Publish data and update last state
	publishData(state)
	gsi.last = state

	return nil
}

// detectRoundEvents detects round-related events
func (gsi *CS2GSI) detectRoundEvents(state *models.State) error {
	last := gsi.last
	if last == nil || last.Round == nil || state.Round == nil {
		return nil
	}

	// Check for round end
	if state.Round.Win_team != "" && last.Round.Win_team == "" {
		winner, loser := gsi.determineWinnerAndLoser(state)
		gsi.updateWinnerScore(winner, last)

		roundScore := &models.Score{
			Winner: winner,
			Loser:  loser,
			Map:    state.Map,
			MapEnd: state.Map.Phase == models.MapPhaseGameOver,
		}

		gsi.logger.Info("Round end detected", "winner", winner.Side, "loser", loser.Side, "score", fmt.Sprintf("%d-%d", winner.Score, loser.Score))
		publishRoundEnd(roundScore)

		// Check for match end
		if roundScore.MapEnd && last.Map.Phase != models.MapPhaseGameOver {
			gsi.logger.Info("Match end detected", "winner", winner.Side, "loser", loser.Side, "score", fmt.Sprintf("%d-%d", winner.Score, loser.Score))
			publishMatchEnd(roundScore)
		}
	}

	return nil
}

// determineWinnerAndLoser determines which team won and lost
func (gsi *CS2GSI) determineWinnerAndLoser(state *models.State) (*models.Team, *models.Team) {
	if state.Round.Win_team == models.CTSide {
		return state.Map.Team_ct, state.Map.Team_t
	}
	return state.Map.Team_t, state.Map.Team_ct
}

// updateWinnerScore updates the winner's score
func (gsi *CS2GSI) updateWinnerScore(winner *models.Team, last *models.State) {
	var oldWinner *models.Team
	if winner.Side == models.CTSide {
		oldWinner = last.Map.Team_ct
	} else {
		oldWinner = last.Map.Team_t
	}

	if winner.Score == oldWinner.Score {
		winner.Score += 1
	}
}

// detectBombEvents detects bomb-related events
func (gsi *CS2GSI) detectBombEvents(state *models.State) error {
	last := gsi.last
	if last == nil {
		return nil
	}

	// Handle bomb state changes
	if last.Bomb != nil && state.Bomb != nil {
		gsi.handleBombStateChanges(last.Bomb, state.Bomb)
	} else if last.Bomb == nil && state.Bomb != nil && state.Bomb.State == models.BombStateExploded {
		gsi.logger.Info("Bomb exploded")
		publishBombExploded(nil)
	}

	return nil
}

// handleBombStateChanges handles changes in bomb state
func (gsi *CS2GSI) handleBombStateChanges(lastBomb, currentBomb *models.Bomb) {
	// Bomb plant stop
	if lastBomb.State == models.BombStatePlanting &&
		currentBomb.State != models.BombStatePlanting &&
		currentBomb.State != models.BombStatePlanted &&
		currentBomb.State != models.BombStateDefusing {
		gsi.logger.Info("Bomb plant stop detected", "player", lastBomb.Player.Name)
		publishBombPlantStop(lastBomb.Player)
	}

	// Bomb planted
	if lastBomb.State == models.BombStatePlanting && currentBomb.State == models.BombStatePlanted {
		gsi.logger.Info("Bomb planted", "player", lastBomb.Player.Name)
		publishBombPlanted(lastBomb.Player)
	}

	// Bomb exploded
	if lastBomb.State != models.BombStateExploded && currentBomb.State == models.BombStateExploded {
		gsi.logger.Info("Bomb exploded")
		publishBombExploded(nil)
	}

	// Bomb defused
	if lastBomb.State != models.BombStateDefused && currentBomb.State == models.BombStateDefused {
		gsi.logger.Info("Bomb defused", "player", lastBomb.Player.Name)
		publishBombDefused(lastBomb.Player)
	}

	// Defuse start
	if lastBomb.State != models.BombStateDefusing && currentBomb.State == models.BombStateDefusing {
		gsi.logger.Info("Defuse start detected", "player", currentBomb.Player.Name)
		publishDefuseStart(currentBomb.Player)
	}

	// Defuse end
	if lastBomb.State == models.BombStateDefusing && currentBomb.State != models.BombStateDefusing {
		gsi.logger.Info("Defuse end detected", "player", lastBomb.Player.Name)
		publishDefuseEnd(lastBomb.Player)
	}

	// Bomb plant start
	if lastBomb.State != models.BombStatePlanting && currentBomb.State == models.BombStatePlanting {
		gsi.logger.Info("Bomb plant start detected", "player", currentBomb.Player.Name)
		publishBombPlantStart(currentBomb.Player)
	}
}

// detectPhaseEvents detects phase-related events
func (gsi *CS2GSI) detectPhaseEvents(state *models.State) error {
	last := gsi.last
	if last == nil {
		return nil
	}

	// Intermission events
	if state.Map.Phase == models.MapPhaseIntermission && last.Map.Phase != models.MapPhaseIntermission {
		gsi.logger.Info("Intermission start detected")
		publishIntermissionStart(nil)
	} else if state.Map.Phase != models.MapPhaseIntermission && last.Map.Phase == models.MapPhaseIntermission {
		gsi.logger.Info("Intermission end detected")
		publishIntermissionEnd(nil)
	}

	// Freezetime events
	phase := state.Phase_countdowns.Phase
	if phase == models.PhaseTypeFreezetime && last.Phase_countdowns.Phase != models.PhaseTypeFreezetime {
		gsi.logger.Info("Freezetime start detected")
		publishFreezetimeStart(nil)
	} else if phase != models.PhaseTypeFreezetime && last.Phase_countdowns.Phase == models.PhaseTypeFreezetime {
		gsi.logger.Info("Freezetime end detected")
		publishFreezetimeEnd(nil)
	}

	return nil
}

// detectTimeoutEvents detects timeout-related events
func (gsi *CS2GSI) detectTimeoutEvents(state *models.State) error {
	last := gsi.last
	if last == nil {
		return nil
	}

	phase := state.Phase_countdowns.Phase
	lastPhase := last.Phase_countdowns.Phase

	// Timeout start
	if strings.HasPrefix(string(phase), "timeout") && !strings.HasPrefix(string(lastPhase), "timeout") {
		var team *models.Team
		if phase == models.PhaseTypeTimeoutCT {
			team = gsi.teams.ct
		} else {
			team = gsi.teams.t
		}
		gsi.logger.Info("Timeout start detected", "team", team.Name, "side", team.Side)
		publishTimeoutStart(team)
	}

	// Timeout end
	if strings.HasPrefix(string(lastPhase), "timeout") && !strings.HasPrefix(string(phase), "timeout") {
		gsi.logger.Info("Timeout end detected")
		publishTimeoutEnd(nil)
	}

	return nil
}

// detectMVPEvents detects MVP-related events
func (gsi *CS2GSI) detectMVPEvents(state *models.State) error {
	last := gsi.last
	if last == nil {
		return nil
	}

	// Check for MVP
	for _, player := range state.AllPlayers {
		if previousPlayer, exists := last.AllPlayers[player.SteamId]; exists {
			if player.Match_stats.Mvps > previousPlayer.Match_stats.Mvps {
				gsi.logger.Info("MVP detected", "player", player.Name)
				publishMvp(player)
				break
			}
		}
	}

	return nil
}
