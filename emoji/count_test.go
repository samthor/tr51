package emoji

import (
	"testing"
)

func TestCount(t *testing.T) {
	type testData struct {
		in    string
		count int
	}

	data := []testData{
		{"🇦🇺", 1},
		{"⛹🏿‍♂️", 1},
		{"⛹🏿‍♂️9⃣", 2},
		{"🔟🏴󠁧󠁢󠁳󠁣󠁴󠁿🇫🇯", 3},
		{"👩‍⚕️", 1}, // has ZWJ
		{"👩‍❤️‍💋‍👩👨‍❤️‍💋‍👨", 2},
		{"👩‍👩‍👧‍👧", 1},
		{"🇺🇳", 1},
		{"🇺🇳🇳", 2}, // technically 1.5
		{"🏴‍☠️", 1},
		{"🏳️‍🌈", 1},
	}
	for _, td := range data {
		if actual := Count(td.in); actual != td.count {
			t.Errorf("for %s, expected %d was %v", td.in, td.count, actual)
		}
	}
}
