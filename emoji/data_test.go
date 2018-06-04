package emoji

import (
	"bytes"
	"testing"

	"github.com/samthor/tr51"
)

func TestData(t *testing.T) {

	raw := `
1F3F4         ; Emoji_Presentation   #  7.0  [1] (🏴)       black flag
1F93C..1F93E  ; Emoji                #  9.0  [3] (🤼..🤾)    people wrestling..person playing handball
1F93C..1F93E  ; Emoji_Presentation   #  9.0  [3] (🤼..🤾)    people wrestling..person playing handball
267E..267F    ; Emoji                #  4.1  [2] (♾️..♿)    infinity..wheelchair symbol
1F3F7..1F4FD  ; Emoji                # [263] (🏷..📽)    LABEL..FILM PROJECTOR
1F680..1F6C5  ; Emoji_Presentation   #  [70] (🚀..🛅)    ROCKET..LEFT LUGGAGE
1F442..1F4FC  ; Emoji_Presentation   # [187] (👂..📼)    EAR..VIDEOCASSETTE
`

	r := tr51.NewReader(bytes.NewBuffer([]byte(raw)))
	ed, err := NewData(r)
	if err != nil {
		t.Fatalf("couldn't BuildEmojiData: %v", err)
	}

	type testData struct {
		in, out string
	}

	data := []testData{
		{"foo", ""},
		{"♾", "♾️"},
		{string([]rune{0x1f1fa, 0xfe0f, 0x1f1f8}), "🇺🇸"},
		{string([]rune{0x1f93e, 0x1f3fe, 0x200d, 0x2642, 0xfe0f}), "🤾"},
		{string([]rune{0x1f3f4, 0xe007f}), string([]rune{0x1f3f4})},
		{string([]rune{0x1f3f4, 0xe0061, 0xe007f}), string([]rune{0x1f3f4, 0xe0061, 0xe007f})},
		{"👨🏼‍🚒", "👨‍🚒"},
	}
	for _, td := range data {
		actual := ed.Normalize(td.in)
		if actual != td.out {
			t.Errorf("for %s, expected %s was `%s` (%v)", td.in, td.out, actual, []rune(actual))
		}
	}
}
