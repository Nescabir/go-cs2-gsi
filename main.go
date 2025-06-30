package cs2gsi

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"

	"github.com/asaskevich/EventBus"
	structs "github.com/nescabir/go-cs2-gsi/structs"
	structsRaw "github.com/nescabir/go-cs2-gsi/structsraw"
)

type CS2GSI struct {
	EventBus              EventBus.Bus
	REGULATION_MAX_ROUNDS int
	OVERTIME_MAX_ROUNDS   int
	Damage                []structs.RoundDamage
	Players               []structs.Player
	Teams                 *Teams
	current               *structs.State
	last                  *structs.State
}

type Teams struct {
	ct *structs.Team
	t  *structs.Team
}

func New(size int) *CS2GSI {
	return &CS2GSI{
		EventBus:              EventBus.New(),
		REGULATION_MAX_ROUNDS: 15,
		OVERTIME_MAX_ROUNDS:   3,
		Damage:                make([]structs.RoundDamage, 0),
		Players:               make([]structs.Player, 0),
		Teams: &Teams{
			ct: nil,
			t:  nil,
		},
		current: nil,
		last:    nil,
	}
}

func (gsi *CS2GSI) Listen(addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		stateRaw := &structsRaw.StateRaw{}

		if err := json.NewDecoder(r.Body).Decode(stateRaw); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}

		gsi.digest(stateRaw)

		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}

