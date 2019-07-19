package tr51

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	raw := []byte(`
# comment ignored

# ^ empty line above ignored
1F60E                                      ; fully-qualified     # 😎 smiling face with sunglasses
1F60D                                      ; fully-qualified     # 😍 smiling face with heart-eyes
  `)

	var count, emoji int
	reader := bytes.NewBuffer(raw)
	err := ReadFunc(reader, func(l Line) error {
		if l.HasEmoji() {
			emoji++
		}
		count++
		return nil
	})
	if err != nil {
		t.Errorf("got err, expected nil: %v", err)
	}
	if expected := 2; emoji != expected {
		t.Errorf("expected %d, was %d emoji", expected, emoji)
	}
	if expected := 4; count != expected {
		t.Errorf("expected %d, was %d lines", expected, count)
	}
}
