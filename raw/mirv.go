package raw

type MIRVUserID struct {
	Value int    `json:"value"`
	Xuid  string `json:"xuid"`
}

type RawKillKeys struct {
	Userid        MIRVUserID `json:"userid"`
	Attacker      MIRVUserID `json:"attacker"`
	Assister      MIRVUserID `json:"assister"`
	Assistedflash bool       `json:"assistedflash"`
	Weapon        string     `json:"weapon"`
	Headshot      bool       `json:"headshot"`
	Attackerblind bool       `json:"attackerblind"`
	Thrusmoke     bool       `json:"thrusmoke"`
	Noscope       bool       `json:"noscope"`
	Penetrated    int        `json:"penetrated"`
	Attackerinair bool       `json:"attackerinair"`
}

type RawKill struct {
	Name       string      `json:"name"`
	ClientTime float64     `json:"clientTime"`
	Keys       RawKillKeys `json:"keys"`
}

type RawHurtKeys struct {
	Userid    MIRVUserID `json:"userid"`
	Attacker  MIRVUserID `json:"attacker"`
	Health    int        `json:"health"`
	Armor     int        `json:"armor"`
	Weapon    string     `json:"weapon"`
	DmgHealth int        `json:"dmg_health"`
	DmgArmor  int        `json:"dmg_armor"`
	Hitgroup  int        `json:"hitgroup"`
}

type RawHurt struct {
	Name       string      `json:"name"`
	ClientTime float64     `json:"clientTime"`
	Keys       RawHurtKeys `json:"keys"`
}
