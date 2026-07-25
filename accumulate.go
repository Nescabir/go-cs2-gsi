package cs2gsi

import (
	"encoding/json"
	"sync"
)

type payloadAcc struct {
	mu      sync.Mutex
	keys    map[string]json.RawMessage
	mapName string
}

func newPayloadAcc() *payloadAcc {
	return &payloadAcc{keys: make(map[string]json.RawMessage)}
}

func (a *payloadAcc) Merge(raw []byte) ([]byte, error) {
	var incoming map[string]json.RawMessage
	if err := json.Unmarshal(raw, &incoming); err != nil {
		return raw, err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if name := mapNameFromPayload(incoming); name != "" && name != a.mapName {
		a.keys = make(map[string]json.RawMessage)
		a.mapName = name
	}

	for key, value := range incoming {
		switch key {
		case "previously", "added":
			continue
		}
		if isJSONFalse(value) {
			delete(a.keys, key)
			continue
		}
		a.keys[key] = value
	}

	ensureRequiredTopLevel(a.keys)

	out := make(map[string]json.RawMessage, len(a.keys)+2)
	for key, value := range a.keys {
		out[key] = value
	}
	if value, ok := incoming["previously"]; ok {
		out["previously"] = value
	}
	if value, ok := incoming["added"]; ok {
		out["added"] = value
	}

	return json.Marshal(out)
}

func ensureRequiredTopLevel(keys map[string]json.RawMessage) {
	defaults := map[string]string{
		"auth":             `{"token":""}`,
		"provider":         `{"name":"CS2","appid":730,"version":0,"steamid":"","timestamp":0}`,
		"round":            `{"phase":"live","bomb":"","win_team":""}`,
		"bomb":             `{"state":"","countdown":"0","position":"0, 0, 0","player":""}`,
		"grenades":         `{}`,
		"phase_countdowns": `{"phase":"live","phase_ends_in":"0.0"}`,
	}
	for key, fallback := range defaults {
		if _, ok := keys[key]; !ok {
			keys[key] = json.RawMessage(fallback)
		}
	}
}

func isJSONFalse(raw json.RawMessage) bool {
	return string(raw) == "false"
}

func mapNameFromPayload(payload map[string]json.RawMessage) string {
	raw, ok := payload["map"]
	if !ok {
		return ""
	}
	var meta struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &meta); err != nil {
		return ""
	}
	return meta.Name
}
