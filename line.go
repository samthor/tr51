package tr51

import (
	"bytes"
	"errors"
	"strconv"
	"unicode"
)

var (
	propertiesSep   = []byte{';'}
	rangeSep        = []byte("..")
	keycapHash      = []byte("keycap: #")
	keycapEscaped   = []byte(`\x{23}`)
	escapedPrefix   = []byte(`\x`)
	versionNAPrefix = []byte("NA ")

	// ErrInvalidRange indicates bad TR51 data.
	ErrInvalidRange = errors.New("invalid emoji range")

	// ErrUnhandledEscape indicates that an unhandled escape was reached.
	// This library special-cases escapes as there's not many of them.
	ErrUnhandledEscape = errors.New("unhandled \\x escape code")
)

// Line represents all possible raw line parts of a TR51 doc.
type Line struct {
	Single     rune     // for single rune emoji
	Low, High  rune     // for low/high pairs e.g., AAAA..BBBB
	Sequence   []rune   // for runs of emoji e.g., 1F468 200D 2764 FE0F
	Version    float32  // unicode version
	Notes      string   // trailing notes as part of comment
	Properties []string // ;-separated properties
}

// HasProperty returns whether this line has the given property.
func (lp *Line) HasProperty(s string) bool {
	for _, v := range lp.Properties {
		if s == v {
			return true
		}
	}
	return false
}

// AsSequence returns single or sequenced emoji as a sequence.
func (lp *Line) AsSequence() []rune {
	if lp.Single != 0 {
		return []rune{lp.Single}
	}
	return lp.Sequence
}

// AsString returns single or sequenced emoji as a string, or returns an empty stirng.
func (lp *Line) AsString() string {
	return string(lp.AsSequence())
}

// AsRange returns single or ranged emoji as a range.
func (lp *Line) AsRange() (low, high rune) {
	if lp.Single != 0 {
		return lp.Single, lp.Single
	}
	return lp.Low, lp.High
}

// Each calls the passed callback for every rune or sequence.
func (lp *Line) Each(fn func([]rune)) {
	if lp.Single != 0 {
		fn([]rune{lp.Single})
	}
	if len(lp.Sequence) != 0 {
		fn(lp.Sequence)
	}
	if lp.High != 0 {
		for i := lp.Low; i <= lp.High; i++ {
			fn([]rune{i})
		}
	}
}

// HasEmoji returns whether this line has any emoji parts.
func (lp *Line) HasEmoji() bool {
	return lp.Single != 0 || len(lp.Sequence) > 0 || lp.High != 0
}

// Parse parses a single line of a TR51 doc.
func Parse(line []byte) (out Line, err error) {
	var comment []byte
	commentIndex := bytes.IndexByte(line, '#')
	if commentIndex != -1 {
		// nb. special-case dealing with "keycap: #" in earlier emoji data (dup comment)
		if commentIndex >= 9 {
			cand := line[commentIndex-8 : commentIndex+1]
			if bytes.Equal(cand, keycapHash) {
				next := bytes.IndexByte(line[commentIndex+1:], '#')
				commentIndex = commentIndex + next + 1
			}
		}

		comment = bytes.TrimSpace(line[commentIndex+1:])
		line = bytes.TrimSpace(line[:commentIndex])
	}

	// special-case the escaped keycap
	// TODO: should support \x{anything}
	line = bytes.Replace(line, keycapEscaped, []byte{'#'}, -1)
	if bytes.Index(line, escapedPrefix) != -1 {
		return out, ErrUnhandledEscape
	}

	left := bytes.Split(line, propertiesSep)
	for i := range left {
		left[i] = bytes.TrimSpace(left[i])
	}
	if len(left) > 1 {
		out.Properties = make([]string, len(left)-1)
		for i, v := range left[1:] {
			out.Properties[i] = string(v)
		}
	}

	// extract points
	points := left[0]
	if pointParts := bytes.Fields(points); len(pointParts) == 1 {
		rangeParts := bytes.Split(pointParts[0], rangeSep)
		if len(rangeParts) == 2 {
			// range e.g. AAAA..BBBB
			out.Low = parsePoint(rangeParts[0])
			out.High = parsePoint(rangeParts[1])
		} else if len(rangeParts) > 2 {
			return out, ErrInvalidRange
		} else {
			// single point only
			out.Single = parsePoint(pointParts[0])
		}
	} else if len(pointParts) == 0 {
		// nothing here, nothing to do
		out.Notes = string(comment)
		return out, nil
	} else {
		// space-separated sequence
		out.Sequence = make([]rune, len(pointParts))
		for i, v := range pointParts {
			out.Sequence[i] = parsePoint(v)
		}
	}

	if len(comment) < 2 {
		return out, nil // weird, but possible
	}

	// strip V if it's part of version
	if comment[0] == 'V' && unicode.IsDigit(rune(comment[1])) {
		comment = comment[1:]
	}

	// special-case "NA", version found in final Emoji 11
	if len(comment) > 3 && bytes.Equal(comment[:3], versionNAPrefix) {
		comment = comment[2:]
	}

	// find version
	cand := numberPrefixOf(comment)
	if len(cand) == 0 && comment[0] == '(' {
		cand = numberPrefixOf(comment[1:])
		if len(cand) != 0 && len(comment) > len(cand)+2 && comment[len(cand)+1] == ')' {
			// ok
		} else {
			cand = nil
		}
	}
	if len(cand) >= 2 { // look for "3.0", not "3"
		v64, err := strconv.ParseFloat(string(cand), 32)
		if err != nil {
			return out, err
		}
		out.Version = float32(v64)
	}

	// work out whether we have a ( before any real chars
	// this allows us to have notes like "AB button (blood type)"
	var hasEarlyBracket bool
	for _, v := range comment {
		if v == '(' {
			hasEarlyBracket = true
			break
		}
		if isASCIILetter(rune(v)) {
			break
		}
	}

	// move to notes
	notesFromIndex := bytes.IndexByte(comment, ')')
	if !hasEarlyBracket || notesFromIndex == -1 {
		// only happens in -test.txt, all other data has ()'s
		spaceIndex := bytes.IndexByte(comment, ' ')
		if spaceIndex != -1 {
			notesFromIndex = spaceIndex
		} else {
			notesFromIndex = bytes.IndexFunc(comment, func(r rune) bool {
				return !(r == 32 || r > 127) // skip over non-ascii and spaces
			})
		}
	} else {
		notesFromIndex++
	}
	notes := bytes.TrimSpace(comment[notesFromIndex:])
	out.Notes = string(notes)

	return out, nil
}

func isASCIILetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// numberPrefixOf returns the byte prefix which contains numeric or dot bytes.
func numberPrefixOf(src []byte) []byte {
	pastNumber := bytes.IndexFunc(src, func(r rune) bool {
		return !unicode.IsNumber(r) && r != '.'
	})
	if pastNumber != -1 {
		return src[:pastNumber]
	}
	return nil
}

// parsePoint parses a hex Unicode code point, returning zero if invalid.
func parsePoint(b []byte) rune {
	point, err := strconv.ParseUint(string(b), 16, 32)
	if err != nil {
		return 0
	}
	return rune(point)
}
