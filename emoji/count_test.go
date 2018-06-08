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
	}
	for _, td := range data {
		if actual := Count(td.in); actual != td.count {
			t.Errorf("for %s, expected %d was %v", td.in, td.count, actual)
		}
	}
}
