package emoji

const (
	runeZWJ       = 0x200d
	runeCap       = 0x20e3
	runeVS16      = 0xfe0f
	runeTagCancel = 0xe007f
)

// IsFlagPart returns whether the passed rune is part of a flag.
func IsFlagPart(r rune) bool {
	return r >= 0x1f1e6 && r <= 0x1f1ff
}

// IsSkinTone returns whether the passed rune is a Fitzpatrick skin tone modifier.
func IsSkinTone(r rune) bool {
	return r >= 0x1f3fb && r <= 0x1f3ff
}

// IsBeforeCap returns whether the passed rune can appear before a keycap.
func IsBeforeCap(r rune) bool {
	return r == '#' || r == '*' || (r >= '0' && r <= '9')
}

// IsTag returns whether the passed rune is a tag character, for tag sequences.
func IsTag(r rune) bool {
	return r >= 0xe0020 && r < 0xe007f
}

// IsTagCancel returns whether the passed rune ends a tag sequence.
func IsTagCancel(r rune) bool {
	return r == runeTagCancel
}

// IsTagBase returns whether the passed rune can have tags following it.
func IsTagBase(r rune) bool {
	return r == 0x1f3f4
}
