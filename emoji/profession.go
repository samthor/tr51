package emoji

// ProfessionFor returns the profession suffix for a profession emoji, or the empty string if none.
func ProfessionFor(s string) string {
	cut := -1

	for i, r := range s {
		if i == 0 {
			if !IsPerson(r) {
				return ""
			}
		} else if r >= 0x1f466 && r <= 0x1f469 {
			// this is a family
			return ""
		} else if r == runeZWJ {
			if cut != -1 {
				return "" // only one ZWJ expected
			}
			cut = 0
		} else if cut == 0 {
			cut = i // cut is "at next position"
		}
	}

	if cut <= 0 {
		return ""
	}

	return s[cut:]
}
