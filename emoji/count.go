package emoji

// Count returns a rough count of the passed emoji string, assumed to be normalized.
func Count(raw string) int {
	var count int
	for _, r := range raw {
		if r == runeZWJ || r == runeVS16 || r == runeCap ||
			IsTag(r) || IsTagCancel(r) || IsSkinTone(r) || IsGender(r) {
			// ignore control characters and gender
			continue
		}
		if IsFlagPart(r) {
			count += 1
		} else {
			count += 2
		}
	}

	// round up but return minimum count if string had content
	out := (count + 1) / 2 // round up
	if out <= 1 && len(raw) > 0 {
		return 1
	}
	return out
}
