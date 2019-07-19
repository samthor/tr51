package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/samthor/tr51"
)

func main() {
	type emojiPart struct {
		name         string
		version      float32
		modifierBase bool
		presentation bool
		profession   bool
		role         bool
		keycap       bool
	}
	type flagKey struct {
		l, r rune
	}
	emojiParts := make(map[rune]emojiPart)
	var emojiFlags []flagKey
	var emojiOthers []string

	// process single emoji data
	dataReader := readTR51("emoji-data.txt")
	for {
		l, err := dataReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("could not read: %v", err)
		}

		isEmoji := l.HasProperty("Emoji")
		isPresentation := l.HasProperty("Emoji_Presentation")
		isModifierBase := l.HasProperty("Emoji_Modifier_Base")

		if !(isEmoji || isPresentation || isModifierBase) {
			continue
		}

		low, high := l.AsRange()
		for r := low; r <= high; r++ {
			ep := emojiParts[r]

			if ep.version != l.Version {
				if !(ep.version == 0.0 && isEmoji) {
					// got inconsistent version
					log.Printf("%c: prop=%+v version=%v was=%v", r, l.Properties, l.Version, ep)
				}
			}
			ep.version = l.Version
			ep.presentation = ep.presentation || isPresentation
			ep.modifierBase = ep.modifierBase || isModifierBase

			emojiParts[r] = ep
		}
	}

	// process test data and find combinations
	testReader := readTR51("emoji-test.txt")
outer:
	for {
		l, err := testReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("could not read: %v", err)
		}
		if l.HasProperty("component") && l.Single != 0 {
			// ok, we'll just name it
		} else if !l.HasProperty("fully-qualified") {
			continue
		}

		// ... yet unqualify it
		raw := []rune(tr51.Unqualify(l.AsString()))
		if len(raw) == 0 {
			log.Fatalf("unqualified emoji is empty: %v", l.AsString())
		}

		r := raw[0]
		var name string
		var isProfession bool
		var isRole bool
		var isKeycap bool

		// 0) insert names
		if len(raw) == 1 {
			name = l.Notes
			goto update
		}

		// ... skip any with skin tone runes
		for _, r := range raw {
			if isSkinTone(r) {
				continue outer
			}
		}

		// 1) look for professions ("man firefighter", "woman firefighter")
		if len(raw) == 3 && isGenderPerson(raw[0]) && raw[1] == 0x200d {
			r = raw[2]
			isProfession = true
			goto update
		}

		// 2) look for roles ("detective", "man detective", "woman detective")
		if len(raw) == 3 && raw[1] == 0x200d && isGender(raw[2]) {
			isRole = true
			goto update
		}

		// 3) look for keycaps
		if len(raw) == 2 && raw[1] == 0x20e3 {
			isKeycap = true
			goto update
		}

		// 4) look for flags
		if len(raw) == 2 && isFlagPart(raw[0]) && isFlagPart(raw[1]) {
			key := flagKey{
				'a' + (raw[0] - 0x1f1e6),
				'a' + (raw[1] - 0x1f1e6),
			}
			emojiFlags = append(emojiFlags, key)
		} else if len(raw) > 1 {
			emojiOthers = append(emojiOthers, string(raw))
		}
		continue

	update:
		ep := emojiParts[r]
		if name != "" {
			ep.name = name
		}
		ep.profession = ep.profession || isProfession
		ep.role = ep.role || isRole
		ep.keycap = ep.keycap || isKeycap
		emojiParts[r] = ep
	}

	count := func(pred func(emojiPart) bool) (out int) {
		for r := range emojiParts {
			if pred(emojiParts[r]) {
				out++
			}
		}
		return out
	}

	log.Printf("professions: %d", count(func(ep emojiPart) bool { return ep.profession }))
	log.Printf("roles: %d", count(func(ep emojiPart) bool { return ep.role }))
	log.Printf("keycaps: %d", count(func(ep emojiPart) bool { return ep.keycap }))

	emojiPartAll := make(runeSlice, 0, len(emojiParts))
	for r := range emojiParts {
		emojiPartAll = append(emojiPartAll, r)
	}
	emojiPartAll.Sort()

	var output struct {
		professions []rune
		roles       []rune
		variation   []rune
		flags       []rune
	}

	for _, r := range emojiPartAll {
		ep := emojiParts[r]
		if ep.profession {
			output.professions = append(output.professions, r)
		}
		if ep.role {
			output.roles = append(output.roles, r)
		}
		if !ep.presentation {
			output.variation = append(output.variation, r)
		}
	}
	for _, flag := range emojiFlags {
		output.flags = append(output.flags, flag.l, flag.r)
	}

	// TODO(samthor): we need emoji-zwj-sequences.txt for coverage of unicode versions

	fmt.Printf(`// Generated on %v
export const professions = "%s";
export const roles = "%s";
export const variation = "%s";
export const flags = "%s";
`, time.Now(), string(output.professions), string(output.roles), string(output.variation), string(output.flags))
}

func isFlagPart(r rune) bool {
	return r >= 0x1f1e6 && r <= 0x1f1ff
}

func isGenderPerson(r rune) bool {
	return r == 0x1f468 || r == 0x1f469 // woman, man
}

func isGender(r rune) bool {
	return r == 0x2640 || r == 0x2642 // woman, man
}

func isSkinTone(r rune) bool {
	return r >= 0x1f3fb && r <= 0x1f3ff
}

func readTR51(filename string) *tr51.Reader {
	all, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("could not read %v: %v", filename, err)
	}

	b := bytes.NewBuffer(all)
	return tr51.NewReader(b)
}
