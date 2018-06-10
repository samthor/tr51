package emoji

// FlagFor returns the country or region code this single flag is for, if any.
func FlagFor(s string) string {
	runes := []rune(s)
	l := len(runes)
	if l < 2 {
		return ""
	}

	if l == 2 && IsFlagPart(runes[0]) && IsFlagPart(runes[1]) {
		// is a 2-char flag
		first := 'a' + (runes[0] - 0x1f1e6)
		second := 'a' + (runes[1] - 0x1f1e6)
		return string([]rune{first, second})
	}

	if runes[0] != runeBlackFlag {
		return ""
	}
	return readTagSequence(runes[1:])
}

// readTagSequence returns the string version of the given tag sequence.
func readTagSequence(raw []rune) string {
	out := make([]rune, 0, len(raw)-1)
	var last bool
	for _, r := range raw {
		if last {
			// hit last without being last
			return ""
		} else if IsTagCancel(r) {
			// should be last rune
			last = true
		} else if IsTag(r) {
			// found a tag!
			out = append(out, ' '+(r-runeTagSpace))
		} else {
			// not a tag, fail early
			return ""
		}
	}

	if !last {
		return ""
	}
	return string(out)
}

// ReadTag returns the string version of the given tag sequence.
func ReadTag(raw string) string {
	return readTagSequence([]rune(raw))
}
