package emoji

import (
	"io"

	"github.com/samthor/tr51"
)

// Images wraps parsed data from a custom images format.
type Images struct {
	emoji   map[string][]string
	vendors []string
}

// NewImages returns a new Images struct.
func NewImages(r *tr51.Reader) (*Images, error) {
	im := &Images{
		emoji: make(map[string][]string),
	}

	var currentVendor string
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !l.HasEmoji() {
			currentVendor = l.Notes
			im.vendors = append(im.vendors, currentVendor)
			continue
		}

		seq := string(l.AsSequence())
		if len(seq) == 0 {
			continue
		}
		im.emoji[seq] = append(im.emoji[seq], currentVendor)
	}

	return im, nil
}

// Vendors returns the list of emoji image vendors.
func (im *Images) Vendors() []string {
	// TODO: splice
	return im.vendors
}

// Get returns the list of codes that this emoji is supported in.
func (im *Images) Get(emoji string) []string {
	key := tr51.Unqualify(emoji)
	return im.emoji[key]
}
