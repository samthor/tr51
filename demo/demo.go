// Package main will print all emoji from the latest revision, except unqualified emoji.
package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/samthor/tr51"
)

const (
	testDataURL       = "https://unicode.org/Public/emoji/latest/emoji-test.txt"
	qualifiedProperty = "fully-qualified"
)

func main() {
	resp, err := http.Get(testDataURL)
	if err != nil {
		log.Fatal(err)
	}
	all, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	r := tr51.NewReader(bytes.NewBuffer(all))
	for {
		out, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if out.HasEmoji() && out.HasProperty(qualifiedProperty) {
			seq := out.AsSequence()
			log.Printf("%v\t%v (%v)", string(seq), out.Notes, seq)
		}
	}
}
