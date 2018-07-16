package emoji

import (
	"io"

	"github.com/samthor/tr51"
)

// Cats wraps parsed data from a custom TR51 format.
type Cats struct {
	emoji   map[string][]string
	vendors []string
}

// NewCats returns a new Cats struct.
func NewCats(r *tr51.Reader) (*Cats, error) {
	im := &Cats{
		emoji: make(map[string][]string),
	}

	var activeTitle string
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !l.HasEmoji() && l.Notes != "" {
			activeTitle = l.Notes
			im.vendors = append(im.vendors, activeTitle)
			continue
		}

		l.Each(func(raw []rune) {
			seq := string(raw)
			im.emoji[seq] = append(im.emoji[seq], activeTitle)
		})
	}

	return im, nil
}

// Titles returns the list of category titles
func (im *Cats) Titles() []string {
	out := make([]string, len(im.vendors))
	copy(out, im.vendors)
	return out
}

// Get returns the list of categories that this emoji is in.
func (im *Cats) Get(emoji string) []string {
	key := tr51.Unqualify(emoji)
	return im.emoji[key]
}
