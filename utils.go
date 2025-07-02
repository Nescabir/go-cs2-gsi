package cs2gsi

import (
	"math"
	"strconv"
	"strings"

	models "github.com/nescabir/go-cs2-gsi/models"
	rawModels "github.com/nescabir/go-cs2-gsi/raw"
)

func (gsi *CS2GSI) initState() *models.State {
	return &models.State{
		Provider: &models.Provider{
			Name:      "",
			AppId:     0,
			Version:   0,
			SteamId:   "",
			Timestamp: 0,
		},
		Map: &models.Map{
			Mode:                      "",
			Name:                      "",
			Phase:                     "",
			Round:                     0,
			Team_ct:                   &models.Team{},
			Team_t:                    &models.Team{},
			Num_matches_to_win_series: 0,
			Current_spectators:        0,
			Souvenirs_total:           0,
			Round_wins:                make(map[string]models.RoundOutcome),
			Rounds:                    make([]models.RoundInfo, 0, 60),
		},
		Round: &models.Round{
			Phase:    "",
			Win_team: "",
			Bomb:     "",
		},
		Player: &models.Player{
			SteamId:     "",
			Clan:        "",
			Name:        "",
			Team:        &models.Team{},
			Activity:    "",
			State:       &models.PlayerState{},
			Weapons:     make(map[string]*models.Weapon),
			Match_stats: &models.PlayerMatchStats{},
			Position:    [3]float32{},
			Forward:     [3]float32{},
			Avatar:      "",
		},
		Observer: &models.Observer{
			Activity:   "",
			Spectarget: "",
			Position:   [3]float32{},
			Forward:    [3]float32{},
		},
		AllPlayers: make(map[string]*models.Player),
		Bomb: &models.Bomb{
			State:     "",
			Countdown: 0,
			Player:    &models.Player{},
			Position:  [3]float32{},
			Site:      "",
		},
		Grenades: make(map[string]*models.Grenade),
		Phase_countdowns: &models.PhaseCountdown{
			Phase:         "",
			Phase_ends_in: 0,
		},
		Auth: &models.Auth{
			Token: "",
		},
	}
}

func (gsi *CS2GSI) parseProvider(raw *rawModels.Provider) *models.Provider {
	return &models.Provider{
		Name:      raw.Name,
		AppId:     raw.AppId,
		Version:   raw.Version,
		SteamId:   raw.SteamId,
		Timestamp: float32(raw.Timestamp),
	}
}

func (gsi *CS2GSI) parseMap(raw *rawModels.Map) *models.Map {
	return &models.Map{
		Mode:                      raw.Mode,
		Name:                      raw.Name,
		Phase:                     parseMapPhase(string(raw.Phase)),
		Round:                     raw.Round,
		Team_ct:                   gsi.parseTeam(raw.Team_ct, models.CTSide),
		Team_t:                    gsi.parseTeam(raw.Team_t, models.TSide),
		Num_matches_to_win_series: raw.Num_matches_to_win_series,
		Current_spectators:        raw.Current_spectators,
		Souvenirs_total:           raw.Souvenirs_total,
		Round_wins:                make(map[string]models.RoundOutcome),
		Rounds:                    make([]models.RoundInfo, 0, 60),
	}
}

func (gsi *CS2GSI) parseRound(raw *rawModels.Round) *models.Round {
	return &models.Round{
		Phase:    parseRoundPhase(string(raw.Phase)),
		Win_team: parseSide(string(raw.Win_team)),
		Bomb:     parseBombRoundState(string(raw.Bomb)),
	}
}

func (gsi *CS2GSI) parsePlayerState(raw *rawModels.PlayerState) *models.PlayerState {
	if raw == nil {
		return &models.PlayerState{
			Health:         0,
			Armor:          0,
			Helmet:         false,
			DefuseKit:      false,
			Flashed:        0,
			Smoked:         0,
			Burning:        0,
			Money:          0,
			Round_kills:    0,
			Round_killhs:   0,
			Round_totaldmg: 0,
			Equip_value:    0,
			Adr:            0,
		}
	}

	return &models.PlayerState{
		Health:         raw.Health,
		Armor:          raw.Armor,
		Helmet:         raw.Helmet,
		DefuseKit:      raw.DefuseKit,
		Flashed:        raw.Flashed,
		Smoked:         raw.Smoked,
		Burning:        raw.Burning,
		Money:          raw.Money,
		Round_kills:    raw.Round_kills,
		Round_killhs:   raw.Round_killhs,
		Round_totaldmg: raw.Round_totaldmg,
		Equip_value:    raw.Equip_value,
		Adr:            0,
	}
}

