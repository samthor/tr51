package emoji

import (
	"testing"
)

func TestFlag(t *testing.T) {
	type testData struct {
		in  string
		out string
	}

	data := []testData{
		{"⛹🏿‍♂️", ""},
		{"🏴󠁧󠁢󠁳󠁣󠁴󠁿", "gbsct"}, // tag sequence flag
		{"🇫🇯", "fj"},         // normal flag
		{"🇫🇯🇫🇯", ""},         // single flag only
	}
	for _, td := range data {
		actual := FlagFor(td.in)

		if actual != td.out {
			t.Errorf("for %s, expected `%v` was `%v`", td.in, td.out, actual)
		}
	}
}
