package main

func overrideControl() []rune {
	return []rune{
		0x1f46a, // 👪, incorrect modifierBase
		0x1f48f, // 💏, incorrect modifierBase
	}
}

func overrideEmojiPart(r rune, ep *emojiPart) {
	switch r {
	case 0x1f46a:
		fallthrough
	case 0x1f48f:
		ep.modifierBase = false
	}
}
