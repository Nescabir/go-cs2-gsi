package structsraw

type Side string

const (
	TSide  Side = "T"
	CTSide Side = "CT"
)

type PlayerActivity string

const (
	PlayerActivityActive    PlayerActivity = "active"
	PlayerActivityMenu      PlayerActivity = "menu"
	PlayerActivityTextInput PlayerActivity = "textinput"
)

type MapPhase string

const (
	MapPhaseWarmup       MapPhase = "warmup"
	MapPhaseLive         MapPhase = "live"
	MapPhaseIntermission MapPhase = "intermission"
	MapPhaseGameOver     MapPhase = "gameover"
)

type RoundOutcome string

const (
	CTWinElimination RoundOutcome = "ct_win_elimination"
	TWinElimination  RoundOutcome = "t_win_elimination"
	CTWinTimeLimit   RoundOutcome = "ct_win_time"
	CTWinDefuse      RoundOutcome = "ct_win_defuse"
	TWinBomb         RoundOutcome = "t_win_bomb"
)

type BombRoundState string

const (
	BombRoundStatePlanted  BombRoundState = "planted"
	BombRoundStateExploded BombRoundState = "exploded"
	BombRoundStateDefused  BombRoundState = "defused"
)

type RoundPhase string

const (
	RoundPhaseFreezeTime RoundPhase = "freezetime"
	RoundPhaseLive       RoundPhase = "live"
	RoundPhaseOver       RoundPhase = "over"
)

type BombState string

const (
	BombStateCarried  BombState = "carried"
	BombStateDropped  BombState = "dropped"
	BombStatePlanted  BombState = "planted"
	BombStateExploded BombState = "exploded"
	BombStateDefused  BombState = "defused"
	BombStateDefusing BombState = "defusing"
	BombStatePlanting BombState = "planting"
)

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
)

type GrenadeType string

const (
	GrenadeTypeFlash      GrenadeType = "flash"
	GrenadeTypeDecoy      GrenadeType = "decoy"
	GrenadeTypeFrag       GrenadeType = "frag"
	GrenadeTypeSmoke      GrenadeType = "smoke"
	GrenadeTypeMolotov    GrenadeType = "firebomb"
	GrenadeTypeIncendiary GrenadeType = "inferno"
)

type WeaponState string

const (
	WeaponStateActive    WeaponState = "active"
	WeaponStateHolstered WeaponState = "holstered"
	WeaponStateReloading WeaponState = "reloading"
)

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
)

type StateRaw struct {
	Provider         *ProviderRaw
	Map              *MapRaw
	Round            *RoundRaw
	Player           *PlayerObservedRaw
	AllPlayers       map[string]*PlayerRaw // allplayers_*: steamid64 ...
	Bomb             *BombRaw
	Grenades         map[string]*GrenadeRaw
	Previously       *StateRaw
	Added            *StateRaw
	Phase_countdowns *PhaseCountdownRaw
	Auth             *AuthRaw
}

// provider
type ProviderRaw struct {
	Name      string
	AppId     int
	Version   int
	SteamId   string
	Timestamp float32
}

// map
type MapRaw struct {
	Mode                      string
	Name                      string
	Phase                     MapPhase
	Round                     int
	Team_ct                   *TeamRaw
	Team_t                    *TeamRaw
	Num_matches_to_win_series int
	Current_spectators        int
	Souvenirs_total           int
	Round_wins                map[string]RoundOutcome
}

// round
type RoundRaw struct {
	Phase    RoundPhase
	Win_team Side
	Bomb     BombRoundState
}

// player_id
type PlayerRaw struct {
	SteamId       string
	Clan          string
	Name          string
	Observer_slot int
	Team          Side
	Activity      PlayerActivity
	State         *PlayerStateRaw
	Weapons       map[string]*WeaponRaw
	Match_stats   *PlayerMatchStatsRaw
	Position      string
	Forward       string
}

type PlayerObservedRaw struct {
	PlayerRaw
	Spectarget string
}

// team
type TeamRaw struct {
	Logo                     string
	Score                    int
	Consecutive_round_losses int
	Timeouts_remaining       int
	Matches_won_this_series  int
	Name                     string
	Flag                     string
}

// player_state
type PlayerStateRaw struct {
	Health         int
	Armor          int
	Helmet         bool
	DefuseKit      bool
	Flashed        int
	Smoked         int
	Burning        int
	Money          int
	Round_kills    int
	Round_killhs   int
	Round_totaldmg int
	Equip_value    int
}

// player_weapons: weapon_0, weapon_1, weapon_2 ...
type WeaponRaw struct {
	Name          string
	PaintKit      string
	Type          WeaponType
	State         WeaponState
	Ammo_clip     int
	Ammo_clip_max int
	Ammo_reserve  int
}

// player_match_stats
type PlayerMatchStatsRaw struct {
	Kills   int
	Assists int
	Deaths  int
	Mvps    int
	Score   int
}

type AuthRaw struct {
	Token string
}

type BombRaw struct {
	State     BombState
	Countdown string
	Player    string
	Position  string
}

type PhaseCountdownRaw struct {
	Phase         PhaseType
	Phase_ends_in string
}

type GrenadeRaw struct {
	Owner      int
	Position   string
	Velocity   string
	Type       GrenadeType
	Lifetime   string
	EffectTime float32
}
