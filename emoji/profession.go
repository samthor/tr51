package emoji

// ProfessionFor returns the profession suffix for a profession emoji, or zero if none.
func ProfessionFor(s string) rune {
	out := make([]rune, 0, 2)

	for i, r := range s {
		if i == 0 {
			if !IsPerson(r) {
				return 0
			}
			continue
		} else if i == 1 && IsSkinTone(r) {
			continue
		} else if r == runeVS16 {
			continue
		}
		out = append(out, r)
	}

	if len(out) != 2 || out[0] != runeZWJ {
		return 0
	}

	cand := out[1]
	if cand >= 0x1f466 && cand <= 0x1f469 {
		return 0 // this is a family
	}
	return cand
}
