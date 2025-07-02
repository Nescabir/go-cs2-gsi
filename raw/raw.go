package raw

import (
	"encoding/json"
)

type Side string

const (
	TSide   Side = "T"
	CTSide  Side = "CT"
	NilSide Side = ""
)

func (s Side) String() string {
	return string(s)
}

type PlayerActivity string

const (
	PlayerActivityActive    PlayerActivity = "active"
	PlayerActivityMenu      PlayerActivity = "menu"
	PlayerActivityTextInput PlayerActivity = "textinput"
	PlayerActivityNil       PlayerActivity = ""
)

func (p PlayerActivity) String() string {
	return string(p)
}

type MapPhase string

const (
	MapPhaseWarmup       MapPhase = "warmup"
	MapPhaseLive         MapPhase = "live"
	MapPhaseIntermission MapPhase = "intermission"
	MapPhaseGameOver     MapPhase = "gameover"
	MapPhaseNil          MapPhase = ""
)

func (m MapPhase) String() string {
	return string(m)
}

type RoundOutcome string

const (
	CTWinElimination RoundOutcome = "ct_win_elimination"
	TWinElimination  RoundOutcome = "t_win_elimination"
	CTWinTimeLimit   RoundOutcome = "ct_win_time"
	CTWinDefuse      RoundOutcome = "ct_win_defuse"
	TWinBomb         RoundOutcome = "t_win_bomb"
	RoundOutcomeNil  RoundOutcome = ""
)

func (r RoundOutcome) String() string {
	return string(r)
}

type BombRoundState string

const (
	BombRoundStatePlanted  BombRoundState = "planted"
	BombRoundStateExploded BombRoundState = "exploded"
	BombRoundStateDefused  BombRoundState = "defused"
	BombRoundStateNil      BombRoundState = ""
)

func (b BombRoundState) String() string {
	return string(b)
}

type RoundPhase string

const (
	RoundPhaseFreezeTime RoundPhase = "freezetime"
	RoundPhaseLive       RoundPhase = "live"
	RoundPhaseOver       RoundPhase = "over"
	RoundPhaseNil        RoundPhase = ""
)

func (r RoundPhase) String() string {
	return string(r)
}

type BombState string

const (
	BombStateCarried  BombState = "carried"
	BombStateDropped  BombState = "dropped"
	BombStatePlanted  BombState = "planted"
	BombStateExploded BombState = "exploded"
	BombStateDefused  BombState = "defused"
	BombStateDefusing BombState = "defusing"
	BombStatePlanting BombState = "planting"
	BombStateNil      BombState = ""
)

func (b BombState) String() string {
	return string(b)
}

type PhaseType string

const (
	PhaseTypeFreezetime PhaseType = "freezetime"
	PhaseTypeBomb       PhaseType = "bomb"
	PhaseTypeWarmup     PhaseType = "warmup"
	PhaseTypeLive       PhaseType = "live"
	PhaseTypeOver       PhaseType = "over"
	PhaseTypeDefuse     PhaseType = "defuse"
	PhaseTypePaused     PhaseType = "paused"
	PhaseTypeTimeoutCT  PhaseType = "timeout_ct"
	PhaseTypeTimeoutT   PhaseType = "timeout_t"
	PhaseTypeNil        PhaseType = ""
)

func (p PhaseType) String() string {
	return string(p)
}

type GrenadeType string

const (
	GrenadeTypeFlash      GrenadeType = "flash"
	GrenadeTypeDecoy      GrenadeType = "decoy"
	GrenadeTypeFrag       GrenadeType = "frag"
	GrenadeTypeSmoke      GrenadeType = "smoke"
	GrenadeTypeMolotov    GrenadeType = "firebomb"
	GrenadeTypeIncendiary GrenadeType = "inferno"
	GrenadeTypeNil        GrenadeType = ""
)

func (g GrenadeType) String() string {
	return string(g)
}

type WeaponState string

