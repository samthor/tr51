package emoji

import (
	"testing"
)

func TestProfession(t *testing.T) {
	type testData struct {
		in  string
		out string
	}

	data := []testData{
		{"👩🏻‍💻", "💻"},
		{"🏴󠁧󠁢󠁳󠁣󠁴󠁿", ""},
		{"👮🏿", ""},
		{"👩🏿‍🚒", "🚒"},
		{"👨‍🚒", "🚒"},
	}
	for _, td := range data {
		actual := ProfessionFor(td.in)

		if actual != td.out {
			t.Errorf("for %s, expected `%v` was `%v`", td.in, td.out, actual)
		}
	}
}
