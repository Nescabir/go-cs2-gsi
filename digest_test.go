package cs2gsi

import (
	"os"
	"sync"
	"testing"

	models "github.com/nescabir/go-cs2-gsi/models"
)

func TestRawEventBeforeData(t *testing.T) {
	gsi := New(NewConfig())
	var order []string
	var mu sync.Mutex

	Subscribe(Raw, func(e Event[[]byte]) {
		mu.Lock()
		order = append(order, "raw")
		mu.Unlock()
	})
	Subscribe(Data, func(e Event[*models.State]) {
		mu.Lock()
		order = append(order, "data")
		mu.Unlock()
	})

	raw, err := os.ReadFile("testdata/gsi/with_deltas.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := gsi.Digest(raw); err != nil {
		t.Fatalf("Digest: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 2 || order[0] != "raw" || order[1] != "data" {
		t.Fatalf("event order = %v, want [raw data]", order)
	}
}

func TestAuthTokenValidation(t *testing.T) {
	raw, err := os.ReadFile("testdata/gsi/with_deltas.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	gsi := New(Config{ExpectedToken: "secret"})
	if err := gsi.Digest(raw); err != nil {
		t.Fatalf("Digest with valid token: %v", err)
	}

	gsiBad := New(Config{ExpectedToken: "wrong"})
	if err := gsiBad.Digest(raw); !IsInvalidToken(err) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestParseDeltas(t *testing.T) {
	raw, err := os.ReadFile("testdata/gsi/with_deltas.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	gsi := New(NewConfig())
	if err := gsi.Digest(raw); err != nil {
		t.Fatalf("Digest: %v", err)
	}

	state := gsi.Snapshot()
	if state.Previously == nil || state.Previously.Player == nil {
		t.Fatal("expected previously.player")
	}
	if state.Previously.Player.State.Health != 100 {
		t.Fatalf("previously health = %d, want 100", state.Previously.Player.State.Health)
	}
	if state.Added == nil || state.Added.Player == nil {
		t.Fatal("expected added.player")
	}
	if state.Added.Player.State.Health != 54 {
		t.Fatalf("added health = %d, want 54", state.Added.Player.State.Health)
	}
}

func TestTimeoutTeam(t *testing.T) {
	raw, err := os.ReadFile("testdata/gsi/with_deltas.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	gsi := New(NewConfig())
	if err := gsi.Digest(raw); err != nil {
		t.Fatalf("Digest: %v", err)
	}

	phase := gsi.Snapshot().Phase_countdowns
	if phase.Timeout_team == nil {
		t.Fatal("expected timeout_team")
	}
	if phase.Timeout_team.Side != models.CTSide {
		t.Fatalf("timeout team side = %q, want CT", phase.Timeout_team.Side)
	}
}

func TestPlayerExtensions(t *testing.T) {
	raw, err := os.ReadFile("testdata/gsi/with_deltas.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	gsi := New(Config{
		PlayerExtensions: []models.PlayerExtension{{
			SteamId:  "76561198000000001",
			Name:     "CloudName",
			Avatar:   "https://example.com/a.png",
			Country:  "US",
			RealName: "Real Player",
		}},
	})
	if err := gsi.Digest(raw); err != nil {
		t.Fatalf("Digest: %v", err)
	}

	player := gsi.Snapshot().AllPlayers["76561198000000001"]
	if player.Name != "CloudName" {
		t.Fatalf("name = %q, want CloudName", player.Name)
	}
	if player.DefaultName != "Player1" {
		t.Fatalf("defaultName = %q, want Player1", player.DefaultName)
	}
	if player.Avatar != "https://example.com/a.png" {
		t.Fatalf("avatar = %q", player.Avatar)
	}
}

func TestDamageExportedOnState(t *testing.T) {
	raw, err := os.ReadFile("testdata/gsi/with_deltas.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	gsi := New(NewConfig())
	if err := gsi.Digest(raw); err != nil {
		t.Fatalf("Digest: %v", err)
	}

	state := gsi.Snapshot()
	if len(state.Damage) == 0 {
		t.Fatal("expected damage history on state")
	}
}