const (
	WeaponStateActive    WeaponState = "active"
	WeaponStateHolstered WeaponState = "holstered"
	WeaponStateReloading WeaponState = "reloading"
	WeaponStateNil       WeaponState = ""
)

func (w WeaponState) String() string {
	return string(w)
}

type WeaponType string

const (
	WeaponTypeKnife         WeaponType = "Knife"
	WeaponTypePistol        WeaponType = "Pistol"
	WeaponTypeGrenade       WeaponType = "Grenade"
	WeaponTypeRifle         WeaponType = "Rifle"
	WeaponTypeSniperRifle   WeaponType = "SniperRifle"
	WeaponTypeC4            WeaponType = "C4"
	WeaponTypeSubmachineGun WeaponType = "Submachine Gun"
	WeaponTypeShotgun       WeaponType = "Shotgun"
	WeaponTypeMachineGun    WeaponType = "Machine Gun"
	WeaponTypeNil           WeaponType = ""
)

func (w WeaponType) String() string {
	return string(w)
}

type State struct {
	Provider         *Provider           `json:"provider"`
	Map              *Map                `json:"map"`
	Round            *Round              `json:"round"`
	Player           *PlayerObserved     `json:"player"`
	AllPlayers       map[string]*Player  `json:"allplayers"`
	Bomb             *Bomb               `json:"bomb"`
	Grenades         map[string]*Grenade `json:"grenades"`
	Previously       *State              `json:"previously"`
	Added            *State              `json:"added"`
	Phase_countdowns *PhaseCountdown     `json:"phase_countdowns"`
	Auth             *Auth               `json:"auth"`
}

