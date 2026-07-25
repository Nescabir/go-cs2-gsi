package cs2gsi

import (
	"strings"
	"testing"
)

func TestNormalizeGSIPayloadNumericOwner(t *testing.T) {
	in := []byte(`{"grenades":{"1":{"owner":76561198047876970,"type":"smoke"}}}`)
	out := string(NormalizeGSIPayload(in))
	if !strings.Contains(out, `"owner": "76561198047876970"`) {
		t.Fatalf("normalize = %s", out)
	}
}

func TestPayloadAccMergeFillsMissingSections(t *testing.T) {
	acc := newPayloadAcc()

	first := []byte(`{
		"auth":{"token":"t"},
		"map":{"name":"de_dust2","phase":"live","round":1,"team_ct":{"score":0},"team_t":{"score":0}},
		"round":{"phase":"live"},
		"phase_countdowns":{"phase":"live","phase_ends_in":"10.0"},
		"player":{"steamid":"1","name":"A","team":"CT"},
		"allplayers":{"1":{"steamid":"1","name":"A","team":"CT","observer_slot":0}},
		"bomb":{"state":"carried","countdown":"0","position":"0, 0, 0","player":"1"},
		"grenades":{}
	}`)
	merged1, err := acc.Merge(NormalizeGSIPayload(first))
	if err != nil {
		t.Fatal(err)
	}

	delta := []byte(`{
		"map":{"name":"de_dust2","phase":"live","round":1,"team_ct":{"score":1},"team_t":{"score":0}},
		"round":{"phase":"live"},
		"player":{"steamid":"1","name":"A","team":"CT"},
		"allplayers":{"1":{"steamid":"1","name":"A","team":"CT","observer_slot":0}}
	}`)
	merged2, err := acc.Merge(NormalizeGSIPayload(delta))
	if err != nil {
		t.Fatal(err)
	}

	gsi := New(NewConfig())
	if err := gsi.Digest(merged1); err != nil {
		t.Fatalf("digest first: %v", err)
	}
	if err := gsi.Digest(merged2); err != nil {
		t.Fatalf("digest delta: %v", err)
	}
	if snap := gsi.Snapshot(); snap == nil || snap.Map == nil {
		t.Fatal("expected snapshot after delta")
	}
}

func TestDigestDeltaTickWithoutBomb(t *testing.T) {
	gsi := New(NewConfig())

	first := []byte(`{
		"provider":{"name":"CS2","appid":730,"version":0,"steamid":"","timestamp":0},
		"auth":{"token":""},
		"map":{"name":"de_mirage","phase":"live","round":3,"team_ct":{"score":1,"name":"CT","consecutive_round_losses":0,"timeouts_remaining":3,"matches_won_this_series":0},"team_t":{"score":2,"name":"T","consecutive_round_losses":0,"timeouts_remaining":3,"matches_won_this_series":0}},
		"round":{"phase":"live","bomb":"","win_team":""},
		"phase_countdowns":{"phase":"live","phase_ends_in":"45.0"},
		"player":{"steamid":"76561198000000001","name":"P1","team":"CT","observer_slot":0,"state":{"health":100,"armor":0,"helmet":false,"flashed":0,"smoked":0,"burning":0,"money":0,"round_kills":0,"round_killhs":0,"round_totaldmg":0,"equip_value":0},"match_stats":{"kills":0,"assists":0,"deaths":0,"mvps":0,"score":0},"weapons":{},"position":"0, 0, 0","forward":"0, 0, 0"},
		"allplayers":{"76561198000000001":{"steamid":"76561198000000001","name":"P1","team":"CT","observer_slot":0,"state":{"health":100,"armor":0,"helmet":false,"flashed":0,"smoked":0,"burning":0,"money":0,"round_kills":0,"round_killhs":0,"round_totaldmg":0,"equip_value":0},"match_stats":{"kills":0,"assists":0,"deaths":0,"mvps":0,"score":0},"weapons":{},"position":"0, 0, 0","forward":"0, 0, 0"}},
		"bomb":{"state":"carried","countdown":"0","position":"0, 0, 0","player":"76561198000000001"},
		"grenades":{}
	}`)
	if err := gsi.Digest(first); err != nil {
		t.Fatalf("first digest: %v", err)
	}

	delta := []byte(`{
		"map":{"name":"de_mirage","phase":"live","round":3,"team_ct":{"score":1,"name":"CT","consecutive_round_losses":0,"timeouts_remaining":3,"matches_won_this_series":0},"team_t":{"score":2,"name":"T","consecutive_round_losses":0,"timeouts_remaining":3,"matches_won_this_series":0}},
		"round":{"phase":"live","bomb":"","win_team":""},
		"player":{"steamid":"76561198000000001","name":"P1","team":"CT","observer_slot":0,"state":{"health":80,"armor":0,"helmet":false,"flashed":0,"smoked":0,"burning":0,"money":0,"round_kills":0,"round_killhs":0,"round_totaldmg":0,"equip_value":0},"match_stats":{"kills":0,"assists":0,"deaths":0,"mvps":0,"score":0},"weapons":{},"position":"0, 0, 0","forward":"0, 0, 0"},
		"allplayers":{"76561198000000001":{"steamid":"76561198000000001","name":"P1","team":"CT","observer_slot":0,"state":{"health":80,"armor":0,"helmet":false,"flashed":0,"smoked":0,"burning":0,"money":0,"round_kills":0,"round_killhs":0,"round_totaldmg":0,"equip_value":0},"match_stats":{"kills":0,"assists":0,"deaths":0,"mvps":0,"score":0},"weapons":{},"position":"0, 0, 0","forward":"0, 0, 0"}}
	}`)
	if err := gsi.Digest(delta); err != nil {
		t.Fatalf("delta digest: %v", err)
	}
	if snap := gsi.Snapshot(); snap.Bomb == nil {
		t.Fatal("expected bomb carried over from accumulated payload")
	}
}
