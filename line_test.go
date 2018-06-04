package tr51

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testdata := map[string]Line{
		// simple/comment cases
		"":             Line{},
		"# blah":       Line{Notes: "blah"},
		"1F61F # blah": Line{Single: 0x1f61f, Notes: "blah"},
		"# v1.0 blah":  Line{Notes: "v1.0 blah"},

		// actual data
		"1F61F ;	emoji ;	L1 ;	secondary ;	x	# V6.1 (😟) WORRIED FACE": Line{
			Single:     0x1f61f,
			Version:    6.1,
			Notes:      "WORRIED FACE",
			Properties: []string{"emoji", "L1", "secondary", "x"},
		},
		"2194..2199    ; Emoji                #   [6] (↔️..↙️)  LEFT RIGHT ARROW..SOUTH WEST ARROW": Line{
			Low:        0x2194,
			High:       0x2199,
			Notes:      "LEFT RIGHT ARROW..SOUTH WEST ARROW",
			Properties: []string{"Emoji"},
		},
		"002A FE0F 20E3; Emoji_Combining_Sequence  ; keycap: *                                                      # 3.0  [1] (*️⃣)": Line{
			Sequence:   []rune{0x002a, 0xfe0f, 0x20e3},
			Version:    3.0,
			Properties: []string{"Emoji_Combining_Sequence", "keycap: *"},
		},
		"0023 FE0E  ; text style;  # (1.1) NUMBER SIGN": Line{
			Sequence:   []rune{0x0023, 0xfe0e},
			Version:    1.1,
			Notes:      "NUMBER SIGN",
			Properties: []string{"text style", ""},
		},
		"1F649                                      ; fully-qualified     # 🙉 hear-no-evil monkey": Line{
			Single:     0x1f649,
			Notes:      "hear-no-evil monkey",
			Properties: []string{"fully-qualified"},
		},
		"1F442..1F4F7  ; Emoji_Presentation   #  6.0[182] (👂..📷)    ear..camera": Line{
			Low:        0x1f442,
			High:       0x1f4f7,
			Version:    6.0,
			Notes:      "ear..camera",
			Properties: []string{"Emoji_Presentation"},
		},
		"0023 FE0F 20E3; Emoji_Combining_Sequence  ; keycap: #                                                      # 3.0  [1] (#️⃣)": Line{
			Sequence:   []rune{0x0023, 0xfe0f, 0x20e3},
			Version:    3.0,
			Properties: []string{"Emoji_Combining_Sequence", "keycap: #"},
		},
		`0023 FE0F 20E3; Emoji_Keycap_Sequence     ; keycap: \x{23}                                                 #  3.0  [1] (#️⃣)`: Line{
			Sequence:   []rune{0x0023, 0xfe0f, 0x20e3},
			Version:    3.0,
			Properties: []string{"Emoji_Keycap_Sequence", `keycap: #`},
		},
		`0031 FE0F 20E3                             ; fully-qualified     # 1️⃣ keycap: 1`: Line{
			Sequence:   []rune{0x0031, 0xfe0f, 0x20e3},
			Notes:      "keycap: 1",
			Properties: []string{"fully-qualified"},
		},
		`1F18E                                      ; fully-qualified     # 🆎 AB button (blood type)`: Line{
			Single:     0x1f18e,
			Notes:      "AB button (blood type)",
			Properties: []string{"fully-qualified"},
		},
		`1F62C         ; Extended_Pictographic#  6.1  [1] (😬)       grimacing face`: Line{
			Single:     0x1f62c,
			Version:    6.1,
			Notes:      "grimacing face",
			Properties: []string{"Extended_Pictographic"},
		},
		`1F93F         ; Extended_Pictographic#   NA  [1] (🤿️)       <reserved-1F93F>`: Line{
			Single:     0x1f93f,
			Notes:      "<reserved-1F93F>",
			Properties: []string{"Extended_Pictographic"},
		},
	}

	for input, expected := range testdata {
		actual, err := Parse([]byte(input))
		if err != nil {
			t.Errorf("got err: %v", err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %+v, was %+v", expected, actual)
		}
	}
}

func TestHasEmoji(t *testing.T) {
	yes := "1F61F ;	emoji ;	L1 ;	secondary ;	x	# V6.1 (😟) WORRIED FACE"
	if actual, err := Parse([]byte(yes)); err != nil || !actual.HasEmoji() {
		t.Errorf("had err or no emoji: %v", err)
	}
	no := "# 6.0 just a comment"
	if actual, err := Parse([]byte(no)); err != nil || actual.HasEmoji() {
		t.Errorf("comment had err or had emoji: %v", err)
	}
}
