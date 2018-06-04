package tr51

import (
	"strings"
)

const (
	vs16 = "\ufe0f"
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
			// ignore
		}
	}
	return out
}

// Unqualify removes all 0xfe0f charaters from the input string.
func Unqualify(o string) string {
	return strings.Replace(o, vs16, "", -1)
}