// UnmarshalJSON implements custom JSON unmarshaling for State
// This handles the specific case where grenades and previously fields contain boolean values
func (s *State) UnmarshalJSON(data []byte) error {
	// Create a temporary struct to unmarshal into
	type StateAlias State
	temp := &struct {
		*StateAlias
		Grenades   map[string]interface{} `json:"grenades"`
		Previously interface{}            `json:"previously"`
		Added      interface{}            `json:"added"`
	}{
		StateAlias: (*StateAlias)(s),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle grenades field manually - only process valid grenade objects
	if temp.Grenades != nil {
		s.Grenades = make(map[string]*Grenade)
		for key, value := range temp.Grenades {
			// Skip boolean values (false means no grenade in slot)
			if boolVal, ok := value.(bool); ok && !boolVal {
				continue
			}

			// Try to unmarshal as grenade object
			if grenadeData, err := json.Marshal(value); err == nil {
				var grenade Grenade
				if err := json.Unmarshal(grenadeData, &grenade); err == nil {
					s.Grenades[key] = &grenade
				}
			}
		}
	}

	// Handle previously field manually - skip boolean values
	if boolVal, ok := temp.Previously.(bool); ok && !boolVal {
		s.Previously = nil
	} else if temp.Previously != nil {
		// Skip processing nested State to avoid infinite recursion
		// The nested State will be processed by the main unmarshaling
		s.Previously = nil
	}

	// Handle added field manually - skip boolean values
	if boolVal, ok := temp.Added.(bool); ok && !boolVal {
		s.Added = nil
	} else if temp.Added != nil {
		// Skip processing nested State to avoid infinite recursion
		// The nested State will be processed by the main unmarshaling
		s.Added = nil
	}

	return nil
}

// provider
type Provider struct {
	Name      string  `json:"name"`
	AppId     int     `json:"appid"`
	Version   int     `json:"version"`
	SteamId   string  `json:"steamid"`
	Timestamp float32 `json:"timestamp"`
}

// map
type Map struct {
	Mode                      string                  `json:"mode"`
	Name                      string                  `json:"name"`
	Phase                     MapPhase                `json:"phase"`
	Round                     int                     `json:"round"`
	Team_ct                   *Team                   `json:"team_ct"`
	Team_t                    *Team                   `json:"team_t"`
	Num_matches_to_win_series int                     `json:"num_matches_to_win_series"`
	Current_spectators        int                     `json:"current_spectators"`
	Souvenirs_total           int                     `json:"souvenirs_total"`
	Round_wins                map[string]RoundOutcome `json:"round_wins"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Map
// This handles the specific case where round_wins field contains boolean values
func (m *Map) UnmarshalJSON(data []byte) error {
	// Create a temporary struct to unmarshal into
	type MapAlias Map
	temp := &struct {
		*MapAlias
		Round_wins map[string]interface{} `json:"round_wins"`
	}{
		MapAlias: (*MapAlias)(m),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle round_wins field manually - only process valid round outcome objects
	m.Round_wins = make(map[string]RoundOutcome)
	for key, value := range temp.Round_wins {
		// Skip boolean values (false means no round outcome)
		if boolVal, ok := value.(bool); ok && !boolVal {
			continue
		}

		// Try to unmarshal as round outcome string
		if outcomeStr, ok := value.(string); ok {
			m.Round_wins[key] = RoundOutcome(outcomeStr)
		}
	}

	return nil
}

// round
type Round struct {
	Phase    RoundPhase     `json:"phase"`
	Win_team Side           `json:"win_team"`
	Bomb     BombRoundState `json:"bomb"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Round
// This handles the specific case where bomb and win_team fields contain boolean values
func (r *Round) UnmarshalJSON(data []byte) error {
	// Create a temporary struct to unmarshal into
	type RoundAlias Round
	temp := &struct {
		*RoundAlias
		Bomb     interface{} `json:"bomb"`
		Win_team interface{} `json:"win_team"`
	}{
		RoundAlias: (*RoundAlias)(r),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle bomb field manually - skip boolean values
	if boolVal, ok := temp.Bomb.(bool); ok && !boolVal {
		r.Bomb = BombRoundState("")
	} else if bombStr, ok := temp.Bomb.(string); ok {
		r.Bomb = BombRoundState(bombStr)
	}

	// Handle win_team field manually - skip boolean values
	if boolVal, ok := temp.Win_team.(bool); ok && !boolVal {
		r.Win_team = Side("")
	} else if winTeamStr, ok := temp.Win_team.(string); ok {
		r.Win_team = Side(winTeamStr)
	}

	return nil
}

// player_id
type Player struct {
	Steamid       string             `json:"steamid"`
	Clan          string             `json:"clan"`
	Name          string             `json:"name"`
	Observer_slot int                `json:"observer_slot"`
	Team          Side               `json:"team"`
	Activity      PlayerActivity     `json:"activity"`
	State         *PlayerState       `json:"state"`
	Weapons       map[string]*Weapon `json:"weapons"`
	Match_stats   *PlayerMatchStats  `json:"match_stats"`
	Position      string             `json:"position"`
	Forward       string             `json:"forward"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Player
// This handles the specific case where weapons, steamid, and observer_slot fields contain boolean values
func (p *Player) UnmarshalJSON(data []byte) error {
	// Create a temporary struct to unmarshal into
	type PlayerAlias Player
	temp := &struct {
		*PlayerAlias
		Weapons       map[string]interface{} `json:"weapons"`
		Steamid       interface{}            `json:"steamid"`
		Observer_slot interface{}            `json:"observer_slot"`
	}{
		PlayerAlias: (*PlayerAlias)(p),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle weapons field manually - only process valid weapon objects
	p.Weapons = make(map[string]*Weapon)
	for key, value := range temp.Weapons {
		// Skip boolean values (false means no weapon in slot)
		if boolVal, ok := value.(bool); ok && !boolVal {
			continue
		}

		// Try to unmarshal as weapon object
		if weaponData, err := json.Marshal(value); err == nil {
			var weapon Weapon
			if err := json.Unmarshal(weaponData, &weapon); err == nil {
				p.Weapons[key] = &weapon
			}
		}
	}

	// Handle steamid field manually - skip boolean values
	if boolVal, ok := temp.Steamid.(bool); ok && !boolVal {
		p.Steamid = ""
	} else if steamidStr, ok := temp.Steamid.(string); ok {
		p.Steamid = steamidStr
	}

	// Handle observer_slot field manually - skip boolean values
	if boolVal, ok := temp.Observer_slot.(bool); ok && !boolVal {
		p.Observer_slot = -1 // Default to -1 for no observer slot
	} else if observerSlotFloat, ok := temp.Observer_slot.(float64); ok {
		p.Observer_slot = int(observerSlotFloat)
	} else if observerSlotInt, ok := temp.Observer_slot.(int); ok {
		p.Observer_slot = observerSlotInt
	}

	return nil
}

type PlayerObserved struct {
	Player
	Spectarget string `json:"spectarget"`
}

// team
type Team struct {
	Logo                     string `json:"logo"`
	Score                    int    `json:"score"`
	Consecutive_round_losses int    `json:"consecutive_round_losses"`
	Timeouts_remaining       int    `json:"timeouts_remaining"`
	Matches_won_this_series  int    `json:"matches_won_this_series"`
	Name                     string `json:"name"`
	Flag                     string `json:"flag"`
}

// player_state
type PlayerState struct {
	Health         int  `json:"health"`
	Armor          int  `json:"armor"`
	Helmet         bool `json:"helmet"`
	DefuseKit      bool `json:"defuse_kit"`
	Flashed        int  `json:"flashed"`
	Smoked         int  `json:"smoked"`
	Burning        int  `json:"burning"`
	Money          int  `json:"money"`
	Round_kills    int  `json:"round_kills"`
	Round_killhs   int  `json:"round_killhs"`
	Round_totaldmg int  `json:"round_totaldmg"`
	Equip_value    int  `json:"equip_value"`
}

// player_weapons: weapon_0, weapon_1, weapon_2 ...
type Weapon struct {
	Name          string      `json:"name"`
	PaintKit      string      `json:"paintkit"`
	Type          WeaponType  `json:"type"`
	State         WeaponState `json:"state"`
	Ammo_clip     int         `json:"ammo_clip"`
	Ammo_clip_max int         `json:"ammo_clip_max"`
	Ammo_reserve  int         `json:"ammo_reserve"`
}

// player_match_stats
type PlayerMatchStats struct {
	Kills   int `json:"kills"`
	Assists int `json:"assists"`
	Deaths  int `json:"deaths"`
	Mvps    int `json:"mvps"`
	Score   int `json:"score"`
}

type Auth struct {
	Token string `json:"token"`
}

type Bomb struct {
	State     BombState `json:"state"`
	Countdown string    `json:"countdown"`
	Player    string    `json:"player"`
	Position  string    `json:"position"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Bomb
// This handles the specific case where countdown and player fields contain boolean values
func (b *Bomb) UnmarshalJSON(data []byte) error {
	// Create a temporary struct to unmarshal into
	type BombAlias Bomb
	temp := &struct {
		*BombAlias
		Countdown interface{} `json:"countdown"`
		Player    interface{} `json:"player"`
	}{
		BombAlias: (*BombAlias)(b),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle countdown field manually - skip boolean values
	if boolVal, ok := temp.Countdown.(bool); ok && !boolVal {
		b.Countdown = ""
	} else if countdownStr, ok := temp.Countdown.(string); ok {
		b.Countdown = countdownStr
	}

	// Handle player field manually - skip boolean values
	if boolVal, ok := temp.Player.(bool); ok && !boolVal {
		b.Player = ""
	} else if playerStr, ok := temp.Player.(string); ok {
		b.Player = playerStr
	}

	return nil
}

type PhaseCountdown struct {
	Phase         PhaseType `json:"phase"`
	Phase_ends_in string    `json:"phase_ends_in"`
}

type Grenade struct {
	Owner      string      `json:"owner"`
	Position   string      `json:"position"`
	Velocity   string      `json:"velocity"`
	Type       GrenadeType `json:"type"`
	Lifetime   string      `json:"lifetime"`
	EffectTime float32     `json:"effect_time"`
}
