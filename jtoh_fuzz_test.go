//go:build go1.18
// +build go1.18

package jtoh_test

import (
	"bytes"
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
