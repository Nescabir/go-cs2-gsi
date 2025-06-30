package cs2gsi

import (
	"math"
	"strconv"
	"strings"

	structs "github.com/nescabir/go-cs2-gsi/structs"
	structsRaw "github.com/nescabir/go-cs2-gsi/structsraw"
)

func (gsi *CS2GSI) parseProvider(raw *structsRaw.ProviderRaw) *structs.Provider {
	return &structs.Provider{
		Name:      raw.Name,
		AppId:     raw.AppId,
		Version:   raw.Version,
		SteamId:   raw.SteamId,
		Timestamp: float32(raw.Timestamp),
	}
}

func (gsi *CS2GSI) parseMap(raw *structsRaw.MapRaw) *structs.Map {
	return &structs.Map{
		Mode:                      raw.Mode,
		Name:                      raw.Name,
		Phase:                     structs.MapPhase(raw.Phase),
		Round:                     raw.Round,
		Team_ct:                   gsi.parseTeam(raw.Team_ct, structs.CTSide),
		Team_t:                    gsi.parseTeam(raw.Team_t, structs.TSide),
		Num_matches_to_win_series: raw.Num_matches_to_win_series,
		Current_spectators:        raw.Current_spectators,
		Souvenirs_total:           raw.Souvenirs_total,
		Round_wins:                make(map[string]structs.RoundOutcome),
		Rounds:                    make([]structs.RoundInfo, 0),
	}
}

func (gsi *CS2GSI) parseRound(raw *structsRaw.RoundRaw) *structs.Round {
	return &structs.Round{
		Phase:    structs.RoundPhase(raw.Phase),
		Win_team: structs.Side(raw.Win_team),
		Bomb:     structs.BombRoundState(raw.Bomb),
	}
}

