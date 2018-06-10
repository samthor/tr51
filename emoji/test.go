package emoji

import (
	"io"
	"strings"

	"github.com/samthor/tr51"
)

type emojiTest struct {
	qualified string
	notes     string
	subgroup  string
}

type subgroupTest struct {
	group string
	emoji []string
}

// Test wraps parsed data from emoji-test.txt.
type Test struct {
	emoji map[string]emojiTest

	groups        map[string][]string      // groups to subgroups
	subgroups     map[string]*subgroupTest // subgroups to emoji
	subgroupOrder []string
}

func (t *Test) subgroup(name, group string) *subgroupTest {
	out := t.subgroups[name]
	if out == nil {
		out = &subgroupTest{group: group}
		t.subgroups[name] = out
		t.subgroupOrder = append(t.subgroupOrder, name)
	}
	return out
}

// NewTest returns a new Test struct, which helps match complex emoji parts. Expects emoji-test.txt
// from Emoji 4.0+.
func NewTest(r *tr51.Reader) (*Test, error) {
	t := &Test{
		emoji:     make(map[string]emojiTest),
		groups:    make(map[string][]string),
		subgroups: make(map[string]*subgroupTest),
	}

	var group, subgroup string
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
				group = parts[1]
				subgroup = ""
			case "subgroup":
				subgroup = parts[1]
				t.groups[group] = append(t.groups[group], subgroup)
			}
			continue
		}

		// non-fully-qualified normally succeeds fully-qualified, but old versions don't always have it
		seq := l.AsSequence()
		qualified := string(seq)
		s := tr51.Unqualify(qualified)
		if _, ok := t.emoji[s]; ok {
			continue
		}
		st := t.subgroup(subgroup, group)
		st.emoji = append(st.emoji, s)

		test := emojiTest{notes: l.Notes, subgroup: subgroup}
		if s != qualified {
			// safe if not the same
			test.qualified = qualified
		}
		t.emoji[s] = test
	}

	return t, nil
}

// TestEach contains data about each emoji.
type TestEach struct {
	Emoji    string
	Notes    string
	Group    string
	Subgroup string
}

// Each enumerates through all found emoji in Test.
func (t *Test) TestEach(fn func(*TestEach)) {
	var each TestEach

	// nb. we order by subgroupOrder to provide consistent ordering through file
	for _, subgroup := range t.subgroupOrder {
		st := t.subgroups[subgroup]
		for _, emoji := range st.emoji {
			test := t.emoji[emoji]

			each.Emoji = test.qualified
			if each.Emoji == "" {
				each.Emoji = emoji
			}

			each.Notes = test.notes

			st := t.subgroups[test.subgroup]
			each.Group = st.group
			each.Subgroup = test.subgroup

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