func (gsi *CS2GSI) digest(raw *structsRaw.StateRaw) *structs.State {
	state := &structs.State{}

	if len(raw.AllPlayers) <= 0 || raw.Map == nil || raw.Phase_countdowns == nil {
		return nil
	}

	state.Provider = gsi.parseProvider(raw.Provider)
	state.Map = gsi.parseMap(raw.Map)
	state.Round = gsi.parseRound(raw.Round)
	state.Phase_countdowns = gsi.parsePhaseCountdown(raw.Phase_countdowns)
	state.Auth = gsi.parseAuth(raw.Auth)
	state.Round = gsi.parseRound(raw.Round)
	state.Grenades = gsi.ParseGrenades(raw.Grenades)

	gsi.Teams = &Teams{
		ct: state.Map.Team_ct,
		t:  state.Map.Team_t,
	}

	state.Bomb = gsi.parseBomb(raw.Bomb, raw.AllPlayers, gsi.Teams, &raw.Map.Name)

	state.Observer = &structs.Observer{
		Activity:   structs.PlayerActivity(raw.Player.Activity),
		Spectarget: raw.Player.Spectarget,
		Position:   ParseVector(raw.Player.Position),
		Forward:    ParseVector(raw.Player.Forward),
	}

	for _, player := range raw.AllPlayers {
		if player.SteamId == raw.Player.PlayerRaw.SteamId {
			state.Player = gsi.parsePlayer(&raw.Player.PlayerRaw, gsi.Teams)
		}
		parsedPlayer := gsi.parsePlayer(player, gsi.Teams)
		state.AllPlayers[player.SteamId] = parsedPlayer
		gsi.Players = append(gsi.Players, *parsedPlayer)
	}

	var rounds []structs.RoundInfo
	if raw.Round != nil && raw.Map != nil && raw.Map.Round_wins != nil {
		var currentRound int = raw.Map.Round + 1

		if raw.Round != nil && (raw.Round.Phase == structsRaw.RoundPhaseOver || raw.Map.Phase == structsRaw.MapPhaseGameOver) {
			currentRound = raw.Map.Round
		}

		for i := 1; i <= currentRound; i++ {
			var result = GetRoundWin(currentRound, gsi.Teams, raw.Map.Round_wins, i, gsi.REGULATION_MAX_ROUNDS, gsi.OVERTIME_MAX_ROUNDS)
			if result == nil {
				continue
			}
			rounds = append(rounds, *result)
		}
	}
	state.Map.Rounds = rounds

	if gsi.last != nil && gsi.last.Map.Name != raw.Map.Name {
		gsi.Damage = make([]structs.RoundDamage, 0)
	}

	var currentRoundForDamage = raw.Map.Round + 1
	if raw.Round != nil && (raw.Round.Phase == structsRaw.RoundPhaseOver || raw.Map.Phase == structsRaw.MapPhaseGameOver) {
		currentRoundForDamage = raw.Map.Round
	}

	var currentRoundDamage *structs.RoundDamage
	for _, damage := range gsi.Damage {
		if damage.Round == currentRoundForDamage {
			currentRoundDamage = &damage
			break
		}
	}

	if currentRoundDamage != nil {
		currentRoundDamage.Round = currentRoundForDamage
		currentRoundDamage.Players = make([]structs.RoundPlayerDamage, 0)

		gsi.Damage = append(gsi.Damage, *currentRoundDamage)
	}

	if raw.Map.Round == 0 && raw.Phase_countdowns.Phase == structsRaw.PhaseTypeFreezetime || raw.Phase_countdowns.Phase == structsRaw.PhaseTypeWarmup {
		gsi.Damage = make([]structs.RoundDamage, 0)
	}

	if currentRoundDamage != nil {
		for _, player := range gsi.Players {
			currentRoundDamage.Players = append(currentRoundDamage.Players, structs.RoundPlayerDamage{
				SteamId: player.SteamId,
				Damage:  player.State.Round_totaldmg,
			})
		}
	}

	for _, player := range gsi.Players {
		if gsi.current == nil {
			continue
		}

		damageForRound := make([]structs.RoundDamage, 0)
		for _, damage := range gsi.Damage {
			if damage.Round < currentRoundForDamage {
				damageForRound = append(damageForRound, damage)
			}
		}

		if len(damageForRound) == 0 {
			continue
		}

		var damageEntries []int
		for _, damage := range damageForRound {
			var playerDamage int
			for _, playerDamageEntry := range damage.Players {
				if playerDamageEntry.SteamId == player.SteamId {
					playerDamage = playerDamageEntry.Damage
					break
				}
			}
			damageEntries = append(damageEntries, playerDamage)
		}

		var totalDamage int
		for _, damage := range damageEntries {
			totalDamage += damage
		}

		adr := float64(totalDamage) / float64(raw.Map.Round)
		if raw.Map.Round == 0 {
			adr = float64(totalDamage)
		}

		player.State.Adr = int(math.Floor(adr))
	}

	state.AllPlayers = make(map[string]*structs.Player)
	for _, player := range gsi.Players {
		state.AllPlayers[player.SteamId] = &player
	}

	gsi.current = state
	if gsi.last != nil {
		gsi.last = state
		gsi.EventBus.Publish(string(structs.Data), state)
		return state
	}

	last := gsi.last

	if last.Round != nil && state.Round != nil {
		var winner *structs.Team
		var loser *structs.Team

		if state.Round.Win_team == structs.CTSide {
			winner = state.Map.Team_ct
			loser = state.Map.Team_t
		} else {
			winner = state.Map.Team_t
			loser = state.Map.Team_ct
		}

		var oldWinner *structs.Team
		if state.Round.Win_team == structs.CTSide {
			oldWinner = last.Map.Team_ct
		} else {
			oldWinner = last.Map.Team_t
		}

		if winner.Score == oldWinner.Score {
			winner.Score += 1
		}

		roundScore := &structs.Score{
			Winner: winner,
			Loser:  loser,
			Map:    state.Map,
			MapEnd: state.Map.Phase == structs.MapPhaseGameOver,
		}
		gsi.EventBus.Publish(string(structs.RoundEnd), roundScore)

		// Match end
		if roundScore.MapEnd && last.Map.Phase != structs.MapPhaseGameOver {
			gsi.EventBus.Publish(string(structs.MatchEnd), roundScore)
		}
	}

	// Bomb actions
	if state.Bomb != nil && last.Bomb != nil {
		if last.Bomb.State == structs.BombStatePlanting && state.Bomb.State != structs.BombStatePlanting && state.Bomb.State != structs.BombStatePlanted && state.Bomb.State != structs.BombStateDefusing {
			gsi.EventBus.Publish(string(structs.BombPlantStop), last.Bomb.Player)
		}

		if last.Bomb.State == structs.BombStatePlanting && state.Bomb.State == structs.BombStatePlanted {
			gsi.EventBus.Publish(string(structs.BombPlanted), last.Bomb.Player)
		} else if last.Bomb.State != structs.BombStateExploded && state.Bomb.State == structs.BombStateExploded {
			gsi.EventBus.Publish(string(structs.BombExploded))
		} else if last.Bomb.State != structs.BombStateDefused && state.Bomb.State == structs.BombStateDefused {
			gsi.EventBus.Publish(string(structs.BombDefused), last.Bomb.Player)
		} else if last.Bomb.State != structs.BombStateDefusing && state.Bomb.State == structs.BombStateDefusing {
			gsi.EventBus.Publish(string(structs.DefuseStart), state.Bomb.Player)
		} else if last.Bomb.State == structs.BombStateDefusing && state.Bomb.State != structs.BombStateDefusing {
			gsi.EventBus.Publish(string(structs.DefuseEnd), last.Bomb.Player)
		} else if last.Bomb.State != structs.BombStatePlanting && state.Bomb.State == structs.BombStatePlanting {
			gsi.EventBus.Publish(string(structs.BombPlantStart), state.Bomb.Player)
		}
	} else if last.Bomb == nil && state.Bomb != nil && state.Bomb.State == structs.BombStateExploded {
		gsi.EventBus.Publish(string(structs.BombExploded))
	}

	// Intermission (between halfs)
	if state.Map.Phase == structs.MapPhaseIntermission && last.Map.Phase != structs.MapPhaseIntermission {
		gsi.EventBus.Publish(string(structs.IntermissionStart))
	} else if state.Map.Phase != structs.MapPhaseIntermission && last.Map.Phase == structs.MapPhaseIntermission {
		gsi.EventBus.Publish(string(structs.IntermissionEnd))
	}

	var phase = state.Phase_countdowns.Phase

	if phase == structs.PhaseTypeFreezetime && last.Phase_countdowns.Phase != structs.PhaseTypeFreezetime {
		gsi.EventBus.Publish(string(structs.FreezetimeStart))
	} else if phase != structs.PhaseTypeFreezetime && last.Phase_countdowns.Phase == structs.PhaseTypeFreezetime {
		gsi.EventBus.Publish(string(structs.FreezetimeEnd))
	}

	// Timeouts
	if strings.HasPrefix(string(phase), "timeout") && !strings.HasPrefix(string(last.Phase_countdowns.Phase), "timeout") {
		var team *structs.Team
		if phase == structs.PhaseTypeTimeoutCT {
			team = gsi.Teams.ct
		} else {
			team = gsi.Teams.t
		}
		gsi.EventBus.Publish(string(structs.TimeoutStart), team)
	} else if strings.HasPrefix(string(last.Phase_countdowns.Phase), "timeout") && !strings.HasPrefix(string(phase), "timeout") {
		gsi.EventBus.Publish(string(structs.TimeoutEnd))
	}

	// MVP
	var mvp *structs.Player
	for _, player := range state.AllPlayers {
		if previousPlayer, exists := last.AllPlayers[player.SteamId]; exists {
			if player.Match_stats.Mvps > previousPlayer.Match_stats.Mvps {
				mvp = player
				break
			}
		}
	}

	if mvp != nil {
		gsi.EventBus.Publish(string(structs.Mvp), mvp)
	}

	gsi.EventBus.Publish(string(structs.Data), state)
	gsi.last = state
	return state
}

func (gsi *CS2GSI) On(event structs.Events, callback func(state *structs.State)) {
	gsi.EventBus.Subscribe(string(event), callback)
}

func (gsi *CS2GSI) Once(event structs.Events, callback func(state *structs.State)) {
	gsi.EventBus.SubscribeOnce(string(event), callback)
}

func (gsi *CS2GSI) Off(event structs.Events, callback func(state *structs.State)) {
	gsi.EventBus.Unsubscribe(string(event), callback)
}

func (gsi *CS2GSI) Emit(event structs.Events, state *structs.State) {
	gsi.EventBus.Publish(string(event), state)
}