func (gsi *CS2GSI) parseWeapon(raw *rawModels.Weapon) *models.Weapon {
	if raw == nil {
		return nil
	}

	return &models.Weapon{
		Name:      raw.Name,
		PaintKit:  raw.PaintKit,
		Type:      parseWeaponType(string(raw.Type)),
		State:     parseWeaponState(string(raw.State)),
		Ammo_clip: raw.Ammo_clip,
	}
}

func (gsi *CS2GSI) parseWeapons(raw map[string]*rawModels.Weapon) map[string]*models.Weapon {
	weapons := make(map[string]*models.Weapon)
	if len(raw) == 0 {
		return weapons
	}

	for _, weapon := range raw {
		if weapon == nil {
			continue
		}

		if weapon.Name == "" {
			continue
		}

		parsedWeapon := gsi.parseWeapon(weapon)
		if parsedWeapon != nil {
			weapons[weapon.Name] = parsedWeapon
		}
	}
	return weapons
}

func (gsi *CS2GSI) parsePlayerMatchStats(raw *rawModels.PlayerMatchStats) *models.PlayerMatchStats {
	if raw == nil {
		return &models.PlayerMatchStats{
			Kills:   0,
			Assists: 0,
			Deaths:  0,
			Mvps:    0,
			Score:   0,
		}
	}

	return &models.PlayerMatchStats{
		Kills:   raw.Kills,
		Assists: raw.Assists,
		Deaths:  raw.Deaths,
		Mvps:    raw.Mvps,
		Score:   raw.Score,
	}
}

func (gsi *CS2GSI) parseAuth(raw *rawModels.Auth) *models.Auth {
	return &models.Auth{
		Token: raw.Token,
	}
}

func (gsi *CS2GSI) parsePhaseCountdown(raw *rawModels.PhaseCountdown) *models.PhaseCountdown {
	if raw == nil {
		return nil
	}

	var phase_ends_in float64 = 0
	var err error

	if raw.Phase_ends_in != "" {
		phase_ends_in, err = strconv.ParseFloat(raw.Phase_ends_in, 64)
		if err != nil {
			gsi.logger.Warn("failed to parse phase_ends_in", "value", raw.Phase_ends_in, "error", err)
			phase_ends_in = 0
		}
	}

	return &models.PhaseCountdown{
		Phase:         parsePhaseType(string(raw.Phase)),
		Phase_ends_in: float32(phase_ends_in),
	}
}

func (gsi *CS2GSI) ParseGrenade(raw *rawModels.Grenade) *models.Grenade {
	if raw == nil {
		return nil
	}

	position := parseVector(raw.Position)
	velocity := parseVector(raw.Velocity)

	var lifetime float64 = 0
	var err error

	if raw.Lifetime != "" {
		lifetime, err = strconv.ParseFloat(raw.Lifetime, 64)
		if err != nil {
			gsi.logger.Warn("failed to parse grenade lifetime", "value", raw.Lifetime, "error", err)
			lifetime = 0
		}
	}

	return &models.Grenade{
		Owner:      raw.Owner,
		Position:   position,
		Velocity:   velocity,
		Type:       parseGrenadeType(string(raw.Type)),
		Lifetime:   float32(lifetime),
		EffectTime: float32(raw.EffectTime),
	}
}

func (gsi *CS2GSI) ParseGrenades(raw map[string]*rawModels.Grenade) map[string]*models.Grenade {
	grenades := make(map[string]*models.Grenade)
	if len(raw) > 0 {
		for _, grenade := range raw {
			grenades[string(grenade.Type)] = gsi.ParseGrenade(grenade)
		}
	}
	return grenades
}

func (gsi *CS2GSI) parseBomb(raw *rawModels.Bomb, allPlayers map[string]*rawModels.Player, teams *teams, mapName *string) *models.Bomb {
	if raw == nil {
		return nil
	}

	if mapName == nil || *mapName == "" {
		gsi.logger.Warn("map name is nil or empty for bomb parsing")
		return nil
	}

	var countdown float32 = 0
	if raw.Countdown != "" {
		cd, err := strconv.ParseFloat(raw.Countdown, 64)
		if err != nil {
			gsi.logger.Warn("failed to parse bomb countdown", "value", raw.Countdown, "error", err)
			countdown = 0
		} else {
			countdown = float32(cd)
		}
	}

	position := parseVector(raw.Position)

	// Validate player reference
	var player *models.Player
	if raw.Player != "" {
		rawPlayer, ok := allPlayers[raw.Player]
		if !ok {
			gsi.logger.Warn("bomb player not found in allPlayers", "player", raw.Player)
		} else {
			player = gsi.parsePlayer(rawPlayer, teams)
		}
	}

	return &models.Bomb{
		State:     parseBombState(string(raw.State)),
		Countdown: countdown,
		Player:    player,
		Position:  position,
		Site:      findSite(*mapName, position),
	}
}

