package models

type Side string

const (
	TSide  Side = "T"
	CTSide Side = "CT"
)

type Orientation string

const (
	OrientationLeft  Orientation = "left"
	OrientationRight Orientation = "right"
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

type BombSite string

const (
	BombSiteA BombSite = "A"
	BombSiteB BombSite = "B"
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

type State struct {
	Provider         *Provider
	Map              *Map
	Round            *Round
	Player           *Player
	Observer         *Observer
	AllPlayers       map[string]*Player // allplayers_*: steamid64 ...
	Bomb             *Bomb
	Grenades         map[string]*Grenade
	Previously       *State
	Added            *State
	Phase_countdowns *PhaseCountdown
	Auth             *Auth
}

// provider
type Provider struct {
	Name      string
	AppId     int
	Version   int
	SteamId   string
	Timestamp float32
}

// map
type Map struct {
	Mode                      string
	Name                      string
	Phase                     MapPhase
	Round                     int
	Team_ct                   *Team
	Team_t                    *Team
	Num_matches_to_win_series int
	Current_spectators        int
	Souvenirs_total           int
	Round_wins                map[string]RoundOutcome
	Rounds                    []RoundInfo
}

// round
type Round struct {
	Phase    RoundPhase
	Win_team Side
	Bomb     BombRoundState
}

type RoundInfo struct {
	Team    *Team
	Round   int
	Side    Side
	Outcome RoundOutcome
}

// player_id
type Player struct {
	SteamId       string
	Clan          string
	Name          string
	Observer_slot int
	Team          *Team
	Activity      PlayerActivity
	State         *PlayerState
	Weapons       map[string]*Weapon
	Match_stats   *PlayerMatchStats
	Position      [3]float32
	Forward       [3]float32
	Avatar        string
}

func (p *Player) IsAlive() bool {
	return p.State.Health > 0
}

type Observer struct {
	Activity   PlayerActivity
	Spectarget string
	Position   [3]float32
	Forward    [3]float32
}

// team
type Team struct {
	Logo                     string
	Score                    int
	Consecutive_round_losses int
	Timeouts_remaining       int
	Matches_won_this_series  int
	Name                     string
	Flag                     string
	Side                     Side
	Orientation              Orientation
}

// player_state
type PlayerState struct {
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
	Adr            int
}

// player_weapons: weapon_0, weapon_1, weapon_2 ...
type Weapon struct {
	Name          string
	PaintKit      string
	Type          WeaponType
	State         WeaponState
	Ammo_clip     int
	Ammo_clip_max int
	Ammo_reserve  int
}

// player_match_stats
type PlayerMatchStats struct {
	Kills   int
	Assists int
	Deaths  int
	Mvps    int
	Score   int
}

type Auth struct {
	Token string
}

type Bomb struct {
	State     BombState
	Countdown float32
	Player    *Player
	Position  [3]float32
	Site      BombSite
}

type PhaseCountdown struct {
	Phase         PhaseType
	Phase_ends_in float32
}

type Grenade struct {
	Owner      string
	Position   [3]float32
	Velocity   [3]float32
	Type       GrenadeType
	Lifetime   float32
	EffectTime float32
}

type RoundPlayerDamage struct {
	SteamId string
	Damage  int
}

type RoundDamage struct {
	Round   int
	Players []RoundPlayerDamage
}

type Score struct {
	Winner *Team
	Loser  *Team
	Map    *Map
	MapEnd bool
}

type KillEvent struct {
	Attacker      *Player
	Victim        *Player
	Weapon        *Weapon
	Assister      *Player
	Flashed       bool
	Headshot      bool
	Wallbang      bool
	AttackerBlind bool
	ThruSmoke     bool
	NoScope       bool
	AttackerInAir bool
}

type HurtEvent struct {
	Attacker  *Player
	Victim    *Player
	Weapon    *Weapon
	Health    int
	Armor     int
	DmgHealth int
	DmgArmor  int
	HitGroup  int
}

type Events string

const (
	Data     Events = "data"
	RoundEnd Events = "roundEnd"
	// RoundStart        Events = "roundStart"
	Kill              Events = "kill"
	Hurt              Events = "hurt"
	TimeoutStart      Events = "timeoutStart"
	TimeoutEnd        Events = "timeoutEnd"
	Mvp               Events = "mvp"
	FreezetimeStart   Events = "freezetimeStart"
	FreezetimeEnd     Events = "freezetimeEnd"
	IntermissionStart Events = "intermissionStart"
	IntermissionEnd   Events = "intermissionEnd"
	DefuseStart       Events = "defuseStart"
	DefuseEnd         Events = "defuseEnd"
	BombPlantStart    Events = "bombPlantStart"
	BombPlantStop     Events = "bombPlantStop"
	BombPlanted       Events = "bombPlanted"
	BombDefused       Events = "bombDefused"
	BombExploded      Events = "bombExploded"
	// MapEnd            Events = "mapEnd"
	// MapStart          Events = "mapStart"
	MatchEnd Events = "matchEnd"
)
