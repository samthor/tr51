package emoji

import (
	"io"

	"github.com/samthor/tr51"
)

type emojiData struct {
	unqualified  bool // whether this needs VS16
	modifierBase bool // whether this can be modified
}

// Data wraps parsed data from emoji-data.txt.
type Data struct {
	emoji map[rune]emojiData
}

// NewData returns a new Data struct, which helps validate raw emoji parts. Expects emoji-data.txt
// from Emoji 2.0+.
func NewData(r *tr51.Reader) (*Data, error) {
	m := make(map[rune]emojiData)
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		isEmoji := l.HasProperty("Emoji")
		isPresentation := l.HasProperty("Emoji_Presentation")

		low, high := l.AsRange()
		if isEmoji || isPresentation {
			for r := low; r <= high; r++ {
				v := m[r]
				v.unqualified = isEmoji
				m[r] = v
			}
		}

		if l.HasProperty("Emoji_Modifier_Base") {
			for r := low; r <= high; r++ {
				v := m[r]
				v.modifierBase = true
				m[r] = v
			}
		}
	}

	return &Data{m}, nil
}

// Normalize returns only the emoji parts of the passed string.
func (ed *Data) Normalize(raw string) string {
	pending := []rune{0}

	// #1: Remove VS16 and other modifiers.
	for _, r := range raw {
		if r == runeVS16 {
			// remove VS16
			continue
		} else if IsSkinTone(r) {
			// TODO(samthor): Optionally retain this if Emoji_Modifier_Base.
			// remove skin tone
			continue
		} else if IsGender(r) {
			// remove gender modifiers
			l := len(pending)
			if pending[l-1] == runeZWJ {
				// ... and drop a previous ZWJ if we find one
				pending = pending[:l-1]
			}
			continue
		}
		pending = append(pending, r)
	}

	// #2: Iterate chars, removing non-emoji.
	out := make([]rune, 0, len(pending))
	var pendingZWJ int
	var allowZWJ int
	for i := 1; i < len(pending); i++ {
		r := pending[i]
		if r == runeZWJ {
			if allowZWJ == i {
				pendingZWJ = i + 1 // add it before valid rune at next index
			}
			continue
		}
		prev := pending[i-1]

		if r == runeCap {
			// allow if previous was number
			if IsBeforeCap(prev) {
				out = append(out, r)
				allowZWJ = i + 1
			}
			continue
		}

		if IsTag(r) {
			// allow if following a base or previous tag
			if IsTagBase(prev) || IsTag(prev) {
				out = append(out, r)
			}
			continue
		}

		if IsTagCancel(r) {
			// allow if following a tag
			if IsTag(prev) {
				out = append(out, r)
			}
			continue
		}

		if IsTag(prev) {
			// cancel the tag sequence if we got this far
			out = append(out, runeTagCancel)
		}

		if IsFlagPart(r) {
			// just allow if flag
			out = append(out, r)
			continue
		}

		if d, ok := ed.emoji[r]; ok {
			if pendingZWJ == i {
				out = append(out, runeZWJ)
			}

			out = append(out, r)
			if d.unqualified {
				// stick a VS16 on the end
				out = append(out, runeVS16)
			}
			allowZWJ = i + 1
			continue
		}
	}

	// #3: Profit!
	return string(out)
}
