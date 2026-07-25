package cs2gsi

import (
	"os"
	"testing"

	models "github.com/nescabir/go-cs2-gsi/models"
)

func seedLastState(gsi *CS2GSI) {
	attacker := &models.Player{
		SteamId: "76561198000000001",
		Name:    "Attacker",
		State:   &models.PlayerState{},
	}
	victim := &models.Player{
		SteamId: "76561198000000002",
		Name:    "Victim",
		State:   &models.PlayerState{},
	}

	gsi.last = &models.State{
		AllPlayers: map[string]*models.Player{
			attacker.SteamId: attacker,
			victim.SteamId:   victim,
		},
	}
}

func TestDigestMIRVKill(t *testing.T) {
	gsi := New(NewConfig())
	seedLastState(gsi)

	raw, err := os.ReadFile("testdata/mirv/player_death.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	result, err := gsi.DigestMIRV(raw, MIRVEventPlayerDeath)
	if err != nil {
		t.Fatalf("DigestMIRV: %v", err)
	}
	if result.Kill == nil {
		t.Fatal("expected kill event")
	}
	if result.Kill.Victim.SteamId != "76561198000000002" {
		t.Fatalf("victim steam id = %q", result.Kill.Victim.SteamId)
	}
	if result.Kill.Attacker.SteamId != "76561198000000001" {
		t.Fatalf("attacker steam id = %q", result.Kill.Attacker.SteamId)
	}
	if !result.Kill.Headshot {
		t.Fatal("expected headshot")
	}
	if !result.Kill.Wallbang {
		t.Fatal("expected wallbang")
	}
	if result.Kill.Weapon.Name != "ak47" {
		t.Fatalf("weapon = %q, want ak47", result.Kill.Weapon.Name)
	}
}

func TestDigestMIRVHurt(t *testing.T) {
	gsi := New(NewConfig())
	seedLastState(gsi)

	raw, err := os.ReadFile("testdata/mirv/player_hurt.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	result, err := gsi.DigestMIRV(raw, MIRVEventPlayerHurt)
	if err != nil {
		t.Fatalf("DigestMIRV: %v", err)
	}
	if result.Hurt == nil {
		t.Fatal("expected hurt event")
	}
	if result.Hurt.DmgHealth != 55 {
		t.Fatalf("dmg_health = %d, want 55", result.Hurt.DmgHealth)
	}
}

func TestDigestMIRVNoPriorState(t *testing.T) {
	gsi := New(NewConfig())

	_, err := gsi.DigestMIRV([]byte(`{}`), MIRVEventPlayerDeath)
	if err != ErrMIRVNoPriorState {
		t.Fatalf("err = %v, want ErrMIRVNoPriorState", err)
	}
}
