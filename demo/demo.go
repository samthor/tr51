// Package main will print all emoji from the latest revision, except unqualified emoji.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/samthor/tr51"
	"github.com/samthor/tr51/emoji"
)

const (
	dataDataURL       = "https://unicode.org/Public/emoji/latest/emoji-data.txt"
	testDataURL       = "https://unicode.org/Public/emoji/latest/emoji-test.txt"
	qualifiedProperty = "fully-qualified"
)

func readURL(url string) *bytes.Buffer {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	all, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return bytes.NewBuffer(all)
}

func main() {
	data, err := emoji.NewData(tr51.NewReader(readURL(dataDataURL)))
	if err != nil {
		log.Fatal(err)
	}

	payload := make(map[string]string)
	r := tr51.NewReader(readURL(testDataURL))
	for {
		out, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if !out.HasEmoji() || !out.HasProperty(qualifiedProperty) {
			continue
		}

		s := string(out.AsSequence())
		norm := data.Normalize(s)
		if _, ok := payload[norm]; ok || len(norm) == 0 {
			continue // dup re: gender or tone
		}
		payload[norm] = out.Notes
	}

	out, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got JSON payload: %d bytes (%d entries)", len(out), len(payload))
	fmt.Printf("%v", string(out))
}
