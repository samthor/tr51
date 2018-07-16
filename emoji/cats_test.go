package emoji

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/samthor/tr51"
)

func TestCats(t *testing.T) {

	raw := `
1F3F4         ; Emoji_Presentation   #  7.0  [1] (🏴)       black flag
#CategoryA
2640          ; Emoji                #  1.1  [1] (♀️)       female sign
2642          ; Emoji                #  1.1  [1] (♂️)       male sign
1F3F4         ; Emoji_Presentation   #  7.0  [1] (🏴)       black flag
1F93C
#CategoryB
1F3F4         ; Emoji_Presentation   #  7.0  [1] (🏴)       black flag
1F93C..1F93E  ; Emoji                #  9.0  [3] (🤼..🤾)    people wrestling..person playing handball
`

	r := tr51.NewReader(bytes.NewBuffer([]byte(raw)))
	cats, err := NewCats(r)
	if err != nil {
		t.Fatalf("couldn't NewCats: %v", err)
	}

	if expected, actual := []string{"CategoryA", "CategoryB"}, cats.Titles(); !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected categories: %+v, was %+v", expected, actual)
	}

	type testData struct {
		seq  string
		cats []string
	}

	data := []testData{
		{"🏴", []string{"", "CategoryA", "CategoryB"}},
		{"♀️", []string{"CategoryA"}},
		{"♂️", []string{"CategoryA"}},
		{"🤼", []string{"CategoryA", "CategoryB"}},
		{"🤽", []string{"CategoryB"}},
		{"🤾", []string{"CategoryB"}},
	}

	for _, td := range data {
		actual := cats.Get(td.seq)
		if !reflect.DeepEqual(actual, td.cats) {
			t.Errorf("for %s, expected (%v) actual (%v)", td.seq, td.cats, actual)
		}
	}
}
