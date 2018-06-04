package tr51

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// Reader reads an io.Reader and passes each Line to the specified method.
func Reader(r io.Reader, fn func(Line) error) error {
	rd := bufio.NewReader(r)

	for {
		raw, _, err := rd.ReadLine()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		raw = bytes.TrimSpace(raw)
		if len(raw) == 0 {
			continue // blank
		}

		out, err := Parse(raw)
		if err != nil {
			return err
		}
		err = fn(out)
		if err != nil {
			return err
		}
	}
	return nil
}

// File reads a file and passes each Line to the specified method.
func File(p string, fn func(Line) error) error {
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()
	return Reader(file, fn)
}
