package cs2gsi

import "regexp"

// CS2 GSI (especially demos) sometimes emits large numeric IDs without JSON
// quotes. Decoding those into string fields fails unless quoted first.
var (
	gsiOwnerNumRe      = regexp.MustCompile(`"owner"\s*:\s*(\d{10,})`)
	gsiSteamIDNumRe    = regexp.MustCompile(`"steamid"\s*:\s*(\d{10,})`)
	gsiSpectargetNumRe = regexp.MustCompile(`"spectarget"\s*:\s*(\d{10,})`)
	gsiBombPlayerNumRe = regexp.MustCompile(`"player"\s*:\s*(\d{10,})`)
)

func NormalizeGSIPayload(raw []byte) []byte {
	if len(raw) == 0 {
		return raw
	}
	s := string(raw)
	s = gsiOwnerNumRe.ReplaceAllString(s, `"owner": "$1"`)
	s = gsiSteamIDNumRe.ReplaceAllString(s, `"steamid": "$1"`)
	s = gsiSpectargetNumRe.ReplaceAllString(s, `"spectarget": "$1"`)
	s = gsiBombPlayerNumRe.ReplaceAllString(s, `"player": "$1"`)
	return []byte(s)
}
