package cs2gsi

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	models "github.com/nescabir/go-cs2-gsi/models"
	rawModels "github.com/nescabir/go-cs2-gsi/raw"
)

func TestParseVector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want [3]float32
	}{
		{"comma space", "100.0, 200.0, 11600.0", [3]float32{100, 200, 11600}},
		{"comma only", "100.0,200.0,11600.0", [3]float32{100, 200, 11600}},
		{"empty", "", [3]float32{0, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseVector(tt.raw)
			if got != tt.want {
				t.Fatalf("parseVector(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestParseGrenadesDualSmoke(t *testing.T) {
	raw := loadGrenadeFixture(t, "testdata/grenades/dual_smoke.json")
	gsi := New(NewConfig())

	grenades := gsi.parseGrenades(raw)

	if len(grenades) != 2 {
		t.Fatalf("len(grenades) = %d, want 2", len(grenades))
	}

	smoke841, ok := grenades["841"]
	if !ok {
		t.Fatal("expected grenade key 841")
	}
	if smoke841.ID != "841" {
		t.Fatalf("smoke841.ID = %q, want 841", smoke841.ID)
	}
	if smoke841.Type != models.GrenadeTypeSmoke {
		t.Fatalf("smoke841.Type = %q, want smoke", smoke841.Type)
	}
	if math.Abs(float64(smoke841.EffectTime)-17.42) > 0.01 {
		t.Fatalf("smoke841.EffectTime = %v, want ~17.42", smoke841.EffectTime)
	}

	smoke842, ok := grenades["842"]
	if !ok {
		t.Fatal("expected grenade key 842")
	}
	if math.Abs(float64(smoke842.EffectTime)-16.10) > 0.01 {
		t.Fatalf("smoke842.EffectTime = %v, want ~16.10", smoke842.EffectTime)
	}

	if _, ok := grenades["smoke"]; ok {
		t.Fatal("grenades must not be keyed by type")
	}
}

func TestParseGrenadesInferno(t *testing.T) {
	raw := loadGrenadeFixture(t, "testdata/grenades/inferno.json")
	gsi := New(NewConfig())

	grenades := gsi.parseGrenades(raw)

	inferno, ok := grenades["900"]
	if !ok {
		t.Fatal("expected grenade key 900")
	}
	if inferno.Type != models.GrenadeTypeIncendiary {
		t.Fatalf("inferno.Type = %q, want inferno", inferno.Type)
	}
	if len(inferno.Flames) != 2 {
		t.Fatalf("len(inferno.Flames) = %d, want 2", len(inferno.Flames))
	}
}

func TestParseGrenadeTypeFlashbangAlias(t *testing.T) {
	t.Parallel()
	if got := parseGrenadeType("flashbang"); got != models.GrenadeTypeFlash {
		t.Fatalf("parseGrenadeType(flashbang) = %q, want flash", got)
	}
}

func loadGrenadeFixture(t *testing.T, path string) map[string]*rawModels.Grenade {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var payload struct {
		Grenades map[string]json.RawMessage `json:"grenades"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	raw := make(map[string]*rawModels.Grenade, len(payload.Grenades))
	for key, value := range payload.Grenades {
		if string(value) == "false" {
			continue
		}
		var grenade rawModels.Grenade
		if err := json.Unmarshal(value, &grenade); err != nil {
			t.Fatalf("unmarshal grenade %s: %v", key, err)
		}
		raw[key] = &grenade
	}
	return raw
}
