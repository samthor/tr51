package tr51

import (
	"strings"
)

const (
	// VS16 is the "Variation Selector 16" character, commonly known as "emoji mode".
	VS16 = "\ufe0f"
)

// Origins converts a space-separated list of origins from the Emoji 1.0 data, to long-form text
// useful for types.
func Origins(o string) []string {
	raw := strings.Fields(o)
	var out []string
	for _, origin := range raw {
		switch origin {
		case "z":
			out = append(out, "Origin_ZDings")
		case "a":
			out = append(out, "Origin_ARIB")
		case "j":
			out = append(out, "Origin_JCarrier")
		case "w":
			out = append(out, "Origin_WDings")
		case "x":
			out = append(out, "Origin_Other")
		default:
			// unknown
		}
	}
	return out
}

// Unqualify removes all VS16 charaters from the input string.
func Unqualify(o string) string {
	return strings.Replace(o, VS16, "", -1)
}
