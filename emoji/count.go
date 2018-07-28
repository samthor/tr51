package emoji

// Count returns a rough count of the passed emoji string, assumed to be normalized. It assumes
// good data and returns ZWJ'ed characters as one.
func Count(raw string) int {
	var halfCount int
	for _, r := range raw {
		if r == runeZWJ {
			halfCount -= 2 // assume good intent / ZWJ joins
		} else if r == runeVS16 || r == runeCap || IsTag(r) || IsTagCancel(r) || IsSkinTone(r) {
			// ignore control characters that don't need ZWJ
		} else if IsFlagPart(r) {
			halfCount += 1 // flag half
		} else {
			halfCount += 2 // normal char
		}
	}

	// round up but return minimum count if string had content
	out := (halfCount + 1) / 2 // round up
	if out <= 1 && len(raw) > 0 {
		return 1
	}
	return out
}
