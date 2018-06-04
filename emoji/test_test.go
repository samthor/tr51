package emoji

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/samthor/tr51"
)

func TestTest(t *testing.T) {
	raw := `
0039 FE0F 20E3                             ; fully-qualified     # 9️⃣ keycap: 9
0039 20E3                                  ; non-fully-qualified # 9⃣ keycap: 9
1F51F                                      ; fully-qualified     # 🔟 keycap: 10
26F9 1F3FF 200D 2642 FE0F                  ; fully-qualified     # ⛹🏿‍♂️ man bouncing ball: dark skin tone
26F9 1F3FF 200D 2642                       ; non-fully-qualified # ⛹🏿‍♂ man bouncing ball: dark skin tone
1F1E6 1F1FA                                ; fully-qualified     # 🇦🇺 Australia
1F1EB 1F1EF                                ; fully-qualified     # 🇫🇯 Fiji
1F1EB 1F1F0                                ; fully-qualified     # 🇫🇰 Falkland Islands
1F1EB 1F1F2                                ; fully-qualified     # 🇫🇲 Micronesia
1F1EB 1F1F4                                ; fully-qualified     # 🇫🇴 Faroe Islands
1F1FA 1F1F2                                ; fully-qualified     # 🇺🇲 U.S. Outlying Islands
1F3F4 E0067 E0062 E0073 E0063 E0074 E007F  ; fully-qualified     # 🏴󠁧󠁢󠁳󠁣󠁴󠁿 Scotland
`

	r := tr51.NewReader(bytes.NewBuffer([]byte(raw)))
	et, err := NewTest(r)
	if err != nil {
		t.Fatalf("couldn't NewTest: %v", err)
	}

	type testData struct {
		in  string
		out []string
	}

	data := []testData{
		{string([]rune{0x1f1e6, 0x1f1fa, 0x1f1f2}), []string{"🇦🇺", string([]rune{0x1f1f2})}},
		{"⛹🏿‍♂️", []string{"⛹🏿‍♂️"}},
		{"⛹🏿‍♂️9⃣", []string{"⛹🏿‍♂️", "9⃣"}},
		{"🔟🏴󠁧󠁢󠁳󠁣󠁴󠁿🇫🇯", []string{"🔟", "🏴󠁧󠁢󠁳󠁣󠁴󠁿", "🇫🇯"}},
	}
	for _, td := range data {
		actual := et.Split(td.in)

		if !reflect.DeepEqual(actual, td.out) {
			t.Errorf("for %s, expected %+v (%d) was %+v (%d)", td.in, td.out, len(td.out), actual, len(actual))
		}
	}

	type testName struct {
		emoji string
		name  string
	}
	names := []testName{
		{"⛹🏿‍♂", "man bouncing ball: dark skin tone"},
		{"🏴󠁧󠁢󠁳󠁣󠁴󠁿", "Scotland"},
		{"🏴󠁧󠁢󠁳󠁴󠁿", ""},
		{"9️⃣", "keycap: 9"},
		{"9⃣", "keycap: 9"},
	}
	for _, n := range names {
		actual := et.Name(n.emoji)

		if !reflect.DeepEqual(actual, n.name) {
			t.Errorf("for %s, expected %v was %v", n.emoji, n.name, actual)
		}
	}
}