func (gsi *CS2GSI) parsePlayer(raw *rawModels.Player, teams *teams) *models.Player {
	if raw == nil {
		return nil
	}

	if teams == nil {
		gsi.logger.Warn("teams is nil for player parsing")
		return nil
	}

	var team *models.Team
	switch raw.Team {
	case rawModels.CTSide:
		team = teams.ct
	case rawModels.TSide:
		team = teams.t
	default:
		gsi.logger.Warn("unknown team for player", "team", raw.Team, "steamId", raw.Steamid)
		team = &models.Team{} // Default empty team
	}

	position := parseVector(raw.Position)
	forward := parseVector(raw.Forward)

	player := &models.Player{
		SteamId:       raw.Steamid,
		Clan:          raw.Clan,
		Name:          raw.Name,
		Observer_slot: raw.Observer_slot,
		Team:          team,
		Activity:      parsePlayerActivity(string(raw.Activity)),
		State:         gsi.parsePlayerState(raw.State),
		Weapons:       gsi.parseWeapons(raw.Weapons),
		Match_stats:   gsi.parsePlayerMatchStats(raw.Match_stats),
		Position:      position,
		Forward:       forward,
		Avatar:        "",
	}

	return player
}

func (gsi *CS2GSI) parseTeam(raw *rawModels.Team, side models.Side) *models.Team {
	if raw == nil {
		return &models.Team{
			Logo:                     "",
			Score:                    0,
			Consecutive_round_losses: 0,
			Timeouts_remaining:       0,
			Matches_won_this_series:  0,
			Name:                     "",
			Flag:                     "",
			Side:                     side,
		}
	}

	return &models.Team{
		Logo:                     raw.Logo,
		Score:                    raw.Score,
		Consecutive_round_losses: raw.Consecutive_round_losses,
		Timeouts_remaining:       raw.Timeouts_remaining,
		Matches_won_this_series:  raw.Matches_won_this_series,
		Name:                     raw.Name,
		Flag:                     raw.Flag,
		Side:                     side,
	}
}

func parseVector(raw string) [3]float32 {
	if raw == "" {
		return [3]float32{0, 0, 0}
	}

	parts := strings.Split(raw, ", ")
	if len(parts) != 3 {
		return [3]float32{0, 0, 0}
	}

	vector := [3]float32{}
	for i := range 3 {
		vec, err := strconv.ParseFloat(strings.TrimSpace(parts[i]), 64)
		if err != nil {
			return [3]float32{0, 0, 0}
		}
		vector[i] = float32(vec)
	}
	return vector
}

func getRoundWin(mapRound int, teams *teams, roundWins map[string]rawModels.RoundOutcome, round int, regulationMR int, overtimeMR int) *models.RoundInfo {
	var indexRound = round
	if mapRound > 2*regulationMR {

		maxOvertimeRounds := 2*overtimeMR*int(math.Floor(float64(mapRound-(2*regulationMR+1))/float64(2*overtimeMR))) + 2*regulationMR

		if round <= int(maxOvertimeRounds) {
			return nil
		}

		roundInOT := ((round - (2*regulationMR + 1)) % (overtimeMR * 2)) + 1
		indexRound = roundInOT
	}

	roundOutcome := roundWins[strconv.Itoa(indexRound)]
	if roundOutcome == "" {
		return nil
	}

	winSide := strings.ToUpper(strings.Split(string(roundOutcome), "_")[0])

	var result = &models.RoundInfo{
		Team:    teams.ct,
		Round:   round,
		Side:    parseSide(winSide),
		Outcome: parseRoundOutcome(string(roundOutcome)),
	}

	if didTeamWinThatRound(teams.ct, round, models.Side(winSide), mapRound, regulationMR, overtimeMR) {
		return result
	}

	result.Team = teams.t

	return result
}

func didTeamWinThatRound(team *models.Team, round int, winSide models.Side, mapRound int, regulationMR int, overtimeMR int) bool {
	currentRound := 1
	currentRoundHalf := getHalfFromRound(currentRound, regulationMR, overtimeMR)
	roundToCheckHalf := getHalfFromRound(round, regulationMR, overtimeMR)

	return (team.Side == winSide) == (currentRoundHalf == roundToCheckHalf)
}

