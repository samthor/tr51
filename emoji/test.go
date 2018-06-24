package emoji

import (
	"io"
	"strings"

	"github.com/samthor/tr51"
)

type emojiTest struct {
	qualified string
	notes     string
}

type groupInfo struct {
	name  string
	emoji []string
}

// Test wraps parsed data from emoji-test.txt.
type Test struct {
	emoji  map[string]emojiTest
	groups []*groupInfo
}

// NewTest returns a new Test struct, which helps match complex emoji parts. Expects emoji-test.txt
// from Emoji 4.0+.
func NewTest(r *tr51.Reader) (*Test, error) {
	t := &Test{
		emoji: make(map[string]emojiTest),
	}

	var currentGroup *groupInfo
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !l.HasEmoji() {
			parts := strings.Split(l.Notes, ": ")
			if len(parts) != 2 {
				continue
			}
			switch parts[0] {
			case "group":
				currentGroup = &groupInfo{name: parts[1]}
				t.groups = append(t.groups, currentGroup)
			}
			continue
		}

		// non-fully-qualified normally succeeds fully-qualified, but old versions don't always have it
		seq := l.AsSequence()
		qualified := string(seq)
		unqualified := tr51.Unqualify(qualified)
		if _, ok := t.emoji[unqualified]; ok {
			continue
		}

		test := emojiTest{notes: l.Notes, qualified: qualified}
		t.emoji[unqualified] = test
		currentGroup.emoji = append(currentGroup.emoji, unqualified)
	}

	return t, nil
}

// TestEach contains data about each emoji.
type TestEach struct {
	Emoji string
	Notes string
	Group string
}

// Groups returns an array of the groups inside the TR51 data.
func (t *Test) Groups() []string {
	out := make([]string, 0, len(t.groups))
	for _, gi := range t.groups {
		out = append(out, gi.name)
	}
	return out
}

// Each enumerates through all found emoji in Test.
func (t *Test) TestEach(fn func(*TestEach)) {
	var each TestEach

	// nb. we use group to provide consistent ordering through file
	for _, gi := range t.groups {
		for _, emoji := range gi.emoji {
			test := t.emoji[emoji]

			each.Emoji = test.qualified
			each.Notes = test.notes
			each.Group = gi.name

			fn(&each)
		}
	}
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
	return t.emoji[tr51.Unqualify(s)].notes
}
