//go:build go1.18
// +build go1.18

package jtoh_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/madlambda/jtoh"
)

func FuzzJTOH(f *testing.F) {
	type seed struct {
		selector string
		input    []byte
	}
	seedCorpus := []seed{
		{
			selector: ":a",
			input:    []byte("string"),
		},
		{
			selector: "/s",
			input:    []byte(" "),
		},
		{
			selector: "aa",
			input:    []byte("{}"),
		},
		{
			selector: ".field4",
			input:    []byte("[]"),
		},
		{
			selector: "%1",
			input:    []byte(`{ "name": "value"}`),
		},
		{
			selector: "<field",
			input:    []byte(`{ "name": 666}`),
		},
		{
			selector: "?xS",
			input:    []byte(`{ "name": true}`),
		},
		{
			selector: "|da",
			input:    []byte(`[{ "name": "value"}]`),
		},
		{
			selector: "^field^field2",
			input:    []byte(`[{ "name": 666}]`),
		},
		{
			selector: "#a#name",
			input:    []byte(`[{ "name": true}]`),
		},
		{
			selector: ":a:B:c",
			input:    []byte(`[{ "name": "value"}, {"name":666}]`),
		},
		{
			selector: ":name",
			input: []byte(`{ "name": "value"}
			{"name":666}`),
		},
		{
			selector: "@a@b",
			input:    []byte("\nmsg\nmsg2\nmsg2\n"),
		},
	}

	for _, seed := range seedCorpus {
		f.Add(seed.selector, seed.input)
	}

	f.Fuzz(func(t *testing.T, selector string, orig []byte) {
		input := bytes.NewReader(orig)
		output := &bytes.Buffer{}

		j, err := jtoh.New(selector)
		if err != nil {
			// It is expected that invalid selector will be generated
			return
		}

		j.Do(input, output)
	})
}

func FuzzJTOHValid(f *testing.F) {
	type seed struct {
		key string
		val string
	}

	seedCorpus := []seed{
		{
			key: "str",
			val: "str",
		},
	}

	for _, seed := range seedCorpus {
		f.Add(seed.key, seed.val)
	}

	f.Fuzz(func(t *testing.T, key string, val string) {
		// Since field selectors have spaces trimmed for now we also
		// trim the key, or else we would not be able to select the key.
		key = strings.TrimSpace(key)
		if key == "" {
			return
		}
		if strings.ContainsAny(key, ".:") {
			// We don't handle nesting/keys with dot on name for now.
			return
		}

		input, err := json.Marshal(map[string]string{key: val})
		if err != nil {
			return
		}

		// key/value may change on marshalling, so we get the actual
		// final key/value from the json encoding/decoding process.
		parsedInput := map[string]string{}
		if err = json.Unmarshal(input, &parsedInput); err != nil {
			t.Fatal(err)
		}

		var selectKey, wantValue string

		for k, v := range parsedInput {
			selectKey = k
			wantValue = v
		}

		selector := ":" + selectKey

		j, err := jtoh.New(selector)
		if err != nil {
			t.Fatal(err)
		}

		// Newlines on values are escaped to avoid breaking the
		// line oriented nature of the output.
		wantValue = strings.Replace(wantValue, "\n", "\\n", -1)
		want := wantValue + "\n"

		testSelection := func(input []byte) {
			output := &bytes.Buffer{}
			j.Do(bytes.NewReader(input), output)

			got := output.String()
			if got != want {
				t.Errorf("input: %q", string(input))
				t.Errorf("str  : got %q != want %q", got, want)
				t.Errorf("bytes: got %v != want %v", []byte(got), []byte(want))
			}
		}

		// Selecting on single document/stream must behave identically to
		// selecting from a list of documents.
		testSelection(input)
		testSelection([]byte("[" + string(input) + "]"))
	})
}