func (gsi *CS2GSI) parsePlayerState(raw *structsRaw.PlayerStateRaw) *structs.PlayerState {
	return &structs.PlayerState{
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

func (gsi *CS2GSI) parseWeapon(raw *structsRaw.WeaponRaw) *structs.Weapon {
	return &structs.Weapon{
		Name:      raw.Name,
		PaintKit:  raw.PaintKit,
		Type:      structs.WeaponType(raw.Type),
		State:     structs.WeaponState(raw.State),
		Ammo_clip: raw.Ammo_clip,
	}
}

func (gsi *CS2GSI) parseWeapons(raw map[string]*structsRaw.WeaponRaw) map[string]*structs.Weapon {
	weapons := make(map[string]*structs.Weapon)
	for _, weapon := range raw {
		weapons[weapon.Name] = gsi.parseWeapon(weapon)
	}
	return weapons
}

func (gsi *CS2GSI) parsePlayerMatchStats(raw *structsRaw.PlayerMatchStatsRaw) *structs.PlayerMatchStats {
	return &structs.PlayerMatchStats{
		Kills:   raw.Kills,
		Assists: raw.Assists,
		Deaths:  raw.Deaths,
		Mvps:    raw.Mvps,
		Score:   raw.Score,
	}
}

func (gsi *CS2GSI) parseAuth(raw *structsRaw.AuthRaw) *structs.Auth {
	return &structs.Auth{
		Token: raw.Token,
	}
}

func (gsi *CS2GSI) parsePhaseCountdown(raw *structsRaw.PhaseCountdownRaw) *structs.PhaseCountdown {
	phase_ends_in, err := strconv.ParseFloat(raw.Phase_ends_in, 64)
	if err != nil {
		return nil
	}

	return &structs.PhaseCountdown{
		Phase:         structs.PhaseType(raw.Phase),
		Phase_ends_in: float32(phase_ends_in),
	}
}

func (gsi *CS2GSI) ParseGrenade(raw *structsRaw.GrenadeRaw) *structs.Grenade {
	position := ParseVector(raw.Position)
	velocity := ParseVector(raw.Velocity)

	lifetime, err := strconv.ParseFloat(raw.Lifetime, 64)
	if err != nil {
		return nil
	}

	return &structs.Grenade{
		Owner:      raw.Owner,
		Position:   position,
		Velocity:   velocity,
		Type:       structs.GrenadeType(raw.Type),
		Lifetime:   float32(lifetime),
		EffectTime: float32(raw.EffectTime),
	}
}

func (gsi *CS2GSI) ParseGrenades(raw map[string]*structsRaw.GrenadeRaw) map[string]*structs.Grenade {
	grenades := make(map[string]*structs.Grenade)
	for _, grenade := range raw {
		grenades[string(grenade.Type)] = gsi.ParseGrenade(grenade)
	}
	return grenades
}

func (gsi *CS2GSI) parseBomb(raw *structsRaw.BombRaw, allPlayers map[string]*structsRaw.PlayerRaw, teams *Teams, mapName *string) *structs.Bomb {
	countdown, err := strconv.ParseFloat(raw.Countdown, 64)
	if err != nil {
		return nil
	}

	position := ParseVector(raw.Position)

	player, ok := allPlayers[raw.Player]
	if !ok {
		return nil
	}

	return &structs.Bomb{
		State:     structs.BombState(raw.State),
		Countdown: float32(countdown),
		Player:    gsi.parsePlayer(player, teams),
		Position:  position,
		Site:      FindSite(*mapName, position),
	}
}

func (gsi *CS2GSI) parsePlayer(raw *structsRaw.PlayerRaw, teams *Teams) *structs.Player {

	var team *structs.Team
	if raw.Team == structsRaw.CTSide {
		team = teams.ct
	} else {
		team = teams.t
	}

	position := ParseVector(raw.Position)
	forward := ParseVector(raw.Forward)

	player := &structs.Player{
		SteamId:     raw.SteamId,
		Clan:        raw.Clan,
		Name:        raw.Name,
		Team:        team,
		Activity:    structs.PlayerActivity(raw.Activity),
		State:       gsi.parsePlayerState(raw.State),
		Weapons:     gsi.parseWeapons(raw.Weapons),
		Match_stats: gsi.parsePlayerMatchStats(raw.Match_stats),
		Position:    position,
		Forward:     forward,
		Avatar:      "",
	}

	return player
}

func (gsi *CS2GSI) parseTeam(raw *structsRaw.TeamRaw, side structs.Side) *structs.Team {
	return &structs.Team{
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

func ParseVector(raw string) [3]float32 {
	parts := strings.Split(raw, ",")
	vector := [3]float64{}
	for i := range 3 {
		vec, err := strconv.ParseFloat(parts[i], 64)
		vector[i] = vec
		if err != nil {
			return [3]float32{0, 0, 0}
		}
	}
	return [3]float32{float32(vector[0]), float32(vector[1]), float32(vector[2])}
}

func GetRoundWin(mapRound int, teams *Teams, roundWins map[string]structsRaw.RoundOutcome, round int, regulationMR int, overtimeMR int) *structs.RoundInfo {
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

	var result = &structs.RoundInfo{
		Team:    teams.ct,
		Round:   round,
		Side:    structs.Side(winSide),
		Outcome: structs.RoundOutcome(roundOutcome),
	}

	if didTeamWinThatRound(teams.ct, round, structs.Side(winSide), mapRound, regulationMR, overtimeMR) {
		return result
	}

	result.Team = teams.t

	return result
}

func didTeamWinThatRound(team *structs.Team, round int, winSide structs.Side, mapRound int, regulationMR int, overtimeMR int) bool {
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

func FindSite(mapName string, position [3]float32) structs.BombSite {
	realMapName := mapName[strings.LastIndex(mapName, "/")+1:]
	switch realMapName {
	case "de_mirage":
		if position[1] < -600 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_cache":
		if position[1] > 0 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_overpass":
		if position[2] > 400 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_nuke":
		if position[2] > -500 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_dust2":
		if position[0] > -500 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_inferno":
		if position[0] > 1400 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_vertigo":
		if position[0] > -1400 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_train":
		if position[1] > -450 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_ancient":
		if position[0] < -500 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	case "de_anubis":
		if position[0] > 0 {
			return structs.BombSiteA
		}
		return structs.BombSiteB
	}
	return structs.BombSite("")
}