func getHalfFromRound(round int, regulationMR int, overtimeMR int) int {
	currentRoundHalf := 1
	if round <= 2*regulationMR {
		if round <= regulationMR {
			currentRoundHalf = 1
		} else {
			currentRoundHalf = 2
		}
	} else {
		roundInOT := ((round - (2*regulationMR + 1)) % (overtimeMR * 2)) + 1
		if roundInOT <= overtimeMR {
			currentRoundHalf = 1
		} else {
			currentRoundHalf = 2
		}
	}
	return currentRoundHalf
}

func findSite(mapName string, position [3]float32) models.BombSite {
	realMapName := mapName[strings.LastIndex(mapName, "/")+1:]
	switch realMapName {
	case "de_mirage":
		if position[1] < -600 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_cache":
		if position[1] > 0 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_overpass":
		if position[2] > 400 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_nuke":
		if position[2] > -500 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_dust2":
		if position[0] > -500 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_inferno":
		if position[0] > 1400 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_vertigo":
		if position[0] > -1400 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_train":
		if position[1] > -450 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_ancient":
		if position[0] < -500 {
			return models.BombSiteA
		}
		return models.BombSiteB
	case "de_anubis":
		if position[0] > 0 {
			return models.BombSiteA
		}
		return models.BombSiteB
	}
	return models.BombSite("")
}

func parseRoundPhase(s string) models.RoundPhase {
	roundPhases := map[rawModels.RoundPhase]models.RoundPhase{
		rawModels.RoundPhaseFreezeTime: models.RoundPhaseFreezeTime,
		rawModels.RoundPhaseLive:       models.RoundPhaseLive,
		rawModels.RoundPhaseOver:       models.RoundPhaseOver,
	}

	rp := rawModels.RoundPhase(s)
	_, ok := roundPhases[rp]
	if !ok {
		return models.RoundPhase("")
	}

	return roundPhases[rp]
}

func parseSide(s string) models.Side {
	roundWinTeams := map[rawModels.Side]models.Side{
		rawModels.CTSide: models.CTSide,
		rawModels.TSide:  models.TSide,
	}

	rwt := rawModels.Side(s)
	_, ok := roundWinTeams[rwt]
	if !ok {
		return models.Side("")
	}

	return roundWinTeams[rwt]
}

func parseBombRoundState(s string) models.BombRoundState {
	roundBombStates := map[rawModels.BombRoundState]models.BombRoundState{
		rawModels.BombRoundStatePlanted:  models.BombRoundStatePlanted,
		rawModels.BombRoundStateExploded: models.BombRoundStateExploded,
		rawModels.BombRoundStateDefused:  models.BombRoundStateDefused,
	}

	rb := rawModels.BombRoundState(s)
	_, ok := roundBombStates[rb]
	if !ok {
		return models.BombRoundState("")
	}

	return roundBombStates[rb]
}

func parsePlayerActivity(s string) models.PlayerActivity {
	playerActivities := map[rawModels.PlayerActivity]models.PlayerActivity{
		rawModels.PlayerActivityActive:    models.PlayerActivityActive,
		rawModels.PlayerActivityMenu:      models.PlayerActivityMenu,
		rawModels.PlayerActivityTextInput: models.PlayerActivityTextInput,
	}

	pa := rawModels.PlayerActivity(s)
	_, ok := playerActivities[pa]
	if !ok {
		return models.PlayerActivity("")
	}

	return playerActivities[pa]
}

func parseMapPhase(s string) models.MapPhase {
	mapPhases := map[rawModels.MapPhase]models.MapPhase{
		rawModels.MapPhaseWarmup:       models.MapPhaseWarmup,
		rawModels.MapPhaseLive:         models.MapPhaseLive,
		rawModels.MapPhaseIntermission: models.MapPhaseIntermission,
		rawModels.MapPhaseGameOver:     models.MapPhaseGameOver,
	}

	mp := rawModels.MapPhase(s)
	_, ok := mapPhases[mp]
	if !ok {
		return models.MapPhase("")
	}

	return mapPhases[mp]
}

func parseRoundOutcome(s string) models.RoundOutcome {
	roundOutcomes := map[rawModels.RoundOutcome]models.RoundOutcome{
		rawModels.CTWinElimination: models.CTWinElimination,
		rawModels.TWinElimination:  models.TWinElimination,
		rawModels.CTWinTimeLimit:   models.CTWinTimeLimit,
		rawModels.CTWinDefuse:      models.CTWinDefuse,
		rawModels.TWinBomb:         models.TWinBomb,
	}

	ro := rawModels.RoundOutcome(s)
	_, ok := roundOutcomes[ro]
	if !ok {
		return models.RoundOutcome("")
	}

	return roundOutcomes[ro]
}

