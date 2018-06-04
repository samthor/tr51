package emoji

import (
	"io"

	"github.com/samthor/tr51"
)

// Test wraps parsed data from emoji-test.txt.
type Test struct {
	emoji map[string]string
}

// NewTest returns a new Test struct, which helps match complex emoji parts. Expects emoji-test.txt
// from Emoji 4.0+.
func NewTest(r *tr51.Reader) (*Test, error) {
	m := make(map[string]string)
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !l.HasEmoji() {
			continue
		}

		// non-fully-qualified normally succeeds fully-qualified, but old versions don't always have it
		seq := l.AsSequence()
		s := tr51.Unqualify(string(seq))
		m[s] = l.Notes
	}

	return &Test{m}, nil
}

// Split splits an input string into component emoji. Assumes the input is well-formed.
func (t *Test) Split(s string) []string {
	var cands []string
	var curr []rune

	next := func() {
		if len(curr) > 0 {
			cands = append(cands, string(curr))
		}
		curr = nil
	}

	for _, c := range s {
		if len(curr) == 0 {
			curr = []rune{c}
			continue
		}
		last := curr[len(curr)-1]

		// flags
		if IsFlagPart(c) {
			if _, ok := t.emoji[string([]rune{last, c})]; ok {
				// we found a flag!
				curr = append(curr, c)
				next()
				continue
			}

			next()
			curr = []rune{c}
			continue
		} else if IsFlagPart(last) {
			next()
		}

		// allow modifiers and etc
		if c == runeZWJ || c == runeCap || c == runeVS16 ||
			IsTag(c) || IsTagCancel(c) || IsSkinTone(c) {
			curr = append(curr, c)
			continue
		}

		if last != runeZWJ {
			next()
		}
		curr = append(curr, c)
	}

	next()
	return cands
}

// Name returns the name for a single emoji. An empty name means there's no match.
func (t *Test) Name(s string) string {
	return t.emoji[tr51.Unqualify(s)]
}
