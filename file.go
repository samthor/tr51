package tr51

import (
	"bufio"
	"bytes"
	"io"
)

// ReadFunc reads an io.Reader and passes each Line to the specified method.
func ReadFunc(r io.Reader, fn func(Line) error) error {
	er := NewReader(r)
	for {
		line, err := er.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		err = fn(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reader allows reading of TR51 data.
type Reader struct {
	r *bufio.Reader
}

// NewReader returns a new Reader for TR51 data.
func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r)}
}

// Read returns the next line of TR51 data. If no line is available, returns an error.
func (r *Reader) Read() (Line, error) {
	var empty Line

	for {
		raw, _, err := r.r.ReadLine()
		if err != nil {
			return empty, err // including io.EOF
		}
		raw = bytes.TrimSpace(raw)
		if len(raw) > 0 {
			return Parse(raw)
		}
	}

	return empty, io.EOF
}
