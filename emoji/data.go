package emoji

import (
	"io"

	"github.com/samthor/tr51"
)

type emojiData struct {
	unqualified  bool    // whether this needs VS16
	modifierBase bool    // whether this can be modified
	version      float32 // unicode version from
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
				v.version = l.Version
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

// StripOpts controls what Normalize will strip.
type StripOpts struct {
	Tone   bool
	Gender bool
}

var (
	stripAll = StripOpts{Tone: true, Gender: true}
)

// Strip returns all the emoji parts of the passed string, removing tone and gender.
func (ed *Data) Strip(raw string) string {
	return ed.Normalize(raw, stripAll)
}

// Normalize returns only the emoji parts of the passed string.
func (ed *Data) Normalize(raw string, opts StripOpts) string {
	pending := []rune{0}

	// #0: Special-case single rune tone modifiers, which appear in test data.
	var singleTone bool
	for i, r := range raw {
		if i == 0 && IsSkinTone(r) {
			singleTone = true
		} else {
			singleTone = false
			break
		}
	}
	if singleTone {
		return raw
	}

	// #1: Remove VS16 and other modifiers.
	for _, r := range raw {
		if r == runeVS16 {
			// remove VS16
			continue
		} else if IsSkinTone(r) {
			if opts.Tone {
				// strip without checking
				continue
			}
			l := len(pending)
			if d, ok := ed.emoji[pending[l-1]]; ok && d.modifierBase {
				// great, skin tone is valid here
				pending = append(pending, r)
			}
			continue
		} else if IsGender(r) && opts.Gender {
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
	pending = append(pending, 0)

	// #2: Iterate chars, removing non-emoji.
	lp := len(pending) - 1
	out := make([]rune, 0, lp)
	var pendingZWJ int
	var allowZWJ int
	for i := 1; i < lp; i++ {
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

		if IsSkinTone(r) {
			// skin tone counts as a VS16, so look for a previous tone
			allowZWJ = i + 1
			l := len(out)
			if out[l-1] == runeVS16 {
				out[l-1] = r
				continue
			}
			out = append(out, r)
			continue
		}

		if IsFlagPart(r) {
			// just allow
			// TODO(samthor): Are these part of the data? Do we need this branch?
			out = append(out, r)
			continue
		}

		if d, ok := ed.emoji[r]; ok {
			if pendingZWJ == i {
				out = append(out, runeZWJ)
			}

			out = append(out, r)
			if d.unqualified {
				if IsSkinTone(pending[i+1]) {
					// do nothing as this acts as a VS16
					continue
				}
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