func parseBombState(s string) models.BombState {
	bombStates := map[rawModels.BombState]models.BombState{
		rawModels.BombStateCarried:  models.BombStateCarried,
		rawModels.BombStateDropped:  models.BombStateDropped,
		rawModels.BombStatePlanted:  models.BombStatePlanted,
		rawModels.BombStateExploded: models.BombStateExploded,
		rawModels.BombStateDefused:  models.BombStateDefused,
		rawModels.BombStateDefusing: models.BombStateDefusing,
		rawModels.BombStatePlanting: models.BombStatePlanting,
	}

	bs := rawModels.BombState(s)
	_, ok := bombStates[bs]
	if !ok {
		return models.BombState("")
	}

	return bombStates[bs]
}

func parsePhaseType(s string) models.PhaseType {
	phaseTypes := map[rawModels.PhaseType]models.PhaseType{
		rawModels.PhaseTypeFreezetime: models.PhaseTypeFreezetime,
		rawModels.PhaseTypeBomb:       models.PhaseTypeBomb,
		rawModels.PhaseTypeWarmup:     models.PhaseTypeWarmup,
		rawModels.PhaseTypeLive:       models.PhaseTypeLive,
		rawModels.PhaseTypeOver:       models.PhaseTypeOver,
		rawModels.PhaseTypeDefuse:     models.PhaseTypeDefuse,
		rawModels.PhaseTypePaused:     models.PhaseTypePaused,
		rawModels.PhaseTypeTimeoutCT:  models.PhaseTypeTimeoutCT,
		rawModels.PhaseTypeTimeoutT:   models.PhaseTypeTimeoutT,
	}

	pt := rawModels.PhaseType(s)
	_, ok := phaseTypes[pt]
	if !ok {
		return models.PhaseType("")
	}

	return phaseTypes[pt]
}

func parseGrenadeType(s string) models.GrenadeType {
	grenadeTypes := map[rawModels.GrenadeType]models.GrenadeType{
		rawModels.GrenadeTypeFlash:      models.GrenadeTypeFlash,
		rawModels.GrenadeTypeDecoy:      models.GrenadeTypeDecoy,
		rawModels.GrenadeTypeFrag:       models.GrenadeTypeFrag,
		rawModels.GrenadeTypeSmoke:      models.GrenadeTypeSmoke,
		rawModels.GrenadeTypeMolotov:    models.GrenadeTypeMolotov,
		rawModels.GrenadeTypeIncendiary: models.GrenadeTypeIncendiary,
	}

	gt := rawModels.GrenadeType(s)
	_, ok := grenadeTypes[gt]
	if !ok {
		return models.GrenadeType("")
	}

	return grenadeTypes[gt]
}

func parseWeaponState(s string) models.WeaponState {
	weaponStates := map[rawModels.WeaponState]models.WeaponState{
		rawModels.WeaponStateActive:    models.WeaponStateActive,
		rawModels.WeaponStateHolstered: models.WeaponStateHolstered,
		rawModels.WeaponStateReloading: models.WeaponStateReloading,
	}

	ws := rawModels.WeaponState(s)
	_, ok := weaponStates[ws]
	if !ok {
		return models.WeaponState("")
	}

	return weaponStates[ws]
}

func parseWeaponType(s string) models.WeaponType {
	weaponTypes := map[rawModels.WeaponType]models.WeaponType{
		rawModels.WeaponTypeKnife:         models.WeaponTypeKnife,
		rawModels.WeaponTypePistol:        models.WeaponTypePistol,
		rawModels.WeaponTypeGrenade:       models.WeaponTypeGrenade,
		rawModels.WeaponTypeRifle:         models.WeaponTypeRifle,
		rawModels.WeaponTypeSniperRifle:   models.WeaponTypeSniperRifle,
		rawModels.WeaponTypeC4:            models.WeaponTypeC4,
		rawModels.WeaponTypeSubmachineGun: models.WeaponTypeSubmachineGun,
		rawModels.WeaponTypeShotgun:       models.WeaponTypeShotgun,
		rawModels.WeaponTypeMachineGun:    models.WeaponTypeMachineGun,
	}

	wt := rawModels.WeaponType(s)
	_, ok := weaponTypes[wt]
	if !ok {
		return models.WeaponType("")
	}

	return weaponTypes[wt]
}
