//go:build go1.18

package jtoh_test

import (
	"bytes"
	"testing"

	"github.com/madlambda/jtoh"
)

func FuzzJTOH(f *testing.F) {
	seedCorpus := [][]byte{
		[]byte("string"),
		[]byte(" "),
		[]byte("{}"),
		[]byte("[]"),
		[]byte(`{ "name": "value"}`),
		[]byte(`{ "name": 666}`),
		[]byte(`{ "name": true}`),
		[]byte(`[{ "name": "value"}]`),
		[]byte(`[{ "name": 666}]`),
		[]byte(`[{ "name": true}]`),
		[]byte(`[{ "name": "value"}, {"name":666}]`),
		[]byte(`{ "name": "value"}
			{"name":666}`),
		[]byte("\nmsg\nmsg2\nmsg2\n"),
	}

	for _, seed := range seedCorpus {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, orig []byte) {
		input := bytes.NewReader(orig)
		output := &bytes.Buffer{}

		j, err := jtoh.New(":selector")
		if err != nil {
			t.Fatal(f)
		}

		j.Do(input, output)
	})
}
